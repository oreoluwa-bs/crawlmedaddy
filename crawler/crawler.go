package crawler

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/oreoluwa-bs/crawlmedaddy/shared"
)

type Store interface {
	Store(ctx context.Context, cd *shared.CrawledData) error
	Retrieve(ctx context.Context, url string) ([]*shared.CrawledData, error)
	RetrieveCrawlable(ctx context.Context, url string) ([]*shared.CrawledData, error)
}

type Crawler struct {
	store Store
}

func New(store Store) *Crawler {

	c := &Crawler{

		store,
	}

	return c
}

// Extracts the links and metadata from the html on a given url
func (c *Crawler) Retriever(url string) *shared.CrawledData {
	cd := &shared.CrawledData{
		Url:  url,
		Meta: make(map[string]string),
	}

	clly := colly.NewCollector()

	// clly.WithTransport(&http.Transport{
	// 	IdleConnTimeout:     90 * time.Second,
	// 	ResponseHeaderTimeout: 30 * time.Second,
	// })

	clly.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")

		// if isRelativeUrl(href) {
		// href = url + href
		// }
		cd.Urls = append(cd.Urls, e.Request.AbsoluteURL(href))
	})

	clly.OnHTML("meta[name]", func(e *colly.HTMLElement) {
		cd.Meta[e.Attr("name")] = e.Attr("content")
	})

	clly.OnHTML("h1, h2, h3, p", func(e *colly.HTMLElement) {
		cd.Content += (e.Text + " ")
	})
	clly.OnHTML("title", func(e *colly.HTMLElement) {
		cd.Title += e.Text
	})

	err := clly.Visit(url)

	if err != nil {
		log.Printf("Error visiting URL %s: %v", url, err)
		return nil
	}

	if cd.Title == "" {
		if metaTitle, ok := cd.Meta["title"]; ok {
			cd.Title = metaTitle
		}
	}
	return cd
}

func (c *Crawler) Crawl(urls []string) {
	const delay = 500 * time.Millisecond
	limiter := time.Tick(delay)
	for _, u := range urls {
		<-limiter
		url := removeTrailingSlash(u)

		cd := c.Retriever(url)
		log.Printf("Retrieved data for URL %s", url)
		// Check if crawled
		visited := c.visited(context.Background(), url)
		if visited {
			continue
		}

		// Rank
		//

		cd.DateCrawled = time.Now()
		cd.NextCrawlDate = time.Now().Add(240 * time.Hour) // 10 days

		// Store crawled
		if err := c.store.Store(context.Background(), cd); err != nil {
			log.Printf("Error storing data for URL %s: %v", url, err)
			continue
		}
		log.Printf("Stored data for URL %s", url)

		// Repeat
		// go func(ctx context.Context, u []string) {
		c.Crawl(cd.Urls)
		// }(ctx, cd.Urls)
	}
}

func (c *Crawler) visited(ctx context.Context, url string) bool {
	existingUrls, err := c.store.RetrieveCrawlable(ctx, url)
	if err != nil {
		log.Panicln(err)

	}

	return len(existingUrls) > 0
}

func isRelativeUrl(url string) bool {
	return strings.HasPrefix(url, "/")
}

func isHashUrl(url string) bool {
	return strings.HasPrefix(url, "#")
}

func removeTrailingSlash(url string) string {
	if strings.HasSuffix(url, "/") {
		return removeTrailingSlash(strings.TrimSuffix(url, "/"))
	}

	return url
}

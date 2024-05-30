package crawler

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/oreoluwa-bs/crawlmedaddy/shared"
)

type CrawlerStore struct {
	db *sql.DB
}

func NewCrawlerStore(connectionString string) *CrawlerStore {

	db, err := sql.Open("sqlite3", connectionString)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	db.Exec(`
	CREATE TABLE IF NOT EXISTS pages (
    	id INTEGER NOT NULL PRIMARY KEY,
    	url TEXT NOT NULL,
    	title TEXT NOT NULL,
     	meta JSON,
      	urls JSON,
    	content LONGTEXT,
     	date_crawled DATETIME DEFAULT CURRENT_TIMESTAMP,
     	next_crawl_date DATETIME NOT NULL
    )
		`)

	cs := &CrawlerStore{
		db,
	}

	return cs
}

func (s *CrawlerStore) Store(ctx context.Context, cd *shared.CrawledData) error {

	jsonurl, err := json.Marshal(cd.Urls)
	if err != nil {
		return err
	}
	jsonmeta, err := json.Marshal(cd.Meta)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, `
	INSERT INTO pages (
		url, title, meta, urls, content, date_crawled, next_crawl_date
	)
	VALUES (?, ?, ?, ?, ?, ?, ?);
	`, cd.Url, cd.Title, string(jsonmeta), string(jsonurl), cd.Content, cd.DateCrawled, cd.NextCrawlDate)

	return err
}

func (s *CrawlerStore) Retrieve(ctx context.Context, url string) ([]*shared.CrawledData, error) {

	rows, err := s.db.QueryContext(ctx, `
	SELECT url, title, date_crawled, next_crawl_date FROM pages
	WHERE url = ?;
	`, url, time.Now())

	if err != nil {
		return nil, err
	}

	cd := make([]*shared.CrawledData, 0)

	for rows.Next() {
		d := shared.CrawledData{}

		rows.Scan(
			d.Url,
			d.Title,
			d.DateCrawled,
			d.NextCrawlDate,
		)

		cd = append(cd, &d)
	}

	defer rows.Close()

	return cd, err

}

func (s *CrawlerStore) RetrieveCrawlable(ctx context.Context, url string) ([]*shared.CrawledData, error) {

	rows, err := s.db.QueryContext(ctx, `
	SELECT url, title, date_crawled, next_crawl_date FROM pages
	WHERE url = ? AND datetime(date_crawled, '+10 days') > ?;
	`, url, time.Now())

	if err != nil {
		return nil, err
	}

	cd := make([]*shared.CrawledData, 0)

	for rows.Next() {
		d := shared.CrawledData{}

		rows.Scan(
			d.Url,
			d.Title,
			d.DateCrawled,
			d.NextCrawlDate,
		)

		cd = append(cd, &d)
	}

	defer rows.Close()

	return cd, err

}

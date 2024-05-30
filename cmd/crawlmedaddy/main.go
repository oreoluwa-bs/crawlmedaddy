package main

import (
	"github.com/oreoluwa-bs/crawlmedaddy/crawler"
)

func main() {

	cs := crawler.NewCrawlerStore("../data/database.db")
	c := crawler.New(cs)

	input := []string{"https://oreoluwabs.com"}


	go func(s []string) {
		c.Crawl(s)
	}(input)

	//
	// Schedule follow up
	//
	//
}

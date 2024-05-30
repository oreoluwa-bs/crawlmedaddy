package shared

import "time"

type CrawledData struct {
	Url           string
	Title         string
	Meta          map[string]string
	Urls          []string
	Content       string
	DateCrawled   time.Time
	NextCrawlDate time.Time
}

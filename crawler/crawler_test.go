package crawler

import (
	"context"
	"testing"
	"time"
)

// Should return a list of urls and metadata from a
func TestRetriever(t *testing.T) {
	cs := NewCrawlerStore("../data/database-test.db")
	c := New(cs)

	input := "https://www.matthiax.com/"

	cd := c.Retriever(input)
	//d

	if len(cd.Meta) < 1 {
		t.Fatalf("crawling %s - expected to find up to %d metatags(s), got %d", input, 1, len(cd.Meta))
	}

	if len(cd.Urls) < 1 {
		t.Fatalf("crawling %s - expected to find up to %d link(s), got %d", input, 1, len(cd.Urls))
	}

}

func TestCrawl(t *testing.T) {
	cs := NewCrawlerStore("../data/database-test.db")
	c := New(cs)

	input := []string{"https://www.matthiax.com/"}

	// go func(s []string) {
	// 	c.Crawl(s)
	// }(input)
	// c.Crawl(input)

	// time.Sleep(10 * time.Second)
	//
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		c.Crawl(input)
	}()

	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			t.Log("Crawl function stopped due to timeout")
		} else {
			t.Fatal(ctx.Err())
		}
	case <-time.After(12 * time.Second):
		t.Fatal("Crawl function did not stop in time")
	}

	data, err := c.store.Retrieve(context.Background(), input[0])
	if err != nil {
		t.Fatal(err)
	}

	if len(data) < 0 {
		t.Fatalf("expects response from the retriever to be greater than %d, got %d", 1, len(data))
	}

}

func TestRemoveTrailingSlash(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/path/", "/path"},
		{"/path//", "/path"},
		{"path/", "path"},
		{"/path/to/resource/", "/path/to/resource"},
		{"/path/to/resource//", "/path/to/resource"},
	}

	for _, test := range tests {
		result := removeTrailingSlash(test.input)
		if result != test.expected {
			t.Fatalf("removeTrailingSlash(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}

package main

import (
	"fmt"
	"time"
)

//GoString returns the string representation of crawler
func (c *Crawler) GoString() string {
	return fmt.Sprintf("Crawler{urlSource: %d, fetcher:%#v}", len(c.urlSource), c.fetcher)
}

//Crawler fetches urls using an assigned fetcher
type Crawler struct {
	fetcher   Fetcher
	urlSource chan crawlable
}

//NewCrawler returns a crawler based on the given fetcher
func NewCrawler(fetcher Fetcher) *Crawler {
	return &Crawler{
		fetcher:   fetcher,
		urlSource: make(chan crawlable),
	}
}

type crawlable struct {
	url   string
	depth int
}

func (c *Crawler) crawl(crawlData *crawlable) {
	println(".")
	if crawlData.depth <= 0 {
		return
	}

	body, urls, err := c.fetcher.Fetch(crawlData.url)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("found: %s %q\n", crawlData.url, body)
	go func() {
		for _, u := range urls {
			c.urlSource <- crawlable{u, crawlData.depth - 1}
		}
	}()

	return
}

func (c *Crawler) run() {
	shouldRun := true
	for shouldRun {
		select {
		case cData := <-c.urlSource:
			if cData.depth == -1 {
				shouldRun = false
				break
			}
			c.crawl(&cData)
		default:
			time.Sleep(5)
		}
	}
	hasRemaining := true
	for hasRemaining {
		_, hasRemaining = <-c.urlSource
	}
	close(c.urlSource)
}

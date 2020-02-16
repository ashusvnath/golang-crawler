package main

import (
	"bytes"
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
	output    *bytes.Buffer
	urlSource chan crawlable
}

//NewCrawler returns a crawler based on the given fetcher
func NewCrawler(fetcher Fetcher, output *bytes.Buffer) *Crawler {
	return &Crawler{
		fetcher:   fetcher,
		output:    output,
		urlSource: make(chan crawlable, 1),
	}
}

type crawlable struct {
	url   string
	depth int
}

func (c *Crawler) visit(crawlData *crawlable) {
	infof("crawling %v at depth %v", crawlData.url, crawlData.depth)
	if crawlData.depth < 0 {
		return
	}

	body, urls, err := c.fetcher.Fetch(crawlData.url)
	if err != nil {
		c.output.WriteString(err.Error() + "\n")
		return
	}

	c.output.WriteString(fmt.Sprintf("found at %s: %q\n", crawlData.url, body))
	go func() {
		for _, u := range urls {
			c.urlSource <- crawlable{u, crawlData.depth - 1}
		}
	}()

	return
}

//Crawl given url to specified depth
func (c *Crawler) Crawl(url string, depth int) {
	c.urlSource <- crawlable{url, depth}
	c.crawl()
}

func (c *Crawler) crawl() {
	shouldRun := true
	for shouldRun {
		select {
		case cData := <-c.urlSource:
			if cData.depth == 0 {
				info("stopping crawl")
				shouldRun = false
				break
			}
			c.visit(&cData)
		default:
			tracef("crawler: %#v", c)
			time.Sleep(1)
		}
	}
}

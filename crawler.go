package main

import (
	"bytes"
	"flag"
	"fmt"
)

var crawlQueueSize = 1

func init() {
	flag.IntVar(&crawlQueueSize, "qlen", 1, "Size of crawl queue")
}

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
		urlSource: make(chan crawlable, crawlQueueSize),
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
	done := make(chan int)
	go c.crawl(done)
	<-done
}

func (c *Crawler) crawl(done chan int) {
	shouldRun := true
	for shouldRun {
		select {
		case cData := <-c.urlSource:
			if cData.depth == 0 {
				info("stopping crawl")
				shouldRun = false
				break
			}
			go c.visit(&cData)
		}
	}
	debugf("crawler: %#v", c)
	done <- 1
}

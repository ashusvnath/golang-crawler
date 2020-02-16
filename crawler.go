package main

import (
	"flag"
	"fmt"
)

var crawlQueueSize = 1
var visitAsync = false

func init() {
	flag.IntVar(&crawlQueueSize, "blen", 1, "Size of crawlable buffer")
	flag.BoolVar(&visitAsync, "aVis", false, "Visit urls asynchronously")
}

//GoString returns the string representation of crawler
func (c *Crawler) GoString() string {
	return fmt.Sprintf("Crawler{urlSource: %d, fetcher:%#v}", len(c.urlSource), c.fetcher)
}

//Crawler fetches urls using an assigned fetcher
type Crawler struct {
	fetcher   Fetcher
	output    chan string
	urlSource chan crawlable
}

//NewCrawler returns a crawler based on the given fetcher
func NewCrawler(fetcher Fetcher, output chan string) *Crawler {
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

	body, urls, err := c.fetcher.Fetch(crawlData.url)
	if err != nil {
		c.output <- err.Error()
		return
	}

	c.output <- fmt.Sprintf("found at %s: %q\n", crawlData.url, body)
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
			if visitAsync {
				go c.visit(&cData)
			} else {
				c.visit(&cData)
			}

		}
	}
	debugf("crawler: %#v", c)
	done <- 1
}

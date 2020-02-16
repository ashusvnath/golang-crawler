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
		urlSource: make(chan crawlable, 1),
	}
}

type crawlable struct {
	url   string
	depth int
}

func (c *Crawler) visit(crawlData *crawlable) {
	println(".", crawlData.depth)
	if crawlData.depth < 0 {
		c.urlSource <- crawlable{"does not matter", -1}
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

//Crawl given url to specified depth
func (c *Crawler) Crawl(url string, depth int) {
	c.urlSource <- crawlable{url, depth}
	c.do()
}

func (c *Crawler) do() {
	shouldRun := true
	for shouldRun {
		select {
		case cData := <-c.urlSource:
			if cData.depth == -1 {
				fmt.Println("Stopping crawl")
				shouldRun = false
				break
			}
			c.visit(&cData)
		default:
			fmt.Printf("crawler: %#v\n", c)
			time.Sleep(5)
		}
	}
}

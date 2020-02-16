package main

import (
	"bytes"
	"fmt"
)

//Fetcher defines an interface for any object that fetches a web url
type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

func main() {
	Crawl("https://golang.org/", 4, NewCachingFetcher(NewFetcher()))
}

//Crawl all urls upto a depth using a fetcher
func Crawl(url string, depth int, fetcher Fetcher) {
	crawlerOutput := bytes.NewBuffer([]byte{})
	crawler := NewCrawler(fetcher, crawlerOutput)

	info("starting crawl")

	crawler.Crawl(url, depth)
	fmt.Printf("Crawl output: \n%v\n", crawlerOutput.String())
}

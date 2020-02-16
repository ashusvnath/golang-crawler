package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"time"
)

//Fetcher defines an interface for any object that fetches a web url
type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

func init() {
	rand.Seed(time.Now().Unix())
	flag.Parse()
}

func main() {
	Crawl("https://golang.org/", 4, NewCachingFetcher(NewFetcher()))
}

//Crawl all urls upto a depth using a fetcher
func Crawl(url string, depth int, fetcher Fetcher) {
	crawlerOutput := new(bytes.Buffer)
	crawler := NewCrawler(fetcher, crawlerOutput)

	crawler.Crawl(url, depth)
	fmt.Printf("Crawl output: \n%v\n", crawlerOutput.String())
}

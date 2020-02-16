package main

import "fmt"

//Fetcher defines an interface for any object that fetches a web url
type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

func main() {
	Crawl("https://golang.org/", 4, NewCachingFetcher(NewFetcher()))
}

//Crawl all urls upto a depth using a fetcher
func Crawl(url string, depth int, fetcher Fetcher) {
	crawler := NewCrawler(fetcher)
	fmt.Printf("crawler: %#v\n", crawler)
	crawler.Crawl(url, depth)
}

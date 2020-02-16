package main

import (
	"fmt"
	"sync"
)

type result struct {
	fakeResult
}

//NewCachingFetcher returns a caching fetcher based on the specified fetcher
func NewCachingFetcher(fetcher Fetcher) Fetcher {
	return &CachingFetcher{
		fetcher:     fetcher,
		mutex:       &sync.Mutex{},
		visitedUrls: make(map[string]*result),
	}
}

//CachingFetcher caches data fetched by a fetcher
type CachingFetcher struct {
	fetcher     Fetcher
	mutex       *sync.Mutex
	visitedUrls map[string]*result
}

//Fetch caches successfully fetched data from the given url
func (cf *CachingFetcher) Fetch(url string) (body string, urls []string, err error) {
	cf.mutex.Lock()
	defer cf.mutex.Unlock()
	if r, ok := cf.visitedUrls[url]; ok {
		return r.body, r.urls, nil
	}

	body, urls, err = cf.fetcher.Fetch(url)

	if err != nil {
		cf.visitedUrls[url] = &result{
			fakeResult{body,
				urls},
		}
	}
	return
}

//GoString returns a representation of CachingFetcher
func (cf *CachingFetcher) GoString() string {
	return fmt.Sprintf("CachingFetcher{visitedUrls: %v, fetcher: %#v}", cf.visitedUrls, cf.fetcher)
}

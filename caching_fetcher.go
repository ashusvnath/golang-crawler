package main

import "sync"

type result struct {
	fakeResult
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

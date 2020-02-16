package main

import (
	"fmt"
	"math/rand"
	"time"
)

type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f *fakeFetcher) Fetch(url string) (string, []string, error) {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	if res, ok := (*f)[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

//GoString returns a representation of fakeFetcher
func (f *fakeFetcher) GoString() string {
	return fmt.Sprintf("fakeFetcher{%v}", f)
}

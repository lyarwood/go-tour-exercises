package main

import (
	"fmt"
	"sort"
	"sync"
)

type FetcherCache struct {

	// Map of found URLs and the associated body
	f  map[string]string
	fm *sync.Mutex

	// Map of URLs fetched already
	c  map[string]bool
	cm *sync.Mutex

	// WaitGroup to sync go routines
	w *sync.WaitGroup
}

func (fc FetcherCache) HasFetch(url string) bool {
	fc.cm.Lock()
	defer fc.cm.Unlock()
	_, ok := fc.c[url]
	return ok
}

func (fc *FetcherCache) AddFetch(url string) {
	fc.cm.Lock()
	defer fc.cm.Unlock()
	fc.c[url] = true
}

func (fc *FetcherCache) AddFound(url string, body string) {
	fc.fm.Lock()
	defer fc.fm.Unlock()
	fc.f[url] = body
}

func (fc FetcherCache) Sort() map[string]string {
	sorted := map[string]string{}
	keys := make([]string, 0, len(fc.f))
	for k := range fc.f {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		sorted[k] = fc.f[k]
	}
	return sorted
}

func (fc FetcherCache) String() string {
	s := ""
	for k, v := range fc.Sort() {
		s = s + fmt.Sprintf("%s %s\n", k, v)
	}
	return s
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fc *FetcherCache) {
	defer fc.w.Done()

	if depth <= 0 || fc.HasFetch(url) {
		return
	}

	fc.AddFetch(url)
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		return
	}
	fc.AddFound(url, body)

	for _, u := range urls {
		fc.w.Add(1)
		go Crawl(u, depth-1, fc)
	}
}

func main() {

	fc := FetcherCache{
		f:  map[string]string{},
		fm: &sync.Mutex{},
		c:  map[string]bool{},
		cm: &sync.Mutex{},
		w:  &sync.WaitGroup{},
	}

	fc.w.Add(1)
	go Crawl("https://golang.org/", 4, &fc)
	fc.w.Wait()
	fmt.Printf("Found the following URLs: \n%+v", &fc)

}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	fmt.Printf("Fetching %s\n", url)
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}

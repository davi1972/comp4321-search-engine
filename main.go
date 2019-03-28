package main

import (
	"github.com/gocolly/colly"
	//"github.com/gocolly/colly/debug"
	"fmt"
	"net/http"
	"sync"
	"comp4321/concurrentMap"
	"comp4321/tokenizer"
)
var idCount int = 2
var pageMap map[string]int

type page struct {
	id int
	title string
	url string
	// size int
	content []string
	children []int
}

var requestID string

func main() {
	var wg = &sync.WaitGroup{}

	rootPage := "https://www.cse.ust.hk"
	// rootPage := "https://apartemen.win/comp4321/page1.html"

	tokenizer.LoadStopWords()

	pageMap := concurrentMap.ConcurrentMap{}
	pageMap.Init()
	pageMap.Set(rootPage, 1)

	pages := []page{}
	
	crawler := colly.NewCollector(		
		colly.MaxDepth(2),
		// colly.Debugger(&debug.LogDebugger{}),
		colly.Async(true),
	)
	
	// Limit the maximum parallelism to 2
	// This is necessary if the goroutines are dynamically
	// created to control the limit of simultaneous requests.
	//
	// Parallelism can be controlled also by spawning fixed
	// number of go routines.
	crawler.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	crawler.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
		fmt.Println("")
	})


	crawler.OnHTML("html", func(e *colly.HTMLElement) {

		temp := page{}
		
		temp.title = e.ChildText("title")
		temp.url = e.Request.URL.String()

		temp.id = pageMap.Get(temp.url).(int)
		temp.content = tokenizer.Tokenize(e.ChildText("body"))

		temp.children = []int{}
		links := e.ChildAttrs("a[href]", "href")

		for _, url := range links {
			url = e.Request.AbsoluteURL(url)
			wg.Add(1)
			go func (url string) {
				defer wg.Done()
				resp, err := http.Get(url)
				if err != nil || resp.StatusCode != 200 {
					return
				}
				defer resp.Body.Close()

				pageMap.Set(url, idCount)

				idCount++

				temp.children = append(temp.children, pageMap.Get(url).(int))

				e.Request.Visit(url)
			}(url)
		}
		wg.Wait()
		pages = append(pages, temp)

	})

	crawler.OnRequest(func(r *colly.Request){
		fmt.Println("Visiting", r.URL)
	})

	crawler.Visit(rootPage)

	crawler.Wait()
	fmt.Println("Pages: ")
	for _, page := range pages {
		fmt.Println("ID:", page.id)
		fmt.Println("Title:", page.title)
		fmt.Println("content:", page.content)
		fmt.Println("url:", page.url)
		fmt.Print("children:")
		for _, child := range page.children {
			fmt.Print(child, " ")
		}
		fmt.Print("\n")
	}

	fmt.Println("\nURL to ID maps:")

	pageMapChan := pageMap.Iter()
	for items := range pageMapChan {
		fmt.Println("ID:", items.Key)
		fmt.Println("URL:", items.Value.(int))
	}
}	
package main

import (
	"github.com/gocolly/colly"
	//"github.com/gocolly/colly/debug"
	"fmt"
	"net/http"
)
var idCount int = 2
var pageMap map[string]int

type page struct {
	id int
	title string
	url string
	// size int
	content string
	children []int
}

var requestID string

func main() {

	pageMap = make(map[string]int)
	pageMap["http://apartemen.win/comp4321/page1.html"] = 1

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

		temp.id = pageMap[temp.url]
		// temp.content = e.ChildText("body")

		temp.children = []int{}
		links := e.ChildAttrs("a[href]", "href")

		for _, url := range links {
			url = e.Request.AbsoluteURL(url)
			resp, err := http.Get(url)
			if err != nil || resp.StatusCode !=200 {
				continue
			}
			defer resp.Body.Close()

			pageMap[url] = idCount
			idCount++

			temp.children = append(temp.children, pageMap[url])

			e.Request.Visit(url)
		}

		pages = append(pages, temp)

	})
	// On every a element which has href attribute call callback
	// crawler.OnHTML("a[href]", func(e *colly.HTMLElement) {
	// 	link := e.Attr("href")
	// 	// Print link
	// 	fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		
	// 	// Visit link found on page
	// 	// Only those links are visited which are in AllowedDomains
	// 	e.Request.Visit(link)

	// })

	crawler.OnRequest(func(r *colly.Request){
		fmt.Println("Visiting", r.URL)
	})

	// crawler.OnHTML("*", func(el *colly.HTMLElement) {
	// 	fmt.Println("Parent Result", el.Text)
	// })

	crawler.Visit("http://apartemen.win/comp4321/page1.html")

	crawler.Wait()
	fmt.Println("Pages: \n")
	for _, page := range pages {
		fmt.Println("ID:", page.id)
		fmt.Println("Title:", page.title)
		fmt.Println("url:", page.url)
		fmt.Print("children:")
		for _, child := range page.children {
			fmt.Print(child, " ")
		}
		fmt.Println("\n")
	}

	fmt.Println("\nURL to ID maps:\n")

	for url, id := range pageMap {
		fmt.Println("ID:", id)
		fmt.Println("URL:", url, "\n")

	}
}	
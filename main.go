package main

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"fmt"
)

func main() {
	crawler := colly.NewCollector(		
		colly.MaxDepth(5),
		colly.Debugger(&debug.LogDebugger{}),
		colly.Async(true),
	)
	
	// Limit the maximum parallelism to 2
	// This is necessary if the goroutines are dynamically
	// created to control the limit of simultaneous requests.
	//
	// Parallelism can be controlled also by spawning fixed
	// number of go routines.
	crawler.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	// On every a element which has href attribute call callback
	crawler.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		e.Request.Visit(link)
	})

	crawler.OnRequest(func(r *colly.Request){
		fmt.Println("Visiting", r.URL)
	})

	crawler.OnHTML("*", func(el *colly.HTMLElement) {
		fmt.Println("Parent Result", el.Text)
	})

	crawler.Visit("http://www.cse.ust.hk")

	crawler.Wait()
}	
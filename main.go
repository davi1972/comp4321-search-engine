package main

import (
	"github.com/gocolly/colly"
	//"github.com/gocolly/colly/debug"
	"comp4321/concurrentMap"
	Indexer "comp4321/indexer"
	"comp4321/tokenizer"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var idCount int = 2
var pageMap map[string]int

type page struct {
	id            int
	title         string
	url           string
	size          int
	content       []string
	children      []int
	date_modified time.Time
}

var requestID string

func main() {
	var wg = &sync.WaitGroup{}
	wd, _ := os.Getwd()

	// rootPage := "https://www.cse.ust.hk"
	rootPage := "https://apartemen.win/comp4321/page1.html"

	tokenizer.LoadStopWords()

	pageMap := concurrentMap.ConcurrentMap{}
	pageMap.Init()
	pageMap.Set(rootPage, 1)

	documentIndexer := &Indexer.MappingIndexer{}
	docErr := documentIndexer.Initialize(wd + "/db/documentIndex")
	if docErr != nil {
		fmt.Printf("error when initializing document indexer: %s\n", docErr)
	}
	defer documentIndexer.Backup()
	defer documentIndexer.Release()

	wordIndexer := &Indexer.MappingIndexer{}
	wordErr := wordIndexer.Initialize(wd + "/db/wordIndex")
	if wordErr != nil {
		fmt.Printf("error when initializing word indexer: %s\n", wordErr)
	}
	defer wordIndexer.Backup()
	defer wordIndexer.Release()

	pagePropertiesIndexer := &Indexer.PagePropetiesIndexer{}
	pagePropertiesErr := pagePropertiesIndexer.Initialize(wd + "/db/pagePropertiesIndex")
	if pagePropertiesErr != nil {
		fmt.Printf("error when initializing page properties: %s\n", pagePropertiesErr)
	}
	defer pagePropertiesIndexer.Backup()
	defer pagePropertiesIndexer.Release()

	contentInvertedIndexer := &Indexer.InvertedFileIndexer{}
	contentInvertedErr := contentInvertedIndexer.Initialize(wd + "/db/contentInvertedIndex")
	if contentInvertedErr != nil {
		fmt.Printf("error when initializing page properties: %s\n", contentInvertedErr)
	}
	defer contentInvertedIndexer.Backup()
	defer contentInvertedIndexer.Release()

	documentForwardIndexer := &Indexer.ForwardIndexer{}
	documentForwardIndexerErr := documentForwardIndexer.Initialize(wd + "/db/documentForwardIndex")
	if documentForwardIndexerErr != nil {
		fmt.Printf("error when initializing page properties: %s\n", contentInvertedErr)
	}
	defer documentForwardIndexer.Backup()
	defer documentForwardIndexer.Release()

	pages := []page{}

	crawler := colly.NewCollector(
		colly.MaxDepth(4),
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

		size, _ := strconv.Atoi(e.Response.Headers.Get("Content-Length"))

		if size == 0 {
			size = len(e.Text)
		}

		temp.size = size

		date := e.Response.Headers.Get("Last-Modified")

		if len(date) != 0 {
			temp.date_modified, _ = time.Parse(time.RFC1123, date)
		} else {
			temp.date_modified = time.Now()
		}

		// Store Document id and properties
		_, err := documentIndexer.AddKeyToIndex(temp.url)
		if err != nil {
			fmt.Println(err)
		}
		id, _ := documentIndexer.GetValueFromKey(temp.url)
		pagePropertiesIndexer.AddKeyToPageProperties(id, Indexer.CreatePage(id, temp.title, temp.url, size, temp.date_modified))

		// temp.id = pageMap.Get(temp.url).(int)
		temp.content = tokenizer.Tokenize(e.ChildText("body"))

		// Check for duplicate words in the document
		wordList := make(map[uint64]*Indexer.InvertedFile)
		for i, v := range temp.content {
			// Add Word to id index
			wordID, err := wordIndexer.GetValueFromKey(v)
			if err != nil {
				wordID, _ = wordIndexer.AddKeyToIndex(v)
			}

			invFile, contain := wordList[wordID]
			if contain {
				invFile.AddWordPositions(uint64(i))
			} else {
				wordList[wordID] = Indexer.CreateInvertedFile(id)
				wordList[wordID].AddWordPositions(uint64(i))
			}
		}

		for k, v := range wordList {
			contentInvertedIndexer.AddKeyToIndexOrUpdate(k, *v)
		}

		// Get List word words in this document in slice
		wordUIntList := make([]uint64, 0, len(wordList))
		for k := range wordList {
			wordUIntList = append(wordUIntList, k)
		}

		documentForwardIndexer.AddWordIdListToKey(id, wordUIntList)

		temp.children = []int{}

		// temp.date_modified = e.Response.Headers.Get("Last-Modified")
		links := e.ChildAttrs("a[href]", "href")

		for _, url := range links {
			url = e.Request.AbsoluteURL(url)
			wg.Add(1)
			go func(url string) {
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

	crawler.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	crawler.Visit(rootPage)

	crawler.Wait()
	fmt.Println("Pages: ")
	for _, page := range pages {
		fmt.Println("ID:", page.id)
		fmt.Println("url:", page.url)
		fmt.Println("Title:", page.title)
		fmt.Println("Size:", page.size)
		fmt.Println("Content:", page.content)

		fmt.Print("children:")
		for _, child := range page.children {
			fmt.Print(child, " ")
		}
		fmt.Print("\n")
		fmt.Println("date_modified:", page.date_modified.Format(time.RFC1123))
		fmt.Print("\n")
	}

	fmt.Println("\nURL to ID maps:")

	pageMapChan := pageMap.Iter()
	for items := range pageMapChan {
		fmt.Println("ID:", items.Key)
		fmt.Println("URL:", items.Value.(int))
	}

	contentInvertedIndexer.Iterate()
	wordIndexer.Iterate()
	documentIndexer.Iterate()
	pagePropertiesIndexer.Iterate()
	documentForwardIndexer.Iterate()
}

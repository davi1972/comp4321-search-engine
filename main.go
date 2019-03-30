package main

import (
	"github.com/gocolly/colly"
	//"github.com/gocolly/colly/debug"
	"github.com/hskrishandi/comp4321/concurrentMap"
	Indexer "github.com/hskrishandi/comp4321/indexer"
	"github.com/hskrishandi/comp4321/tokenizer"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
	"strings"
)

type page struct {
	id       uint64
	children []uint64
	parent   []uint64
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

	titleInvertedIndexer := &Indexer.InvertedFileIndexer{}
	titleInvertedErr := titleInvertedIndexer.Initialize(wd + "/db/titleInvertedIndex")
	if titleInvertedErr != nil {
		fmt.Printf("error when initializing page properties: %s\n", titleInvertedErr)
	}
	defer titleInvertedIndexer.Backup()
	defer titleInvertedIndexer.Release()

	contentInvertedIndexer := &Indexer.InvertedFileIndexer{}
	contentInvertedErr := contentInvertedIndexer.Initialize(wd + "/db/contentInvertedIndex")
	if contentInvertedErr != nil {
		fmt.Printf("error when initializing page properties: %s\n", contentInvertedErr)
	}
	defer contentInvertedIndexer.Backup()
	defer contentInvertedIndexer.Release()

	documentWordForwardIndexer := &Indexer.ForwardIndexer{}
	documentWordForwardIndexerErr := documentWordForwardIndexer.Initialize(wd + "/db/documentWordForwardIndex")
	if documentWordForwardIndexerErr != nil {
		fmt.Printf("error when initializing document -> word forward Indexer: %s\n", documentWordForwardIndexerErr)
	}
	defer documentWordForwardIndexer.Backup()
	defer documentWordForwardIndexer.Release()

	parentChildDocumentForwardIndexer := &Indexer.ForwardIndexer{}
	parentChildDocumentForwardIndexerErr := parentChildDocumentForwardIndexer.Initialize(wd + "/db/parentChildDocumentForwardIndex")
	if parentChildDocumentForwardIndexerErr != nil {
		fmt.Printf("error when initializing parentDocument -> childDocument forward Indexer: %s\n", parentChildDocumentForwardIndexerErr)
	}
	defer parentChildDocumentForwardIndexer.Backup()
	defer parentChildDocumentForwardIndexer.Release()

	childParentDocumentForwardIndexer := &Indexer.ForwardIndexer{}
	childParentDocumentForwardIndexerErr := childParentDocumentForwardIndexer.Initialize(wd + "/db/childParentDocumentForwardIndex")
	if childParentDocumentForwardIndexerErr != nil {
		fmt.Printf("error when initializing childDocument -> parentDocument forward Indexer: %s\n", childParentDocumentForwardIndexerErr)
	}
	defer childParentDocumentForwardIndexer.Backup()
	defer childParentDocumentForwardIndexer.Release()

	pages := []page{}

	crawler := colly.NewCollector(
		colly.MaxDepth(3),
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

		title := e.ChildText("title")
		url := e.Request.URL.String()

		size, _ := strconv.Atoi(e.Response.Headers.Get("Content-Length"))

		if size == 0 {
			size = len(e.Text)
		}

		date := e.Response.Headers.Get("Last-Modified")
		dateTime := time.Time{}
		if len(date) != 0 {
			dateTime, _ = time.Parse(time.RFC1123, date)
		} else {
			dateTime = time.Now()
		}

		// Store Document id and properties
		_, err := documentIndexer.AddKeyToIndex(url)
		if err != nil {
			fmt.Println(err)
		}
		id, _ := documentIndexer.GetValueFromKey(url)
		pagePropertiesIndexer.AddKeyToPageProperties(id, Indexer.CreatePage(id, title, url, size, dateTime))

		text := e.ChildText("body")

		// Remove javascripts and styles in page text
		e.ForEach("script", func(_ int, elem *colly.HTMLElement) {
			text = strings.Replace(text, elem.Text, " ", 1)  
		})
		e.ForEach("style", func(_ int, elem *colly.HTMLElement) {
			text = strings.Replace(text, elem.Text, " ", 1)  
		})

		// Preprocess page text
		content := tokenizer.Tokenize(e.ChildText(text))

		processedTitle := tokenizer.Tokenize(title)

		titleWordList := make(map[uint64]*Indexer.InvertedFile)
		for i, v := range processedTitle {
			// Add Word to id index
			wordID, err := wordIndexer.GetValueFromKey(v)
			if err != nil {
				wordID, _ = wordIndexer.AddKeyToIndex(v)
			}

			invFile, contain := titleWordList[wordID]
			if contain {
				invFile.AddWordPositions(uint64(i))
			} else {
				titleWordList[wordID] = Indexer.CreateInvertedFile(id)
				titleWordList[wordID].AddWordPositions(uint64(i))
			}
		}

		for k, v := range titleWordList {
			titleInvertedIndexer.AddKeyToIndexOrUpdate(k, *v)
		}

		// Check for duplicate words in the document
		contentWordList := make(map[uint64]*Indexer.InvertedFile)
		for i, v := range content {
			// Add Word to id index
			wordID, err := wordIndexer.GetValueFromKey(v)
			if err != nil {
				wordID, _ = wordIndexer.AddKeyToIndex(v)
			}

			invFile, contain := contentWordList[wordID]
			if contain {
				invFile.AddWordPositions(uint64(i))
			} else {
				contentWordList[wordID] = Indexer.CreateInvertedFile(id)
				contentWordList[wordID].AddWordPositions(uint64(i))
			}
		}

		for k, v := range contentWordList {
			contentInvertedIndexer.AddKeyToIndexOrUpdate(k, *v)
		}

		// Get List word words in this document in slice
		wordUIntList := make([]uint64, 0, len(contentWordList))
		for k := range contentWordList {
			wordUIntList = append(wordUIntList, k)
		}

		documentWordForwardIndexer.AddIdListToKey(id, wordUIntList)

		temp := page{}

		temp.children = []uint64{}

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
				childID, _ := documentIndexer.GetValueFromKey(url)
				if temp.children != temp.children.Contains(childID) {
					temp.children = append(temp.children, childID)
				}

				e.Request.Visit(url)
			}(url)
		}
		wg.Wait()
		pages = append(pages, temp)
		fmt.Println(temp.children)
		parentChildDocumentForwardIndexer.AddIdListToKey(id, temp.children)
	})

	// After finished, iterate of

	crawler.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	crawler.Visit(rootPage)

	crawler.Wait()

	contentInvertedIndexer.Iterate()
	wordIndexer.Iterate()
	documentIndexer.Iterate()
	pagePropertiesIndexer.Iterate()
	documentWordForwardIndexer.Iterate()
	parentChildDocumentForwardIndexer.Iterate()
}

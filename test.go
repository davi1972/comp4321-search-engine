package main

import (
	Indexer "github.com/hskrishandi/comp4321/indexer"
	"fmt"
	"os"
)

func main() {

	wd, _ := os.Getwd()

	pagePropertiesIndexer := &Indexer.PagePropetiesIndexer{}
	pagePropertiesErr := pagePropertiesIndexer.Initialize(wd + "/db/pagePropertiesIndex")
	if pagePropertiesErr != nil {
		fmt.Printf("error when initializing page properties: %s\n", pagePropertiesErr)
	}
	defer pagePropertiesIndexer.Backup()
	defer pagePropertiesIndexer.Release()

	documentWordForwardIndexer := &Indexer.DocumentWordForwardIndexer{}
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

	reverseDocumentIndexer := &Indexer.ReverseMappingIndexer{}
	reverseDocumentIndexerErr := reverseDocumentIndexer.Initialize(wd + "/db/reverseDocumentIndexer")
	if reverseDocumentIndexerErr != nil {
		fmt.Printf("error when initializing reverse document indexer: %s\n", reverseDocumentIndexerErr)
	}
	defer reverseDocumentIndexer.Backup()
	defer reverseDocumentIndexer.Release()

	reverseWordindexer := &Indexer.ReverseMappingIndexer{}
	reverseWordindexerErr := reverseWordindexer.Initialize(wd + "/db/reverseWordIndexer")
	if reverseWordindexerErr != nil {
		fmt.Printf("error when initializing reverse word indexer: %s\n", reverseWordindexerErr)
	}
	defer reverseWordindexer.Backup()
	defer reverseWordindexer.Release()

	reverseDocumentIndexer.Iterate()


	pages, err := pagePropertiesIndexer.All()

	if(err != nil){
		fmt.Println(err)
	}

	fmt.Println()

	for _, page := range pages {

		if(page.GetSize()==0){
			continue
		}

		fmt.Println(page.GetTitle())
		fmt.Println(page.GetUrl())
		fmt.Println(page.GetDate()+",", page.GetSize(), "B")

		termFreq, _ := documentWordForwardIndexer.GetWordFrequencyListFromKey(page.GetId())

		fmt.Println(termFreq)

		for _, tf := range termFreq {
			word, _ := reverseWordindexer.GetValueFromKey(tf.GetID())
			fmt.Print(word, tf.GetFrequency(), ", ")
		}

		fmt.Println("Children:")
		
		children, _ := parentChildDocumentForwardIndexer.GetIdListFromKey(page.GetId())

		fmt.Println(children)

		for _, child := range children {
			
			childUrl, _ := reverseDocumentIndexer.GetValueFromKey(child)
			fmt.Println(childUrl)
		}

		fmt.Println("------------------------------------------------------------------------")
	}

}

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

	pages, err := pagePropertiesIndexer.All()

	if(err != nil){
		fmt.Println(err)
	}

	for _, page := range pages {
		fmt.Println(page.GetTitle())
		fmt.Println(page.GetUrl())
		fmt.Println(page.GetDate()+",", page.GetSize(), "B")

		// wordIds, err := documentWordForwardIndexer.GetIdListFromKey(page.GetId())

		// if(err != nil){
		// 	fmt.Println(err)
		// }

		fmt.Println()
		fmt.Println("Children:")
		
		children, _ := parentChildDocumentForwardIndexer.GetIdListFromKey(page.GetId())
		fmt.Println(children)
		for _, child := range children {
			
			childPage, _ := pagePropertiesIndexer.GetPagePropertiesFromKey(child)
			fmt.Println(childPage.GetUrl())
		}

		fmt.Println("------------------------------------------------------------------------")
	}

}
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


	documentWordForwardIndexer := &Indexer.ForwardIndexer{}
	documentWordForwardIndexerErr := documentWordForwardIndexer.Initialize(wd + "/db/documentWordForwardIndex")
	if documentWordForwardIndexerErr != nil {
		fmt.Printf("error when initializing document -> word forward Indexer: %s\n", documentWordForwardIndexerErr)
	}
	defer documentWordForwardIndexer.Backup()
	defer documentWordForwardIndexer.Release()

	pages, err := pagePropertiesIndexer.All()

	if(err != nil){
		fmt.Println(err)
	}

	for _, page := range pages {
		fmt.Println(page.GetTitle())
		fmt.Println(page.GetUrl())
		fmt.Println(page.GetDate()+",", page.GetSize(), "B")

		

		fmt.Println("------------------------------------------------------------------------")
	}

}
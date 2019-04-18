package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	Indexer "github.com/davi1972/comp4321-search-engine/indexer"
)

func main() {
	wd, _ := os.Getwd()

	documentIndexer := &Indexer.MappingIndexer{}
	docErr := documentIndexer.Initialize(wd + "/db/documentIndex")
	if docErr != nil {
		fmt.Printf("error when initializing document indexer: %s\n", docErr)
	}
	defer documentIndexer.Backup()
	defer documentIndexer.Release()

	reverseDocumentIndexer := &Indexer.ReverseMappingIndexer{}
	reverseDocumentIndexerErr := reverseDocumentIndexer.Initialize(wd + "/db/reverseDocumentIndexer")
	if reverseDocumentIndexerErr != nil {
		fmt.Printf("error when initializing reverse document indexer: %s\n", reverseDocumentIndexerErr)
	}
	defer reverseDocumentIndexer.Backup()
	defer reverseDocumentIndexer.Release()

	wordIndexer := &Indexer.MappingIndexer{}
	wordErr := wordIndexer.Initialize(wd + "/db/wordIndex")
	if wordErr != nil {
		fmt.Printf("error when initializing word indexer: %s\n", wordErr)
	}
	defer wordIndexer.Backup()
	defer wordIndexer.Release()

	reverseWordindexer := &Indexer.ReverseMappingIndexer{}
	reverseWordindexerErr := reverseWordindexer.Initialize(wd + "/db/reverseWordIndexer")
	if reverseWordindexerErr != nil {
		fmt.Printf("error when initializing reverse word indexer: %s\n", reverseWordindexerErr)
	}
	defer reverseWordindexer.Backup()
	defer reverseWordindexer.Release()

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

	childParentDocumentForwardIndexer := &Indexer.ForwardIndexer{}
	childParentDocumentForwardIndexerErr := childParentDocumentForwardIndexer.Initialize(wd + "/db/childParentDocumentForwardIndex")
	if childParentDocumentForwardIndexerErr != nil {
		fmt.Printf("error when initializing childDocument -> parentDocument forward Indexer: %s\n", childParentDocumentForwardIndexerErr)
	}
	defer childParentDocumentForwardIndexer.Backup()
	defer childParentDocumentForwardIndexer.Release()

	fmt.Println("Select 1 to 9:")
	fmt.Println("1 - documentIndexer \n2 - reverseDocumentIndexer \n3 - contentInvertedIndexer \n4 - wordIndexer \n5 - reverseWordIndexer \n6 - pagePropertiesIndexer \n7 - documentWordForwardIndexer \n8 - parentChildDocumentForwardIndexer \n9 - childParentDocumentForwardIndexer")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println("Select 1 to 10:")
		fmt.Println("1 - documentIndexer \n2 - reverseDocumentIndexer \n3 - contentInvertedIndexer \n4 - wordIndexer \n5 - reverseWordIndexer \n6 - pagePropertiesIndexer \n7 - documentWordForwardIndexer \n8 - parentChildDocumentForwardIndexer \n9 - childParentDocumentForwardIndexer\n10 - Input Doc ID")

		i, _ := strconv.Atoi(scanner.Text())
		switch i {
		case 1:
			documentIndexer.Iterate()
		case 2:
			reverseDocumentIndexer.Iterate()
		case 3:
			contentInvertedIndexer.Iterate()
		case 4:
			wordIndexer.Iterate()
		case 5:
			reverseWordindexer.Iterate()
		case 6:
			pagePropertiesIndexer.Iterate()
		case 7:
			documentWordForwardIndexer.Iterate()
		case 8:
			parentChildDocumentForwardIndexer.Iterate()
		case 9:
			childParentDocumentForwardIndexer.Iterate()
		}

		if i == 10 {
			scanner.Scan()
			id, _ := strconv.ParseUint(scanner.Text(), 10, 64)
			wflist, _ := documentWordForwardIndexer.GetWordFrequencyListFromKey(id)

			for j := range wflist {
				fmt.Println("key: " + fmt.Sprint(wflist[j].GetID()) + " value: " + fmt.Sprint(wflist[j].GetFrequency()) + " wordID: " + fmt.Sprint(wflist[j].GetWordID()))

			}
		}
		// documentIndexer.Iterate()
		// reverseDocumentIndexer.Iterate()
		// contentInvertedIndexer.Iterate()
		// wordIndexer.Iterate()
		// reverseWordindexer.Iterate()
		// pagePropertiesIndexer.Iterate()
		// documentWordForwardIndexer.Iterate()
		// parentChildDocumentForwardIndexer.Iterate()
		// childParentDocumentForwardIndexer.Iterate()
	}
}

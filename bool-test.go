package main

import (
	"bufio"
	"fmt"
	"os"

	"./boolsearch"

	Indexer "./indexer"
	"./vsm"
	"github.com/davi1972/comp4321-search-engine/tokenizer"
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

	titleWordForwardIndexer := &Indexer.DocumentWordForwardIndexer{}
	titleWordForwardIndexerErr := titleWordForwardIndexer.Initialize(wd + "/db/titleWordForwardIndex")
	if titleWordForwardIndexerErr != nil {
		fmt.Printf("error when initializing title word forward Indexer: %s\n", titleWordForwardIndexerErr)
	}
	defer titleWordForwardIndexer.Backup()
	defer titleWordForwardIndexer.Release()

	v := &vsm.VSM{
		DocumentIndexer:                   documentIndexer,
		WordIndexer:                       wordIndexer,
		ReverseDocumentIndexer:            reverseDocumentIndexer,
		ReverseWordIndexer:                reverseWordindexer,
		PagePropertiesIndexer:             pagePropertiesIndexer,
		TitleInvertedIndexer:              titleInvertedIndexer,
		ContentInvertedIndexer:            contentInvertedIndexer,
		DocumentWordForwardIndexer:        documentWordForwardIndexer,
		ParentChildDocumentForwardIndexer: parentChildDocumentForwardIndexer,
		ChildParentDocumentForwardIndexer: childParentDocumentForwardIndexer,
		TitleWordForwardIndexer:           titleWordForwardIndexer,
	}

	bs := &boolsearch.BoolSearch{
		ContentInvertedIndexer: contentInvertedIndexer,
		Vsm:                    v,
	}

	fmt.Println("\nboolSearch.go")
	fmt.Println("--------------")
	fmt.Printf("Search: (enter keywords and press enter)\n")

	scanner := bufio.NewScanner(os.Stdin)
	var q string
	for scanner.Scan() {
		q = scanner.Text()

		fmt.Printf("Results for %s:\n", q)
		query := tokenizer.Tokenize(q)
		barr := bs.FindBoolean(query)

		for k, v := range barr {
			fmt.Println("Boolean Search:")
			fmt.Println("k: ", k, " v: ", v)

		}
		fmt.Println("\nboolSearch.go")
		fmt.Println("--------------")
		fmt.Printf("Search: (enter keywords and press enter)\n")

	}
}

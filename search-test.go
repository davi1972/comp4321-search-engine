package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	Indexer "github.com/davi1972/comp4321-search-engine/indexer"
	"github.com/davi1972/comp4321-search-engine/tokenizer"

	//"github.com/davi1972/comp4321-search-engine/vsm"
	"./vsm"
)

func main() {
	fmt.Println("search-test.go")
	fmt.Println("--------------")

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
	}

	fmt.Println("\n\nsearch-test.go")
	fmt.Println("--------------\n\n")

	// stringToWordID test
	fmt.Println("StringToWordID Test:")
	fmt.Println("--------------------")
	query := tokenizer.Tokenize("paragraph document")
	stwid1, err1 := v.StringToWordID(query[0]) // 4 - ok
	fmt.Println(stwid1)
	fmt.Println(err1)
	fmt.Println("\n--------------------")
	stwid2, err2 := v.StringToWordID(query[1]) // 78 - ok
	fmt.Println(stwid2)
	fmt.Println(err2)
	fmt.Println("\n--------------------")

	// InverseDocFreq test
	fmt.Println("InverseDocumentFreq Test:")
	fmt.Println("-------------------------")
	idf1, err3 := v.InverseDocumentFreq(query[0]) // 1.222392421336448
	fmt.Println(idf1)
	fmt.Println(err3)
	fmt.Println("\n-------------------------")
	idf2, err4 := v.InverseDocumentFreq(query[1]) // 2.807354922057604
	fmt.Println(idf2)
	fmt.Println(err4)
	fmt.Println("\n-------------------------")

	// TermFreq test
	fmt.Println("TermFreq Test:")
	fmt.Println("--------------")
	tf1, err5 := v.TermFreq(query[0], 0) // doc 0: 4 - paragraph
	fmt.Println(tf1)
	fmt.Println(err5)
	fmt.Println("\n--------------")
	tf2, err6 := v.TermFreq(query[1], 0) // doc 1: 0 - document
	fmt.Println(tf2)
	fmt.Println(err6)
	fmt.Println("\n--------------")

	// ComputeTermWeight test
	fmt.Println("ComputeTermWeight Test:")
	fmt.Println("-----------------------")
	ctw1 := v.ComputeTermWeight(query[0], 0) // 0.9779139370691585 - doc 0
	fmt.Println(ctw1)
	//fmt.Println(err7)
	fmt.Println("\n-----------------------")
	ctw2 := v.ComputeTermWeight(query[1], 0) // 0 - doc 0
	fmt.Println(ctw2)
	//fmt.Println(err8)
	fmt.Println("\n-----------------------")

	// MaxTermFreq test
	fmt.Println("maxTermFreq Test:")
	fmt.Println("-----------------")
	mtf1 := v.MaxTermFreq(0) // doc 0 - 5
	fmt.Println(mtf1)
	fmt.Println("\n-----------------")
	mtf2 := v.MaxTermFreq(1) // doc 1 - 3
	fmt.Println(mtf2)
	fmt.Println("\n-----------------")

	// CosSimilarity test
	fmt.Println("CosSimilarity Test:")
	fmt.Println("-------------------")
	coss1 := v.CosSimilarity("paragraph document", 0) // 0.3992213689757632
	fmt.Println(coss1)
	//fmt.Println(err7)
	fmt.Println("\n-------------------")
	coss2 := v.ComputeTermWeight("paragraph document", 1) // 0.602451640685868
	fmt.Println(coss2)
	//fmt.Println(err8)
	fmt.Println("\n-------------------")

	// ComputeCosineScore test
	fmt.Println("ComputeCosineScore Test:")
	fmt.Println("------------------------")
	scores := v.ComputeCosineScore("paragraph document")
	fmt.Println(len(scores))
	fmt.Println("\n------------------------")
	fmt.Println(scores[0]) // 0.004435792988619592
	fmt.Println("\n------------------------")
	fmt.Println(scores[1]) // 0.019961068448788165
	fmt.Println("\n------------------------")
	fmt.Println(scores[2]) // 0.25
	fmt.Println("\n------------------------")

	fmt.Println("\n\nsearch-test.go")
	fmt.Println("--------------")
	fmt.Println("Search: (enter keywords below and press enter)")

	scanner := bufio.NewScanner(os.Stdin)
	var q string = ""
	for scanner.Scan() {
		q = scanner.Text()
		fmt.Printf("Results for %s:\n", q)
		score := v.ComputeCosineScore(q)
		for i := range score {
			fmt.Println("Doc" + strconv.Itoa(i) + " CosSim Score: " + fmt.Sprint(score[i]))
		}
		fmt.Println("\nsearch-test.go")
		fmt.Println("--------------")
		fmt.Println("Search: (enter keywords below and press enter)")

	}
}

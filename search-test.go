package main

import (
	"fmt"

	"github.com/davi1972/comp4321-search-engine/tokenizer"
	"github.com/davi1972/comp4321-search-engine/vsm"
)

// func computeCosineScore(query string) []float64 { // score normalized by doclength
// 	scores := make([]float64, 3)
// 	lengths := make([]float64, 3)

// 	docIndex := Indexer.DocumentWordForwardIndexer{}

// 	// init length list
// 	//var i uint64
// 	var length float64
// 	for i := 0; i < 3; i++ {
// 		length = 0
// 		docList, _ := docIndex.GetWordFrequencyListFromKey(uint64(i))
// 		for j := range docList {
// 			length += float64(docList[j].GetFrequency())
// 		}
// 		lengths[i] = length
// 	}

// 	for c := 0; c < 3; c++ {
// 		scores[c] = vsm.CosSimilarity(query, uint64(c))
// 		scores[c] = scores[c] / lengths[c]
// 	}

// 	return scores
// }

func main() {
	fmt.Println("search-test.go")
	fmt.Println("--------------")

	// wd, _ := os.Getwd()

	// documentIndexer := &Indexer.MappingIndexer{}
	// docErr := documentIndexer.Initialize(wd + "/db/documentIndex")
	// if docErr != nil {
	// 	fmt.Printf("error when initializing document indexer: %s\n", docErr)
	// }
	// defer documentIndexer.Backup()
	// defer documentIndexer.Release()

	// reverseDocumentIndexer := &Indexer.ReverseMappingIndexer{}
	// reverseDocumentIndexerErr := reverseDocumentIndexer.Initialize(wd + "/db/reverseDocumentIndexer")
	// if reverseDocumentIndexerErr != nil {
	// 	fmt.Printf("error when initializing reverse document indexer: %s\n", reverseDocumentIndexerErr)
	// }
	// defer reverseDocumentIndexer.Backup()
	// defer reverseDocumentIndexer.Release()

	// wordIndexer := &Indexer.MappingIndexer{}
	// wordErr := wordIndexer.Initialize(wd + "/db/wordIndex")
	// if wordErr != nil {
	// 	fmt.Printf("error when initializing word indexer: %s\n", wordErr)
	// }
	// defer wordIndexer.Backup()
	// defer wordIndexer.Release()

	// reverseWordindexer := &Indexer.ReverseMappingIndexer{}
	// reverseWordindexerErr := reverseWordindexer.Initialize(wd + "/db/reverseWordIndexer")
	// if reverseWordindexerErr != nil {
	// 	fmt.Printf("error when initializing reverse word indexer: %s\n", reverseWordindexerErr)
	// }
	// defer reverseWordindexer.Backup()
	// defer reverseWordindexer.Release()

	// pagePropertiesIndexer := &Indexer.PagePropetiesIndexer{}
	// pagePropertiesErr := pagePropertiesIndexer.Initialize(wd + "/db/pagePropertiesIndex")
	// if pagePropertiesErr != nil {
	// 	fmt.Printf("error when initializing page properties: %s\n", pagePropertiesErr)
	// }
	// defer pagePropertiesIndexer.Backup()
	// defer pagePropertiesIndexer.Release()

	// titleInvertedIndexer := &Indexer.InvertedFileIndexer{}
	// titleInvertedErr := titleInvertedIndexer.Initialize(wd + "/db/titleInvertedIndex")
	// if titleInvertedErr != nil {
	// 	fmt.Printf("error when initializing page properties: %s\n", titleInvertedErr)
	// }
	// defer titleInvertedIndexer.Backup()
	// defer titleInvertedIndexer.Release()

	// contentInvertedIndexer := &Indexer.InvertedFileIndexer{}
	// contentInvertedErr := contentInvertedIndexer.Initialize(wd + "/db/contentInvertedIndex")
	// if contentInvertedErr != nil {
	// 	fmt.Printf("error when initializing page properties: %s\n", contentInvertedErr)
	// }
	// defer contentInvertedIndexer.Backup()
	// defer contentInvertedIndexer.Release()

	// documentWordForwardIndexer := &Indexer.DocumentWordForwardIndexer{}
	// documentWordForwardIndexerErr := documentWordForwardIndexer.Initialize(wd + "/db/documentWordForwardIndex")
	// if documentWordForwardIndexerErr != nil {
	// 	fmt.Printf("error when initializing document -> word forward Indexer: %s\n", documentWordForwardIndexerErr)
	// }
	// defer documentWordForwardIndexer.Backup()
	// defer documentWordForwardIndexer.Release()

	// parentChildDocumentForwardIndexer := &Indexer.ForwardIndexer{}
	// parentChildDocumentForwardIndexerErr := parentChildDocumentForwardIndexer.Initialize(wd + "/db/parentChildDocumentForwardIndex")
	// if parentChildDocumentForwardIndexerErr != nil {
	// 	fmt.Printf("error when initializing parentDocument -> childDocument forward Indexer: %s\n", parentChildDocumentForwardIndexerErr)
	// }
	// defer parentChildDocumentForwardIndexer.Backup()
	// defer parentChildDocumentForwardIndexer.Release()

	// childParentDocumentForwardIndexer := &Indexer.ForwardIndexer{}
	// childParentDocumentForwardIndexerErr := childParentDocumentForwardIndexer.Initialize(wd + "/db/childParentDocumentForwardIndex")
	// if childParentDocumentForwardIndexerErr != nil {
	// 	fmt.Printf("error when initializing childDocument -> parentDocument forward Indexer: %s\n", childParentDocumentForwardIndexerErr)
	// }
	// defer childParentDocumentForwardIndexer.Backup()
	// defer childParentDocumentForwardIndexer.Release()

	// // fmt.Println("\n\nsearch-test.go")
	// // fmt.Println("--------------")
	// // fmt.Println("Search: (enter keywords below and press enter)")

	// // scanner := bufio.NewScanner(os.Stdin)
	// // var query string = ""
	// // for scanner.Scan() {
	// // 	query = scanner.Text()
	// // 	fmt.Println("Results:")
	// // 	scores := computeCosineScore(query)
	// // 	for i := range scores {
	// // 		fmt.Println("Doc" + strconv.Itoa(i) + " CosSim Score: " + fmt.Sprint(scores[i]))
	// // 	}
	// // 	fmt.Println("\nsearch-test.go")
	// // 	fmt.Println("--------------")
	// // 	fmt.Println("Search: (enter keywords below and press enter)")

	// // }

	// stringToWordID test
	query := tokenizer.Tokenize("paragraph")
	fmt.Sprintln(vsm.StringToWordID(query[0]))

}

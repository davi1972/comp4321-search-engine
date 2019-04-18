package vsm

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

	Indexer "github.com/davi1972/comp4321-search-engine/indexer"
	"github.com/davi1972/comp4321-search-engine/tokenizer"
)

// Returns a wordid given a (tokenized) term.
func StringToWordID(qterm string) uint64 {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error when getting wd: %s\n", err)
	}
	parent := filepath.Dir(wd)
	wordsIndex := &Indexer.MappingIndexer{}
	wordsErr := wordsIndex.Initialize(parent + "/db/wordIndex")
	if wordsErr != nil {
		fmt.Printf("Error when initializing word indexer: %s\n", wordsErr)
	}
	wordid, _ := wordsIndex.GetValueFromKey(qterm)
	return wordid
}

// Returns the inverse document frequency of a string.
func InverseDocumentFreq(qterm string) float64 {
	// log_2(N/df of term)
	N := 0
	df := 0
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error when getting wd: %s\n", err)
	}
	parent := filepath.Dir(wd)

	pageIndex := &Indexer.PagePropetiesIndexer{}
	pageErr := pageIndex.Initialize(parent + "/db/pagePropertiesIndex")
	if pageErr != nil {
		fmt.Printf("Error when initializing page properties: %s\n", pageErr)
	}
	pages, _ := pageIndex.All()

	contentInvertedIndex := &Indexer.InvertedFileIndexer{}
	contentInvertedErr := contentInvertedIndex.Initialize(parent + "/db/contentInvertedIndex")
	if contentInvertedErr != nil {
		fmt.Printf("Error when initializing page properties: %s\n", contentInvertedErr)
	}
	invertedFile, _ := contentInvertedIndex.GetInvertedFileFromKey(StringToWordID(qterm))

	N = len(pages)
	df = len(invertedFile)
	return math.Log2(float64(N) / float64(df))
}

// Returns the term frequency of a term in document (ID).
func TermFreq(qterm string, documentID uint64) uint64 {
	// frequency of term j in document i
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error when getting wd: %s\n", err)
	}
	parent := filepath.Dir(wd)

	docWords := &Indexer.DocumentWordForwardIndexer{}
	docWordsErr := docWords.Initialize(parent + "/db/documentWordForwardIndex")
	if docWordsErr != nil {
		fmt.Printf("error when initializing Document -> Word Forward Indexer: %s\n", docWordsErr)
	}
	words, _ := docWords.GetWordFrequencyListFromKey(documentID)

	wordsIndex := &Indexer.MappingIndexer{}
	wordsIndexErr := wordsIndex.Initialize(parent + "/db/wordIndex")
	if wordsIndexErr != nil {
		fmt.Printf("Error when initializing word indexer: %s\n", wordsIndexErr)
	}
	index, _ := wordsIndex.GetValueFromKey(qterm) // word id

	// iterate through doc's word IDs
	for i := range words {
		if words[i].GetID() == index {
			return words[i].GetFrequency()
		}
	}
	return 0
}

// Returns the computed term weight of a (tokenized) term given a string and document (ID).
func ComputeTermWeight(qterm string, documentID uint64) float64 {
	return float64(TermFreq(qterm, documentID)/MaxTermFreq(documentID)) * InverseDocumentFreq(qterm)
}

// Returns the maximum term frequency of a term in a document ID.
func MaxTermFreq(documentID uint64) uint64 {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error when getting wd: %s\n", err)
	}
	parent := filepath.Dir(wd)

	docWords := &Indexer.DocumentWordForwardIndexer{}
	docWordsErr := docWords.Initialize(parent + "/db/documentWordForwardIndex")
	if docWordsErr != nil {
		fmt.Printf("error when initializing Document -> Word Forward Indexer: %s\n", docWordsErr)
	}
	words, _ := docWords.GetWordFrequencyListFromKey(documentID)
	//wordID := stringToWordID(qterm)

	wf := words[0]
	for i := range words[1:] {
		if words[i].GetFrequency() > wf.GetFrequency() {
			wf = words[i]
		}
	}

	return wf.GetFrequency()
}

// Returns the cosine similarity between query and document ID.
func CosSimilarity(query string, documentID uint64) float64 {
	terms := tokenizer.Tokenize(query)
	termWeights := make(map[string]float64)
	queryFreq := make(map[string]int)

	for i := range terms {
		termWeights[terms[i]] = ComputeTermWeight(terms[i], documentID)
		queryFreq[terms[i]]++
	}
	// dik is weight of term k in doc i, qk is weight of term k in query
	innerPro := 0.0
	sumD := 0.0
	sumQ := 0.0

	for i := 0; i < len(queryFreq); i++ {
		invDocFreq := InverseDocumentFreq(terms[i])
		innerPro += termWeights[terms[i]] * (float64(queryFreq[terms[i]]) * invDocFreq)
		sumD += termWeights[terms[i]] * termWeights[terms[i]]
		sumQ += (float64(queryFreq[terms[i]]) * invDocFreq) * (float64(queryFreq[terms[i]]) * invDocFreq)
	}

	return innerPro / (math.Sqrt(sumD) * math.Sqrt(sumQ))
}

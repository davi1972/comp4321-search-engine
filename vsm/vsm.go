package vsm

import (
	"math"

	Indexer "github.com/davi1972/comp4321-search-engine/indexer"
	"github.com/davi1972/comp4321-search-engine/tokenizer"
)

func stringToWordID(qterm string) uint64 {
	wordsindex := &Indexer.MappingIndexer{}
	wordid, _ := wordsindex.GetValueFromKey(qterm)
	return wordid
}

// Returns the inverse document frequency
func inverseDocumentFreq(qterm string, documentId uint64) float64 {
	// log_2(N/df of term)
	N := 0
	df := 0
	pageindex := &Indexer.PagePropetiesIndexer{}
	pages, _ := pageindex.All()

	// get N
	for i := range pages {
		N++
	}

	// get no. of docs with qterm
	invertedIndex := &Indexer.InvertedFileIndexer{}
	invertedFile, _ := invertedIndex.GetInvertedFileFromKey(stringToWordID(qterm))

	for i := range invertedFile {
		df++
	}

	return math.Log2(float64(N) / float64(df))
}

// Returns the term frequency of a term in document
func termFreq(qterm string, documentID uint64) uint64 {
	// frequency of term j in document i
	docwords := &Indexer.DocumentWordForwardIndexer{}
	words, _ := docwords.GetWordFrequencyListFromKey(documentID)
	wordsindex := &Indexer.MappingIndexer{}
	index, _ := wordsindex.GetValueFromKey(qterm) // word id

	// iterate through doc's word IDs
	for i := range words {
		if words[i].GetID() == index {
			return words[i].GetFrequency()
		}
	}
	return 0
}

func computeTermWeight(qterm string, documentID uint64) float64 {
	maxtf := maxTermFreq(qterm, documentID)
	return float64(termFreq(qterm, documentID)/maxtf) * inverseDocumentFreq(qterm, documentID)
}

func computeQueryWeight(qterm string) float64 {

}

func maxTermFreq(qterm string, documentID uint64) uint64 {
	docwords := &Indexer.DocumentWordForwardIndexer{}
	words, _ := docwords.GetWordFrequencyListFromKey(documentID)
	wordID := stringToWordID(qterm)

	wf := words[0]
	for i := range words[1:] {
		if words[i].GetFrequency() > wf.GetFrequency() {
			wf = words[i]
		}
	}

	return wf.GetFrequency()
}

// Returns the cosine similarity between query and document ID
func cosSimilarity(query string, documentID uint64) float64 {
	// use doc id to find terms
	// forea(dw * qw) / forea(doclen) * forea(qlen)
	similarity := 0.0
	terms := tokenizer.Tokenize(query)
	docwords := &Indexer.DocumentWordForwardIndexer{}
	words, _ := docwords.GetWordFrequencyListFromKey(documentID)

	// dik is weight of term k in doc i, qk is weight of term k in query
	var innerPro float64 = 0.0
	for i := range terms {
		innerPro += computeTermWeight(terms[i], documentID) + computeQueryWeight(terms[i])
	}

}

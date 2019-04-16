package vsm

import (
	"math"

	Indexer "github.com/davi1972/comp4321-search-engine/indexer"
	"github.com/davi1972/comp4321-search-engine/tokenizer"
)

// Returns a wordid given a (tokenized) term.
func StringToWordID(qterm string) uint64 {
	wordsindex := &Indexer.MappingIndexer{}
	wordid, _ := wordsindex.GetValueFromKey(qterm)
	return wordid
}

// Returns the inverse document frequency of a string.
func InverseDocumentFreq(qterm string) float64 {
	// log_2(N/df of term)
	N := 0
	df := 0

	pageindex := &Indexer.PagePropetiesIndexer{}
	pages, _ := pageindex.All()

	invertedIndex := &Indexer.InvertedFileIndexer{}
	invertedFile, _ := invertedIndex.GetInvertedFileFromKey(StringToWordID(qterm))

	N = len(pages)
	df = len(invertedFile)
	return math.Log2(float64(N) / float64(df))
}

// Returns the term frequency of a term in document (ID).
func TermFreq(qterm string, documentID uint64) uint64 {
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

// Returns the computed term weight of a (tokenized) term given a string and document (ID).
func ComputeTermWeight(qterm string, documentID uint64) float64 {
	return float64(TermFreq(qterm, documentID)/MaxTermFreq(documentID)) * InverseDocumentFreq(qterm)
}

// Returns the maximum term frequency of a term in a document ID.
func MaxTermFreq(documentID uint64) uint64 {
	docwords := &Indexer.DocumentWordForwardIndexer{}
	words, _ := docwords.GetWordFrequencyListFromKey(documentID)
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

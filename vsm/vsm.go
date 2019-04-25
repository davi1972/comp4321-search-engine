package vsm

import (
	"fmt"
	"math"

	Indexer "github.com/davi1972/comp4321-search-engine/indexer"
	"github.com/davi1972/comp4321-search-engine/tokenizer"
)

type VSM struct {
	documentIndexer                   *Indexer.MappingIndexer
	wordIndexer                       *Indexer.MappingIndexer
	reverseDocumentIndexer            *Indexer.ReverseMappingIndexer
	reverseWordIndexer                *Indexer.ReverseMappingIndexer
	pagePropertiesIndexer             *Indexer.PagePropetiesIndexer
	titleInvertedIndexer              *Indexer.InvertedFileIndexer
	contentInvertedIndexer            *Indexer.InvertedFileIndexer
	documentWordForwardIndexer        *Indexer.DocumentWordForwardIndexer
	parentChildDocumentForwardIndexer *Indexer.ForwardIndexer
	childParentDocumentForwardIndexer *Indexer.ForwardIndexer
}

// Returns a wordid given a (tokenized) term.
func (vsm *VSM) StringToWordID(qterm string) (uint64, error) {
	wordid, err := vsm.wordIndexer.GetValueFromKey(qterm)
	return wordid, err
}

// Returns the inverse document frequency of a string.
func (vsm *VSM) InverseDocumentFreq(qterm string) (float64, error) {
	N := vsm.documentIndexer.GetSize()
	wordid, err := vsm.StringToWordID(qterm)

	if err != nil {
		err = fmt.Errorf("Error when getting value from key: %s", err)
	}

	df, err2 := vsm.contentInvertedIndexer.GetDocFreq(wordid)

	if err2 != nil {
		err2 = fmt.Errorf("Error when getting inverted file from key: %s", err2)
	}

	return math.Log2(float64(N) / float64(df)), err
}

// Returns the term frequency of a term in document (ID).
func (vsm *VSM) TermFreq(qterm string, documentID uint64) (uint64, error) {
	// frequency of term j in document i
	words, err := vsm.documentWordForwardIndexer.GetWordFrequencyListFromKey(documentID)

	if err != nil {
		err = fmt.Errorf("Error when getting word frequency list from key: %s", err)
	}

	index, err2 := vsm.wordIndexer.GetValueFromKey(qterm) // word id
	if err2 != nil {
		err2 = fmt.Errorf("Error when getting value from key transaction: %s", err2)
	}

	// iterate through doc's word IDs
	for i := range words {
		if words[i].GetID() == index {
			return words[i].GetFrequency(), err
		}
	}
	return 0, err
}

// Returns the computed term weight of a (tokenized) term given a string and document (ID).
func (vsm *VSM) ComputeTermWeight(qterm string, documentID uint64) float64 {
	tf, _ := vsm.TermFreq(qterm, documentID)
	maxtf := vsm.MaxTermFreq(documentID)
	infreq, _ := vsm.InverseDocumentFreq(qterm)
	return float64(tf) / float64(maxtf) * float64(infreq)
}

// Returns the maximum term frequency of a term in a document ID.
func (vsm *VSM) MaxTermFreq(documentID uint64) uint64 {
	words, _ := vsm.documentWordForwardIndexer.GetWordFrequencyListFromKey(documentID)

	wf := words[0]
	for i := range words[1:] {
		if words[i].GetFrequency() > wf.GetFrequency() {
			wf = words[i]
		}
	}
	return wf.GetFrequency()
}

// Returns the cosine similarity between query and document ID.
func (vsm *VSM) CosSimilarity(query string, documentID uint64) float64 {
	terms := tokenizer.Tokenize(query)
	termWeights := make(map[string]float64)
	queryFreq := make(map[string]int)

	for i := range terms {
		termWeights[terms[i]] = vsm.ComputeTermWeight(terms[i], documentID)
		queryFreq[terms[i]]++
	}
	// dik is weight of term k in doc i, qk is weight of term k in query
	innerPro := 0.0
	sumD := 0.0
	sumQ := 0.0

	for i := 0; i < len(queryFreq); i++ {
		invDocFreq, _ := vsm.InverseDocumentFreq(terms[i])
		innerPro += termWeights[terms[i]] * (float64(queryFreq[terms[i]]) * invDocFreq)
		sumD += termWeights[terms[i]] * termWeights[terms[i]]
		sumQ += (float64(queryFreq[terms[i]]) * invDocFreq) * (float64(queryFreq[terms[i]]) * invDocFreq)
	}

	return innerPro / (math.Sqrt(sumD) * math.Sqrt(sumQ))
}

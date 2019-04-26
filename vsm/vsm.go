package vsm

import (
	"fmt"
	"math"

	Indexer "github.com/davi1972/comp4321-search-engine/indexer"
	//Indexer "../indexer"
	"github.com/davi1972/comp4321-search-engine/tokenizer"
)

type VSM struct {
	DocumentIndexer                   *Indexer.MappingIndexer
	WordIndexer                       *Indexer.MappingIndexer
	ReverseDocumentIndexer            *Indexer.ReverseMappingIndexer
	ReverseWordIndexer                *Indexer.ReverseMappingIndexer
	PagePropertiesIndexer             *Indexer.PagePropetiesIndexer
	TitleInvertedIndexer              *Indexer.InvertedFileIndexer
	ContentInvertedIndexer            *Indexer.InvertedFileIndexer
	DocumentWordForwardIndexer        *Indexer.DocumentWordForwardIndexer
	ParentChildDocumentForwardIndexer *Indexer.ForwardIndexer
	ChildParentDocumentForwardIndexer *Indexer.ForwardIndexer
}

// Returns a wordid given a (tokenized) term.
func (vsm *VSM) StringToWordID(qterm string) (uint64, error) {
	wordid, err := vsm.WordIndexer.GetValueFromKey(qterm)
	return wordid, err
}

// Returns the inverse document frequency of a string.
func (vsm *VSM) InverseDocumentFreq(qterm string) (float64, error) {
	N := vsm.DocumentWordForwardIndexer.GetSize()
	fmt.Printf("N = %d\n", N)

	wordid, err := vsm.StringToWordID(qterm)

	if err != nil {
		err = fmt.Errorf("Error when getting value from key: %s", err)
	}

	df, err2 := vsm.ContentInvertedIndexer.GetDocFreq(wordid)
	fmt.Printf("df = %d\n", df)

	if err2 != nil {
		err2 = fmt.Errorf("Error when getting inverted file from key: %s", err2)
	}

	return math.Log2(float64(N) / float64(df)), err
}

// Returns the term frequency of a term in document (ID).
func (vsm *VSM) TermFreq(qterm string, documentID uint64) (uint64, error) {
	// frequency of term j in document i
	words, err := vsm.DocumentWordForwardIndexer.GetWordFrequencyListFromKey(documentID)

	if err != nil {
		err = fmt.Errorf("Error when getting word frequency list from key: %s", err)
	}

	index, err2 := vsm.WordIndexer.GetValueFromKey(qterm) // word id
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
	words, _ := vsm.DocumentWordForwardIndexer.GetWordFrequencyListFromKey(documentID)

	if len(words) > 0 {
		wf := words[0]
		for i := range words[1:] {
			if words[i].GetFrequency() > wf.GetFrequency() {
				wf = words[i]
			}
		}
		return wf.GetFrequency()
	}
	return 0
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
	res := innerPro / (math.Sqrt(sumD) * math.Sqrt(sumQ))

	if math.IsNaN(res) {
		res = 0
	}
	return res
}

// Returns a float array with scores starting with doc 0 as index
func (vsm *VSM) ComputeCosineScore(query string) []float64 {
	size := int(vsm.DocumentWordForwardIndexer.GetSize())
	fmt.Printf("N = %d\n", 0)
	scores := make([]float64, 0)
	lengths := make([]float64, 0)
	docIDList, err1 := vsm.DocumentWordForwardIndexer.GetDocIDList()

	if err1 != nil {
		fmt.Errorf("Error getting doc ID list: %s", err1)
		return nil
	}

	var length float64
	for i := 0; i < len(docIDList); i++ {
		length = 0
		docList, err2 := vsm.DocumentWordForwardIndexer.GetWordFrequencyListFromKey(docIDList[i])
		if err2 != nil {
			fmt.Errorf("Error getting word frequency list: %s", err2)
			return nil
		}
		if len(docList) > 0 {
			for j := range docList {
				length += float64(docList[j].GetFrequency())
			}
		}
		lengths[i] = length
	}

	for c := 0; c < size; c++ {
		scores[c] = vsm.CosSimilarity(query, uint64(c))
		scores[c] = scores[c] / lengths[c]
	}
	return scores
}

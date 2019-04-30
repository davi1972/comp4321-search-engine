package vsm

import (
	"math"

	//Indexer "github.com/davi1972/comp4321-search-engine/indexer"
	Indexer "github.com/davi1972/comp4321-search-engine/indexer"
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
	TitleWordForwardIndexer           *Indexer.DocumentWordForwardIndexer
}

// Returns a wordid given a (tokenized) term.
func (vsm *VSM) StringToWordID(qterm string) (uint64, error) {
	wordid, err := vsm.WordIndexer.GetValueFromKey(qterm)
	return wordid, err
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

// Returns a float array with scores starting with doc 0 as index
func (vsm *VSM) ComputeCosineScore(query string) (map[uint64]float64, error) {
	//fmt.Printf("N = %d\n", 0)
	scores := make(map[uint64]float64)
	queryFreq := make(map[string]int)

	terms := tokenizer.Tokenize(query)
	queryLength := 0.0
	docLength := 0.0
	for _, term := range terms {
		wordID, wordIDErr := vsm.WordIndexer.GetValueFromKey(term)
		if wordIDErr != nil {
			continue
		}
		invFileListContent, _ := vsm.ContentInvertedIndexer.GetInvertedFileFromKey(wordID)
		for _, invFile := range invFileListContent {
			tf := len(invFile.GetWordPositions())

			maxtf := vsm.MaxTermFreq(invFile.GetPageID())
			N := vsm.DocumentWordForwardIndexer.GetSize()
			df := len(invFileListContent)
			infreq := math.Log2(float64(N) / float64(df))
			scores[invFile.GetPageID()] += (float64(tf) / float64(maxtf) * float64(infreq))
			docLength += (float64(tf) * float64(infreq) * float64(tf) * float64(infreq))
		}

		invFileListTitle, _ := vsm.TitleInvertedIndexer.GetInvertedFileFromKey(wordID)
		for _, invFile := range invFileListTitle {
			tf := len(invFile.GetWordPositions())

			maxtf := vsm.MaxTermFreq(invFile.GetPageID())
			N := vsm.TitleWordForwardIndexer.GetSize()
			df := len(invFileListTitle)
			infreq := math.Log2(float64(N) / float64(df))
			scores[invFile.GetPageID()] += (float64(tf) / float64(maxtf) * float64(infreq)) * 1.5 // Special consideration
			docLength += (float64(tf) * float64(infreq) * float64(tf) * float64(infreq)) * 1.5 * 1.5
		}
		queryFreq[term]++
	}

	// Compute query weight
	for k := range queryFreq {
		queryLength += float64(queryFreq[k] * queryFreq[k])
	}
	queryLength = math.Sqrt(queryLength)

	docLength = math.Sqrt(docLength)

	for k := range scores {

		scores[k] /= (docLength * queryLength)

	}

	return scores, nil
}

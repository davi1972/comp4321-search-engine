package phrasalSearch

import (
	"sort"

	"github.com/davi1972/comp4321-search-engine/boolsearch"
	"github.com/davi1972/comp4321-search-engine/vsm"

	//Indexer "github.com/davi1972/comp4321-search-engine/indexer"
	Indexer "github.com/davi1972/comp4321-search-engine/indexer"
)

type PhrasalSearch struct {
	TitleInvertedIndexer    *Indexer.InvertedFileIndexer
	ContentInvertedIndexer  *Indexer.InvertedFileIndexer
	TitleWordForwardIndexer *Indexer.DocumentWordForwardIndexer
	V                       *vsm.VSM
	Bs                      *boolsearch.BoolSearch
}

type BiGram struct {
	phrase1 string
	phrase2 string
}

// Returns an array of phrases (bi-grams)
func (phrases *PhrasalSearch) splitToPhrase(query []string) []BiGram {
	var phrs []BiGram
	for i := 0; i+1 < len(query); i++ {
		phr := BiGram{query[i], query[i+1]}
		phrs = append(phrs, phr)
	}
	return phrs
}

func (phrases *PhrasalSearch) hasPhraseInTitle(documentID uint64, b BiGram) bool {
	wid1, err := phrases.V.StringToWordID(b.phrase1)
	if err != nil {
		return false
	}
	doc1, _ := phrases.TitleInvertedIndexer.GetInvertedFileFromKey(wid1)
	var doc1pos []uint64
	for i := range doc1 {
		if doc1[i].GetPageID() == documentID {
			doc1pos = doc1[i].GetWordPositions()
			break
		}
	}

	wid2, err2 := phrases.V.StringToWordID(b.phrase2)
	if err2 != nil {
		return false
	}
	doc2, _ := phrases.TitleInvertedIndexer.GetInvertedFileFromKey(wid2)
	var doc2pos []uint64
	for i := range doc1 {
		if doc2[i].GetPageID() == documentID {
			doc2pos = doc1[i].GetWordPositions()
			break
		}
	}

	for i := range doc2 {
		doc2pos[i]--
	}
	docsIntersect := phrases.Bs.FindIntersect(doc1pos, doc2pos)

	if len(docsIntersect) > 0 {
		return true
	}
	return false
}

func (phrases *PhrasalSearch) hasPhraseInBody(documentID uint64, b BiGram) bool {
	wid1, err := phrases.V.StringToWordID(b.phrase1)
	if err != nil {
		return false
	}
	doc1, _ := phrases.ContentInvertedIndexer.GetInvertedFileFromKey(wid1)
	var doc1pos []uint64
	for i := range doc1 {
		if doc1[i].GetPageID() == documentID {
			doc1pos = doc1[i].GetWordPositions()
			break
		}
	}

	wid2, err2 := phrases.V.StringToWordID(b.phrase2)
	if err2 != nil {
		return false
	}
	doc2, _ := phrases.ContentInvertedIndexer.GetInvertedFileFromKey(wid2)
	var doc2pos []uint64
	for i := range doc1 {
		if doc2[i].GetPageID() == documentID {
			doc2pos = doc1[i].GetWordPositions()
			break
		}
	}

	for i := range doc2 {
		doc2pos[i]--
	}
	docsIntersect := phrases.Bs.FindIntersect(doc1pos, doc2pos)

	if len(docsIntersect) > 0 {
		return true
	}
	return false
}

// Returns whether the document IDs that has the matching phrase
func (phrases *PhrasalSearch) hasPhrase(b BiGram, query []string) []uint64 {
	docs := phrases.Bs.FindBoolean(query)
	res := make([]uint64, 0)

	for _, did := range docs {
		// bigram in title
		bgTitle := phrases.hasPhraseInTitle(did, b)
		// bigram in body
		bgBody := phrases.hasPhraseInBody(did, b)

		// if present in either title / body, then append
		if bgTitle || bgBody {
			res = append(res, did)
		}

	}
	return res
}

// Returns the document IDs with phrases stated in query
func (phrases *PhrasalSearch) GetPhraseDocuments(query []string) []uint64 {
	if len(query) <= 1 {
		return phrases.Bs.FindBoolean(query)
	}
	ps := phrases.splitToPhrase(query)
	bgDocs := make([][]uint64, 0)
	for _, phr := range ps {
		tempPhr := phrases.hasPhrase(phr, query)
		bgDocs = append(bgDocs, tempPhr)
	}

	sort.Slice(bgDocs, func(i, j int) bool {
		return len(bgDocs[i]) < len(bgDocs[j])
	})

	docs := make([]uint64, 0)
	docs = append(docs, bgDocs[0]...)

	for k, v := range bgDocs {
		if k == 0 {
			continue
		}
		docs = phrases.Bs.FindIntersect(docs, v)
	}

	return docs
}

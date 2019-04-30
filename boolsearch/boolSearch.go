package boolsearch

import (
	"fmt"
	"sort"

	Indexer "github.com/davi1972/comp4321-search-engine/indexer"
	"github.com/davi1972/comp4321-search-engine/vsm"
)

type BoolSearch struct {
	ContentInvertedIndexer *Indexer.InvertedFileIndexer
	Vsm                    *vsm.VSM
}

// Return back an array of doc ids that exist in both arrays
// Intersection algorithm based on https://github.com/juliangruber/go-intersect
func (bs *BoolSearch) FindIntersect(docs1, docs2 []uint64) []uint64 {
	result := make([]uint64, 0)
	for i := 0; i < len(docs1); i++ {
		el := docs1[i]
		idx := sort.Search(len(docs2), func(i int) bool {
			return docs2[i] == el
		})
		if idx < len(docs2) && docs2[idx] == el {
			result = append(result, el)
		}
	}
	return result
}

// Return back an array of page IDs containing query terms
func (bs *BoolSearch) FindBoolean(query []string) []uint64 {
	if len(query) == 0 {
		return nil
	}
	boolMap := make(map[string][]uint64)
	docs := make([]uint64, 0)

	// return page ids with word
	for i := range query {
		wid, err := bs.Vsm.StringToWordID(query[i])
		fmt.Printf("wordID: %d\n", wid)
		if err != nil {
			continue
		}
		invFiles, err := bs.ContentInvertedIndexer.GetInvertedFileFromKey(wid)
		pages := make([]uint64, 0)
		for j := range invFiles {
			pages = append(pages, invFiles[j].GetPageID())
		}
		boolMap[query[i]] = pages
	}

	sort.Slice(query, func(i, j int) bool {
		return len(boolMap[query[i]]) < len(boolMap[query[j]])
	})

	docs = append(docs, boolMap[query[0]]...)
	for k, v := range query {
		if k == 0 {
			continue
		}
		// check for intersection
		docs = bs.FindIntersect(docs, boolMap[v])
	}
	return docs
}

package pageRank

import (
	"fmt"
	"math"

	Indexer "github.com/davi1972/comp4321-search-engine/indexer"
)

type PageRank struct {
	parents                           map[uint64][]uint64
	pageRanks                         map[uint64]float64
	numOutlinks                       map[uint64]int
	documentIndexer                   *Indexer.MappingIndexer
	reverseDocumentIndexer            *Indexer.ReverseMappingIndexer
	parentChildDocumentForwardIndexer *Indexer.ForwardIndexer
	childParentDocumentForwardIndexer *Indexer.ForwardIndexer
	pageRankIndexer                   *Indexer.PageRankIndexer
}

var i = 1

func (pageRank *PageRank) Initialize(mapping *Indexer.MappingIndexer, reverseMapping *Indexer.ReverseMappingIndexer, childParent *Indexer.ForwardIndexer, parentChild *Indexer.ForwardIndexer, page *Indexer.PageRankIndexer) {
	pageRank.parents = make(map[uint64][]uint64)
	pageRank.pageRanks = make(map[uint64]float64)
	pageRank.numOutlinks = make(map[uint64]int)
	pageRank.documentIndexer = mapping
	pageRank.reverseDocumentIndexer = reverseMapping
	pageRank.parentChildDocumentForwardIndexer = parentChild
	pageRank.childParentDocumentForwardIndexer = childParent
	pageRank.pageRankIndexer = page
}

func (pageRank *PageRank) ProcessPageRank() {

	pageIds, err := pageRank.documentIndexer.All()

	pageRank.documentIndexer.Iterate()

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Page IDs")
	fmt.Println(pageIds)

	for _, id := range pageIds {

		_, err := pageRank.reverseDocumentIndexer.GetValueFromKey(id)

		if err != nil {
			fmt.Println(err)
			continue
		}

		pageRank.pageRanks[id] = 1

		inlinks, err := pageRank.childParentDocumentForwardIndexer.GetIdListFromKey(id)

		if len(inlinks) == 0 {
			continue
		}

		if err != nil {
			fmt.Printf("Error getting parent from id: %s", err)
		}

		pageRank.parents[id] = inlinks

		children, _ := pageRank.parentChildDocumentForwardIndexer.GetIdListFromKey(id)

		if len(children) == 0 {
			continue
		} else {
			pageRank.numOutlinks[id] = len(children)
		}

	}

	pageRank.calculatePageRank(0.85, 0.0001)

	for k, v := range pageRank.pageRanks {

		err := pageRank.pageRankIndexer.AddKeyToIndex(k, v)

		if err != nil {
			fmt.Printf("Error inserting pagerank value: %s", err)
		}

	}

}

func (pageRank *PageRank) calculatePageRank(damping float64, threshold float64) {
	fmt.Println("Iteration:", i)
	fmt.Println(pageRank.pageRanks)
	oldRanks := make(map[uint64]float64)

	for key, value := range pageRank.pageRanks {
		oldRanks[key] = value
	}

	stop := true

	for key, value := range oldRanks {

		myParents := pageRank.parents[key]

		var sum = 0.0

		if len(myParents) != 0 {

			for _, id := range myParents {
				parentsChild := float64(pageRank.numOutlinks[id])
				sum += oldRanks[id] / parentsChild
			}

		}

		pr := (1 - damping) + damping*sum

		pageRank.pageRanks[key] = pr

		if math.Abs(pr-value) > threshold {
			stop = false
		}
	}

	i++

	if !stop {
		pageRank.calculatePageRank(damping, threshold)
	} else {
		return
	}

}

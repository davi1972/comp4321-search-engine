package main

import (
	"os"
	"path/filepath"
	"fmt"
	"math"
	Indexer "github.com/davi1972/comp4321-search-engine/indexer"
)
var i = 1
var parents = make(map[uint64][]uint64)

var pageRanks = make(map[uint64]float64)

var numOutlinks = make(map[uint64]int)

func main() {

	wd, _ := os.Getwd()
	parent := filepath.Dir(wd)
	documentIndexer := &Indexer.MappingIndexer{}
	docErr := documentIndexer.Initialize(parent + "/db/documentIndex")
	if docErr != nil {
		fmt.Printf("error when initializing document indexer: %s\n", docErr)
	}
	defer documentIndexer.Backup()
	defer documentIndexer.Release()

	reverseDocumentIndexer := &Indexer.ReverseMappingIndexer{}
	reverseDocumentIndexerErr := reverseDocumentIndexer.Initialize(parent + "/db/reverseDocumentIndexer")
	if reverseDocumentIndexerErr != nil {
		fmt.Printf("error when initializing reverse document indexer: %s\n", reverseDocumentIndexerErr)
	}
	defer reverseDocumentIndexer.Backup()
	defer reverseDocumentIndexer.Release()

	parentChildDocumentForwardIndexer := &Indexer.ForwardIndexer{}
	parentChildDocumentForwardIndexerErr := parentChildDocumentForwardIndexer.Initialize(parent + "/db/parentChildDocumentForwardIndex")
	if parentChildDocumentForwardIndexerErr != nil {
		fmt.Printf("error when initializing parentDocument -> childDocument forward Indexer: %s\n", parentChildDocumentForwardIndexerErr)
	}
	defer parentChildDocumentForwardIndexer.Backup()
	defer parentChildDocumentForwardIndexer.Release()

	childParentDocumentForwardIndexer := &Indexer.ForwardIndexer{}
	childParentDocumentForwardIndexerErr := childParentDocumentForwardIndexer.Initialize(parent + "/db/childParentDocumentForwardIndex")
	if childParentDocumentForwardIndexerErr != nil {
		fmt.Printf("error when initializing childDocument -> parentDocument forward Indexer: %s\n", childParentDocumentForwardIndexerErr)
	}
	defer childParentDocumentForwardIndexer.Backup()
	defer childParentDocumentForwardIndexer.Release()

	pageIds, err := documentIndexer.All()

	if (err != nil){
		fmt.Println(err)
	}

	fmt.Println(pageIds)

	for _, id := range pageIds {

		_, err := reverseDocumentIndexer.GetValueFromKey(id)

		if(err!=nil){
			fmt.Println(err)
			continue
		}

		pageRanks[id] = 1

		inlinks, err := childParentDocumentForwardIndexer.GetIdListFromKey(id)
		
		if(len(inlinks)==0){
			continue
		}

		if (err != nil){
			fmt.Printf("Error getting parent from id: %s", err)
		}

		parents[id] = inlinks

		children, _ := parentChildDocumentForwardIndexer.GetIdListFromKey(id)

		if(len(children)==0){
			continue
		} else {
			numOutlinks[id] = len(children)
		}

	}

	fmt.Println(parents)
	fmt.Println(pageRanks)
	fmt.Println(numOutlinks)

	calculatePageRank(1, 0.001)

}

func calculatePageRank(damping float64, threshold float64){
	fmt.Println("Iteration i:", i)
	fmt.Println(pageRanks)
	oldRanks := make(map[uint64]float64)

	for key, value := range pageRanks {
		oldRanks[key] = value
	}

	stop := true

	for key, value := range oldRanks {

		myParents := parents[key]

		var sum = 0.0

		if(len(myParents)!=0){

			for _, id := range myParents {
				parentsChild := float64(numOutlinks[id])
				sum += oldRanks[id]/parentsChild
			}

		}

		pr := (1-damping) + damping*sum

		pageRanks[key] = pr

		if(math.Abs(pr - value) > threshold){
			stop = false
		}
	}

	i++

	if(!stop){
		calculatePageRank(damping, threshold)
	} else {
		return
	}

}
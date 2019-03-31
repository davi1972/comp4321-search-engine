package main

import (
	"bufio"
	Indexer "github.com/davi1972/comp4321-search-engine/indexer"
	"fmt"
	"os"
	"strconv"
)

func main() {

	wd, _ := os.Getwd()

	pagePropertiesIndexer := &Indexer.PagePropetiesIndexer{}
	pagePropertiesErr := pagePropertiesIndexer.Initialize(wd + "/db/pagePropertiesIndex")
	if pagePropertiesErr != nil {
		fmt.Printf("error when initializing page properties: %s\n", pagePropertiesErr)
	}
	defer pagePropertiesIndexer.Backup()
	defer pagePropertiesIndexer.Release()

	documentWordForwardIndexer := &Indexer.DocumentWordForwardIndexer{}
	documentWordForwardIndexerErr := documentWordForwardIndexer.Initialize(wd + "/db/documentWordForwardIndex")
	if documentWordForwardIndexerErr != nil {
		fmt.Printf("error when initializing document -> word forward Indexer: %s\n", documentWordForwardIndexerErr)
	}
	defer documentWordForwardIndexer.Backup()
	defer documentWordForwardIndexer.Release()

	parentChildDocumentForwardIndexer := &Indexer.ForwardIndexer{}
	parentChildDocumentForwardIndexerErr := parentChildDocumentForwardIndexer.Initialize(wd + "/db/parentChildDocumentForwardIndex")
	if parentChildDocumentForwardIndexerErr != nil {
		fmt.Printf("error when initializing parentDocument -> childDocument forward Indexer: %s\n", parentChildDocumentForwardIndexerErr)
	}
	defer parentChildDocumentForwardIndexer.Backup()
	defer parentChildDocumentForwardIndexer.Release()

	reverseDocumentIndexer := &Indexer.ReverseMappingIndexer{}
	reverseDocumentIndexerErr := reverseDocumentIndexer.Initialize(wd + "/db/reverseDocumentIndexer")
	if reverseDocumentIndexerErr != nil {
		fmt.Printf("error when initializing reverse document indexer: %s\n", reverseDocumentIndexerErr)
	}
	defer reverseDocumentIndexer.Backup()
	defer reverseDocumentIndexer.Release()

	reverseWordindexer := &Indexer.ReverseMappingIndexer{}
	reverseWordindexerErr := reverseWordindexer.Initialize(wd + "/db/reverseWordIndexer")
	if reverseWordindexerErr != nil {
		fmt.Printf("error when initializing reverse word indexer: %s\n", reverseWordindexerErr)
	}
	defer reverseWordindexer.Backup()
	defer reverseWordindexer.Release()

	pages, err := pagePropertiesIndexer.All()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println()

	output := ""

	i := 0

	for _, page := range pages {

		if i>=30 {
			break
		}
		output += strconv.Itoa(i+1)+"\n"
		output += page.GetTitle()
		output += "\n"
		output += page.GetUrl()
		output += "\n"
		output += page.GetDateString() + ", " + strconv.Itoa(page.GetSize()) + " B"
		output += "\n"

		termFreq, _ := documentWordForwardIndexer.GetWordFrequencyListFromKey(page.GetId())

		freqText := ""

		for _, tf := range termFreq {
			word, _ := reverseWordindexer.GetValueFromKey(tf.GetID())
			freqText += word + " " + strconv.FormatUint(tf.GetFrequency(), 10) + ", "
		}

		output += freqText
		output += "\n"
		output += "Children:\n"

		children, _ := parentChildDocumentForwardIndexer.GetIdListFromKey(page.GetId())

		for _, child := range children {

			childUrl, _ := reverseDocumentIndexer.GetValueFromKey(child)
			output += childUrl
			output += "\n"
		}
		output += "------------------------------------------------------------------------\n"

		file, err := os.Create("spider_result.txt")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		w := bufio.NewWriter(file)
		_, err = w.WriteString(output)
		
		i++
	}

}

package tokenizer

import (
	"bufio"
	"os"
	"log"
	"strings"
	"regexp"
	"github.com/reiver/go-porterstemmer"
)

var stopwords string

func LoadStopWords(){

	stopwords = ""

	file, err := os.Open("stopwords.txt")

	if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
		stopwords += scanner.Text()+"|"
	}

	stopwords = stopwords[:len(stopwords)-1]

}

func Tokenize(text string) []string{

	reg, err := regexp.Compile("[^a-zA-Z ]+")
    if err != nil {
        log.Fatal(err)
	}

	regStop, err := regexp.Compile("\\b(" + stopwords + ")\\b")
	if err != nil {
        log.Fatal(err)
	}

	// Remove all special characters
	words := strings.ToLower(reg.ReplaceAllString(text, " "))

	// Remove stopwords
	words = regStop.ReplaceAllString(words, " ")

	tokens := strings.Fields(words)

	// Run Porter's Stemming algorithm
	for i, token:= range tokens {
		tokens[i] = porterstemmer.StemString(token)
	}

	return tokens
}

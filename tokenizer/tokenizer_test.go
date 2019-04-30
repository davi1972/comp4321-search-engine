package tokenizer

import (
	"testing"
)

func TestLoadStopwords(t *testing.T) {
	LoadStopWords()
}

func TestTokenize(t *testing.T) {
	LoadStopWords()
	testStr := "Test tEST arrival"
	result := Tokenize(testStr)
	if result[0] != "test" || result[2] != "arriv" {
		t.Fail()
	}
}

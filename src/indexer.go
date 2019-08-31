package main

import (
	"regexp"
	"strings"
	"io"
	"bytes"
)

// IndexResults defines the output type of the index engine
type IndexResults map[string]int

func (ir IndexResults) Add(word string, value int) {
	ir[strings.ToLower(word)] += value
}

func (ir IndexResults) Get(word string) (int, bool) {
	stat, ok := ir[strings.ToLower(word)]

	return stat, ok
}

// WordIndexer defines an interface by index engine
type WordIndexer interface {
	Index() (IndexResults, error)
}

const delimiter = byte(' ')
const wordsCount = 1

// NewWordIndexer creates a string indexer
func NewWordIndexer(input []byte) WordIndexer {
	gen := NewWordSplitGenerator(bytes.NewBuffer(input), delimiter)
	
	return &wordIndexer{gen: gen}
}

type wordIndexer struct {
	gen BatchesGenerator
}

// Index takes bytes slice and counts the words
func (wi *wordIndexer) Index() (IndexResults, error) {
	statistic := IndexResults{}
	var err error
	var word []byte
	re := regexp.MustCompile(`[A-Za-z0-9]+`)

	for {
		word, err = wi.gen.Next()
		word = re.Find(word)

		if len(word) == 0 {
			break
		}

		statistic.Add(string(word), 1)

		if err != nil {
			break
		}
	}

	if err != io.EOF {
		return statistic, err
	}

	return statistic, nil
}
package main

import (
	"io/ioutil"
	"bytes"
	"io"
	"bufio"
)

// BatchesGenerator defines the interface for text analyzer engine
type BatchesGenerator interface {
	Next() ([]byte, error)
}


// NewSplitGenerator creates new generator that splits input by delimiter
func NewWordSplitGenerator(r io.Reader, delimiter byte) BatchesGenerator {
	return &wordsSplitGenerator{r: bufio.NewReader(r), delimiter: delimiter}
}

func NewBatchGenerator(r io.Reader, delimiter byte, batchSize int64) BatchesGenerator {
	return &bytesBatchGenerator{r: bufio.NewReader(r), delimiter: delimiter, batchSize: batchSize}
}

type wordsSplitGenerator struct {
	r *bufio.Reader
	delimiter byte
}

func (gen *wordsSplitGenerator) Next() (word []byte, err error) {
	for {
		word, err = gen.r.ReadBytes(gen.delimiter)

		if len(word) != 1 || word[0] != gen.delimiter || err != nil {
			break
		}
	}

	if err == nil {
		// Remove the delimiter
		word = word[:len(word)-1]
	}

	return word, err
}


type bytesBatchGenerator struct {
	r *bufio.Reader
	delimiter byte
	batchSize int64
}

func (gen *bytesBatchGenerator) Next() ([]byte, error) {
	var buf []byte
	var err error

	for ; len(buf) == 0; {
		buf = make([]byte, gen.batchSize)
		n, err := gen.r.Read(buf)
		buf = bytes.Trim(buf, "\x00")
		buf = bytes.Trim(buf, string(gen.delimiter))

		if n == 0 || err != nil {
			return buf, err
		}
	}

	if len(buf) == 0 {
		return buf, err
	}

	ending, err := gen.r.ReadBytes(gen.delimiter)

	buf = append(buf, ending...)
	buf = bytes.Trim(buf, string(gen.delimiter))

	return buf, err
}

func uploadBody(body io.ReadCloser) (string, error) {
	tmpfile, err := ioutil.TempFile("", "input")

	if err != nil {
		return tmpfile.Name(), err
	}
	
	io.Copy(tmpfile, body)

	if err := tmpfile.Close(); err != nil {
		return tmpfile.Name(), err
	}

	return tmpfile.Name(), nil
}
package main

import (
	"crypto/tls"
	"net/http"
	"os"
	"log"
	"io"
)

const defaultMaxRetries = 1
const defaultBatchSizeBytes = 100

func NewBatchingJob(input string, inputType InputType) *AsyncJob {
	j := &batchingJob{inputType: inputType, input: input, batchSize: defaultBatchSizeBytes, broker: ChannelBroker}
	return NewAsyncJob(j, defaultMaxRetries, ChannelBroker)
}

func NewReportStatisticJob(statistic IndexResults) *AsyncJob {
	j := &reportStatisticJob{statistic: statistic }
	return NewAsyncJob(j, defaultMaxRetries, ChannelBroker)
}

func NewBatchIndexerJob(batch []byte) *AsyncJob {
	j := &batchIndexerJob{batch: batch, broker: ChannelBroker}
	return NewAsyncJob(j, defaultMaxRetries, ChannelBroker)
}

type InputType string

func (jt InputType) String() string {
	return string(jt)
}

const (
	Text InputType = "text"
	URL InputType = "url"
	File InputType = "file"
)

const wordDelimiter = byte(' ')

type batchingJob struct {
	inputType InputType
	input string
	batchSize int64
	broker Broker
}

func (j *batchingJob) Execute() error {
	log.Println("Got batching task for input type", j.inputType, j.input)
	var reader io.Reader

	switch j.inputType {
	case File:
		f, err := os.Open(j.input)

		if err != nil {
			return err
		}

		reader = f

		defer f.Close()
	case URL:
		transCfg := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
		}
		client := &http.Client{Transport: transCfg}
		resp, err := client.Get(j.input)

		if err != nil {
			return err
		}

		defer resp.Body.Close()
		reader = resp.Body

	default:
		log.Println("Got unknown input type", j.inputType)
		return nil
	}

	return j.execute(reader)
}

func (j *batchingJob) execute(r io.Reader) error {
	gen := NewBatchGenerator(r, wordDelimiter, j.batchSize)
	var err error
	
	for {
		wordsBatch, err := gen.Next()
		
		if len(wordsBatch) == 0 {
			break
		}

		err = NewBatchIndexerJob(wordsBatch).Dispatch()

		if err != nil {
			break
		}
	}

	if err != nil && err != io.EOF {
		log.Println("Batching ended with error", err)
		return err
	}

	log.Println("Batching ended")
	return nil
}

type batchIndexerJob struct {
	batch []byte
	broker Broker
}

func (j *batchIndexerJob) Execute() error {
	log.Println("Got batch to index with len", len(j.batch))

	statistic, err := NewWordIndexer(j.batch).Index()

	if err != nil {
		log.Println("Got to indexer error", err)
		return err
	}


	err = NewReportStatisticJob(statistic).Dispatch()

	return nil
}

type reportStatisticJob struct {
	statistic IndexResults
}

func (j *reportStatisticJob) Execute() error {
	log.Println("got statistic report")

	return StatsStore.UpdateStats(j.statistic)
}
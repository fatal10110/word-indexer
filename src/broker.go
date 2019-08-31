package main

import (
	"reflect"
	"log"
)

// JobExecutor defines the inteface for async job runtime
type JobExecutor interface {
	Execute() error
}

type Broker interface {
	Connect() error
	Disconnect() error
	Handle(func(*AsyncJob))
	Put(*AsyncJob)
}

// NewChannelBroker creats new broker instance base on channels
func NewChannelBroker() Broker {
	return &channelBroker{}
}

// NewAsyncJob wraps the JobExecutor with AsyncJob that can be dispatched
func NewAsyncJob(job JobExecutor, maxRetries int, broker Broker) *AsyncJob {
	return &AsyncJob{Job: job, MaxRetries: maxRetries, broker: broker}
}

type AsyncJob struct {
	Job JobExecutor
	MaxRetries int
	RetriesCount int
	broker Broker
}

func (j *AsyncJob) Dispatch() error {
	go j.broker.Put(j)

	return nil
}

// Retry checks if job retries expired, if not it dispatches itself again
func (j *AsyncJob) Retry() error {
	if j.MaxRetries > j.RetriesCount {
		log.Println("Retring job", j)
		j.RetriesCount++
		return j.Dispatch()
	}

	log.Println("Stop retring", j)

	return nil
}

type channelBroker struct {
	c chan *AsyncJob
	handler func(*AsyncJob)
}

// Handle takes a worker function that should ge an AsyncJob and execute it
func (cb *channelBroker) Handle(handler func(*AsyncJob)) {
	log.Println("Handling started")

	for job := range cb.c {
		log.Println("Executing job", reflect.TypeOf(job.Job))
		handler(job)
	}

	log.Println("Handling ended")
}

func (cb *channelBroker) Put(job *AsyncJob) {
	log.Println("Scheduling new job", reflect.TypeOf(job.Job))

	cb.c <- job
}

func (cb *channelBroker) Connect() error {
	cb.c = make(chan *AsyncJob)

	return nil
}

func (cb *channelBroker) Disconnect() error {
	close(cb.c)

	return nil
}
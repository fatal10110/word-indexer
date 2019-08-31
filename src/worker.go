package main

import (
	"log"
)

func CreateWorker() func(asyncJob *AsyncJob) {
	return func(asyncJob *AsyncJob) {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Trying to recover from error", r)
				asyncJob.Retry()
			}
		}()
	
		if err := asyncJob.Job.Execute(); err != nil {
			log.Println("Job error", asyncJob.Job, err)
			asyncJob.Retry()
		}
	}	
}
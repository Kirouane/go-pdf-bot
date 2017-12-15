package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

type worker struct {
	JobQueue chan job
}

func newWorker(j chan job) worker {
	return worker{
		JobQueue: j,
	}
}

func (w *worker) start() {
	go func() {
		for {
			select {
			case job := <-w.JobQueue:
				fmt.Println("Job received :", job.ID)

				d1 := []byte(job.HTML)
				htmlFile := "/tmp/src.html"
				ioutil.WriteFile(htmlFile, d1, 0644)
				c := NewChrome("file:///"+htmlFile, job.ID)
				pdf := c.run()
				log.Println(job.Webhook)
				if job.Webhook != "" {
					w := webhook{}
					w.push(pdf, job)
				}
			}
		}
	}()
}

func (w *worker) addJob(j job) {
	w.JobQueue <- j
	fmt.Println("Job sent :", j.ID)
}

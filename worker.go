package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
				dir, err := os.Getwd()
				if err != nil {
					log.Fatal(err)
				}
				htmlFile := dir + "/src.html"
				err = ioutil.WriteFile(htmlFile, d1, 0644)
				if err != nil {
					log.Fatal(err)
				}
				headless := NewHeadless()
				pdf := headless.PrintPdf(job.ID, "file:///"+htmlFile)
				headless.cancel()
				err = ioutil.WriteFile(dir+"/storage/pdf/"+pdf.Filename, pdf.Content, 0644)

				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}()
}

func (w *worker) addJob(j job) {
	w.JobQueue <- j
	fmt.Println("Job sent :", j.ID)
}

package main

import (
	"time"

	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/satori/go.uuid"
)

/**
 * Model
 */
type job struct {
	ID      string
	Date    time.Time
	HTML    string
	Webhook string
}

/**
 * Controler
 */
type jobCreateController struct {
	Worker worker
}

func (controller jobCreateController) action(params map[string]string) map[string]string {
	storage := &storage{}
	storage.connect()

	j := &job{
		ID:      uuid.NewV4().String(),
		Date:    time.Now(),
		HTML:    params["html"],
		Webhook: params["webhook"],
	}

	controller.Worker.addJob(*j)

	//storage.writeJob(j)
	return map[string]string{}
}

/**
 * Storage
 */
type storage struct {
	db *scribble.Driver
}

func (s *storage) connect() {
	db, err := scribble.New("storage/", nil)
	if err != nil {
		panic(err)
	}
	s.db = db
}

func (s *storage) writeJob(j *job) {
	err := s.db.Write("job", j.ID, j)

	if err != nil {
		panic(err)
	}
}

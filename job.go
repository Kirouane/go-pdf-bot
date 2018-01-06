package main

import (
	"time"

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
	id, _ := uuid.NewV4()
	j := &job{
		ID:      id.String(),
		Date:    time.Now(),
		HTML:    params["html"],
		Webhook: params["webhook"],
	}

	controller.Worker.addJob(*j)
	return map[string]string{}
}

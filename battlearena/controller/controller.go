package controller

import "github.com/google/uuid"

type Controller struct {
	Uuid           uuid.UUID
	Assigned       bool
	ControllerName string
}

// New
func New() *Controller {
	return &Controller{
		Uuid:     uuid.New(),
		Assigned: false,
	}
}

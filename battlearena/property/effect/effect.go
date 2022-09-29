package effect

import "github.com/ecumeurs/upsilonbattle/battlearena/property"

type Effect struct {
	Properties []property.Property
	Name       string
}

// New
func New() *Effect {
	return &Effect{
		Properties: []property.Property{},
		Name:       "New Effect",
	}
}

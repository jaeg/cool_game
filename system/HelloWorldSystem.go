package system

import (
	"log"

	"github.com/jaeg/cool_game/component"
)

// HelloWorldSystem .
type HelloWorldSystem struct {
}

// Update .
func (HelloWorldSystem) Update(a *component.HelloWorldComponent) {
	log.Println("hello world")
}

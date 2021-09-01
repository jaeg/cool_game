package system

import (
	"log"

	"github.com/jaeg/cool_game/components"
)

// HelloWorldSystem .
type HelloWorldSystem struct {
}

// Update .
func (HelloWorldSystem) Update(a *components.HelloWorldComponent) {
	log.Println("hello world")
}

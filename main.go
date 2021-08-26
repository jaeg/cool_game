package main

import (
	"log"

	"github.com/jaeg/cool_game/entity"
	"github.com/jaeg/cool_game/game"
)

const (
	screenWidth  = 800
	screenHeight = 640
)

func main() {
	entity.FactoryLoad("entities.blueprints")
	g, err := game.NewGame("Cool Game", screenWidth, screenHeight)
	if err != nil {
		log.Fatal(err)
	}

	err = g.Run()

	if err != nil {
		log.Fatal(err)
	}
}

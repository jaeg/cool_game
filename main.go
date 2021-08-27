package main

import (
	"log"

	"github.com/jaeg/cool_game/entity"
	"github.com/jaeg/cool_game/game"
)

func main() {
	entity.FactoryLoad("entities.blueprints")
	g, err := game.NewGame("Cool Game")
	if err != nil {
		log.Fatal(err)
	}

	err = g.Run()

	if err != nil {
		log.Fatal(err)
	}
}

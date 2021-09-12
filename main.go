package main

import (
	"log"

	"github.com/jaeg/cool_game/factory"
	"github.com/jaeg/cool_game/game"
	"github.com/jaeg/game-engine/entity"
)

func main() {
	factory.FactoryLoad("entities.blueprints")
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

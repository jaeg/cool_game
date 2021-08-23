package main

import (
	"log"

	"github.com/jaeg/cool_game/game"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

func main() {
	g, err := game.NewGame("Cool Game", screenWidth, screenHeight)
	if err != nil {
		log.Fatal(err)
	}

	err = g.Run()
	if err != nil {
		log.Fatal(err)
	}
}

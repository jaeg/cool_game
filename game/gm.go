package game

import (
	"log"
	"math/rand"

	"github.com/jaeg/cool_game/config"
	"github.com/jaeg/cool_game/entity"
	"github.com/jaeg/cool_game/world"
)

const foodMinimum = 100
const hostileMinimum = 100
const foodInitial = 200
const hostileInitial = 200

var hostiles = []string{"snake", "skeleton", "spider", "zombie"}
var rareHostiles = []string{"giant", "centaur", "griffon"}
var epicHostiles = []string{"darkknight"}
var foods = []string{"rat", "dog", "cat", "roach"}

//GameMaster The Game master manages what's going on in the game world and balances the difficulty.
type GameMaster struct {
	level *world.Level
}

//Init Initial the game master
func (gm *GameMaster) Init(level *world.Level) {
	gm.level = level

	log.Println("Placing food")
	//Random food
	for i := 0; i < foodInitial; i++ {
		x := rand.Intn(config.WorldGenSizeW)
		y := rand.Intn(config.WorldGenSizeH)
		tile := level.GetTileAt(x, y)
		tries := 0
		for tile.Type == 2 || tile.Type == 4 || level.GetEntityAt(x, y) != nil {
			x = rand.Intn(config.WorldGenSizeW)
			y = rand.Intn(config.WorldGenSizeH)
			tile = level.GetTileAt(x, y)
			tries++
			if tries > 10 {
				break
			}
		}
		if tries > 10 {
			continue
		}
		blueprint := foods[getRandom(0, len(foods))]
		food, err := entity.Create(blueprint, x, y)
		if err == nil {
			level.AddEntity(food)
		}
	}

	log.Println("Placing hostiles")
	//Random hostiles
	for i := 0; i < hostileInitial; i++ {
		x := rand.Intn(config.WorldGenSizeW)
		y := rand.Intn(config.WorldGenSizeH)
		tile := level.GetTileAt(x, y)
		tries := 0
		for tile.Type == 2 || tile.Type == 4 || level.GetEntityAt(x, y) != nil {
			x = rand.Intn(config.WorldGenSizeW)
			y = rand.Intn(config.WorldGenSizeH)
			tile = level.GetTileAt(x, y)
			tries++
			if tries > 10 {
				break
			}
		}
		if tries > 10 {
			continue
		}
		blueprint := hostiles[getRandom(0, len(hostiles))]
		if getRandom(0, 100) == 0 {
			log.Println("Spawn a rare hostile enemy!")
			blueprint = rareHostiles[getRandom(0, len(rareHostiles))]
		}
		food, err := entity.Create(blueprint, x, y)
		if err == nil {
			level.AddEntity(food)
		}
	}
}

//Update Update the game master
func (gm *GameMaster) Update() {
	foodCount := 0
	hostileCount := 0
	dragonPresent := false

	//Gather stats
	for _, e := range gm.level.Entities {
		if e.HasComponent("FoodComponent") {
			foodCount++
		}

		if e.HasComponent("HostileAIComponent") {
			hostileCount++
		}

		if e.Blueprint == "dragon" {
			dragonPresent = true
		}
	}

	// Handle food count
	if foodCount < foodMinimum {
		log.Println("Below minimum number of food entities... Placing food")
		//Random food
		for i := 0; i < foodMinimum-foodCount; i++ {
			x := rand.Intn(config.WorldGenSizeW)
			y := rand.Intn(config.WorldGenSizeH)
			tile := gm.level.GetTileAt(x, y)
			tries := 0
			for tile.Type == 2 || tile.Type == 4 || gm.level.GetEntityAt(x, y) != nil {
				x = rand.Intn(config.WorldGenSizeW)
				y = rand.Intn(config.WorldGenSizeH)
				tile = gm.level.GetTileAt(x, y)
				tries++
				if tries > 10 {
					break
				}
			}
			if tries > 10 {
				continue
			}
			blueprint := foods[getRandom(0, len(foods))]
			food, err := entity.Create(blueprint, x, y)
			if err == nil {
				gm.level.AddEntity(food)
			}
		}
	}

	// Handle hostile count
	if hostileCount < hostileMinimum {
		log.Println("Below minimum number of hostile entities... Placing hostiles")
		//Random snakes
		for i := 0; i < hostileInitial-hostileCount; i++ {
			x := rand.Intn(config.WorldGenSizeW)
			y := rand.Intn(config.WorldGenSizeH)
			tile := gm.level.GetTileAt(x, y)
			tries := 0
			for tile.Type == 2 || tile.Type == 4 || gm.level.GetEntityAt(x, y) != nil {
				x = rand.Intn(config.WorldGenSizeW)
				y = rand.Intn(config.WorldGenSizeH)
				tile = gm.level.GetTileAt(x, y)
				tries++
				if tries > 10 {
					break
				}
			}
			if tries > 10 {
				continue
			}

			blueprint := hostiles[getRandom(0, len(hostiles))]

			if getRandom(0, 500) == 0 {
				log.Println("Spawn a rare hostile enemy!")
				blueprint = rareHostiles[getRandom(0, len(rareHostiles))]
				if getRandom(0, 500) == 0 {
					log.Println("!!Spawned an epic hostile enemy instead!!")
					blueprint = epicHostiles[getRandom(0, len(epicHostiles))]
				}
			}
			food, err := entity.Create(blueprint, x, y)
			if err == nil {
				gm.level.AddEntity(food)
			}
		}
	}

	// DRAGONS!!
	if dragonPresent {
		if getRandom(0, 1000) == 0 {
			log.Println("A dragon has flown in!")
			x := rand.Intn(config.WorldGenSizeW)
			y := rand.Intn(config.WorldGenSizeH)
			tile := gm.level.GetTileAt(x, y)
			tries := 0
			for tile.Type == 2 || tile.Type == 4 || gm.level.GetEntityAt(x, y) != nil {
				x = rand.Intn(config.WorldGenSizeW)
				y = rand.Intn(config.WorldGenSizeH)
				tile = gm.level.GetTileAt(x, y)
				tries++
				if tries > 10 {
					break
				}
			}
			if tries < 10 {
				dragon, err := entity.Create("dragon", x, y)
				if err == nil {
					gm.level.AddEntity(dragon)
				}
			}
		}
	}
}

func getRandom(low int, high int) int {
	if low == high {
		return low
	}
	return (rand.Intn((high - low))) + low
}

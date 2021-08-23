package world

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/jaeg/cool_game/component"
	"github.com/jaeg/cool_game/entity"

	"github.com/aquilax/go-perlin"
	"github.com/beefsack/go-astar"
)

const (
	alpha = 6.
	beta  = 5.
	n     = 2
)

// Level .
type Level struct {
	data                  [][]Tile
	Entities              []*entity.Entity
	Width, Height         int
	id                    int
	left, right, up, down int
	theme                 string
	Hour                  int
}

//Tile Types
//1 - open
//2 - solid
//3 - stairs [level id, to x, to y]
//4 - water
type Tile struct {
	Type      int
	SpriteX   int
	SpriteY   int
	Data      []int
	X         int
	Y         int
	Elevation int
	level     *Level
	Entities  []*entity.Entity
}

func (t *Tile) PathNeighbors() []astar.Pather {
	neighbors := []astar.Pather{}
	for _, offset := range [][]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
	} {
		if n := t.level.GetTileAt(t.X+offset[0], t.Y+offset[1]); n != nil &&
			n.Type != 2 {
			neighbors = append(neighbors, n)
		}
	}
	return neighbors
}

func (t *Tile) PathNeighborCost(to astar.Pather) float64 {
	toTile, ok := to.(*Tile)
	if !ok {
		return 100
	}
	if toTile == nil {
		return 100
	}
	if t == nil {
		return 100
	}
	cost := math.Abs(float64(t.Elevation - toTile.Elevation))
	if toTile.Type == 2 {
		cost = 100
	}

	if toTile.level.GetSolidEntityAt(toTile.X, toTile.Y) != nil {
		cost = 100
	}
	return cost
}

func (t *Tile) PathEstimatedCost(to astar.Pather) float64 {
	if to == nil {
		return -1
	}
	t2, ok := to.(*Tile)
	if !ok {
		return -1
	}
	if t2 == nil {
		return -1
	}

	if t == nil {
		return -1
	}

	first := math.Pow(float64(t2.X-t.X), 2)
	second := math.Pow(float64(t2.Y-t.Y), 2)
	return first + second
}

func newLevel(width int, height int) (level *Level) {
	level = &Level{left: -1, right: -1, up: -1, down: -1, Width: width, Height: height, Hour: 9}

	data := make([][]Tile, width, height)
	for x := 0; x < width; x++ {
		col := []Tile{}
		for y := 0; y < height; y++ {
			col = append(col, Tile{Type: 4, X: x, Y: y, SpriteX: 16, SpriteY: 128, level: level})
		}
		data[x] = append(data[x], col...)
	}

	level.data = data
	return
}

func NewOverworldSection(width int, height int) (level *Level) {
	log.Println("Creating new random level")
	level = newLevel(width, height)

	p := perlin.NewPerlin(alpha, beta, n, time.Now().UnixNano())
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			tile := level.GetTileAt(x, y)
			value := int(p.Noise2D(float64(x)/100, float64(y)/100) * 10)
			tile.Elevation = value

			//Beach
			if value == -2 {
				tile.Type = 1
				tile.SpriteX = 112
				tile.SpriteY = 112
				rn := getRandom(0, 4)
				if rn == 1 {
					tile.SpriteX = 128
				} else if rn == 2 {
					tile.SpriteX = 144
				} else if rn == 3 {
					tile.SpriteX = 160
				}
			}

			//grass
			if value > -2 {

				tile.Type = 1
				tile.SpriteX = 128
				tile.SpriteY = 80
				if getRandom(0, 3) == 2 {
					tile.SpriteX = 144
				}

			}

			//Mountain
			if value >= 2 {

				tile.Type = 2
				tile.SpriteX = 0
				tile.SpriteY = 80
				rn := getRandom(0, 4)
				if rn == 1 {
					tile.SpriteX = 16
				} else if rn == 2 {
					tile.SpriteX = 32
				} else if rn == 3 {
					tile.SpriteX = 48
				}
			}
		}
	}
	//Generate flower formations
	for i := 0; i < 50; i++ {
		x := getRandom(1, width)
		y := getRandom(1, height)

		level.createCluster(x, y, 10, 176, 80, false, false)
	}

	//Generate tree formations
	for i := 0; i < 50; i++ {
		x := getRandom(1, width)
		y := getRandom(1, height)

		level.CreateClusterOfTrees(x, y, 10)
	}

	//Generate ponds
	for i := 0; i < 10; i++ {
		x := getRandom(1, width)
		y := getRandom(1, height)

		level.createCluster(x, y, getRandom(20, 300), 16, 128, false, true)
	}

	// //Generate Mountains
	// for i := 0; i < 3; i++ {
	// 	x := getRandom(1, width)
	// 	y := getRandom(1, height)
	// 	rn := getRandom(0, 4)
	// 	sX := 0
	// 	if rn == 1 {
	// 		sX = 16
	// 	} else if rn == 2 {
	// 		sX = 32
	// 	} else if rn == 3 {
	// 		sX = 48
	// 	}
	// 	level.createCluster(x, y, 2000, int32(sX), 80, true, false)
	// }

	return
}

func (level *Level) GetTileAt(x int, y int) (tile *Tile) {
	if x < level.Width && y < level.Height && x >= 0 && y >= 0 {
		tile = &level.data[x][y]
	}
	return
}

func (level *Level) PlaceEntity(x int, y int, entity *entity.Entity) {
	if x < level.Width && y < level.Height && x >= 0 && y >= 0 {
		tile := &level.data[x][y]
		pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
		oldTile := &level.data[pc.GetX()][pc.GetY()]
		for i := 0; i < len(oldTile.Entities); i++ {
			if oldTile.Entities[i] == entity {
				oldTile.Entities = append(oldTile.Entities[:i], oldTile.Entities[i+1:]...)
			}
		}
		tile.Entities = append(tile.Entities, entity)
		pc.SetPosition(x, y)
	}
}

//Get's the view frustum with the player in the center
func (level *Level) GetView(aX int, aY int, width int, height int, blind bool, centered bool) (data [][]*Tile) {
	left := aX - width/2
	right := aX + width/2
	up := aY - height/2
	down := aY + height/2

	if centered == false {
		left = aX
		right = aX + width - 1
		up = aY
		down = aY + height
	}

	data = make([][]*Tile, width+1-width%2)

	cX := 0
	for x := left; x <= right; x++ {
		col := []*Tile{}
		for y := up; y <= down; y++ {
			currentTile := level.GetTileAt(x, y)
			if blind {
				if y < aY-height/4 || y > aY+height/4 || x > aX+width/4 || x < aX-width/4 {
					currentTile = nil
				}
			}
			col = append(col, currentTile)
		}
		data[cX] = append(data[cX], col...)
		cX++
	}
	return
}

func (level *Level) GetEntitiesAround(x int, y int, width int, height int) (entities []*entity.Entity) {
	left := x - width/2
	right := x + width/2
	up := y - height/2
	down := y + height/2

	for x := left; x < right; x++ {
		for y := up; y < down; y++ {
			tile := level.GetTileAt(x, y)
			if tile != nil {
				if len(tile.Entities) > 0 {
					entity := tile.Entities[0]

					if entity.HasComponent("PositionComponent") {
						pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
						if pc.GetX() >= left && pc.GetX() <= right && pc.GetY() >= up && pc.GetY() <= down {
							entities = append(entities, entity)
						}
					}
				}
			}
		}
	}
	return
}

func (level *Level) GetPlayersAround(x int, y int, width int, height int) (entities []*entity.Entity) {
	left := x - width/2
	right := x + width/2
	up := y - height/2
	down := y + height/2

	entitiesLen := len(level.Entities)

	for i := 0; i < entitiesLen; i++ {
		entity := level.Entities[i]

		if entity.HasComponent("PositionComponent") {
			pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
			if pc.GetX() >= left && pc.GetX() <= right && pc.GetY() >= up && pc.GetY() <= down {
				if entity.HasComponent("PlayerComponent") {
					entities = append(entities, entity)
				}
			}
		}
	}

	return
}

func (level *Level) GetEntityAt(x int, y int) (entity *entity.Entity) {
	if x < level.Width && y < level.Height && x >= 0 && y >= 0 {
		tile := &level.data[x][y]
		if len(tile.Entities) > 0 {
			return tile.Entities[0]
		}
	}

	entity = nil
	return
}

func (level *Level) GetSolidEntityAt(x int, y int) (entity *entity.Entity) {
	if x < level.Width && y < level.Height && x >= 0 && y >= 0 {
		tile := &level.data[x][y]
		if len(tile.Entities) > 0 {
			if tile.Entities[0].HasComponent("SolidComponent") {

				return tile.Entities[0]
			}
		}
	}

	entity = nil
	return
}

func (level *Level) GetInteractableEntityAt(x int, y int) (entity *entity.Entity) {
	for i := 0; i < len(level.Entities); i++ {
		entity = level.Entities[i]
		if entity.HasComponent("PositionComponent") {
			if entity.HasComponent("InteractComponent") {
				pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
				if pc.GetX() == x && pc.GetY() == y {
					return
				}
			}
		}
	}
	entity = nil
	return
}

func (level *Level) NextHour() {
	level.Hour++
	if level.Hour >= 24 {
		level.Hour = 0
	}
}

func (level *Level) AddEntity(entity *entity.Entity) {
	level.Entities = append(level.Entities, entity)
	if entity.HasComponent("PositionComponent") {
		pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
		x := pc.GetX()
		y := pc.GetY()
		level.PlaceEntity(x, y, entity)
	}
}

func (level *Level) RemoveEntity(entity *entity.Entity) {
	if entity.HasComponent("PositionComponent") {
		pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
		x := pc.GetX()
		y := pc.GetY()

		if x < level.Width && y < level.Height && x >= 0 && y >= 0 {
			tile := &level.data[pc.GetX()][pc.GetY()]
			for i := 0; i < len(tile.Entities); i++ {
				if tile.Entities[i] == entity {
					tile.Entities = append(tile.Entities[:i], tile.Entities[i+1:]...)
				}
			}
		}
	}
	for i := 0; i < len(level.Entities); i++ {
		if level.Entities[i] == entity {
			level.Entities = append(level.Entities[:i], level.Entities[i+1:]...)
			return
		}
	}
}

func (level *Level) createCluster(x int, y int, size int, spriteX int, spriteY int, solid bool, water bool) {
	for i := 0; i < size; i++ {
		n := getRandom(1, 6)
		e := getRandom(1, 6)
		s := getRandom(1, 6)
		w := getRandom(1, 6)

		if n == 1 {
			x += 1
		}

		if s == 1 {
			x--
		}

		if e == 1 {
			y++
		}

		if w == 1 {
			y--
		}

		if level.GetTileAt(x, y) != nil {
			tile := level.GetTileAt(x, y)
			tile.SpriteX = spriteX
			tile.SpriteY = spriteY
			if solid {
				tile.Type = 2
			} else if water {
				tile.Type = 4
			} else {
				tile.Type = 1
			}

		}
	}
}

func (level *Level) CreateClusterOfTrees(x int, y int, size int) {
	for i := 0; i < size; i++ {
		n := getRandom(1, 6)
		e := getRandom(1, 6)
		s := getRandom(1, 6)
		w := getRandom(1, 6)

		if n == 1 {
			x += 1
		}

		if s == 1 {
			x--
		}

		if e == 1 {
			y++
		}

		if w == 1 {
			y--
		}

		tries := 0
		if level.GetTileAt(x, y) != nil {
			tile := level.GetTileAt(x, y)
			if tile.Type != 2 && tile.Type != 4 && level.GetEntityAt(x, y) == nil {
				tree, err := entity.Create("tree", x, y)
				if err == nil {
					level.AddEntity(tree)
				}
			} else {
				i--
				tries++
			}
			if tries > 10 {
				continue
			}
		}
	}
}

func (level *Level) CreateClusterOfGoblins(x int, y int, size int) {
	for i := 0; i < size; i++ {
		n := getRandom(1, 6)
		e := getRandom(1, 6)
		s := getRandom(1, 6)
		w := getRandom(1, 6)

		if n == 1 {
			x += 1
		}

		if s == 1 {
			x--
		}

		if e == 1 {
			y++
		}

		if w == 1 {
			y--
		}

		tries := 0
		if level.GetTileAt(x, y) != nil {
			tile := level.GetTileAt(x, y)
			if tile.Type != 2 && tile.Type != 4 && level.GetEntityAt(x, y) == nil {
				goblin, err := entity.Create("goblin", x, y)
				if err == nil {
					level.AddEntity(goblin)
				}
			} else {
				i--
				tries++
			}
			if tries > 10 {
				continue
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

func Sgn(a int) int {
	switch {
	case a < 0:
		return -1
	case a > 0:
		return +1
	}
	return 0
}

//Ported from http://www.roguebasin.com/index.php?title=Simple_Line_of_Sight
func los(pX int, pY int, tX int, tY int, level *Level) bool {
	deltaX := pX - tX
	deltaY := pY - tY

	absDeltaX := math.Abs(float64(deltaX))
	absDeltaY := math.Abs(float64(deltaY))

	signX := Sgn(deltaX)
	signY := Sgn(deltaY)

	if absDeltaX > absDeltaY {
		t := absDeltaY*2 - absDeltaX
		for {
			if t >= 0 {
				tY += signY
				t -= absDeltaX * 2
			}

			tX += signX
			t += absDeltaY * 2

			if tX == pX && tY == pY {
				return true
			}
			if level.GetTileAt(tX, tY).Type == 2 {
				break
			}
		}
		return false
	}

	t := absDeltaX*2 - absDeltaY

	for {
		if t >= 0 {
			tX += signX
			t -= absDeltaY * 2
		}
		tY += signY
		t += absDeltaX * 2
		if tX == pX && tY == pY {
			return true
		}

		if level.GetTileAt(tX, tY).Type == 2 {
			break
		}
	}

	return false

}

func distance(x1 int, y1 int, x2 int, y2 int) int {
	var dy int
	if y1 > y2 {
		dy = y1 - y2
	} else {
		dy = (y2 - y1)
	}

	var dx int
	if x1 > x2 {
		dx = x1 - x2
	} else {
		dx = x2 - x1
	}

	var d int
	if dy > dx {
		d = dy + (dx >> 1)
	} else {
		d = dx + (dy >> 1)
	}

	return d
}

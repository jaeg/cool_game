package game

import (
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jaeg/cool_game/component"
	"github.com/jaeg/cool_game/config"
	"github.com/jaeg/cool_game/entity"
	"github.com/jaeg/cool_game/system"
	"github.com/jaeg/cool_game/world"
)

type Game struct {
	title            string
	worldTileset     *ebiten.Image
	characterTileset *ebiten.Image
	uiTileset        *ebiten.Image
	fxTileset        *ebiten.Image
	minimap          *ebiten.Image
	level            *world.Level
	Width            int
	Height           int
	CameraX          int
	CameraY          int
	keys             []ebiten.Key
	Cursor           *Cursor
	systems          []system.System
	gm               *GameMaster
}

type Cursor struct {
	X     int
	Y     int
	State int
}

func (c *Cursor) Update() {
	c.State = 0
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		c.State = 1
	}

	c.X, c.Y = ebiten.CursorPosition()
}

func NewGame(title string, width int, height int) (*Game, error) {
	game := &Game{Width: width, Height: height, title: title, Cursor: &Cursor{}}
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle(title)

	//Load assets
	img, err := LoadImage("assets/tiny_dungeon_world.png")
	if err != nil {
		return nil, err
	}
	game.worldTileset = img

	img, err = LoadImage("assets/tiny_dungeon_monsters.png")
	if err != nil {
		return nil, err
	}
	game.characterTileset = img

	img, err = LoadImage("assets/tiny_dungeon_interface.png")
	if err != nil {
		return nil, err
	}
	game.uiTileset = img

	img, err = LoadImage("assets/tiny_dungeon_fx.png")
	if err != nil {
		return nil, err
	}
	game.fxTileset = img

	game.level = world.NewOverworldSection(config.WorldGenSizeW, config.WorldGenSizeH)

	game.systems = make([]system.System, 0)

	//Initiative System
	game.systems = append(game.systems, system.InitiativeSystem{})

	//AI System
	game.systems = append(game.systems, system.AISystem{})

	//StatusCondition System
	game.systems = append(game.systems, system.StatusConditionSystem{})

	game.gm = &GameMaster{}
	game.gm.Init(game.level)

	return game, nil
}

func (g *Game) Run() error {
	err := ebiten.RunGame(g)
	return err
}

func (g *Game) Update() error {
	g.keys = inpututil.PressedKeys()
	for _, k := range g.keys {
		if k.String() == "W" {
			if g.CameraY > 0 {
				g.CameraY--
			}
		}
		if k.String() == "S" {
			if (g.CameraY + g.Height/config.TileSizeH) < config.WorldGenSizeH {
				g.CameraY++
			}
		}
		if k.String() == "A" {
			if g.CameraX > 0 {
				g.CameraX--
			}
		}
		if k.String() == "D" {
			if (g.CameraX + g.Width/config.TileSizeW) < config.WorldGenSizeW {
				g.CameraX++
			}
		}
	}

	g.Cursor.Update()

	fps := ebiten.CurrentFPS()
	ebiten.SetWindowTitle(g.title + " FPS: " + strconv.FormatFloat(fps, 'f', 1, 64))
	g.gm.Update()
	for _, entity := range g.level.Entities {
		for s := range g.systems {
			g.systems[s].Update(g.level, entity)
		}
	}
	cs := system.CleanUpSystem{}
	cs.Update(g.level)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	//Draw world
	worldImage := ebiten.NewImage(config.World_W, config.World_H)
	g.DrawLevel(worldImage, g.CameraX, g.CameraY, config.World_W/config.TileSizeW, config.World_H/config.TileSizeH, false, false)

	/*view := g.level.GetView(g.CameraX, g.CameraY, config.World_W/config.TileSizeW, config.World_H/config.TileSizeH, false, false)
	for x := 0; x < len(view); x++ {
		for y := 0; y < len(view[x]); y++ {
			tX := float64(x * config.SpriteSizeW)
			tY := float64(y * config.SpriteSizeH)
			tile := view[x][y]

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(tX, tY)
			op.GeoM.Scale(float64(config.TileSizeW/config.SpriteSizeW), float64(config.TileSizeH/config.SpriteSizeH))

			if tile == nil {
				worldImage.DrawImage(g.worldTileset.SubImage(image.Rect(0, 112, config.SpriteSizeW, 112+config.SpriteSizeH)).(*ebiten.Image), op)
				continue
			} else {
				worldImage.DrawImage(g.worldTileset.SubImage(image.Rect(tile.SpriteX, tile.SpriteY, tile.SpriteX+config.SpriteSizeW, tile.SpriteY+config.SpriteSizeH)).(*ebiten.Image), op)
			}

			//Draw entity on tile.
			entity := g.level.GetEntityAt(tile.X, tile.Y)
			if entity != nil {
				if entity.HasComponent("AppearanceComponent") {
					ac := entity.GetComponent("AppearanceComponent").(*component.AppearanceComponent)
					dir := 0
					if entity.HasComponent("DirectionComponent") {
						dc := entity.GetComponent("DirectionComponent").(*component.DirectionComponent)
						dir = dc.Direction
					}

					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(tX, tY)
					op.GeoM.Scale(float64(config.TileSizeW/config.SpriteSizeW), float64(config.TileSizeH/config.SpriteSizeH))

					if entity.HasComponent("InanimateComponent") {
						worldImage.DrawImage(g.worldTileset.SubImage(image.Rect(ac.SpriteX+dir*config.SpriteSizeW, ac.SpriteY, ac.SpriteX+config.SpriteSizeW+dir*config.SpriteSizeW, ac.SpriteY+config.SpriteSizeH)).(*ebiten.Image), op)
					} else {
						worldImage.DrawImage(g.characterTileset.SubImage(image.Rect(ac.SpriteX+dir*config.SpriteSizeW, ac.SpriteY, ac.SpriteX+config.SpriteSizeW+dir*config.SpriteSizeW, ac.SpriteY+config.SpriteSizeH)).(*ebiten.Image), op)
					}
				}
			}
		}
	}*/
	screen.DrawImage(worldImage, nil)
	g.DrawUI(screen)

	g.DrawCursor(screen)

	//Draw Minimap
	op := &ebiten.DrawImageOptions{}
	//op.GeoM.Scale(.2, .2)
	op.GeoM.Translate(config.World_W+5, 16)
	if g.minimap == nil {
		g.minimap = g.GetMinimap(0, 0, config.WorldGenSizeW, config.WorldGenSizeH, 150, 150)
	}
	//g.minimap = g.GetMinimap(g.CameraX, g.CameraY, 300, 300, 150, 150)
	screen.DrawImage(g.minimap, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.Width, g.Height
}

func LoadImage(path string) (*ebiten.Image, error) {
	imgFile, err := ebitenutil.OpenFile(path)
	if err != nil {
		fmt.Println("Error opening tileset " + path)
		return nil, errors.New("error opening tileset " + path)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}

func (g *Game) DrawUI(screen *ebiten.Image) {
	//Draw menu
	for x := config.World_W; x < g.Width; x += 16 {
		for y := 0; y < g.Height; y += 16 {
			sX := 127
			sY := 16
			//Left Top
			if x == config.World_W && y == 0 {
				sY = 0
				sX = 144
			} else if x == g.Width-16 && y == 0 { //Right top
				sY = 0
				sX = 176
			} else if x == config.World_W && y == g.Height-16 { //Left bottom
				sY = 32
				sX = 144
			} else if x == g.Width-16 && y == g.Height-16 { //Right bottom
				sY = 32
				sX = 176
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			//s.drawSpriteEx(int32(x), int32(y), sX, sY, 32, 32, 255, 255, 255, 255, s.uiTexture)
			screen.DrawImage(g.uiTileset.SubImage(image.Rect(sX, sY, sX+config.SpriteSizeW, sY+config.SpriteSizeH)).(*ebiten.Image), op)

		}
	}
}

//GetMinimap
//Width and Height is in number of tiles.
func (g *Game) GetMinimap(sX int, sY int, width int, height int, imageWidth int, imageHeight int) *ebiten.Image {
	worldImage := ebiten.NewImage(imageWidth, imageHeight)

	view := g.level.GetView(sX, sY, width, height, false, false)
	for x := 0; x < len(view); x++ {
		for y := 0; y < len(view[x]); y++ {
			tX := float64(x * imageWidth / width)
			tY := float64(y * imageHeight / height)
			tile := view[x][y]

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(tX, tY)
			//op.GeoM.Scale(float64(config.TileSizeW/config.SpriteSizeW), float64(config.TileSizeH/config.SpriteSizeH))

			if tile == nil {
				worldImage.DrawImage(g.worldTileset.SubImage(image.Rect(0, 112, config.SpriteSizeW, 112+config.SpriteSizeH)).(*ebiten.Image), op)
				continue
			} else {
				worldImage.DrawImage(g.worldTileset.SubImage(image.Rect(tile.SpriteX, tile.SpriteY, tile.SpriteX+config.SpriteSizeW, tile.SpriteY+config.SpriteSizeH)).(*ebiten.Image), op)
			}
		}
	}

	return worldImage
}

func (g *Game) DrawCursor(screen *ebiten.Image) {
	//Cursor logic

	//pX := g.Cursor.X/config.TileSizeW + g.CameraX
	//pY := g.Cursor.Y/config.TileSizeW + g.CameraY
	var cursorY = 128
	if g.Cursor.State == 1 {
		cursorY = 144
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(config.TileSizeW/config.SpriteSizeW), float64(config.TileSizeH/config.SpriteSizeH))

	if g.Cursor.X > config.World_W {
		op.GeoM.Translate(float64(g.Cursor.X), float64(g.Cursor.Y))
		screen.DrawImage(g.uiTileset.SubImage(image.Rect(64, cursorY, 64+config.SpriteSizeW, cursorY+config.SpriteSizeH)).(*ebiten.Image), op)
		//s.drawSprite(int32(g.Cursor.X), int32(g.Cursor.Y), 64, cursorY, 255, 255, 255, g.uiTexture) //Cursor?
	} else {
		//This works because the math is being done on ints then turned into a float giving us a nice even number.
		op.GeoM.Translate(float64((g.Cursor.X/config.TileSizeW)*config.TileSizeW), float64((g.Cursor.Y/config.TileSizeH)*config.TileSizeH))
		screen.DrawImage(g.uiTileset.SubImage(image.Rect(128, cursorY, 128+config.SpriteSizeW, cursorY+config.SpriteSizeH)).(*ebiten.Image), op)
	}
}

func (g *Game) DrawLevel(screen *ebiten.Image, aX int, aY int, width int, height int, blind bool, centered bool) {
	left := aX - width/2
	right := aX + width/2
	up := aY - height/2
	down := aY + height/2

	if !centered {
		left = aX
		right = aX + width - 1
		up = aY
		down = aY + height
	}

	screenX := 0
	screenY := 0
	for x := left; x <= right; x++ {
		screenY = 0
		for y := up; y <= down; y++ {
			tile := g.level.GetTileAt(x, y)
			if blind {
				if y < aY-height/4 || y > aY+height/4 || x > aX+width/4 || x < aX-width/4 {
					tile = nil
				}
			}

			//Draw tile
			tX := float64(screenX * config.TileSizeW)
			tY := float64(screenY * config.TileSizeH)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(float64(config.TileSizeW/config.SpriteSizeW), float64(config.TileSizeH/config.SpriteSizeH))
			op.GeoM.Translate(tX, tY)

			if tile == nil {
				screen.DrawImage(g.worldTileset.SubImage(image.Rect(0, 112, config.SpriteSizeW, 112+config.SpriteSizeH)).(*ebiten.Image), op)
				continue
			} else {
				screen.DrawImage(g.worldTileset.SubImage(image.Rect(tile.SpriteX, tile.SpriteY, tile.SpriteX+config.SpriteSizeW, tile.SpriteY+config.SpriteSizeH)).(*ebiten.Image), op)
			}

			//Draw entity on tile.  We do this here to prevent yet another loop. ;)
			entity := g.level.GetEntityAt(tile.X, tile.Y)
			g.DrawEntity(screen, entity, tX, tY)

			screenY++
		}
		screenX++
	}
}

func (g *Game) DrawEntity(screen *ebiten.Image, entity *entity.Entity, x float64, y float64) {
	//Draw entity on tile.
	if entity != nil {
		if entity.HasComponent("AppearanceComponent") {
			ac := entity.GetComponent("AppearanceComponent").(*component.AppearanceComponent)
			dir := 0
			if entity.HasComponent("DirectionComponent") {
				dc := entity.GetComponent("DirectionComponent").(*component.DirectionComponent)
				dir = dc.Direction
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(float64(config.TileSizeW/config.SpriteSizeW), float64(config.TileSizeH/config.SpriteSizeH))

			if entity.HasComponent("DeadComponent") {
				op.GeoM.Rotate(180)
			}

			op.GeoM.Translate(x, y)

			if entity.HasComponent("InanimateComponent") {
				screen.DrawImage(g.worldTileset.SubImage(image.Rect(ac.SpriteX+dir*config.SpriteSizeW, ac.SpriteY, ac.SpriteX+config.SpriteSizeW+dir*config.SpriteSizeW, ac.SpriteY+config.SpriteSizeH)).(*ebiten.Image), op)
			} else {
				screen.DrawImage(g.characterTileset.SubImage(image.Rect(ac.SpriteX+dir*config.SpriteSizeW, ac.SpriteY, ac.SpriteX+config.SpriteSizeW+dir*config.SpriteSizeW, ac.SpriteY+config.SpriteSizeH)).(*ebiten.Image), op)
			}

			//Draw FX
			if entity.HasComponent("AttackComponent") {
				attackC := entity.GetComponent("AttackComponent").(*component.AttackComponent)
				if attackC.Frame == 3 {
					entity.RemoveComponent("AttackComponent")
				} else {
					xOffset := attackC.SpriteX + (attackC.Frame * config.SpriteSizeW)
					screen.DrawImage(g.fxTileset.SubImage(image.Rect(xOffset, attackC.SpriteY, xOffset+config.SpriteSizeW, attackC.SpriteY+config.SpriteSizeH)).(*ebiten.Image), op)
					attackC.Frame++
				}
			}
		}

	}
}

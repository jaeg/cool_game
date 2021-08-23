package game

import (
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jaeg/cool_game/world"
)

//TileSizeW width of the tile when rendered
const TileSizeW = 16

//TileSizeH height of the tile when rendered
const TileSizeH = 16

//SpriteSizeH Height of the sprite in the tileset.
const SpriteSizeH = 16

//SpriteSizeW Width of the sprite in the tileset.
const SpriteSizeW = 16

type Game struct {
	worldTileset *ebiten.Image
	level        *world.Level
	Width        int
	Height       int
	CameraX      int
	CameraY      int
	keys         []ebiten.Key
}

func NewGame(title string, width int, height int) (*Game, error) {
	game := &Game{Width: width, Height: height}
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle(title)

	//Load assets
	imgFile, err := ebitenutil.OpenFile("assets/tiny_dungeon_world.png")
	if err != nil {
		fmt.Println("Error opening tileset")
		return nil, errors.New("error opening tileset")
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Fatal(err)
	}
	game.worldTileset = ebiten.NewImageFromImage(img)

	game.level = world.NewOverworldSection(100, 100)

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
			g.CameraY--
		}
		if k.String() == "S" {
			g.CameraY++
		}
		if k.String() == "A" {
			g.CameraX--
		}
		if k.String() == "D" {
			g.CameraX++
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//Draw world
	view := g.level.GetView(g.CameraX, g.CameraY, 100, 100, false, false)
	for x := 0; x < len(view); x++ {
		for y := 0; y < len(view[x]); y++ {
			tX := float64(x * SpriteSizeW)
			tY := float64(y * SpriteSizeH)
			tile := view[x][y]

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(tX, tY)
			op.GeoM.Scale(float64(TileSizeW/SpriteSizeW), float64(TileSizeH/SpriteSizeH))

			if tile == nil {
				screen.DrawImage(g.worldTileset.SubImage(image.Rect(0, 112, SpriteSizeH, 112+SpriteSizeH)).(*ebiten.Image), op)
			} else {
				screen.DrawImage(g.worldTileset.SubImage(image.Rect(tile.SpriteX, tile.SpriteY, tile.SpriteX+SpriteSizeH, tile.SpriteY+SpriteSizeH)).(*ebiten.Image), op)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.Width, g.Height
}

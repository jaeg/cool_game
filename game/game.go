package game

import (
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	worldTileset *ebiten.Image
	Width        int
	Height       int
}

func NewGame(title string, width int, height int) (*Game, error) {
	game := &Game{Width: width, Height: height}
	ebiten.SetWindowSize(width*2, width*2)
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

	return game, nil
}

func (g *Game) Run() error {
	err := ebiten.RunGame(g)
	return err
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	screen.DrawImage(g.worldTileset.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image), op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.Width, g.Height
}

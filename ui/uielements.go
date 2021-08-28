package ui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jaeg/cool_game/config"
	"github.com/jaeg/cool_game/resource"
)

type Button struct {
	X      int
	Y      int
	Width  int
	Height int
	Text   string
	IconX  int
	IconY  int
}

func (b Button) Draw(screen *ebiten.Image) {
	for x := b.X; x < b.X+b.Width; x += 16 {
		for y := b.Y; y < b.Y+b.Height; y += 16 {
			sX := 127
			sY := 16
			//Left Top
			if x == b.X && y == b.Y {
				sY = 0
				sX = 144
			} else if x == b.X+b.Width-16 && y == b.Y { //Right top
				sY = 0
				sX = 176
			} else if x == b.X && y == b.Y+b.Height-16 { //Left bottom
				sY = 32
				sX = 144
			} else if x == b.X+b.Width-16 && y == b.Y+b.Height-16 { //Right bottom
				sY = 32
				sX = 176
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			//s.drawSpriteEx(int32(x), int32(y), sX, sY, 32, 32, 255, 255, 255, 255, s.uiTexture)
			screen.DrawImage(resource.Textures["ui"].SubImage(image.Rect(sX, sY, sX+config.SpriteSizeW, sY+config.SpriteSizeH)).(*ebiten.Image), op)

		}
	}
}

func (b Button) IsWithin(cX int, cY int) bool {
	if cX >= b.X && cX <= b.X+b.Width && cY >= b.Y && cY <= b.Height+b.Y {
		return true
	}
	return false
}

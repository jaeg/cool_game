package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jaeg/cool_game/config"
	"github.com/jaeg/cool_game/resource"
	"github.com/jaeg/cool_game/state"
	"github.com/jaeg/cool_game/ui"
)

// Main gui
type GUIViewMain struct {
	ui.GUIViewBase
	minimap *ebiten.Image
	x       int
}

func (g *GUIViewMain) Update(s state.StateInterface) {
	g.x++
	mainState, ok := s.(*MainState)
	if ok {
		if g.minimap == nil {
			g.minimap = mainState.GetMinimap(0, 0, config.WorldGenSizeW, config.WorldGenSizeH, 150, 150)
		}
	}

}

func (g *GUIViewMain) Draw(screen *ebiten.Image, s state.StateInterface) {
	//Draw sidebar
	for x := config.World_W; x < config.ScreenWidth; x += 16 {
		for y := 0; y < config.ScreenHeight; y += 16 {
			sX := 127
			sY := 16
			//Left Top
			if x == config.World_W && y == 0 {
				sY = 0
				sX = 144
			} else if x == config.ScreenWidth-16 && y == 0 { //Right top
				sY = 0
				sX = 176
			} else if x == config.World_W && y == config.ScreenHeight-16 { //Left bottom
				sY = 32
				sX = 144
			} else if x == config.ScreenWidth-16 && y == config.ScreenHeight-16 { //Right bottom
				sY = 32
				sX = 176
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			//s.drawSpriteEx(int32(x), int32(y), sX, sY, 32, 32, 255, 255, 255, 255, s.uiTexture)
			screen.DrawImage(resource.Textures["ui"].SubImage(image.Rect(sX, sY, sX+config.SpriteSizeW, sY+config.SpriteSizeH)).(*ebiten.Image), op)
		}
	}

	//Draw Minimap
	if g.minimap != nil {
		op := &ebiten.DrawImageOptions{}
		//op.GeoM.Scale(.2, .2)
		op.GeoM.Translate(config.World_W+5, 16)
		screen.DrawImage(g.minimap, op)
	}

	//Draw buttons
	for _, b := range g.Buttons {
		b.Draw(screen)
	}
}

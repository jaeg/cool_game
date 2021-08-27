package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jaeg/cool_game/config"
	"github.com/jaeg/cool_game/resource"
)

//Base GUIView interface.
//Since we are dealing with interfaces the GUIView is being passed around by value instead of reference
type GUIViewInterface interface {
	Update(game *Game) GUIViewInterface
	Draw(screen *ebiten.Image, game *Game)
}

//GUIViewBase gives views some basic functionality when inherited.
type GUIViewBase struct {
	Buttons []*Button
}

func (g *GUIViewBase) AddButton(button *Button) {
	if g.Buttons == nil {
		g.Buttons = make([]*Button, 0)
	}
	g.Buttons = append(g.Buttons, button)
}

// Main gui
type GUIViewMain struct {
	GUIViewBase
	minimap *ebiten.Image
	x       int
}

func (g GUIViewMain) Update(game *Game) GUIViewInterface {
	g.x++
	if g.minimap == nil {
		g.minimap = game.GetMinimap(0, 0, config.WorldGenSizeW, config.WorldGenSizeH, 150, 150)
	}
	return g
}

func (g GUIViewMain) Draw(screen *ebiten.Image, game *Game) {
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

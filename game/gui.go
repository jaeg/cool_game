package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jaeg/cool_game/config"
	"github.com/jaeg/cool_game/resource"
)

//GUI Main struct that manages the gui for the game. Includes the cursor
type GUI struct {
	State GUIViewInterface
}

func NewGUI() *GUI {
	mainState := &GUIViewMain{}
	button := &Button{X: 16, Y: 16, Width: 64, Height: 32, Text: "Test"}
	mainState.AddButton(button)
	return &GUI{State: mainState}
}

func (g *GUI) Update(game *Game) {
	if g.State != nil {
		g.State.Update(game)
	}
}

func (g *GUI) Draw(screen *ebiten.Image, game *Game) {
	g.State.Draw(screen, game)

	g.DrawCursor(screen, game)
}

// GetMinimap
// Generates a minimap image of specified size and returns the image.
// Width and Height are in tiles not pixels.
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
				worldImage.DrawImage(resource.Textures["world"].SubImage(image.Rect(0, 112, config.SpriteSizeW, 112+config.SpriteSizeH)).(*ebiten.Image), op)
				continue
			} else {
				worldImage.DrawImage(resource.Textures["world"].SubImage(image.Rect(tile.SpriteX, tile.SpriteY, tile.SpriteX+config.SpriteSizeW, tile.SpriteY+config.SpriteSizeH)).(*ebiten.Image), op)
			}
		}
	}

	return worldImage
}

func (g *GUI) DrawCursor(screen *ebiten.Image, game *Game) {
	//Cursor logic

	cX, cY := ebiten.CursorPosition()

	var cursorY = 128
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		cursorY = 144
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(config.TileSizeW/config.SpriteSizeW), float64(config.TileSizeH/config.SpriteSizeH))

	if cX > config.World_W {
		op.GeoM.Translate(float64(cX), float64(cY))
		screen.DrawImage(resource.Textures["ui"].SubImage(image.Rect(64, cursorY, 64+config.SpriteSizeW, cursorY+config.SpriteSizeH)).(*ebiten.Image), op)
		//s.drawSprite(int32(g.Cursor.X), int32(g.Cursor.Y), 64, cursorY, 255, 255, 255, g.uiTexture) //Cursor?
	} else {
		//This works because the math is being done on ints then turned into a float giving us a nice even number.
		op.GeoM.Translate(float64((cX/config.TileSizeW)*config.TileSizeW), float64((cY/config.TileSizeH)*config.TileSizeH))
		screen.DrawImage(resource.Textures["ui"].SubImage(image.Rect(128, cursorY, 128+config.SpriteSizeW, cursorY+config.SpriteSizeH)).(*ebiten.Image), op)
	}
}

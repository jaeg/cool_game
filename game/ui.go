package game

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

func (b Button) Draw(screen *ebiten.Image, game *Game) {
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

//Base GUIState interface.
//Since we are dealing with interfaces the guistate is being passed around by value instead of reference

type GUIState interface {
	Update(game *Game) GUIState
	Draw(screen *ebiten.Image, game *Game)
}

type GUIStateMain struct {
	x       int
	buttons []*Button
}

func (g GUIStateMain) Update(game *Game) GUIState {
	g.x++
	return g
}

func (g GUIStateMain) Draw(screen *ebiten.Image, game *Game) {
	for _, b := range g.buttons {
		b.Draw(screen, game)
	}
}

//GUI Main struct that manages the gui for the game. Includes the cursor
type GUI struct {
	State GUIState
}

func NewGUI() *GUI {
	mainState := GUIStateMain{buttons: make([]*Button, 0)}
	button := &Button{X: 16, Y: 16, Width: 64, Height: 32, Text: "Test"}
	mainState.buttons = append(mainState.buttons, button)
	return &GUI{State: mainState}
}

func (g *GUI) Update(game *Game) {
	if g.State != nil {
		g.State = g.State.Update(game)
	}
}

func (g *GUI) DrawUI(screen *ebiten.Image, game *Game) {
	//Draw menu
	for x := config.World_W; x < game.Width; x += 16 {
		for y := 0; y < game.Height; y += 16 {
			sX := 127
			sY := 16
			//Left Top
			if x == config.World_W && y == 0 {
				sY = 0
				sX = 144
			} else if x == game.Width-16 && y == 0 { //Right top
				sY = 0
				sX = 176
			} else if x == config.World_W && y == game.Height-16 { //Left bottom
				sY = 32
				sX = 144
			} else if x == game.Width-16 && y == game.Height-16 { //Right bottom
				sY = 32
				sX = 176
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			//s.drawSpriteEx(int32(x), int32(y), sX, sY, 32, 32, 255, 255, 255, 255, s.uiTexture)
			screen.DrawImage(resource.Textures["ui"].SubImage(image.Rect(sX, sY, sX+config.SpriteSizeW, sY+config.SpriteSizeH)).(*ebiten.Image), op)

		}
	}

	g.State.Draw(screen, game)

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

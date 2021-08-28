package game

import (
	"image"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jaeg/cool_game/component"
	"github.com/jaeg/cool_game/config"
	"github.com/jaeg/cool_game/entity"
	"github.com/jaeg/cool_game/resource"
	"github.com/jaeg/cool_game/system"
	"github.com/jaeg/cool_game/ui"
	"github.com/jaeg/cool_game/world"
)

type MainState struct {
	level   *world.Level
	CameraX int
	CameraY int
	keys    []ebiten.Key
	gui     *ui.GUI
	systems []system.System
	gm      *GameMaster
}

func NewMainState() (*MainState, error) {
	s := &MainState{gui: ui.NewGUI(&GUIViewMain{})}

	s.level = world.NewOverworldSection(config.WorldGenSizeW, config.WorldGenSizeH)

	s.systems = make([]system.System, 0)

	//Initiative System
	s.systems = append(s.systems, system.InitiativeSystem{})

	//AI System
	s.systems = append(s.systems, system.AISystem{})

	//StatusCondition System
	s.systems = append(s.systems, system.StatusConditionSystem{})

	s.gm = &GameMaster{}
	s.gm.Init(s.level)
	return s, nil
}

func (s *MainState) Update() {
	s.keys = inpututil.PressedKeys()
	for _, k := range s.keys {
		if k.String() == "W" {
			if s.CameraY > 0 {
				s.CameraY--
			}
		}
		if k.String() == "S" {
			if (s.CameraY + config.ScreenHeight/config.TileSizeH) < config.WorldGenSizeH {
				s.CameraY++
			}
		}
		if k.String() == "A" {
			if s.CameraX > 0 {
				s.CameraX--
			}
		}
		if k.String() == "D" {
			if (s.CameraX + config.ScreenWidth/config.TileSizeW) < config.WorldGenSizeW {
				s.CameraX++
			}
		}
	}

	s.gui.Update(s)

	fps := ebiten.CurrentFPS()
	ebiten.SetWindowTitle(config.Title + " FPS: " + strconv.FormatFloat(fps, 'f', 1, 64))
	s.gm.Update()
	for _, entity := range s.level.Entities {
		for system := range s.systems {
			s.systems[system].Update(s.level, entity)
		}
	}
	cs := system.CleanUpSystem{}
	cs.Update(s.level)
}

func (s *MainState) Draw(screen *ebiten.Image) {
	//Draw world
	worldImage := ebiten.NewImage(config.World_W, config.World_H)
	s.DrawLevel(worldImage, s.CameraX, s.CameraY, config.World_W/config.TileSizeW, config.World_H/config.TileSizeH, false, false)

	screen.DrawImage(worldImage, nil)
	s.gui.Draw(screen, s)
}

func (s *MainState) DrawLevel(screen *ebiten.Image, aX int, aY int, width int, height int, blind bool, centered bool) {
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
			tile := s.level.GetTileAt(x, y)
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
				screen.DrawImage(resource.Textures["world"].SubImage(image.Rect(0, 112, config.SpriteSizeW, 112+config.SpriteSizeH)).(*ebiten.Image), op)
				continue
			} else {
				screen.DrawImage(resource.Textures["world"].SubImage(image.Rect(tile.SpriteX, tile.SpriteY, tile.SpriteX+config.SpriteSizeW, tile.SpriteY+config.SpriteSizeH)).(*ebiten.Image), op)
			}

			//Draw entity on tile.  We do this here to prevent yet another loop. ;)
			entity := s.level.GetEntityAt(tile.X, tile.Y)
			s.DrawEntity(screen, entity, tX, tY)

			screenY++
		}
		screenX++
	}
}

func (s *MainState) DrawEntity(screen *ebiten.Image, entity *entity.Entity, x float64, y float64) {
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
				op.GeoM.Scale(1, -1)
				op.GeoM.Translate(0, config.TileSizeH)
			}
			op.GeoM.Translate(x, y)

			// TODO - I don't like this.  The appearance component should specify the resource.
			if entity.HasComponent("InanimateComponent") {
				screen.DrawImage(resource.Textures["world"].SubImage(image.Rect(ac.SpriteX+dir*config.SpriteSizeW, ac.SpriteY, ac.SpriteX+config.SpriteSizeW+dir*config.SpriteSizeW, ac.SpriteY+config.SpriteSizeH)).(*ebiten.Image), op)
			} else {
				screen.DrawImage(resource.Textures["characters"].SubImage(image.Rect(ac.SpriteX+dir*config.SpriteSizeW, ac.SpriteY, ac.SpriteX+config.SpriteSizeW+dir*config.SpriteSizeW, ac.SpriteY+config.SpriteSizeH)).(*ebiten.Image), op)
			}

			//Draw FX
			if entity.HasComponent("AttackComponent") {
				attackC := entity.GetComponent("AttackComponent").(*component.AttackComponent)
				if attackC.Frame == 3 {
					entity.RemoveComponent("AttackComponent")
				} else {
					xOffset := attackC.SpriteX + (attackC.Frame * config.SpriteSizeW)
					screen.DrawImage(resource.Textures["fx"].SubImage(image.Rect(xOffset, attackC.SpriteY, xOffset+config.SpriteSizeW, attackC.SpriteY+config.SpriteSizeH)).(*ebiten.Image), op)
					attackC.Frame++
				}
			}
		}

	}
}

// GetMinimap
// Generates a minimap image of specified size and returns the image.
// Width and Height are in tiles not pixels.
func (g *MainState) GetMinimap(sX int, sY int, width int, height int, imageWidth int, imageHeight int) *ebiten.Image {
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

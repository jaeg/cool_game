package game

import (
	"image"
	_ "image/png"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jaeg/cool_game/component"
	"github.com/jaeg/cool_game/config"
	"github.com/jaeg/cool_game/entity"
	"github.com/jaeg/cool_game/resource"
	"github.com/jaeg/cool_game/system"
	"github.com/jaeg/cool_game/world"
)

type Game struct {
	title   string
	level   *world.Level
	Width   int
	Height  int
	CameraX int
	CameraY int
	keys    []ebiten.Key
	gui     *GUI
	systems []system.System
	gm      *GameMaster
}

func NewGame(title string, width int, height int) (*Game, error) {
	game := &Game{Width: width, Height: height, title: title, gui: NewGUI()}
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle(title)

	//Load assets
	err := resource.LoadImageAsTexture("world", "assets/tiny_dungeon_world.png")
	if err != nil {
		return nil, err
	}

	err = resource.LoadImageAsTexture("characters", "assets/tiny_dungeon_monsters.png")
	if err != nil {
		return nil, err
	}

	err = resource.LoadImageAsTexture("ui", "assets/tiny_dungeon_interface.png")
	if err != nil {
		return nil, err
	}

	err = resource.LoadImageAsTexture("fx", "assets/tiny_dungeon_fx.png")
	if err != nil {
		return nil, err
	}

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

	g.gui.Update(g)

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

	screen.DrawImage(worldImage, nil)
	g.gui.Draw(screen, g)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.Width, g.Height
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
				screen.DrawImage(resource.Textures["world"].SubImage(image.Rect(0, 112, config.SpriteSizeW, 112+config.SpriteSizeH)).(*ebiten.Image), op)
				continue
			} else {
				screen.DrawImage(resource.Textures["world"].SubImage(image.Rect(tile.SpriteX, tile.SpriteY, tile.SpriteX+config.SpriteSizeW, tile.SpriteY+config.SpriteSizeH)).(*ebiten.Image), op)
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

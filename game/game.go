package game

import (
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jaeg/cool_game/config"
	"github.com/jaeg/cool_game/resource"
	"github.com/jaeg/cool_game/state"
)

type Game struct {
	title        string
	CurrentState state.StateInterface
}

func NewGame(title string) (*Game, error) {
	game := &Game{title: title}
	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)
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

	err = resource.LoadFont("main", "assets/Roboto-Regular.ttf")
	if err != nil {
		return nil, err
	}

	game.CurrentState, err = NewMainState()
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (g *Game) Run() error {
	err := ebiten.RunGame(g)
	return err
}

func (g *Game) Update() error {
	g.CurrentState.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.CurrentState.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}

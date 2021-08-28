package state

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type StateInterface interface {
	Update()
	Draw(screen *ebiten.Image)
}

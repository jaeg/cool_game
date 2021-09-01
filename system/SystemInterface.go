package system

import (
	"github.com/jaeg/cool_game/world"
	"github.com/jaeg/game-engine/entity"
)

// System base system interface
type System interface {
	Update(*world.Level, *entity.Entity)
}

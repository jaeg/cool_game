package system

import (
	"github.com/jaeg/cool_game/entity"
	"github.com/jaeg/cool_game/world"
)

// System base system interface
type System interface {
	Update(*world.Level, *entity.Entity) *world.Level
}

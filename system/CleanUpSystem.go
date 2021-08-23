package system

import (
	"github.com/jaeg/cool_game/component"
	"github.com/jaeg/cool_game/world"
)

type CleanUpSystem struct {
}

// CleanUpSystem .
func (s CleanUpSystem) Update(level *world.Level) *world.Level {
	for _, entity := range level.Entities {
		if entity.HasComponent("MyTurnComponent") {
			entity.RemoveComponent("MyTurnComponent")
		}

		if entity.HasComponent("DeadComponent") {
			if entity.HasComponent("FoodComponent") {
				fc := entity.GetComponent("FoodComponent").(*component.FoodComponent)
				if fc.Amount <= 0 {
					level.RemoveEntity(entity)
				}
			} else {
				level.RemoveEntity(entity)
			}
		}

	}

	return level
}

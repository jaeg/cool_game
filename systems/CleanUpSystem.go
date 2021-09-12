package systems

import (
	"github.com/jaeg/cool_game/components"
	"github.com/jaeg/cool_game/world"
)

type CleanUpSystem struct {
}

// CleanUpSystem .
func (s CleanUpSystem) Update(level *world.Level) {
	for _, entity := range level.Entities {
		if entity.HasComponent("MyTurnComponent") {
			entity.RemoveComponent("MyTurnComponent")
		}

		if entity.HasComponent("DeadComponent") {
			if entity.HasComponent("FoodComponent") {
				fc := entity.GetComponent("FoodComponent").(*components.FoodComponent)
				if fc.Amount <= 0 {
					level.RemoveEntity(entity)
				}
			} else {
				level.RemoveEntity(entity)
			}
		}

	}

}

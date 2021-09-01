package system

import (
	"github.com/jaeg/cool_game/components"
	"github.com/jaeg/cool_game/world"
	"github.com/jaeg/game-engine/entity"
)

type StatusConditionSystem struct {
}

var statusConditions = []string{"Poisoned", "Alerted"}

// StatusConditionSystem .
func (s StatusConditionSystem) Update(level *world.Level, entity *entity.Entity) {

	for _, statusCondition := range statusConditions {
		if entity.HasComponent(statusCondition + "Component") {
			pc := entity.GetComponent(statusCondition + "Component").(components.DecayingComponent)

			if pc.Decay() {
				entity.RemoveComponent(statusCondition + "Component")
			}
		}
	}

}

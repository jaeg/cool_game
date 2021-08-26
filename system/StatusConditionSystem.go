package system

import (
	"github.com/jaeg/cool_game/component"
	"github.com/jaeg/cool_game/entity"
	"github.com/jaeg/cool_game/world"
)

type StatusConditionSystem struct {
}

var statusConditions = []string{"Poisoned", "Alerted"}

// StatusConditionSystem .
func (s StatusConditionSystem) Update(level *world.Level, entity *entity.Entity) {

	for _, statusCondition := range statusConditions {
		if entity.HasComponent(statusCondition + "Component") {
			pc := entity.GetComponent(statusCondition + "Component").(component.DecayingComponent)

			if pc.Decay() {
				entity.RemoveComponent(statusCondition + "Component")
			}
		}
	}

}

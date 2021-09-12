package systems

import (
	"github.com/jaeg/cool_game/components"
	"github.com/jaeg/game-engine/entity"
)

type StatusConditionSystem struct {
}

var statusConditions = []string{"Poisoned", "Alerted"}

// StatusConditionSystem .
func (s StatusConditionSystem) Update(levelInterface interface{}, entity *entity.Entity) error {
	for _, statusCondition := range statusConditions {
		if entity.HasComponent(statusCondition + "Component") {
			pc := entity.GetComponent(statusCondition + "Component").(components.DecayingComponent)

			if pc.Decay() {
				entity.RemoveComponent(statusCondition + "Component")
			}
		}
	}
	return nil
}

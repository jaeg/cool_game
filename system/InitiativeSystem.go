package system

import (
	"github.com/jaeg/cool_game/components"
	"github.com/jaeg/cool_game/world"
	"github.com/jaeg/game-engine/entity"
)

type InitiativeSystem struct {
}

// InitiativeSystem .
func (s InitiativeSystem) Update(level *world.Level, entity *entity.Entity) {

	if entity.HasComponent("InitiativeComponent") {
		ic := entity.GetComponent("InitiativeComponent").(*components.InitiativeComponent)
		ic.Ticks--

		if ic.Ticks <= 0 {
			ic.Ticks = ic.DefaultValue
			if ic.OverrideValue > 0 {
				ic.Ticks = ic.OverrideValue
			}

			if !entity.HasComponent("MyTurnComponent") {

				//Handle sleep schedules.
				canGo := false
				if entity.HasComponent("NocturnalComponent") {
					if level.Hour >= 20 || (level.Hour >= 0 && level.Hour <= 7) {
						canGo = true
					}
				} else {
					if level.Hour <= 20 && level.Hour >= 7 {
						canGo = true
					}
				}

				//Can't sleep if alerted.
				if entity.HasComponent("AlertedComponent") {
					canGo = true
				}

				//Never sleeps
				if entity.HasComponent("NeverSleepComponent") {
					canGo = true
				}

				if canGo {
					mTC := &components.MyTurnComponent{}
					entity.AddComponent(mTC)
				}
			}
		}
	}
}

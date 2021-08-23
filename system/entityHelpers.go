package system

import (
	"github.com/jaeg/cool_game/component"
	"github.com/jaeg/cool_game/entity"
	"github.com/jaeg/cool_game/world"
)

func hit(entity *entity.Entity, entityHit *entity.Entity) {
	if entityHit != entity {
		//Attack it
		if entityHit.HasComponent("HealthComponent") {
			damage := 1
			if entityHit.HasComponent("DamageComponent") {
				dc := entityHit.GetComponent("DamageComponent").(*component.DamageComponent)
				damage = dc.Amount
			}
			ehc := entityHit.GetComponent("HealthComponent").(*component.HealthComponent)
			if ehc.Health > 0 {
				ehc.Health -= damage
			}
		}

		//Create visual of attack
		entityHit.AddComponent(&component.AttackComponent{SpriteX: 192, SpriteY: 80})

		// Trigger their defenses
		if entityHit.HasComponent("DefensiveAIComponent") {
			daic := entityHit.GetComponent("DefensiveAIComponent").(*component.DefensiveAIComponent)
			pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
			daic.Attacked = true
			daic.AttackerX = pc.GetX()
			daic.AttackerY = pc.GetY()
		}

		if !entityHit.HasComponent("AlertedComponent") {
			entityHit.AddComponent(&component.AlertedComponent{Duration: 120})
		}

		if entity.HasComponent("PoisonousComponent") {
			if !entityHit.HasComponent("PoisonedComponent") {
				poisonousComponent := entity.GetComponent("PoisonousComponent").(*component.PoisonousComponent)
				entityHit.AddComponent(&component.PoisonedComponent{Duration: poisonousComponent.Duration})
			}
		}
	}
}

//Returns true if successfully ate.
func eat(entity *entity.Entity, entityHit *entity.Entity) bool {
	if entityHit != entity {
		if entityHit.HasComponent("FoodComponent") {
			fc := entityHit.GetComponent("FoodComponent").(*component.FoodComponent)
			fc.Amount--
			return true
		}
	}
	return false
}

func face(entity *entity.Entity, deltaX int, deltaY int) {
	dc := entity.GetComponent("DirectionComponent").(*component.DirectionComponent)
	if deltaY > 0 {
		dc.Direction = 1
	}
	if deltaY < 0 {
		dc.Direction = 2
	}
	if deltaX < 0 {
		dc.Direction = 3
	}
	if deltaX > 0 {
		dc.Direction = 0
	}
}

func handleDeath(entity *entity.Entity) bool {
	if entity.HasComponent("HealthComponent") {
		hc := entity.GetComponent("HealthComponent").(*component.HealthComponent)
		if hc.Health <= 0 {
			entity.AddComponent(&component.DeadComponent{})

			return true
		}
	}
	return false
}

//Returns true if a solid entity is in the way.
func move(entity *entity.Entity, level *world.Level, deltaX int, deltaY int) bool {
	pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
	entityHit := level.GetSolidEntityAt(pc.GetX()+deltaX, pc.GetY()+deltaY)
	if entityHit == nil {
		tile := level.GetTileAt(pc.GetX()+deltaX, pc.GetY()+deltaY)
		if tile == nil {
		} else if tile.Type != 2 && tile.Type != 4 {
			level.PlaceEntity(pc.GetX()+deltaX, pc.GetY()+deltaY, entity)

		}
		return false
	} else {
		//If we are massive we destroy what we moved into if it's not also massive.
		if entity.HasComponent("MassiveComponent") && !entityHit.HasComponent("MassiveComponent") {
			entityHit.AddComponent(&component.DeadComponent{})
		}
	}
	return true
}

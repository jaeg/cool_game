package components

// DirectionComponent .
type DirectionComponent struct {
	Direction int
}

func (pc DirectionComponent) GetType() string {
	return "DirectionComponent"
}

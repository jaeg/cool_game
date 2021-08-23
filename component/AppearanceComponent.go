package component

// MyTurnComponent .
type AppearanceComponent struct {
	SpriteX  int32
	SpriteY  int32
	Resource string
	R, G, B  uint8
}

func (pc AppearanceComponent) GetType() string {
	return "AppearanceComponent"
}

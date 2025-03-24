package components

import "github.com/yohamta/donburi"

type GridPositionData struct {
	Col int
	Row int
}

var GridPos = donburi.NewComponentType[GridPositionData]()

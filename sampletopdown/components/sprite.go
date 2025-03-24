package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

var Sprite = donburi.NewComponentType[ebiten.Image]()

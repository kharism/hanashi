package components

import (
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func init() {
	InteractibleMap = map[int]Interact{}
}

type HanashiSceneSetter interface {
	SetHanashiScene(scene *core.Scene)
	UnsetHanashiScene()
}
type Interact func(ecs *ecs.ECS, sceneSetter HanashiSceneSetter)

var Interactibles = donburi.NewComponentType[Interact]()

var InteractibleMap map[int]Interact

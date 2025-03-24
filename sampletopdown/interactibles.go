package main

import (
	"github.com/kharism/hanashi/core"
	"github.com/kharism/hanashi/sampletopdown/components"
	"github.com/yohamta/donburi/ecs"
)

var scene *core.Scene

func ChestInteract(ecs *ecs.ECS, sceneSetter components.HanashiSceneSetter) {
	if scene == nil {
		scene = core.NewScene()
	}

	scene.Events = []core.Event{
		core.NewDialogueEvent("", "It's Empty", face),
	}
	scene.Done = func() {
		sceneSetter.UnsetHanashiScene()
	}
	scene.Events[0].Execute(scene)
	sceneSetter.SetHanashiScene(scene)
}
func DoorInteract(ecs *ecs.ECS, sceneSetter components.HanashiSceneSetter) {
	if scene == nil {
		scene = core.NewScene()
	}

	scene.Events = []core.Event{
		core.NewDialogueEvent("", "It's Locked", face),
	}
	scene.Done = func() {
		sceneSetter.UnsetHanashiScene()
	}
	scene.Events[0].Execute(scene)
	sceneSetter.SetHanashiScene(scene)
}
func InitInteactible() {
	// 42 is the index of the chest
	components.InteractibleMap[42] = ChestInteract
	components.InteractibleMap[84] = DoorInteract
}

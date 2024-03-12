package core

// jump
type BranchingJumpEvent struct {
	SceneDataKey  string
	BranchingJump map[string]int
}

func NewBranchingJumpEvent(SceneDataKey string, BranchingJump map[string]int) Event {
	return &BranchingJumpEvent{SceneDataKey: SceneDataKey, BranchingJump: BranchingJump}
}
func (c *BranchingJumpEvent) Execute(scene *Scene) {
	if value, ok := scene.SceneData[c.SceneDataKey]; ok {
		scene.EventIndex = c.BranchingJump[value.(string)]
		scene.Events[scene.EventIndex].Execute(scene)
	}
}

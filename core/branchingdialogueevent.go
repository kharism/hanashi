package core

// create branching dialogue based on SceneDataKey.
// this code assume the value of SceneData is string
// if no matching string data is found it will show empty string
type BranchingDialogueEvent struct {
	Name            string
	SceneDataKey    string
	BranchingStruct map[string]string
}

func NewBranchingDialogueEvent(name, sceneKey string, options map[string]string) Event {
	return &BranchingDialogueEvent{Name: name, SceneDataKey: sceneKey, BranchingStruct: options}
}

func (b *BranchingDialogueEvent) Execute(scene *Scene) {
	scene.CurCharName = b.Name
	defaultText := ""
	jj, ok := scene.SceneData[b.SceneDataKey].(string)
	if ok {
		defaultText = b.BranchingStruct[jj]
	}
	scene.CurDialog = defaultText
}

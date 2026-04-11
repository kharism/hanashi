package core

// this struct is to ignore any click event
type PreventsClickEvent struct {
}

// this struct is to enable the click event back
type EnableClickEvent struct {
}

func NewPreventClickEvent() Event {
	return &PreventsClickEvent{}
}

func NewEnableClickEvent() Event {
	return &EnableClickEvent{}
}

func (c *PreventsClickEvent) Execute(scene *Scene) {
	scene.PreventsClick = true
}

func (c *EnableClickEvent) Execute(scene *Scene) {
	scene.PreventsClick = false
}

package core

// Event that SHOWS you selection option substate, by default using
// OptionPickerState
type OptionSelectEvent struct {
	Question string
	Options  []string
	Selected string
}

// create OptionSelectEvent
func NewOptionSelectEvent(question string, options ...string) Event {
	return &OptionSelectEvent{Question: question, Options: options}
}
func (b *OptionSelectEvent) Execute(scene *Scene) {
	optionsSubstate := &OptionPickerState{Question: b.Question, Options: b.Options, Scene: scene}
	optionsSubstate.InitYStart()
	scene.CurrentSubState = optionsSubstate
}

package core

// Event that SHOWS you selection option substate, by default using
// OptionPickerState
type OptionSelectEvent struct {
	QuestionId string
	Question   string
	Options    []string
	Selected   string
}

// create OptionSelectEvent
func NewOptionSelectEvent(questionId, question string, options ...string) Event {
	return &OptionSelectEvent{QuestionId: questionId, Question: question, Options: options}
}
func (b *OptionSelectEvent) Execute(scene *Scene) {
	optionsSubstate := &OptionPickerState{QId: b.QuestionId, Question: b.Question, Options: b.Options, Scene: scene}
	optionsSubstate.InitYStart()
	scene.CurrentSubState = optionsSubstate
}

package core

// complex event is event that is complex, we use it if we want to do several
// things in single click like change bg and show dialoge at the same time
type ComplexEvent struct {
	Events []Event
}

func (c *ComplexEvent) Execute(scene *Scene) {
	for _, s := range c.Events {
		s.Execute(scene)
	}
}

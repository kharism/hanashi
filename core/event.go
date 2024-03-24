package core

// This is something that must happened in a scene, whether they are dialogue.
// moving image, some fancy stuff
type Event interface {
	// do something on the screen
	Execute(scene *Scene)
}

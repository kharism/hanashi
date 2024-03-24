package core

import (
	"log"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// event that shows dialogue/monologue. The location of the Name and text
// will follow Layouter on the scene.
//
// Name is character name who tells the Dialogue. It does not
// need to be registered on AvailableCharacter on Scene. it
// can be made empty string
//
// Dialogue is the text shown.
type DialogueEvent struct {
	Name     string
	Dialogue string
	FontFace font.Face
}

var (
	// the default font used for many things. Overwrite this var
	// to use different font
	DefaultFont font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	DefaultFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
}
func NewDialogueEvent(name, dialogue string, fontface font.Face) Event {
	return &DialogueEvent{Name: name, Dialogue: dialogue, FontFace: fontface}
}
func (b *DialogueEvent) Execute(scene *Scene) {
	scene.FontFace = b.FontFace
	scene.CurCharName = b.Name
	scene.CurDialog = b.Dialogue
}

package core

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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
	FontFace text.Face
}

var (
	// the default font used for many things. Overwrite this var
	// to use different font
	DefaultFont text.Face
)

func init() {
	// tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	// DefaultFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
	// 	Size:    24,
	// 	DPI:     dpi,
	// 	Hinting: font.HintingVertical,
	// })
	DefaultFont = &text.GoTextFace{
		Source: s,
		Size:   24,
	}
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
func NewDialogueEvent(name, dialogue string, fontface text.Face) Event {
	return &DialogueEvent{Name: name, Dialogue: dialogue, FontFace: fontface}
}
func (b *DialogueEvent) Execute(scene *Scene) {
	if b.FontFace != nil {
		scene.FontFace = b.FontFace
	}
	scene.CurCharName = b.Name
	scene.CurDialog = b.Dialogue
}

package core

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type Scene struct {
	EventIndex  int
	Events      []Event
	CurrentBg   *MovableImage
	txtBgGetter TxtBgGetter
	getLayouter GetLayouter

	// text stuff
	TxtBg       *ebiten.Image
	CurCharName string
	CurDialog   string
	FontFace    font.Face

	//characters/sprite stuff
	// Characters used in the scene
	Characters []*Character
	// Characters viewable in the scene
	// Any character that needs to be shown will be moved here
	ViewableCharacters []*Character
}

func NewScene() *Scene {
	return &Scene{Events: []Event{}, Characters: []*Character{}}
}

type GetLayouter interface {
	// return width and height of the scene
	GetLayout() (width, height int)
	// return the starting text position where the box containing name of the character appear on the scene
	// return negative number if no such box needed
	GetNamePosition() (x, y int)
	// get the starting position of the text
	GetTextPosition() (x, y int)
}

type TxtBgGetter interface {
	GetTxtBg() *ebiten.Image
}

func (g *Scene) SetLayouter(getLayouter GetLayouter) {
	g.getLayouter = getLayouter
}

// copy Character from Character to ViewableCharacter
func (g *Scene) AddViewableCharacters(name string, moveParam *MoveParam, scaleParam *ScaleParam) {
	for _, c := range g.Characters {
		if c.Name == name {
			c.Img.x = moveParam.Sx
			c.Img.y = moveParam.Sy
			c.Img.ScaleParam = scaleParam
			c.Img.AddAnimation(&MoveAnimation{tx: moveParam.Tx, ty: moveParam.Ty, Speed: moveParam.Speed})
			g.ViewableCharacters = append(g.ViewableCharacters, c)
			break
		}
	}

}

// remove character from ViewableCharacter
func (g *Scene) RemoveVieableCharacter(characterName string) {
	//g.Characters = append(g.Characters, character)
	newChars := []*Character{}
	for idx, s := range g.Characters {
		if !s.CheckName(characterName) {
			newChars = append(newChars, s)
		} else {
			newChars = append(newChars, g.Characters[idx+1:]...)
			break
		}
	}
	g.ViewableCharacters = newChars

}
func (g *Scene) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		g.EventIndex += 1
		g.Events[g.EventIndex].Execute(g)
	}
	g.CurrentBg.Update()
	for _, c := range g.ViewableCharacters {
		c.Img.Update()
	}
	return nil

}
func (g *Scene) Draw(screen *ebiten.Image) {
	// ebitenutil.DebugPrint(screen, "Hello, World!")
	if g.CurrentBg != nil {
		// drawOptions := ebiten.DrawImageOptions{}
		// drawOptions.GeoM.Scale()
		// drawOptions.GeoM.Reset()
		// drawOptions.GeoM.Scale(0.5, 0.5)
		// drawOptions.GeoM.Translate(100, 100)
		// screen.DrawImage(g.CurrentBg, &drawOptions)
		g.CurrentBg.Draw(screen)
	}
	for _, c := range g.ViewableCharacters {
		c.Img.Draw(screen)
	}
	nameX, nameY := g.getLayouter.GetNamePosition()
	dialogueX, dialogueY := g.getLayouter.GetTextPosition()

	if g.TxtBg != nil {
		// w, h := g.getLayouter.GetLayout()
		opt := ebiten.DrawImageOptions{}
		opt.GeoM.Translate(float64(nameX), float64(nameY))
		screen.DrawImage(g.TxtBg, &opt)
	}

	curFont := DefaultFont
	if g.FontFace != nil {
		curFont = g.FontFace
	}

	text.Draw(screen, g.CurCharName, curFont, nameX, nameY, color.White)
	text.Draw(screen, g.CurDialog, curFont, dialogueX, dialogueY, color.White)
}

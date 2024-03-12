package core

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

// Scene is whatever happened on the screen. It has several Event that is loaded in order
// It implement Draw(*ebiten.Image) and Update() so it should be simple to be splashed into
// default ebitengine projects.
//
// use some statemanagement framework like github.com/joelschutz/stagehand to manage state/scene better
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

	SceneData map[string]any

	//characters/sprite stuff
	// Characters used in the scene
	Characters []*Character
	// Characters viewable in the scene
	// Any character that needs to be shown will be moved here
	ViewableCharacters []*Character

	//substates and its management
	CurrentSubState     SubState
	OptionPikerSubstate SubState
	// this function is executed after the scene complete
	Done func()
}

func NewScene() *Scene {
	return &Scene{Events: []Event{}, Characters: []*Character{}, SceneData: map[string]any{}}
}

type SubState interface {
	Draw(screen *ebiten.Image)
	Update()
}
type OptionPickerState struct {
	Scene         *Scene
	Question      string
	Options       []string
	OptionsYStart []int
}

func (s *OptionPickerState) Update() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		_, yCurInt := ebiten.CursorPosition()
		for i := len(s.OptionsYStart) - 1; i >= 0; i-- {
			if yCurInt > s.OptionsYStart[i] {
				fmt.Println("Clicked", s.Options[i])
				s.Scene.SetSceneData(s.Question, s.Options[i])
				// s.Scene.EventIndex += 1
				s.Scene.CurrentSubState = nil
				pp := len(s.Scene.Events)
				if s.Scene.EventIndex == pp-1 {
					s.Scene.Done()
				} else {
					s.Scene.EventIndex += 1
					s.Scene.Events[s.Scene.EventIndex].Execute(s.Scene)
				}

				break
			}
		}
	}
}

var textHeight = 50
var textOffset = 10

// init the y location of option buttons
func (s *OptionPickerState) InitYStart() {
	_, height := s.Scene.getLayouter.GetLayout()

	totalOptionsHeight := len(s.Options)*textHeight + (len(s.Options)-1)*textOffset
	y0Pos := (height / 2) - (totalOptionsHeight / 2)
	kk := []int{y0Pos}
	for i := 1; i < len(s.Options); i++ {
		jj := y0Pos + i*(textHeight+textOffset)
		kk = append(kk, jj)
	}
	s.OptionsYStart = kk
}
func (s *OptionPickerState) Draw(screen *ebiten.Image) {
	width, height := s.Scene.getLayouter.GetLayout()
	transBg := ebiten.NewImage(width, height)
	transBg.Fill(color.NRGBA{R: 100, G: 100, B: 100, A: 100})
	screen.DrawImage(transBg, nil)
	optBg := ebiten.NewImage(width, textHeight)
	optBg.Fill(color.RGBA{120, 0, 255, 255})
	optDraw := ebiten.DrawImageOptions{}
	optDraw.GeoM.Translate(0, float64(s.OptionsYStart[0]-textOffset-textHeight))
	screen.DrawImage(optBg, &optDraw)
	curFont := DefaultFont
	if s.Scene.FontFace != nil {
		curFont = s.Scene.FontFace
	}
	text.Draw(screen, s.Question, curFont, 0, s.OptionsYStart[0]-textOffset-textHeight+20, color.Black)
	for idx, opt := range s.Options {
		// draw the box
		optBg := ebiten.NewImage(width, textHeight)
		optBg.Fill(color.RGBA{120, 255, 0, 255})
		optDraw := ebiten.DrawImageOptions{}
		optDraw.GeoM.Translate(0, float64(s.OptionsYStart[idx]))
		screen.DrawImage(optBg, &optDraw)
		// draw the text

		text.Draw(screen, opt, curFont, 0, s.OptionsYStart[idx]+20, color.Black)
	}
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

func (g *Scene) SetSceneData(key string, value any) {
	g.SceneData[key] = value
}
func (g *Scene) GetSceneData(key string) any {
	return g.SceneData[key]
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
	for idx, s := range g.ViewableCharacters {
		if !s.CheckName(characterName) {
			newChars = append(newChars, s)
		} else {
			newChars = append(newChars, g.ViewableCharacters[idx+1:]...)
			break
		}
	}
	g.ViewableCharacters = newChars

}
func (g *Scene) Update() error {

	g.CurrentBg.Update()
	for _, c := range g.ViewableCharacters {
		c.Img.Update()
	}
	if g.CurrentSubState == nil {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
			g.EventIndex += 1
			if g.EventIndex >= len(g.Events) {
				g.Done()
				return nil
			} else {
				g.Events[g.EventIndex].Execute(g)
			}
		}
	} else {
		g.CurrentSubState.Update()
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
	if g.CurrentSubState == nil {
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
	} else {
		g.CurrentSubState.Draw(screen)
	}

}

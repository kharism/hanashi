package core

import (
	"bytes"
	"fmt"
	"image/color"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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
	TxtBg         *ebiten.Image
	CurCharName   string
	CurDialog     string
	VisibleDialog string
	FontFace      text.Face

	SceneData map[string]any

	AudioInterface AudioInterface

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

// An interface to the Scene. The Scene will draw either draw its substate
// or dialog scene
type SubState interface {
	Draw(screen *ebiten.Image)
	Update()
}
type OptionPickerState struct {
	Scene         *Scene
	QId           string
	Question      string
	curoptions    int
	Options       []string
	OptionsYStart []int
}

func (s *OptionPickerState) Update() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		_, yCurInt := ebiten.CursorPosition()
		for i := len(s.OptionsYStart) - 1; i >= 0; i-- {
			if yCurInt > s.OptionsYStart[i] {
				fmt.Println("Clicked", s.Options[i])
				s.Scene.SetSceneData(s.QId, s.Options[i])
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
	if inpututil.IsKeyJustReleased(ebiten.KeyDown) {
		s.curoptions += 1
		if s.curoptions > len(s.Options) {
			s.curoptions -= 1
		}
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyUp) {
		s.curoptions -= 1
		if s.curoptions < 0 {
			s.curoptions += 1
		}
	}
	if DetectKeyboardNext != nil && DetectKeyboardNext() {
		fmt.Println("Pick through keyboard", s.Options[s.curoptions])
		s.Scene.SetSceneData(s.QId, s.Options[s.curoptions])
		s.Scene.CurrentSubState = nil
		pp := len(s.Scene.Events)
		if s.Scene.EventIndex == pp-1 {
			s.Scene.Done()
		} else {
			s.Scene.EventIndex += 1
			s.Scene.Events[s.Scene.EventIndex].Execute(s.Scene)
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
	textOpt := text.DrawOptions{}
	textOpt.GeoM.Translate(0, float64(s.OptionsYStart[0]-textOffset-textHeight+20))
	textOpt.ColorScale.ScaleWithColor(color.Black)
	// text.Draw(screen, s.Question, curFont, 0, s.OptionsYStart[0]-textOffset-textHeight+20, color.Black)
	text.Draw(screen, s.Question, curFont, &textOpt)
	for idx, opt := range s.Options {
		// draw the box
		optBg := ebiten.NewImage(width, textHeight)
		optBg.Fill(color.RGBA{120, 255, 0, 255})
		optDraw := ebiten.DrawImageOptions{}
		optDraw.GeoM.Translate(0, float64(s.OptionsYStart[idx]))
		screen.DrawImage(optBg, &optDraw)
		addDist := 0.0
		//draw cursor if keyboard input is there
		if DetectKeyboardNext != nil && idx == s.curoptions {
			optDraw := ebiten.DrawImageOptions{}
			optDraw.GeoM.Translate(0, float64(s.OptionsYStart[idx]))
			screen.DrawImage(Cursor, &optDraw)
			addDist += 60
		}
		// draw the text
		textOpt := text.DrawOptions{}
		textOpt.ColorScale.ScaleWithColor(color.Black)
		textOpt.GeoM.Translate(0+addDist, float64(s.OptionsYStart[idx]+20))
		text.Draw(screen, opt, curFont, &textOpt)

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

// Set text background. Unused
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
func (g *Scene) AddViewableCharacters(name string, moveParam MoveParam, scaleParam *ScaleParam) {
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

type Pos struct {
	X float64
	Y float64
}

//go:embed assets/cursor.png
var cursor []byte
var (
	TouchIDs []ebiten.TouchID
	TouchPos map[ebiten.TouchID]*Pos
	Cursor   *ebiten.Image
)

func init() {
	TouchIDs = []ebiten.TouchID{}
	TouchPos = map[ebiten.TouchID]*Pos{}
	imgReader := bytes.NewReader(cursor)
	Cursor, _, _ = ebitenutil.NewImageFromReader(imgReader)
}

// return whether a click or tap is happened, and its location if it happened
func IsClickedOrTap() (res bool, posX int, posY int) {
	posX = -1
	posY = -1
	mouseReleased := inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0)
	if mouseReleased {
		posX, posY = ebiten.CursorPosition()
		return true, posX, posY
	}
	// fmt.Println("Check", TouchIDs)
	for _, id := range TouchIDs {
		// fmt.Println("Check", id, inpututil.TouchPressDuration(id))
		if inpututil.IsTouchJustReleased(id) {

			posX, posY = int(TouchPos[id].X), int(TouchPos[id].Y)
			// fmt.Println("touch released", id, posX, posY)

			return true, posX, posY
		}
	}
	TouchIDs = inpututil.AppendJustPressedTouchIDs(TouchIDs[:0])
	for _, id := range TouchIDs {
		x, y := ebiten.TouchPosition(id)
		TouchPos[id] = &Pos{
			X: float64(x),
			Y: float64(y),
		}
	}
	TouchIDs = ebiten.AppendTouchIDs(TouchIDs[:0])
	// fmt.Println(TouchIDs)
	return false, posX, posY
}

// a signature function to detect keyboard input to proceed to next
// item in scenario list.
type KeyboardNext func() bool

// The global function to detect KeyboardNext
var DetectKeyboardNext KeyboardNext

func (g *Scene) Update() error {
	if g.AudioInterface != nil {
		g.AudioInterface.Update()
	}
	if g.CurrentBg != nil {
		g.CurrentBg.Update()
	}

	for _, c := range g.ViewableCharacters {
		c.Img.Update()
	}
	if g.CurrentSubState == nil {
		clicked, _, _ := IsClickedOrTap()
		if clicked || (DetectKeyboardNext != nil && DetectKeyboardNext()) {
			g.EventIndex += 1
			if g.EventIndex >= len(g.Events) {
				g.Done()
				return nil
			} else {
				g.Events[g.EventIndex].Execute(g)
				// g.CurDialog = ""
			}

			g.VisibleDialog = ""
		}
	} else {
		g.CurrentSubState.Update()
		// g.CurDialog = ""
		g.VisibleDialog = ""
	}
	if g.CurDialog != g.VisibleDialog {
		// fmt.Println(g.VisibleDialog, g.CurDialog)
		g.VisibleDialog = g.CurDialog[:len(g.VisibleDialog)+1]
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
		textOpt1 := text.DrawOptions{}
		textOpt1.GeoM.Translate(float64(nameX), float64(nameY))
		textOpt1.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, g.CurCharName, curFont, &textOpt1)
		textOpt2 := text.DrawOptions{}
		textOpt2.GeoM.Translate(float64(dialogueX), float64(dialogueY))
		textOpt2.ColorScale.ScaleWithColor(color.White)
		textOpt2.LineSpacing = 24
		text.Draw(screen, g.VisibleDialog, curFont, &textOpt2)
	} else {
		g.CurrentSubState.Draw(screen)
	}

}

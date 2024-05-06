package main

import (
	"image/color"

	"github.com/kharism/hanashi/core"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type CombatSubstate interface {
	Draw(screen *ebiten.Image)
	Update()
	OnLoad()
}
type MenuButton struct {
	*core.MovableImage
	Label       string
	cursorIn    bool
	onClickFunc func()
}

func (b *MenuButton) Draw(screen *ebiten.Image) {
	if b.cursorIn {
		b.ScaleParam.Sx = 1.1
	} else {
		b.ScaleParam.Sx = 1
	}
	btnX, btnY := b.MovableImage.GetPos()
	b.MovableImage.Draw(screen)
	txtOpt := text.DrawOptions{}
	txtOpt.GeoM.Translate(btnX+10, btnY+10)
	txtOpt.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, b.Label, core.DefaultFont, &txtOpt)
}
func (b *MenuButton) Update() {
	curX, curY := ebiten.CursorPosition()
	butPosX, butPosY := b.GetPos()
	width, height := b.GetSize()
	// fmt.Println(width, height)
	if curX > int(butPosX) && curX < int(butPosX+width) && curY > int(butPosY) && curY < int(butPosY+height) {
		b.cursorIn = true
		// fmt.Println("Cursor In")
	} else {
		b.cursorIn = false
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		if b.cursorIn && b.onClickFunc != nil {
			b.onClickFunc()
		}
	}
}

type MainCombatMenu struct {
	combatScene *CombatScene

	buttons []*MenuButton
}

func (b *MainCombatMenu) OnLoad() {

}
func NewMainCombatMenu(combatScene *CombatScene) CombatSubstate {
	jj := MainCombatMenu{combatScene: combatScene}
	btn, _ := imgPool.GetImage("icon/blue_button00.png")
	// attack button
	atkBtnImgParam := core.NewMovableImageParams().WithMoveParam(core.MoveParam{Sy: 80, Sx: 0})
	atkButton := &MenuButton{MovableImage: core.NewMovableImage(btn, atkBtnImgParam), Label: "Attack"}
	atkButton.onClickFunc = func() {
		cc := CombatCommand{CharacterIdx: combatScene.CurrentChrIdx, Command: COMMAND_ATTACK}
		combatScene.CombatQueue = append(combatScene.CombatQueue, cc)
		combatScene.CurrentChrIdx += 1
		for true {
			if combatScene.CurrentChrIdx >= len(combatScene.Characters) {
				combatScene.CurrentChrIdx = -1
				combatScene.BeginCombat()
				break
			}
			if combatScene.Characters[combatScene.CurrentChrIdx].HP <= 0 {
				combatScene.CurrentChrIdx += 1
				continue
			}
		}

	}
	// back button
	backBtnImgParam := core.NewMovableImageParams().WithMoveParam(core.MoveParam{Sy: 120, Sx: 0})
	backButton := &MenuButton{MovableImage: core.NewMovableImage(btn, backBtnImgParam), Label: "Back"}
	backButton.onClickFunc = func() {
		combatScene.CurrentChrIdx -= 1
		combatScene.CombatQueue = combatScene.CombatQueue[:len(combatScene.CombatQueue)-1]
	}
	jj.buttons = append(jj.buttons, atkButton, backButton)
	return &jj
}
func (mm *MainCombatMenu) Draw(screen *ebiten.Image) {
	for _, b := range mm.buttons {
		b.Draw(screen)
	}
}
func (mm *MainCombatMenu) Update() {
	for _, b := range mm.buttons {
		b.Update()
	}
}

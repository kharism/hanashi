package main

import (
	"fmt"
	"github/kharism/hanashi/core"
	"image/color"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type CombatLogSubstate struct {
	combatScene *CombatScene
	queueIndex  int
	BattleLog   string
}

func NewCombatLogSubState(combatScene *CombatScene) CombatSubstate {
	hh := &CombatLogSubstate{combatScene: combatScene}
	return hh
}
func (c *CombatLogSubstate) Draw(screen *ebiten.Image) {
	LogBox := ebiten.NewImage(768, 100)
	LogBox.Fill(color.RGBA{R: 0, G: 0, B: 197, A: 0})
	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(0, 300)
	screen.DrawImage(LogBox, &opt)
	text.Draw(screen, c.BattleLog, core.DefaultFont, 0, 300+20, color.White)
}
func (v *CombatLogSubstate) BeginCombat() {

}

// reset queueIndex
func (v *CombatLogSubstate) OnLoad() {
	v.queueIndex = 0
	// do the combat log immediately
	v.DoCombat()
}
func (v *CombatLogSubstate) DoCombat() {
	if v.queueIndex == len(v.combatScene.CombatQueue) {
		// v.queueIndex = 0
		// v.combatScene.SwitchMenuSubstate(mainCombatMenu)
		return
	}
	cc := v.combatScene.CombatQueue[v.queueIndex]
	switch cc.Command {
	case COMMAND_ATTACK:
		character := v.combatScene.Characters[cc.CharacterIdx]
		damage := character.AtkDamage()
		cplx_anim := &ComplexAnim{cs: v.combatScene, doneFunc: func(cs *CombatScene) {
			cs.OppTakeDamage(damage)
			DoneAnim(cs)
		}, animations: []CombatAnimation{
			&MoveAttackAnim{cs: v.combatScene, PosX: float64(OPP_POS_X), PosY: float64(OPP_POS_Y)},
			&AttackAnim{cs: v.combatScene},
			&BlinkAnim{cs: v.combatScene, blinkCount: 5},
		},
		}
		v.combatScene.attackAnim.Done = func() {
			cplx_anim.idx += 1
		}
		go func() {
			v.combatScene.animationQueue <- cplx_anim
		}()

		v.BattleLog = fmt.Sprintf("%s deals %d of damage", character.Name, damage)
	case COMMAND_OPP_ATTACK:
		cc.Routine.Execute(v.combatScene, v)
	case COMMAND_END_WIN:
		v.BattleLog = fmt.Sprintf("opponent defeated")
	}
	v.queueIndex += 1

}
func (v *CombatLogSubstate) Update() {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) && v.combatScene.currentAnim == nil {
		if v.queueIndex == len(v.combatScene.CombatQueue) {
			// v.queueIndex = 0
			if v.combatScene.oppHp == 0 {
				v.combatScene.currentAnim = nil
				v.combatScene.DoneCombat()
			} else {
				v.combatScene.SwitchMenuSubstate(mainCombatMenu)
			}
			// return
		} else {
			v.DoCombat()
		}
	}
}

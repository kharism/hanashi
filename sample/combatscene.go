package main

import (
	"fmt"
	"github/kharism/hanashi/core"
	"image/color"
	"strconv"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/joelschutz/stagehand"
)

type CombatScene struct {
	state        MyState
	background   *core.MovableImage
	opp          *core.MovableImage
	hpIcon       *core.MovableImage
	oppHp        int
	attackAnim   *core.AnimatedImage
	onAttackAnim bool
	onIsBlinking bool
	counter      int
	MenuSubState CombatSubstate

	CombatQueue    []CombatCommand
	Characters     []*CombatCharacter
	CurrentChrIdx  int
	currentAnim    CombatAnimation
	animationQueue chan CombatAnimation

	DoneCombat func()
}

var (
	mainCombatMenu    CombatSubstate
	combatLogSubState CombatSubstate
)

const (
	COMMAND_ATTACK     = "CMD_ATK"
	COMMAND_OPP_ATTACK = "CMD_ATK_2"
	COMMAND_END_WIN    = "CMD_END_WIN"
	COMMAND_END_LOOSE  = "CMD_END_LOOSE"
)

type CombatCommand struct {
	CharacterIdx int
	Command      string
	Target       string
	Routine      CombatRoutine
}

// some animation stuff.
type CombatAnimation interface {
	Draw(screen *ebiten.Image)
	Update()
	GetDoneFunc() func(cs *CombatScene)
	SetDoneFunc(func(cs *CombatScene))
}

// default DoneFunc for animation
func DoneAnim(cs *CombatScene) {
	cs.currentAnim = nil
}

// switch menu. Any combat menu is applied here
func (v *CombatScene) SwitchMenuSubstate(newState CombatSubstate) {
	v.MenuSubState = newState
	v.MenuSubState.OnLoad()
	switch newState.(type) {
	case *MainCombatMenu:
		v.CurrentChrIdx = 0
		for true {
			if v.Characters[v.CurrentChrIdx].HP == 0 {
				v.CurrentChrIdx += 1
			} else {
				break
			}
		}
		v.CombatQueue = []CombatCommand{}
	case *CombatLogSubstate:
		// add AI attack procedure here
		jj := CombatCommand{Command: COMMAND_OPP_ATTACK, Routine: NewAttackRandomly(func() int {
			return 1
		})}
		v.CombatQueue = append(v.CombatQueue, jj)
	}
}
func (v *CombatScene) PlayerTakeDamage(player *CombatCharacter, damage int) {
	if damage >= player.HP {
		player.HP = 0
		logger := combatLogSubState.(*CombatLogSubstate)
		logger.BattleLog = fmt.Sprintf("%s is defeated", player.Name)
		allDead := true
		for _, c := range v.Characters {
			if c.HP > 0 {
				allDead = false
				break
			}
		}
		if allDead {
			commandEnd := CombatCommand{Command: COMMAND_END_LOOSE}
			v.CombatQueue = append(v.CombatQueue[:logger.queueIndex], commandEnd)
		}
	} else {
		player.HP -= damage
	}
}
func (v *CombatScene) OppTakeDamage(dmg int) {
	if dmg >= v.oppHp {
		v.oppHp = 0
		// v.CombatQueue
		logger := combatLogSubState.(*CombatLogSubstate)
		commandEnd := CombatCommand{Command: COMMAND_END_WIN}
		v.CombatQueue = append(v.CombatQueue[:logger.queueIndex], commandEnd)
		// logger.queueIndex
	} else {
		v.oppHp -= dmg
	}

}

var (
	OPP_POS_X = (768 / 2) - (512 / 8)
	OPP_POS_Y = 80

	PLAYER_POS_X_START = 20
	PLAYER_POS_Y_START = 450
)

func (v *CombatScene) Load(state MyState, sm *stagehand.SceneManager[MyState]) {
	v.state = state
	v.counter = 0
	img, _ := imgPool.GetImage(state.monsterName)
	if v.animationQueue == nil {
		v.animationQueue = make(chan CombatAnimation, 20)
	}
	if mainCombatMenu == nil {
		mainCombatMenu = NewMainCombatMenu(v)
	}
	if combatLogSubState == nil {
		combatLogSubState = NewCombatLogSubState(v)
	}
	// v.MenuSubState = mainCombatMenu

	attackSprites, _ := imgPool.GetImage("anim/attack2.png")
	hpIcon, _ := imgPool.GetImage("icon/heart.png")
	v.oppHp = state.monsterHp
	v.Characters = state.CombatCharacters
	hpIconparam := core.NewMovableImageParams().WithMoveParam(core.MoveParam{Sx: (768 / 2) - (512 / 8), Sy: 182})
	v.hpIcon = core.NewMovableImage(hpIcon, hpIconparam)
	attackAnimParam := core.NewMovableImageParams().WithMoveParam(core.MoveParam{Sx: float64(OPP_POS_X), Sy: float64(OPP_POS_Y)}).WithScale(&core.ScaleParam{Sx: 0.25, Sy: 0.25})
	v.attackAnim = &core.AnimatedImage{
		MovableImage:   core.NewMovableImage(attackSprites, attackAnimParam),
		SubImageStartX: 0,
		SubImageStartY: 0,
		SubImageWidth:  512,
		SubImageHeight: 512,
		FrameCount:     16,
		Modulo:         5,
	}
	v.attackAnim.Done = func() {
		DoneAnim(v)
	}
	v.SwitchMenuSubstate(mainCombatMenu)
	bgImg, _ := imgPool.GetImage(state.backgroundName)
	bgImageParam := core.NewMovableImageParams()
	v.background = core.NewMovableImage(bgImg, bgImageParam)
	oppImageParam := core.NewMovableImageParams().WithMoveParam(core.MoveParam{Sx: float64(OPP_POS_X), Sy: float64(OPP_POS_Y)}).
		WithScale(&core.ScaleParam{Sx: 0.25, Sy: 0.25})
	v.opp = core.NewMovableImage(img, oppImageParam)

}
func (g *CombatScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// return 768, 512
	return 768, 512
}
func (v *CombatScene) Unload() MyState {
	return v.state
}
func (v *CombatScene) BeginCombat() {
	v.SwitchMenuSubstate(combatLogSubState)
}
func (s *CombatScene) Update() error {
	s.counter = (s.counter + 1) % 10000
	// s.opp.Update()

	if inpututil.IsKeyJustReleased(ebiten.KeySpace) && !s.onAttackAnim && !s.onIsBlinking {
		// s.onAttackAnim = true
		// s.attackAnim.Done = func() {
		// 	s.onAttackAnim = false
		// 	s.onIsBlinking = true
		// 	s.counter = 0
		// }

		// s.animationQueue <- &AttackAnim{cs: s, doneFunc: DoneAnim}

		cplx_anim := &ComplexAnim{cs: s, doneFunc: DoneAnim, animations: []CombatAnimation{
			&MoveAttackAnim{cs: s, PosX: float64(PLAYER_POS_X_START), PosY: float64(PLAYER_POS_Y_START - 30)},
			&AttackAnim{cs: s},
			// &BlinkAnim{cs: s, blinkCount: 5},
		},
		}
		// s.attackAnim.Done = func() {
		// 	cplx_anim.idx += 1
		// }
		s.animationQueue <- cplx_anim
	}
	if s.currentAnim == nil {
		select {
		case aa := <-s.animationQueue:
			s.currentAnim = aa
		default:
			break
		}
	}
	if s.currentAnim != nil {
		s.currentAnim.Update()
	}
	if s.MenuSubState != nil {
		s.MenuSubState.Update()
	}
	// if s.onAttackAnim {
	// 	s.attackAnim.Update()
	// }
	return nil
}
func (s *CombatScene) DrawCharacter(screen *ebiten.Image) {
	curFont := core.DefaultFont

	for idx, c := range s.Characters {
		box := ebiten.NewImage(160, 80)
		box.Fill(color.RGBA{0, 169, 0, 255})
		opt := ebiten.DrawImageOptions{}
		boxPosX := float64(PLAYER_POS_X_START + idx*170)
		text.Draw(box, c.Name, curFont, 0, 15, color.White)
		text.Draw(box, "HP "+strconv.Itoa(c.HP), curFont, 0, 40, color.White)
		opt.GeoM.Translate(boxPosX, 450)
		if idx == s.CurrentChrIdx {
			box2 := ebiten.NewImage(170, 100)
			box2.Fill(color.RGBA{169, 169, 169, 255})
			opt := ebiten.DrawImageOptions{}
			boxPosX := float64(15 + idx*170)
			opt.GeoM.Translate(boxPosX, float64(PLAYER_POS_Y_START))
			screen.DrawImage(box2, &opt)
		}
		screen.DrawImage(box, &opt)
	}

}
func (s *CombatScene) Draw(screen *ebiten.Image) {
	// your draw code
	s.background.Draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("OnBlink %s\nCounter %d", s.onIsBlinking, s.counter))

	s.hpIcon.Draw(screen)
	text.Draw(screen, strconv.Itoa(s.oppHp), core.DefaultFont, (768/2)-(512/8)+50, 220, color.White)
	// draw combat character
	s.DrawCharacter(screen)
	// draw combat button
	if s.MenuSubState != nil {
		s.MenuSubState.Draw(screen)
	}
	if s.currentAnim != nil {
		s.currentAnim.Draw(screen)
	} else {
		s.opp.Draw(screen)

	}
}

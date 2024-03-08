package main

import (
	"fmt"
	"github/kharism/hanashi/core"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/joelschutz/stagehand"
)

type CombatScene struct {
	state          MyState
	background     *core.MovableImage
	opp            *core.MovableImage
	attackAnim     *core.AnimatedImage
	onAttackAnim   bool
	onIsBlinking   bool
	counter        int
	blinkCount     int
	currentAnim    CombatAnimation
	animationQueue chan CombatAnimation
}

type CombatAnimation interface {
	Draw(screen *ebiten.Image)
	Update()
	GetDoneFunc() func(cs *CombatScene)
	SetDoneFunc(func(cs *CombatScene))
}

func DoneAnim(cs *CombatScene) {
	cs.currentAnim = nil
}
func (v *CombatScene) Load(state MyState, sm *stagehand.SceneManager[MyState]) {
	v.state = state
	v.counter = 0
	img, _ := imgPool.GetImage(state.monsterName)
	if v.animationQueue == nil {
		v.animationQueue = make(chan CombatAnimation, 20)
	}
	attackSprites, _ := imgPool.GetImage("anim/attack2.png")

	v.attackAnim = &core.AnimatedImage{
		MovableImage:   core.NewMovableImage(attackSprites, core.MoveParam{Sx: (768 / 2) - (512 / 8), Sy: 80}, core.ScaleParam{Sx: 0.25, Sy: 0.25}, nil),
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
	bgImg, _ := imgPool.GetImage(state.backgroundName)
	v.background = core.NewMovableImage(bgImg, core.MoveParam{}, core.ScaleParam{Sx: 1, Sy: 1}, nil)
	v.opp = core.NewMovableImage(img, core.MoveParam{Sx: (768 / 2) - (512 / 8), Sy: 80}, core.ScaleParam{Sx: 0.25, Sy: 0.25}, nil)
}
func (g *CombatScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// return 768, 512
	return 768, 512
}
func (v *CombatScene) Unload() MyState {
	return v.state
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
			&AttackAnim{cs: s},
			&BlinkAnim{cs: s, blinkCount: 5},
		},
		}
		s.attackAnim.Done = func() {
			cplx_anim.idx += 1
		}
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

	// if s.onAttackAnim {
	// 	s.attackAnim.Update()
	// }
	return nil
}

func (s *CombatScene) Draw(screen *ebiten.Image) {
	// your draw code
	// s.background.Draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("OnBlink %s\nCounter %d", s.onIsBlinking, s.counter))
	if s.currentAnim != nil {
		s.currentAnim.Draw(screen)
	} else {
		s.opp.Draw(screen)
	}
}

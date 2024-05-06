package main

import (
	"log"
	"os"

	"github.com/kharism/hanashi/core"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
)

var imgPool ImagePool
var sm *stagehand.SceneDirector[MyState]
var gameoverScene *GameOverScene

func init() {
	imgPool = ImagePool{Map: map[string]*ebiten.Image{}}
	Characters = map[string]*core.Character{}
}

const (
	Trigger1 stagehand.SceneTransitionTrigger = iota
	Trigger2
)

type MyState struct {
	monsterName    string
	backgroundName string
	CustomData     map[string]any

	CombatCharacters []*CombatCharacter
	monsterHp        int
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	// game := NewGame(GetScene())
	scene1 := GetScene()
	combatScene := &CombatScene{}
	state := MyState{}
	state.monsterName = "opp/slime.png"
	state.backgroundName = "bg/alley.png"
	state.monsterHp = 20
	state.CombatCharacters = []*CombatCharacter{
		NewCombatCharacter("sven", 0, 9).SetAtkDamage(func() int {
			return Dice(1, 6)
		}),
		NewCombatCharacter("shizuku", 1, 9).SetAtkDamage(func() int {
			return Dice(1, 5)
		}),
	}
	trans := stagehand.NewSlideTransition[MyState](stagehand.LeftToRight, 0.05)

	runScene := RunScene1(combatScene)
	gameoverScene = &GameOverScene{}
	rs := map[stagehand.Scene[MyState]][]stagehand.Directive[MyState]{
		scene1: []stagehand.Directive[MyState]{
			stagehand.Directive[MyState]{Dest: combatScene, Trigger: Trigger1, Transition: trans},
			stagehand.Directive[MyState]{Dest: runScene, Trigger: Trigger2},
		},
		runScene: []stagehand.Directive[MyState]{
			stagehand.Directive[MyState]{Dest: combatScene, Trigger: Trigger1, Transition: trans},
		},
		combatScene: []stagehand.Directive[MyState]{
			stagehand.Directive[MyState]{
				Dest:       gameoverScene,
				Trigger:    Trigger1,
				Transition: trans,
			},
		},
	}
	sm = stagehand.NewSceneDirector[MyState](scene1, state, rs) //stagehand.NewSceneManager[MyState](scene1, state)

	// when combat is go to scene
	combatScene.DoneCombat = func() {
		os.Exit(0)
	}
	if err := ebiten.RunGame(sm); err != nil {
		log.Fatal(err)
	}
}

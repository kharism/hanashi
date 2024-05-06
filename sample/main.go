package main

import (
	"log"
	"os"

	"github.com/kharism/hanashi/core"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
)

var imgPool ImagePool
var sm *stagehand.SceneManager[MyState]
var gameoverScene *GameOverScene

func init() {
	imgPool = ImagePool{Map: map[string]*ebiten.Image{}}
	Characters = map[string]*core.Character{}
}

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
	// runScene := RunScene1(combatScene, sm)
	sm = stagehand.NewSceneManager[MyState](scene1, state)
	gameoverScene = &GameOverScene{}
	// set Done function to tell the scene what to do after
	scene1.Done = func() {
		// check whether user decided to fight or not
		if scene1.GetSceneData("Fight it?").(string) == "yes" {
			scene1.StateDecorator = func(ms MyState) MyState {
				ms.monsterName = "opp/slime.png"
				ms.backgroundName = "bg/alley.png"
				ms.monsterHp = 10
				ms.CombatCharacters = []*CombatCharacter{NewCombatCharacter("sven", 9, 9).SetAtkDamage(func() int {
					return Dice(1, 6)
				})}
				return ms
			}
			sm.SwitchWithTransition(combatScene, stagehand.NewFadeTransition[MyState](0.05))
		} else {
			runScene := RunScene1(combatScene, sm)
			sm.SwitchTo(runScene)
		}

	}
	// when combat is go to scene
	combatScene.DoneCombat = func() {
		os.Exit(0)
	}
	if err := ebiten.RunGame(sm); err != nil {
		log.Fatal(err)
	}
}

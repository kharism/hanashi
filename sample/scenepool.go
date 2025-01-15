package main

import (
	"image/color"
	_ "image/jpeg"

	"github.com/kharism/hanashi/core"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
)

type CustomScene struct {
	*core.Scene
	state          MyState
	StateDecorator StateDecorator
	sm             *stagehand.SceneDirector[MyState]
}
type StateDecorator func(MyState) MyState

func (v *CustomScene) Load(state MyState, sm stagehand.SceneController[MyState]) {
	v.sm = sm.(*stagehand.SceneDirector[MyState])
	v.Scene.Events[0].Execute(v.Scene)
}
func (v *CustomScene) Unload() MyState {
	if v.StateDecorator != nil {
		return v.StateDecorator(v.state)
	}
	return v.state
}
func (g *CustomScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// return 768, 512
	return 768, 512
}
func (g *CustomScene) GetLayout() (screenWidth, screenHeight int) {
	return 768, 512
}
func (g *CustomScene) GetNamePosition() (x, y int) {
	return 0, 512 - 100
}
func (g *CustomScene) GetTextPosition() (x, y int) {
	return 0, 512 - 70
}

// an example how we structure a scene. We have Characters which we going to use
// then an order of events to tell scene what to do.
func GetScene() *CustomScene {
	h := CustomScene{Scene: core.NewScene()}
	h.Characters = []*core.Character{
		core.NewCharacter("sven", "chr/8009774b-2341-4b31-b63d-e172b525841e.png", &imgPool),
		core.NewCharacter("shizuku", "chr/ec40a408-8263-4a26-9738-e02370d3280e.png", &imgPool),
		core.NewCharacter("slime", "opp/slime.png", &imgPool),
	}

	h.Scene.Events = []core.Event{
		// 0th event. Set up bg
		core.NewBGChangeEventFromPath("bg/livingroom.png", core.MoveParam{0, 0, -30, 0, 1}, &imgPool, &core.ShaderParam{ShaderName: core.GRAYSCALE_SHADER}),
		// 1st event, create complex event compromised of add character and dialogue
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterAddEvent("sven", core.MoveParam{-100, 200, 0, 200, 10}, &core.ScaleParam{0.5, 0.5, 0, 0}),
			// core.NewCharacterAddShaderEvent("sven", &core.ShaderParam{ShaderName: core.DARKER_SHADER}),
			core.NewDialogueEvent("sven", "(My living room)", nil),
		}},
		// 2nd event, we don't need to add new character. Just place new dialogue
		core.NewDialogueEvent("sven", "(Finally I can lay back after some office\nonline meeting )", nil),
		// 3rd event, we add new character moving from right to left into view then make sven darker using shader
		// then put in dialogue. All in single click
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterAddEvent("shizuku", core.MoveParam{800, 200, 480, 200, 10}, &core.ScaleParam{0.75, 0.75, 0, 0}),
			core.NewCharacterAddShaderEvent("sven", &core.ShaderParam{ShaderName: core.DARKER_SHADER}),
			core.NewDialogueEvent("shizuku", "Honey, What do you want for dinner?", nil),
		}},
		core.NewDialogueEvent("shizuku", "curry or spaghetti?", nil),
		// 5th event. Add in options for user selection. "What's for dinner?" is used as both
		// question and the options selected will be stored in State.StateData["What's for dinner?"]
		core.NewOptionSelectEvent("dinner", "What's for dinner?", "curry", "spaghetti"),
		// 6th event. Show dialogue based on answer of "What's for dinner?". We don't need to access the answer imemdiately
		// As long as the scenedata is passed on from scene to scene we can pick up what user answers from
		// previous scene
		core.NewBranchingDialogueEvent("shizuku", "dinner", map[string]string{
			"curry":     "We need more potatoes and carrot.",
			"spaghetti": "We need more garlic and basil",
		}),
		core.NewDialogueEvent("shizuku", "Can you get more of them?", nil),
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterAddShaderEvent("shizuku", &core.ShaderParam{ShaderName: core.DARKER_SHADER}),
			core.NewCharacterAddShaderEvent("sven", nil), //remove DARKER_SHADER on sven
			core.NewDialogueEvent("sven", "Sure sweetie, this cold weather is nothing", nil),
		}},
		&core.ComplexEvent{Events: []core.Event{
			// core.NewCharacterAddEvent("shizuku", &core.MoveParam{800, 200, 480, 200, 10}, &core.ScaleParam{0.75, 0.75}),
			core.NewCharacterAddShaderEvent("sven", &core.ShaderParam{ShaderName: core.DARKER_SHADER}),
			core.NewDialogueEvent("shizuku", "The store has moved to new location", nil),
		}},
		//10th
		core.NewDialogueEvent("shizuku", "Do you know the new location?", nil),
		core.NewOptionSelectEvent("Location", "Ask for new store location?", "yes", "no"),
		core.NewBranchingJumpEvent("Location", map[string]int{
			"yes": 13,
			"no":  14,
		}),
		&core.ComplexEvent{Events: []core.Event{
			// core.NewCharacterAddEvent("shizuku", &core.MoveParam{800, 200, 480, 200, 10}, &core.ScaleParam{0.75, 0.75}),
			core.NewCharacterAddShaderEvent("sven", &core.ShaderParam{ShaderName: core.DARKER_SHADER}),
			core.NewDialogueEvent("shizuku", "From the old location go north on the 2nd alley\non the left", nil),
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterAddShaderEvent("shizuku", &core.ShaderParam{ShaderName: core.DARKER_SHADER}),
			core.NewCharacterAddShaderEvent("sven", nil), //remove DARKER_SHADER on sven
			core.NewDialogueEvent("sven", "Alright I'm going", nil),
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("shizuku"),
			core.NewBGChangeEventFromPath("bg/snowvillage.jpg", core.MoveParam{0, 0, 0, 0, 1}, &imgPool, nil),
			core.NewDialogueEvent("sven", "(me and my big mouth, this weather is more chilling\nthan I expected)", nil),
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewDialogueEvent("sven", "(what's this?)", nil),
			core.NewCharacterAddEvent("slime", core.MoveParam{200, -400, 200, 50, 10}, &core.ScaleParam{0.75, 0.75, 0, 0}),
		}},
		core.NewDialogueEvent("sven", "(wait, it's attacking me?)", nil),
		core.NewOptionSelectEvent("Fight", "Fight it?", "yes", "no"),
		//core.NewDialogueEvent("sven", "(A room, Where we do our living)", nil),

	}

	txtBgImage := ebiten.NewImage(768, 300)
	txtBgImage.Fill(color.NRGBA{0, 0, 255, 255})
	h.TxtBg = txtBgImage
	h.SetLayouter(&h)
	h.Done = func() {
		if h.GetSceneData("Fight").(string) == "yes" {
			h.StateDecorator = func(ms MyState) MyState {
				ms.monsterName = "opp/slime.png"
				ms.backgroundName = "bg/alley.png"
				ms.monsterHp = 10
				ms.CombatCharacters = []*CombatCharacter{NewCombatCharacter("sven", 9, 9).SetAtkDamage(func() int {
					return Dice(1, 6)
				})}
				return ms
			}
			h.sm.ProcessTrigger(Trigger1)
		} else {
			h.sm.ProcessTrigger(Trigger2)
		}

	}
	return &h
}
func RunScene1(cs *CombatScene) *CustomScene {
	h := CustomScene{Scene: core.NewScene()}
	h.Characters = []*core.Character{
		core.NewCharacter("sven", "chr/8009774b-2341-4b31-b63d-e172b525841e.png", &imgPool),
		core.NewCharacter("slime", "opp/slime.png", &imgPool),
	}
	h.Scene.Events = []core.Event{
		&core.ComplexEvent{Events: []core.Event{
			core.NewBGChangeEventFromPath("bg/snowvillage.jpg", core.MoveParam{0, 0, 0, 0, 1}, &imgPool, nil),
			core.NewCharacterAddEvent("sven", core.MoveParam{0, 200, 0, 200, 10}, &core.ScaleParam{0.5, 0.5, 0, 0}),
			core.NewCharacterAddEvent("slime", core.MoveParam{200, 50, 200, 50, 10}, &core.ScaleParam{0.75, 0.75, 0, 0}),
			// core.NewCharacterAddShaderEvent("sven", &core.ShaderParam{ShaderName: core.DARKER_SHADER}),
			core.NewDialogueEvent("sven", "(I should run for now)", nil),
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewDialogueEvent("", "", nil), //empty the curDialog buffer
			core.NewCharacterMoveEvent("slime", core.MoveParam{200, 50, 800, 50, 50}),
		}},

		core.NewDialogueEvent("sven", "(*pant* is it still after me?)", nil),
	}
	txtBgImage := ebiten.NewImage(768, 300)
	txtBgImage.Fill(color.NRGBA{0, 0, 255, 255})
	h.TxtBg = txtBgImage
	h.SetLayouter(&h)
	h.StateDecorator = func(ms MyState) MyState {
		ms.monsterName = "opp/slime.png"
		ms.backgroundName = "bg/alley.png"
		ms.monsterHp = 10
		ms.CombatCharacters = []*CombatCharacter{NewCombatCharacter("sven", 9, 9).SetAtkDamage(func() int {
			return Dice(1, 6)
		})}
		return ms
	}
	h.Done = func() {
		// sceneManager.SwitchWithTransition(cs, stagehand.NewFadeTransition[MyState](0.05))
		h.sm.ProcessTrigger(Trigger1)
	}
	return &h
}

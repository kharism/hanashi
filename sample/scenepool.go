package main

import (
	"github/kharism/hanashi/core"
	"image/color"
	_ "image/jpeg"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
)

type CustomScene struct {
	*core.Scene
	state          MyState
	StateDecorator StateDecorator
}
type StateDecorator func(MyState) MyState

func (v *CustomScene) Load(state MyState, sm *stagehand.SceneManager[MyState]) {
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
func GetScene() *CustomScene {
	h := CustomScene{Scene: core.NewScene()}
	h.Scene.Events = []core.Event{
		core.NewBGChangeEventFromPath("bg/livingroom.png", core.MoveParam{0, 0, -30, 0, 1}, &imgPool, &core.ShaderParam{ShaderName: core.GRAYSCALE_SHADER}),
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterAddEvent("sven", &core.MoveParam{-100, 200, 0, 200, 10}, &core.ScaleParam{0.5, 0.5}),
			// core.NewCharacterAddShaderEvent("sven", &core.ShaderParam{ShaderName: core.DARKER_SHADER}),
			core.NewDialogueEvent("sven", "(My living room)", nil),
		}},
		core.NewDialogueEvent("sven", "(Finally I can lay back after some office online \n meeting )", nil),
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterAddEvent("shizuku", &core.MoveParam{800, 200, 480, 200, 10}, &core.ScaleParam{0.75, 0.75}),
			core.NewCharacterAddShaderEvent("sven", &core.ShaderParam{ShaderName: core.DARKER_SHADER}),
			core.NewDialogueEvent("shizuku", "Honey, What do you want for dinner?", nil),
		}},
		core.NewDialogueEvent("shizuku", "curry or spaghetti?", nil),
		core.NewOptionSelectEvent("What's for dinner?", "curry", "spaghetti"),
		core.NewBranchingDialogueEvent("shizuku", "What's for dinner?", map[string]string{
			"curry":     "We need more potatoes.",
			"spaghetti": "We need more garlic",
		}),
		core.NewDialogueEvent("shizuku", "Can you get more of those?", nil),
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterAddShaderEvent("shizuku", &core.ShaderParam{ShaderName: core.DARKER_SHADER}),
			core.NewCharacterAddShaderEvent("sven", nil), //remove DARKER_SHADER on sven
			core.NewDialogueEvent("sven", "Sure sweetie, this cold weather is nothing", nil),
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("shizuku"),
			core.NewBGChangeEventFromPath("bg/snowvillage.jpg", core.MoveParam{0, 0, 0, 0, 1}, &imgPool, nil),
			core.NewDialogueEvent("sven", "(me and my big mouth, this weather is more chilling than I expected)", nil),
		}},
		core.NewDialogueEvent("sven", "(what's this?)", nil),
		core.NewCharacterAddEvent("slime", &core.MoveParam{200, -400, 200, 50, 10}, &core.ScaleParam{0.75, 0.75}),
		core.NewDialogueEvent("sven", "(wait, it's attacking me?)", nil),
		core.NewOptionSelectEvent("Fight it?", "yes", "no"),
		//core.NewDialogueEvent("sven", "(A room, Where we do our living)", nil),

	}
	h.Characters = []*core.Character{
		core.NewCharacter("sven", "chr/8009774b-2341-4b31-b63d-e172b525841e.png", &imgPool),
		core.NewCharacter("shizuku", "chr/ec40a408-8263-4a26-9738-e02370d3280e.png", &imgPool),
		core.NewCharacter("slime", "opp/slime.png", &imgPool),
	}
	txtBgImage := ebiten.NewImage(768, 300)
	txtBgImage.Fill(color.NRGBA{0, 0, 255, 255})
	h.TxtBg = txtBgImage
	h.SetLayouter(&h)
	return &h
}
func RunScene1(cs *CombatScene, sceneManager *stagehand.SceneManager[MyState]) *CustomScene {
	h := CustomScene{Scene: core.NewScene()}
	h.Characters = []*core.Character{
		core.NewCharacter("sven", "chr/8009774b-2341-4b31-b63d-e172b525841e.png", &imgPool),
		core.NewCharacter("slime", "opp/slime.png", &imgPool),
	}
	h.Scene.Events = []core.Event{
		&core.ComplexEvent{Events: []core.Event{
			core.NewBGChangeEventFromPath("bg/snowvillage.jpg", core.MoveParam{0, 0, 0, 0, 1}, &imgPool, nil),
			core.NewCharacterAddEvent("sven", &core.MoveParam{0, 200, 0, 200, 10}, &core.ScaleParam{0.5, 0.5}),
			core.NewCharacterAddEvent("slime", &core.MoveParam{200, 50, 200, 50, 10}, &core.ScaleParam{0.75, 0.75}),
			// core.NewCharacterAddShaderEvent("sven", &core.ShaderParam{ShaderName: core.DARKER_SHADER}),
			core.NewDialogueEvent("sven", "(I should run for now)", nil),
		}},
		core.NewCharacterMoveEvent("slime", core.MoveParam{200, 50, 800, 50, 50}),
		core.NewDialogueEvent("sven", "(*pant* is it still after me?)", nil),
	}
	txtBgImage := ebiten.NewImage(768, 300)
	txtBgImage.Fill(color.NRGBA{0, 0, 255, 255})
	h.TxtBg = txtBgImage
	h.SetLayouter(&h)
	h.StateDecorator = func(ms MyState) MyState {
		ms.monsterName = "opp/slime.png"
		ms.backgroundName = "bg/alley.png"
		return ms
	}
	h.Done = func() {
		sceneManager.SwitchWithTransition(cs, stagehand.NewFadeTransition[MyState](0.05))
	}
	return &h
}

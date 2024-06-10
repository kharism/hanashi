package main

import (
	"image/color"
	"os"

	"github.com/kharism/hanashi/core"

	"github.com/hajimehoshi/ebiten/v2"
)

func Scene1(layouter core.GetLayouter) *core.Scene {
	scene := core.NewScene()
	scene.Characters = []*core.Character{
		core.NewCharacter("sven", "../sample/chr/8009774b-2341-4b31-b63d-e172b525841e.png", &imgPool),
	}
	scene.Events = []core.Event{
		core.NewBGChangeEventFromPath("../sample/bg/village.png", core.MoveParam{Sx: 0, Sy: 0, Tx: 0, Ty: -200, Speed: 1}, &imgPool, nil),
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterAddEvent("sven", &core.MoveParam{-100, 200, 0, 200, 10}, &core.ScaleParam{0.4, 0.4, 0, 0}),
			// core.NewCharacterAddShaderEvent("sven", &core.ShaderParam{ShaderName: core.DARKER_SHADER}),
			core.NewDialogueEvent("sven", "(What a wonderful scenery)", nil),
		}},
		core.NewDialogueEvent("sven", "(I still have time before dusk to find a way home)", nil),
	}
	scene.Events[0].Execute(scene)
	txtBgImage := ebiten.NewImage(768, 300)
	txtBgImage.Fill(color.NRGBA{0, 0, 255, 255})
	scene.TxtBg = txtBgImage
	scene.SetLayouter(layouter)
	scene.Done = func() {
		os.Exit(0)
	}
	return scene
}

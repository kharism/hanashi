package main

import (
	"github/kharism/hanashi/core"
	"image/color"
	_ "image/jpeg"

	"github.com/hajimehoshi/ebiten/v2"
)

func GetScene() *core.Scene {
	h := core.Scene{}

	h.Events = []core.Event{
		core.NewBGChangeEventFromPath("bg/livingroom.png", core.MoveParam{0, 0, -30, 0, 1}, &imgPool, &core.ShaderParam{ShaderName: core.SEPIA_SHADER}),
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterAddEvent("sven", &core.MoveParam{-100, 200, 0, 200, 10}, &core.ScaleParam{0.5, 0.5}),
			// core.NewCharacterAddShaderEvent("sven", &core.ShaderParam{ShaderName: core.DARKER_SHADER}),
			core.NewDialogueEvent("sven", "(My living room)", nil),
		}},
		core.NewDialogueEvent("sven", "(A room, Where we do our living)", nil),
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterAddEvent("shizuku", &core.MoveParam{800, 200, 480, 200, 10}, &core.ScaleParam{0.75, 0.75}),
			core.NewCharacterAddShaderEvent("sven", &core.ShaderParam{ShaderName: core.DARKER_SHADER}),
			core.NewDialogueEvent("shizuku", "Honey, can you get some milk and egg?", nil),
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterAddShaderEvent("shizuku", &core.ShaderParam{ShaderName: core.DARKER_SHADER}),
			core.NewCharacterAddShaderEvent("sven", nil), //remove DARKER_SHADER on sven
			core.NewDialogueEvent("sven", "Under this cold weather? sure", nil),
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("shizuku"),
			core.NewBGChangeEventFromPath("bg/snowvillage.jpg", core.MoveParam{0, 0, 0, 0, 1}, &imgPool, nil),
			core.NewDialogueEvent("", "", nil),
		}},
		//core.NewDialogueEvent("sven", "(A room, Where we do our living)", nil),

	}
	h.Characters = []*core.Character{
		core.NewCharacter("sven", "chr/8009774b-2341-4b31-b63d-e172b525841e.png", &imgPool),
		core.NewCharacter("shizuku", "chr/ec40a408-8263-4a26-9738-e02370d3280e.png", &imgPool),
	}
	txtBgImage := ebiten.NewImage(768, 300)
	txtBgImage.Fill(color.NRGBA{0, 0, 255, 255})
	h.TxtBg = txtBgImage
	return &h
}

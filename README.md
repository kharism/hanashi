# HANASHI

this package provide some limited tools to convey story on ebitengine. As for now it only handle VN-like story telling with some way to handle branching options. Look up on sample directory for some example on how to use this package


## Installation

```
go get github.com/kharism/hanashi
```

## How to use

core.Scene is basically an object we use to store anything to do with the story we want to covey.
```
scene := core.NewScene()
```
each scene will have its own set of characters,event, and scene data. We basically put scene in our Update() and Draw(*ebiten.Image) function

```
func NewScene() *Scene {
	return &Scene{Events: []Event{}, Characters: []*Character{}, SceneData: map[string]any{}}
}
```
Characters is characters we want to move around in a scene. We need to register each characters first before we can use them to tell a dialogue or monologue.

```
scene.Characters = []*core.Character{
		core.NewCharacter("sven", "../sample/chr/8009774b-2341-4b31-b63d-e172b525841e.png", &imgPool),
	}
```
After that we set events as if we write a screenplay
```
scene.Events = []core.Event{
		core.NewBGChangeEventFromPath("../sample/bg/village.png", core.MoveParam{Sx: 0, Sy: 0, Tx: 0, Ty: -200, Speed: 1}, &imgPool, nil),
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterAddEvent("sven", &core.MoveParam{-100, 200, 0, 200, 10}, &core.ScaleParam{0.4, 0.4}),
			// core.NewCharacterAddShaderEvent("sven", &core.ShaderParam{ShaderName: core.DARKER_SHADER}),
			core.NewDialogueEvent("sven", "(What a wonderful scenery)", nil),
		}},
		core.NewDialogueEvent("sven", "(I still have time before dusk to find a way home)", nil),
	}
```
after that we execute the 1st event
```
scene.Events[0].Execute(scene)
```
Set background of the text and set layouter
```
txtBgImage := ebiten.NewImage(768, 300)
txtBgImage.Fill(color.NRGBA{0, 0, 255, 255})
scene.TxtBg = txtBgImage
scene.SetLayouter(layouter)
```
after that we set Done Callback
```
scene.Done = func() {
		os.Exit(0)
	}
```
The full example of this example is in sample2 directory. Sample directory is for more complex example where we use multiple scene, have scene transition, and dialogue option 
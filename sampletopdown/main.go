package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kharism/hanashi/core"
	"github.com/kharism/hanashi/sampletopdown/components"
	"github.com/kharism/hanashi/sampletopdown/system"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"

	_ "embed"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
)

const mapPath = "test.tmx"

var gameMap *tiled.Map

//go:embed assets/PixelOperator8.ttf
var PixelFontTTF []byte

//go:embed assets/player.png
var player []byte

var PixelFont *text.GoTextFaceSource
var face *text.GoTextFace
var PlayerImg *ebiten.Image

const tileSize = 16

type GridPos struct {
	Col int
	Row int
}
type Game struct {
	bg          *ebiten.Image
	tiledLayout *tiled.Map
	PlayerPos   GridPos
	ecs         *ecs.ECS
	Scene       *core.Scene
}

func (g *Game) Update() error {
	// g.Scene.Update()
	if g.Scene == nil {
		g.ecs.Update()
	} else {
		g.Scene.Update()
	}

	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	g.ecs.DrawLayer(0, screen)
	g.ecs.DrawLayer(1, screen)
	if g.Scene != nil {
		g.Scene.Draw(screen)
	}
}
func (g *Game) Layout(width, height int) (int, int) {
	return 15 * tileSize, 11 * tileSize
}

// return width and height of the scene
func (g *Game) GetLayout() (width, height int) {
	return 15 * tileSize, 11 * tileSize
}

// return the starting text position where the box containing name of the character appear on the scene
// return negative number if no such box needed
func (g *Game) GetNamePosition() (x, y int) {
	return 0, 6 * tileSize
}

// get the starting position of the text
func (g *Game) GetTextPosition() (x, y int) {
	return 0, 7 * tileSize
}
func (g *Game) SetHanashiScene(scene *core.Scene) {
	scene.SetLayouter(g)
	scene.VisibleDialog = ""
	scene.TxtBg = ebiten.NewImage(1024-128, 128)
	scene.TxtBg.Fill(color.RGBA{R: 0x4f, G: 0x8f, B: 0xba, A: 255})
	g.Scene = scene
}
func (g *Game) UnsetHanashiScene() {
	g.Scene = nil
}
func main() {
	var err error
	// gameMap is 10x10 grid each grid is 16x16 pixel
	// it has several layers, the 0th layer define wall and ground. 1st layer define interactibles
	gameMap, err = tiled.LoadFile(mapPath)
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}
	// this functions init how player will interact with interactibles
	// for now this is hardcoded
	InitInteactible()

	core.DetectKeyboardNext = func() bool {
		return inpututil.IsKeyJustReleased(ebiten.KeySpace)
	}

	renderer, err := render.NewRenderer(gameMap)
	if err != nil {
		fmt.Printf("map unsupported for rendering: %s", err.Error())
		os.Exit(2)
	}
	err = renderer.RenderLayer(0)
	if err != nil {
		fmt.Printf("layer unsupported for rendering: %s", err.Error())
		os.Exit(2)
	}
	// Get a reference to the Renderer's output, an image.NRGBA struct.
	img := renderer.Result

	ebiBg := ebiten.NewImageFromImage(img)

	PlayerImg, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(player))

	if err != nil {
		fmt.Printf("Error loading player spirte: %s", err.Error())
		os.Exit(2)
	}

	// // Render just layer 0 to the Renderer.
	// err = renderer.RenderLayer(1)
	// if err != nil {
	// 	fmt.Printf("layer unsupported for rendering: %s", err.Error())
	// 	os.Exit(2)
	// }

	s, err := text.NewGoTextFaceSource(bytes.NewReader(PixelFontTTF))
	if err != nil {
		log.Fatal(err)
	}
	PixelFont = s
	face = &text.GoTextFace{
		Source: PixelFont,
		Size:   8,
	}
	ebiten.SetWindowSize(35*tileSize, 30*tileSize)
	ebiten.SetWindowTitle("test")
	game := &Game{}
	game.bg = ebiBg
	world := donburi.NewWorld()
	game.ecs = ecs.NewECS(world)

	player := world.Create(components.GridPos, components.Sprite)
	playerEntry := world.Entry(player)
	gridPos := components.GridPos.Get(playerEntry)
	gridPos.Col = 4
	gridPos.Row = 4
	components.Sprite.Set(playerEntry, PlayerImg)
	bgRenderer := &system.BgRenderer{
		Player:      playerEntry,
		MapRenderer: renderer,
	}
	spriteRenderer := &system.SpriteRenderer{
		Player: playerEntry,
		Query: donburi.NewQuery(
			filter.Contains(
				components.Sprite,
				components.GridPos,
			),
		),
	}
	game.ecs.AddRenderer(0, bgRenderer.RenderBg)
	game.ecs.AddRenderer(1, spriteRenderer.RenderBg)
	playerMovement := &system.PlayerMovementSystem{Map: gameMap, Player: playerEntry, SceneSetter: game}
	game.ecs.AddSystem(playerMovement.Update)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

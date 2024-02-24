package main

import (
	"github/kharism/hanashi/core"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Scene *core.Scene
}

func NewGame(Scene *core.Scene) Game {
	g := Game{Scene: Scene}
	Scene.SetLayouter(&g)
	g.Scene.Events[0].Execute(g.Scene)
	return g
}
func (g *Game) Update() error {
	g.Scene.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// ebitenutil.DebugPrint(screen, "Hello, World!")
	g.Scene.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// return 768, 512
	return 768, 512
}
func (g *Game) GetLayout() (screenWidth, screenHeight int) {
	return 768, 512
}
func (g *Game) GetNamePosition() (x, y int) {
	return 0, 512 - 100
}
func (g *Game) GetTextPosition() (x, y int) {
	return 0, 512 - 70
}

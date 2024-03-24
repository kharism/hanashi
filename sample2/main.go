package main

import (
	"github/kharism/hanashi/core"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	Scene *core.Scene
}

func (g *Game) Update() error {
	g.Scene.Update()
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	g.Scene.Draw(screen)
}
func (g *Game) Layout(width, height int) (int, int) {
	return 768, 512
}

// return width and height of the scene
func (g *Game) GetLayout() (width, height int) {
	return 768, 512
}

// return the starting text position where the box containing name of the character appear on the scene
// return negative number if no such box needed
func (g *Game) GetNamePosition() (x, y int) {
	return 0, 512 - 100
}

// get the starting position of the text
func (g *Game) GetTextPosition() (x, y int) {
	return 0, 512 - 70
}

type ImagePool struct {
	Map map[string]*ebiten.Image
}

func (m *ImagePool) GetImage(key string) (*ebiten.Image, error) {
	if _, ok := m.Map[key]; ok {
		return m.Map[key], nil
	}
	img, _, err := ebitenutil.NewImageFromFile(key)
	if err != nil {
		return nil, err
	}
	m.Map[key] = img
	return img, nil
}

var (
	imgPool ImagePool
)

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	imgPool = ImagePool{Map: map[string]*ebiten.Image{}}
	game := &Game{}
	scene := Scene1(game)
	// scene.SetLayouter(game)
	game.Scene = scene
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}

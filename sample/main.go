package main

import (
	"github/kharism/hanashi/core"
	"log"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

var imgPool ImagePool

func init() {
	imgPool = ImagePool{Map: map[string]*ebiten.Image{}}
	Characters = map[string]*core.Character{}
}
func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	game := NewGame(GetScene())

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

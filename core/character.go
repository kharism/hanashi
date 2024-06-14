package core

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

type Character struct {
	Name string
	Img  *MovableImage
}

func NewCharacterImage(name string, image *ebiten.Image) *Character {
	newChar := &Character{Name: name, Img: &MovableImage{image: image, mutex: &sync.Mutex{}}}
	return newChar
}
func NewCharacter(name, imgpath string, loader ImageLoader) *Character {
	newChar := &Character{Name: name}
	if imgpath != "" {
		img, _ := loader.GetImage(imgpath)
		newMovableImage := &MovableImage{image: img, mutex: &sync.Mutex{}}
		newChar.Img = newMovableImage
	}
	return newChar
}
func (c *Character) CheckName(name string) bool {
	return c.Name == name
}
func (c *Character) AddAnimation(ma *MoveAnimation) {
	c.Img.AddAnimation(ma)
}

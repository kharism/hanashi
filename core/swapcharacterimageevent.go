package core

import "github.com/hajimehoshi/ebiten/v2"

type CharacterImageSwapEvent struct {
	Name     string
	NewImage *ebiten.Image
}

func (c *CharacterImageSwapEvent) Execute(scene *Scene) {
	for idx, v := range scene.Characters {
		if v.Name == c.Name {
			scene.Characters[idx].Img.image = c.NewImage
		}
	}
	for idx, v := range scene.ViewableCharacters {
		if v.Name == c.Name {
			scene.ViewableCharacters[idx].Img.image = c.NewImage
		}
	}
}

package core

import "github.com/hajimehoshi/ebiten/v2"

// Replace Image of a character with other image
// usually used to replace an image of a character with other image of that character
// showing different emotion
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

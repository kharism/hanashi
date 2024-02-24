package core

import "github.com/hajimehoshi/ebiten/v2"

type ImageLoader interface {
	GetImage(key string) (*ebiten.Image, error)
}

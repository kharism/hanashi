package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

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

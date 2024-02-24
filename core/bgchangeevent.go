package core

import (
	"sync"

	// "github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/v2"
)

// this event change background of the scene and animate it a bit
type BGChangeEvent struct {
	bg *ebiten.Image
	// start position
	Sx float64
	Sy float64

	// target position
	Tx     float64
	Ty     float64
	Speed  float64
	Shader *ebiten.Shader
}

func NewBGChangeEventFromPath(imagepath string, moveParam MoveParam, imageloader ImageLoader, shader *ShaderParam) Event {
	image, _ := imageloader.GetImage(imagepath)
	if shader == nil {
		return &BGChangeEvent{bg: image, Shader: nil, Sx: moveParam.Sx, Sy: moveParam.Sy, Tx: moveParam.Tx, Ty: moveParam.Ty, Speed: moveParam.Speed}
	}

	if shader.Shader != nil {
		return &BGChangeEvent{bg: image, Shader: shader.Shader, Sx: moveParam.Sx, Sy: moveParam.Sy, Tx: moveParam.Tx, Ty: moveParam.Ty, Speed: moveParam.Speed}
	} else {
		sh, _ := GetShaderPool().GetShader(shader.ShaderName)
		return &BGChangeEvent{bg: image, Shader: sh, Sx: moveParam.Sx, Sy: moveParam.Sy, Tx: moveParam.Tx, Ty: moveParam.Ty, Speed: moveParam.Speed}
	}

}
func (b *BGChangeEvent) Execute(scene *Scene) {
	img := MovableImage{image: b.bg, mutex: &sync.Mutex{}, Shader: b.Shader}
	moveAnim := MoveAnimation{
		tx:    b.Tx,
		ty:    b.Ty,
		Speed: b.Speed,
	}
	img.AddAnimation(&moveAnim)
	scene.CurrentBg = &img
}

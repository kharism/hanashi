package core

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// for image that do animation with frame
type AnimatedImage struct {
	*MovableImage
	SubImageStartX int
	SubImageStartY int

	SubImageWidth  int
	SubImageHeight int

	FrameCount int

	// internal counter. Auto increase every update()
	counter uint
	// Move to next frame everytime counter reach this number
	// smaller number will make the animation faster
	Modulo         int
	FlipHorizontal bool
	Done           func()
}

func (a *AnimatedImage) Update() {
	a.MovableImage.Update()
	a.counter = (a.counter + 1)
	subImageIndex := (int(a.counter) / a.Modulo) % a.FrameCount
	if subImageIndex == a.FrameCount-1 {
		if a.Done != nil {
			a.Done()
		}
		a.counter = 0
	}
}

func (e *AnimatedImage) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	// op.GeoM.Scale(0.25, 0.25)
	op.GeoM.Translate(0, 0)
	if e.ScaleParam != nil {
		op.GeoM.Scale(e.ScaleParam.Sx, e.ScaleParam.Sy)
	}
	subImageIndex := (int(e.counter) / e.Modulo) % e.FrameCount
	subImageStartX := e.SubImageStartX + subImageIndex*e.SubImageWidth
	subImageStartY := e.SubImageStartY

	subImageEndX := e.SubImageStartX + (subImageIndex+1)*e.SubImageWidth
	subImageEndY := e.SubImageStartY + e.SubImageHeight

	subImage := e.image.SubImage(image.Rect(subImageStartX, subImageStartY, subImageEndX, subImageEndY))
	ebitenSubImg := ebiten.NewImageFromImage(subImage)
	if e.FlipHorizontal {
		op.GeoM.Scale(-1, 1)
		bound := subImage.Bounds()
		op.GeoM.Translate(float64(bound.Dx()), 0)
	}
	op.GeoM.Translate(float64(e.x), float64(e.y))
	if e.Shader != nil {
		opts := &ebiten.DrawRectShaderOptions{}
		if e.ScaleParam != nil {
			opts.GeoM.Scale(e.ScaleParam.Sx, e.ScaleParam.Sy)
		}
		opts.GeoM.Translate(float64(e.x), float64(e.y))
		opts.Images[0] = ebitenSubImg
		bounds := ebitenSubImg.Bounds()
		// e.image.DrawRectShader(bounds.Dx(), bounds.Dy(), e.Shader, nil)
		screen.DrawRectShader(bounds.Dx(), bounds.Dy(), e.Shader, opts)
		// screen.DrawImage(e.image, op)
	} else {
		screen.DrawImage(ebitenSubImg, op)
	}

}

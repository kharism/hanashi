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
	Modulo int

	Done func()
}

func (a *AnimatedImage) Update() {
	a.MovableImage.Update()
	a.counter = (a.counter + 1)
	subImageIndex := (int(a.counter) / 5) % a.FrameCount
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
	subImageIndex := (int(e.counter) / 5) % e.FrameCount
	subImageStartX := e.SubImageStartX + subImageIndex*e.SubImageWidth
	subImageStartY := e.SubImageStartY

	subImageEndX := e.SubImageStartX + (subImageIndex+1)*e.SubImageWidth
	subImageEndY := e.SubImageStartY + e.SubImageHeight

	subImage := e.image.SubImage(image.Rect(subImageStartX, subImageStartY, subImageEndX, subImageEndY))
	ebitenSubImg := ebiten.NewImageFromImage(subImage)
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

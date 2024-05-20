package core

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	csg "github.com/kharism/golang-csg/core"
)

// an image that you can move
type MovableImage struct {
	image     *ebiten.Image
	oriWidth  int
	oriHeight int
	//store the current scale of the image
	ScaleParam *ScaleParam
	// current position
	x float64
	y float64
	// velocity of card movement
	vx float64
	vy float64
	// target position if card moved
	tx float64
	ty float64

	//scale target
	tsx float64
	tsy float64
	vsx float64
	vsy float64

	mutex *sync.Mutex
	// animation stuff
	CurrMove       Animation
	AnimationQueue []Animation
	Shader         *ebiten.Shader
}

type MoveParam struct {
	Sx float64
	Sy float64
	Tx float64
	Ty float64

	Speed float64
}

type ScaleParam struct {
	Sx float64
	Sy float64

	// Ty float64
	// Tx float64
}
type MovableImageParams struct {
	MoveParam     MoveParam
	ScaleParam    *ScaleParam
	ShaderOptions *ShaderParam
}

func (e *MovableImage) GetPos() (float64, float64) {
	return e.x, e.y
}

// this function is to immediately move the image to pos x,y
func (e *MovableImage) SetPos(x, y float64) {
	e.x = x
	e.y = y
}

func (e *MovableImage) SetImage(i *ebiten.Image) {
	e.image = i
}
func (e *MovableImage) GetSize() (float64, float64) {
	return float64(e.image.Bounds().Dx()) * e.ScaleParam.Sx, float64(e.image.Bounds().Dy()) * e.ScaleParam.Sy
}
func NewMovableImageParams() *MovableImageParams {
	return &MovableImageParams{ScaleParam: &ScaleParam{Sx: 1, Sy: 1}}
}
func (p *MovableImageParams) WithMoveParam(param MoveParam) *MovableImageParams {
	p.MoveParam = param
	return p
}
func (p *MovableImageParams) WithShader(param *ShaderParam) *MovableImageParams {
	p.ShaderOptions = param
	return p
}
func (p *MovableImageParams) WithScale(param *ScaleParam) *MovableImageParams {
	p.ScaleParam = param
	return p
}
func NewMovableImage(image *ebiten.Image, param *MovableImageParams) *MovableImage {
	mov := &MovableImage{image: image, x: param.MoveParam.Sx, y: param.MoveParam.Sy, ScaleParam: param.ScaleParam, mutex: &sync.Mutex{}}
	if param.ShaderOptions != nil {
		if param.ShaderOptions.Shader != nil {
			mov.Shader = param.ShaderOptions.Shader
		} else {
			mov.Shader, _ = shaderPool.GetShader(param.ShaderOptions.ShaderName)
		}
	}
	if mov.ScaleParam != nil {
		mov.tsx = mov.ScaleParam.Sx
		mov.tsy = mov.ScaleParam.Sy
	}
	return mov
}

// parameter to use shader. You can fill Shader or Shadername and the function
// that takes it will determine wheter to assign shader directly or uses shader
// already registered on shaderpool
type ShaderParam struct {
	Shader     *ebiten.Shader
	ShaderName string
}

// Interface to apply animation to an image
type Animation interface {
	Apply(image *MovableImage)
}

// move animation
type MoveAnimation struct {
	// target x
	tx float64
	// target y
	ty        float64
	Speed     float64
	SleepPre  time.Duration
	SleepPost time.Duration
	DoneFunc  func()
}

func NewMoveAnimationFromParam(param MoveParam) *MoveAnimation {
	return &MoveAnimation{tx: param.Tx, ty: param.Ty, Speed: param.Speed}
}
func (h *MoveAnimation) Apply(e *MovableImage) {
	e.CurrMove = h
	e.tx = h.tx
	e.ty = h.ty
	vx := float64(e.tx - e.x)
	vy := float64(e.ty - e.y)
	if vx != 0 || vy != 0 {
		speedVector := csg.NewVector(vx, vy, 0)
		speedVector = speedVector.Normalize().MultiplyScalar(h.Speed)
		e.vx = speedVector.X
		e.vy = speedVector.Y
	} else {
		e.vx = 0
		e.vy = 0
	}
}

type ScaleAnimation struct {
	// target x
	Tsx float64
	// target y
	Tsy    float64
	SpeedX float64
	SpeedY float64
}

func (s *ScaleAnimation) Apply(img *MovableImage) {
	img.tsx = s.Tsx
	img.tsy = s.Tsy
	img.vsx = s.SpeedX
	img.vsy = s.SpeedY
}

func (e *MovableImage) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	// op.GeoM.Scale(0.25, 0.25)
	op.GeoM.Translate(0, 0)
	if e.ScaleParam != nil {
		op.GeoM.Scale(e.ScaleParam.Sx, e.ScaleParam.Sy)
	}

	op.GeoM.Translate(float64(e.x), float64(e.y))
	if e.Shader != nil {
		opts := &ebiten.DrawRectShaderOptions{}
		if e.ScaleParam != nil {
			opts.GeoM.Scale(e.ScaleParam.Sx, e.ScaleParam.Sy)
		}
		opts.GeoM.Translate(float64(e.x), float64(e.y))
		opts.Images[0] = e.image
		bounds := e.image.Bounds()
		// e.image.DrawRectShader(bounds.Dx(), bounds.Dy(), e.Shader, nil)
		screen.DrawRectShader(bounds.Dx(), bounds.Dy(), e.Shader, opts)
		// screen.DrawImage(e.image, op)
	} else {
		screen.DrawImage(e.image, op)
	}

}

func (e *MovableImage) AddAnimation(animation ...Animation) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.AnimationQueue = append(e.AnimationQueue, animation...)
}

// replace current animation with new one
func (e *MovableImage) ReplaceCurrentAnim(animation *MoveAnimation) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	animation.Apply(e)
	// e.CurrMove = animation
	// e.tx = e.CurrMove.tx
	// e.ty = e.CurrMove.ty
	// vx := float64(e.tx - e.x)
	// vy := float64(e.ty - e.y)
	// if vx != 0 || vy != 0 {
	// 	speedVector := csg.NewVector(vx, vy, 0)
	// 	speedVector = speedVector.Normalize().MultiplyScalar(e.CurrMove.Speed)
	// 	e.vx = speedVector.X
	// 	e.vy = speedVector.Y
	// } else {
	// 	e.vx = 0
	// 	e.vy = 0
	// }
}

func (e *MovableImage) Update() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if e.CurrMove == nil && len(e.AnimationQueue) > 0 {
		e.CurrMove = e.AnimationQueue[0]
		e.AnimationQueue = e.AnimationQueue[1:]
		e.CurrMove.Apply(e)
		// fmt.Println("animation queue", e.card.GetName(), e.CurrMove)
		// if e.CurrMove.SleepPre != 0 {
		// 	time.Sleep(e.CurrMove.SleepPre)
		// }
		// e.tx = e.CurrMove.tx
		// e.ty = e.CurrMove.ty
		// vx := float64(e.tx - e.x)
		// vy := float64(e.ty - e.y)
		// if vx != 0 || vy != 0 {
		// 	speedVector := csg.NewVector(vx, vy, 0)
		// 	speedVector = speedVector.Normalize().MultiplyScalar(e.CurrMove.Speed)
		// 	e.vx = speedVector.X
		// 	e.vy = speedVector.Y
		// } else {
		// 	e.vx = 0
		// 	e.vy = 0
		// }

	}
	e.x += e.vx
	e.y += e.vy
	if e.ScaleParam != nil {
		fmt.Println(e.ScaleParam, e.tsx, math.Abs(e.ScaleParam.Sx-e.tsx))
		if math.Abs(e.ScaleParam.Sx-e.tsx) > 0.01 {
			e.ScaleParam.Sx += e.vsx
		}
		if math.Abs(e.ScaleParam.Sy-e.tsy) >= 0.01 {
			e.ScaleParam.Sy += e.vsy
		}
	}
	// fmt.Println(e.x, e.y)
	if math.Abs(float64(e.tx-e.x))+math.Abs(float64(e.ty-e.y)) < 15 {
		// if e.CurrMove != nil && e.CurrMove.DoneFunc != nil {
		// 	if e.CurrMove.SleepPost != 0 {
		// 		//time.Sleep(e.CurrMove.SleepPost)
		// 	}
		// 	e.CurrMove.DoneFunc()
		// }
		if len(e.AnimationQueue) == 0 {
			e.x = e.tx
			e.y = e.ty
			e.vx = 0
			e.vy = 0
			e.CurrMove = nil
		} else {

			e.CurrMove = e.AnimationQueue[0]
			e.AnimationQueue = e.AnimationQueue[1:]
			e.CurrMove.Apply(e)

			// if e.CurrMove.SleepPre != 0 {
			// 	//time.Sleep(e.CurrMove.SleepPre)
			// }
			// e.tx = e.CurrMove.tx
			// e.ty = e.CurrMove.ty
			// vx := float64(e.tx - e.x)
			// vy := float64(e.ty - e.y)
			// if vy != 0 || vx != 0 {
			// 	speedVector := csg.NewVector(vx, vy, 0)
			// 	speedVector = speedVector.Normalize().MultiplyScalar(e.CurrMove.Speed)
			// 	e.vx = speedVector.X
			// 	e.vy = speedVector.Y
			// } else {
			// 	e.vx = 0
			// 	e.vy = 0
			// }

		}

	}
}

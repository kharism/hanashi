package core

import (
	"math"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	csg "github.com/kharism/golang-csg/core"
)

// an image that you can move
type MovableImage struct {
	image      *ebiten.Image
	oriWidth   int
	oriHeight  int
	ScaleParam *ScaleParam
	// current position
	x float64
	y float64
	// velocity of card movement
	vx float64
	vy float64
	// target position if card moved
	tx    float64
	ty    float64
	mutex *sync.Mutex
	// animation stuff
	CurrMove       *MoveAnimation
	AnimationQueue []*MoveAnimation
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

// parameter to use shader. You can fill Shader or Shadername and the function
// that takes it will determine wheter to assign shader directly or uses shader
// already registered on shaderpool
type ShaderParam struct {
	Shader     *ebiten.Shader
	ShaderName string
}
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

func (e *MovableImage) AddAnimation(animation ...*MoveAnimation) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.AnimationQueue = append(e.AnimationQueue, animation...)
}
func (e *MovableImage) ReplaceCurrentAnim(animation *MoveAnimation) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.CurrMove = animation
	e.tx = e.CurrMove.tx
	e.ty = e.CurrMove.ty
	vx := float64(e.tx - e.x)
	vy := float64(e.ty - e.y)
	if vx != 0 || vy != 0 {
		speedVector := csg.NewVector(vx, vy, 0)
		speedVector = speedVector.Normalize().MultiplyScalar(e.CurrMove.Speed)
		e.vx = speedVector.X
		e.vy = speedVector.Y
	} else {
		e.vx = 0
		e.vy = 0
	}
}

func (e *MovableImage) Update() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if e.CurrMove == nil && len(e.AnimationQueue) > 0 {
		e.CurrMove = e.AnimationQueue[0]
		e.AnimationQueue = e.AnimationQueue[1:]
		// fmt.Println("animation queue", e.card.GetName(), e.CurrMove)
		if e.CurrMove.SleepPre != 0 {
			time.Sleep(e.CurrMove.SleepPre)
		}
		e.tx = e.CurrMove.tx
		e.ty = e.CurrMove.ty
		vx := float64(e.tx - e.x)
		vy := float64(e.ty - e.y)
		if vx != 0 || vy != 0 {
			speedVector := csg.NewVector(vx, vy, 0)
			speedVector = speedVector.Normalize().MultiplyScalar(e.CurrMove.Speed)
			e.vx = speedVector.X
			e.vy = speedVector.Y
		} else {
			e.vx = 0
			e.vy = 0
		}

	}
	e.x += e.vx
	e.y += e.vy
	// fmt.Println(e.x, e.y)
	if math.Abs(float64(e.tx-e.x))+math.Abs(float64(e.ty-e.y)) < 15 {
		if e.CurrMove != nil && e.CurrMove.DoneFunc != nil {
			if e.CurrMove.SleepPost != 0 {
				//time.Sleep(e.CurrMove.SleepPost)
			}
			e.CurrMove.DoneFunc()
		}
		if len(e.AnimationQueue) == 0 {
			e.x = e.tx
			e.y = e.ty
			e.vx = 0
			e.vy = 0
			e.CurrMove = nil
		} else {
			e.CurrMove = e.AnimationQueue[0]
			e.AnimationQueue = e.AnimationQueue[1:]
			if e.CurrMove.SleepPre != 0 {
				//time.Sleep(e.CurrMove.SleepPre)
			}
			e.tx = e.CurrMove.tx
			e.ty = e.CurrMove.ty
			vx := float64(e.tx - e.x)
			vy := float64(e.ty - e.y)
			if vy != 0 || vx != 0 {
				speedVector := csg.NewVector(vx, vy, 0)
				speedVector = speedVector.Normalize().MultiplyScalar(e.CurrMove.Speed)
				e.vx = speedVector.X
				e.vy = speedVector.Y
			} else {
				e.vx = 0
				e.vy = 0
			}

		}

	}
}

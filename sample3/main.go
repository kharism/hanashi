package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/kharism/hanashi/core"

	_ "embed"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	Scene *core.Scene
}

func (g *Game) Update() error {
	HeartImg.Update()
	HeartImg2.Update()
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	// g.Scene.Draw(screen)
	HeartImg.Draw(screen)
	HeartImg2.Draw(screen)
}
func (g *Game) Layout(width, height int) (int, int) {
	return 640, 480
}

// return width and height of the scene
func (g *Game) GetLayout() (width, height int) {
	return 768, 512
}

// return the starting text position where the box containing name of the character appear on the scene
// return negative number if no such box needed
func (g *Game) GetNamePosition() (x, y int) {
	return 0, 512 - 100
}

// get the starting position of the text
func (g *Game) GetTextPosition() (x, y int) {
	return 0, 512 - 70
}

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

var (
	imgPool ImagePool
)

//go:embed heart.png
var Heart []byte
var HeartImg *core.MovableImage
var HeartImg2 *core.MovableImage

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	imgPool = ImagePool{Map: map[string]*ebiten.Image{}}
	fmt.Println(len(Heart))
	heartImg, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(Heart))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	bounds := heartImg.Bounds()
	param := core.MovableImageParams{}
	originX := bounds.Dx() / 2 //float64(640-bounds.Dx()) / 2
	originY := bounds.Dy() / 2 //float64(480-bounds.Dy()) / 2
	param.WithMoveParam(core.MoveParam{Sx: 100, Sy: 100})
	param.WithScale(&core.ScaleParam{Sx: 0.5, Sy: 0.5})
	HeartImg = core.NewMovableImage(heartImg, &param)
	scaleAnim := core.ScaleAnimation{Tsx: 2.0, Tsy: 2.0, SpeedX: 0.01, SpeedY: 0.01, CenterX: float64(originX), CenterY: float64(originY)}
	// scaleAnim := &core.ScaleCenterAnim{Tsx: 5.0, Tsy: 5.0, SpeedX: 0.01, SpeedY: 0.01, OriginX: originX, OriginY: originY}
	HeartImg.AddAnimation(&scaleAnim)

	param2 := core.MovableImageParams{}
	param2.WithMoveParam(core.MoveParam{Sx: 200, Sy: 100})
	param2.WithRotation(&core.RotationParam{Rotation: 0, RotSpeed: 0, CenterX: float64(originX), CenterY: float64(originY)})
	// box := ebiten.NewImage(bounds.Dx(), bounds.Dy())
	// box.Fill(color.White)
	HeartImg2 = core.NewMovableImage(heartImg, &param2)
	rotAnim := core.RotationAnimation{Trot: 70, RotSpeed: 0.1}
	HeartImg2.AddAnimation(&rotAnim)

	game := &Game{}
	// scene := Scene1(game)
	// scene.SetLayouter(game)
	// game.Scene = scene
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}

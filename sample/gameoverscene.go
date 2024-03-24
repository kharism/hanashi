package main

import (
	"image/color"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
)

type GameOverScene struct {
	state          MyState
	StateDecorator StateDecorator
	alpha          float64
	alphaMove      float64
	gameoverImage  *ebiten.Image
	bg             *ebiten.Image
}

func (v *GameOverScene) Load(state MyState, sm *stagehand.SceneManager[MyState]) {

	v.gameoverImage, _ = imgPool.GetImage("bg/game_over.png")
	v.bg = ebiten.NewImage(768, 512)
	v.bg.Fill(color.RGBA{70, 70, 70, 100})
}
func (v *GameOverScene) Unload() MyState {
	if v.StateDecorator != nil {
		return v.StateDecorator(v.state)
	}
	return v.state
}
func (g *GameOverScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// return 768, 512
	return 768, 512
}
func (c *GameOverScene) Update() error {
	// c.alphaMove = 1
	if c.alpha < 1.0 {
		c.alpha += c.alphaMove
		c.alphaMove += 0.0001
		// fmt.Println(c.alpha, c.alphaMove)
	}
	return nil
}

func (c *GameOverScene) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(0, 0)
	screen.DrawImage(c.bg, op)
	op.ColorScale.Scale(1, 1, 1, float32(c.alpha))
	screen.DrawImage(c.gameoverImage, op)

	// if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
	// 	mainGame.(*MainGameState).stateChanger.ChangeState(STATE_MAIN_MENU)
	// }
}

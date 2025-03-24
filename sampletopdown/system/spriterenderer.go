package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/hanashi/sampletopdown/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type SpriteRenderer struct {
	Player *donburi.Entry
	Query  *donburi.Query
}

func (b *SpriteRenderer) RenderBg(ecs *ecs.ECS, screen *ebiten.Image) {
	playerPos := components.GridPos.Get(b.Player)
	renderStartCol := float64(playerPos.Col - 7)
	renderStartRow := float64(playerPos.Row - 5)
	basicTranslateX := -16 * renderStartCol
	basicTranslateY := -16 * renderStartRow
	b.Query.Each(ecs.World, func(e *donburi.Entry) {
		gridPos := components.GridPos.Get(e)
		scrPosX := gridPos.Col * 16
		scrPosY := gridPos.Row * 16
		geom := ebiten.GeoM{}
		geom.Translate(float64(scrPosX)+basicTranslateX, float64(scrPosY)+basicTranslateY)
		s := components.Sprite.Get(b.Player)
		screen.DrawImage(s, &ebiten.DrawImageOptions{GeoM: geom})
	})

}

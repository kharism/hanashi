package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/hanashi/sampletopdown/components"
	"github.com/lafriks/go-tiled/render"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type BgRenderer struct {
	Player      *donburi.Entry
	MapRenderer *render.Renderer
}

// render background so that the player always in center
func (b *BgRenderer) RenderBg(ecs *ecs.ECS, screen *ebiten.Image) {
	playerPos := components.GridPos.Get(b.Player)
	b.MapRenderer.RenderLayer(0)
	ebiBg := ebiten.NewImageFromImage(b.MapRenderer.Result)
	b.MapRenderer.RenderLayer(1)
	ebiBg2 := ebiten.NewImageFromImage(b.MapRenderer.Result)
	renderStartCol := float64(playerPos.Col - 7)
	renderStartRow := float64(playerPos.Row - 5)
	geom := ebiten.GeoM{}
	geom.Translate(-16*renderStartCol, -16*renderStartRow)
	screen.DrawImage(ebiBg, &ebiten.DrawImageOptions{GeoM: geom})
	screen.DrawImage(ebiBg2, &ebiten.DrawImageOptions{GeoM: geom})
}

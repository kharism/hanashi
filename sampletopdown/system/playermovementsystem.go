package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kharism/hanashi/sampletopdown/components"
	"github.com/lafriks/go-tiled"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type PlayerMovementSystem struct {
	Map         *tiled.Map
	Player      *donburi.Entry
	SceneSetter components.HanashiSceneSetter
}

func (s *PlayerMovementSystem) Update(ecs *ecs.ECS) {
	playerGridPos := components.GridPos.Get(s.Player)
	checkCol := playerGridPos.Col
	checkRow := playerGridPos.Row
	groundLayer := s.Map.Layers[0]
	interactibles := s.Map.Layers[1]
	isMoved := false
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		checkCol -= 1
		isMoved = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		checkCol += 1
		isMoved = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		checkRow -= 1
		isMoved = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		checkRow += 1
		isMoved = true
	}
	rawIdx := checkRow*s.Map.Width + checkCol
	// fmt.Println(groundLayer.Tiles[rawIdx].ID)
	if isMoved && groundLayer.Tiles[rawIdx].ID == 6 && interactibles.Tiles[rawIdx].ID == 0 {
		playerGridPos.Col = checkCol
		playerGridPos.Row = checkRow
	} else if isMoved && interactibles.Tiles[rawIdx].ID != 0 {
		if _, ok := components.InteractibleMap[rawIdx]; ok {
			components.InteractibleMap[rawIdx](ecs, s.SceneSetter)
		}

	}

}

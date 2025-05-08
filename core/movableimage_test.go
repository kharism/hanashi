package core

import (
	"testing"
)

// simple horizontal and vertical movement
func TestMovement1(t *testing.T) {
	j := NewMovableImage(nil, NewMovableImageParams().WithMoveParam(MoveParam{Sx: 10, Sy: 10}))
	j.AddAnimation(NewMoveAnimationFromParam(MoveParam{Tx: 30, Ty: 10, Speed: 5}))
	j.Update()
	xPos, yPos := j.GetPos()
	// t.Log(xPos, yPos)
	if xPos != 15 && yPos != 10 {
		// t.Fail()
		t.FailNow()
	}
	j.Update()
	xPos, yPos = j.GetPos()
	// t.Log(j.GetPos())
	if xPos != 20 && yPos != 10 {
		// t.Fail()
		t.FailNow()
	}
	j.ReplaceCurrentAnim(NewMoveAnimationFromParam(MoveParam{Tx: 20, Ty: 30, Speed: 5}))
	j.Update()
	xPos, yPos = j.GetPos()
	// t.Log(j.GetPos())
	if xPos != 20 && yPos != 15 {
		// t.Fail()
		t.FailNow()
	}
}

// test move diagonal 45 degree
func TestMovement2(t *testing.T) {
	j := NewMovableImage(nil, NewMovableImageParams().WithMoveParam(MoveParam{Sx: 10, Sy: 10}))
	j.AddAnimation(NewMoveAnimationFromParam(MoveParam{Tx: 30, Ty: 30, Speed: 5}))
	j.Update()
	xPos, yPos := j.GetPos()
	// t.Log(xPos, yPos)
	if xPos != 13.535533905932738 && yPos != 13.535533905932738 {
		// t.Fail()
		t.FailNow()
	}
	j.Update()
	xPos, yPos = j.GetPos()
	// t.Log(j.GetPos())
	if xPos != 17.071067811865476 && yPos != 17.071067811865476 {
		// t.Fail()
		t.FailNow()
	}
}

func TestScale(t *testing.T) {
	j := NewMovableImage(nil, NewMovableImageParams().WithMoveParam(MoveParam{Sx: 10, Sy: 10}))
	j.AddAnimation(&ScaleAnimation{Tsx: 1.2, Tsy: 1.2, SpeedX: 0.01, SpeedY: 0.01})
	j.Update()
	t.Log(j.ScaleParam.Sx, j.ScaleParam.Sy)
	if j.ScaleParam.Sx != 1.01 || j.ScaleParam.Sy != 1.01 {
		t.FailNow()
	}
}

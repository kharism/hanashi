package main

import (
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

// play attack animation
type AttackAnim struct {
	counter  int
	cs       *CombatScene
	doneFunc func(cs *CombatScene)
}

func (a *AttackAnim) GetDoneFunc() func(cs *CombatScene) {
	return a.doneFunc
}
func (a *AttackAnim) SetDoneFunc(g func(cs *CombatScene)) {
	a.doneFunc = g
}
func (a *AttackAnim) Update() {
	a.cs.attackAnim.Update()
}
func (a *AttackAnim) Draw(screen *ebiten.Image) {
	a.cs.opp.Draw(screen)
	a.cs.attackAnim.Draw(screen)
}

// move attack animation
type MoveAttackAnim struct {
	counter  int
	cs       *CombatScene
	doneFunc func(cs *CombatScene)

	PosX float64
	PosY float64
}

func (a *MoveAttackAnim) GetDoneFunc() func(cs *CombatScene) {
	return a.doneFunc
}
func (a *MoveAttackAnim) SetDoneFunc(g func(cs *CombatScene)) {
	a.doneFunc = g
}
func (a *MoveAttackAnim) Update() {
	a.cs.attackAnim.SetPos(a.PosX, a.PosY)

	// a.cs.attackAnim.Update()
}
func (a *MoveAttackAnim) Draw(screen *ebiten.Image) {
	a.cs.opp.Draw(screen)
	// a.cs.attackAnim.Draw(screen)
	if a.doneFunc != nil {
		a.doneFunc(a.cs)
	}
}

// blink the opp sprite
type BlinkAnim struct {
	counter int

	// how many blinking
	blinkCount int
	cs         *CombatScene
	doneFunc   func(cs *CombatScene)
}

func (a *BlinkAnim) GetDoneFunc() func(cs *CombatScene) {
	return a.doneFunc
}
func (a *BlinkAnim) SetDoneFunc(g func(cs *CombatScene)) {
	a.doneFunc = g
}

func (a *BlinkAnim) Update() {
	// a.cs.attackAnim.Update()
	a.counter = (a.counter + 1) //% (5 * a.blinkCount)
}
func (a *BlinkAnim) Draw(screen *ebiten.Image) {
	// a.cs.attackAnim.Draw(screen)
	if (a.counter/2)%2 != 0 {
		a.cs.opp.Draw(screen)
	}

	if a.counter == 2*a.blinkCount {
		a.doneFunc(a.cs)
	}
}

// complex animation where multiple animation is queued as single animation.
// useful when we want to ensure the order of animation if multiple routine
// is accessing the same channel
type ComplexAnim struct {
	// index of the current animation, always init this value to -1
	idx        int
	cs         *CombatScene
	animations []CombatAnimation
	doneFunc   func(cs *CombatScene)
}

func (a *ComplexAnim) GetDoneFunc() func(cs *CombatScene) {
	return a.doneFunc
}
func (a *ComplexAnim) SetDoneFunc(g func(cs *CombatScene)) {
	a.doneFunc = g
}
func (c *ComplexAnim) Draw(screen *ebiten.Image) {
	c.animations[c.idx].Draw(screen)
}
func (c *ComplexAnim) Update() {
	if c.idx == len(c.animations) {
		c.doneFunc(c.cs)
		return
	}
	// fmt.Println(c.idx)
	// prev := c.idx
	c.animations[c.idx].Update()
	// fmt.Println("=====")
	// fmt.Println(c.idx)
	if c.idx == len(c.animations) {
		c.doneFunc(c.cs)
		return
	}
	if c.animations[c.idx].GetDoneFunc() == nil {
		c.animations[c.idx].SetDoneFunc(func(cs *CombatScene) {
			c.idx += 1
		})
	}

}

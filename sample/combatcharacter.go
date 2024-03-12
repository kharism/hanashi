package main

import "github/kharism/hanashi/core"

// struct used to define character in combatscene
//
type CombatCharacter struct {
	Name       string
	HP         int
	HPMax      int
	AttackAnim *core.AnimatedImage
}

func NewCombatCharacter(name string, Hp, HpMax int) *CombatCharacter {
	return &CombatCharacter{Name: name, HP: Hp, HPMax: HpMax}
}

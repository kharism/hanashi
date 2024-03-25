package main

import (
	"math/rand"

	"github.com/kharism/hanashi/core"
)

// struct used to define character in combatscene
type CombatCharacter struct {
	Name       string
	HP         int
	HPMax      int
	AtkDamage  GetAtkDamage
	AttackAnim *core.AnimatedImage
}

func Dice(num, face int) int {
	total := 0
	for i := 0; i < num; i++ {
		total += rand.Int()%face + 1
	}
	return total
}

type GetAtkDamage func() int

func (h *CombatCharacter) SetAtkDamage(j GetAtkDamage) *CombatCharacter {
	h.AtkDamage = j
	return h
}
func NewCombatCharacter(name string, Hp, HpMax int) *CombatCharacter {
	return &CombatCharacter{Name: name, HP: Hp, HPMax: HpMax}
}

package main

import (
	"fmt"
	"math/rand"
)

// how opp bot will do during combat
type CombatRoutine interface {
	Execute(cs *CombatScene, ls *CombatLogSubstate)
}

type AttackRandomly struct {
	atkDamage GetAtkDamage
}

func NewAttackRandomly(a GetAtkDamage) CombatRoutine {
	return &AttackRandomly{atkDamage: a}
}
func (a *AttackRandomly) Execute(cs *CombatScene, ls *CombatLogSubstate) {
	damage := a.atkDamage()
	target := rand.Int() % len(cs.Characters)
	character := cs.Characters[target]
	for true {
		if character.HP == 0 {
			target = rand.Int() % len(cs.Characters)
			character = cs.Characters[target]
		} else {
			break
		}
	}

	// damage := character.AtkDamage()
	cplx_anim := &ComplexAnim{cs: cs, doneFunc: func(cs *CombatScene) {
		// cs.OppTakeDamage(damage)
		if damage > character.HP {
			character.HP = 0
		} else {
			character.HP -= damage
		}
		DoneAnim(cs)
	}, animations: []CombatAnimation{
		&MoveAttackAnim{cs: cs, PosX: float64(PLAYER_POS_X_START + target*170), PosY: float64(PLAYER_POS_Y_START - 30)},
		&AttackAnim{cs: cs},
		// &BlinkAnim{cs: cs, blinkCount: 5},
	},
	}
	cs.attackAnim.Done = func() {
		if cplx_anim.idx < len(cplx_anim.animations) {
			cplx_anim.idx += 1
		}

	}
	go func() {
		cs.animationQueue <- cplx_anim
	}()
	ls.BattleLog = fmt.Sprintf("opponent deals %d of damage to %s", damage, character.Name)
}

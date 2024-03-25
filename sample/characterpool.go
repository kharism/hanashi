package main

import "github.com/kharism/hanashi/core"

var Characters map[string]*core.Character

type CharPool struct {
}

func (c *CharPool) GetCharacter(name string) *core.Character {
	return Characters[name]
}

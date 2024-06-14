package core

import "github.com/hajimehoshi/ebiten/v2"

type CharacterViewEvent struct {
	Name       string
	MoveParam  *MoveParam
	ScaleParam *ScaleParam
}

// even to add character to the scene. name is character name that's already
// registered on characters field on scene
// moveParam determine how the character should get into scene
// scaleParam determine how
func NewCharacterAddEvent(name string, moveParam *MoveParam, scaleParam *ScaleParam) Event {
	return &CharacterViewEvent{name, moveParam, scaleParam}
}
func (s *CharacterViewEvent) Execute(Scene *Scene) {
	Scene.AddViewableCharacters(s.Name, s.MoveParam, s.ScaleParam)
}

type CharacterRemoveEvent struct {
	Name string
}

func NewCharacterRemoveEvent(name string) Event {
	return &CharacterRemoveEvent{name}
}

func (s *CharacterRemoveEvent) Execute(Scene *Scene) {
	Scene.RemoveVieableCharacter(s.Name)
}

type CharacterMoveEvent struct {
	Name string
	MoveParam
}

func NewCharacterMoveEvent(name string, moveParam MoveParam) Event {
	return &CharacterMoveEvent{Name: name, MoveParam: moveParam}
}

func (s *CharacterMoveEvent) Execute(scene *Scene) {
	for _, c := range scene.ViewableCharacters {
		if c.Name == s.Name {
			ma := MoveAnimation{}
			ma.tx = s.MoveParam.Tx
			ma.ty = s.MoveParam.Ty
			ma.Speed = s.MoveParam.Speed
			c.Img.AddAnimation(&ma)
		}
	}
}

type CharacterComplexMoveEvent struct {
	Name           string
	AnimationQueue []Animation
}

func (s *CharacterComplexMoveEvent) Execute(scene *Scene) {
	for _, c := range scene.ViewableCharacters {
		if c.Name == s.Name {
			for _, ma := range s.AnimationQueue {
				c.Img.AddAnimation(ma)
			}
		}
	}
}

type CharacterAddShaderEvent struct {
	Name   string
	Shader *ebiten.Shader
}

// add shader to a character, set shaderparam to nil if you want to remove any shader
func NewCharacterAddShaderEvent(name string, shaderParam *ShaderParam) Event {
	if shaderParam == nil {
		return &CharacterAddShaderEvent{Name: name, Shader: nil}
	}
	if shaderParam.Shader == nil {
		sh, _ := GetShaderPool().GetShader(shaderParam.ShaderName)
		return &CharacterAddShaderEvent{Name: name, Shader: sh}
	}
	return &CharacterAddShaderEvent{Name: name, Shader: shaderParam.Shader}
}
func (s *CharacterAddShaderEvent) Execute(Scene *Scene) {
	// Scene.AddViewableCharacters(s.Name, s.MoveParam, s.ScaleParam)
	for _, c := range Scene.ViewableCharacters {
		if c.Name == s.Name {
			c.Img.Shader = s.Shader
			break
		}
	}
}

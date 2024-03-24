package core

import "testing"

func TestSceneMM(t *testing.T) {
	s := NewScene()
	vv := s.GetSceneData("DDD")
	if vv != nil {
		t.FailNow()
	}
}

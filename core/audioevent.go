package core

type PlayBgmEvent struct {
	Audio *[]byte
	Type  MusicType
}

func (f *PlayBgmEvent) Execute(scene *Scene) {
	scene.AudioInterface.PlayBgm(*f.Audio, f.Type)
}

type StopBgmEvent struct{}

func (f *StopBgmEvent) Execute(scene *Scene) {
	scene.AudioInterface.StopBgm()
}

type PlaySfxEvent struct {
	Audio *[]byte
	Type  MusicType
	Name  string
}

func (f *PlaySfxEvent) Execute(scene *Scene) {
	scene.AudioInterface.PlaySfx(*f.Audio, f.Type, f.Name)
}

type StopSfxEvent struct {
	Name string
}

func (f *StopSfxEvent) Execute(scene *Scene) {
	scene.AudioInterface.StopSfx(f.Name)
}

type SetVolumeEvent struct {
	// 0 for BGM, 1 for SFX
	Type int
	// 0-128 only
	Vol int
}

func (f *SetVolumeEvent) Execute(scene *Scene) {
	if f.Type == 0 {
		scene.AudioInterface.SetBGMVolume(f.Vol)
	} else if f.Type == 1 {
		scene.AudioInterface.SetSFXVolume(f.Vol)
	}
}

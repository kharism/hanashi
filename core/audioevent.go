package core

type PlayBgmEvent struct {
	Audio *[]byte
	Type  MusicType
}

func (f *PlayBgmEvent) Execute(scene *Scene) {
	scene.AudioInterface.PlayBgm(*f.Audio, f.Type)
}

type PlaySfxEvent struct {
	Audio *[]byte
	Type  MusicType
}

func (f *PlaySfxEvent) Execute(scene *Scene) {
	scene.AudioInterface.PlaySfx(*f.Audio, f.Type)
}

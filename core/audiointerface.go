package core

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

// a basic audio interface to handle audio based events
// creating an interface so client program can overwrite it using their own audio system
type AudioInterface interface {
	// set current bgm, audio is raw bytes
	PlayBgm(audio []byte, format MusicType)
	// stop current bgm
	StopBgm()
	// this gets called in update function
	Update() error
	// use this to play short sfx
	PlaySfx(audio []byte, format MusicType)
	// StopSfx()
}
type MusicType int

const (
	TypeOgg MusicType = iota
	TypeMP3
)

var audioContext *audio.Context

const (
	sampleRate = 48000
)

func init() {
	audioContext = audio.NewContext(sampleRate)
}

type DefaultAudioInterface struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
	current      time.Duration
	// total        time.Duration
	seBytes []byte
	seCh    chan []byte

	volume128 int

	musicType MusicType
	playSfx   bool
}

func (p *DefaultAudioInterface) Close() error {
	return p.audioPlayer.Close()
}
func (p *DefaultAudioInterface) Update() error {
	select {
	case p.seBytes = <-p.seCh:
		// fmt.Println("SFX detected")
		// close(p.seCh)
		// p.playSfx = true

		// p.seCh = nil
	default:
	}
	if !p.audioPlayer.IsPlaying() {
		p.audioPlayer.Rewind()
		p.audioPlayer.Play()
	}
	p.PlaySEIfNeeded()
	return nil
}
func (p *DefaultAudioInterface) ShouldPlaySE() bool {
	if p.seBytes == nil {
		// Bytes for the SE is not loaded yet.
		return false
	}
	// fmt.Println(p.seCh)
	return p.seCh != nil
}

func (p *DefaultAudioInterface) PlaySEIfNeeded() {
	if !p.ShouldPlaySE() {
		return
	}
	sePlayer := p.audioContext.NewPlayerFromBytes(p.seBytes)
	sePlayer.Play()
	p.seBytes = nil
}
func (p *DefaultAudioInterface) PlayBgm(audio []byte, format MusicType) {
	// var b []byte
	var err error
	var s io.Reader
	if format == TypeMP3 {
		s, err = mp3.DecodeWithoutResampling(bytes.NewReader(audio))
		// p.seCh <- param
		if err != nil {
			log.Fatal(err)
			return
		}
		// b, err = io.ReadAll(s)
	} else if format == TypeOgg {
		s, err = mp3.DecodeWithoutResampling(bytes.NewReader(audio))
		// p.seCh <- param
		if err != nil {
			log.Fatal(err)
			return
		}
		// b, err = io.ReadAll(s)
	}
	pl, err := audioContext.NewPlayer(s)
	if err != nil {
		return
	}
	p.audioPlayer = pl
	p.audioPlayer.Play()
}
func (p *DefaultAudioInterface) StopBgm() {
	p.audioPlayer.Pause()
}
func NewDefaultAudioInterfacer() (AudioInterface, error) {
	type audioStream interface {
		io.ReadSeeker
		Length() int64
	}
	const bytesPerSample = 4 // TODO: This should be defined in audio package

	var s audioStream
	// audio, err := os.Open(audioPath)
	// if err != nil {
	// 	return nil, err
	// }
	// defer audio.Close()

	p, err := audioContext.NewPlayer(s)
	if err != nil {
		return nil, err
	}

	player := &DefaultAudioInterface{
		audioContext: audioContext,
		audioPlayer:  p,
		// total:        time.Second * time.Duration(s.Length()) / bytesPerSample / sampleRate,
		volume128: 2,
		seCh:      make(chan []byte, 100),
		seBytes:   []byte{},
		// musicType:    musicType,
	}
	// if player.total == 0 {
	// 	player.total = 1
	// }

	// player.audioPlayer.Play()

	return player, nil
}
func (p *DefaultAudioInterface) PlaySfx(audio []byte, format MusicType) {
	go func() {
		var b []byte
		var err error
		var s io.Reader
		if format == TypeMP3 {
			s, err = mp3.DecodeWithoutResampling(bytes.NewReader(audio))
			// p.seCh <- param
			if err != nil {
				log.Fatal(err)
				return
			}
			b, err = io.ReadAll(s)
		} else if format == TypeOgg {
			s, err = vorbis.DecodeWithoutResampling(bytes.NewReader(audio))
			// p.seCh <- param
			if err != nil {
				log.Fatal(err)
				return
			}
			b, err = io.ReadAll(s)
		}

		if err != nil {
			log.Fatal(err)
			return
		}
		p.seCh <- b
	}()

}

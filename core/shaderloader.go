package core

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed shaders/invert.kage
var InvertShader []byte

const INVERT_SHADER = "INVERT_SHADER"

//go:embed shaders/sepia.kage
var SepiaShader []byte

const SEPIA_SHADER = "SEPIA_SHADER"

//go:embed shaders/darker.kage
var DarkerShader []byte

const DARKER_SHADER = "DARKER_SHADER"

//go:embed shaders/darker.kage
var GrayscaleShader []byte

const GRAYSCALE_SHADER = "DARKER_SHADER"

//go:embed shaders/paletteswapbrg.kage
var PaletteSwapBrgShader []byte

const PALETTESWAPBRG_SHADER = "BRG_SHADER"

// pool that store compiled shader. It is a singleton
// register new shader here so you can use it alongside ShaderParam
type ShaderPool struct {
	mutex     *sync.Mutex
	shaderMap map[string]*ebiten.Shader
}

var shaderPool *ShaderPool

// the default way to get shaderPool object
func GetShaderPool() *ShaderPool {
	if shaderPool == nil {
		shaderPool = &ShaderPool{mutex: &sync.Mutex{}, shaderMap: map[string]*ebiten.Shader{}}
	}
	return shaderPool
}
func (s *ShaderPool) RegisterShader(name string, shader *ebiten.Shader) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.shaderMap[name] = shader
}

var ShaderNotFoundError = errors.New("Shader Not Found")

func (s *ShaderPool) GetShader(name string) (*ebiten.Shader, error) {
	if _, ok := s.shaderMap[name]; ok {
		return s.shaderMap[name], nil
	}
	return nil, ShaderNotFoundError
}
func init() {
	pool := GetShaderPool()
	sepia, err := ebiten.NewShader(SepiaShader)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}
	invert, err := ebiten.NewShader(InvertShader)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}
	darker, err := ebiten.NewShader(DarkerShader)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}
	grayscale, err := ebiten.NewShader(GrayscaleShader)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}

	psbrg, err := ebiten.NewShader(PaletteSwapBrgShader)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}

	pool.RegisterShader(SEPIA_SHADER, sepia)
	pool.RegisterShader(INVERT_SHADER, invert)
	pool.RegisterShader(DARKER_SHADER, darker)
	pool.RegisterShader(GRAYSCALE_SHADER, grayscale)
	pool.RegisterShader(PALETTESWAPBRG_SHADER, psbrg)
}

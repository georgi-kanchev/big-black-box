package internal

import (
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var shader *ebiten.Shader

func ShaderCache() {
	var sh, err = ebiten.NewShader(shaderCode)
	if err != nil {
		log.Println(err)
	}
	shader = sh
}

// private ========================================================

//go:embed shader.kage
var shaderCode []byte

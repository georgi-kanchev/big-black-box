// stores the loaded resources from the assets package for everybody to access

package internal

import (
	_ "embed"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Font byte
type Image uint16

var Fonts []*ebiten.Image = make([]*ebiten.Image, 0, 64)
var Images []*ebiten.Image = make([]*ebiten.Image, 0, 64)

var White1x1 = ebiten.NewImage(1, 1)

// private ========================================================

var shader *ebiten.Shader

//go:embed shader.kage
var shaderCode []byte

func initAssets() {
	White1x1.Set(0, 0, color.White)

	var sh, err = ebiten.NewShader(shaderCode)
	if err != nil {
		log.Println(err)
	}
	shader = sh
}

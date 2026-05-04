// stores the loaded resources from the assets package for everybody to access

package internal

import (
	_ "embed"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type FontId byte
type ImageId uint16

var Fonts []*ebiten.Image = make([]*ebiten.Image, 0, 64)
var Images []*ebiten.Image = make([]*ebiten.Image, 0, 64)

var White1x1 = ebiten.NewImage(1, 1)

func AddImage(img *ebiten.Image) ImageId {
	var newId int

	if len(freeImageIds) > 0 { // any IDs to reuse?
		var lastIndex = len(freeImageIds) - 1
		newId = freeImageIds[lastIndex]
		freeImageIds = freeImageIds[:lastIndex] // reslice to remove it

		Images[newId] = img
	} else { // no free IDs, grow the main slice
		newId = len(Images)
		Images = append(Images, img)
	}

	return ImageId(newId + 1) // 0 is default image (white1x1)
}
func RemoveImage(id ImageId) {
	if int(id) >= len(Images) || Images[id] == nil {
		return
	}

	Images[id].Deallocate()
	Images[id] = nil
	freeImageIds = append(freeImageIds, int(id))
}

func AddFont(img *ebiten.Image) FontId {
	var newId int

	if len(freeFontIds) > 0 { // any IDs to reuse?
		var lastIndex = len(freeFontIds) - 1
		newId = freeFontIds[lastIndex]
		freeFontIds = freeFontIds[:lastIndex] // reslice to remove it

		Fonts[newId] = img
	} else { // no free IDs, grow the main slice
		newId = len(Fonts)
		Fonts = append(Fonts, img)
	}

	return FontId(newId + 1) // 0 is default font (ebiten)
}
func RemoveFont(id FontId) {
	if int(id) >= len(Fonts) || Fonts[id] == nil {
		return
	}

	Fonts[id].Deallocate()
	Fonts[id] = nil
	freeFontIds = append(freeFontIds, int(id))
}

// private ========================================================

var freeImageIds []int = make([]int, 0, 64)
var freeFontIds []int = make([]int, 0, 64)

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

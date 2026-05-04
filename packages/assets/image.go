package assets

import (
	"big-black-box/packages/internal"
	"big-black-box/packages/utility/file"
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Image internal.ImageId

const DefaultImage Image = 0

func LoadImage(imagePath string) Image {
	var img, _, err = image.Decode(bytes.NewReader(file.LoadBytes(imagePath)))
	if err != nil {
		log.Println(err)
		return 0
	}

	return Image(internal.AddImage(ebiten.NewImageFromImage(img)))
}

func (i Image) Size() (width, height int) {
	var img = internal.Images[i]
	var bounds = img.Bounds()
	return bounds.Dx(), bounds.Dy()
}
func (i Image) Unload() {

}

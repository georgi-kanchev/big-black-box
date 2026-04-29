package assets

import (
	"big-black-box/packages/internal"
	"big-black-box/packages/utility/file"
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Image internal.Image

func LoadImage(imagePath string) Image {
	var img, _, err = image.Decode(bytes.NewReader(file.LoadBytes(imagePath)))
	if err != nil {
		log.Println(err)
		return 0
	}

	var asset = ebiten.NewImageFromImage(img)
	internal.Images = append(internal.Images, asset)
	return Image(len(internal.Images))
}

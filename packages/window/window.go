// holds some data that internal/window reads from to update ebiten

package window

import "github.com/hajimehoshi/ebiten/v2"

type Mode byte

var Title = "game"
var VSync = true
var CurrentMode Mode = 1
var Monitor byte = 0
var PixelScale float32 = 1

var Width, Height float32 // Read-only

func IsFocused() bool {
	return ebiten.IsFocused()
}

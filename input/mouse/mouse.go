package mouse

import (
	"unsafe"

	"big-black-box/internal"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Button ebiten.MouseButton
type Cursor ebiten.CursorShapeType

//=================================================================

func SetCursorVisibility(visible bool) {
	if visible {
		ebiten.SetCursorMode(ebiten.CursorModeVisible)
	} else {
		ebiten.SetCursorMode(ebiten.CursorModeHidden)
	}
}
func SetCursor(cursor Cursor) {
	ebiten.SetCursorShape(ebiten.CursorShapeType(cursor))
}

//=================================================================

func CursorDelta() (x, y float32) {
	return internal.MouseDeltaX, internal.MouseDeltaY
}
func GetCursor() Cursor {
	return Cursor(ebiten.CursorShape())
}

func Scroll() float32 {
	return internal.Scroll
}
func ScrollSmooth() float32 {
	return internal.SmoothScroll
}

func Pressed() []Button {
	internal.ButtonsPressed = internal.ButtonsPressed[:0]
	for b := ebiten.MouseButton0; b <= ebiten.MouseButtonMax; b++ {
		if ebiten.IsMouseButtonPressed(b) {
			internal.ButtonsPressed = append(internal.ButtonsPressed, b)
		}
	}
	return *(*[]Button)(unsafe.Pointer(&internal.ButtonsPressed))
}

func IsPressed(button Button) bool {
	return ebiten.IsMouseButtonPressed(ebiten.MouseButton(button))
}
func IsJustPressed(button Button) bool {
	return inpututil.IsMouseButtonJustPressed(ebiten.MouseButton(button))
}
func IsJustReleased(button Button) bool {
	return inpututil.IsMouseButtonJustReleased(ebiten.MouseButton(button))
}

func IsAnyPressed() bool {
	return len(Pressed()) > 0
}
func IsAnyJustPressed() bool {
	for b := ebiten.MouseButton0; b <= ebiten.MouseButtonMax; b++ {
		if inpututil.IsMouseButtonJustPressed(b) {
			return true
		}
	}
	return false
}
func IsAnyJustReleased() bool {
	for b := ebiten.MouseButton0; b <= ebiten.MouseButtonMax; b++ {
		if inpututil.IsMouseButtonJustReleased(b) {
			return true
		}
	}
	return false
}

//=================================================================

func (b Button) IsPressed() bool      { return IsPressed(b) }
func (b Button) IsJustPressed() bool  { return IsJustPressed(b) }
func (b Button) IsJustReleased() bool { return IsJustReleased(b) }

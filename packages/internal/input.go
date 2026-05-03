// caches the input values this tick that other packages can use
// and the game can access through the input package

package internal

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var MouseX, MouseY, MouseDeltaX, MouseDeltaY, SmoothScroll, Scroll float32
var ButtonsPressed = make([]ebiten.MouseButton, 0, 5)

var KeysPressed = make([]ebiten.Key, 0, 5)
var InputBuffer = make([]rune, 0, 16)
var InputString string
var AnyKeyJustPressed, AnyKeyJustReleased bool

// private ========================================================

var justPressed = make([]ebiten.Key, 0, 4)
var justReleased = make([]ebiten.Key, 0, 4)

func cacheInput() {
	if !ebiten.IsFocused() {
		KeysPressed = KeysPressed[:0]
		ButtonsPressed = ButtonsPressed[:0]
		InputString = ""
		return
	}

	var mx, my = ebiten.CursorPosition()
	var currX, currY = float32(mx), float32(my)

	MouseDeltaX = currX - MouseX
	MouseDeltaY = currY - MouseY
	MouseX, MouseY = currX, currY

	var _, scroll = ebiten.Wheel()
	Scroll = float32(scroll)

	var tps = ebiten.TPS()
	var dt = float32(1.0 / tps)
	const scrollAccel, scrollDecay = 600.0, -8.0

	SmoothScroll += Scroll * scrollAccel * dt
	SmoothScroll *= float32(math.Exp(scrollDecay * 1 / float64(tps)))

	if SmoothScroll != 0 && (SmoothScroll < 0.0001 && SmoothScroll > -0.0001) {
		SmoothScroll = 0
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) || inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButton1) || inpututil.IsMouseButtonJustReleased(ebiten.MouseButton1) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButton2) || inpututil.IsMouseButtonJustReleased(ebiten.MouseButton2) {
		SmoothScroll = 0
	}

	//=================================================================

	InputBuffer = ebiten.AppendInputChars(InputBuffer[:0])
	InputString = string(InputBuffer)

	AnyKeyJustPressed = false
	justReleased = justReleased[:0]
	inpututil.AppendJustPressedKeys(justReleased)
	if len(justReleased) > 0 {
		AnyKeyJustPressed = true
	}

	AnyKeyJustReleased = false
	justReleased = justReleased[:0]
	inpututil.AppendJustReleasedKeys(justReleased)
	if len(justReleased) > 0 {
		AnyKeyJustReleased = true
	}
}

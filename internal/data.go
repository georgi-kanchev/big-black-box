package internal

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Data struct {
	Engine Engine
	Window Window
}

var State Data
var GameLoop func()

var MouseX, MouseY, MouseDeltaX, MouseDeltaY, SmoothScroll, Scroll float32
var ButtonsPressed = make([]ebiten.MouseButton, 0, 5)

var KeysPressed = make([]ebiten.Key, 0, 5)
var InputBuffer = make([]rune, 0, 16)
var InputString string
var AnyKeyJustPressed, AnyKeyJustReleased bool

//=================================================================

func (d *Data) Update() error {
	if d.Engine.Exiting {
		return ebiten.Termination
	}
	cacheInput()
	GameLoop()
	return nil
}
func (d *Data) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %f", ebiten.ActualFPS()), 0, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %f", ebiten.ActualTPS()), 0, 16)
}
func (d *Data) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth / 4, outsideHeight / 4
}

// private =================================================================

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

	var tps = State.Engine.TargetTickRate
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
	var tempKeys = inpututil.AppendJustPressedKeys(nil)
	if len(tempKeys) > 0 {
		AnyKeyJustPressed = true
	}

	AnyKeyJustReleased = false
	var tempReleased = inpututil.AppendJustReleasedKeys(nil)
	if len(tempReleased) > 0 {
		AnyKeyJustReleased = true
	}
}

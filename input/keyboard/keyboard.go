package keyboard

import (
	"big-black-box/internal"
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Key ebiten.Key

// Input returns the text characters typed this frame.
func Input() string {
	return internal.InputString
}

func Pressed() []Key {
	internal.KeysPressed = internal.KeysPressed[:0]
	internal.KeysPressed = inpututil.AppendPressedKeys(internal.KeysPressed)
	return *(*[]Key)(unsafe.Pointer(&internal.KeysPressed))
}

func IsPressed(key Key) bool {
	return ebiten.IsKeyPressed(ebiten.Key(key))
}
func IsHeld(key Key) bool {
	var d = inpututil.KeyPressDuration(ebiten.Key(key))
	if d == 1 { // Just pressed
		return true
	}
	if d >= 24 && (d-24)%6 == 0 { // Standard repeating interval
		return true
	}
	return false
}

func IsJustPressed(key Key) bool {
	return inpututil.IsKeyJustPressed(ebiten.Key(key))
}
func IsJustReleased(key Key) bool {
	return inpututil.IsKeyJustReleased(ebiten.Key(key))
}

func IsAnyPressed() bool {
	internal.KeysPressed = internal.KeysPressed[:0]
	return len(inpututil.AppendPressedKeys(internal.KeysPressed)) > 0
}
func IsAnyJustPressed() bool {
	return internal.AnyKeyJustPressed
}
func IsAnyJustReleased() bool {
	return internal.AnyKeyJustReleased
}

func IsComboJustPressed(keys ...Key) bool {
	if !IsJustPressed(keys[len(keys)-1]) {
		return false
	}
	return combo(keys)
}
func IsComboHeld(keys ...Key) bool {
	if !IsHeld(keys[len(keys)-1]) {
		return false
	}
	return combo(keys)
}

//=================================================================

func (k Key) IsPressed() bool      { return IsPressed(k) }
func (k Key) IsHeld() bool         { return IsHeld(k) }
func (k Key) IsJustPressed() bool  { return IsJustPressed(k) }
func (k Key) IsJustReleased() bool { return IsJustReleased(k) }

// private =================================================================

func combo(keys []Key) bool {
	internal.KeysPressed = internal.KeysPressed[:0]
	internal.KeysPressed = inpututil.AppendPressedKeys(internal.KeysPressed)

	if len(internal.KeysPressed) != len(keys) {
		return false
	}

	for _, k := range keys {
		found := false
		for _, p := range internal.KeysPressed {
			if k == Key(p) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	for i := 0; i < len(keys)-1; i++ {
		if inpututil.KeyPressDuration(ebiten.Key(keys[i])) <= inpututil.KeyPressDuration(ebiten.Key(keys[i+1])) {
			return false
		}
	}

	return true
}

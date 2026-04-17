package keyboard

import (
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Key ebiten.Key

// Input returns the text characters typed this frame.
func Input() string {
	return string(ebiten.AppendInputChars(nil))
}

func Pressed() []Key {
	pressed = pressed[:0] // Reset length, keep capacity
	pressed = inpututil.AppendPressedKeys(pressed)
	return *(*[]Key)(unsafe.Pointer(&pressed))
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
	pressed = pressed[:0]
	return len(inpututil.AppendPressedKeys(pressed)) > 0
}
func IsAnyJustPressed() bool {
	justPressed = justPressed[:0]
	return len(inpututil.AppendJustPressedKeys(justPressed)) > 0
}
func IsAnyJustReleased() bool {
	justReleased = justReleased[:0]
	return len(inpututil.AppendJustReleasedKeys(justReleased)) > 0
}

// IsComboJustPressed checks if the keys are all held,
// they were pressed in the specific order provided,
// and the final key was just pressed this frame.
func IsComboJustPressed(keys ...Key) bool {
	if !IsJustPressed(keys[len(keys)-1]) {
		return false
	}
	return combo(keys)
}

// IsComboHeld checks if the keys are held in the specific order
// and follows the repeat-logic for the final key.
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
// private

var pressed, justPressed, justReleased = make([]ebiten.Key, 0, 8), make([]ebiten.Key, 0, 8), make([]ebiten.Key, 0, 8)

//=================================================================

func combo(keys []Key) bool {
	pressed = pressed[:0]
	pressed = inpututil.AppendPressedKeys(pressed)

	if len(pressed) != len(keys) {
		return false
	}

	for _, k := range keys {
		found := false
		for _, p := range pressed {
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

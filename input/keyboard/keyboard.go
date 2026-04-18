package keyboard

import (
	"big-black-box/internal"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Key ebiten.Key

// Input returns the text characters typed this frame.
func Input() string {
	return internal.InputString
}

func Pressed() [5]Key {
	internal.KeysPressed = internal.KeysPressed[:0]
	internal.KeysPressed = inpututil.AppendPressedKeys(internal.KeysPressed)

	var result [5]Key
	for i := 0; i < len(internal.KeysPressed) && i < 5; i++ {
		result[i] = Key(internal.KeysPressed[i])
	}
	return result
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

func IsComboPressed(key1, key2, key3 Key) bool {
	if key1 != keyNone && !IsPressed(key1) {
		return false
	}
	if key2 != keyNone && !IsPressed(key2) {
		return false
	}
	if key3 != keyNone && !IsPressed(key3) {
		return false
	}

	return checkOrder(key1, key2, key3)
}
func IsComboJustPressed(key1, key2, key3 Key) bool {
	var trigger Key
	if key3 != keyNone {
		trigger = key3
	} else if key2 != keyNone {
		trigger = key2
	} else {
		trigger = key1
	}

	if !IsJustPressed(trigger) {
		return false
	}

	return IsComboPressed(key1, key2, key3)
}
func IsComboHeld(key1, key2, key3 Key) bool {
	var trigger Key
	if key3 != keyNone {
		trigger = key3
	} else if key2 != keyNone {
		trigger = key2
	} else {
		trigger = key1
	}

	if !IsHeld(trigger) {
		return false
	}

	return IsComboPressed(key1, key2, key3)
}

//=================================================================

func (k Key) IsPressed() bool      { return IsPressed(k) }
func (k Key) IsHeld() bool         { return IsHeld(k) }
func (k Key) IsJustPressed() bool  { return IsJustPressed(k) }
func (k Key) IsJustReleased() bool { return IsJustReleased(k) }

// private =================================================================

const keyNone = -1

func checkOrder(k1, k2, k3 Key) bool {
	var d1 = inpututil.KeyPressDuration(ebiten.Key(k1))
	var d2 = inpututil.KeyPressDuration(ebiten.Key(k2))
	var d3 = inpututil.KeyPressDuration(ebiten.Key(k3))

	if k2 != keyNone && d1 <= d2 {
		return false
	}
	if k3 != keyNone && d2 <= d3 {
		return false
	}
	return true
}

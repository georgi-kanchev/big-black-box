package key

import (
	"big-black-box/packages/input/keyboard"
	
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	None                        = -1
	A              keyboard.Key = keyboard.Key(ebiten.KeyA)
	B              keyboard.Key = keyboard.Key(ebiten.KeyB)
	C              keyboard.Key = keyboard.Key(ebiten.KeyC)
	D              keyboard.Key = keyboard.Key(ebiten.KeyD)
	E              keyboard.Key = keyboard.Key(ebiten.KeyE)
	F              keyboard.Key = keyboard.Key(ebiten.KeyF)
	G              keyboard.Key = keyboard.Key(ebiten.KeyG)
	H              keyboard.Key = keyboard.Key(ebiten.KeyH)
	I              keyboard.Key = keyboard.Key(ebiten.KeyI)
	J              keyboard.Key = keyboard.Key(ebiten.KeyJ)
	K              keyboard.Key = keyboard.Key(ebiten.KeyK)
	L              keyboard.Key = keyboard.Key(ebiten.KeyL)
	M              keyboard.Key = keyboard.Key(ebiten.KeyM)
	N              keyboard.Key = keyboard.Key(ebiten.KeyN)
	O              keyboard.Key = keyboard.Key(ebiten.KeyO)
	P              keyboard.Key = keyboard.Key(ebiten.KeyP)
	Q              keyboard.Key = keyboard.Key(ebiten.KeyQ)
	R              keyboard.Key = keyboard.Key(ebiten.KeyR)
	S              keyboard.Key = keyboard.Key(ebiten.KeyS)
	T              keyboard.Key = keyboard.Key(ebiten.KeyT)
	U              keyboard.Key = keyboard.Key(ebiten.KeyU)
	V              keyboard.Key = keyboard.Key(ebiten.KeyV)
	W              keyboard.Key = keyboard.Key(ebiten.KeyW)
	X              keyboard.Key = keyboard.Key(ebiten.KeyX)
	Y              keyboard.Key = keyboard.Key(ebiten.KeyY)
	Z              keyboard.Key = keyboard.Key(ebiten.KeyZ)
	AltLeft        keyboard.Key = keyboard.Key(ebiten.KeyAltLeft)
	AltRight       keyboard.Key = keyboard.Key(ebiten.KeyAltRight)
	ArrowDown      keyboard.Key = keyboard.Key(ebiten.KeyArrowDown)
	ArrowLeft      keyboard.Key = keyboard.Key(ebiten.KeyArrowLeft)
	ArrowRight     keyboard.Key = keyboard.Key(ebiten.KeyArrowRight)
	ArrowUp        keyboard.Key = keyboard.Key(ebiten.KeyArrowUp)
	Backquote      keyboard.Key = keyboard.Key(ebiten.KeyBackquote)
	Backslash      keyboard.Key = keyboard.Key(ebiten.KeyBackslash)
	Backspace      keyboard.Key = keyboard.Key(ebiten.KeyBackspace)
	BracketLeft    keyboard.Key = keyboard.Key(ebiten.KeyBracketLeft)
	BracketRight   keyboard.Key = keyboard.Key(ebiten.KeyBracketRight)
	CapsLock       keyboard.Key = keyboard.Key(ebiten.KeyCapsLock)
	Comma          keyboard.Key = keyboard.Key(ebiten.KeyComma)
	ContextMenu    keyboard.Key = keyboard.Key(ebiten.KeyContextMenu)
	ControlLeft    keyboard.Key = keyboard.Key(ebiten.KeyControlLeft)
	ControlRight   keyboard.Key = keyboard.Key(ebiten.KeyControlRight)
	Delete         keyboard.Key = keyboard.Key(ebiten.KeyDelete)
	Digit0         keyboard.Key = keyboard.Key(ebiten.KeyDigit0)
	Digit1         keyboard.Key = keyboard.Key(ebiten.KeyDigit1)
	Digit2         keyboard.Key = keyboard.Key(ebiten.KeyDigit2)
	Digit3         keyboard.Key = keyboard.Key(ebiten.KeyDigit3)
	Digit4         keyboard.Key = keyboard.Key(ebiten.KeyDigit4)
	Digit5         keyboard.Key = keyboard.Key(ebiten.KeyDigit5)
	Digit6         keyboard.Key = keyboard.Key(ebiten.KeyDigit6)
	Digit7         keyboard.Key = keyboard.Key(ebiten.KeyDigit7)
	Digit8         keyboard.Key = keyboard.Key(ebiten.KeyDigit8)
	Digit9         keyboard.Key = keyboard.Key(ebiten.KeyDigit9)
	End            keyboard.Key = keyboard.Key(ebiten.KeyEnd)
	Enter          keyboard.Key = keyboard.Key(ebiten.KeyEnter)
	Equal          keyboard.Key = keyboard.Key(ebiten.KeyEqual)
	Escape         keyboard.Key = keyboard.Key(ebiten.KeyEscape)
	F1             keyboard.Key = keyboard.Key(ebiten.KeyF1)
	F2             keyboard.Key = keyboard.Key(ebiten.KeyF2)
	F3             keyboard.Key = keyboard.Key(ebiten.KeyF3)
	F4             keyboard.Key = keyboard.Key(ebiten.KeyF4)
	F5             keyboard.Key = keyboard.Key(ebiten.KeyF5)
	F6             keyboard.Key = keyboard.Key(ebiten.KeyF6)
	F7             keyboard.Key = keyboard.Key(ebiten.KeyF7)
	F8             keyboard.Key = keyboard.Key(ebiten.KeyF8)
	F9             keyboard.Key = keyboard.Key(ebiten.KeyF9)
	F10            keyboard.Key = keyboard.Key(ebiten.KeyF10)
	F11            keyboard.Key = keyboard.Key(ebiten.KeyF11)
	F12            keyboard.Key = keyboard.Key(ebiten.KeyF12)
	Home           keyboard.Key = keyboard.Key(ebiten.KeyHome)
	Insert         keyboard.Key = keyboard.Key(ebiten.KeyInsert)
	Minus          keyboard.Key = keyboard.Key(ebiten.KeyMinus)
	NumLock        keyboard.Key = keyboard.Key(ebiten.KeyNumLock)
	Numpad0        keyboard.Key = keyboard.Key(ebiten.KeyNumpad0)
	Numpad1        keyboard.Key = keyboard.Key(ebiten.KeyNumpad1)
	Numpad2        keyboard.Key = keyboard.Key(ebiten.KeyNumpad2)
	Numpad3        keyboard.Key = keyboard.Key(ebiten.KeyNumpad3)
	Numpad4        keyboard.Key = keyboard.Key(ebiten.KeyNumpad4)
	Numpad5        keyboard.Key = keyboard.Key(ebiten.KeyNumpad5)
	Numpad6        keyboard.Key = keyboard.Key(ebiten.KeyNumpad6)
	Numpad7        keyboard.Key = keyboard.Key(ebiten.KeyNumpad7)
	Numpad8        keyboard.Key = keyboard.Key(ebiten.KeyNumpad8)
	Numpad9        keyboard.Key = keyboard.Key(ebiten.KeyNumpad9)
	NumpadAdd      keyboard.Key = keyboard.Key(ebiten.KeyNumpadAdd)
	NumpadDecimal  keyboard.Key = keyboard.Key(ebiten.KeyNumpadDecimal)
	NumpadDivide   keyboard.Key = keyboard.Key(ebiten.KeyNumpadDivide)
	NumpadEnter    keyboard.Key = keyboard.Key(ebiten.KeyNumpadEnter)
	NumpadEqual    keyboard.Key = keyboard.Key(ebiten.KeyNumpadEqual)
	NumpadMultiply keyboard.Key = keyboard.Key(ebiten.KeyNumpadMultiply)
	NumpadSubtract keyboard.Key = keyboard.Key(ebiten.KeyNumpadSubtract)
	PageDown       keyboard.Key = keyboard.Key(ebiten.KeyPageDown)
	PageUp         keyboard.Key = keyboard.Key(ebiten.KeyPageUp)
	Pause          keyboard.Key = keyboard.Key(ebiten.KeyPause)
	Period         keyboard.Key = keyboard.Key(ebiten.KeyPeriod)
	PrintScreen    keyboard.Key = keyboard.Key(ebiten.KeyPrintScreen)
	Quote          keyboard.Key = keyboard.Key(ebiten.KeyQuote)
	ScrollLock     keyboard.Key = keyboard.Key(ebiten.KeyScrollLock)
	Semicolon      keyboard.Key = keyboard.Key(ebiten.KeySemicolon)
	ShiftLeft      keyboard.Key = keyboard.Key(ebiten.KeyShiftLeft)
	ShiftRight     keyboard.Key = keyboard.Key(ebiten.KeyShiftRight)
	Slash          keyboard.Key = keyboard.Key(ebiten.KeySlash)
	Space          keyboard.Key = keyboard.Key(ebiten.KeySpace)
	Tab            keyboard.Key = keyboard.Key(ebiten.KeyTab)
	Alt            keyboard.Key = keyboard.Key(ebiten.KeyAlt)
	Control        keyboard.Key = keyboard.Key(ebiten.KeyControl)
	Shift          keyboard.Key = keyboard.Key(ebiten.KeyShift)
)

func FromName(name string) keyboard.Key {
	var k, _ = nameToKey[name]
	return k
}
func ToName(key keyboard.Key) string {
	name, _ := keyToName[key]
	return name
}

//=================================================================

var nameToKey = map[string]keyboard.Key{
	"None": None, "A": A, "B": B, "C": C, "D": D, "E": E, "F": F, "G": G, "H": H, "I": I, "J": J, "K": K, "L": L, "M": M,
	"N": N, "O": O, "P": P, "Q": Q, "R": R, "S": S, "T": T, "U": U, "V": V, "W": W, "X": X, "Y": Y, "Z": Z,
	"AltLeft": AltLeft, "AltRight": AltRight,
	"ArrowDown": ArrowDown, "ArrowLeft": ArrowLeft, "ArrowRight": ArrowRight, "ArrowUp": ArrowUp,
	"Backquote": Backquote, "Backslash": Backslash, "Backspace": Backspace,
	"BracketLeft": BracketLeft, "BracketRight": BracketRight,
	"CapsLock": CapsLock, "Comma": Comma, "ContextMenu": ContextMenu, "ControlLeft": ControlLeft, "ControlRight": ControlRight,
	"Delete": Delete, "Digit0": Digit0, "Digit1": Digit1, "Digit2": Digit2, "Digit3": Digit3, "Digit4": Digit4, "Digit5": Digit5,
	"Digit6": Digit6, "Digit7": Digit7, "Digit8": Digit8, "Digit9": Digit9,
	"End": End, "Enter": Enter, "Equal": Equal, "Escape": Escape,
	"F1": F1, "F2": F2, "F3": F3, "F4": F4, "F5": F5, "F6": F6, "F7": F7, "F8": F8, "F9": F9, "F10": F10, "F11": F11, "F12": F12,
	"Home": Home, "Insert": Insert, "Minus": Minus, "NumLock": NumLock,
	"Numpad0": Numpad0, "Numpad1": Numpad1, "Numpad2": Numpad2, "Numpad3": Numpad3, "Numpad4": Numpad4, "Numpad5": Numpad5,
	"Numpad6": Numpad6, "Numpad7": Numpad7, "Numpad8": Numpad8, "Numpad9": Numpad9,
	"NumpadAdd": NumpadAdd, "NumpadDecimal": NumpadDecimal, "NumpadDivide": NumpadDivide, "NumpadEnter": NumpadEnter,
	"NumpadEqual": NumpadEqual, "NumpadMultiply": NumpadMultiply, "NumpadSubtract": NumpadSubtract,
	"PageDown": PageDown, "PageUp": PageUp, "Pause": Pause, "Period": Period, "PrintScreen": PrintScreen, "Quote": Quote,
	"ScrollLock": ScrollLock, "Semicolon": Semicolon,
	"ShiftLeft": ShiftLeft, "ShiftRight": ShiftRight, "Slash": Slash, "Space": Space, "Tab": Tab,
	"Alt": Alt, "Control": Control, "Shift": Shift,
}

var keyToName = map[keyboard.Key]string{
	None: "None", A: "A", B: "B", C: "C", D: "D", E: "E", F: "F", G: "G", H: "H", I: "I", J: "J", K: "K", L: "L", M: "M",
	N: "N", O: "O", P: "P", Q: "Q", R: "R", S: "S", T: "T", U: "U", V: "V", W: "W", X: "X", Y: "Y", Z: "Z",
	AltLeft: "AltLeft", AltRight: "AltRight",
	ArrowDown: "ArrowDown", ArrowLeft: "ArrowLeft", ArrowRight: "ArrowRight", ArrowUp: "ArrowUp",
	Backquote: "Backquote", Backslash: "Backslash", Backspace: "Backspace",
	BracketLeft: "BracketLeft", BracketRight: "BracketRight",
	CapsLock: "CapsLock", Comma: "Comma", ContextMenu: "ContextMenu", ControlLeft: "ControlLeft", ControlRight: "ControlRight",
	Delete: "Delete", Digit0: "Digit0", Digit1: "Digit1", Digit2: "Digit2", Digit3: "Digit3", Digit4: "Digit4", Digit5: "Digit5",
	Digit6: "Digit6", Digit7: "Digit7", Digit8: "Digit8", Digit9: "Digit9",
	End: "End", Enter: "Enter", Equal: "Equal", Escape: "Escape",
	F1: "F1", F2: "F2", F3: "F3", F4: "F4", F5: "F5", F6: "F6", F7: "F7", F8: "F8", F9: "F9", F10: "F10", F11: "F11", F12: "F12",
	Home: "Home", Insert: "Insert", Minus: "Minus", NumLock: "NumLock",
	Numpad0: "Numpad0", Numpad1: "Numpad1", Numpad2: "Numpad2", Numpad3: "Numpad3", Numpad4: "Numpad4", Numpad5: "Numpad5",
	Numpad6: "Numpad6", Numpad7: "Numpad7", Numpad8: "Numpad8", Numpad9: "Numpad9",
	NumpadAdd: "NumpadAdd", NumpadDecimal: "NumpadDecimal", NumpadDivide: "NumpadDivide", NumpadEnter: "NumpadEnter",
	NumpadEqual: "NumpadEqual", NumpadMultiply: "NumpadMultiply", NumpadSubtract: "NumpadSubtract",
	PageDown: "PageDown", PageUp: "PageUp", Pause: "Pause", Period: "Period", PrintScreen: "PrintScreen", Quote: "Quote",
	ScrollLock: "ScrollLock", Semicolon: "Semicolon",
	ShiftLeft: "ShiftLeft", ShiftRight: "ShiftRight", Slash: "Slash", Space: "Space", Tab: "Tab",
	Alt: "Alt", Control: "Control", Shift: "Shift",
}

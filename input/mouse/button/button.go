package button

import (
	"big-black-box/input/mouse"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	Left   mouse.Button = mouse.Button(ebiten.MouseButtonLeft)
	Middle mouse.Button = mouse.Button(ebiten.MouseButtonMiddle)
	Right  mouse.Button = mouse.Button(ebiten.MouseButtonRight)
	Extra1 mouse.Button = mouse.Button(ebiten.MouseButton3)
	Extra2 mouse.Button = mouse.Button(ebiten.MouseButton4)
)

func FromName(name string) mouse.Button {
	var b, _ = nameToButton[name]
	return b
}
func ToName(button mouse.Button) string {
	var name, _ = buttonToName[button]
	return name
}

//=================================================================

var nameToButton = map[string]mouse.Button{"Left": Left, "Middle": Middle, "Right": Right, "Extra1": Extra1, "Extra2": Extra2}
var buttonToName = map[mouse.Button]string{Left: "Left", Middle: "Middle", Right: "Right", Extra1: "Extra1", Extra2: "Extra2"}

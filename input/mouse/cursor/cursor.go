package cursor

import (
	"big-black-box/input/mouse"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	Default    mouse.Cursor = mouse.Cursor(ebiten.CursorShapeDefault)
	Arrow                   = Default
	Text       mouse.Cursor = mouse.Cursor(ebiten.CursorShapeText)
	Crosshair  mouse.Cursor = mouse.Cursor(ebiten.CursorShapeCrosshair)
	Pointer    mouse.Cursor = mouse.Cursor(ebiten.CursorShapePointer)
	Hand                    = Pointer
	ResizeEW   mouse.Cursor = mouse.Cursor(ebiten.CursorShapeEWResize)
	ResizeNS   mouse.Cursor = mouse.Cursor(ebiten.CursorShapeNSResize)
	ResizeNESW mouse.Cursor = mouse.Cursor(ebiten.CursorShapeNESWResize)
	ResizeNWSE mouse.Cursor = mouse.Cursor(ebiten.CursorShapeNWSEResize)
	Move       mouse.Cursor = mouse.Cursor(ebiten.CursorShapeMove)
	NotAllowed mouse.Cursor = mouse.Cursor(ebiten.CursorShapeNotAllowed)
)

package window

import (
	"big-black-box/internal"

	"github.com/hajimehoshi/ebiten/v2"
)

type Mode byte

const (
	ModeFloating Mode = iota
	ModeMaximized
	ModeFullscreen
)

func GetTitle() string {
	return internal.State.Window.Title
}
func SetTitle(title string) {
	internal.State.Window.Title = title
	ebiten.SetWindowTitle(title)
}

func IsVsynced() bool {
	return internal.State.Window.IsVsynced
}
func SetVsync(enabled bool) {
	internal.State.Window.IsVsynced = true
	ebiten.SetVsyncEnabled(enabled)
}

func GetMode() Mode {
	return Mode(internal.State.Window.Mode)
}
func SetMode(mode Mode) {
	internal.State.Window.Mode = byte(mode)
	if mode != ModeFullscreen && ebiten.IsFullscreen() {
		ebiten.SetFullscreen(false)
	}

	switch mode {
	case ModeFloating:
		ebiten.RestoreWindow()
	case ModeMaximized:
		ebiten.MaximizeWindow()
	case ModeFullscreen:
		ebiten.SetFullscreen(true)
	}
}

func SetMonitor(monitor byte) {
	var monitors []*ebiten.MonitorType
	monitors = ebiten.AppendMonitors(monitors)

	if int(monitor) >= len(monitors) {
		monitor = 0
	}

	internal.State.Window.Monitor = monitor
	ebiten.SetMonitor(monitors[monitor])
}
func GetMonitor() byte {
	var monitors []*ebiten.MonitorType
	monitors = ebiten.AppendMonitors(monitors)

	var cur = ebiten.Monitor()
	for i, m := range monitors {
		if m == cur {
			return byte(i)
		}
	}
	return 0
}

func IsHovered() bool {
	var x, y = ebiten.CursorPosition()
	var w, h = internal.State.Layout(ebiten.WindowSize())
	return x >= 0 && y >= 0 && x < w && y < h
}
func IsFocused() bool {
	return ebiten.IsFocused()
}

func GetPixelScale() float32 {
	return internal.State.Window.PixelScale
}
func SetPixelScale(pixelScale float32) {
	internal.State.Window.PixelScale = pixelScale
}

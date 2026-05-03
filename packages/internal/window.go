// tracks the window package data and maps it to ebiten when it changes

package internal

import (
	"big-black-box/packages/window"

	"github.com/hajimehoshi/ebiten/v2"
)

var prevTitle string
var prevVSync bool
var prevMode, prevMonitor byte = 255, 255 // default to 255 to trigger an update

func readWindowValues() {
	windowSizeDuctapeFix()

	if prevTitle != window.Title {
		prevTitle = window.Title
		ebiten.SetWindowTitle(prevTitle)
	}
	if prevVSync != window.VSync {
		prevVSync = window.VSync
		ebiten.SetVsyncEnabled(prevVSync)
	}
	if prevMonitor != window.Monitor {
		var monitors []*ebiten.MonitorType
		monitors = ebiten.AppendMonitors(monitors)

		if int(prevMonitor) >= len(monitors) {
			prevMonitor = 0
		}

		prevMonitor = window.Monitor
		ebiten.SetMonitor(monitors[prevMonitor])
	}
	if prevMode != byte(window.CurrentMode) {
		prevMode = byte(window.CurrentMode)
		if prevMode != 2 && ebiten.IsFullscreen() {
			ebiten.SetFullscreen(false)
		}

		switch prevMode {
		case 0: // floating
			if ebiten.IsWindowMaximized() || ebiten.IsWindowMinimized() {
				ebiten.RestoreWindow()
			}
		case 1: // maximized
			ebiten.MaximizeWindow()
		case 2: // fullscreen
			ebiten.SetFullscreen(true)
		}
	}

	var w, h = Engine.Layout(ebiten.WindowSize())
	window.Width, window.Height = float32(w), float32(h)
}

func windowSizeDuctapeFix() {
	if ebiten.Tick() == 1 && window.CurrentMode > 0 {
		ebiten.SetFullscreen(false)
		if ebiten.IsWindowMaximized() || ebiten.IsWindowMinimized() {
			ebiten.RestoreWindow()
		}
		prevMode = 255
	} // for some reason, ebiten doesn't update the window size immediately
} // and requires a resize after init

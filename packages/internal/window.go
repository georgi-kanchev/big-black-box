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
	if prevTitle != window.Title {
		prevTitle = window.Title
		ebiten.SetWindowTitle(prevTitle)
	}
	if prevVSync != window.VSync {
		prevVSync = window.VSync
		ebiten.SetVsyncEnabled(prevVSync)
	}
	if prevMode != byte(window.CurrentMode) {
		prevMode = byte(window.CurrentMode)
		if prevMode != 2 && ebiten.IsFullscreen() {
			ebiten.SetFullscreen(false)
		}

		switch prevMode {
		case 0: // floating
			ebiten.RestoreWindow()
		case 1: // maximized
			ebiten.MaximizeWindow()
		case 2: // fullscreen
			ebiten.SetFullscreen(true)
		}
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
}

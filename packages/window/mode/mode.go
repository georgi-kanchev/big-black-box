// not inside the parent package for a shorter API & non-cluttered autocomplete

package mode

import "big-black-box/packages/window"

const (
	Floating window.Mode = iota
	Maximized
	Fullscreen
)

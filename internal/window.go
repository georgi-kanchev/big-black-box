package internal

type Window struct {
	Title   string
	Mode    byte // see window.Mode
	Monitor byte

	IsMaximized, IsVsynced bool
}

package internal

import "strings"

// Global stack of builders to allow for nested string building without
// mid-operation interference. This is throwaway state.
var builders [8]strings.Builder
var builderDepth int

// BuilderPush starts a new building operation at a deeper level.
func BuilderPush() {
	builders[builderDepth].Reset()
	builderDepth++
}

// BuilderPop ends the current building operation and returns to the previous level.
func BuilderPop() {
	builderDepth--
}

// BuilderWriteString writes a string to the current builder level.
func BuilderWriteString(s string) {
	builders[builderDepth-1].WriteString(s)
}

// BuilderWriteByte writes a byte to the current builder level.
func BuilderWriteByte(b byte) {
	builders[builderDepth-1].WriteByte(b)
}

// BuilderWriteRune writes a rune to the current builder level.
func BuilderWriteRune(r rune) {
	builders[builderDepth-1].WriteRune(r)
}

// BuilderResult returns the accumulated string from the current builder level.
func BuilderResult() string {
	return builders[builderDepth-1].String()
}

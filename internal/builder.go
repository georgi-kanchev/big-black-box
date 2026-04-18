package internal

import (
	"strconv"
	"strings"
)

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

func BuilderWriteInt(i int64) {
	var buf [20]byte
	var b = strconv.AppendInt(buf[:0], i, 10)
	builders[builderDepth-1].Write(b)
}
func BuilderWriteFloat(f float64, prec int) {
	var buf [32]byte
	var b = strconv.AppendFloat(buf[:0], f, 'f', prec, 64)
	builders[builderDepth-1].Write(b)
}

// BuilderWriteString writes a string to the current builder level.
func BuilderWriteString(value string) {
	builders[builderDepth-1].WriteString(value)
}

// BuilderWriteByte writes a byte to the current builder level.
func BuilderWriteByte(value byte) {
	builders[builderDepth-1].WriteByte(value)
}

// BuilderWriteRune writes a rune to the current builder level.
func BuilderWriteRune(value rune) {
	builders[builderDepth-1].WriteRune(value)
}

// BuilderResult returns the accumulated string from the current builder level.
func BuilderResult() string {
	return builders[builderDepth-1].String()
}

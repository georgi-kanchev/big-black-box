// Helper functions for bit-masks. Useful for storing up to 64 flags (bool values) in a single integer
// where each bit represents each flag (on/off).
package flag

import "big-black-box/utility/number"

func IsOn[T number.Integer](allFlags, flag T) bool {
	return allFlags&flag != 0
}
func TurnOn[T number.Integer](allFlags, flag T) T {
	return allFlags | flag
}
func Toggle[T number.Integer](allFlags, flag T) T {
	return allFlags ^ flag
}
func TurnOff[T number.Integer](allFlags, flag T) T {
	return allFlags &^ flag
}

func FromBit[T number.Integer](bitPosition int) T {
	return T(1 << bitPosition)
}

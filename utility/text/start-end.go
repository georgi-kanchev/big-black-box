package text

import (
	"strconv"
	"strings"
)

func Pad(text string, length int, pad string) string {
	var textLen = Length(text)
	var spaces = length - textLen
	if spaces <= 0 {
		return text
	}
	var left = spaces / 2
	return PadRight(PadLeft(text, textLen+left, pad), length, pad)
}
func PadLeft(text string, length int, pad string) string {
	var textLen = Length(text)
	var padding = length - textLen
	if padding <= 0 || pad == "" {
		return text
	}
	return repeatPad(pad, padding) + text
}
func PadRight(text string, length int, pad string) string {
	var textLen = Length(text)
	var padding = length - textLen
	if padding <= 0 || pad == "" {
		return text
	}
	return text + repeatPad(pad, padding)
}
func PadZeros(number float32, count int) string {
	if count == 0 {
		return Start().Float(number).End()
	}

	if count < 0 {
		var width = -count
		var s = strconv.Itoa(int(number))
		for Length(s) < width {
			s = "0" + s
		}
		return s
	}

	return strconv.FormatFloat(float64(number), 'f', count, 32)
}

func Trim(text string) string {
	return TrimEnd(TrimStart(text))
}
func TrimStart(text string) string {
	return strings.TrimLeft(text, " \r\n")
}
func TrimEnd(text string) string {
	return strings.TrimRight(text, " \r\n")
}

// Surrounds a text with the given start part and end part.
// If end part is empty, it uses the start part for both sides.
func SurroundWith(text, startPart string, endPart string) string {
	if endPart == "" {
		endPart = startPart
	}
	return startPart + text + endPart
}

// Adds the part to the start of the text only if it doesn't already have it.
func EnsureStart(text, part string) string {
	if !strings.HasPrefix(text, part) {
		return part + text
	}
	return text
}

// Adds the part to the end of the text only if it doesn't already have it.
func EnsureEnd(text, part string) string {
	if !strings.HasSuffix(text, part) {
		return text + part
	}
	return text
}

// Removes the given start part and end part only if both are present.
// If end part is empty, it looks for the start part on both sides.
func Chop(text, startPart string, endPart string) string {
	if endPart == "" {
		endPart = startPart
	}

	if strings.HasPrefix(text, startPart) && strings.HasSuffix(text, endPart) {
		return text[len(startPart) : len(text)-len(endPart)]
	}
	return text
}

// Removes the start part from the text if it exists.
func ChopStart(text, part string) string {
	return strings.TrimPrefix(text, part)
}

// Removes the end part from the text if it exists.
func ChopEnd(text, part string) string {
	return strings.TrimSuffix(text, part)
}

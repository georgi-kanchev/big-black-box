package text

import (
	"big-black-box/internal"
	"strconv"
)

// Text is a primitive-based alias that allows method chaining for string construction.
type Text int

// Start starts a new string building operation.
// It must always be paired with .End() to free the internal text.
func Start() Text {
	internal.BuilderPush()
	return 0
}

// String appends a string to the text.
func (t Text) String(v string) Text {
	internal.BuilderWriteString(v)
	return t
}

// Int appends an integer to the text.
func (t Text) Int(v int) Text {
	internal.BuilderWriteString(strconv.Itoa(v))
	return t
}

// Float appends a float32 to the text.
func (t Text) Float(v float32) Text {
	internal.BuilderWriteString(strconv.FormatFloat(float64(v), 'f', -1, 32))
	return t
}

// Bool appends a boolean to the text.
func (t Text) Bool(v bool) Text {
	if v {
		internal.BuilderWriteString("true")
	} else {
		internal.BuilderWriteString("false")
	}
	return t
}

// Rune appends a rune to the text.
func (t Text) Rune(v rune) Text {
	internal.BuilderWriteRune(v)
	return t
}

// End finishes the building operation and returns the resulting string.
func (t Text) End() string {
	var s = internal.BuilderResult()
	internal.BuilderPop()
	return s
}

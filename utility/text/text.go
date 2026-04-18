// Another important package, similar to the number one. It has many helper functions that operate on a string -
// transformations, checks, formatting or executing a result on it. Also wraps some standard string functions
// to make them more digestible and clarify their API.
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
func (t Text) String(value string) Text {
	internal.BuilderWriteString(value)
	return t
}

// Int appends an integer to the text.
func (t Text) Int(value int) Text {
	internal.BuilderWriteString(strconv.Itoa(value))
	return t
}

// Float appends a float32 to the text.
func (t Text) Float(value float32) Text {
	internal.BuilderWriteString(strconv.FormatFloat(float64(value), 'f', -1, 32))
	return t
}

// Bool appends a boolean to the text.
func (t Text) Bool(value bool) Text {
	if value {
		internal.BuilderWriteString("true")
	} else {
		internal.BuilderWriteString("false")
	}
	return t
}

// Rune appends a rune to the text.
func (t Text) Rune(value rune) Text {
	internal.BuilderWriteRune(value)
	return t
}

// End finishes the building operation and returns the resulting string.
func (t Text) End() string {
	var s = internal.BuilderResult()
	internal.BuilderPop()
	return s
}

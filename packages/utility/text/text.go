// Another important package, similar to the number one. It has many helper functions that operate on a string -
// transformations, checks, formatting or executing a result on it. Also wraps some standard string functions
// to make them more digestible and clarify their API.
package text

// Text is a primitive-based alias that allows method chaining for string construction and
// wraps the Builder functionality.
type Text int

// Start starts a new string building operation.
// It must always be paired with .End() to free the internal text.
func Start() Text {
	BuilderPush()
	return 0
}

// String appends a string to the text.
func (t Text) String(value string) Text {
	BuilderWriteString(value)
	return t
}

// Int appends an integer to the text.
func (t Text) Int(value int) Text {
	BuilderWriteInt(int64(value))
	return t
}

// Float appends a float32 to the text.
func (t Text) Float(value float32) Text {
	BuilderWriteFloat(float64(value), -1)
	return t
}

// Bool appends a boolean to the text.
func (t Text) Bool(value bool) Text {
	if value {
		BuilderWriteString("true")
	} else {
		BuilderWriteString("false")
	}
	return t
}

// Rune appends a rune to the text.
func (t Text) Rune(value rune) Text {
	BuilderWriteRune(value)
	return t
}

// End finishes the building operation and returns the resulting string.
func (t Text) End() string {
	var s = BuilderResult()
	BuilderPop()
	return s
}

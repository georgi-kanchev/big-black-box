package text

import (
	"big-black-box/internal"
	"big-black-box/utility/number"
	b64 "encoding/base64"
	"math"
	"strconv"
	"strings"
)

func ToNumber[T number.Number](text string) T {
	var res T
	switch p := any(&res).(type) {
	case *float32:
		if f, err := strconv.ParseFloat(text, 32); err == nil {
			*p = float32(f)
		} else {
			*p = float32(math.NaN())
		}
	case *float64:
		if f, err := strconv.ParseFloat(text, 64); err == nil {
			*p = f
		} else {
			*p = math.NaN()
		}
	case *int:
		if i, err := strconv.ParseInt(text, 10, 0); err == nil {
			*p = int(i)
		}
	case *int64:
		if i, err := strconv.ParseInt(text, 10, 64); err == nil {
			*p = i
		}
	case *int32:
		if i, err := strconv.ParseInt(text, 10, 32); err == nil {
			*p = int32(i)
		}
	case *int16:
		if i, err := strconv.ParseInt(text, 10, 16); err == nil {
			*p = int16(i)
		}
	case *int8:
		if i, err := strconv.ParseInt(text, 10, 8); err == nil {
			*p = int8(i)
		}
	case *uint:
		if u, err := strconv.ParseUint(text, 10, 0); err == nil {
			*p = uint(u)
		}
	case *uint64:
		if u, err := strconv.ParseUint(text, 10, 64); err == nil {
			*p = u
		}
	case *uint32:
		if u, err := strconv.ParseUint(text, 10, 32); err == nil {
			*p = uint32(u)
		}
	case *uint16:
		if u, err := strconv.ParseUint(text, 10, 16); err == nil {
			*p = uint16(u)
		}
	case *uint8:
		if u, err := strconv.ParseUint(text, 10, 8); err == nil {
			*p = uint8(u)
		}
	}
	return res
}
func FormatByteSize(byteSize int) string {
	const unit = 1024

	internal.BuilderPush()
	defer internal.BuilderPop()

	if byteSize < unit {
		// strconv.AppendInt can write to a byte slice,
		// but since we need to use your BuilderWriteString:
		internal.BuilderWriteInt(int64(byteSize))
		internal.BuilderWriteString(" B")
		return internal.BuilderResult()
	}

	const units = "KMGTPE"
	var div, exp = int64(unit), 0
	for n := int64(byteSize) / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	var val = float64(byteSize) / float64(div)

	// Formatting the float without fmt.Sprintf
	// We use strconv.FormatFloat to get a string then write it to the builder
	internal.BuilderWriteFloat(val, 3)
	internal.BuilderWriteByte(' ')
	internal.BuilderWriteByte(units[exp])
	internal.BuilderWriteByte('B')

	return internal.BuilderResult()
}

func ToBase64(text string) string {
	return b64.StdEncoding.EncodeToString([]byte(text))
}
func FromBase64(base64 string) string {
	var decodedBytes, err = b64.StdEncoding.DecodeString(base64)
	if err != nil {
		return ""
	}
	return string(decodedBytes)
}

func Token(text, divider string, index int) string {
	if index < 0 {
		return ""
	}

	var start = 0
	for i := 0; i < index; i++ {
		var pos = strings.Index(text[start:], divider)
		if pos == -1 {
			return "" // Index out of bounds
		}
		start += pos + len(divider)
	}

	var end = strings.Index(text[start:], divider)
	if end == -1 {
		return text[start:] // Last token in string
	}
	return text[start : start+end]
}

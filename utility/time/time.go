// Helper functions for time. Provides time conversions. As well as runtime and real time stats.
// Provides a way to slow down or speed up those runtime stats. All values are in seconds, unless
// a unit is explicitly specified otherwise. Useful for when dealing with motion, debugging or displaying
// timers/clocks etc.
package time

import (
	"big-black-box/internal"
	"big-black-box/utility/flag"
	"big-black-box/utility/number"
	"big-black-box/utility/time/unit"
	"strconv"
	"time"
)

func AsClock24(seconds float32, divider string, units int) string {
	var ts = time.Duration(seconds * float32(time.Second))
	return formatTimeParts(ts, divider, units, false, false)
}
func AsClock12(seconds float32, divider string, units int, amPm bool) string {
	var ts = time.Duration(seconds * float32(time.Second))
	return formatTimeParts(ts, divider, units, true, amPm)
}

//=================================================================

func FPS() float32       { return internal.FPS }
func TPS() float32       { return internal.TPS }
func FrameCount() uint64 { return internal.FrameCount }
func Running() float32   { return internal.Runtime }
func Clock() float32     { return internal.Clock }

func ToMilliseconds(seconds float32) float32 { return seconds * 1000 }
func ToMinutes(secodns float32) float32      { return secodns / 60 }
func ToHours(seconds float32) float32        { return seconds / 3600 }
func ToDays(seconds float32) float32         { return seconds / 86400 }
func ToWeeks(seconds float32) float32        { return seconds / 604800 }

func FromMilliseconds(milliseconds float32) float32 { return milliseconds / 1000 }
func FromMinutes(minutes float32) float32           { return minutes * 60 }
func FromHours(hours float32) float32               { return hours * 3600 }
func FromDays(days float32) float32                 { return days * 86400 }
func FromWeeks(weeks float32) float32               { return weeks * 604800 }

// private ========================================================

func formatTimeParts(ts time.Duration, divider string, units int, is12Hour, amPm bool) string {
	internal.BuilderPush()
	defer internal.BuilderPop()

	var counter = 0

	if flag.IsOn(units, unit.Day) {
		var val = int(ts.Hours() / 24)
		writePaddedInt(val, 2) // manual padding to avoid temporary string allocations
		counter++
	}

	if flag.IsOn(units, unit.Hour) {
		writeSep(counter, divider)
		var val int
		var h = int((ts % (24 * time.Hour)) / time.Hour)
		if is12Hour {
			val = int(number.Wrap(float32(h), 0, 12))
		} else {
			val = h
		}
		writePaddedInt(val, 2)
		counter++
	}

	if flag.IsOn(units, unit.Minute) {
		writeSep(counter, divider)
		var val = int((ts % time.Hour) / time.Minute)
		writePaddedInt(val, 2)
		counter++
	}

	if flag.IsOn(units, unit.Second) {
		writeSep(counter, divider)
		var val = int((ts % time.Minute) / time.Second)
		writePaddedInt(val, 2)
		counter++
	}

	if flag.IsOn(units, unit.Millisecond) {
		var val = int((ts % time.Second) / time.Millisecond)
		if flag.IsOn(units, unit.Second) {
			internal.BuilderWriteByte('.')
		} else if counter > 0 {
			internal.BuilderWriteString(divider)
		}
		internal.BuilderWriteString(strconv.Itoa(val))
		counter++
	}

	if is12Hour && amPm {
		internal.BuilderWriteByte(' ')
		if int(ts.Hours())%24 >= 12 {
			internal.BuilderWriteString("PM")
		} else {
			internal.BuilderWriteString("AM")
		}
	}

	return internal.BuilderResult()
}
func writePaddedInt(val int, width int) {
	var str = strconv.Itoa(val)
	for i := len(str); i < width; i++ {
		internal.BuilderWriteByte('0')
	}
	internal.BuilderWriteString(str)
}
func writeSep(counter int, divider string) {
	if counter > 0 {
		internal.BuilderWriteString(divider)
	}
}

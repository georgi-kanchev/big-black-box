package debug

import (
	"big-black-box/utility/text"
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2"
)

func LinesOfCode() string {
	var directory, _ = os.Getwd()
	var cmd = exec.Command("bash", "-c", fmt.Sprintf(`find "%s" -name "*.go" -type f -exec wc -l {} +`, directory))
	var cmdOut bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdOut
	if err := cmd.Run(); err != nil {
		return ""
	}

	var results = make(map[string]int)
	var scanner = bufio.NewScanner(&cmdOut)
	for scanner.Scan() {
		var line = strings.TrimSpace(scanner.Text())
		var parts = strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		var count, _ = strconv.ParseInt(parts[0], 10, 32)
		var path = parts[1]
		var rel, _ = filepath.Rel(directory, path)
		results[rel] = int(count)
	}

	var dirTotals = make(map[string]int)
	for path, count := range results {
		dirTotals[path] = count
		var dir = filepath.Dir(path)
		for dir != path {
			dirTotals[dir] += count
			dir = filepath.Dir(dir)
			if dir == "." {
				break
			}
		}
	}

	var allPaths []string
	for p := range dirTotals {
		if p != "." {
			allPaths = append(allPaths, p)
		}
	}
	sort.Strings(allPaths)

	var t = text.Start().String("Lines of code in:\n")

	var printTree func(path, prefix string, isLast bool)
	printTree = func(path, prefix string, isLast bool) {
		var connector = "├"
		if isLast {
			connector = "└"
		}

		var name = filepath.Base(path)
		var displayCount = ""
		if _, ok := results[path]; ok {
			displayCount = strconv.Itoa(dirTotals[path])
		} else {
			displayCount = "[" + strconv.Itoa(dirTotals[path]) + "]"
		}

		if name == "." {
			t.String("[").String(displayCount).String("] ").String(directory).String("\n")
		} else {
			// Basic padding for displayCount
			var padding = ""
			for i := 0; i < 6-len(displayCount); i++ {
				padding += " "
			}
			t.String(padding).String(displayCount).String(" ").String(prefix).String(connector).String(name).String("\n")
		}

		var children []string
		for _, p := range allPaths {
			if filepath.Dir(p) == path {
				children = append(children, p)
			}
		}
		sort.Strings(children)

		for i, c := range children {
			var newPrefix = prefix
			if isLast {
				newPrefix += "  "
			} else {
				newPrefix += "│ "
			}
			printTree(c, newPrefix, i == len(children)-1)
		}
	}

	var topLevel []string
	for _, p := range allPaths {
		if !strings.Contains(p, string(filepath.Separator)) {
			topLevel = append(topLevel, p)
		}
	}
	sort.Strings(topLevel)
	for i, tl := range topLevel {
		printTree(tl, "", i == len(topLevel)-1)
	}

	return t.End()
}
func MemoryUsage() string {
	if ebiten.Tick()%20 != 0 {
		return str
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memBuf = memBuf[:0]

	memBuf = append(memBuf, "Memory:\n"...)
	memBuf = append(memBuf, "UsedNow = "...)
	memBuf = appendByteSize(memBuf, int(m.Alloc))
	memBuf = append(memBuf, " (current heap in use)\n"...)
	memBuf = append(memBuf, "UsedTotal = "...)
	memBuf = appendByteSize(memBuf, int(m.TotalAlloc))
	memBuf = append(memBuf, " (total allocated since start)\n"...)
	memBuf = append(memBuf, "FromOS = "...)
	memBuf = appendByteSize(memBuf, int(m.Sys))
	memBuf = append(memBuf, " (memory reserved from OS)\n"...)

	memBuf = append(memBuf, "\nHeap:\n"...)
	memBuf = append(memBuf, "Used = "...)
	memBuf = appendByteSize(memBuf, int(m.HeapAlloc))
	memBuf = append(memBuf, " \n"...)
	memBuf = append(memBuf, "Reserved = "...)
	memBuf = appendByteSize(memBuf, int(m.HeapSys))
	memBuf = append(memBuf, " \n"...)
	memBuf = append(memBuf, "Idle = "...)
	memBuf = appendByteSize(memBuf, int(m.HeapIdle))
	memBuf = append(memBuf, " (not used but still reserved)\n"...)
	memBuf = append(memBuf, "Active = "...)
	memBuf = appendByteSize(memBuf, int(m.HeapInuse))
	memBuf = append(memBuf, " (actively in use)\n"...)
	memBuf = append(memBuf, "Released = "...)
	memBuf = appendByteSize(memBuf, int(m.HeapReleased))
	memBuf = append(memBuf, " (given back to OS)\n"...)

	memBuf = append(memBuf, "\nStack:\n"...)
	memBuf = append(memBuf, "Used = "...)
	memBuf = appendByteSize(memBuf, int(m.StackInuse))
	memBuf = append(memBuf, "\n"...)
	memBuf = append(memBuf, "Reserved = "...)
	memBuf = appendByteSize(memBuf, int(m.StackSys))
	memBuf = append(memBuf, "\n"...)
	memBuf = append(memBuf, "Other = "...)
	memBuf = appendByteSize(memBuf, int(m.OtherSys))
	memBuf = append(memBuf, " (misc runtime overhead)\n"...)

	memBuf = append(memBuf, "\nObjects:\n"...)
	memBuf = append(memBuf, "Allocs = "...)
	memBuf = appendThousands(memBuf, m.Mallocs)
	memBuf = append(memBuf, " (objects allocated)\n"...)
	memBuf = append(memBuf, "Frees = "...)
	memBuf = appendThousands(memBuf, m.Frees)
	memBuf = append(memBuf, " (objects freed)\n"...)
	memBuf = append(memBuf, "Live = "...)
	memBuf = appendThousands(memBuf, m.HeapObjects)
	memBuf = append(memBuf, " (currently alive)\n"...)

	memBuf = append(memBuf, "\nGarbage Collection:\n"...)
	memBuf = append(memBuf, "Total = "...)
	memBuf = appendThousands(memBuf, uint64(m.NumGC))
	memBuf = append(memBuf, " (total collections)\n"...)
	memBuf = append(memBuf, "Forced = "...)
	memBuf = strconv.AppendUint(memBuf, uint64(m.NumForcedGC), 10)
	memBuf = append(memBuf, " (manual triggers)\n"...)
	memBuf = append(memBuf, "Next = "...)
	memBuf = appendByteSize(memBuf, int(m.NextGC))
	memBuf = append(memBuf, " (target heap size of the next GC)\n"...)
	memBuf = append(memBuf, "PauseTotal = "...)
	memBuf = strconv.AppendFloat(memBuf, float64(m.PauseTotalNs)/1e9, 'f', 2, 64)
	memBuf = append(memBuf, " s (total time spent in GC)\n"...)
	if m.LastGC == 0 {
		memBuf = append(memBuf, "SinceLast = never\n"...)
	} else {
		memBuf = append(memBuf, "SinceLast = "...)
		memBuf = strconv.AppendFloat(memBuf, time.Since(time.Unix(0, int64(m.LastGC))).Seconds(), 'f', 2, 64)
		memBuf = append(memBuf, " s\n"...)
	}

	str = unsafe.String(unsafe.SliceData(memBuf), len(memBuf))
	return str
}
func ProfileAllocations(seconds float32) {
	go func() {
		var ts = time.Now().Format("2006-01-02_15-04-05")
		var profileFile = fmt.Sprintf("allocs_%s.prof", ts)

		log.Printf("Allocation profiling: capturing for %.2f seconds...\n", seconds)

		var duration = time.Duration(float64(seconds) * float64(time.Second))
		time.Sleep(duration)

		runtime.GC() // flush pending frees so the snapshot is accurate

		var f, err = os.Create(profileFile)
		if err != nil {
			log.Println("could not create allocs profile:", err)
			return
		}
		defer f.Close()

		if err := pprof.Lookup("allocs").WriteTo(f, 0); err != nil {
			log.Println("could not write allocs profile:", err)
			return
		}

		log.Println("Allocation profile saved at", profileFile)
		log.Println("Opening browser at http://localhost:8081 ...")

		exec.Command("go", "tool", "pprof", "-http=:8081", profileFile).Start()
	}()
}
func ProfileCPU(seconds float32) {
	go func() {
		// timestamp for filenames
		var ts = time.Now().Format("2006-01-02_15-04-05")
		var profileFile = fmt.Sprintf("cpu_%s.prof", ts)
		var svgFile = fmt.Sprintf("cpu_%s.svg", ts)

		var f, err = os.Create(profileFile)
		if err != nil {
			log.Println("could not create profile:", err)
			return
		}
		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			log.Println("could not start CPU profile:", err)
			return
		}
		log.Printf("CPU profiling started for %.2f seconds...\n", seconds)

		// convert float32 seconds → duration
		var duration = time.Duration(float64(seconds) * float64(time.Second))
		time.Sleep(duration)

		pprof.StopCPUProfile()
		log.Println("CPU profiling stopped. Profile saved at", profileFile)

		// Generate SVG via `go tool pprof`
		var cmd = exec.Command("go", "tool", "pprof", "-svg", profileFile)
		var out, err2 = cmd.CombinedOutput() // Captures both Stdout and Stderr
		if err2 != nil {
			log.Printf("failed to generate svg: %v. Output: %s", err2, string(out))
			return
		}
		if err2 := os.WriteFile(svgFile, out, 0644); err2 != nil {
			log.Println("failed to save svg:", err2)
			return
		}

		log.Println("SVG generated at", svgFile)

		exec.Command("xdg-open", svgFile).Start()
	}()
}
func ProfileTrace(seconds float32) {
	go func() {
		var ts = time.Now().Format("2006-01-02_15-04-05")
		var traceFile = fmt.Sprintf("trace_%s.out", ts)
		var f, err = os.Create(traceFile)
		if err != nil {
			log.Println("could not create trace file:", err)
			return
		}
		defer f.Close()

		// Start tracing
		if err := trace.Start(f); err != nil {
			log.Println("could not start trace:", err)
			return
		}

		log.Printf("Execution tracing: capturing for %.2f seconds...\n", seconds)

		// Capture for the specified duration
		var duration = time.Duration(float64(seconds) * float64(time.Second))
		time.Sleep(duration)

		trace.Stop()
		log.Println("Tracing stopped. Trace saved at", traceFile)

		log.Println("Opening trace viewer at http://localhost:8082 ...")

		// Launch the trace tool viewer
		// Note: trace tool requires its own server to parse the binary data
		exec.Command("go", "tool", "trace", "-http=:8082", traceFile).Start()
	}()
}

// private ========================================================

var memBuf []byte
var str = ""

func appendByteSize(buf []byte, n int) []byte {
	const unit = 1024
	if n < unit {
		buf = strconv.AppendInt(buf, int64(n), 10)
		return append(buf, " B"...)
	}
	var div, exp = int(unit), 0
	for v := n / unit; v >= unit; v /= unit {
		div *= unit
		exp++
	}
	buf = strconv.AppendFloat(buf, float64(n)/float64(div), 'f', 3, 64)
	buf = append(buf, ' ')
	buf = append(buf, "KMGTPE"[exp])
	return append(buf, 'B')
}
func appendThousands(buf []byte, n uint64) []byte {
	var tmp [32]byte
	var s = strconv.AppendUint(tmp[:0], n, 10)
	var length = len(s)
	for i, c := range s {
		if i > 0 && (length-i)%3 == 0 {
			buf = append(buf, ' ')
		}
		buf = append(buf, c)
	}
	return buf
}

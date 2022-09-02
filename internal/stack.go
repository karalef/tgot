package internal

import (
	"runtime"
	"strconv"
	"strings"
)

// BackTrace returns a formatted info about function invocations on the calling goroutine's stack.
func BackTrace(skip int) string {
	pc := make([]uintptr, 0, 16)
	for {
		n := runtime.Callers(2+skip+len(pc), pc[len(pc):cap(pc)])
		pc = pc[:len(pc)+n]
		if len(pc) < cap(pc) {
			pc = pc[:len(pc)-1]
			break
		}

		newpc := make([]uintptr, len(pc)*2)
		copy(newpc, pc)
		pc = newpc[:len(pc)]
	}

	var s strings.Builder
	frames := runtime.CallersFrames(pc)

	for {
		f, more := frames.Next()
		s.WriteString(f.Function + "\n\t" + f.File)
		s.WriteString(":" + strconv.Itoa(f.Line))
		if !more {
			break
		}
		s.WriteByte('\n')
	}

	return s.String()
}

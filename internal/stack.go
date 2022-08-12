package internal

import (
	"runtime"
	"strconv"
	"strings"
)

// BackTrace returns info about function invocations on the calling goroutine's stack.
func BackTrace(skip int, length int) []runtime.Frame {
	pc := make([]uintptr, length)
	n := runtime.Callers(2+skip, pc[:])

	stack := make([]runtime.Frame, 0, n)
	frames := runtime.CallersFrames(pc[:n])

	for {
		f, more := frames.Next()
		stack = append(stack, f)
		if !more {
			break
		}
	}

	return stack
}

// Caller returns runtime frame with info about func invocation.
func Caller(skip int) runtime.Frame {
	return BackTrace(2+skip, 1)[0]
}

// FramesString ...
func FramesString(frames []runtime.Frame) string {
	s := new(strings.Builder)
	for i, f := range frames {
		s.WriteString(f.Function + "\n\t" + f.File)
		s.WriteString(":" + strconv.Itoa(f.Line))
		if i < len(frames)-1 {
			s.WriteByte('\n')
		}
	}
	return s.String()
}

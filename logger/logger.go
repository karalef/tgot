package logger

import (
	"fmt"
	"io"
	"os"
)

// Discard is a Logger that does nothing on write calls.
var Discard = New(io.Discard, "")

// New creates new logger instance.
func New(w io.Writer, name string, useColor ...bool) *Logger {
	return &Logger{
		out: &writer{
			out:      w,
			useColor: len(useColor) > 0 && useColor[0],
		},
		name: name,
	}
}

// Default creates new logger instance with stderr output.
func Default(name string, useColor ...bool) *Logger {
	return New(os.Stderr, name, useColor...)
}

// File creates new logger instance with file output.
func File(path string, name string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	return New(f, name), nil
}

// Logger struct.
type Logger struct {
	out  *writer
	name string
}

// Child ...
func (l Logger) Child(name string) *Logger {
	l.name += "::" + name
	return &l
}

// SetOutput sets log output.
func (l *Logger) SetOutput(w io.Writer, useColor ...bool) {
	l.out = &writer{
		out:      w,
		useColor: len(useColor) > 0 && useColor[0],
	}
}

const namecol = green

func (l Logger) log(pref string, pcol ansiColor, f string, v ...interface{}) {
	if f == "" {
		return
	}

	w := l.out
	w.mut.Lock()
	w.writeBuf(pref, pcol)
	if l.name != "" {
		w.writeBuf(l.name, namecol)
	}
	w.writeBuf(fmt.Sprintf(f, v...), white)
	w.write()
	w.mut.Unlock()
}

// Info ...
func (l Logger) Info(f string, v ...interface{}) {
	l.log("INFO", white, f, v...)
}

// Warn ...
func (l Logger) Warn(f string, v ...interface{}) {
	l.log("WARN", yellow, f, v...)
}

// Error ...
func (l Logger) Error(f string, v ...interface{}) {
	l.log("ERROR", red, f, v...)
}

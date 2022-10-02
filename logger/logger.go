package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Discard is a Logger that does nothing on write calls.
var Discard = New(io.Discard, "")

// New creates new logger instance.
func New(w io.Writer, name string, color ...ColorConfig) *Logger {
	l := Logger{name: name}
	l.SetOutput(w, color...)
	return &l
}

// Default creates new logger instance with stderr output.
func Default(name string, color ...ColorConfig) *Logger {
	return New(os.Stderr, name, color...)
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
	name string
	out  *writer
}

// Child ...
func (l Logger) Child(name string) *Logger {
	l.name += "::" + name
	return &l
}

// SetOutput sets log output.
func (l *Logger) SetOutput(out io.Writer, color ...ColorConfig) {
	w := writer{out: out}
	if len(color) > 0 {
		w.color = color[0]
	}
	l.out = &w
}

func (l Logger) log(pref string, pcol *Color, f string, v ...any) {
	if f == "" {
		return
	}
	w := l.out
	w.mut.Lock()
	w.writeBuf(time.Now().Format("02.01.2006T15:04:05"), &w.color.Time)
	w.writeBuf(pref, pcol)
	w.writeBuf(l.name, &w.color.Name)
	w.writeBuf(fmt.Sprintf(f, v...), &w.color.Text)
	w.write()
	w.mut.Unlock()
}

// Info ...
func (l Logger) Info(f string, v ...interface{}) {
	l.log("INFO", &l.out.color.Info, f, v...)
}

// Warn ...
func (l Logger) Warn(f string, v ...interface{}) {
	l.log("WARN", &l.out.color.Warn, f, v...)
}

// Error ...
func (l Logger) Error(f string, v ...interface{}) {
	l.log("ERROR", &l.out.color.Error, f, v...)
}

type writer struct {
	out   io.Writer
	buf   []byte
	mut   sync.Mutex
	color ColorConfig
}

func (w *writer) writeBuf(text string, col *Color) {
	w.buf = append(w.buf, col.wrap(text)...)
	w.buf = append(w.buf, ' ')
}

func (w *writer) write() {
	w.buf = append(w.buf, '\n')
	w.out.Write(w.buf)
	w.buf = w.buf[:0]
}

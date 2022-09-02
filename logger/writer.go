package logger

import (
	"io"
	"sync"
	"time"
)

type writer struct {
	out      io.Writer
	buf      []byte
	mut      sync.Mutex
	last     time.Time
	useColor bool
}

func (w *writer) writeBuf(text string, col ansiColor) {
	if w.useColor {
		text = col.wrap(text)
	}
	w.buf = append(w.buf, text...)
	w.buf = append(w.buf, ' ')
}

const timecol = magenta

func (w *writer) write() {
	t := time.Now()
	if t.Minute() != w.last.Minute() {
		w.last = t
		f := t.Format("02.01.2006 15:04\n")
		if w.useColor {
			f = timecol.wrap(f)
		}
		w.out.Write([]byte(f))
	}

	w.buf = append(w.buf, '\n')
	w.out.Write(w.buf)
	w.buf = w.buf[:0]
}

type ansiColor string

func (c ansiColor) wrap(text string) string {
	return string(c) + text + string(resetColor)
}

const (
	red        ansiColor = "\033[1;31m"
	green      ansiColor = "\033[1;32m"
	yellow     ansiColor = "\033[1;33m"
	blue       ansiColor = "\033[1;34m"
	magenta    ansiColor = "\033[1;35m"
	cyan       ansiColor = "\033[1;36m"
	white      ansiColor = "\033[1;37m"
	resetColor ansiColor = "\033[0m"
)

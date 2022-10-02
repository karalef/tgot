package logger

// ColorConfig contains colors for name, time, text and prefixes.
type ColorConfig struct {
	Name, Time, Text  Color
	Info, Warn, Error Color
}

// DefaultColorConfig is default ColorConfig.
var DefaultColorConfig = ColorConfig{
	Name:  Green,
	Time:  Magenta,
	Text:  NoColor,
	Info:  White,
	Warn:  Yellow,
	Error: Red,
}

// Color represents color.
type Color uint8

// colors.
const (
	NoColor Color = iota
	White
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
)

func (c Color) wrap(text string) string {
	if ansi, ok := ansiTable[c]; ok {
		return ansi + text + reset
	}
	return text
}

var ansiTable = map[Color]string{
	White:   white,
	Red:     red,
	Green:   green,
	Yellow:  yellow,
	Blue:    blue,
	Magenta: magenta,
	Cyan:    cyan,
}

const (
	red     = "\033[1;31m"
	green   = "\033[1;32m"
	yellow  = "\033[1;33m"
	blue    = "\033[1;34m"
	magenta = "\033[1;35m"
	cyan    = "\033[1;36m"
	white   = "\033[1;37m"
	reset   = "\033[0m"
)

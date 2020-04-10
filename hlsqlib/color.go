package hlsqlib

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// Reset is the control code for resetting the color
	Reset = "\x1b[0m"
)

// Color codes
const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	LightGray
)

// Color codes
const (
	DarkGray Color = iota + 90
	LightRed
	LightGreen
	LightYellow
	LightBlue
	LightMagenta
	LightCyan
	White
)

var (
	colorNames = map[string]Color{
		"black":         Black,
		"red":           Red,
		"green":         Green,
		"yellow":        Yellow,
		"blue":          Blue,
		"magenta":       Magenta,
		"cyan":          Cyan,
		"light-gray":    LightGray,
		"dark-gray":     DarkGray,
		"light-red":     LightRed,
		"light-green":   LightGreen,
		"light-yellow":  LightYellow,
		"light-blue":    LightBlue,
		"light-magenta": LightMagenta,
		"light-cyan":    LightCyan,
		"white":         White,
	}
)

// Color wraps the int value for a color
type Color int

// ParseColor takes a string & converts to a Color
func ParseColor(s string) (Color, error) {
	s = strings.ToLower(s)
	if c, ok := colorNames[s]; ok {
		return c, nil
	}
	return Black, fmt.Errorf("No color found for %q", s)
}

// Control outputs the color's control code
func (c Color) Control() string { return "\x1b[0;" + strconv.Itoa(int(c)) + "m" }

// Colorizer aliases a control code
type Colorizer string

// NewColorizer creates a new colorizer from a color
func NewColorizer(c Color) Colorizer { return Colorizer(c.Control()) }

// S wraps a string with the control code
func (c Colorizer) S(s string) string { return string(c) + s + Reset }

// B wraps a byte slice with the control code
func (c Colorizer) B(b []byte) []byte {
	b = append([]byte(c), b...)
	b = append(b, Reset...)
	return b
}

// ColorSettings controls the colors that are used when serializing
type ColorSettings struct{ Tag, Attr Color }

package hlsqlib

import (
	"fmt"
	"strconv"
)

// Raw is an unadorned string
type Raw string

// Tag represents a # prefixed tag
type Tag struct {
	Name     string
	Attrs    []Attr
	Trailing []Raw
}

// AttrValue is a generic interface for attribute value
type AttrValue interface{ fmt.Stringer }

// AttrEnum is a string without quotes
type AttrEnum string

func (a AttrEnum) String() string { return string(a) }

// AttrString wraps a quoted string
type AttrString string

func (a AttrString) String() string { return `"` + string(a) + `"` }

// AttrBool YES or NO
type AttrBool bool

func (a AttrBool) String() string {
	if bool(a) {
		return "YES"
	}
	return "NO"
}

// AttrFloat attribe of a float
type AttrFloat float64

func (a AttrFloat) String() string { return fmt.Sprintf("%.3f", a) }

// AttrInt integer
type AttrInt int

func (a AttrInt) String() string { return strconv.Itoa(int(a)) }

// Attr Key+Value tuple
type Attr struct {
	Key   string
	Value AttrValue
}

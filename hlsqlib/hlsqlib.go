package hlsqlib

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Line interface for a line of the m3u8
type Line interface{}

// Raw is an unadorned string
type Raw string

// Tag represents a # prefixed tag
type Tag struct {
	Name  string
	Attrs []Attr
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

// ParseAttr take a key & value string to convert to a tuple
func ParseAttr(k, v string) Attr {
	var (
		attr = Attr{Key: k}
		i    = regexp.MustCompile(`^\d+$`)
		f    = regexp.MustCompile(`^\d+\.\d+$`)
	)
	switch {
	default:
		attr.Value = AttrEnum(v)
	case v == "":
		return attr
	case i.MatchString(v):
		i, _ := strconv.Atoi(v)
		attr.Value = AttrInt(i)
	case f.MatchString(v):
		f, _ := strconv.ParseFloat(v, 64)
		attr.Value = AttrFloat(f)
	case v == "YES" || v == "NO":
		attr.Value = AttrBool(v == "YES")
	case len(v) >= 2 && v[0] == '"' && v[len(v)-1] == '"':
		attr.Value = AttrString(v[1 : len(v)-1])
	}
	return attr
}

type runeBuffer struct{ runes []rune }

func (b runeBuffer) peek() rune {
	if len(b.runes) > 0 {
		return b.runes[0]
	}
	return '\n'
}

func (b *runeBuffer) takeUntil(t rune) string {
	var escaped bool
	for i, r := range b.runes {
		if r == t && !escaped {
			escaped = false
			out := b.runes[:i]
			b.runes = b.runes[i+1:]
			return string(out)
		}
		escaped = r == '\\'
	}
	out := string(b.runes)
	b.runes = nil
	return out
}
func (b *runeBuffer) pop() (rune, bool) {
	if len(b.runes) == 0 {
		return '0', false
	}
	out := b.runes[0]
	b.runes = b.runes[1:]
	return out, true
}

// ParseLine reads one line & parses it into the right structure
func ParseLine(s string) (Line, error) {
	if !strings.HasPrefix(s, "#") {
		return Raw(s), nil
	}
	parts := strings.SplitN(s, ":", 2)
	tag := Tag{Name: parts[0]}
	if len(parts) == 1 {
		return tag, nil
	}
	buf := runeBuffer{runes: []rune(parts[1])}
	for key := buf.takeUntil('='); key != ""; key = buf.takeUntil('=') {
		var v string
		if buf.peek() == '"' {
			buf.pop()
			v = "\"" + buf.takeUntil('"') + "\""
		}
		v += buf.takeUntil(',')
		tag.Attrs = append(tag.Attrs, ParseAttr(key, v))
	}
	return tag, nil
}

// SerializeOption is a function that is applied during serialization
type SerializeOption func(Line) (Line, bool)

// Chomp is an option that removes whitespace
var Chomp = SerializeOption(func(l Line) (Line, bool) {
	if raw, ok := l.(Raw); ok && strings.TrimSpace(string(raw)) == "" {
		return l, false
	}
	return l, true
})

// ColorLines returns a SerializeOption to color the output
func ColorLines(settings ColorSettings) SerializeOption {
	var (
		tag  = NewColorizer(settings.Tag)
		attr = NewColorizer(settings.Attr)
	)
	return SerializeOption(func(l Line) (Line, bool) {
		if t, ok := l.(Tag); ok {
			t.Name = tag.S(t.Name)
			for i, a := range t.Attrs {
				t.Attrs[i].Key = attr.S(a.Key)
			}
			l = t
		}
		return l, true
	})
}

// Serialize writes one line to the output
func Serialize(o io.Writer, l Line, opts ...SerializeOption) {
	var ok bool
	for _, opt := range opts {
		if l, ok = opt(l); !ok {
			return
		}
	}
	switch line := l.(type) {
	case Tag:
		out := line.Name
		if len(line.Attrs) > 0 {
			out += ":"
			for i, attr := range line.Attrs {
				if i > 0 {
					out += ","
				}
				out += attr.Key
				if attr.Value != nil {
					out += "=" + attr.Value.String()
				}
			}
		}
		fmt.Fprintln(o, out)
	case Raw:
		fmt.Fprintln(o, string(line))
	}
}

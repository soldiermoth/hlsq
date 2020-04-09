package hlsqlib

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type Line interface{}

type Raw string

type Tag struct {
	Name  string
	Attrs []Attr
}

type AttrValue interface{ fmt.Stringer }
type AttrEnum string

func (a AttrEnum) String() string { return string(a) }

type AttrString string

func (a AttrString) String() string { return `"` + string(a) + `"` }

type AttrBool bool

func (a AttrBool) String() string {
	if bool(a) {
		return "YES"
	}
	return "NO"
}

type AttrFloat float64

func (a AttrFloat) String() string { return fmt.Sprintf("%f", a) }

type AttrInt int

func (a AttrInt) String() string { return strconv.Itoa(int(a)) }

type Attr struct {
	Key   string
	Value AttrValue
}

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
		v := buf.takeUntil(',')
		tag.Attrs = append(tag.Attrs, ParseAttr(key, v))
	}
	return tag, nil
}

func Serialize(o io.Writer, lines ...Line) {
	for _, l := range lines {
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
}

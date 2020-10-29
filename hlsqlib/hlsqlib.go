package hlsqlib

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

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
	case strings.EqualFold(v, "YES") || strings.EqualFold(v, "NO"):
		attr.Value = AttrBool(strings.EqualFold(v, "YES"))
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

// Scanner reqds through an M3U8 file & parses it into Tags
type Scanner struct {
	last          bool
	err           error
	current, next Tag
	buf           *bufio.Scanner
}

// NewScanner creates a new scanner from a Reader
func NewScanner(r io.Reader) *Scanner {
	s := &Scanner{buf: bufio.NewScanner(r)}
	// Read once to get ball rolling
	s.Scan()
	return s
}

// Tag returns the current tag after a Scan()
func (s Scanner) Tag() Tag { return s.current }

// Scan advances the scanner to the next Tag
func (s *Scanner) Scan() bool {
	for s.buf.Scan() {
		// Get to next Tag
		line := s.buf.Text()
		if !strings.HasPrefix(line, "#") {
			s.next.Trailing = append(s.next.Trailing, Raw(line))
			continue
		}
		s.current = s.next
		s.next, s.err = ParseTag(line)
		return true
	}
	if !s.last {
		s.last = true
		s.current = s.next
		return true
	}
	return false
}

// ParseTag reads one line & parses it into the right structure
func ParseTag(s string) (Tag, error) {
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
type SerializeOption func(*Tag) bool

// Chomp is an option that removes whitespace
var Chomp = SerializeOption(func(t *Tag) bool {
	var newRaw []Raw
	for _, raw := range t.Trailing {
		if strings.TrimSpace(string(raw)) == "" {
			continue
		}
		newRaw = append(newRaw, raw)
	}
	t.Trailing = newRaw
	return true
})

// ColorLines returns a SerializeOption to color the output
func ColorLines(settings ColorSettings) SerializeOption {
	var (
		tag  = NewColorizer(settings.Tag)
		attr = NewColorizer(settings.Attr)
	)
	return SerializeOption(func(t *Tag) bool {
		t.Name = tag.S(t.Name)
		for i, a := range t.Attrs {
			t.Attrs[i].Key = attr.S(a.Key)
		}
		return true
	})
}

// AttrMatch is a matcher func that will filter tags that have an attribute that satisfies the function
func AttrMatch(f func(Attr) bool) SerializeOption {
	return SerializeOption(func(t *Tag) bool {
		for _, a := range t.Attrs {
			if f(a) {
				return true
			}
		}
		return false
	})
}

// Serialize writes one line to the output
func Serialize(o io.Writer, line Tag, opts ...SerializeOption) {
	var ok bool
	for _, opt := range opts {
		if ok = opt(&line); !ok {
			return
		}
	}
	out := line.Name
	if len(line.Attrs) > 0 {
		out += ":"
	}
	for i, attr := range line.Attrs {
		if i > 0 {
			out += ","
		}
		out += attr.Key
		if attr.Value != nil {
			out += "=" + attr.Value.String()
		}
	}
	fmt.Fprintln(o, out)
	for _, raw := range line.Trailing {
		fmt.Fprintln(o, string(raw))
	}
}

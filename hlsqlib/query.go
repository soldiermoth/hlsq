package hlsqlib

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	numOps = map[string]func(float64, AttrValue) bool{
		">":  func(t float64, a AttrValue) bool { f, ok := attrAsNum(a); return ok && t > f },
		">=": func(t float64, a AttrValue) bool { f, ok := attrAsNum(a); return ok && t >= f },
		"<":  func(t float64, a AttrValue) bool { f, ok := attrAsNum(a); return ok && t < f },
		"<=": func(t float64, a AttrValue) bool { f, ok := attrAsNum(a); return ok && t <= f },
		"=":  func(t float64, a AttrValue) bool { f, ok := attrAsNum(a); return ok && t == f },
		"!=": func(t float64, a AttrValue) bool { f, ok := attrAsNum(a); return ok && t != f },
	}
	strOps = map[string]func(string, AttrValue) bool{
		"=":  func(t string, a AttrValue) bool { s, ok := attrAsString(a); return ok && strings.EqualFold(s, t) },
		"!=": func(t string, a AttrValue) bool { s, ok := attrAsString(a); return ok && !strings.EqualFold(s, t) },
		"~": func(t string, a AttrValue) bool {
			t = strings.ToLower(t)
			s, ok := attrAsString(a)
			s = strings.ToLower(s)
			return ok && strings.Contains(s, t)
		},
		"!~": func(t string, a AttrValue) bool {
			t = strings.ToLower(t)
			s, ok := attrAsString(a)
			s = strings.ToLower(s)
			return ok && !strings.Contains(s, t)
		},
		"rlike": func(t string, a AttrValue) bool {
			regex := regexp.MustCompile(t)
			s, ok := attrAsString(a)
			return ok && regex.MatchString(s)
		},
	}
)

func attrAsNum(t AttrValue) (float64, bool) {
	switch n := t.(type) {
	case AttrInt:
		return float64(n), true
	case AttrFloat:
		return float64(n), true
	}
	return 0, false
}

func attrAsString(t AttrValue) (string, bool) {
	switch s := t.(type) {
	case AttrString:
		return string(s), true
	case AttrEnum:
		return string(s), true
	case AttrBool:
		return s.String(), true
	}
	return "", false
}

// Query is a function that returns true when matched against an attribute
type Query func(Attr) bool

// ParseQuery takes a string in form of {attr} {op} {value} & turns it into a Query func
func ParseQuery(q string) (Query, error) {
	pieces := strings.Split(q, " ")
	if len(pieces) != 3 {
		return nil, fmt.Errorf("expected '{name} {op} {value}' got %q", q)
	}
	var (
		match func(AttrValue) bool
		t     = ParseAttr(pieces[0], pieces[2])
		op    = pieces[1]
	)
	switch v := t.Value.(type) {
	case AttrInt:
		if opFunc, ok := numOps[op]; ok {
			match = func(a AttrValue) bool { return opFunc(float64(v), a) }
		}
	case AttrFloat:
		if opFunc, ok := numOps[op]; ok {
			match = func(a AttrValue) bool { return opFunc(float64(v), a) }
		}
	case AttrBool:
		if opFunc, ok := strOps[op]; ok {
			match = func(a AttrValue) bool { return opFunc(v.String(), a) }
		}
	case AttrEnum:
		if opFunc, ok := strOps[op]; ok {
			match = func(a AttrValue) bool { return opFunc(string(v), a) }
		}
	case AttrString:
		if opFunc, ok := strOps[op]; ok {
			match = func(a AttrValue) bool { return opFunc(string(v), a) }
		}
	}
	if match == nil {
		return nil, fmt.Errorf("no operator %q defined on type %T", op, t.Value)
	}
	return Query(func(a Attr) bool {
		return strings.EqualFold(a.Key, t.Key) && match(a.Value)
	}), nil
}

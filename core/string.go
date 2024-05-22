package core

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// A sequence of bytes, usually containing utf-8.
//
//lace:export
type String interface {
	any
	IndexCounted
	Seqable

	S() string
	AppendTo(str String) String
}

type GoString string

var _ String = GoString("")

func (s GoString) WithInfo(info *ObjectInfo) any {
	return s
}

func (s GoString) S() string {
	return string(s)
}

func (s GoString) ToString(env *Env, escape bool) (string, error) {
	printReadably := ToBool(env.printReadably.GetStatic())
	if printReadably || escape {
		return escapeString(s.S()), nil
	}
	return s.S(), nil
}

func MakeString(s string) String {
	if s == ":private" {
		panic("no")
	}
	return GoString(s)
}

func (s GoString) Equals(env *Env, other interface{}) bool {
	switch other := other.(type) {
	case String:
		return s == other
	default:
		return false
	}
}

func (s GoString) Native() interface{} {
	return s.S()
}

func (s GoString) Hash(env *Env) (uint32, error) {
	h := getHash()
	h.Write([]byte(s))
	return h.Sum32(), nil
}

func (s GoString) AppendTo(str String) String {
	return GoString(string(s) + str.S())
}

func (s GoString) Count() int {
	return utf8.RuneCountInString(s.S())
}

func (s GoString) Seq() Seq {
	runes := make([]any, 0, len(s))
	for _, r := range s {
		runes = append(runes, NewChar(r))
	}
	return &ArraySeq{arr: runes}
}

func (s GoString) Nth(env *Env, i int) (any, error) {
	if i < 0 {
		return nil, env.NewError(fmt.Sprintf("Negative index: %d", i))
	}
	j := 0
	var r rune

	for j, r = range s {
		if i == j {
			return NewChar(r), nil
		}
	}

	return nil, env.NewError(fmt.Sprintf("Index %d exceeds string's length %d", i, j+1))
}

func (s GoString) TryNth(env *Env, i int, d any) (any, error) {
	if i < 0 {
		return d, nil
	}
	for j, r := range s {
		if i == j {
			return NewChar(r), nil
		}
	}
	return d, nil
}

func (s GoString) Compare(env *Env, other any) (int, error) {
	s2, err := AssertString(env, other, "Cannot compare String and "+TypeName(other))
	if err != nil {
		return 0, err
	}

	return strings.Compare(s.S(), s2.S()), nil
}

func (s GoString) GetInfo() *ObjectInfo {
	return nil
}

type Rope struct {
	segments []string
}

var _ String = (*Rope)(nil)

func (s *Rope) S() string {
	return strings.Join(s.segments, "")
}

func (s *Rope) AppendTo(str String) String {
	r := *s

	switch o := str.(type) {
	case GoString:
		r.segments = append(s.segments, string(o))
	case *Rope:
		r.segments = append(s.segments, o.segments...)
	default:
		r.segments = append(s.segments, o.S())
	}

	return &r
}

func (s *Rope) ToString(env *Env, escape bool) (string, error) {
	printReadably := ToBool(env.printReadably.GetStatic())
	if printReadably || escape {
		return escapeString(s.S()), nil
	}
	return s.S(), nil
}

func (s *Rope) Equals(env *Env, other interface{}) bool {
	switch other := other.(type) {
	case String:
		return s == other
	default:
		return false
	}
}

func (s *Rope) Native() interface{} {
	return s.S()
}

func (s *Rope) Hash(env *Env) (uint32, error) {
	h := getHash()
	for _, p := range s.segments {
		h.Write([]byte(p))
	}
	return h.Sum32(), nil
}

func (s *Rope) Count() int {
	return utf8.RuneCountInString(s.S())
}

func (s *Rope) Seq() Seq {
	x := s.S()
	runes := make([]any, 0, len(x))
	for _, r := range x {
		runes = append(runes, NewChar(r))
	}
	return &ArraySeq{arr: runes}
}

func (s *Rope) Nth(env *Env, i int) (any, error) {
	if i < 0 {
		return nil, env.NewError(fmt.Sprintf("Negative index: %d", i))
	}
	j := 0
	var r rune

	for j, r = range s.S() {
		if i == j {
			return NewChar(r), nil
		}
	}

	return nil, env.NewError(fmt.Sprintf("Index %d exceeds string's length %d", i, j+1))
}

func (s *Rope) TryNth(env *Env, i int, d any) (any, error) {
	if i < 0 {
		return d, nil
	}
	for j, r := range s.S() {
		if i == j {
			return NewChar(r), nil
		}
	}
	return d, nil
}

func (s *Rope) Compare(env *Env, other any) (int, error) {
	s2, err := AssertString(env, other, "Cannot compare String and "+TypeName(other))
	if err != nil {
		return 0, err
	}

	return strings.Compare(s.S(), s2.S()), nil
}

var emptyString = GoString("")

// Combine many values into a single string.
//
//lace:export
func CombineToString(env *Env, args []any) (any, error) {
	if len(args) == 0 {
		return emptyString, nil
	}

	segments := make([]string, 0, len(args))

	for _, obj := range args {
		if obj != NIL {
			switch sv := obj.(type) {
			case GoString:
				segments = append(segments, sv.S())
			case *Rope:
				segments = append(segments, sv.segments...)
			case Char:
				segments = append(segments, string(sv.Ch()))
			case *Regex:
				segments = append(segments, sv.R.String())
			default:
				s, err := ToString(env, obj)
				if err != nil {
					return nil, err
				}
				segments = append(segments, s)
			}
		}
	}

	return &Rope{segments: segments}, nil
}

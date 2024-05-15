package core

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type String string

func (s String) S() string {
	return string(s)
}

func (s String) ToString(env *Env, escape bool) (string, error) {
	if escape {
		return escapeString(s.S()), nil
	}
	return s.S(), nil
}

func MakeString(s string) String {
	return String(s)
}

func (s String) Equals(env *Env, other interface{}) bool {
	switch other := other.(type) {
	case String:
		return s == other
	default:
		return false
	}
}

func (s String) GetType() *Type {
	return TYPE.String
}

func (s String) Native() interface{} {
	return s.S()
}

func (s String) Hash(env *Env) (uint32, error) {
	h := getHash()
	h.Write([]byte(s))
	return h.Sum32(), nil
}

func (s String) Count() int {
	return utf8.RuneCountInString(s.S())
}

func (s String) Seq() Seq {
	runes := make([]Object, 0, len(s))
	for _, r := range s {
		runes = append(runes, Char{Ch: r})
	}
	return &ArraySeq{arr: runes}
}

func (s String) Nth(env *Env, i int) (Object, error) {
	if i < 0 {
		return nil, env.NewError(fmt.Sprintf("Negative index: %d", i))
	}
	j := 0
	var r rune

	for j, r = range s {
		if i == j {
			return Char{Ch: r}, nil
		}
	}

	return nil, env.NewError(fmt.Sprintf("Index %d exceeds string's length %d", i, j+1))
}

func (s String) TryNth(env *Env, i int, d Object) (Object, error) {
	if i < 0 {
		return d, nil
	}
	for j, r := range s {
		if i == j {
			return Char{Ch: r}, nil
		}
	}
	return d, nil
}

func (s String) Compare(env *Env, other Object) (int, error) {
	os, err := other.GetType().ToString(env, false)
	if err != nil {
		return 0, err
	}
	s2, err := AssertString(env, other, "Cannot compare String and "+os)
	if err != nil {
		return 0, err
	}

	return strings.Compare(s.S(), s2.S()), nil
}

func (s String) GetInfo() *ObjectInfo {
	return nil
}

var emptyString = String("")

// Combine many values into a single string.
//
//lace:export
func CombineToString(env *Env, args []Object) (Object, error) {
	if len(args) == 0 {
		return emptyString, nil
	}

	var buffer strings.Builder
	for _, obj := range args {
		if !obj.Equals(env, NIL) {
			t := obj.GetType()
			// TODO: this is a hack. Rethink escape parameter in ToString
			escaped := (t == TYPE.String) || (t == TYPE.Char) || (t == TYPE.Regex)
			s, err := obj.ToString(env, !escaped)
			if err != nil {
				return nil, err
			}
			buffer.WriteString(s)
		}
	}
	return MakeString(buffer.String()), nil
}

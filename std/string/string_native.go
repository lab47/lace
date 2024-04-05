package string

import (
	"bytes"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	. "github.com/candid82/joker/core"
)

var newLine *regexp.Regexp

func padRight(s, pad string, n int) (string, error) {
	toAdd := n - utf8.RuneCountInString(s)
	if toAdd <= 0 {
		return s, nil
	}
	c := utf8.RuneCountInString(pad)
	d := toAdd / c
	r := toAdd % c
	for i := 0; i < d; i++ {
		s += pad
	}
	if r > 0 {
		s += string([]rune(pad)[:r])
	}
	return s, nil
}

func padLeft(s, pad string, n int) (string, error) {
	toAdd := n - utf8.RuneCountInString(s)
	if toAdd <= 0 {
		return s, nil
	}
	c := utf8.RuneCountInString(pad)
	d := toAdd / c
	r := toAdd % c
	for i := 0; i < d; i++ {
		s = pad + s
	}
	if r > 0 {
		s = string([]rune(pad)[c-r:]) + s
	}
	return s, nil
}

func split(s string, r *regexp.Regexp, n int) (Object, error) {
	indexes := r.FindAllStringIndex(s, n-1)
	lastStart := 0
	result := EmptyVector()
	var err error
	for _, el := range indexes {
		result, err = result.Conjoin(String{S: s[lastStart:el[0]]})
		if err != nil {
			return nil, err
		}
		lastStart = el[1]
	}
	result, err = result.Conjoin(String{S: s[lastStart:]})
	return result, err
}

func splitOnStringOrRegex(s string, sep Object, n int) (Object, error) {
	switch sep := sep.(type) {
	case String:
		v := strings.Split(s, sep.S)
		result := EmptyVector()
		var err error
		for _, el := range v {
			result, err = result.Conjoin(String{S: el})
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	case *Regex:
		return split(s, sep.R, n)
	default:
		panic(StubNewArgTypeError(1, sep, "String or Regex"))
	}
}

func join(sep string, seqable Seqable) (string, error) {
	seq := seqable.Seq()
	var b bytes.Buffer
	for !seq.IsEmpty() {
		b.WriteString(seq.First().ToString(false))
		seq = seq.Rest()
		if !seq.IsEmpty() {
			b.WriteString(sep)
		}
	}
	return b.String(), nil
}

func isBlank(env *Env, s Object) (bool, error) {
	if s.Equals(NIL) {
		return true, nil
	}
	str, err := AssertString(env, s, "")
	if err != nil {
		return false, err
	}
	for _, r := range str.S {
		if !unicode.IsSpace(r) {
			return false, nil
		}
	}
	return true, nil
}

func capitalize(s string) (string, error) {
	if len(s) < 2 {
		return strings.ToUpper(s), nil
	}
	return strings.ToUpper(string([]rune(s)[:1])) + strings.ToLower(string([]rune(s)[1:])), nil
}

func escape(env *Env, s string, cmap Callable) (string, error) {
	var b bytes.Buffer
	for _, r := range s {
		obj, err := cmap.Call(env, []Object{Char{Ch: r}})
		if err != nil {
			return "", err
		}
		if !obj.Equals(NIL) {
			b.WriteString(obj.ToString(false))
		} else {
			b.WriteRune(r)
		}
	}
	return b.String(), nil
}

func indexOf(s string, value Object, from int) (Object, error) {
	var res int
	if from != 0 {
		s = string([]rune(s)[from:])
	}
	switch value := value.(type) {
	case Char:
		res = strings.IndexRune(s, value.Ch)
	case String:
		res = strings.Index(s, value.S)
	default:
		return nil, StubNewArgTypeError(1, value, "String or Char")
	}
	if res == -1 {
		return NIL, nil
	}
	return MakeInt(utf8.RuneCountInString(s[:res]) + from), nil
}

func lastIndexOf(s string, value Object, from int) (Object, error) {
	var res int
	if from != 0 {
		s = string([]rune(s)[:from])
	}
	switch value := value.(type) {
	case Char:
		res = strings.LastIndex(s, string(value.Ch))
	case String:
		res = strings.LastIndex(s, value.S)
	default:
		return nil, StubNewArgTypeError(1, value, "String or Char")
	}
	if res == -1 {
		return NIL, nil
	}
	return MakeInt(utf8.RuneCountInString(s[:res])), nil
}

func replace(s string, match Object, repl string) (string, error) {
	switch match := match.(type) {
	case String:
		return strings.Replace(s, match.S, repl, -1), nil
	case *Regex:
		return match.R.ReplaceAllString(s, repl), nil
	default:
		return "", StubNewArgTypeError(1, match, "String or Regex")
	}
}

func replaceFirst(s string, match Object, repl string) (string, error) {
	switch match := match.(type) {
	case String:
		return strings.Replace(s, match.S, repl, 1), nil
	case *Regex:
		m := match.R.FindStringIndex(s)
		if m == nil {
			return s, nil
		}
		return s[:m[0]] + repl + s[m[1]:], nil
	default:
		return "", StubNewArgTypeError(1, match, "String or Regex")
	}
}

func reverse(s string) (string, error) {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes), nil
}

func init() {
	newLine, _ = regexp.Compile("\r?\n")
}

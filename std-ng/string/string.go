package string

import (
	"bytes"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/lab47/lace/core"
)

func Setup(env *core.Env) error {
	b := core.NewNSBuilder(env, "lace.string")

	b.Defn(&core.DefnInfo{
		Name:  "ends-with?",
		Args:  []string{"s", "substr"},
		Doc:   "True if s ends with substr.",
		Added: "1.0",
		Tag:   "Boolean",
		Fn:    strings.HasSuffix,
	})

	b.Defn(&core.DefnInfo{
		Name:  "starts-with?",
		Args:  []string{"s", "substr"},
		Doc:   "True if s starts with substr.",
		Added: "1.0",
		Tag:   "Boolean",
		Fn:    strings.HasPrefix,
	})

	b.Defn(&core.DefnInfo{
		Name:  "pad-right",
		Args:  []string{"s", "pad", "n"},
		Doc:   "Returns s padded with pad at the end to length n.",
		Added: "1.0",
		Tag:   "String",
		Fn:    padRight,
	})

	b.Defn(&core.DefnInfo{
		Name:  "pad-left",
		Args:  []string{"s", "pad", "n"},
		Doc:   "Returns s padded with pad at the beginning to length n.",
		Added: "1.0",
		Tag:   "String",
		Fn:    padLeft,
	})

	b.Defn(&core.DefnInfo{
		Name:  "split",
		Doc:   "Splits string on a string or regular expression. Returns vector of the splits.",
		Added: "1.0",
		Fns: []core.ArityFn{
			{
				Args: []string{"s", "sep"},
				Fn:   splitOnStringOrRegex2,
			},
			{
				Args: []string{"s", "sep", "n"},
				Fn:   splitOnStringOrRegex3,
			},
		},
	})

	newLine := regexp.MustCompile("\r?\n")

	b.Defn(&core.DefnInfo{
		Name:  "split-lines",
		Doc:   "Splits string on \\n or \\r\\n. Returns vector of the splits.",
		Added: "1.0",
		Args:  []string{"s"},
		Fn: func(s string) (core.Object, error) {
			return split(s, newLine, 0)
		},
	})

	b.Defn(&core.DefnInfo{
		Name:  "join",
		Doc:   "Returns a string of all elements in coll, as returned by (seq coll), separated by an optional separator.",
		Added: "1.0",
		Tag:   "String",
		Fns: []core.ArityFn{
			{
				Args: []string{"coll"},
				Fn: func(env *core.Env, seq core.Seqable) (string, error) {
					return join(env, "", seq)
				},
			},
			{
				Args: []string{"separator", "coll"},
				Fn:   join,
			},
		},
	})

	b.Defn(&core.DefnInfo{
		Name: "replace",
		Doc: `Replaces all instances of match (String or Regex) with string repl in string s.

  If match is Regex, $1, $2, etc. in the replacement string repl are
  substituted with the string that matched the corresponding
  parenthesized group in the pattern.
  `,
		Added: "1.0",
		Tag:   "String",
		Args:  []string{"s", "match", "repl"},
		Fn:    replace,
	})

	b.Defn(&core.DefnInfo{
		Name: "replace-first",
		Doc: `Replaces first instance of match (String or Regex) with string repl in string s.

  If match is Regex, $1, $2, etc. in the replacement string repl are
  substituted with the string that matched the corresponding
  parenthesized group in the pattern.
  `,
		Added: "1.0",
		Tag:   "String",
		Args:  []string{"s", "match", "repl"},
		Fn:    replaceFirst,
	})

	b.Defn(&core.DefnInfo{
		Name:  "trim",
		Doc:   "Removes whitespace from both ends of string.",
		Added: "1.0",
		Tag:   "String",
		Args:  []string{"s"},
		Fn:    strings.TrimSpace,
	})

	b.Defn(&core.DefnInfo{
		Name:  "trim-newline",
		Doc:   "Removes all trailing newline \\n or return \\r characters from string.",
		Added: "1.0",
		Tag:   "String",
		Args:  []string{"s"},
		Fn: func(s string) string {
			return strings.TrimRight(s, "\n\r")
		},
	})

	b.Defn(&core.DefnInfo{
		Name:    "trim-left",
		Doc:     "Removes whitespace from the left side of string.",
		Added:   "1.0",
		Tag:     "String",
		Args:    []string{"s"},
		Aliases: []string{"triml"},
		Fn: func(s string) string {
			return strings.TrimLeftFunc(s, unicode.IsSpace)
		},
	})

	b.Defn(&core.DefnInfo{
		Name:    "trim-right",
		Doc:     "Removes whitespace from the right side of string.",
		Added:   "1.0",
		Tag:     "String",
		Args:    []string{"s"},
		Aliases: []string{"trimr"},
		Fn: func(s string) string {
			return strings.TrimRightFunc(s, unicode.IsSpace)
		},
	})

	b.Defn(&core.DefnInfo{
		Name:  "blank?",
		Doc:   "True if s is nil, empty, or contains only whitespace.",
		Added: "1.0",
		Tag:   "Boolean",
		Args:  []string{"s"},
		Fn:    isBlank,
	})

	b.Defn(&core.DefnInfo{
		Name:  "capitalize",
		Doc:   "Converts first character of the string to upper-case, all other characters to lower-case.",
		Added: "1.0",
		Tag:   "String",
		Args:  []string{"s"},
		Fn:    capitalize,
	})

	b.Defn(&core.DefnInfo{
		Name: "escape",
		Doc: `Return a new string, using cmap to escape each character ch
  from s as follows:

  If (cmap ch) is nil, append ch to the new string.
  If (cmap ch) is non-nil, append (str (cmap ch)) instead.`,
		Added: "1.0",
		Tag:   "String",
		Args:  []string{"s", "cmap"},
		Fn:    escape,
	})

	b.Defn(&core.DefnInfo{
		Name:  "includes?",
		Doc:   "True if s includes substr.",
		Added: "1.0",
		Tag:   "Boolean",
		Args:  []string{"s", "substr"},
		Fn:    strings.Contains,
	})

	b.Defn(&core.DefnInfo{
		Name:  "index-of",
		Doc:   "Return index of value (string or char) in s, optionally searching forward from from or nil if not found.",
		Added: "1.0",
		Fns: []core.ArityFn{
			{
				Args: []string{"s", "value"},
				Fn: func(s string, val core.Object) (core.Object, error) {
					return indexOf(s, val, 0)
				},
			},
			{
				Args: []string{"s", "value", "num"},
				Fn:   indexOf,
			},
		},
	})

	b.Defn(&core.DefnInfo{
		Name:  "last-index-of",
		Doc:   "Return last index of value (string or char) in s, optionally searching forward from from or nil if not found.",
		Added: "1.0",
		Fns: []core.ArityFn{
			{
				Args: []string{"s", "value"},
				Fn: func(s string, val core.Object) (core.Object, error) {
					return lastIndexOf(s, val, 0)
				},
			},
			{
				Args: []string{"s", "value", "num"},
				Fn:   lastIndexOf,
			},
		},
	})

	b.Defn(&core.DefnInfo{
		Name:  "lower-case",
		Doc:   "Converts string to all lower-case.",
		Added: "1.0",
		Tag:   "String",
		Args:  []string{"s"},
		Fn:    strings.ToLower,
	})

	b.Defn(&core.DefnInfo{
		Name:  "upper-case",
		Doc:   "Converts string to all upper-case.",
		Added: "1.0",
		Tag:   "String",
		Args:  []string{"s"},
		Fn:    strings.ToUpper,
	})

	b.Defn(&core.DefnInfo{
		Name:  "reverse",
		Doc:   "Returns s with its characters reversed.",
		Added: "1.0",
		Tag:   "String",
		Args:  []string{"s"},
		Fn:    reverse,
	})

	b.Defn(&core.DefnInfo{
		Name:  "re-quote",
		Doc:   "Returns an instance of Regex that matches the string exactly",
		Added: "1.0",
		Tag:   "Regex",
		Args:  []string{"s"},
		Fn: func(s string) (*regexp.Regexp, error) {
			return regexp.Compile(regexp.QuoteMeta(s))
		},
	})

	return nil
}

func init() {
	core.AddNativeNamespace("lace.string", Setup)
}

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

func split(s string, r *regexp.Regexp, n int) (core.Object, error) {
	indexes := r.FindAllStringIndex(s, n-1)
	lastStart := 0
	result := core.EmptyVector()
	var err error
	for _, el := range indexes {
		result, err = result.Conjoin(core.MakeString(s[lastStart:el[0]]))
		if err != nil {
			return nil, err
		}
		lastStart = el[1]
	}
	result, err = result.Conjoin(core.MakeString(s[lastStart:]))
	return result, err
}

func splitOnStringOrRegex2(s string, sep core.Object) (core.Object, error) {
	return splitOnStringOrRegex3(s, sep, 0)
}

func splitOnStringOrRegex3(s string, sep core.Object, n int) (core.Object, error) {
	switch sep := sep.(type) {
	case core.String:
		v := strings.Split(s, sep.S())
		result := core.EmptyVector()
		var err error
		for _, el := range v {
			result, err = result.Conjoin(core.MakeString(el))
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	case *core.Regex:
		return split(s, sep.R, n)
	default:
		panic(core.StubNewArgTypeError(1, sep, "String or Regex"))
	}
}

func join(env *core.Env, sep string, seqable core.Seqable) (string, error) {
	seq := seqable.Seq()
	var b bytes.Buffer
	for {
		empty, err := seq.IsEmpty(env)
		if err != nil {
			return "", err
		}
		if empty {
			break
		}
		f, err := seq.First(env)
		if err != nil {
			return "", err
		}
		s, err := f.ToString(env, false)
		if err != nil {
			return "", err
		}
		b.WriteString(s)
		seq, err = seq.Rest(env)
		if err != nil {
			return "", err
		}
		empty, err = seq.IsEmpty(env)
		if err != nil {
			return "", err
		}
		if !empty {
			b.WriteString(sep)
		}
	}
	return b.String(), nil
}

func isBlank(env *core.Env, s core.Object) (bool, error) {
	if s.Equals(env, core.NIL) {
		return true, nil
	}
	str, err := core.AssertString(env, s, "")
	if err != nil {
		return false, err
	}
	for _, r := range str.S() {
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

func escape(env *core.Env, s string, cmap core.Callable) (string, error) {
	var b bytes.Buffer
	for _, r := range s {
		obj, err := cmap.Call(env, []core.Object{core.Char{Ch: r}})
		if err != nil {
			return "", err
		}
		if !obj.Equals(env, core.NIL) {
			s, err := obj.ToString(env, false)
			if err != nil {
				return "", err
			}
			b.WriteString(s)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String(), nil
}

func indexOf(s string, value core.Object, from int) (core.Object, error) {
	var res int
	if from != 0 {
		s = string([]rune(s)[from:])
	}
	switch value := value.(type) {
	case core.Char:
		res = strings.IndexRune(s, value.Ch)
	case core.String:
		res = strings.Index(s, value.S())
	default:
		return nil, core.StubNewArgTypeError(1, value, "String or Char")
	}
	if res == -1 {
		return core.NIL, nil
	}
	return core.MakeInt(utf8.RuneCountInString(s[:res]) + from), nil
}

func lastIndexOf(s string, value core.Object, from int) (core.Object, error) {
	var res int
	if from != 0 {
		s = string([]rune(s)[:from])
	}
	switch value := value.(type) {
	case core.Char:
		res = strings.LastIndex(s, string(value.Ch))
	case core.String:
		res = strings.LastIndex(s, value.S())
	default:
		return nil, core.StubNewArgTypeError(1, value, "String or Char")
	}
	if res == -1 {
		return core.NIL, nil
	}
	return core.MakeInt(utf8.RuneCountInString(s[:res])), nil
}

func replace(s string, match core.Object, repl string) (string, error) {
	switch match := match.(type) {
	case core.String:
		return strings.Replace(s, match.S(), repl, -1), nil
	case *core.Regex:
		return match.R.ReplaceAllString(s, repl), nil
	default:
		return "", core.StubNewArgTypeError(1, match, "String or Regex")
	}
}

func replaceFirst(s string, match core.Object, repl string) (string, error) {
	switch match := match.(type) {
	case core.String:
		return strings.Replace(s, match.S(), repl, 1), nil
	case *core.Regex:
		m := match.R.FindStringIndex(s)
		if m == nil {
			return s, nil
		}
		return s[:m[0]] + repl + s[m[1]:], nil
	default:
		return "", core.StubNewArgTypeError(1, match, "String or Regex")
	}
}

func reverse(s string) (string, error) {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes), nil
}

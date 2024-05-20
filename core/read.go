package core

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"regexp"
	"strconv"
	"sync/atomic"
	"unicode"
	"unicode/utf8"
)

type (
	ReadError struct {
		line     int
		column   int
		filename string
		msg      string
	}
	ReadFunc func(reader *Reader) any
	pos      struct {
		line   int
		column int
	}
)

const EOF = -1

var (
	LINTER_MODE   bool = false
	PROBLEM_COUNT      = 0
	DIALECT       Dialect
	LINTER_CONFIG *Var
)

func pushPos(reader *Reader) {
	reader.posStack = append(reader.posStack, pos{line: reader.line, column: reader.column})
}

func popPos(reader *Reader) pos {
	p := reader.posStack[len(reader.posStack)-1]
	reader.posStack = reader.posStack[:len(reader.posStack)-1]
	return p
}

func escapeRune(r rune) string {
	switch r {
	case ' ':
		return "\\space"
	case '\n':
		return "\\newline"
	case '\t':
		return "\\tab"
	case '\r':
		return "\\return"
	case '\b':
		return "\\backspace"
	case '\f':
		return "\\formfeed"
	default:
		return "\\" + string(r)
	}
}

func escapeString(str string) string {
	var b bytes.Buffer
	b.WriteRune('"')
	for _, r := range str {
		switch r {
		case '"':
			b.WriteString("\\\"")
		case '\\':
			b.WriteString("\\\\")
		case '\t':
			b.WriteString("\\t")
		case '\r':
			b.WriteString("\\r")
		case '\n':
			b.WriteString("\\n")
		case '\f':
			b.WriteString("\\f")
		case '\b':
			b.WriteString("\\b")
		default:
			b.WriteRune(r)
		}
	}
	b.WriteRune('"')
	return b.String()
}

func MakeReadError3(env *Env, reader *Reader, msg string, obj any) ReadError {
	if obj != nil {
		s, err := ToString(env, obj)
		if err != nil {
			s = fmt.Sprintf("%s(%p)", TypeName(obj), obj)
		}

		msg = msg + ": " + s
	}

	return ReadError{
		line:     reader.line,
		column:   reader.column,
		filename: reader.filename,
		msg:      msg,
	}
}

func MakeReadError(reader *Reader, msg string) ReadError {
	return ReadError{
		line:     reader.line,
		column:   reader.column,
		filename: reader.filename,
		msg:      msg,
	}
}

func MakeReadObject(reader *Reader, obj any) any {
	p := popPos(reader)
	return SetInfo(obj, &ObjectInfo{Position: Position{
		startColumn: p.column,
		startLine:   p.line,
		endLine:     reader.line,
		endColumn:   reader.column,
		filename:    reader.filename,
	}})
}

func DeriveReadObject(base any, obj any) any {
	baseInfo := GetInfo(base)
	if baseInfo != nil {
		bi := *baseInfo
		return SetInfo(obj, &bi)
	}
	return obj
}

func (err ReadError) Message() any {
	return MakeString(err.msg)
}

func (err ReadError) Error() string {
	return fmt.Sprintf("%s:%d:%d: Read error: %s", filename(err.filename), err.line, err.column, err.msg)
}

func isDelimiter(r rune) bool {
	switch r {
	case '(', ')', '[', ']', '{', '}', '"', ';', EOF, '\\':
		return true
	}
	return isWhitespace(r)
}

func eatString(reader *Reader, str string) error {
	for _, sr := range str {
		r, err := reader.Get()
		if err != nil {
			return err
		}

		if r != sr {
			return MakeReadError(reader, fmt.Sprintf("Unexpected character %U", r))
		}
	}
	return nil
}

func peekExpectedDelimiter(reader *Reader) error {
	r := reader.Peek()
	if !isDelimiter(r) {
		return MakeReadError(reader, "Character not followed by delimiter")
	}
	return nil
}

func readSpecialCharacter(reader *Reader, ending string, r rune) (any, error) {
	if err := eatString(reader, ending); err != nil {
		return nil, err
	}

	if err := peekExpectedDelimiter(reader); err != nil {
		return nil, err
	}

	return MakeReadObject(reader, NewChar(r)), nil
}

func isWhitespace(r rune) bool {
	return unicode.IsSpace(r) || r == ','
}

func eatWhitespace(env *Env, reader *Reader) error {
	r, err := reader.Get()
	if err != nil {
		return err
	}
	for r != EOF {
		if isWhitespace(r) {
			r, err = reader.Get()
			if err != nil {
				return err
			}
			continue
		}
		if r == ';' || (r == '#' && reader.Peek() == '!') {
			for r != '\n' && r != EOF {
				r, err = reader.Get()
				if err != nil {
					return err
				}
			}
			r, err = reader.Get()
			if err != nil {
				return err
			}
			continue
		}
		if r == '#' && reader.Peek() == '_' {
			_, err := reader.Get()
			if err != nil {
				return err
			}
			_, _, err = Read(env, reader)
			if err != nil {
				return err
			}

			r, err = reader.Get()
			if err != nil {
				return err
			}
			continue
		}
		reader.Unget()
		break
	}

	return nil
}

func readUnicodeCharacter(reader *Reader, length, base int) (any, error) {
	var b bytes.Buffer

	n, err := reader.Get()
	if err != nil {
		return nil, err
	}

	for !isDelimiter(n) {
		b.WriteRune(n)
		n, err = reader.Get()
		if err != nil {
			return nil, err
		}
	}

	reader.Unget()
	str := b.String()
	if len(str) != length {
		return nil, MakeReadError(reader, "Invalid unicode character: \\o"+str)
	}
	i, err := strconv.ParseInt(str, base, 32)
	if err != nil {
		return nil, MakeReadError(reader, "Invalid unicode character: \\o"+str)
	}
	err = peekExpectedDelimiter(reader)
	if err != nil {
		return nil, err
	}
	return MakeReadObject(reader, NewChar(rune(i))), nil
}

func readCharacter(reader *Reader) (any, error) {
	r, err := reader.Get()
	if err != nil {
		return nil, err
	}
	if r == EOF {
		return nil, MakeReadError(reader, "Incomplete character literal")
	}
	switch r {
	case 's':
		if reader.Peek() == 'p' {
			return readSpecialCharacter(reader, "pace", ' ')
		}
	case 'n':
		if reader.Peek() == 'e' {
			return readSpecialCharacter(reader, "ewline", '\n')
		}
	case 't':
		if reader.Peek() == 'a' {
			return readSpecialCharacter(reader, "ab", '\t')
		}
	case 'f':
		if reader.Peek() == 'o' {
			return readSpecialCharacter(reader, "ormfeed", '\f')
		}
	case 'b':
		if reader.Peek() == 'a' {
			return readSpecialCharacter(reader, "ackspace", '\b')
		}
	case 'r':
		if reader.Peek() == 'e' {
			return readSpecialCharacter(reader, "eturn", '\r')
		}
	case 'u':
		if !isDelimiter(reader.Peek()) {
			return readUnicodeCharacter(reader, 4, 16)
		}
	case 'o':
		if !isDelimiter(reader.Peek()) {
			return readUnicodeCharacter(reader, 3, 8)
		}
	}
	err = peekExpectedDelimiter(reader)
	if err != nil {
		return nil, err
	}
	return MakeReadObject(reader, NewChar(r)), nil
}

func scanBigInt(str string, base int, err error, reader *Reader) (any, error) {
	var bi big.Int
	if _, ok := bi.SetString(str, base); !ok {
		return nil, err
	}
	res := BigInt{b: bi}
	return MakeReadObject(reader, &res), nil
}

func scanRatio(str string, err error, reader *Reader) (any, error) {
	var rat big.Rat
	if _, ok := rat.SetString(str); !ok {
		return nil, err
	}
	r, err := ratioOrInt(&rat)
	if err != nil {
		return nil, err
	}
	return MakeReadObject(reader, r), nil
}

func scanBigFloat(str string, err error, reader *Reader) (any, error) {
	var bf big.Float
	if _, ok := bf.SetPrec(256).SetString(str); !ok {
		return nil, err
	}
	res := BigFloat{b: bf}
	return MakeReadObject(reader, &res), nil
}

func scanInt(str string, base int, err error, reader *Reader) (any, error) {
	i, e := strconv.ParseInt(str, base, 0)
	if e != nil {
		return scanBigInt(str, base, err, reader)
	}
	// TODO: 32-bit issue
	return MakeReadObject(reader, MakeInt(int(i))), nil
}

func readNumber(reader *Reader) (any, error) {
	var b bytes.Buffer
	isDouble, isHex, isExp, isRatio, base, nonDigits := false, false, false, false, "", 0
	d, err := reader.Get()
	if err != nil {
		return nil, err
	}
	last := d
	for !isDelimiter(d) {
		switch d {
		case '.':
			isDouble = true
		case '/':
			isRatio = true
		case 'x', 'X':
			isHex = true
		case 'e', 'E':
			isExp = true
		case 'r', 'R':
			if base == "" {
				base = b.String()
				b.Reset()
				last = d
				d, err = reader.Get()
				if err != nil {
					return nil, err
				}
				continue
			}
		}
		if !unicode.IsDigit(d) {
			nonDigits++
		}
		b.WriteRune(d)
		last = d
		d, err = reader.Get()
		if err != nil {
			return nil, err
		}
	}
	reader.Unget()
	str := b.String()
	if base != "" {
		invalidNumberError := MakeReadError(reader, fmt.Sprintf("Invalid number: %s", base+"r"+str))
		baseInt, err := strconv.ParseInt(base, 0, 0)
		if err != nil {
			return nil, invalidNumberError
		}
		if base[0] == '-' {
			baseInt = -baseInt
			str = "-" + str
		}
		if baseInt < 2 || baseInt > 36 {
			return nil, invalidNumberError
		}
		return scanInt(str, int(baseInt), invalidNumberError, reader)
	}
	invalidNumberError := MakeReadError(reader, fmt.Sprintf("Invalid number: %s", str))
	if isRatio {
		if nonDigits > 2 || nonDigits > 1 && str[0] != '-' && str[0] != '+' {
			return nil, invalidNumberError
		}
		return scanRatio(str, invalidNumberError, reader)
	}
	if last == 'N' {
		b.Truncate(b.Len() - 1)
		return scanBigInt(b.String(), 0, invalidNumberError, reader)
	}
	if last == 'M' {
		b.Truncate(b.Len() - 1)
		return scanBigFloat(b.String(), invalidNumberError, reader)
	}
	if isDouble || (!isHex && isExp) {
		dbl, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, invalidNumberError
		}
		return MakeReadObject(reader, Double{D: dbl}), nil
	}
	return scanInt(str, 0, invalidNumberError, reader)
}

func isSymbolInitial(r rune) bool {
	switch r {
	case '*', '+', '!', '-', '_', '?', ':', '=', '<', '>', '&', '.', '%', '$', '|':
		return true
	}
	return unicode.IsLetter(r) || r > 255
}

func isSymbolRune(r rune) bool {
	return isSymbolInitial(r) || unicode.IsDigit(r) || r == '#' || r == '/' || r == '\''
}

func readSymbol(env *Env, reader *Reader, first rune) (any, error) {
	var b bytes.Buffer
	if first != ':' {
		b.WriteRune(first)
	}
	var lastAdded rune
	r, err := reader.Get()
	if err != nil {
		return nil, err
	}
	for isSymbolRune(r) {
		if r == ':' {
			if lastAdded == ':' {
				return nil, MakeReadError(reader, "Invalid use of ':' in symbol name")
			}
		}
		b.WriteRune(r)
		lastAdded = r
		r, err = reader.Get()
		if err != nil {
			return nil, err
		}
	}
	if lastAdded == ':' || lastAdded == '/' {
		return nil, MakeReadError(reader, fmt.Sprintf("Invalid use of %c in symbol name", lastAdded))
	}
	reader.Unget()
	str := b.String()
	switch {
	case str == "":
		return nil, MakeReadError(reader, "Invalid keyword: :")
	case first == ':':
		if str[0] == '/' {
			return nil, MakeReadError(reader, "Blank namespaces are not allowed")
		}
		if str[0] == ':' {
			sym := MakeSymbol(str[1:])
			ns := env.NamespaceFor(env.CurrentNamespace(), sym)
			if ns == nil {
				msg := fmt.Sprintf("Unable to resolve namespace %s in keyword %s", sym.Namespace(), ":"+str)
				if LINTER_MODE {
					printReadWarning(reader, msg)
					return MakeReadObject(reader, MakeKeyword(sym.Name())), nil
				}
				return nil, MakeReadError(reader, msg)
			}
			ns.isUsed = true
			ns.isGloballyUsed = true
			return MakeReadObject(reader, MakeKeyword(ns.Name.Name()+"/"+sym.Name())), nil
		}
		return MakeReadObject(reader, MakeKeyword(str)), nil
	case str == "nil":
		return MakeReadObject(reader, NIL), nil
	case str == "true":
		return MakeReadObject(reader, Boolean(true)), nil
	case str == "false":
		return MakeReadObject(reader, Boolean(false)), nil
	default:
		return MakeReadObject(reader, MakeSymbol(str)), nil
	}
}

func readRegex(reader *Reader) (any, error) {
	var b bytes.Buffer
	r, err := reader.Get()
	if err != nil {
		return nil, err
	}
	for r != '"' {
		if r == EOF {
			return nil, MakeReadError(reader, "Non-terminated regex literal")
		}
		b.WriteRune(r)
		if r == '\\' {
			r, err = reader.Get()
			if err != nil {
				return nil, err
			}
			if r == EOF {
				return nil, MakeReadError(reader, "Non-terminated regex literal")
			}
			b.WriteRune(r)
		}
		r, err = reader.Get()
		if err != nil {
			return nil, err
		}
	}
	regex, err := regexp.Compile(b.String())
	if err != nil {
		if LINTER_MODE {
			return MakeReadObject(reader, &Regex{}), nil
		}
		return nil, MakeReadError(reader, "Invalid regex: "+err.Error())
	}
	return MakeReadObject(reader, &Regex{R: regex}), nil
}

func readUnicodeCharacterInString(reader *Reader, initial rune, length, base int, exactLength bool) (rune, error) {
	n := initial
	var b bytes.Buffer
	var err error
	for i := 0; i < length && n != '"'; i++ {
		b.WriteRune(n)
		n, err = reader.Get()
		if err != nil {
			return 0, err
		}
	}
	reader.Unget()
	str := b.String()
	if exactLength && len(str) != length {
		return 0, MakeReadError(reader, fmt.Sprintf("Invalid character length: %d, should be: %d", len(str), length))
	}
	i, err := strconv.ParseInt(str, base, 32)
	if err != nil {
		return 0, MakeReadError(reader, "Invalid unicode code: "+str)
	}
	return rune(i), nil
}

func readString(reader *Reader) (any, error) {
	var b bytes.Buffer
	r, err := reader.Get()
	if err != nil {
		return nil, err
	}
	for r != '"' {
		if r == '\\' {
			r, err = reader.Get()
			if err != nil {
				return nil, err
			}
			switch r {
			case '\\':
				r = '\\'
			case '"':
				r = '"'
			case 'n':
				r = '\n'
			case 't':
				r = '\t'
			case 'r':
				r = '\r'
			case 'b':
				r = '\b'
			case 'f':
				r = '\f'
			case 'u':
				n, err := reader.Get()
				if err != nil {
					return nil, err
				}
				r, err = readUnicodeCharacterInString(reader, n, 4, 16, true)
				if err != nil {
					return nil, err
				}
			default:
				if unicode.IsDigit(r) {
					r, err = readUnicodeCharacterInString(reader, r, 3, 8, false)
					if err != nil {
						return nil, err
					}
				} else {
					return nil, MakeReadError(reader, "Unsupported escape character: \\"+string(r))
				}
			}
		}
		if r == EOF {
			return nil, MakeReadError(reader, "Non-terminated string literal")
		}
		b.WriteRune(r)
		r, err = reader.Get()
		if err != nil {
			return nil, err
		}
	}
	return MakeReadObject(reader, MakeString(b.String())), nil
}

func readList(env *Env, reader *Reader) (any, error) {
	s := make([]any, 0, 10)
	err := eatWhitespace(env, reader)
	if err != nil {
		return nil, err
	}
	r := reader.Peek()
	for r != ')' {
		obj, multi, err := Read(env, reader)
		if err != nil {
			return nil, err
		}

		if multi {
			var v *Vector
			if err := Cast(env, obj, &v); err != nil {
				return nil, err
			}
			for i := 0; i < v.Count(); i++ {
				s = append(s, v.at(i))
			}
		} else {
			s = append(s, obj)
		}
		err = eatWhitespace(env, reader)
		if err != nil {
			return nil, err
		}
		r = reader.Peek()
	}
	_, err = reader.Get()
	if err != nil {
		return nil, err
	}
	list := EmptyList
	for i := len(s) - 1; i >= 0; i-- {
		list = list.conj(s[i])
	}
	res := MakeReadObject(reader, list)
	return res, nil
}

func readVector(env *Env, reader *Reader) (any, error) {
	res := EmptyVector()
	err := eatWhitespace(env, reader)
	if err != nil {
		return nil, err
	}
	r := reader.Peek()
	for r != ']' {
		obj, multi, err := Read(env, reader)
		if err != nil {
			return nil, err
		}

		if multi {
			var v *Vector
			if err := Cast(env, obj, &v); err != nil {
				return nil, err
			}
			for i := 0; i < v.Count(); i++ {
				res, err = res.Conjoin(v.at(i))
				if err != nil {
					return nil, err
				}
			}
		} else {
			res, err = res.Conjoin(obj)
			if err != nil {
				return nil, err
			}
		}
		err = eatWhitespace(env, reader)
		if err != nil {
			return nil, err
		}
		r = reader.Peek()
	}
	_, err = reader.Get()
	if err != nil {
		return nil, err
	}
	return MakeReadObject(reader, res), nil
}

func resolveKey(key any, nsname string) any {
	if nsname == "" {
		return key
	}
	switch key := key.(type) {
	case Keyword:
		if key.Namespace() == "" {
			return DeriveReadObject(key, MakeKeyword(nsname+"/"+key.Name()))
		}
		if key.Namespace() == "_" {
			return DeriveReadObject(key, MakeKeyword(key.Name()))
		}
	case Symbol:
		if key.Namespace() == "" {
			return DeriveReadObject(key, MakeSymbol(nsname+"/"+key.Name()))
		}
		if key.Namespace() == "_" {
			return DeriveReadObject(key, MakeSymbol(key.Name()))
		}
	}
	return key
}

func readMap(env *Env, reader *Reader) (any, error) {
	return readMapWithNamespace(env, reader, "")
}

func readMapWithNamespace(env *Env, reader *Reader, nsname string) (any, error) {
	err := eatWhitespace(env, reader)
	if err != nil {
		return nil, err
	}
	r := reader.Peek()
	objs := []any{}
	for r != '}' {
		obj, multi, err := Read(env, reader)
		if err != nil {
			return nil, err
		}
		if !multi {
			objs = append(objs, obj)
		} else {
			var v *Vector
			if err := Cast(env, obj, &v); err != nil {
				return nil, err
			}
			for i := 0; i < v.Count(); i++ {
				objs = append(objs, v.at(i))
			}
		}
		err = eatWhitespace(env, reader)
		if err != nil {
			return nil, err
		}
		r = reader.Peek()
	}
	_, err = reader.Get()
	if err != nil {
		return nil, err
	}
	if len(objs)%2 != 0 {
		return nil, MakeReadError(reader, "Map literal must contain an even number of forms")
	}
	if int64(len(objs)) >= HASHMAP_THRESHOLD {
		hashMap, err := NewHashMap(env)
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(objs); i += 2 {
			key := resolveKey(objs[i], nsname)
			if hashMap.containsKey(env, key) {
				return nil, MakeReadError3(env, reader, "Duplicate key", key)
			}
			v, err := hashMap.Assoc(env, key, objs[i+1])
			if err != nil {
				return nil, err
			}
			if err := Cast(env, v, &hashMap); err != nil {
				return nil, err
			}
		}
		return MakeReadObject(reader, hashMap), nil
	}
	m := EmptyArrayMap()
	for i := 0; i < len(objs); i += 2 {
		key := resolveKey(objs[i], nsname)
		if !m.Add(env, key, objs[i+1]) {
			return nil, MakeReadError3(env, reader, "Duplicate key", key)
		}
	}
	return MakeReadObject(reader, m), nil
}

func readSet(env *Env, reader *Reader) (any, error) {
	set := EmptySet()
	err := eatWhitespace(env, reader)
	if err != nil {
		return nil, err
	}
	r := reader.Peek()
	for r != '}' {
		obj, multi, err := Read(env, reader)
		if err != nil {
			return nil, err
		}
		if !multi {
			ok, err := set.Add(env, obj)
			if err != nil {
				return nil, err
			}
			if !ok {
				return nil, MakeReadError3(env, reader, "Duplicate set element ", obj)
			}
		} else {
			var v *Vector
			if err := Cast(env, obj, &v); err != nil {
				return nil, err
			}
			for i := 0; i < v.Count(); i++ {
				ok, err := set.Add(env, v.at(i))
				if err != nil {
					return nil, err
				}
				if !ok {
					return nil, MakeReadError3(env, reader, "Duplicate set element ", v.at(i))
				}
			}
		}
		err = eatWhitespace(env, reader)
		if err != nil {
			return nil, err
		}
		r = reader.Peek()
	}
	_, err = reader.Get()
	if err != nil {
		return nil, err
	}
	return MakeReadObject(reader, set), nil
}

func makeQuote(obj any, quote Symbol) any {
	res := NewListFrom(quote, obj)
	return DeriveReadObject(obj, res)
}

func readMeta(env *Env, reader *Reader) (*ArrayMap, error) {
	obj, err := readFirst(env, reader)
	if err != nil {
		return nil, err
	}
	switch v := obj.(type) {
	case *ArrayMap:
		return v, nil
	case String, Symbol:
		return &ArrayMap{arr: []any{DeriveReadObject(obj, criticalKeywords.tag), obj}}, nil
	case Keyword:
		return &ArrayMap{arr: []any{obj, DeriveReadObject(obj, Boolean(true))}}, nil
	default:
		return nil, MakeReadError(reader, "Metadata must be Symbol, Keyword, String or Map")
	}
}

func fillInMissingArgs(args map[int]Symbol) {
	max := 0
	for k := range args {
		if k > max {
			max = k
		}
	}
	for i := 1; i < max; i++ {
		if _, ok := args[i]; !ok {
			args[i] = generateSymbol("p__")
		}
	}
}

func makeFnForm(args map[int]Symbol, body any) (any, error) {
	fillInMissingArgs(args)

	a := make([]Symbol, len(args))
	for key, value := range args {
		if key != -1 {
			a[key-1] = value
		}
	}
	if v, ok := args[-1]; ok {
		a[len(args)-1] = criticalSymbols.amp
		a = append(a, v)
	}
	argVector := EmptyVector()
	var err error
	for _, v := range a {
		argVector, err = argVector.Conjoin(v)
		if err != nil {
			return nil, err
		}
	}
	return DeriveReadObject(body, NewListFrom(MakeSymbol("lace.core/fn"), argVector, body)), nil
}

func isTerminatingMacro(r rune) bool {
	switch r {
	case '"', ';', '@', '^', '`', '~', '(', ')', '[', ']', '{', '}', '\\':
		return true
	default:
		return false
	}
}

var genSymCounter = atomic.Int64{}

func genSym(prefix string, postfix string) Symbol {
	val := genSymCounter.Add(1)
	return MakeSymbol(fmt.Sprintf("%s%d%s", prefix, val, postfix))
}

func generateSymbol(prefix string) Symbol {
	return genSym(prefix, "#")
}

func registerArg(r *Reader, index int) Symbol {
	if s, ok := r.args[index]; ok {
		return s
	}
	r.args[index] = generateSymbol("p__")
	return r.args[index]
}

func readArgSymbol(env *Env, reader *Reader) (any, error) {
	r := reader.Peek()
	if isWhitespace(r) || isTerminatingMacro(r) {
		return MakeReadObject(reader, registerArg(reader, 1)), nil
	}
	obj, err := readFirst(env, reader)
	if err != nil {
		return nil, err
	}
	if Equals(env, obj, criticalSymbols.amp) {
		return MakeReadObject(reader, registerArg(reader, -1)), nil
	}
	switch n := obj.(type) {
	case Number:
		return MakeReadObject(reader, registerArg(reader, n.Int().I())), nil
	default:
		return nil, MakeReadError(reader, "Arg literal must be %, %& or %integer")
	}
}

func isSelfEvaluating(obj any) bool {
	if obj == EmptyList {
		return true
	}
	switch obj.(type) {
	case Boolean, Double, Int, Char, Keyword, String:
		return true
	default:
		return false
	}
}

func isCall(env *Env, obj any, name Symbol) (bool, error) {
	switch seq := obj.(type) {
	case Seq:
		f, err := seq.First(env)
		if err != nil {
			return false, err
		}
		return name.Is(f), nil
	default:
		return false, nil
	}
}

func syntaxQuoteSeq(tenv *Env, seq Seq, env map[string]Symbol, reader *Reader) (*ArraySeq, error) {
	res := make([]any, 0)
	for iter := iter(seq); iter.HasNext(tenv); {
		obj, err := iter.Next(tenv)
		if err != nil {
			return nil, err
		}
		ok, err := isCall(tenv, obj, criticalSymbols.unquoteSplicing)
		if err != nil {
			return nil, err
		}
		if ok {
			var seq Seq
			if err := Cast(tenv, obj, &seq); err != nil {
				return nil, err
			}
			r, err := seq.Rest(tenv)
			if err != nil {
				return nil, err
			}

			f, err := r.First(tenv)
			if err != nil {
				return nil, err
			}
			res = append(res, f)
		} else {
			q, err := makeSyntaxQuote(tenv, obj, env, reader)
			if err != nil {
				return nil, err
			}
			res = append(res, DeriveReadObject(q, NewListFrom(criticalSymbols.list, q)))
		}
	}
	return &ArraySeq{arr: res}, nil
}

func syntaxQuoteColl(tenv *Env, seq Seq, env map[string]Symbol, reader *Reader, ctor Symbol, info *ObjectInfo) (any, error) {
	q, err := syntaxQuoteSeq(tenv, seq, env, reader)
	if err != nil {
		return nil, err
	}

	concat := q.Cons(criticalSymbols.concatS)
	seqList := NewListFrom(criticalSymbols.seq, concat)
	var res any = seqList
	if ctor != criticalSymbols.emptySymbol {
		res = NewListFrom(ctor, seqList).Cons(criticalSymbols.apply)
	}
	return SetInfo(res, info), nil
}

func makeSyntaxQuote(tenv *Env, obj any, env map[string]Symbol, reader *Reader) (any, error) {
	if isSelfEvaluating(obj) {
		return obj, nil
	}
	if IsSpecialSymbol(obj) {
		return makeQuote(obj, criticalSymbols.quote), nil
	}
	info := GetInfo(obj)
	switch s := obj.(type) {
	case Symbol:
		str := s.Name()
		if r, _ := utf8.DecodeLastRuneInString(str); r == '#' && s.Namespace() == "" {
			sym, ok := env[s.Name()]
			if !ok {
				sym = generateSymbol(str[:len(str)-1] + "__")
				env[s.Name()] = sym
			}
			obj = DeriveReadObject(obj, sym)
		} else {
			v, err := tenv.ResolveSymbol(s)
			if err != nil {
				return nil, err
			}
			obj = DeriveReadObject(obj, v)
		}
		return makeQuote(obj, criticalSymbols.quote), nil
	case Seq:
		ok, err := isCall(tenv, obj, criticalSymbols.unquote)
		if err != nil {
			return nil, err
		}

		if ok {
			return Second(tenv, s)
		}

		ok, err = isCall(tenv, obj, criticalSymbols.unquoteSplicing)
		if err != nil {
			return nil, err
		}
		if ok {
			return nil, MakeReadError(reader, "Splice not in list")
		}
		return syntaxQuoteColl(tenv, s, env, reader, criticalSymbols.emptySymbol, info)
	case *Vector:
		return syntaxQuoteColl(tenv, s.Seq(), env, reader, criticalSymbols.vector, info)
	case *ArrayMap:
		return syntaxQuoteColl(tenv, ArraySeqFromArrayMap(s), env, reader, criticalSymbols.hashMap, info)
	case *MapSet:
		return syntaxQuoteColl(tenv, s.Seq(), env, reader, criticalSymbols.hashSet, info)
	default:
		return obj, nil
	}
}

func filename(f string) string {
	if f != "" {
		return f
	}
	return "<file>"
}

func handleNoReaderError(env *Env, reader *Reader, s Symbol) (any, error) {
	if LINTER_MODE {
		if DIALECT != EDN {
			printReadWarning(reader, "No reader function for tag "+s.String())
		}
		return readFirst(env, reader)
	}
	return nil, MakeReadError3(env, reader, "No reader function for tag", s)
}

func readTagged(env *Env, reader *Reader) (any, error) {
	obj, err := readFirst(env, reader)
	if err != nil {
		return nil, err
	}

	switch s := obj.(type) {
	case Symbol:
		readersVar, ok := env.CoreNamespace.LookupVar(criticalSymbols.defaultDataReaders.Name())
		if !ok {
			return handleNoReaderError(env, reader, s)
		}
		readersMap, ok := readersVar.GetStatic().(Map)
		if !ok {
			return handleNoReaderError(env, reader, s)
		}
		ok, readFunc := readersMap.GetEqu(s)
		if !ok {
			return handleNoReaderError(env, reader, s)
		}
		v, err := AssertVar(env, readFunc, "")
		if err != nil {
			return nil, err
		}
		o, err := readFirst(env, reader)
		if err != nil {
			return nil, err
		}
		return v.Call(env, []any{o})
	default:
		return nil, MakeReadError(reader, "Reader tag must be a symbol")
	}
}

func readConditional(env *Env, reader *Reader) (any, bool, error) {
	isSplicing := false
	if reader.Peek() == '@' {
		_, err := reader.Get()
		if err != nil {
			return nil, false, err
		}
		isSplicing = true
	}
	err := eatWhitespace(env, reader)
	if err != nil {
		return nil, false, err
	}
	r, err := reader.Get()
	if err != nil {
		return nil, false, err
	}
	if r != '(' {
		return nil, false, MakeReadError(reader, "Reader conditional body must be a list")
	}
	v, err := readList(env, reader)
	if err != nil {
		return nil, false, err
	}
	var cond *List
	if err := Cast(env, v, &cond); err != nil {
		return nil, false, err
	}
	if cond.count%2 != 0 {
		if LINTER_MODE {
			printReadError(reader, "Reader conditional requires an even number of forms")
		} else {
			return nil, false, MakeReadError(reader, "Reader conditional requires an even number of forms")
		}
	}
	for cond.count > 0 {
		ok, _, err := env.Features.Get(env, cond.first)
		if err != nil {
			return nil, false, err
		}
		if ok {
			v, err := Second(env, cond)
			if err != nil {
				return nil, false, err
			}
			if isSplicing {
				s, ok := v.(Seqable)
				if !ok {
					msg := "Spliced form in reader conditional must be Seqable, got " + TypeName(v)
					if LINTER_MODE {
						printReadError(reader, msg)
						return EmptyVector(), true, nil
					} else {
						return nil, false, MakeReadError(reader, msg)
					}
				}
				vec, err := NewVectorFromSeq(env, s.Seq())
				if err != nil {
					return nil, false, err
				}
				return vec, true, nil
			}
			return v, false, nil
		}
		cond = cond.rest.rest
	}
	return EmptyVector(), true, nil
}

func readNamespacedMap(env *Env, reader *Reader) (any, error) {
	r, err := reader.Get()
	if err != nil {
		return nil, err
	}
	auto := r == ':'
	if !auto {
		reader.Unget()
	}
	var sym any
	r, err = reader.Get()
	if err != nil {
		return nil, err
	}
	if isWhitespace(r) {
		if !auto {
			reader.Unget()
			return nil, MakeReadError(reader, "Namespaced map must specify a namespace")
		}
		for isWhitespace(r) {
			r, err = reader.Get()
			if err != nil {
				return nil, err
			}
		}
		if r != '{' {
			reader.Unget()
			return nil, MakeReadError(reader, "Namespaced map must specify a namespace")
		}
	} else if r != '{' {
		reader.Unget()
		sym, _, err = Read(env, reader)
		if err != nil {
			return nil, err
		}
		r, err = reader.Get()
		if err != nil {
			return nil, err
		}
		for isWhitespace(r) {
			r, err = reader.Get()
			if err != nil {
				return nil, err
			}
		}
	}
	if r != '{' {
		return nil, MakeReadError(reader, "Namespaced map must specify a map")
	}
	var nsname string
	if auto {
		if sym == nil {
			nsname = env.CurrentNamespace().Name.Name()
		} else {
			sym, ok := sym.(Symbol)
			if !ok || sym.Namespace() != "" {
				return nil, MakeReadError3(env, reader, "Namespaced map must specify a valid namespace", sym)
			}
			ns := env.CurrentNamespace().aliases[sym.Name()]
			if ns == nil {
				ns = env.FindNamespace(sym)
			}
			if ns == nil {
				return nil, MakeReadError3(env, reader, "Unknown auto-resolved namespace alias", sym)
			}
			ns.isUsed = true
			ns.isGloballyUsed = true
			nsname = ns.Name.Name()
		}
	} else {
		if sym == nil {
			return nil, MakeReadError(reader, "Namespaced map must specify a valid namespace")
		}
		sym, ok := sym.(Symbol)
		if !ok || sym.Namespace() != "" {
			return nil, MakeReadError3(env, reader, "Namespaced map must specify a valid namespace", sym)
		}
		nsname = sym.Name()
	}
	return readMapWithNamespace(env, reader, nsname)
}

func readDispatch(env *Env, reader *Reader) (any, bool, error) {
	r, err := reader.Get()
	if err != nil {
		return nil, false, err
	}
	switch r {
	case '"':
		re, err := readRegex(reader)
		if err != nil {
			return nil, false, err
		}
		return re, false, nil
	case '\'':
		popPos(reader)
		nextObj, err := readFirst(env, reader)
		if err != nil {
			return nil, false, err
		}
		return DeriveReadObject(nextObj, NewListFrom(DeriveReadObject(nextObj, criticalSymbols._var), nextObj)), false, nil
	case '^':
		popPos(reader)
		v, err := readWithMeta(env, reader)
		if err != nil {
			return nil, false, err
		}
		return v, false, nil
	case '{':
		s, err := readSet(env, reader)
		if err != nil {
			return nil, false, err
		}
		return s, false, nil
	case '(':
		popPos(reader)
		reader.Unget()
		old := reader.args
		reader.args = make(map[int]Symbol)
		fn, err := readFirst(env, reader)
		if err != nil {
			return nil, false, err
		}
		res, err := makeFnForm(reader.args, fn)
		if err != nil {
			return nil, false, err
		}
		reader.args = old
		return res, false, nil
	case '?':
		return readConditional(env, reader)
	case ':':
		m, err := readNamespacedMap(env, reader)
		if err != nil {
			return nil, false, err
		}

		return m, false, nil
	}
	popPos(reader)
	reader.Unget()
	v, err := readTagged(env, reader)
	if err != nil {
		return nil, false, err
	}
	return v, false, nil
}

func readWithMeta(env *Env, reader *Reader) (any, error) {
	meta, err := readMeta(env, reader)
	if err != nil {
		return nil, err
	}
	nextObj, err := readFirst(env, reader)
	if err != nil {
		return nil, err
	}
	switch v := nextObj.(type) {
	case Meta:
		m, err := v.WithMeta(env, meta)
		if err != nil {
			return nil, err
		}
		return DeriveReadObject(nextObj, m), nil
	default:
		return nil, MakeReadError3(env, reader, "Metadata cannot be applied to", v)
	}
}

func readFirst(env *Env, reader *Reader) (any, error) {
	obj, multi, err := Read(env, reader)
	if err != nil {
		return nil, err
	}
	if !multi {
		return obj, nil
	}
	var v *Vector
	if err := Cast(env, obj, &v); err != nil {
		return nil, err
	}
	if v.Count() == 0 {
		return readFirst(env, reader)
	}
	return v.at(0), nil
}

func Read(env *Env, reader *Reader) (any, bool, error) {
	err := eatWhitespace(env, reader)
	if err != nil {
		return nil, false, err
	}
	r, err := reader.Get()
	if err != nil {
		return nil, false, err
	}
	pushPos(reader)
	switch {
	case r == '\\':
		c, err := readCharacter(reader)
		if err != nil {
			return nil, false, err
		}
		return c, false, nil
	case unicode.IsDigit(r):
		reader.Unget()
		o, err := readNumber(reader)
		if err != nil {
			return nil, false, err
		}
		return o, false, nil
	case r == '-' || r == '+':
		if unicode.IsDigit(reader.Peek()) {
			reader.Unget()
			o, err := readNumber(reader)
			if err != nil {
				return nil, false, err
			}
			return o, false, nil
		}
		o, err := readSymbol(env, reader, r)
		if err != nil {
			return nil, false, err
		}

		return o, false, nil
	case r == '%' && reader.args != nil:
		v, err := readArgSymbol(env, reader)
		if err != nil {
			return nil, false, err
		}
		return v, false, nil
	case isSymbolInitial(r):
		v, err := readSymbol(env, reader, r)
		if err != nil {
			return nil, false, err
		}
		return v, false, nil
	case r == '"':
		v, err := readString(reader)
		if err != nil {
			return nil, false, err
		}
		return v, false, nil
	case r == '(':
		v, err := readList(env, reader)
		if err != nil {
			return nil, false, err
		}
		return v, false, nil
	case r == '[':
		v, err := readVector(env, reader)
		if err != nil {
			return nil, false, err
		}
		return v, false, nil
	case r == '{':
		v, err := readMap(env, reader)
		if err != nil {
			return nil, false, err
		}
		return v, false, nil
	case r == '/' && isDelimiter(reader.Peek()):
		return MakeReadObject(reader, criticalSymbols.backslash), false, nil
	case r == '\'':
		popPos(reader)
		nextObj, err := readFirst(env, reader)
		if err != nil {
			return nil, false, err
		}
		return makeQuote(nextObj, criticalSymbols.quote), false, nil
	case r == '@':
		popPos(reader)
		nextObj, err := readFirst(env, reader)
		if err != nil {
			return nil, false, err
		}
		return DeriveReadObject(nextObj, NewListFrom(DeriveReadObject(nextObj, criticalSymbols.deref), nextObj)), false, nil
	case r == '~':
		popPos(reader)
		if reader.Peek() == '@' {
			_, err := reader.Get()
			if err != nil {
				return nil, false, err
			}
			nextObj, err := readFirst(env, reader)
			if err != nil {
				return nil, false, err
			}
			return makeQuote(nextObj, criticalSymbols.unquoteSplicing), false, nil
		}
		nextObj, err := readFirst(env, reader)
		if err != nil {
			return nil, false, err
		}
		return makeQuote(nextObj, criticalSymbols.unquote), false, nil
	case r == '`':
		popPos(reader)
		nextObj, err := readFirst(env, reader)
		if err != nil {
			return nil, false, err
		}
		sq, err := makeSyntaxQuote(env, nextObj, make(map[string]Symbol), reader)
		if err != nil {
			return nil, false, err
		}
		return sq, false, nil
	case r == '^':
		popPos(reader)
		m, err := readWithMeta(env, reader)
		if err != nil {
			return nil, false, err
		}

		return m, false, nil
	case r == '#':
		return readDispatch(env, reader)
	case r == EOF:
		return nil, false, MakeReadError(reader, "Unexpected end of file")
	}
	return nil, false, MakeReadError(reader, fmt.Sprintf("Unexpected %c", r))
}

func TryRead(env *Env, reader *Reader) (obj any, err error) {
	for {
		err := eatWhitespace(env, reader)
		if err != nil {
			return nil, err
		}
		if reader.Peek() == EOF {
			return NIL, io.EOF
		}
		obj, multi, err := Read(env, reader)
		if err != nil {
			return nil, err
		}
		if !multi {
			return obj, nil
		}
		var v *Vector
		if err := Cast(env, obj, &v); err != nil {
			return nil, err
		}
		if v.Count() > 0 {
			PROBLEM_COUNT++
			return NIL, MakeReadError(reader, "Reader conditional splicing not allowed at the top level.")
		}
	}
}

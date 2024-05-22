package core

import (
	"cmp"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/exp/slices"
	"golang.org/x/term"
)

type ExInfo struct {
	ArrayMap

	err error
}

var (
	_ any   = (*ExInfo)(nil)
	_ error = (*ExInfo)(nil)
)

func (exInfo *ExInfo) ToString(env *Env, escape bool) (string, error) {
	return exInfo.Error(), nil
}

func (exInfo *ExInfo) Equals(env *Env, other interface{}) bool {
	return exInfo == other
}

func (exInfo *ExInfo) Hash(env *Env) (uint32, error) {
	return HashPtr(exInfo), nil
}

func (exInfo *ExInfo) Message() any {
	if ok, res := exInfo.GetEqu(criticalKeywords.message); ok {
		return res
	}
	return NIL
}

func (e *ExInfo) Unwrap() error {
	return e.err
}

func (exInfo *ExInfo) Error() string {
	var pos Position
	prefix := "Exception"

	_, data := exInfo.GetEqu(criticalKeywords.data)
	dm, ok := data.(Map)
	if ok {
		ok, form := dm.GetEqu(criticalKeywords.form)
		if ok {
			pos = GetPosition(form)
		}
		if ok, pr := dm.GetEqu(criticalKeywords._prefix); ok {
			prefix = SimpleToString(pr)
		}
	}
	_, msg := exInfo.GetEqu(criticalKeywords.message)
	var strMsg string

	if sv, ok := msg.(String); ok {
		strMsg = sv.S()
	} else {
		strMsg = "no proper message"
	}

	return fmt.Sprintf("%s:%d:%d: %s: %s", pos.Filename(), pos.startLine, pos.startColumn, prefix, strMsg)
}

// The standard error thrown when evalution of a program encounters an error.
//
//lace:export
type EvalError struct {
	err  error
	hash uint32

	cat string

	stackTrace *VMStacktrace

	Map
}

func (e *EvalError) Category() string {
	return cmp.Or(e.cat, "Error")
}

func (err *EvalError) Unwrap() error {
	return err.err
}

func (err *EvalError) Is(target error) bool {
	_, ok := target.(*EvalError)
	return ok
}

func (err *EvalError) ToString(env *Env, escape bool) (string, error) {
	return err.Error(), nil
}

func (err *EvalError) Hash(env *Env) (uint32, error) {
	return err.hash, nil
}

func (err *EvalError) WithInfo(info *ObjectInfo) any {
	return err
}

func (err *EvalError) Message() any {
	return MakeString(err.err.Error())
}

func (err *EvalError) Error() string {
	return err.err.Error()
}

func (err *EvalError) ErrorData() Map {
	return err.Map
}

var ErrCustomError = errors.New("custom error")

type ErrorData interface {
	ErrorData() Map
}

type HasCategory interface {
	Category() string
}

func Errorf(env *Env, str string, args ...any) error {
	return env.populateStackTrace(fmt.Errorf(str, args...))
}

func SError(env *Env, cat, str string, args ...any) error {
	if len(args)%2 != 0 {
		args = append(args, NIL)
	}

	var bits []any

	for i := 0; i < len(args); i += 2 {
		var ko any
		switch sv := args[i].(type) {
		case string:
			ko = MakeString(sv)
		case any:
			ko = sv
		default:
			continue
		}

		var vo any

		switch sv := args[i+1].(type) {
		case string:
			vo = MakeString(sv)
		case int:
			vo = MakeInt(sv)
		case any:
			vo = sv
		default:
			vo = MakeString(fmt.Sprint(sv))
		}

		bits = append(bits, ko, vo)
	}

	d, _ := NewHashMap(env, bits...)

	return env.populateStackTrace(
		&EvalError{
			err: errors.New(str),
			cat: cat,
			Map: d,
		},
	)
}

func hashStringU32(str string) uint32 {
	h := fnv.New32()
	h.Write([]byte(str))

	return h.Sum32()
}

func WrapError(env *Env, err error) *EvalError {
	if err == nil {
		return nil
	}

	if !errors.Is(err, &EvalError{}) {
		err = &EvalError{
			err:  err,
			hash: hashStringU32(err.Error()),
		}
	}

	return env.populateStackTrace(err)
}

func AddContext(env *Env, err error, str string, args ...any) error {
	e := WrapError(env, err)

	obj := MakeString(fmt.Sprintf(str, args...))

	if e.Map == nil {
		m, err := NewArrayMap(MakeSymbol("context"), obj)
		if err == nil {
			e.Map = m.(Map)
		}
		return e
	}

	m, err := e.Map.Assoc(env, MakeSymbol("context"), obj)
	if err == nil {
		e.Map = m.(Map)
	}
	return e
}

func NewEvalError(env *Env, str string) *EvalError {
	err := &EvalError{
		err:  fmt.Errorf(str),
		hash: hashStringU32(str),
	}

	return env.populateStackTrace(err)
}

func (e *EvalError) AddData(env *Env, obj any) {
	if e.Map == nil {
		m, err := NewArrayMap(MakeKeyword("data"), obj)
		if err == nil {
			e.Map = m.(Map)
		}
		return
	}

	m, err := e.Map.Assoc(env, MakeKeyword("data"), obj)
	if err == nil {
		e.Map = m.(Map)
	}
}

func DisplayError(env *Env, err error) {
	var ee *EvalError

	if !errors.As(err, &ee) {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return
	}

	fmt.Fprintf(os.Stderr, "An unexpected error occurred.\n\nStacktrace:\n")
	if ee.stackTrace != nil {
		ee.stackTrace.PrintTo(env, os.Stderr)
		fmt.Fprintln(os.Stderr)
	}
	category := ee.Category()

	fmt.Fprintf(Stderr, "%s: %s\n", category, err.Error())

	if ed, ok := err.(ErrorData); ok {
		data := ed.ErrorData()

		if data != nil {
			var keys []string

			vals := make(map[string]string)

			m := data.Iter()
			for m.HasNext() {
				p := m.Next()

				key, _ := ToString(env, p.Key)
				val, _ := ToString(env, p.Value)

				keys = append(keys, key)
				vals[key] = val
			}

			sort.Strings(keys)

			for _, k := range keys {
				val := vals[k]
				fmt.Fprintf(os.Stderr, "        %s: %s\n", keyColor(k), val)
			}
		}
	}

	fmt.Fprintln(os.Stderr)
}

type VMStacktrace struct {
	upper      error
	StackTrace any
	pcs        []uintptr
	treeStack  []Expr
}

func (v *VMStacktrace) Unwrap() error {
	return v.upper
}

func (v *VMStacktrace) Error() string {
	return v.upper.Error()
}

func (v *VMStacktrace) Is(other error) bool {
	_, ok := other.(*VMStacktrace)
	return ok
}

func (env *Env) populateStackTrace(err error) *EvalError {
	var ee *EvalError
	if !errors.As(err, &ee) {
		ee = &EvalError{
			err:  err,
			hash: hashStringU32(err.Error()),
		}
	}

	if ee.stackTrace != nil {
		return ee
	}

	pcs := make([]uintptr, 256)
	cnt := runtime.Callers(2, pcs)

	var ts []Expr
	if len(env.treeEvalStack) > 0 {
		ts = slices.Clone(env.treeEvalStack)
	}

	if env.Engine == nil {
		return ee
	}

	ee.stackTrace = &VMStacktrace{
		upper:      err,
		StackTrace: env.Engine.makeStackTrace(),
		pcs:        pcs[:cnt],
		treeStack:  ts,
	}

	return ee
}

type outputFrame struct {
	name string
	loc  string
	lace bool
}

func (vs *VMStacktrace) renderFrame(env *Env, ele any) outputFrame {
	var str string

	switch sv := ele.(type) {
	case String:
		return outputFrame{name: sv.S()}
	case IndexCounted:
		if sv.Count() >= 2 {
			a, _ := sv.Nth(env, 0)
			b, _ := sv.Nth(env, 1)

			var (
				fn *Fn
				ip Int
			)

			if cmp.Or(
				Cast(env, a, &fn),
				Cast(env, b, &ip),
			) == nil {
				var name string
				if fn.meta != nil {
					if ok, val := fn.meta.GetEqu(criticalKeywords.name); ok {
						if sym, ok := val.(Symbol); ok {
							name = sym.String()
						}
					}

					if ok, val := fn.meta.GetEqu(criticalKeywords.ns); ok {
						if ns, ok := val.(Symbol); ok {
							name = ns.Name() + "/" + name
						}
					}
				}

				codeFile := fn.code.fileForIp(ip.I())
				if codeFile != fn.code.filename {
					macroLine := fn.code.macroLineForIp(ip.I())

					return outputFrame{
						lace: true,
						name: name,
						loc:  fmt.Sprintf("%s:%d (from %s:%d)", codeFile, macroLine, fn.code.filename, fn.code.lineForIp(ip.I())),
					}
				} else {
					return outputFrame{
						lace: true,
						name: name,
						loc:  fmt.Sprintf("%s:%d", fn.code.filename, fn.code.lineForIp(ip.I())),
					}
				}
			}
		}
	}

	var err error
	if str == "" {
		str, err = ToString(env, ele)
		if err != nil {
			str = fmt.Sprintf("error decoding stacktrace: %s\n", err)
		}
	}
	return outputFrame{
		name: str,
	}
}

const bcName = "github.com/lab47/lace/core.(*Engine).RunBC"

func extractMod(path string) string {
	idx := strings.Index(path, "pkg/mod/")
	modOn := path[idx+len("pkg/mod/"):]

	at := strings.IndexByte(modOn, '@')

	if at != -1 {
		pkg := modOn[:at]
		rest := modOn[at:]

		slash := strings.IndexByte(rest, '/')
		if slash != -1 {
			return "<" + pkg + ">" + rest[slash:]
		}
	}

	return path
}

var laceDir string

func init() {
	_, fileName, _, _ := runtime.Caller(0)

	laceDir = filepath.Dir(filepath.Dir(fileName))
}

func splitName(name string) (string, string) {
	i := len(name) - 1
	for ; i > 0; i-- {
		if name[i] == '/' {
			break
		}
	}
	for ; i < len(name); i++ {
		if name[i] == '.' {
			break
		}
	}
	return name[:i], name[i:]
}

func trimName(fn *runtime.Func) string {
	if fn == nil {
		return ""
	}

	pkg, name := splitName(fn.Name())

	pkg = strings.ReplaceAll(pkg, "github.com/lab47/lace", "lace")

	return pkg + name
}

var (
	goColor   = color.New(color.FgBlue).Sprintf
	laceColor = color.New(color.FgHiWhite).Sprintf
	locColor  = color.New(color.FgWhite).Sprintf
	sepColor  = color.New(color.Faint).Sprintf
	keyColor  = color.New(color.Bold).Sprintf
)

var ignoreFuncs = map[string]struct{}{
	"github.com/lab47/lace/core.WrapToProc3_2[...].func1": {},
	"runtime.goexit": {},
	"runtime.main":   {},
}

func cleanupPath(path string) string {
	clean := extractMod(path)

	if strings.HasPrefix(clean, laceDir) {
		clean = "<lace>" + clean[len(laceDir):]
	}

	return clean
}

func (vs *VMStacktrace) PrintTo(env *Env, w io.Writer) {
	frames := runtime.CallersFrames(vs.pcs)

	if st, ok := vs.StackTrace.(Seq); ok {
		it := iter(st)

		var oframes []outputFrame

		for {
			fr, more := frames.Next()

			var ofr outputFrame

			if fr.Func.Name() == bcName {
				ele, err := it.Next(env)
				if err != nil {
					ofr = outputFrame{name: fmt.Sprintf("error decoding stackframe: %s", err)}
				} else {
					ofr = vs.renderFrame(env, ele)
				}
				oframes = append(oframes, ofr)

			} else {
				if _, skip := ignoreFuncs[fr.Func.Name()]; !skip {
					ofr = outputFrame{
						name: trimName(fr.Func),
						loc:  fmt.Sprintf("%s:%d", cleanupPath(fr.File), fr.Line),
					}
					oframes = append(oframes, ofr)

				}
			}

			if !more {
				break
			}
		}

		width := 0

		for _, ofr := range oframes {
			if len(ofr.name) > width {
				width = len(ofr.name)
			}
		}

		var maxWidth int

		if f, ok := w.(*os.File); ok {
			maxWidth, _, _ = term.GetSize(int(f.Fd()))
		}

		pad := strings.Repeat(" ", width)

		if len(vs.treeStack) > 0 {
			fmt.Fprintf(w, "%s  Macro evalution trace:\n", pad)
			var prev string
			for _, e := range vs.treeStack {
				cur := e.Pos().String()
				if cur == prev {
					continue
				}

				prev = cur
				fmt.Fprintf(w, "%s  %s %s\n", pad, sepColor("@"), locColor(cur))
			}
			fmt.Fprintf(w, "%s  -----------------------\n", pad)
		}

		for _, ofr := range oframes {
			padWidth := len(pad) - len(ofr.name)
			visSize := len(ofr.name) + len(ofr.loc) + 2 + padWidth

			if visSize >= maxWidth {
				padWidth = maxWidth - visSize
			}

			cw := goColor
			if ofr.lace {
				cw = laceColor
			}

			str := fmt.Sprintf("%s %s %s", cw(ofr.name), sepColor("@"), locColor(ofr.loc))

			if padWidth <= 0 {
				fmt.Fprintf(w, " %s\n", str)
			} else {
				fmt.Fprintf(w, "%s %s\n", pad[:padWidth], str)
			}
		}

		/*
			for it.HasNext(env) {
				ele, err := it.Next(env)
				if err != nil {
					fmt.Fprintf(w, "error decoding stacktrace: %s\n", err)
				}

				str := vs.renderFrame(env, ele)

				fmt.Fprintln(w, str)
			}
		*/

		return
	}

	str, err := ToString(env, vs.StackTrace)
	if err == nil {
		fmt.Fprintln(w, str)
	}
}

func StubNewError(msg string) *EvalError {
	err := errors.New(msg)

	h := fnv.New32()
	h.Write([]byte(err.Error()))

	return &EvalError{
		err:  err,
		hash: h.Sum32(),
	}
}

func StubNewArgTypeError(index int, obj any, expectedType string) *EvalError {
	return StubNewError(fmt.Sprintf("Arg[%d] of <<func_name>> must have type %s, got %s", index, expectedType, TypeName(obj)))
}

func (e *Env) NewError(msg string, args ...any) *EvalError {
	return WrapError(e, fmt.Errorf(msg, args...))
}

func TypeError[T any](env *Env, obj any) *EvalError {
	ts := reflect.TypeFor[T]().String()
	ee := env.NewError(fmt.Sprintf("object must have type %s, got %s", ts, TypeName(obj)))
	return env.populateStackTrace(ee)
}

type TCContext struct {
	Context string
	Index   int
}

func (e *Env) NewArgTypeError(index int, obj any, expectedType string) *EvalError {
	if index >= 0 {
		return e.NewError(fmt.Sprintf("Arg[%d] must have type %s, got %s", index, expectedType, TypeName(obj)))
	} else {
		return e.NewError(fmt.Sprintf("Value must have type %s, got %s", expectedType, TypeName(obj)))
	}
}

func (e *Env) TypeError(ctx TCContext, obj any, expectedType string) *EvalError {
	if ctx.Context != "" {
		if ctx.Index >= 0 {
			return e.NewError(fmt.Sprintf("%s[%d] must have type %s, got %s", ctx.Context, ctx.Index, expectedType, TypeName(obj)))
		} else {
			return e.NewError(fmt.Sprintf("%s must have type %s, got %s", ctx.Context, expectedType, TypeName(obj)))
		}
	} else {
		return e.NewError(fmt.Sprintf("Value must have type %s, got %s", expectedType, TypeName(obj)))
	}
}

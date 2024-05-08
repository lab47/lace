package core

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/exp/slices"
	"golang.org/x/term"
)

var (
	Stdin          io.Reader = os.Stdin
	Stdout         io.Writer = os.Stdout
	Stderr         io.Writer = os.Stderr
	VerbosityLevel           = 0
)

type (
	Env struct {
		Namespaces    map[string]*Namespace
		CoreNamespace *Namespace
		stdout        *Var
		stdin         *Var
		stderr        *Var
		printReadably *Var
		file          *Var
		MainFile      *Var
		args          *Var
		classPath     *Var
		ns            *Var
		NS_VAR        *Var
		IN_NS_VAR     *Var
		version       *Var
		Features      Set
		CurrentVar    Associative

		Context        context.Context
		cycleDetection map[[2]Object]struct{}

		RT *Runtime

		Engine *Engine

		DebugBytecode bool

		treeEvalStack []Expr
	}
)

func (env *Env) enableCycleDetection() func() {
	if env.cycleDetection != nil {
		return func() {}
	}

	env.cycleDetection = make(map[[2]Object]struct{})
	return func() {
		env.cycleDetection = nil
	}
}

func (env *Env) cycling(a, b Object) bool {
	if env.cycleDetection == nil {
		return false
	}

	key := [2]Object{a, b}

	_, ok := env.cycleDetection[key]
	if ok {
		return true
	}

	env.cycleDetection[key] = struct{}{}

	return false
}

func versionMap(env *Env) Map {
	res := EmptyArrayMap()
	parts := strings.Split(VERSION[1:], ".")
	i, _ := strconv.ParseInt(parts[0], 10, 64)
	res.Add(env, MakeKeyword("major"), Int{I: int(i)})
	i, _ = strconv.ParseInt(parts[1], 10, 64)
	res.Add(env, MakeKeyword("minor"), Int{I: int(i)})
	i, _ = strconv.ParseInt(parts[2], 10, 64)
	res.Add(env, MakeKeyword("incremental"), Int{I: int(i)})
	return res
}

func (env *Env) SetEnvArgs(newArgs []string) {
	args := EmptyVector()
	for _, arg := range newArgs {
		args, _ = args.Conjoin(MakeString(arg))
	}
	if args.Count() > 0 {
		env.args.Value = args.Seq()
	} else {
		env.args.Value = NIL
	}
}

/*
This runs after invariant initialization, which includes calling

	NewEnv().
*/
func (env *Env) SetClassPath(cp string) {
	cpArray := filepath.SplitList(cp)
	cpVec := EmptyVector()
	for _, cpelem := range cpArray {
		cpVec, _ = cpVec.Conjoin(MakeString(cpelem))
	}
	if cpVec.Count() == 0 {
		cpVec, _ = cpVec.Conjoin(MakeString(""))
	}
	env.classPath.Value = cpVec
}

/*
This runs after invariant initialization, which includes calling

	NewEnv().
*/
func (env *Env) InitEnv(stdin io.Reader, stdout, stderr io.Writer, args []string) {
	env.stdin.Value = MakeBufferedReader(stdin)
	env.stdout.Value = MakeIOWriter(stdout)
	env.stderr.Value = MakeIOWriter(stderr)
	env.SetEnvArgs(args)
}

func (env *Env) SetStdIO(stdin, stdout, stderr Object) {
	env.stdin.Value = stdin
	env.stdout.Value = stdout
	env.stderr.Value = stderr
}

func (env *Env) StdIO() (stdin, stdout, stderr Object) {
	return env.stdin.Value, env.stdout.Value, env.stderr.Value
}

/*
This runs after invariant initialization, which includes calling

	NewEnv().
*/
func (env *Env) SetMainFilename(filename string) {
	env.MainFile.Value = MakeString(filename)
}

/*
This runs after invariant initialization, which includes calling

	NewEnv().
*/
func (env *Env) SetFilename(obj Object) {
	env.file.Value = obj
}

func (env *Env) IsStdIn(obj Object) bool {
	return env.stdin.Value == obj
}

func (env *Env) CurrentNamespace() *Namespace {
	ns, err := AssertNamespace(env, env.ns.Value, "")
	if err != nil {
		panic(err) // this is extremely rare, we should probably make it not possible
	}

	return ns
}

func (env *Env) SetCurrentNamespace(ns *Namespace) {
	env.ns.Value = ns
}

func (env *Env) EnsureNamespace(sym Symbol) *Namespace {
	if sym.ns != "" {
		panic(env.RT.NewError("Namespace's name cannot be qualified: " + sym.String()))
	}
	var err error
	if env.Namespaces[sym.name] == nil {
		env.Namespaces[sym.name], err = NewNamespace(env, sym)
		if err != nil {
			panic(err)
		}
		if setup, ok := builtinNSSetup[sym.name]; ok {
			err := setup(env)
			if err != nil {
				panic(err)
			}
		} else {
			_, err = PopulateNativeNamespaceToEnv(env, sym.name)
			if err != nil {
				panic(err)
			}
		}
	}
	return env.Namespaces[sym.name]
}

func (env *Env) ensureNamespace(sym Symbol) *Namespace {
	if sym.ns != "" {
		panic(env.RT.NewError("Namespace's name cannot be qualified: " + sym.String()))
	}
	var err error
	if env.Namespaces[sym.name] == nil {
		env.Namespaces[sym.name], err = NewNamespace(env, sym)
		if err != nil {
			panic(err)
		}
	}
	return env.Namespaces[sym.name]
}

func (env *Env) NamespaceFor(ns *Namespace, s Symbol) *Namespace {
	var res *Namespace
	if s.ns == "" {
		res = ns
	} else {
		res = ns.aliases[s.ns]
		if res == nil {
			res = env.Namespaces[s.ns]
		}
	}
	if res != nil {
		res.MaybeLazy(env, "NamespaceFor")
	}
	return res
}

func (env *Env) ResolveIn(n *Namespace, s Symbol) (*Var, bool) {
	ns := env.NamespaceFor(n, s)
	if ns == nil {
		return nil, false
	}
	if v, ok := ns.mappings[s.name]; ok {
		return v, true
	}
	if s.Is(env.IN_NS_VAR.name) {
		return env.IN_NS_VAR, true
	}
	if s.Is(env.NS_VAR.name) {
		return env.NS_VAR, true
	}
	return nil, false
}

func (env *Env) Resolve(s Symbol) (*Var, bool) {
	ns := env.CurrentNamespace()
	return env.ResolveIn(ns, s)
}

func (env *Env) MakeVar(s Symbol) (*Var, error) {
	ns := env.CurrentNamespace()
	return ns.Intern(env, s)
}

func (env *Env) FindNamespace(s Symbol) *Namespace {
	if s.ns != "" {
		return nil
	}
	ns := env.Namespaces[s.name]
	if ns != nil {
		ns.MaybeLazy(env, "FindNameSpace")
	} else {
		if _, ok := builtinNSSetup[s.name]; ok {
			// don't call setup! just create the namespace because EnsureNamespace will call
			// the setup.
			ns = env.EnsureNamespace(s)
		} else {
			_, err := PopulateNativeNamespaceToEnv(env, s.name)
			if err != nil {
				panic(err)
			}
		}
	}

	//if ns == nil {
	//panic("nope " + *s.name)
	//}
	return ns
}

func (env *Env) RemoveNamespace(s Symbol) *Namespace {
	if s.ns != "" {
		return nil
	}
	if s.Is(criticalSymbols.lace_core) {
		panic(env.RT.NewError("Cannot remove core namespace"))
	}
	ns := env.Namespaces[s.name]
	delete(env.Namespaces, s.name)
	return ns
}

func (env *Env) ResolveSymbol(s Symbol) (Symbol, error) {
	if strings.ContainsRune(s.name, '.') {
		return s, nil
	}
	if s.ns == "" && TYPES[s.name] != nil {
		return s, nil
	}
	currentNs := env.CurrentNamespace()

	if s.ns != "" {
		ns := env.NamespaceFor(currentNs, s)
		if ns == nil || ns.Name.name == s.ns {
			if ns != nil {
				ns.isUsed = true
				ns.isGloballyUsed = true
			}
			return s, nil
		}
		ns.isUsed = true
		ns.isGloballyUsed = true
		return Symbol{
			name: s.name,
			ns:   ns.Name.name,
		}, nil
	}
	vr, ok := currentNs.mappings[s.name]
	if !ok {
		return Symbol{
			name: s.name,
			ns:   currentNs.Name.name,
		}, nil
	}
	vr.isUsed = true
	vr.isGloballyUsed = true
	vr.ns.isUsed = true
	vr.ns.isGloballyUsed = true
	return Symbol{
		name: vr.name.name,
		ns:   vr.ns.Name.name,
	}, nil
}

func (env *Env) Eval(str string) (Object, error) {
	reader := NewReader(strings.NewReader(str), "<expr>")
	return ProcessReader(env, reader, "")
}

type VMStacktrace struct {
	upper      error
	StackTrace Object
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

func (env *Env) populateStackTrace(err error) error {
	if errors.Is(err, &VMStacktrace{}) {
		return err
	}

	pcs := make([]uintptr, 256)
	cnt := runtime.Callers(2, pcs)

	var ts []Expr
	if len(env.treeEvalStack) > 0 {
		ts = slices.Clone(env.treeEvalStack)
	}

	return &VMStacktrace{
		upper:      err,
		StackTrace: env.Engine.makeStackTrace(),
		pcs:        pcs[:cnt],
		treeStack:  ts,
	}
}

type outputFrame struct {
	name string
	loc  string
	lace bool
}

func (vs *VMStacktrace) renderFrame(env *Env, ele Object) outputFrame {
	var str string

	switch sv := ele.(type) {
	case String:
		return outputFrame{name: sv.S}
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

				codeFile := fn.code.fileForIp(ip.I)
				if codeFile != fn.code.filename {
					macroLine := fn.code.macroLineForIp(ip.I)

					return outputFrame{
						lace: true,
						name: name,
						loc:  fmt.Sprintf("%s:%d (from %s:%d)", codeFile, macroLine, fn.code.filename, fn.code.lineForIp(ip.I)),
					}
				} else {
					return outputFrame{
						lace: true,
						name: name,
						loc:  fmt.Sprintf("%s:%d", fn.code.filename, fn.code.lineForIp(ip.I)),
					}
				}
			}
		}
	}

	var err error
	if str == "" {
		str, err = ele.ToString(env, false)
		if err != nil {
			str = fmt.Sprintf("error decoding stacktrace: %s\n", err)
		}
	}
	return outputFrame{
		name: str,
	}
}

const bcName = "github.com/lab47/lace/core.(*Engine).RunBC"

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
)

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
			} else {
				ofr = outputFrame{
					name: trimName(fr.Func),
					loc:  fmt.Sprintf("%s:%d", fr.File, fr.Line),
				}
			}

			oframes = append(oframes, ofr)

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

	str, err := vs.StackTrace.ToString(env, false)
	if err == nil {
		fmt.Fprintln(w, str)
	}
}

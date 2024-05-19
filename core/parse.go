package core

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/davecgh/go-spew/spew"
)

type (
	Expr interface {
		Eval(genv *Env, env *LocalEnv) (Object, error)
		InferType() *Type
		Pos() Position
		Dump(includePosition bool) Map
	}
	LiteralExpr struct {
		Position
		obj         Object
		isSurrogate bool
	}
	VectorExpr struct {
		Position
		v []Expr
	}
	MapExpr struct {
		Position
		keys   []Expr
		values []Expr
	}
	SetExpr struct {
		Position
		elements []Expr
	}
	IfExpr struct {
		Position
		cond     Expr
		positive Expr
		negative Expr
	}
	DefExpr struct {
		Position
		vr               *Var
		name             Symbol
		value            Expr
		meta             Expr
		isCreatedByMacro bool
	}
	CallExpr struct {
		Position
		callable Expr
		args     []Expr
	}
	MacroCallExpr struct {
		Position
		macro Callable
		args  []Object
		name  string
	}
	RecurExpr struct {
		Position
		args []Expr
	}
	VarRefExpr struct {
		Position
		vr *Var
	}
	BindingExpr struct {
		Position
		name    Symbol
		binding *Binding
	}
	MetaExpr struct {
		Position
		meta *MapExpr
		expr Expr
	}
	DoExpr struct {
		Position
		body             []Expr
		isCreatedByMacro bool
	}
	FnArityExpr struct {
		Position
		args       []Symbol
		body       []Expr
		taggedType *Type
	}
	FnExpr struct {
		Position
		arities  []FnArityExpr
		variadic *FnArityExpr
		self     Symbol

		compiled *Code
		closure  *fnClosure
	}
	LetExpr struct {
		Position
		names  []Symbol
		values []Expr
		body   []Expr
	}
	LoopExpr  LetExpr
	ThrowExpr struct {
		Position
		e Expr
	}
	CatchExpr struct {
		Position
		excType   *Type
		excSymbol Symbol
		body      []Expr
	}
	TryExpr struct {
		Position
		body        []Expr
		catches     []*CatchExpr
		finallyExpr []Expr
	}
	SetMacroExpr struct {
		Position
		vr *Var
	}
	ParseError struct {
		obj Object
		msg string
	}
	Callable interface {
		Call(env *Env, args []Object) (Object, error)
	}
	Binding struct {
		name   Symbol
		index  int
		frame  int
		isUsed bool
	}
	Bindings struct {
		bindings map[string]*Binding
		parent   *Bindings
		frame    int
	}
	LocalEnv struct {
		bindings []Object
		parent   *LocalEnv
		frame    int
	}
	ParseContext struct {
		Env                    *Env
		localBindings          *Bindings
		loopBindings           [][]Symbol
		linterBindings         *Bindings
		recur                  bool
		noRecurAllowed         bool
		isUnknownCallableScope bool
	}
	Warnings struct {
		ifWithoutElse           bool
		unusedFnParameters      bool
		fnWithEmptyBody         bool
		ignoredUnusedNamespaces Set
		IgnoredFileRegexes      []*regexp.Regexp
		entryPoints             Set
	}
	Keywords struct {
		tag                Keyword
		skipUnused         Keyword
		private            Keyword
		line               Keyword
		column             Keyword
		file               Keyword
		ns                 Keyword
		macro              Keyword
		message            Keyword
		form               Keyword
		data               Keyword
		cause              Keyword
		arglist            Keyword
		doc                Keyword
		added              Keyword
		meta               Keyword
		knownMacros        Keyword
		rules              Keyword
		ifWithoutElse      Keyword
		unusedFnParameters Keyword
		fnWithEmptyBody    Keyword
		_prefix            Keyword
		pos                Keyword
		startLine          Keyword
		endLine            Keyword
		startColumn        Keyword
		endColumn          Keyword
		filename           Keyword
		object             Keyword
		type_              Keyword
		var_               Keyword
		value              Keyword
		vector             Keyword
		name               Keyword
		dynamic            Keyword
	}
	Symbols struct {
		lace_core          Symbol
		underscore         Symbol
		catch              Symbol
		finally            Symbol
		amp                Symbol
		_if                Symbol
		quote              Symbol
		fn_                Symbol
		fn                 Symbol
		let_               Symbol
		letfn_             Symbol
		loop_              Symbol
		recur              Symbol
		setMacro_          Symbol
		def                Symbol
		defLinter          Symbol
		_var               Symbol
		do                 Symbol
		throw              Symbol
		try                Symbol
		unquoteSplicing    Symbol
		list               Symbol
		concat             Symbol
		concatS            Symbol
		seq                Symbol
		apply              Symbol
		emptySymbol        Symbol
		unquote            Symbol
		vector             Symbol
		hashMap            Symbol
		hashSet            Symbol
		defaultDataReaders Symbol
		backslash          Symbol
		deref              Symbol
	}
	Str struct {
		_if          string
		quote        string
		fn_          string
		let_         string
		letfn_       string
		loop_        string
		recur        string
		setMacro_    string
		def          string
		defLinter    string
		_var         string
		do           string
		throw        string
		try          string
		coreFilename string
	}
)

var (
	LOCAL_BINDINGS *Bindings = nil
	KNOWN_MACROS   *Var
	REQUIRE_VAR    *Var
	ALIAS_VAR      *Var
	REFER_VAR      *Var
	CREATE_NS_VAR  *Var
	IN_NS_VAR      *Var
	WARNINGS       = Warnings{
		fnWithEmptyBody: true,
		entryPoints:     EmptySet(),
	}
)

func (b *Bindings) ToMap(env *Env) (Map, error) {
	var res Map = EmptyArrayMap()
	for b != nil {
		for _, v := range b.bindings {
			v, err := res.Assoc(env, v.name, NIL)
			if err != nil {
				return nil, err
			}
			if err := Cast(env, v, &res); err != nil {
				return nil, err
			}
		}
		b = b.parent
	}
	return res, nil
}

func (localEnv *LocalEnv) addEmptyFrame(capacity int) *LocalEnv {
	res := LocalEnv{
		bindings: make([]Object, 0, capacity),
		parent:   localEnv,
	}
	if localEnv != nil {
		res.frame = localEnv.frame + 1
	}
	return &res
}

func (localEnv *LocalEnv) addBinding(obj Object) {
	localEnv.bindings = append(localEnv.bindings, obj)
}

func (localEnv *LocalEnv) addFrame(values []Object) *LocalEnv {
	res := LocalEnv{
		bindings: values,
		parent:   localEnv,
	}
	if localEnv != nil {
		res.frame = localEnv.frame + 1
	}
	return &res
}

func (localEnv *LocalEnv) replaceFrame(values []Object) *LocalEnv {
	res := LocalEnv{
		bindings: values,
		parent:   localEnv.parent,
		frame:    localEnv.frame,
	}
	return &res
}

func (ctx *ParseContext) PushLoopBindings(bindings []Symbol) {
	ctx.loopBindings = append(ctx.loopBindings, bindings)
}

func (ctx *ParseContext) PopLoopBindings() {
	ctx.loopBindings = ctx.loopBindings[:len(ctx.loopBindings)-1]
}

func (ctx *ParseContext) GetLoopBindings() []Symbol {
	n := len(ctx.loopBindings)
	if n == 0 {
		return nil
	}
	return ctx.loopBindings[n-1]
}

func (b *Bindings) PushFrame() *Bindings {
	frame := 0
	if b != nil {
		frame = b.frame + 1
	}
	return &Bindings{
		bindings: make(map[string]*Binding),
		parent:   b,
		frame:    frame,
	}
}

func (b *Bindings) PopFrame() *Bindings {
	return b.parent
}

func (b *Bindings) AddBinding(sym Symbol, index int, skipUnused bool) {
	if LINTER_MODE && !skipUnused {
		old := b.bindings[sym.Name()]
		if old != nil && needsUnusedWarning(old) {
			printParseWarning(GetPosition(old.name), "Unused binding: "+old.name.String())
		}
	}
	b.bindings[sym.Name()] = &Binding{
		name:  sym,
		frame: b.frame,
		index: index,
	}
}

func (ctx *ParseContext) PushEmptyLocalFrame() {
	ctx.localBindings = ctx.localBindings.PushFrame()
}

func (ctx *ParseContext) PushLocalFrame(names []Symbol) {
	ctx.PushEmptyLocalFrame()
	for i, sym := range names {
		ctx.localBindings.AddBinding(sym, i, true)
	}
}

func (ctx *ParseContext) PopLocalFrame() {
	ctx.localBindings = ctx.localBindings.PopFrame()
}

func (b *Bindings) GetBinding(sym Symbol) *Binding {
	env := b
	for env != nil {
		if b, ok := env.bindings[sym.Name()]; ok {
			return b
		}
		env = env.parent
	}
	return nil
}

func (ctx *ParseContext) GetLocalBinding(sym Symbol) *Binding {
	if sym.Namespace() != "" {
		return nil
	}
	return ctx.localBindings.GetBinding(sym)
}

func (pos Position) Pos() Position {
	return pos
}

func printError(pos Position, msg string) {
	PROBLEM_COUNT++
	fmt.Fprintf(Stderr, "%s:%d:%d: %s\n", pos.Filename(), pos.startLine, pos.startColumn, msg)
}

func printParseWarning(pos Position, msg string) {
	printError(pos, "Parse warning: "+msg)
}

func printParseError(pos Position, msg string) {
	printError(pos, "Parse error: "+msg)
}

func printReadWarning(reader *Reader, msg string) {
	pos := Position{
		filename:    reader.filename,
		startColumn: reader.column,
		startLine:   reader.line,
	}
	printError(pos, "Read warning: "+msg)
}

func printReadError(reader *Reader, msg string) {
	pos := Position{
		filename:    reader.filename,
		startColumn: reader.column,
		startLine:   reader.line,
	}
	printError(pos, "Read error: "+msg)
}

func isIgnoredUnusedNamespace(ns *Namespace) bool {
	if WARNINGS.ignoredUnusedNamespaces == nil {
		return false
	}
	return WARNINGS.ignoredUnusedNamespaces.Has(ns.Name)
}

func isEntryPointNs(ns *Namespace) bool {
	return WARNINGS.entryPoints.Has(ns.Name)
}

func WarnOnGloballyUnusedNamespaces(env *Env) {
	var names []string
	positions := make(map[string]Position)

	for _, ns := range env.Namespaces {
		if !ns.isGloballyUsed && !isIgnoredUnusedNamespace(ns) && !isEntryPointNs(ns) {
			pos := GetInfo(ns.Name)
			if pos != nil && pos.Filename() != "<lace.core>" && pos.Filename() != "<user>" {
				name := ns.Name.String()
				names = append(names, name)
				positions[name] = pos.Position
			}
		}
	}

	sort.Strings(names)
	for _, name := range names {
		printParseWarning(positions[name], "globally unused namespace "+name)
	}
}

func WarnOnUnusedNamespaces(env *Env) {
	var names []string
	positions := make(map[string]Position)

	for _, ns := range env.Namespaces {
		if ns != env.CurrentNamespace() && !ns.isUsed && !isIgnoredUnusedNamespace(ns) {
			pos := GetInfo(ns.Name)
			if pos != nil && pos.Filename() != "<lace.core>" && pos.Filename() != "<user>" {
				name := ns.Name.String()
				names = append(names, name)
				positions[name] = pos.Position
			}
		}
	}

	sort.Strings(names)
	for _, name := range names {
		printParseWarning(positions[name], "unused namespace "+name)
	}
}

func isEntryPointVar(vr *Var) bool {
	if isEntryPointNs(vr.ns) {
		return true
	}
	sym := AssembleSymbol(vr.ns.Name.Name(), vr.name.Name())
	return WARNINGS.entryPoints.Has(sym)
}

func NewLiteralExpr(obj Object) *LiteralExpr {
	res := LiteralExpr{obj: obj}
	info := GetInfo(obj)
	if info != nil {
		res.Position = info.Position
	}
	return &res
}

func NewSurrogateExpr(obj Object) *LiteralExpr {
	res := NewLiteralExpr(obj)
	res.isSurrogate = true
	return res
}

var _ Object = &ParseError{}

func (err *ParseError) ToString(env *Env, escape bool) (string, error) {
	return err.Error(), nil
}

func (err *ParseError) Equals(env *Env, other interface{}) bool {
	return err == other
}

func (err *ParseError) GetInfo() *ObjectInfo {
	return nil
}

func (err *ParseError) GetType() *Type {
	return TYPE.ParseError
}

func (err *ParseError) Hash(env *Env) (uint32, error) {
	return HashPtr(err), nil
}

func (err *ParseError) WithInfo(info *ObjectInfo) Object {
	return err
}

func (err *ParseError) Message() Object {
	return MakeString(err.msg)
}

func (err ParseError) Error() string {
	line, column, filename := 0, 0, "<file>"
	info := GetInfo(err.obj)
	if info != nil {
		line, column, filename = info.startLine, info.startColumn, info.Filename()
	}
	return fmt.Sprintf("%s:%d:%d: Parse error: %s", filename, line, column, err.msg)
}

func parseSeq(seq Seq, ctx *ParseContext) ([]Expr, error) {
	res := make([]Expr, 0)
	for {
		empty, err := seq.IsEmpty(ctx.Env)
		if err != nil {
			return nil, err
		}

		if empty {
			break
		}

		v, err := seq.First(ctx.Env)
		if err != nil {
			return nil, err
		}
		pv, err := Parse(v, ctx)
		if err != nil {
			return nil, err
		}

		res = append(res, pv)
		seq, err = seq.Rest(ctx.Env)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func parseVector(v *Vector, pos Position, ctx *ParseContext) (Expr, error) {
	r := make([]Expr, v.count)
	var err error
	for i := 0; i < v.count; i++ {
		r[i], err = Parse(v.at(i), ctx)
		if err != nil {
			return nil, err
		}
	}
	return &VectorExpr{
		v:        r,
		Position: pos,
	}, nil
}

func parseMap(m Map, pos Position, ctx *ParseContext) (*MapExpr, error) {
	res := &MapExpr{
		keys:     make([]Expr, m.Count()),
		values:   make([]Expr, m.Count()),
		Position: pos,
	}
	var err error
	for iter, i := m.Iter(), 0; iter.HasNext(); i++ {
		p := iter.Next()
		res.keys[i], err = Parse(p.Key, ctx)
		if err != nil {
			return nil, err
		}
		res.values[i], err = Parse(p.Value, ctx)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func parseSet(s *MapSet, pos Position, ctx *ParseContext) (Expr, error) {
	res := &SetExpr{
		elements: make([]Expr, s.m.Count()),
		Position: pos,
	}
	for iter, i := iter(s.Seq()), 0; iter.HasNext(ctx.Env); i++ {
		v, err := iter.Next(ctx.Env)
		if err != nil {
			return nil, err
		}
		res.elements[i], err = Parse(v, ctx)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func checkForm(env *Env, obj Object, min int, max int) (int, error) {
	var seq Seq
	if err := Cast(env, obj, &seq); err != nil {
		return 0, err
	}

	c, err := SeqCount(env, seq)
	if err != nil {
		return 0, err
	}
	if c < min {
		v, err := seq.First(env)
		if err != nil {
			return 0, err
		}
		s, err := v.ToString(env, false)
		if err != nil {
			return 0, err
		}
		return 0, &ParseError{obj: obj, msg: "Too few arguments to " + s}
	}
	if c > max {
		v, err := seq.First(env)
		if err != nil {
			return 0, err
		}
		s, err := v.ToString(env, false)
		if err != nil {
			return 0, err
		}
		return 0, &ParseError{obj: obj, msg: "Too many arguments to " + s}
	}
	return c, nil
}

func GetPosition(obj Object) Position {
	info := GetInfo(obj)
	if info != nil {
		return info.Position
	}
	return Position{}
}

func updateVar(vr *Var, info *ObjectInfo, valueExpr Expr, sym Symbol) {
	vr.WithInfo(info)
	meta := GetMeta(sym)
	if meta != nil {
		if ok, p := meta.GetEqu(criticalKeywords.private); ok {
			vr.isPrivate = ToBool(p)
		}
		if ok, p := meta.GetEqu(criticalKeywords.dynamic); ok {
			vr.isDynamic = ToBool(p)
		}
		vr.taggedType = getTaggedType(sym)
	}
}

func isCreatedByMacro(env *Env, formSeq Seq) bool {
	f, err := formSeq.First(env)
	if err != nil {
		return false
	}
	info := GetInfo(f)
	return info != nil && info.Pos().filename == STR.coreFilename
}

func parseDef(obj Object, ctx *ParseContext, isForLinter bool) (*DefExpr, error) {
	count, err := checkForm(ctx.Env, obj, 2, 4)
	if err != nil {
		return nil, err
	}
	var seq Seq
	if err := Cast(ctx.Env, obj, &seq); err != nil {
		return nil, err
	}
	s, err := Second(ctx.Env, seq)
	if err != nil {
		return nil, err
	}
	var meta Map
	switch sym := Unwrap(s).(type) {
	case Symbol:
		if sym.Namespace() != "" && (AssembleSymbol("", sym.Namespace()) != ctx.Env.CurrentNamespace().Name) {
			return nil, &ParseError{
				msg: "Can't create defs outside of current ns",
				obj: obj,
			}
		}
		symWithoutNs := AssembleSymbol("", sym.Name())
		vr, err := ctx.Env.CurrentNamespace().Intern(ctx.Env, symWithoutNs)
		if err != nil {
			return nil, err
		}
		if isForLinter {
			vr.isGloballyUsed = true
		}
		res := &DefExpr{
			vr:               vr,
			name:             sym,
			value:            nil,
			Position:         GetPosition(obj),
			isCreatedByMacro: isCreatedByMacro(ctx.Env, seq),
		}
		meta = GetMeta(sym)
		if count == 3 {
			v, err := Third(ctx.Env, seq)
			if err != nil {
				return nil, err
			}
			res.value, err = Parse(v, ctx)
			if err != nil {
				return nil, err
			}
		} else if count == 4 {
			v, err := Fourth(ctx.Env, seq)
			if err != nil {
				return nil, err
			}
			res.value, err = Parse(v, ctx)
			if err != nil {
				return nil, err
			}
			docstring, err := Third(ctx.Env, seq)
			if err != nil {
				return nil, err
			}

			switch Unwrap(docstring).(type) {
			case String:
				if meta != nil {
					v, err := meta.Assoc(ctx.Env, criticalKeywords.doc, docstring)
					if err != nil {
						return nil, err
					}
					if err := Cast(ctx.Env, v, &meta); err != nil {
						return nil, err
					}
				} else {
					v, err := EmptyArrayMap().Assoc(ctx.Env, criticalKeywords.doc, docstring)
					if err != nil {
						return nil, err
					}

					if err := Cast(ctx.Env, v, &meta); err != nil {
						return nil, err
					}
				}
			default:
				return nil, &ParseError{obj: docstring, msg: "Docstring must be a string"}
			}
		}
		updateVar(vr, GetInfo(obj), res.value, sym)
		if meta != nil {
			res.meta, err = Parse(DeriveReadObject(obj, meta), ctx)
			if err != nil {
				return nil, err
			}
		}
		return res, nil
	default:
		return nil, &ParseError{obj: s, msg: "First argument to def must be a Symbol"}
	}
}

func parseBody(seq Seq, ctx *ParseContext) ([]Expr, error) {
	recur := ctx.recur
	ctx.recur = false
	defer func() { ctx.recur = recur }()
	res := make([]Expr, 0)
	for {
		empty, err := seq.IsEmpty(ctx.Env)
		if err != nil {
			return nil, err
		}

		if empty {
			break
		}

		ro, err := seq.First(ctx.Env)
		if err != nil {
			return nil, err
		}

		expr, err := Parse(ro, ctx)
		if err != nil {
			return nil, err
		}
		seq, err = seq.Rest(ctx.Env)
		if err != nil {
			return nil, err
		}
		empty, err = seq.IsEmpty(ctx.Env)
		if err != nil {
			return nil, err
		}

		if ctx.recur && !empty && !LINTER_MODE {
			return nil, &ParseError{obj: ro, msg: "Can only recur from tail position"}
		}
		res = append(res, expr)
		if LINTER_MODE {
			if defExpr, ok := expr.(*DefExpr); ok && !defExpr.isCreatedByMacro {
				printParseWarning(defExpr.Pos(), "inline def")
			} else if doExpr, ok := expr.(*DoExpr); ok && !doExpr.isCreatedByMacro {
				printParseWarning(doExpr.Pos(), "redundant do form")
			}
		}
	}
	return res, nil
}

func parseParams(env *Env, params Object) ([]Symbol, bool, error) {
	res := make([]Symbol, 0)
	var v *Vector
	if err := Cast(env, params, &v); err != nil {
		return nil, false, err
	}
	for i := 0; i < v.count; i++ {
		ro := v.at(i)
		sym := ro
		if !IsSymbol(sym) {
			if LINTER_MODE {
				sym = generateSymbol("linter")
			} else {
				s, err := sym.ToString(env, false)
				if err != nil {
					return nil, false, err
				}
				return nil, false, &ParseError{obj: ro, msg: "Unsupported binding form 1: " + s}
			}
		}
		if criticalSymbols.amp.Equals(env, sym) {
			if v.count > i+2 {
				ro := v.at(i + 2)
				s, err := ro.ToString(env, false)
				if err != nil {
					return nil, false, err
				}
				return nil, false, &ParseError{obj: ro, msg: "Unexpected parameter: " + s}
			}
			if v.count == i+2 {
				variadic := v.at(i + 1)
				if !IsSymbol(variadic) {
					if LINTER_MODE {
						variadic = generateSymbol("linter")
					} else {
						s, err := variadic.ToString(env, false)
						if err != nil {
							return nil, false, err
						}
						return nil, false, &ParseError{obj: variadic, msg: "Unsupported binding form 2: " + s}
					}
				}
				res = append(res, variadic.(Symbol))
				return res, true, nil
			} else {
				return res, false, nil
			}
		}
		res = append(res, sym.(Symbol))
	}
	return res, false, nil
}

func needsUnusedWarning(b *Binding) bool {
	return !b.isUsed &&
		!strings.HasPrefix(b.name.Name(), "_") &&
		!strings.HasPrefix(b.name.Name(), "&form") &&
		!strings.HasPrefix(b.name.Name(), "&env") &&
		!isSkipUnused(b.name)
}

func addArity(fn *FnExpr, sig Seq, ctx *ParseContext) error {
	params, err := sig.First(ctx.Env)
	if err != nil {
		return err
	}
	body, err := sig.Rest(ctx.Env)
	if err != nil {
		return err
	}
	args, isVariadic, err := parseParams(ctx.Env, params)
	if err != nil {
		return err
	}
	ctx.PushLocalFrame(args)
	defer ctx.PopLocalFrame()
	ctx.PushLoopBindings(args)
	defer ctx.PopLoopBindings()

	noRecurAllowed := ctx.noRecurAllowed
	ctx.noRecurAllowed = false
	defer func() { ctx.noRecurAllowed = noRecurAllowed }()

	bodye, err := parseBody(body, ctx)
	if err != nil {
		return err
	}

	var meta Object
	if err := Cast(ctx.Env, params, &meta); err != nil {
		return err
	}

	arity := FnArityExpr{
		Position:   GetPosition(sig),
		args:       args,
		body:       bodye,
		taggedType: getTaggedType(meta),
	}
	if isVariadic {
		if fn.variadic != nil {
			return &ParseError{obj: params, msg: "Can't have more than 1 variadic overload"}
		}
		for _, arity := range fn.arities {
			if len(arity.args) >= len(args) {
				return &ParseError{obj: params, msg: "Can't have fixed arity function with more params than variadic function"}
			}
		}
		fn.variadic = &arity
	} else {
		for _, arity := range fn.arities {
			if len(arity.args) == len(args) {
				return &ParseError{obj: params, msg: "Can't have 2 overloads with same arity"}
			}
		}
		if fn.variadic != nil && len(args) >= len(fn.variadic.args) {
			return &ParseError{obj: params, msg: "Can't have fixed arity function with more params than variadic function"}
		}
		fn.arities = append(fn.arities, arity)
	}

	if LINTER_MODE {
		if WARNINGS.fnWithEmptyBody {
			if len(arity.body) == 0 {
				printParseWarning(arity.Position, "fn form with empty body")
			}
		}

		if WARNINGS.unusedFnParameters {
			var unused []Symbol
			for _, b := range ctx.localBindings.bindings {
				if needsUnusedWarning(b) {
					unused = append(unused, b.name)
				}
			}
			sort.Sort(BySymbolName(unused))
			for _, u := range unused {
				printParseWarning(GetPosition(u), "unused parameter: "+u.String())
			}
		}
	}

	return nil
}

func wrapWithMeta(fnExpr *FnExpr, obj Object, ctx *ParseContext) (Expr, error) {
	var metao Meta
	if err := Cast(ctx.Env, obj, &metao); err != nil {
		return nil, err
	}

	meta := metao.GetMeta()
	if meta != nil {
		m, err := parseMap(meta, fnExpr.Pos(), ctx)
		if err != nil {
			return nil, err
		}
		return &MetaExpr{
			meta:     m,
			expr:     fnExpr,
			Position: fnExpr.Pos(),
		}, nil
	}
	return fnExpr, nil
}

// Examples:
// (fn f [] 1 2)
// (fn f ([] 1 2)
//
//	([a] a 3)
//	([a & b] a b))
func parseFn(obj Object, ctx *ParseContext) (Expr, error) {
	var seq Seq
	if err := Cast(ctx.Env, obj, &seq); err != nil {
		return nil, err
	}

	res := &FnExpr{Position: GetPosition(obj)}
	bodies, err := seq.Rest(ctx.Env)
	if err != nil {
		return nil, err
	}
	p, err := bodies.First(ctx.Env)
	if err != nil {
		return nil, err
	}
	if IsSymbol(p) { // self reference
		res.self = p.(Symbol)
		bodies, err = bodies.Rest(ctx.Env)
		if err != nil {
			return nil, err
		}
		p, err = bodies.First(ctx.Env)
		if err != nil {
			return nil, err
		}

		ctx.PushLocalFrame([]Symbol{res.self})
		defer ctx.PopLocalFrame()
	}
	if IsVector(p) { // single arity
		if err := addArity(res, bodies, ctx); err != nil {
			return nil, err
		}
	} else {
		// multiple arities
		ok, err := bodies.IsEmpty(ctx.Env)
		if err != nil {
			return nil, err
		}

		if ok {
			spew.Dump(bodies)
			return nil, &ParseError{obj: p, msg: "Parameter declaration missing"}
		}

		for {
			empty, err := bodies.IsEmpty(ctx.Env)
			if err != nil {
				return nil, err
			}

			if empty {
				break
			}

			body, err := bodies.First(ctx.Env)
			if err != nil {
				return nil, err
			}

			switch s := Unwrap(body).(type) {
			case Seq:
				params, err := s.First(ctx.Env)
				if err != nil {
					return nil, err
				}

				if !IsVector(params) {
					s, err := params.ToString(ctx.Env, false)
					if err != nil {
						return nil, err
					}

					return nil, &ParseError{obj: params, msg: "Parameter declaration must be a vector. Got: " + s}
				}
				if err := addArity(res, s, ctx); err != nil {
					return nil, err
				}
			default:
				ss, err := s.ToString(ctx.Env, false)
				if err != nil {
					return nil, err
				}

				return nil, &ParseError{obj: body, msg: "Function body must be a list. Got: " + ss}
			}
			bodies, err = bodies.Rest(ctx.Env)
			if err != nil {
				return nil, err
			}
		}
	}

	fn := &Fn{
		fnExpr: res,
	}

	closure, err := compileFn(ctx.Env, fn, nil)
	if err != nil {
		return nil, err
	}

	res.compiled = fn.code
	res.closure = closure

	return wrapWithMeta(res, obj, ctx)
}

func isCatch(env *Env, obj Object) bool {
	seq, ok := obj.(Seq)
	if !ok {
		return false
	}

	v, err := seq.First(env)
	if err != nil {
		return false
	}

	return criticalSymbols.catch.Is(v)
}

func isFinally(env *Env, obj Object) bool {
	seq, ok := obj.(Seq)
	if !ok {
		return false
	}

	v, err := seq.First(env)
	if err != nil {
		return false
	}

	return criticalSymbols.finally.Is(v)
}

func resolveType(obj Object, ctx *ParseContext) (*Type, error) {
	excType, err := Parse(obj, ctx)
	if err != nil {
		return nil, err
	}

	switch excType := excType.(type) {
	case *LiteralExpr:
		switch t := excType.obj.(type) {
		case *Type:
			return t, nil
		}
	}
	if LINTER_MODE {
		return TYPE.Error, nil
	}
	s, err := obj.ToString(ctx.Env, false)
	if err != nil {
		return nil, err
	}
	return nil, &ParseError{obj: obj, msg: "Unable to resolve type: " + s}
}

func parseCatch(obj Object, ctx *ParseContext) (*CatchExpr, error) {
	var seqo Seq
	if err := Cast(ctx.Env, obj, &seqo); err != nil {
		return nil, err
	}
	seq, err := seqo.Rest(ctx.Env)
	if err != nil {
		return nil, err
	}

	empty, err := seq.IsEmpty(ctx.Env)
	if err != nil {
		return nil, err
	}

	if empty {
		return nil, &ParseError{obj: obj, msg: "catch requires at least two arguments: type symbol and binding symbol"}
	}

	r, err := seq.Rest(ctx.Env)
	if err != nil {
		return nil, err
	}

	empty, err = r.IsEmpty(ctx.Env)
	if err != nil {
		return nil, err
	}

	if empty {
		return nil, &ParseError{obj: obj, msg: "catch requires at least two arguments: type symbol and binding symbol"}
	}
	excSymbol, err := Second(ctx.Env, seq)
	if err != nil {
		return nil, err
	}

	v, err := seq.First(ctx.Env)
	if err != nil {
		return nil, err
	}

	excType, err := resolveType(v, ctx)
	if err != nil {
		return nil, err
	}

	if !IsSymbol(excSymbol) {
		s, err := excSymbol.ToString(ctx.Env, false)
		if err != nil {
			return nil, err
		}

		return nil, &ParseError{obj: excSymbol, msg: "Bad binding form, expected symbol, got: " + s}
	}

	ctx.PushLocalFrame([]Symbol{excSymbol.(Symbol)})
	defer ctx.PopLocalFrame()
	r1, err := seq.Rest(ctx.Env)
	if err != nil {
		return nil, err
	}
	r2, err := r1.Rest(ctx.Env)
	if err != nil {
		return nil, err
	}
	bodye, err := parseBody(r2, ctx)
	if err != nil {
		return nil, err
	}
	return &CatchExpr{
		Position:  GetPosition(obj),
		excType:   excType,
		excSymbol: excSymbol.(Symbol),
		body:      bodye,
	}, nil
}

func parseFinally(body Seq, ctx *ParseContext) ([]Expr, error) {
	return parseBody(body, ctx)
}

func parseTry(obj Object, ctx *ParseContext) (*TryExpr, error) {
	const (
		Regular = iota
		Catch   = iota
		Finally = iota
	)
	res := &TryExpr{Position: GetPosition(obj)}
	lastType := Regular
	var seqo Seq
	if err := Cast(ctx.Env, obj, &seqo); err != nil {
		return nil, err
	}

	seq, err := seqo.Rest(ctx.Env)
	if err != nil {
		return nil, err
	}

	noRecurAllowed := ctx.noRecurAllowed
	ctx.noRecurAllowed = true
	defer func() { ctx.noRecurAllowed = noRecurAllowed }()

	for {
		empty, err := seq.IsEmpty(ctx.Env)
		if err != nil {
			return nil, err
		}

		if empty {
			break
		}

		obj, err = seq.First(ctx.Env)
		if err != nil {
			return nil, err
		}

		if lastType == Finally {
			return nil, &ParseError{obj: obj, msg: "finally clause must be last in try expression"}
		}
		if isCatch(ctx.Env, obj) {
			c, err := parseCatch(obj, ctx)
			if err != nil {
				return nil, err
			}
			res.catches = append(res.catches, c)
			lastType = Catch
		} else if isFinally(ctx.Env, obj) {
			var seqo Seq
			if err := Cast(ctx.Env, obj, &seqo); err != nil {
				return nil, err
			}

			r, err := seqo.Rest(ctx.Env)
			if err != nil {
				return nil, err
			}

			res.finallyExpr, err = parseFinally(r, ctx)
			if err != nil {
				return nil, err
			}
			lastType = Finally
		} else {
			if lastType == Catch {
				return nil, &ParseError{obj: obj, msg: "Only catch or finally clause can follow catch in try expression"}
			}
			b, err := Parse(obj, ctx)
			if err != nil {
				return nil, err
			}
			res.body = append(res.body, b)
		}
		seq, err = seq.Rest(ctx.Env)
		if err != nil {
			return nil, err
		}
	}
	if LINTER_MODE {
		if res.body == nil {
			printParseWarning(res.Pos(), "try form with empty body")
		}
		if res.catches == nil && res.finallyExpr == nil {
			printParseWarning(res.Pos(), "try form without catch or finally")
		}
		if res.finallyExpr != nil && len(res.finallyExpr) == 0 {
			printParseWarning(GetPosition(obj), "finally form with empty body")
		}
	}
	return res, nil
}

func parseLet(obj Object, ctx *ParseContext) (*LetExpr, error) {
	return parseLetLoop(obj, "let", ctx)
}

func parseLoop(obj Object, ctx *ParseContext) (*LoopExpr, error) {
	e, err := parseLetLoop(obj, "loop", ctx)
	if err != nil {
		return nil, err
	}
	return (*LoopExpr)(e), nil
}

func parseLetfn(obj Object, ctx *ParseContext) (*LoopExpr, error) {
	e, err := parseLetLoop(obj, "letfn", ctx)
	if err != nil {
		return nil, err
	}
	return (*LoopExpr)(e), nil
}

func isSkipUnused(obj Object) bool {
	if m := GetMeta(obj); m != nil {
		if ok, v := m.GetEqu(criticalKeywords.skipUnused); ok {
			return ToBool(v)
		}
	}
	return false
}

func parseLetLoop(obj Object, formName string, ctx *ParseContext) (*LetExpr, error) {
	res := &LetExpr{
		Position: GetPosition(obj),
	}

	var seqo Seq
	if err := Cast(ctx.Env, obj, &seqo); err != nil {
		return nil, err
	}

	bindings, err := Second(ctx.Env, seqo)
	if err != nil {
		return nil, err
	}

	switch b := Unwrap(bindings).(type) {
	case *Vector:
		if b.count%2 != 0 {
			return nil, &ParseError{obj: bindings, msg: formName + " requires an even number of forms in binding vector"}
		}
		if LINTER_MODE && formName != "loop" && b.count == 0 {
			pos := GetPosition(obj)
			printParseWarning(pos, formName+" form with empty bindings vector")
		}
		skipUnused := isSkipUnused(b)
		res.names = make([]Symbol, b.count/2)
		res.values = make([]Expr, b.count/2)
		ctx.PushEmptyLocalFrame()
		defer ctx.PopLocalFrame()

		var err error

		for i := 0; i < b.count/2; i++ {
			s := b.at(i * 2)
			switch sym := Unwrap(s).(type) {
			case Symbol:
				if sym.Namespace() != "" {
					msg := "Can't let qualified name: " + sym.String()
					if LINTER_MODE {
						printParseError(GetPosition(s), msg)
					} else {
						return nil, &ParseError{obj: s, msg: msg}
					}
				}
				res.names[i] = sym
			default:
				if LINTER_MODE {
					res.names[i] = generateSymbol("linter")
				} else {
					ss, err := sym.ToString(ctx.Env, false)
					if err != nil {
						return nil, err
					}

					return nil, &ParseError{obj: s, msg: "Unsupported binding form 3: " + ss}
				}
			}
			if formName != "letfn" {
				res.values[i], err = Parse(b.at(i*2+1), ctx)
				if err != nil {
					return nil, err
				}
			}
			ctx.localBindings.AddBinding(res.names[i], i, skipUnused)
		}

		if formName == "letfn" {
			for i := 0; i < b.count/2; i++ {
				res.values[i], err = Parse(b.at(i*2+1), ctx)
				if err != nil {
					return nil, err
				}
			}
		}

		if formName == "loop" {
			ctx.PushLoopBindings(res.names)
			defer ctx.PopLoopBindings()

			noRecurAllowed := ctx.noRecurAllowed
			ctx.noRecurAllowed = false
			defer func() { ctx.noRecurAllowed = noRecurAllowed }()
		}

		var seqo Seq
		if err := Cast(ctx.Env, obj, &seqo); err != nil {
			return nil, err
		}
		r1, err := seqo.Rest(ctx.Env)
		if err != nil {
			return nil, err
		}

		r2, err := r1.Rest(ctx.Env)
		if err != nil {
			return nil, err
		}

		res.body, err = parseBody(r2, ctx)
		if err != nil {
			return nil, err
		}

		if LINTER_MODE {
			if len(res.body) == 0 {
				pos := GetPosition(obj)
				printParseWarning(pos, formName+" form with empty body")
			}

			if !skipUnused {
				var unused []Symbol
				for _, b := range ctx.localBindings.bindings {
					if needsUnusedWarning(b) {
						unused = append(unused, b.name)
					}
				}
				sort.Sort(BySymbolName(unused))
				for _, u := range unused {
					s, err := u.ToString(ctx.Env, false)
					if err != nil {
						return nil, err
					}

					printParseWarning(GetPosition(u), "unused binding: "+s)
				}
			}
		}

	default:
		return nil, &ParseError{
			obj: obj,
			msg: fmt.Sprintf("%s requires a vector for its bindings, got %T", formName, bindings),
		}
	}
	return res, nil
}

func parseRecur(obj Object, ctx *ParseContext) (*RecurExpr, error) {
	if ctx.noRecurAllowed {
		return nil, &ParseError{obj: obj, msg: "Cannot recur across try"}
	}
	loopBindings := ctx.GetLoopBindings()
	if loopBindings == nil && !LINTER_MODE {
		return nil, &ParseError{obj: obj, msg: "No recursion point for recur"}
	}

	var seq Seq
	if err := Cast(ctx.Env, obj, &seq); err != nil {
		return nil, err
	}

	r, err := seq.Rest(ctx.Env)
	if err != nil {
		return nil, err
	}

	args, err := parseSeq(r, ctx)
	if err != nil {
		return nil, err
	}
	if len(loopBindings) != len(args) && !LINTER_MODE {
		return nil, &ParseError{obj: obj, msg: fmt.Sprintf("Mismatched argument count to recur, expected: %d args, got: %d", len(loopBindings), len(args))}
	}
	ctx.recur = true
	return &RecurExpr{
		args:     args,
		Position: GetPosition(obj),
	}, nil
}

func resolveMacro(obj Object, ctx *ParseContext) *Var {
	switch sym := Unwrap(obj).(type) {
	case Symbol:
		if ctx.GetLocalBinding(sym) != nil {
			return nil
		}
		vr, ok := ctx.Env.Resolve(sym)
		if !ok || !vr.isMacro || vr.GetStatic() == nil {
			return nil
		}
		vr.isUsed = true
		vr.isGloballyUsed = true
		vr.ns.isUsed = true
		vr.ns.isGloballyUsed = true
		return vr
	default:
		return nil
	}
}

func fixInfo(env *Env, obj Object, info *ObjectInfo) (Object, error) {
	switch s := Unwrap(obj).(type) {
	case Nil:
		return obj, nil
	case Seq:
		objs := make([]Object, 0, 8)
		for {
			empty, err := s.IsEmpty(env)
			if err != nil {
				return nil, err
			}

			if empty {
				break
			}

			v, err := s.First(env)
			if err != nil {
				return nil, err
			}
			t, err := fixInfo(env, v, info)
			if err != nil {
				return nil, err
			}
			objs = append(objs, t)
			s, err = s.Rest(env)
			if err != nil {
				return nil, err
			}
		}
		res := NewListFrom(objs...)
		if objInfo := GetInfo(obj); objInfo != nil {
			return res.WithInfo(objInfo), nil
		}
		return res.WithInfo(info), nil
	case *Vector:
		res := EmptyVector()
		for i := 0; i < s.count; i++ {
			t, err := fixInfo(env, s.at(i), info)
			if err != nil {
				return nil, err
			}
			res, _ = res.Conjoin(t)
		}
		res.meta = s.meta
		if objInfo := GetInfo(obj); objInfo != nil {
			return res.WithInfo(objInfo), nil
		}
		return res.WithInfo(info), nil
	case Map:
		res := EmptyArrayMap()
		iter := s.Iter()
		for iter.HasNext() {
			p := iter.Next()
			key, err := fixInfo(env, p.Key, info)
			if err != nil {
				return nil, err
			}
			value, err := fixInfo(env, p.Value, info)
			if err != nil {
				return nil, err
			}
			res.Add(env, key, value)
		}

		var metao Meta
		if err := Cast(env, s, &metao); err != nil {
			return nil, err
		}

		res.meta = metao.GetMeta()
		if objInfo := GetInfo(obj); objInfo != nil {
			return res.WithInfo(objInfo), nil
		}
		return res.WithInfo(info), nil
	default:
		return obj, nil
	}
}

func macroexpand1(env *Env, seq Seq, ctx *ParseContext) (Object, error) {
	op, err := seq.First(env)
	if err != nil {
		return nil, err
	}

	vr := resolveMacro(op, ctx)
	if vr != nil {
		m, err := ctx.localBindings.ToMap(env)
		if err != nil {
			return nil, err
		}
		r, err := seq.Rest(env)
		if err != nil {
			return nil, err
		}

		slice, err := ToSlice(env, r.Cons(m).Cons(seq))
		if err != nil {
			return nil, err
		}

		var callable Callable
		if err := Cast(env, vr.GetStatic(), &callable); err != nil {
			return nil, err
		}

		expr := &MacroCallExpr{
			Position: GetPosition(seq),
			macro:    callable,
			args:     slice,
			name:     varCallableString(vr),
		}
		v, err := Eval(env, expr, nil)
		if err != nil {
			return nil, err
		}
		return fixInfo(env, v, GetInfo(seq))
	} else {
		return seq, nil
	}
}

func reportNotAFunction(pos Position, name string) {
	printParseWarning(pos, name+" is not a function")
}

func getTaggedType(obj Object) *Type {
	if m := GetMeta(obj); m != nil {
		if ok, typeName := m.GetEqu(criticalKeywords.tag); ok {
			if typeSym, ok := typeName.(Symbol); ok {
				if t := TYPES[typeSym.Name()]; t != nil {
					return t
				}
			}
		}
	}
	return nil
}

func getTaggedTypes(obj Object) []*Type {
	var res []*Type
	if m := GetMeta(obj); m != nil {
		if ok, typeName := m.GetEqu(criticalKeywords.tag); ok {
			switch typeDecl := Unwrap(typeName).(type) {
			case Symbol:
				if t := TYPES[typeDecl.Name()]; t != nil {
					res = append(res, t)
				}
			case String:
				parts := strings.Split(typeDecl.S(), "|")
				for _, p := range parts {
					if t := TYPES[MakeSymbol(p).Name()]; t != nil {
						res = append(res, t)
					}
				}
			}
		}
	}
	return res
}

func isTypeOneOf(abstractTypes []*Type, concreteType *Type) bool {
	for _, t := range abstractTypes {
		if IsEqualOrImplements(t, concreteType) {
			return true
		}
	}
	return false
}

func typesString(env *Env, types []*Type) (string, error) {
	var b bytes.Buffer
	for i, t := range types {
		s, err := t.ToString(env, false)
		if err != nil {
			return "", err
		}
		b.WriteString(s)
		if i < len(types)-1 {
			b.WriteString(" or ")
		}
	}
	return b.String(), nil
}

func checkTypes(env *Env, declaredArgs []Symbol, call *CallExpr) (bool, error) {
	res := false
	for i, da := range declaredArgs {
		if declaredTypes := getTaggedTypes(da); len(declaredTypes) > 0 {
			passedType := call.args[i].InferType()
			if passedType != nil {
				if !isTypeOneOf(declaredTypes, passedType) {
					ts, err := typesString(env, declaredTypes)
					if err != nil {
						return false, err
					}
					printParseWarning(call.args[i].Pos(), fmt.Sprintf("arg[%d] of %s must have type %s, got %s", i, call.Name(), ts, passedType.Name()))
					res = true
				}
			}
		}
	}
	return res, nil
}

func selectArity(expr *FnExpr, passedArgsCount int) *FnArityExpr {
	for _, arity := range expr.arities {
		if len(arity.args) == passedArgsCount {
			return &arity
		}
	}
	if expr.variadic != nil && passedArgsCount >= len(expr.variadic.args)-1 {
		return expr.variadic
	}
	return nil
}

func reportWrongArity(env *Env, expr *FnExpr, isMacro bool, call *CallExpr, pos Position) (bool, error) {
	passedArgsCount := len(call.args)
	if isMacro {
		passedArgsCount += 2
	}
	if v := selectArity(expr, passedArgsCount); v != nil {
		return checkTypes(env, v.args, call)
	}
	printParseWarning(pos, fmt.Sprintf("Wrong number of args (%d) passed to %s", len(call.args), call.Name()))
	return true, nil
}

func checkArglist(env *Env, arglist Seq, passedArgsCount int) (bool, error) {
	for {
		empty, err := arglist.IsEmpty(env)
		if err != nil {
			return false, err
		}

		if empty {
			break
		}

		f, err := arglist.First(env)
		if err != nil {
			return false, err
		}
		if v, ok := f.(*Vector); ok {
			if v.Count() == passedArgsCount {
				return true, nil
			}

			if v.Count() >= 2 {
				n, err := v.Nth(env, v.Count()-2)
				if err != nil {
					return false, err
				}

				if n.Equals(env, criticalSymbols.amp) && passedArgsCount >= (v.Count()-2) {
					return true, nil
				}
			}
		}
		arglist, err = arglist.Rest(env)
		if err != nil {
			return false, err
		}
	}
	return false, nil
}

func setMacroMeta(env *Env, vr *Var) error {
	var err error
	var ass Associative
	if vr.meta == nil {
		ass, err = EmptyArrayMap().Assoc(env, criticalKeywords.macro, Boolean(true))
	} else {
		ass, err = vr.meta.Assoc(env, criticalKeywords.macro, Boolean(true))
	}

	if err := Cast(env, ass, &vr.meta); err != nil {
		return err
	}

	return err
}

func parseSetMacro(env *Env, obj Object, ctx *ParseContext) (Expr, error) {
	var seq Seq
	if err := Cast(env, obj, &seq); err != nil {
		return nil, err
	}

	s, err := Second(env, seq)
	if err != nil {
		return nil, err
	}

	expr, err := Parse(s, ctx)
	if err != nil {
		return nil, err
	}

	switch expr := expr.(type) {
	case *LiteralExpr:
		switch vr := Unwrap(expr.obj).(type) {
		case *Var:
			res := &SetMacroExpr{
				vr: vr,
			}
			_, err = res.Eval(env, nil)
			return res, err
		}
	}
	return nil, &ParseError{obj: obj, msg: "set-macro__ argument must be a var"}
}

func isKnownMacros(env *Env, sym Symbol) (bool, Seq) {
	if KNOWN_MACROS == nil {
		knownMacros := env.CoreNamespace.Resolve("*known-macros*")
		if knownMacros == nil {
			return false, nil
		}
		KNOWN_MACROS = knownMacros
	}

	var m Map
	if err := Cast(env, KNOWN_MACROS.GetStatic(), &m); err != nil {
		return false, nil
	}

	if ok, v := m.GetEqu(sym); ok {
		switch v := Unwrap(v).(type) {
		case Seqable:
			return true, v.Seq()
		default:
			return true, nil
		}
	}
	return false, nil
}

func isUnknownCallable(env *Env, expr Expr) (bool, Seq) {
	if !LINTER_MODE {
		return false, nil
	}
	if c, ok := expr.(*VarRefExpr); ok {
		if c.vr.isMacro {
			return true, nil
		}
		var sym Symbol
		if c.vr.ns != env.CurrentNamespace() && c.vr.ns != env.CoreNamespace {
			sym = AssembleSymbol(c.vr.ns.Name.Name(), c.vr.name.Name())
		} else {
			sym = MakeSymbol(c.vr.name.Name())
		}
		b, s := isKnownMacros(env, sym)
		if b {
			return b, s
		}
		if sym.Namespace() == "" && c.vr.ns != env.CoreNamespace {
			return true, nil
		}
	}
	return false, nil
}

func areAllLiteralExprs(exprs []Expr) bool {
	for _, expr := range exprs {
		if _, ok := expr.(*LiteralExpr); !ok {
			return false
		}
	}
	return true
}

func getRequireVar(ctx *ParseContext) *Var {
	if REQUIRE_VAR == nil {
		REQUIRE_VAR = ctx.Env.CoreNamespace.Resolve("require")
	}
	return REQUIRE_VAR
}

func getReferVar(ctx *ParseContext) *Var {
	if REFER_VAR == nil {
		REFER_VAR = ctx.Env.CoreNamespace.Resolve("refer")
	}
	return REFER_VAR
}

func getAliasVar(ctx *ParseContext) *Var {
	if ALIAS_VAR == nil {
		ALIAS_VAR = ctx.Env.CoreNamespace.Resolve("alias")
	}
	return ALIAS_VAR
}

func getCreateNsVar(ctx *ParseContext) *Var {
	if CREATE_NS_VAR == nil {
		CREATE_NS_VAR = ctx.Env.CoreNamespace.Resolve("create-ns")
	}
	return CREATE_NS_VAR
}

func getInNsVar(ctx *ParseContext) *Var {
	if IN_NS_VAR == nil {
		IN_NS_VAR = ctx.Env.CoreNamespace.Resolve("in-ns")
	}
	return IN_NS_VAR
}

func checkCall(env *Env, expr Expr, isMacro bool, call *CallExpr, pos Position) {
	argsCount := len(call.args)
	switch expr := expr.(type) {
	case *FnExpr:
		_, err := reportWrongArity(env, expr, isMacro, call, pos)
		if err != nil {
			printParseWarning(pos, err.Error())
		}
	case *MapExpr:
		if argsCount == 0 || argsCount > 2 {
			printParseWarning(pos, fmt.Sprintf("Wrong number of args (%d) passed to a map", argsCount))
		}
	case *SetExpr:
		if argsCount == 0 || argsCount > 1 {
			printParseWarning(pos, fmt.Sprintf("Wrong number of args (%d) passed to a set", argsCount))
		}
	case *LiteralExpr:
		if _, ok := expr.obj.(Callable); !ok && !expr.isSurrogate {
			reportNotAFunction(pos, call.Name())
			return
		}
		switch Unwrap(expr.obj).(type) {
		case Keyword:
			if argsCount == 0 || argsCount > 2 {
				printParseWarning(pos, fmt.Sprintf("Wrong number of args (%d) passed to %s", argsCount, call.Name()))
			}
		}
	case *RecurExpr:
		reportNotAFunction(pos, call.Name())
	case *ThrowExpr:
		reportNotAFunction(pos, call.Name())
	}
}

func parseList(env *Env, obj Object, ctx *ParseContext) (Expr, error) {
	var seq Seq
	if err := Cast(env, obj, &seq); err != nil {
		return nil, err
	}

	expanded, err := macroexpand1(env, seq, ctx)
	if err != nil {
		return nil, err
	}
	if expanded != obj {
		//str, _ := expanded.ToString(env, true)
		//fmt.Println(str)

		return Parse(expanded, ctx)
	}
	empty, err := seq.IsEmpty(env)
	if err != nil {
		return nil, err
	}

	if empty {
		return NewLiteralExpr(obj), nil
	}

	currentIsUnknownCallableScope := ctx.isUnknownCallableScope
	defer func() {
		ctx.isUnknownCallableScope = currentIsUnknownCallableScope
	}()

	ctx.isUnknownCallableScope = false

	pos := GetPosition(obj)
	first, err := seq.First(env)
	if err != nil {
		return nil, err
	}

	if v, ok := first.(Symbol); ok && v.Namespace() == "" {
		switch v.Name() {
		case STR.quote:
			sec, err := Second(env, seq)
			if err != nil {
				return nil, err
			}
			return NewLiteralExpr(sec), nil
		case STR._if:
			if _, err := checkForm(env, obj, 3, 4); err != nil {
				return nil, err
			}

			if LINTER_MODE {
				c, err := SeqCount(env, seq)
				if err != nil {
					return nil, err
				}
				if c < 4 && WARNINGS.ifWithoutElse {
					printParseWarning(pos, "missing else branch")
				}
			}
			sec, err := Second(env, seq)
			if err != nil {
				return nil, err
			}
			cond, err := Parse(sec, ctx)
			if err != nil {
				return nil, err
			}
			thi, err := Third(env, seq)
			if err != nil {
				return nil, err
			}
			positive, err := Parse(thi, ctx)
			if err != nil {
				return nil, err
			}
			fou, err := Fourth(env, seq)
			if err != nil {
				return nil, err
			}
			negative, err := Parse(fou, ctx)
			if err != nil {
				return nil, err
			}
			return &IfExpr{
				cond:     cond,
				positive: positive,
				negative: negative,
				Position: pos,
			}, nil
		case STR.fn_:
			return parseFn(obj, ctx)
		case STR.let_:
			return parseLet(obj, ctx)
		case STR.letfn_:
			return parseLetfn(obj, ctx)
		case STR.loop_:
			return parseLoop(obj, ctx)
		case STR.recur:
			return parseRecur(obj, ctx)

		// Vars' isMacro has to be properly set during parse stage
		// for linter mode to correctly handle arguments count.
		case STR.setMacro_:
			return parseSetMacro(env, obj, ctx)

		case STR.def:
			return parseDef(obj, ctx, false)
		case STR.defLinter:
			return parseDef(obj, ctx, true)
		case STR._var:
			if _, err := checkForm(env, obj, 2, 2); err != nil {
				return nil, err
			}

			obj, err := Second(env, seq)
			if err != nil {
				return nil, err
			}

			switch sym := Unwrap(obj).(type) {
			case Symbol:
				vr, ok := ctx.Env.Resolve(sym)
				if !ok {
					if !LINTER_MODE {
						return nil, &ParseError{obj: obj, msg: "Unable to resolve var " + sym.String() + " in this context"}
					}
					symNs := ctx.Env.NamespaceFor(ctx.Env.CurrentNamespace(), sym)
					if !ctx.isUnknownCallableScope {
						if symNs == nil || symNs == ctx.Env.CurrentNamespace() {
							printParseError(GetInfo(obj).Pos(), "Unable to resolve symbol 1: "+sym.String())
						}
					}
					vr, err = InternFakeSymbol(ctx.Env, symNs, sym)
					if err != nil {
						return nil, err
					}
				}
				vr.isUsed = true
				vr.isGloballyUsed = true
				vr.ns.isUsed = true
				vr.ns.isGloballyUsed = true
				return &LiteralExpr{
					obj:      vr,
					Position: pos,
				}, nil
			default:
				return nil, &ParseError{obj: obj, msg: "var's argument must be a symbol"}
			}
		case STR.do:
			r, err := seq.Rest(env)
			if err != nil {
				return nil, err
			}
			body, err := parseBody(r, ctx)
			if err != nil {
				return nil, err
			}
			res := &DoExpr{
				body:             body,
				Position:         pos,
				isCreatedByMacro: isCreatedByMacro(env, seq),
			}
			if LINTER_MODE {
				if len(res.body) == 0 {
					printParseWarning(pos, "do form with empty body")
				} else if len(res.body) == 1 {
					printParseWarning(pos, "redundant do form")
				}
			}
			return res, nil
		case STR.throw:
			sec, err := Second(env, seq)
			if err != nil {
				return nil, err
			}
			e, err := Parse(sec, ctx)
			if err != nil {
				return nil, err
			}
			return &ThrowExpr{
				Position: pos,
				e:        e,
			}, nil
		case STR.try:
			return parseTry(obj, ctx)
		}
	}

	if sym, ok := first.(Symbol); ok && sym.Namespace() == "" && strings.HasPrefix(sym.Name(), ".") {
		r, err := seq.Rest(env)
		if err != nil {
			return nil, err
		}
		args, err := parseSeq(r, ctx)
		if err != nil {
			return nil, err
		}

		if len(args) == 0 {
			return nil, fmt.Errorf("method expression must have at least 1 argument")
		}

		return &MethodExpr{
			Position: GetPosition(obj),
			name:     sym,
			method:   convertMethodName(sym.Name()[1:]),
			obj:      args[0],
			args:     args[1:],
		}, nil
	}

	ctx.isUnknownCallableScope = currentIsUnknownCallableScope
	callable, err := Parse(first, ctx)
	if err != nil {
		return nil, err
	}
	unknown, syms := isUnknownCallable(env, callable)
	if unknown {
		ctx.isUnknownCallableScope = true
		if syms != nil {
			ctx.linterBindings = ctx.linterBindings.PushFrame()
			defer func() {
				ctx.linterBindings = ctx.linterBindings.PopFrame()
			}()
			for {
				empty, err := syms.IsEmpty(env)
				if err != nil {
					return nil, err
				}

				if empty {
					break
				}

				v, err := syms.First(env)
				if err != nil {
					return nil, err
				}
				if sym, ok := v.(Symbol); ok {
					ctx.linterBindings.AddBinding(sym, 0, true)
				}
				syms, err = syms.Rest(env)
				if err != nil {
					return nil, err
				}
			}
		}
	} else {
		ctx.isUnknownCallableScope = false
	}

	r, err := seq.Rest(env)
	if err != nil {
		return nil, err
	}
	args, err := parseSeq(r, ctx)
	if err != nil {
		return nil, err
	}

	res := &CallExpr{
		callable: callable,
		args:     args,
		Position: pos,
	}

	return res, nil
}

func convertMethodName(sym string) string {
	capNext := true
	return strings.Map(func(r rune) rune {
		if capNext {
			capNext = false
			return unicode.ToUpper(r)
		}

		if r == '-' {
			capNext = true
			return -1
		}

		return r
	}, sym)
}

func InternFakeSymbol(env *Env, ns *Namespace, sym Symbol) (*Var, error) {
	if ns != nil {
		fakeSym := AssembleSymbol("", sym.Name())
		return ns.Intern(env, fakeSym)
	}
	fakeSym := AssembleSymbol("", sym.String())
	return env.CurrentNamespace().Intern(env, fakeSym)
}

func isInteropSymbol(sym Symbol) bool {
	return sym.Namespace() == "" && (strings.HasPrefix(sym.Name(), ".") || strings.HasSuffix(sym.Name(), ".") || strings.Contains(sym.Name(), "$"))
}

func isRecordConstructor(sym Symbol) bool {
	return sym.Namespace() == "" && (strings.HasPrefix(sym.Name(), "->") || strings.HasPrefix(sym.Name(), "map->"))
}

var fullClassNameRe = regexp.MustCompile(`.+\..+\.[A-Z].+`)

func isJavaSymbol(sym Symbol) bool {
	s := sym.Name()
	if sym.Namespace() != "" {
		s = sym.Namespace()
	}
	return fullClassNameRe.MatchString(s)
}

func MakeVarRefExpr(vr *Var, obj Object) *VarRefExpr {
	vr.isUsed = true
	vr.isGloballyUsed = true
	vr.ns.isUsed = true
	vr.ns.isGloballyUsed = true
	return &VarRefExpr{
		vr:       vr,
		Position: GetPosition(obj),
	}
}

func parseSymbol(obj Object, ctx *ParseContext) (Expr, error) {
	var sym Symbol
	if err := Cast(ctx.Env, obj, &sym); err != nil {
		return nil, err
	}

	b := ctx.GetLocalBinding(sym)
	if b != nil {
		b.isUsed = true
		return &BindingExpr{
			binding:  b,
			name:     sym,
			Position: GetPosition(obj),
		}, nil
	}
	if vr, ok := ctx.Env.Resolve(sym); ok {
		return MakeVarRefExpr(vr, obj), nil
	}
	if sym.Namespace() == "" && TYPES[sym.Name()] != nil {
		return &LiteralExpr{
			Position: GetPosition(obj),
			obj:      TYPES[sym.Name()],
		}, nil
	}

	if !LINTER_MODE {
		return nil, &ParseError{obj: obj, msg: "Unable to resolve symbol 2: " + sym.String()}
	}
	if DIALECT == CLJS && sym.Namespace() == "" {
		// Check if this is a "callable namespace"
		ns := ctx.Env.FindNamespace(sym)
		if ns == nil {
			ns = ctx.Env.CurrentNamespace().aliases[sym.Name()]
		}
		if ns != nil {
			ns.isUsed = true
			ns.isGloballyUsed = true
			return NewSurrogateExpr(obj), nil
		}
		// See if this is a JS interop (i.e. Math.PI)
		parts := strings.Split(sym.Name(), ".")
		if len(parts) > 1 && parts[0] != "" && parts[len(parts)-1] != "" {
			return parseSymbol(DeriveReadObject(obj, MakeSymbol(strings.Join(parts[:len(parts)-1], "."))), ctx)
		}
		// Check if this is a constructor call
		if len(parts) == 2 && parts[0] != "" && parts[len(parts)-1] == "" {
			if vr, ok := ctx.Env.Resolve(MakeSymbol(parts[0])); ok {
				return MakeVarRefExpr(vr, obj), nil
			}
		}
	}
	symNs := ctx.Env.NamespaceFor(ctx.Env.CurrentNamespace(), sym)
	if symNs == nil || symNs == ctx.Env.CurrentNamespace() {
		if isInteropSymbol(sym) || isJavaSymbol(sym) {
			return NewSurrogateExpr(sym), nil
		}
		if !ctx.isUnknownCallableScope {
			if ctx.linterBindings.GetBinding(sym) == nil {
				printParseError(GetInfo(obj).Pos(), "Unable to resolve symbol 3: "+sym.String())
			}
		}
	}
	vr, err := InternFakeSymbol(ctx.Env, symNs, sym)
	if err != nil {
		return nil, err
	}
	return MakeVarRefExpr(vr, obj), nil
}

func Parse(obj Object, ctx *ParseContext) (Expr, error) {
	pos := GetPosition(obj)
	var res Expr
	var err error
	canHaveMeta := false
	switch v := Unwrap(obj).(type) {
	case Int, String, Char, Double, *BigInt, *BigFloat, Boolean, Nil, *Ratio, Keyword, *Regex, *Type:
		res = NewLiteralExpr(obj)
	case *Vector:
		canHaveMeta = true
		res, err = parseVector(v, pos, ctx)
	case Map:
		canHaveMeta = true
		res, err = parseMap(v, pos, ctx)
	case *MapSet:
		canHaveMeta = true
		res, err = parseSet(v, pos, ctx)
	case Seq:
		res, err = parseList(ctx.Env, obj, ctx)
	case Symbol:
		res, err = parseSymbol(obj, ctx)
	default:
		s, err := obj.ToString(ctx.Env, false)
		if err != nil {
			return nil, err
		}
		return nil, &ParseError{obj: obj, msg: "Cannot parse form: " + s}
	}

	if err != nil {
		return nil, err
	}

	if canHaveMeta {
		var metao Meta
		if err := Cast(ctx.Env, obj, &metao); err != nil {
			return nil, err
		}

		meta := metao.GetMeta()
		if meta != nil {
			meta, err := parseMap(meta, pos, ctx)
			if err != nil {
				return nil, err
			}
			return &MetaExpr{
				meta:     meta,
				expr:     res,
				Position: pos,
			}, nil
		}
	}
	return res, nil
}

func TryParse(obj Object, ctx *ParseContext) (expr Expr, err error) {
	expr, err = Parse(obj, ctx)
	if err != nil {
		PROBLEM_COUNT++
	}

	return expr, err
}

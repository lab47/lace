//go:generate go run gen/gen_types.go assert Comparable *Vector Char String Symbol Keyword *Regex Boolean Time Number Seqable Callable *Type Meta Int Double Stack Map Set Associative Reversible Named Comparator *Ratio *Namespace *Var Error *Fn Deref *Atom Ref KVReduce Pending *File io.Reader io.Writer StringReader io.RuneReader *Channel
//go:generate go run gen/gen_types.go info *List *ArrayMapSeq *ArrayMap *HashMap *ExInfo *Fn *Var Nil *Ratio *BigInt *BigFloat Char Double Int Boolean Time Keyword *Regex Symbol String *LazySeq *MappingSeq *ArraySeq *ConsSeq *NodeSeq *ArrayNodeSeq *MapSet *Vector *VectorSeq *VectorRSeq
//go:generate go run -tags gen_data gen_data/gen_data.go

package core

import (
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"hash"
	"hash/fnv"
	"io"
	"math"
	"math/big"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"unsafe"
)

// interfaces
type (
	Object interface {
		Equality
		ToString(escape bool) string
		GetInfo() *ObjectInfo
		WithInfo(*ObjectInfo) Object
		GetType() *Type
		Hash() uint32
	}
	Equality interface {
		Equals(interface{}) bool
	}
	Conjable interface {
		Object
		Conj(obj Object) (Conjable, error)
	}
	Counted interface {
		Count() int
	}
	Error interface {
		error
		Object
		Message() Object
	}
	Meta interface {
		GetMeta() Map
		WithMeta(Map) (Object, error)
	}
	Ref interface {
		AlterMeta(env *Env, fn *Fn, args []Object) (Map, error)
		ResetMeta(m Map) Map
	}
	Sequential interface {
		sequential()
	}
	Comparable interface {
		Compare(env *Env, other Object) (int, error)
	}
	Comparator interface {
		Compare(env *Env, a, b Object) (int, error)
	}
	Indexed interface {
		Nth(i int) Object
		TryNth(i int, d Object) Object
	}
	Stack interface {
		Peek() Object
		Pop() Stack
	}
	Gettable interface {
		Get(key Object) (bool, Object)
	}
	Associative interface {
		Conjable
		Gettable
		EntryAt(key Object) (*Vector, error)
		Assoc(key, val Object) (Associative, error)
	}
	Reversible interface {
		Rseq() Seq
	}
	Named interface {
		Name() string
		Namespace() string
	}
	Printer interface {
		Print(writer io.Writer, printReadably bool)
	}
	Pprinter interface {
		Pprint(writer io.Writer, indent int) int
	}
	Collection interface {
		Object
		Counted
		Seqable
		Empty() Collection
	}
	Deref interface {
		Deref(env *Env) (Object, error)
	}
	Native interface {
		Native() interface{}
	}
	KVReduce interface {
		kvreduce(env *Env, c Callable, init Object) (Object, error)
	}
	Pending interface {
		IsRealized() bool
	}
)

// implementations
type (
	Position struct {
		endLine     int
		endColumn   int
		startLine   int
		startColumn int
		filename    *string
	}
	Atom struct {
		MetaHolder
		value Object
	}
	Type struct {
		MetaHolder
		name        string
		reflectType reflect.Type
	}
	MetaHolder struct {
		meta Map
	}
	ObjectInfo struct {
		Position
	}
	InfoHolder struct {
		info *ObjectInfo
	}
	Char struct {
		InfoHolder
		Ch rune
	}
	Double struct {
		InfoHolder
		D float64
	}
	Int struct {
		InfoHolder
		I int
	}
	BigInt struct {
		InfoHolder
		b big.Int
	}
	BigFloat struct {
		InfoHolder
		b big.Float
	}
	Ratio struct {
		InfoHolder
		r big.Rat
	}
	Boolean struct {
		InfoHolder
		B bool
	}
	Nil struct {
		InfoHolder
		n struct{}
	}
	Keyword struct {
		InfoHolder
		ns   *string
		name *string
		hash uint32
	}
	Symbol struct {
		InfoHolder
		MetaHolder
		ns   *string
		name *string
		hash uint32
	}
	String struct {
		InfoHolder
		S string
	}
	Regex struct {
		InfoHolder
		R *regexp.Regexp
	}
	Time struct {
		InfoHolder
		T time.Time
	}
	Var struct {
		InfoHolder
		MetaHolder
		ns             *Namespace
		name           Symbol
		Value          Object
		expr           Expr
		isMacro        bool
		isPrivate      bool
		isDynamic      bool
		isUsed         bool
		isGloballyUsed bool
		taggedType     *Type
	}
	ProcFn func(env *Env, args []Object) (Object, error)
	Proc   struct {
		Fn      ProcFn
		Name    string
		Package string // "" for core (this package), else e.g. "std/string"
		File    string
		Line    int
	}
	Fn struct {
		InfoHolder
		MetaHolder
		isMacro bool
		fnExpr  *FnExpr
		env     *LocalEnv
	}
	ExInfo struct {
		ArrayMap
		rt *Runtime
	}
	RecurBindings []Object
	Delay         struct {
		fn    Callable
		value Object
	}
	SortableSlice struct {
		env *Env
		s   []Object
		cmp Comparator
		err error
	}
	Types struct {
		Associative    *Type
		Callable       *Type
		Collection     *Type
		Comparable     *Type
		Comparator     *Type
		Counted        *Type
		Deref          *Type
		Channel        *Type
		Error          *Type
		Gettable       *Type
		Indexed        *Type
		IOReader       *Type
		IOWriter       *Type
		KVReduce       *Type
		Map            *Type
		Meta           *Type
		Named          *Type
		Number         *Type
		Pending        *Type
		Ref            *Type
		Reversible     *Type
		Seq            *Type
		Seqable        *Type
		Sequential     *Type
		Set            *Type
		Stack          *Type
		ArrayMap       *Type
		ArrayMapSeq    *Type
		ArrayNodeSeq   *Type
		ArraySeq       *Type
		MapSet         *Type
		Atom           *Type
		BigFloat       *Type
		BigInt         *Type
		Boolean        *Type
		Time           *Type
		Buffer         *Type
		Char           *Type
		ConsSeq        *Type
		Delay          *Type
		Double         *Type
		EvalError      *Type
		ExInfo         *Type
		Fn             *Type
		File           *Type
		BufferedReader *Type
		HashMap        *Type
		Int            *Type
		Keyword        *Type
		LazySeq        *Type
		List           *Type
		Opaque         *Type
		MappingSeq     *Type
		Namespace      *Type
		Nil            *Type
		NodeSeq        *Type
		ParseError     *Type
		Proc           *Type
		ProcFn         *Type
		Ratio          *Type
		RecurBindings  *Type
		Regex          *Type
		String         *Type
		Symbol         *Type
		Type           *Type
		Var            *Type
		Vector         *Type
		VectorRSeq     *Type
		VectorSeq      *Type
	}
)

// interface checks
var (
	_ Conjable = &HashMap{}
	_ Conjable = &Vector{}
	_ Conjable = Nil{}
	_ Conjable = &ArrayMap{}
	_ Conjable = &MapSet{}
	_ Conjable = &List{}

	_ Counted = &Vector{}
	_ Counted = &List{}
	_ Counted = Nil{}
	_ Counted = String{}
	_ Counted = &HashMap{}
	_ Counted = &ArrayMap{}
	_ Counted = &MapSet{}

	_ Meta = &Vector{}
	_ Meta = &VectorSeq{}
	_ Meta = &VectorRSeq{}
	_ Meta = &Namespace{}
	_ Meta = &List{}
	_ Meta = &Atom{}
	_ Meta = &Fn{}
	_ Meta = Symbol{}
	_ Meta = &ArrayNodeSeq{}
	_ Meta = &NodeSeq{}
	_ Meta = &HashMap{}
	_ Meta = &MappingSeq{}
	_ Meta = &LazySeq{}
	_ Meta = &ArraySeq{}
	_ Meta = &ConsSeq{}
	_ Meta = &ArrayMapSeq{}
	_ Meta = &ArrayMap{}
	_ Meta = &MapSet{}

	_ Ref = &Namespace{}
	_ Ref = &Atom{}
	_ Ref = &Var{}

	_ Sequential = &VectorSeq{}
	_ Sequential = &VectorRSeq{}
	_ Sequential = &Vector{}
	_ Sequential = &List{}
	_ Sequential = &MappingSeq{}
	_ Sequential = &LazySeq{}
	_ Sequential = &ArraySeq{}
	_ Sequential = &ConsSeq{}
	_ Sequential = &ArrayMapSeq{}
	_ Sequential = &ArrayNodeSeq{}
	_ Sequential = &NodeSeq{}

	_ Comparable = &Ratio{}
	_ Comparable = &BigInt{}
	_ Comparable = &BigFloat{}
	_ Comparable = Char{}
	_ Comparable = Double{}
	_ Comparable = Int{}
	_ Comparable = Boolean{}
	_ Comparable = Time{}
	_ Comparable = Keyword{}
	_ Comparable = Symbol{}
	_ Comparable = &Vector{}

	_ Comparator = &Fn{}
	_ Comparator = Proc{}

	_ Indexed = String{}
	_ Indexed = &Vector{}

	_ Stack = &Vector{}
	_ Stack = &List{}

	_ Gettable = &HashMap{}
	_ Gettable = Nil{}
	_ Gettable = &Vector{}
	_ Gettable = &ArrayMap{}
	_ Gettable = &MapSet{}

	_ Associative = Nil{}
	_ Associative = &Vector{}
	_ Associative = &HashMap{}
	_ Associative = &ArrayMap{}

	_ Reversible = &Vector{}

	_ Named = Keyword{}
	_ Named = Symbol{}

	_ Printer = &Namespace{}
	_ Printer = &Regex{}

	_ Pprinter = &List{}
	_ Pprinter = &VectorSeq{}
	_ Pprinter = &VectorRSeq{}
	_ Pprinter = &Vector{}
	_ Pprinter = &ArrayNodeSeq{}
	_ Pprinter = &NodeSeq{}
	_ Pprinter = &HashMap{}
	_ Pprinter = &MappingSeq{}
	_ Pprinter = &LazySeq{}
	_ Pprinter = &ArraySeq{}
	_ Pprinter = &ConsSeq{}
	_ Pprinter = &ArrayMapSeq{}
	_ Pprinter = &ArrayMap{}
	_ Pprinter = &MapSet{}

	_ Collection = &Vector{}
	_ Collection = &List{}
	_ Collection = &HashMap{}
	_ Collection = &ArrayMap{}
	_ Collection = &MapSet{}

	_ Deref = &Atom{}
	_ Deref = &Delay{}
	_ Deref = &Var{}

	_ KVReduce = &Vector{}

	_ Pending = &Delay{}
	_ Pending = &LazySeq{}
)

func (pos Position) Filename() string {
	if pos.filename == nil {
		return "<file>"
	}
	return *pos.filename
}

var hasher hash.Hash32 = fnv.New32a()

func newIteratorError() error {
	return errors.New("iterator reached the end of collection")
}

func uint32ToBytes(i uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, i)
	return b
}

func getHash() hash.Hash32 {
	hasher.Reset()
	return hasher
}

func hashSymbol(ns, name *string) uint32 {
	h := getHash()
	if ns != nil {
		h.Write([]byte(*ns))
	}
	h.Write([]byte("/" + *name))
	return h.Sum32()
}

func MakeSymbol(nsname string) Symbol {
	index := strings.IndexRune(nsname, '/')
	if index == -1 || nsname == "/" {
		return Symbol{
			ns:   nil,
			name: STRINGS.Intern(nsname),
		}
	}
	return Symbol{
		ns:   STRINGS.Intern(nsname[0:index]),
		name: STRINGS.Intern(nsname[index+1:]),
	}
}

type BySymbolName []Symbol

func (s BySymbolName) Len() int {
	return len(s)
}
func (s BySymbolName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s BySymbolName) Less(i, j int) bool {
	return s[i].ToString(false) < s[j].ToString(false)
}

const KeywordHashMask uint32 = 0x7334c790

func MakeKeyword(nsname string) Keyword {
	index := strings.IndexRune(nsname, '/')
	if index == -1 || nsname == "/" {
		name := STRINGS.Intern(nsname)
		return Keyword{
			ns:   nil,
			name: name,
			hash: hashSymbol(nil, name) ^ KeywordHashMask,
		}
	}
	ns := STRINGS.Intern(nsname[0:index])
	name := STRINGS.Intern(nsname[index+1:])
	return Keyword{
		ns:   ns,
		name: name,
		hash: hashSymbol(ns, name) ^ KeywordHashMask,
	}
}

func ErrorArity(env *Env, n int) error {
	name := env.RT.topName()
	return env.RT.NewError(fmt.Sprintf("Wrong number of args (%d) passed to %s", n, name))
}

func rangeString(min, max int) string {
	if min == max {
		return strconv.Itoa(min)
	}
	if min+1 == max {
		return strconv.Itoa(min) + " or " + strconv.Itoa(max)
	}
	if min+2 == max {
		return strconv.Itoa(min) + ", " + strconv.Itoa(min+1) + ", or " + strconv.Itoa(max)
	}
	if max >= 999 {
		return "at least " + strconv.Itoa(min)
	}
	return "between " + strconv.Itoa(min) + " and " + strconv.Itoa(max) + ", inclusive"
}

func ErrorArityMinMax(env *Env, n, min, max int) error {
	name := env.RT.topName()
	return env.RT.NewError(fmt.Sprintf("Wrong number of args (%d) passed to %s; expects %s", n, name, rangeString(min, max)))
}

func ReturnArityMinMax(env *Env, n, min, max int) error {
	name := env.RT.topName()
	return env.RT.NewError(fmt.Sprintf("Wrong number of args (%d) passed to %s; expects %s", n, name, rangeString(min, max)))
}

func CheckArity(env *Env, args []Object, min int, max int) error {
	n := len(args)
	if n < min || n > max {
		return ReturnArityMinMax(env, n, min, max)
	}
	return nil
}

func getMap(env *Env, k Object, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 2); err != nil {
		return nil, err
	}

	switch m := args[0].(type) {
	case Map:
		ok, v := m.Get(k)
		if ok {
			return v, nil
		}
	}
	if len(args) == 2 {
		return args[1], nil
	}
	return NIL, nil
}

func (s *SortableSlice) Len() int {
	return len(s.s)
}

func (s *SortableSlice) Swap(i, j int) {
	s.s[i], s.s[j] = s.s[j], s.s[i]
}

func (s *SortableSlice) Less(i, j int) bool {
	cmp, err := s.cmp.Compare(s.env, s.s[i], s.s[j])
	if err != nil {
		s.err = err
		return false
	}

	return cmp == -1
}

func HashPtr(ptr uintptr) uint32 {
	h := getHash()
	b := make([]byte, unsafe.Sizeof(ptr))
	b[0] = byte(ptr)
	b[1] = byte(ptr >> 8)
	b[2] = byte(ptr >> 16)
	b[3] = byte(ptr >> 24)
	if unsafe.Sizeof(ptr) == 8 {
		b[4] = byte(ptr >> 32)
		b[5] = byte(ptr >> 40)
		b[6] = byte(ptr >> 48)
		b[7] = byte(ptr >> 56)
	}
	h.Write(b)
	return h.Sum32()
}

func hashGobEncoder(e gob.GobEncoder) (uint32, error) {
	h := getHash()
	b, err := e.GobEncode()
	if err != nil {
		return 0, err
	}
	h.Write(b)
	return h.Sum32(), nil
}

func equalsNumbers(x Number, y interface{}) bool {
	switch y := y.(type) {
	case Number:
		return category(x) == category(y) && numbersEq(x, y)
	default:
		return false
	}
}

func (a *Atom) ToString(escape bool) string {
	return "#object[Atom {:val " + a.value.ToString(escape) + "}]"
}

func (a *Atom) Equals(other interface{}) bool {
	return a == other
}

func (a *Atom) GetInfo() *ObjectInfo {
	return nil
}

func (a *Atom) GetType() *Type {
	return TYPE.Atom
}

func (a *Atom) Hash() uint32 {
	return HashPtr(uintptr(unsafe.Pointer(a)))
}

func (a *Atom) WithInfo(info *ObjectInfo) Object {
	return a
}

func (a *Atom) WithMeta(meta Map) (Object, error) {
	res := *a
	m, err := SafeMerge(res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return &res, nil
}

func (a *Atom) ResetMeta(newMeta Map) Map {
	a.meta = newMeta
	return a.meta
}

func (a *Atom) AlterMeta(env *Env, fn *Fn, args []Object) (Map, error) {
	return AlterMeta(env, &a.MetaHolder, fn, args)
}

func (a *Atom) Deref(env *Env) (Object, error) {
	return a.value, nil
}

func (d *Delay) ToString(escape bool) string {
	return "#object[Delay]"
}

func (d *Delay) Equals(other interface{}) bool {
	return d == other
}

func (d *Delay) GetInfo() *ObjectInfo {
	return nil
}

func (d *Delay) GetType() *Type {
	return TYPE.Delay
}

func (d *Delay) Hash() uint32 {
	return HashPtr(uintptr(unsafe.Pointer(d)))
}

func (d *Delay) WithInfo(info *ObjectInfo) Object {
	return d
}

func (d *Delay) Force(env *Env) (Object, error) {
	if d.value == nil {
		val, err := d.fn.Call(env, []Object{})
		if err != nil {
			return nil, err
		}
		d.value = val
	}
	return d.value, nil
}

func (d *Delay) Deref(env *Env) (Object, error) {
	return d.Force(env)
}

func (d *Delay) IsRealized() bool {
	return d.value != nil
}

func (t *Type) ToString(escape bool) string {
	return t.name
}

func (t *Type) Equals(other interface{}) bool {
	return t == other
}

func (t *Type) GetInfo() *ObjectInfo {
	return nil
}

func (t *Type) GetType() *Type {
	return TYPE.Type
}

func (t *Type) Hash() uint32 {
	return HashPtr(uintptr(unsafe.Pointer(t)))
}

func (rb RecurBindings) ToString(escape bool) string {
	return "#object[RecurBindings]"
}

func (rb RecurBindings) Equals(other interface{}) bool {
	return false
}

func (rb RecurBindings) GetInfo() *ObjectInfo {
	return nil
}

func (rb RecurBindings) GetType() *Type {
	return TYPE.RecurBindings
}

func (rb RecurBindings) Hash() uint32 {
	return 0
}

func (exInfo *ExInfo) ToString(escape bool) string {
	return exInfo.Error()
}

func (exInfo *ExInfo) Equals(other interface{}) bool {
	return exInfo == other
}

func (exInfo *ExInfo) GetType() *Type {
	return TYPE.ExInfo
}

func (exInfo *ExInfo) Hash() uint32 {
	return HashPtr(uintptr(unsafe.Pointer(exInfo)))
}

func (exInfo *ExInfo) Message() Object {
	if ok, res := exInfo.Get(criticalKeywords.message); ok {
		return res
	}
	return NIL
}

func (exInfo *ExInfo) Error() string {
	var pos Position
	_, data := exInfo.Get(criticalKeywords.data)
	ok, form := data.(Map).Get(criticalKeywords.form)
	if ok {
		if form.GetInfo() != nil {
			pos = form.GetInfo().Pos()
		}
	}
	prefix := "Exception"
	if ok, pr := data.(Map).Get(criticalKeywords._prefix); ok {
		prefix = pr.ToString(false)
	}
	_, msg := exInfo.Get(criticalKeywords.message)
	if len(exInfo.rt.callstack.frames) > 0 && !LINTER_MODE {
		return fmt.Sprintf("%s:%d:%d: %s: %s\nStacktrace:\n%s", pos.Filename(), pos.startLine, pos.startColumn, prefix, msg.(String).S, exInfo.rt.stacktrace())
	} else {
		return fmt.Sprintf("%s:%d:%d: %s: %s", pos.Filename(), pos.startLine, pos.startColumn, prefix, msg.(String).S)
	}
}

func (fn *Fn) ToString(escape bool) string {
	return "#object[Fn]"
}

func (fn *Fn) Equals(other interface{}) bool {
	switch other := other.(type) {
	case *Fn:
		return fn == other
	default:
		return false
	}
}

func (fn *Fn) WithMeta(meta Map) (Object, error) {
	res := *fn
	m, err := SafeMerge(res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return &res, nil
}

func (fn *Fn) GetType() *Type {
	return TYPE.Fn
}

func (fn *Fn) Hash() uint32 {
	return HashPtr(uintptr(unsafe.Pointer(fn)))
}

func (fn *Fn) Call(env *Env, args []Object) (Object, error) {
	min := math.MaxInt32
	max := -1
	for _, arity := range fn.fnExpr.arities {
		a := len(arity.args)
		if a == len(args) {
			//env.RT.pushFrame()
			//defer env.RT.popFrame()
			return evalLoop(env, arity.body, fn.env.addFrame(args))
		}
		if min > a {
			min = a
		}
		if max < a {
			max = a
		}
	}
	v := fn.fnExpr.variadic
	if v == nil || len(args) < len(v.args)-1 {
		if v != nil {
			min = len(v.args)
			max = math.MaxInt32
		}
		c := len(args)
		if fn.isMacro {
			c -= 2
			min -= 2
			if max != math.MaxInt32 {
				max -= 2
			}
		}
		return nil, ErrorArityMinMax(env, c, min, max)
	}
	var restArgs Object = NIL
	if len(v.args)-1 < len(args) {
		restArgs = &ArraySeq{arr: args, index: len(v.args) - 1}
	}
	vargs := make([]Object, len(v.args))
	for i := 0; i < len(vargs)-1; i++ {
		vargs[i] = args[i]
	}
	vargs[len(vargs)-1] = restArgs
	//env.RT.pushFrame()
	//defer env.RT.popFrame()
	return evalLoop(env, v.body, fn.env.addFrame(vargs))
}

func compare(env *Env, c Callable, a, b Object) (int, error) {
	val, err := c.Call(env, []Object{a, b})
	if err != nil {
		return 0, err
	}

	switch r := val.(type) {
	case Boolean:
		if r.B {
			return -1, nil
		}

		v, err := c.Call(env, []Object{b, a})
		if err != nil {
			return 0, err
		}

		b, err := AssertBoolean(env, v, "")
		if err != nil {
			return 0, err
		}

		if b.B {
			return 1, nil
		}
		return 0, nil
	default:
		a, err := AssertNumber(env, r, "Function is not a comparator since it returned a non-integer value")
		if err != nil {
			return 0, err
		}

		return a.Int().I, nil
	}
}

func (fn *Fn) Compare(env *Env, a, b Object) (int, error) {
	return compare(env, fn, a, b)
}

func (p Proc) Call(env *Env, args []Object) (Object, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr,
				"\nPanic from proc: %s at %s:%d\nerror: %s\n\n",
				p.Name, p.File, p.Line, err,
			)

			panic(err)
		}
	}()
	return p.Fn(env, args)
}

var _ Callable = (*Fn)(nil)
var _ Callable = Proc{}

func (p Proc) Compare(env *Env, a, b Object) (int, error) {
	return compare(env, p, a, b)
}

func (p Proc) ToString(escape bool) string {
	pkg := p.Package
	if pkg != "" {
		pkg += "."
	}

	file := p.File
	if file == "" {
		file = "<unknown>"
	}

	return fmt.Sprintf("#object[Proc:%s%s %s:%d]", pkg, p.Name, file, p.Line)
}

func (p Proc) Equals(other interface{}) bool {
	switch other := other.(type) {
	case Proc:
		return reflect.ValueOf(p.Fn).Pointer() == reflect.ValueOf(other.Fn).Pointer()
	}
	return false
}

func (p Proc) GetInfo() *ObjectInfo {
	return nil
}

func (p Proc) WithInfo(*ObjectInfo) Object {
	return p
}

func (p Proc) GetType() *Type {
	return TYPE.Proc
}

func (p Proc) Hash() uint32 {
	return HashPtr(reflect.ValueOf(p.Fn).Pointer())
}

func (i InfoHolder) GetInfo() *ObjectInfo {
	return i.info
}

func (m MetaHolder) GetMeta() Map {
	return m.meta
}

func AlterMeta(env *Env, m *MetaHolder, fn *Fn, args []Object) (Map, error) {
	meta := m.meta
	if meta == nil {
		meta = NIL
	}
	fargs := append([]Object{meta}, args...)

	v, err := fn.Call(env, fargs)
	if err != nil {
		return nil, err
	}

	mm, err := AssertMap(env, v, "")
	if err != nil {
		return nil, err
	}
	m.meta = mm
	return m.meta, nil
}

func (sym Symbol) WithMeta(meta Map) (Object, error) {
	res := sym
	m, err := SafeMerge(res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return res, nil
}

func (v *Var) Name() string {
	return v.ns.Name.ToString(false) + "/" + v.name.ToString(false)
}

func (v *Var) ToString(escape bool) string {
	return "#'" + v.Name()
}

func (v *Var) Equals(other interface{}) bool {
	// TODO: revisit this
	return v == other
}

func (v *Var) WithMeta(meta Map) (Object, error) {
	res := *v
	m, err := SafeMerge(res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return &res, nil
}

func (v *Var) ResetMeta(newMeta Map) Map {
	v.meta = newMeta
	return v.meta
}

func (v *Var) AlterMeta(env *Env, fn *Fn, args []Object) (Map, error) {
	return AlterMeta(env, &v.MetaHolder, fn, args)
}

func (v *Var) GetType() *Type {
	return TYPE.Var
}

func (v *Var) Hash() uint32 {
	return HashPtr(uintptr(unsafe.Pointer(v)))
}

func (v *Var) Resolve() Object {
	if v.Value == nil {
		return NIL
	}
	return v.Value
}

func (v *Var) Call(env *Env, args []Object) (Object, error) {
	vl := v.Resolve()
	call, err := AssertCallable(env,
		vl,
		"Var "+v.ToString(false)+" resolves to "+vl.ToString(false)+", which is not a Fn")
	if err != nil {
		return nil, err
	}

	return call.Call(env, args)
}

var _ Callable = (*Var)(nil)

func (v *Var) Deref(env *Env) (Object, error) {
	return v.Resolve(), nil
}

func (n Nil) ToString(escape bool) string {
	return "nil"
}

func (n Nil) Equals(other interface{}) bool {
	switch other.(type) {
	case Nil:
		return true
	default:
		return false
	}
}

func (n Nil) GetType() *Type {
	return TYPE.Nil
}

func (n Nil) Hash() uint32 {
	return 0
}

func (n Nil) Seq() Seq {
	return n
}

func (n Nil) First() Object {
	return NIL
}

func (n Nil) Rest() Seq {
	return NIL
}

func (n Nil) IsEmpty() bool {
	return true
}

func (n Nil) Cons(obj Object) Seq {
	return NewListFrom(obj)
}

func (n Nil) Conj(obj Object) (Conjable, error) {
	return NewListFrom(obj), nil
}

func (n Nil) Without(key Object) Map {
	return n
}

func (n Nil) Count() int {
	return 0
}

func (n Nil) Iter() MapIterator {
	return emptyMapIterator
}

func (n Nil) Merge(other Map) (Map, error) {
	return other, nil
}

func (n Nil) Assoc(key, value Object) (Associative, error) {
	return EmptyArrayMap().Assoc(key, value)
}

func (n Nil) EntryAt(key Object) (*Vector, error) {
	return nil, nil
}

func (n Nil) Get(key Object) (bool, Object) {
	return false, NIL
}

func (n Nil) Disjoin(key Object) Set {
	return n
}

func (n Nil) Keys() Seq {
	return NIL
}

func (n Nil) Vals() Seq {
	return NIL
}

func (rat *Ratio) ToString(escape bool) string {
	return rat.r.String()
}

func (rat *Ratio) Equals(other interface{}) bool {
	return equalsNumbers(rat, other)
}

func (rat *Ratio) GetType() *Type {
	return TYPE.Ratio
}

func (rat *Ratio) Hash() uint32 {
	h, _ := hashGobEncoder(&rat.r)

	return h
}

func (rat *Ratio) Compare(env *Env, other Object) (int, error) {
	n, err := AssertNumber(env, other, "Cannot compare Ratio and "+other.GetType().ToString(false))
	if err != nil {
		return 0, err
	}

	return CompareNumbers(rat, n), nil
}

func MakeBigInt(bi int64) *BigInt {
	return &BigInt{b: *big.NewInt(bi)}
}

func (bi *BigInt) ToString(escape bool) string {
	return bi.b.String() + "N"
}

func (bi *BigInt) Equals(other interface{}) bool {
	return equalsNumbers(bi, other)
}

func (bi *BigInt) GetType() *Type {
	return TYPE.BigInt
}

func (bi *BigInt) Hash() uint32 {
	h, _ := hashGobEncoder(&bi.b)
	return h
}

func (bi *BigInt) Compare(env *Env, other Object) (int, error) {
	n, err := AssertNumber(env, other, "Cannot compare BigInt and "+other.GetType().ToString(false))
	if err != nil {
		return 0, err
	}
	return CompareNumbers(bi, n), nil
}

func (bf *BigFloat) ToString(escape bool) string {
	return bf.b.Text('g', -1) + "M"
}

func (bf *BigFloat) Equals(other interface{}) bool {
	return equalsNumbers(bf, other)
}

func (bf *BigFloat) GetType() *Type {
	return TYPE.BigFloat
}

func (bf *BigFloat) Hash() uint32 {
	h, _ := hashGobEncoder(&bf.b)
	return h
}

func (bf *BigFloat) Compare(env *Env, other Object) (int, error) {
	n, err := AssertNumber(env, other, "Cannot compare BigFloat and "+other.GetType().ToString(false))
	if err != nil {
		return 0, err
	}

	return CompareNumbers(bf, n), nil
}

func (c Char) ToString(escape bool) string {
	if escape {
		return escapeRune(c.Ch)
	}
	return string(c.Ch)
}

func (c Char) Equals(other interface{}) bool {
	switch other := other.(type) {
	case Char:
		return c.Ch == other.Ch
	default:
		return false
	}
}

func (c Char) GetType() *Type {
	return TYPE.Char
}

func (c Char) Native() interface{} {
	return c.Ch
}

func (c Char) Hash() uint32 {
	h := getHash()
	h.Write([]byte(string(c.Ch)))
	return h.Sum32()
}

func (c Char) Compare(env *Env, other Object) (int, error) {
	c2, err := AssertChar(env, other, "Cannot compare Char and "+other.GetType().ToString(false))
	if err != nil {
		return 0, err
	}
	if c.Ch < c2.Ch {
		return -1, nil
	}
	if c2.Ch < c.Ch {
		return 1, nil
	}
	return 0, nil
}

func MakeBoolean(b bool) Boolean {
	return Boolean{B: b}
}

func MakeTime(t time.Time) Time {
	return Time{T: t}
}

func MakeDouble(d float64) Double {
	return Double{D: d}
}

func (d Double) ToString(escape bool) string {
	return fmt.Sprintf("%g", d.D)
}

func (d Double) Equals(other interface{}) bool {
	return equalsNumbers(d, other)
}

func (d Double) GetType() *Type {
	return TYPE.Double
}

func (d Double) Native() interface{} {
	return d.D
}

func (d Double) Hash() uint32 {
	h := getHash()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, math.Float64bits(d.D))
	h.Write(b)
	return h.Sum32()
}

func (d Double) Compare(env *Env, other Object) (int, error) {
	n, err := AssertNumber(env, other, "Cannot compare Double and "+other.GetType().ToString(false))
	if err != nil {
		return 0, err
	}
	return CompareNumbers(d, n), nil
}

func (i Int) ToString(escape bool) string {
	return fmt.Sprintf("%d", i.I)
}

func MakeInt(i int) Int {
	return Int{I: i}
}

func (i Int) Equals(other interface{}) bool {
	return equalsNumbers(i, other)
}

func (i Int) GetType() *Type {
	return TYPE.Int
}

func (i Int) Native() interface{} {
	return i.I
}

func (i Int) Hash() uint32 {
	h := getHash()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i.I))
	h.Write(b)
	return h.Sum32()
}

func (i Int) Compare(env *Env, other Object) (int, error) {
	n, err := AssertNumber(env, other, "Cannot compare Int and "+other.GetType().ToString(false))
	if err != nil {
		return 0, err
	}
	return CompareNumbers(i, n), nil
}

func (b Boolean) ToString(escape bool) string {
	return fmt.Sprintf("%t", b.B)
}

func (b Boolean) Equals(other interface{}) bool {
	switch other := other.(type) {
	case Boolean:
		return b.B == other.B
	default:
		return false
	}
}

func (b Boolean) GetType() *Type {
	return TYPE.Boolean
}

func (b Boolean) Native() interface{} {
	return b.B
}

func (b Boolean) Hash() uint32 {
	h := getHash()
	var bs = make([]byte, 1)
	if b.B {
		bs[0] = 1
	} else {
		bs[0] = 0
	}
	h.Write(bs)
	return h.Sum32()
}

func (b Boolean) Compare(env *Env, other Object) (int, error) {
	b2, err := AssertBoolean(env, other, "Cannot compare Boolean and "+other.GetType().ToString(false))
	if err != nil {
		return 0, err
	}
	if b.B == b2.B {
		return 0, nil
	}
	if b.B {
		return 1, nil
	}
	return -1, nil
}

func (t Time) ToString(escape bool) string {
	return t.T.String()
}

func (t Time) Equals(other interface{}) bool {
	switch other := other.(type) {
	case Time:
		return t.T.Equal(other.T)
	default:
		return false
	}
}

func (t Time) GetType() *Type {
	return TYPE.Time
}

func (t Time) Native() interface{} {
	return t.T
}

func (t Time) Hash() uint32 {
	h, _ := hashGobEncoder(t.T)
	return h
}

func (t Time) Compare(env *Env, other Object) (int, error) {
	t2, err := AssertTime(env, other, "Cannot compare Time and "+other.GetType().ToString(false))
	if err != nil {
		return 0, err
	}
	if t.T.Equal(t2.T) {
		return 0, nil
	}
	if t2.T.Before(t.T) {
		return 1, nil
	}
	return -1, nil
}

func (k Keyword) ToString(escape bool) string {
	if k.ns != nil {
		return ":" + *k.ns + "/" + *k.name
	}
	return ":" + *k.name
}

func (k Keyword) Name() string {
	return *k.name
}

func (k Keyword) Namespace() string {
	if k.ns != nil {
		return *k.ns
	}
	return ""
}

func (k Keyword) Equals(other interface{}) bool {
	switch other := other.(type) {
	case Keyword:
		return k.ns == other.ns && k.name == other.name
	default:
		return false
	}
}

func (k Keyword) GetType() *Type {
	return TYPE.Keyword
}

func (k Keyword) Hash() uint32 {
	return k.hash
}

func (k Keyword) Compare(env *Env, other Object) (int, error) {
	k2, err := AssertKeyword(env, other, "Cannot compare Keyword and "+other.GetType().ToString(false))
	if err != nil {
		return 0, err
	}
	return strings.Compare(k.ToString(false), k2.ToString(false)), nil
}

func (k Keyword) Call(env *Env, args []Object) (Object, error) {
	return getMap(env, k, args)
}

var _ Callable = Keyword{}

func MakeRegex(r *regexp.Regexp) *Regex {
	return &Regex{R: r}
}

func (rx *Regex) ToString(escape bool) string {
	if escape {
		return "#\"" + rx.R.String() + "\""
	}
	return rx.R.String()
}

func (rx *Regex) Print(w io.Writer, printReadably bool) {
	fmt.Fprint(w, rx.ToString(true))
}

func (rx *Regex) Equals(other interface{}) bool {
	switch other := other.(type) {
	case *Regex:
		return rx.R == other.R
	default:
		return false
	}
}

func (rx *Regex) GetType() *Type {
	return TYPE.Regex
}

func (rx *Regex) Hash() uint32 {
	return HashPtr(uintptr(unsafe.Pointer(rx.R)))
}

func (s Symbol) ToString(escape bool) string {
	if s.ns != nil {
		return *s.ns + "/" + *s.name
	}
	return *s.name
}

func (s Symbol) Name() string {
	return *s.name
}

func (s Symbol) Namespace() string {
	if s.ns != nil {
		return *s.ns
	}
	return ""
}

func (s Symbol) Equals(other interface{}) bool {
	switch other := other.(type) {
	case Symbol:
		return s.ns == other.ns && s.name == other.name
	default:
		return false
	}
}

func (s Symbol) GetType() *Type {
	return TYPE.Symbol
}

func (s Symbol) Hash() uint32 {
	return hashSymbol(s.ns, s.name) + 0x9e3779b9
}

func (s Symbol) Compare(env *Env, other Object) (int, error) {
	s2, err := AssertSymbol(env, other, "Cannot compare Symbol and "+other.GetType().ToString(false))
	if err != nil {
		return 0, err
	}
	return strings.Compare(s.ToString(false), s2.ToString(false)), nil
}

func (s Symbol) Call(env *Env, args []Object) (Object, error) {
	return getMap(env, s, args)
}

var _ Callable = Symbol{}

func (s String) ToString(escape bool) string {
	if escape {
		return escapeString(s.S)
	}
	return s.S
}

func MakeString(s string) String {
	return String{S: s}
}

func MakeStringVector(ss []string) *Vector {
	res := EmptyVector()
	for _, s := range ss {
		res, _ = res.Conjoin(MakeString(s))
	}
	return res
}

func (s String) Equals(other interface{}) bool {
	switch other := other.(type) {
	case String:
		return s.S == other.S
	default:
		return false
	}
}

func (s String) GetType() *Type {
	return TYPE.String
}

func (s String) Native() interface{} {
	return s.S
}

func (s String) Hash() uint32 {
	h := getHash()
	h.Write([]byte(s.S))
	return h.Sum32()
}

func (s String) Count() int {
	return utf8.RuneCountInString(s.S)
}

func (s String) Seq() Seq {
	runes := make([]Object, 0, len(s.S))
	for _, r := range s.S {
		runes = append(runes, Char{Ch: r})
	}
	return &ArraySeq{arr: runes}
}

func (s String) Nth(i int) Object {
	if i < 0 {
		panic(StubNewError(fmt.Sprintf("Negative index: %d", i)))
	}
	j, r := 0, 't'
	for j, r = range s.S {
		if i == j {
			return Char{Ch: r}
		}
	}
	panic(StubNewError(fmt.Sprintf("Index %d exceeds string's length %d", i, j+1)))
}

func (s String) TryNth(i int, d Object) Object {
	if i < 0 {
		return d
	}
	for j, r := range s.S {
		if i == j {
			return Char{Ch: r}
		}
	}
	return d
}

func (s String) Compare(env *Env, other Object) (int, error) {
	s2, err := AssertString(env, other, "Cannot compare String and "+other.GetType().ToString(false))
	if err != nil {
		return 0, err
	}

	return strings.Compare(s.S, s2.S), nil
}

func IsSymbol(obj Object) bool {
	switch obj.(type) {
	case Symbol:
		return true
	default:
		return false
	}
}

func IsVector(obj Object) bool {
	switch obj.(type) {
	case *Vector:
		return true
	default:
		return false
	}
}

func IsSeq(obj Object) bool {
	switch obj.(type) {
	case Seq:
		return true
	default:
		return false
	}
}

func (x *Type) WithInfo(info *ObjectInfo) Object {
	return x
}

func (x RecurBindings) WithInfo(info *ObjectInfo) Object {
	return x
}

func IsEqualOrImplements(abstractType *Type, concreteType *Type) bool {
	if abstractType.reflectType.Kind() == reflect.Interface {
		return concreteType.reflectType.Implements(abstractType.reflectType)
	} else {
		return concreteType.reflectType == abstractType.reflectType
	}
}

func IsInstance(t *Type, obj Object) bool {
	if obj.Equals(NIL) {
		return false
	}
	return IsEqualOrImplements(t, obj.GetType())
}

var specialSymbols = make(map[*string]struct{})

func init() {
	specialSymbols[criticalSymbols._if.name] = struct{}{}
	specialSymbols[criticalSymbols.quote.name] = struct{}{}
	specialSymbols[criticalSymbols.fn_.name] = struct{}{}
	specialSymbols[criticalSymbols.let_.name] = struct{}{}
	specialSymbols[criticalSymbols.letfn_.name] = struct{}{}
	specialSymbols[criticalSymbols.loop_.name] = struct{}{}
	specialSymbols[criticalSymbols.recur.name] = struct{}{}
	specialSymbols[criticalSymbols.setMacro_.name] = struct{}{}
	specialSymbols[criticalSymbols.def.name] = struct{}{}
	specialSymbols[criticalSymbols.defLinter.name] = struct{}{}
	specialSymbols[criticalSymbols._var.name] = struct{}{}
	specialSymbols[criticalSymbols.do.name] = struct{}{}
	specialSymbols[criticalSymbols.throw.name] = struct{}{}
	specialSymbols[criticalSymbols.try.name] = struct{}{}
	specialSymbols[criticalSymbols.catch.name] = struct{}{}
	specialSymbols[criticalSymbols.finally.name] = struct{}{}
}

func IsSpecialSymbol(obj Object) bool {
	switch obj := obj.(type) {
	case Symbol:
		if obj.ns != nil {
			return false
		}

		_, found := specialSymbols[obj.name]
		return found
	default:
		return false
	}
}

func MakeMeta(arglists Seq, docstring string, added string) *ArrayMap {
	res := EmptyArrayMap()
	if arglists != nil {
		res.Add(criticalKeywords.arglist, arglists)
	}
	res.Add(criticalKeywords.doc, String{S: docstring})
	res.Add(criticalKeywords.added, String{S: added})
	return res
}

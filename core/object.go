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
		ToString(env *Env, escape bool) (string, error)
		GetInfo() *ObjectInfo
		WithInfo(*ObjectInfo) Object
		GetType() *Type
		Hash(env *Env) (uint32, error)
	}
	Equality interface {
		Equals(env *Env, other interface{}) bool
	}
	Conjable interface {
		Object
		Conj(env *Env, obj Object) (Conjable, error)
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
		WithMeta(*Env, Map) (Object, error)
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
		Nth(env *Env, i int) (Object, error)
		TryNth(env *Env, i int, d Object) (Object, error)
	}
	IndexCounted interface {
		Indexed
		Counted
	}
	Stack interface {
		Object
		Peek(env *Env) (Object, error)
		Pop(env *Env) (Stack, error)
	}
	Gettable interface {
		Get(env *Env, key Object) (bool, Object, error)
	}
	Associative interface {
		Object
		Conjable
		Gettable
		EntryAt(env *Env, key Object) (*Vector, error)
		Assoc(env *Env, key, val Object) (Associative, error)
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
		Pprint(env *Env, writer io.Writer, indent int) (int, error)
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
		filename    string
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
	}
	Keyword struct {
		InfoHolder
		ns   string
		name string
		hash uint32
	}
	Symbol struct {
		InfoHolder
		MetaHolder
		ns   string
		name string
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

		code           *Code
		importedUpvals []*NamedPair
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
		ReflectType    *Type
		ReflectValue   *Type
		NodeSeq        *Type
		ParseError     *Type
		NamedPair      *Type
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
	_ Object = Time{}

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
	if pos.filename == "" {
		return "<file>"
	}
	return pos.filename
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

func hashSymbol(ns, name string) uint32 {
	h := getHash()
	if ns != "" {
		h.Write([]byte(ns))
	}
	h.Write([]byte("/" + name))
	return h.Sum32()
}

func AssembleSymbol(ns, name string) Symbol {
	if ns == "" {
		if strings.ContainsRune(name, '/') {
			panic("bad symbol name")
		}

		return Symbol{
			name: name,
		}
	}
	return Symbol{
		ns:   ns,
		name: name,
	}
}

func MakeSymbol(nsname string) Symbol {
	index := strings.IndexRune(nsname, '/')
	if index == -1 || nsname == "/" {
		return Symbol{
			name: nsname,
		}
	}
	return Symbol{
		ns:   nsname[0:index],
		name: nsname[index+1:],
	}
}

func MakeSymbolWithMeta(nsname string, m Map) Symbol {
	index := strings.IndexRune(nsname, '/')
	var sym Symbol
	if index == -1 || nsname == "/" {
		sym = Symbol{
			name: nsname,
		}
	} else {
		sym = Symbol{
			ns:   nsname[0:index],
			name: nsname[index+1:],
		}
	}

	sym.meta = m

	return sym
}

func MakeTaggedSymbol(nsname string, tag Symbol) Symbol {
	var sym Symbol

	index := strings.IndexRune(nsname, '/')
	if index == -1 || nsname == "/" {
		sym = Symbol{
			name: nsname,
		}
	} else {
		sym = Symbol{
			ns:   nsname[0:index],
			name: nsname[index+1:],
		}
	}

	m := EmptyArrayMap()
	m.AddEqu(criticalKeywords.tag, tag)

	sym.meta = m

	return sym
}

type BySymbolName []Symbol

func (s BySymbolName) Len() int {
	return len(s)
}
func (s BySymbolName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s BySymbolName) Less(i, j int) bool {
	return s[i].String() < s[j].String()
}

const KeywordHashMask uint32 = 0x7334c790

func MakeKeyword(nsname string) Keyword {
	index := strings.IndexRune(nsname, '/')
	if index == -1 || nsname == "/" {
		name := nsname
		return Keyword{
			name: name,
			hash: hashSymbol("", name) ^ KeywordHashMask,
		}
	}
	ns := nsname[0:index]
	name := nsname[index+1:]
	return Keyword{
		ns:   ns,
		name: name,
		hash: hashSymbol(ns, name) ^ KeywordHashMask,
	}
}

func ErrorArity(env *Env, n int) error {
	return env.NewError(fmt.Sprintf("Wrong number of args (%d)", n))
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
	return env.NewError(fmt.Sprintf("Wrong number of args (%d); expects %s", n, rangeString(min, max)))
}

func ReturnArityMinMax(env *Env, n, min, max int) error {
	return env.NewError(fmt.Sprintf("Wrong number of args (%d); expects %s", n, rangeString(min, max)))
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
		ok, v, err := m.Get(env, k)
		if err != nil {
			return nil, err
		}

		if ok {
			return v, nil
		}
	default:
		return nil, env.NewArgTypeError(1, args[0], "Map")
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

func (a *Atom) ToString(env *Env, escape bool) (string, error) {
	v, err := a.value.ToString(env, escape)
	if err != nil {
		return "", err
	}
	return "#object[Atom {:val " + v + "}]", nil
}

func (a *Atom) Equals(env *Env, other interface{}) bool {
	return a == other
}

func (a *Atom) GetInfo() *ObjectInfo {
	return nil
}

func (a *Atom) GetType() *Type {
	return TYPE.Atom
}

func (a *Atom) Hash(env *Env) (uint32, error) {
	return HashPtr(uintptr(unsafe.Pointer(a))), nil
}

func (a *Atom) WithInfo(info *ObjectInfo) Object {
	return a
}

func (a *Atom) WithMeta(env *Env, meta Map) (Object, error) {
	res := *a
	m, err := SafeMerge(env, res.meta, meta)
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

func (d *Delay) ToString(env *Env, escape bool) (string, error) {
	return "#object[Delay]", nil
}

func (d *Delay) Equals(env *Env, other interface{}) bool {
	return d == other
}

func (d *Delay) GetInfo() *ObjectInfo {
	return nil
}

func (d *Delay) GetType() *Type {
	return TYPE.Delay
}

func (d *Delay) Hash(env *Env) (uint32, error) {
	return HashPtr(uintptr(unsafe.Pointer(d))), nil
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

func (t *Type) ToString(env *Env, escape bool) (string, error) {
	return t.name, nil
}

func (t *Type) Name() string {
	return t.name
}

func (t *Type) Equals(env *Env, other interface{}) bool {
	return t == other
}

func (t *Type) GetInfo() *ObjectInfo {
	return nil
}

func (t *Type) GetType() *Type {
	return TYPE.Type
}

func (t *Type) Hash(env *Env) (uint32, error) {
	return HashPtr(uintptr(unsafe.Pointer(t))), nil
}

func (rb RecurBindings) ToString(env *Env, escape bool) (string, error) {
	return "#object[RecurBindings]", nil
}

func (rb RecurBindings) Equals(env *Env, other interface{}) bool {
	return false
}

func (rb RecurBindings) GetInfo() *ObjectInfo {
	return nil
}

func (rb RecurBindings) GetType() *Type {
	return TYPE.RecurBindings
}

func (rb RecurBindings) Hash(env *Env) (uint32, error) {
	return 0, nil
}

func (fn *Fn) ToString(env *Env, escape bool) (string, error) {
	return "#object[Fn]", nil
}

func (fn *Fn) String() string {
	if fn.code != nil {
		return fmt.Sprintf("<fn bc @ %s:%d>", fn.code.filename, fn.code.lineForIp(0))
	} else {
		pos := fn.fnExpr.Pos()
		return fmt.Sprintf("<fn tree @ %s:%d>", pos.filename, pos.startLine)
	}
}

func (fn *Fn) Equals(env *Env, other interface{}) bool {
	switch other := other.(type) {
	case *Fn:
		return fn == other
	default:
		return false
	}
}

func (fn *Fn) WithMeta(env *Env, meta Map) (Object, error) {
	res := *fn
	m, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return &res, nil
}

func (fn *Fn) GetType() *Type {
	return TYPE.Fn
}

func (fn *Fn) Hash(env *Env) (uint32, error) {
	return HashPtr(uintptr(unsafe.Pointer(fn))), nil
}

func (fn *Fn) Call(env *Env, args []Object) (Object, error) {
	obj, err := env.Engine.RunWithArgs(env, fn, args)
	if err != nil {
		return nil, err
	}

	return obj, nil
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
	ret, err := p.Fn(env, args)
	if err != nil {
		err = env.populateStackTrace(err)
	}

	return ret, err
}

var _ Callable = (*Fn)(nil)
var _ Callable = Proc{}

func (p Proc) Compare(env *Env, a, b Object) (int, error) {
	return compare(env, p, a, b)
}

func (p Proc) ToString(env *Env, escape bool) (string, error) {
	pkg := p.Package
	if pkg != "" {
		pkg += "."
	}

	file := p.File
	if file == "" {
		file = "<unknown>"
	}

	return fmt.Sprintf("#object[Proc:%s%s %s:%d]", pkg, p.Name, file, p.Line), nil
}

func (p Proc) Equals(env *Env, other interface{}) bool {
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

func (p Proc) Hash(env *Env) (uint32, error) {
	return HashPtr(reflect.ValueOf(p.Fn).Pointer()), nil
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

func (sym Symbol) WithMeta(env *Env, meta Map) (Object, error) {
	res := sym
	m, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return res, nil
}

func (v *Var) Name() string {
	return v.ns.Name.String() + "/" + v.name.String()
}

func (v *Var) ToString(env *Env, escape bool) (string, error) {
	return "#'" + v.Name(), nil
}

func (v *Var) String() string {
	return "#'" + v.Name()
}

func (v *Var) Equals(env *Env, other interface{}) bool {
	// TODO: revisit this
	return v == other
}

func (v *Var) WithMeta(env *Env, meta Map) (Object, error) {
	res := *v
	m, err := SafeMerge(env, res.meta, meta)
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

func (v *Var) Hash(env *Env) (uint32, error) {
	return HashPtr(uintptr(unsafe.Pointer(v))), nil
}

func (v *Var) Resolve() Object {
	if v.Value == nil {
		return NIL
	}
	return v.Value
}

func (v *Var) Call(env *Env, args []Object) (Object, error) {
	vl := v.Resolve()
	vs, err := v.ToString(env, false)
	if err != nil {
		return nil, err
	}

	vls, err := v.ToString(env, false)
	if err != nil {
		return nil, err
	}

	call, err := AssertCallable(env,
		vl,
		"Var "+vs+" resolves to "+vls+", which is not a Fn")
	if err != nil {
		return nil, err
	}

	return call.Call(env, args)
}

var _ Callable = (*Var)(nil)

func (v *Var) Deref(env *Env) (Object, error) {
	return v.Resolve(), nil
}

func (n Nil) ToString(env *Env, escape bool) (string, error) {
	return "nil", nil
}

func (n Nil) Equals(env *Env, other interface{}) bool {
	switch other.(type) {
	case Nil:
		return true
	default:
		return false
	}
}

func (n Nil) GetEqu(key Equ) (bool, Object) {
	return false, NIL
}

func (n Nil) GetType() *Type {
	return TYPE.Nil
}

func (n Nil) Hash(env *Env) (uint32, error) {
	return 0, nil
}

func (n Nil) Seq() Seq {
	return n
}

func (n Nil) First(env *Env) (Object, error) {
	return NIL, nil
}

func (n Nil) Rest(env *Env) (Seq, error) {
	return NIL, nil
}

func (n Nil) IsEmpty(env *Env) (bool, error) {
	return true, nil
}

func (n Nil) Cons(obj Object) Seq {
	return NewListFrom(obj)
}

func (n Nil) Conj(env *Env, obj Object) (Conjable, error) {
	return NewListFrom(obj), nil
}

func (n Nil) Without(env *Env, key Object) (Map, error) {
	return n, nil
}

func (n Nil) Count() int {
	return 0
}

func (n Nil) Iter() MapIterator {
	return emptyMapIterator
}

func (n Nil) Merge(env *Env, other Map) (Map, error) {
	return other, nil
}

func (n Nil) Assoc(env *Env, key, value Object) (Associative, error) {
	return EmptyArrayMap().Assoc(env, key, value)
}

func (n Nil) EntryAt(env *Env, key Object) (*Vector, error) {
	return nil, nil
}

func (n Nil) Get(env *Env, key Object) (bool, Object, error) {
	return false, NIL, nil
}

func (n Nil) Disjoin(env *Env, key Object) (Set, error) {
	return n, nil
}

func (n Nil) SetIter() SetIter {
	return emptySetIterator
}

func (n Nil) Has(key Equ) bool {
	return false
}

func (n Nil) Keys() Seq {
	return NIL
}

func (n Nil) Vals() Seq {
	return NIL
}

func MakeRatio(x, y *big.Int) *Ratio {
	r := big.NewRat(x.Int64(), y.Int64())
	return &Ratio{r: *r}
}

func (rat *Ratio) ToString(env *Env, escape bool) (string, error) {
	return rat.r.String(), nil
}

func (rat *Ratio) Equals(env *Env, other interface{}) bool {
	return equalsNumbers(rat, other)
}

func (rat *Ratio) GetType() *Type {
	return TYPE.Ratio
}

func (rat *Ratio) Hash(env *Env) (uint32, error) {
	return hashGobEncoder(&rat.r)
}

func (rat *Ratio) Compare(env *Env, other Object) (int, error) {
	os, err := other.GetType().ToString(env, false)
	if err != nil {
		return 0, err
	}
	n, err := AssertNumber(env, other, "Cannot compare Ratio and "+os)
	if err != nil {
		return 0, err
	}

	return CompareNumbers(rat, n), nil
}

func MakeBigInt(bi int64) *BigInt {
	return &BigInt{b: *big.NewInt(bi)}
}

func MakeBigIntFrom(bi *big.Int) *BigInt {
	return &BigInt{b: *bi}
}

func (bi *BigInt) ToString(env *Env, escape bool) (string, error) {
	return bi.b.String() + "N", nil
}

func (bi *BigInt) Equals(env *Env, other interface{}) bool {
	return equalsNumbers(bi, other)
}

func (bi *BigInt) GetType() *Type {
	return TYPE.BigInt
}

func (bi *BigInt) Hash(env *Env) (uint32, error) {
	return hashGobEncoder(&bi.b)
}

func (bi *BigInt) Compare(env *Env, other Object) (int, error) {
	os, err := other.GetType().ToString(env, false)
	if err != nil {
		return 0, err
	}
	n, err := AssertNumber(env, other, "Cannot compare BigInt and "+os)
	if err != nil {
		return 0, err
	}
	return CompareNumbers(bi, n), nil
}

func MakeBigFloatFrom(bi *big.Float) *BigFloat {
	return &BigFloat{b: *bi}
}

func (bf *BigFloat) ToString(env *Env, escape bool) (string, error) {
	return bf.b.Text('g', -1) + "M", nil
}

func (bf *BigFloat) Equals(env *Env, other interface{}) bool {
	return equalsNumbers(bf, other)
}

func (bf *BigFloat) GetType() *Type {
	return TYPE.BigFloat
}

func (bf *BigFloat) Hash(env *Env) (uint32, error) {
	return hashGobEncoder(&bf.b)
}

func (bf *BigFloat) Compare(env *Env, other Object) (int, error) {
	os, err := other.GetType().ToString(env, false)
	if err != nil {
		return 0, err
	}
	n, err := AssertNumber(env, other, "Cannot compare BigFloat and "+os)
	if err != nil {
		return 0, err
	}

	return CompareNumbers(bf, n), nil
}

func (c Char) ToString(env *Env, escape bool) (string, error) {
	if escape {
		return escapeRune(c.Ch), nil
	}
	return string(c.Ch), nil
}

func (c Char) Equals(env *Env, other interface{}) bool {
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

func (c Char) Hash(env *Env) (uint32, error) {
	h := getHash()
	h.Write([]byte(string(c.Ch)))
	return h.Sum32(), nil
}

func (c Char) Compare(env *Env, other Object) (int, error) {
	os, err := other.GetType().ToString(env, false)
	if err != nil {
		return 0, err
	}

	c2, err := AssertChar(env, other, "Cannot compare Char and "+os)
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

func (d Double) ToString(env *Env, escape bool) (string, error) {
	return fmt.Sprintf("%g", d.D), nil
}

func (d Double) Equals(env *Env, other interface{}) bool {
	return equalsNumbers(d, other)
}

func (d Double) GetType() *Type {
	return TYPE.Double
}

func (d Double) Native() interface{} {
	return d.D
}

func (d Double) Hash(env *Env) (uint32, error) {
	h := getHash()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, math.Float64bits(d.D))
	h.Write(b)
	return h.Sum32(), nil
}

func (d Double) Compare(env *Env, other Object) (int, error) {
	os, err := other.GetType().ToString(env, false)
	if err != nil {
		return 0, err
	}

	n, err := AssertNumber(env, other, "Cannot compare Double and "+os)
	if err != nil {
		return 0, err
	}
	return CompareNumbers(d, n), nil
}

func (i Int) ToString(env *Env, escape bool) (string, error) {
	return fmt.Sprintf("%d", i.I), nil
}

func MakeInt(i int) Int {
	return Int{I: i}
}

func (i Int) Equals(env *Env, other interface{}) bool {
	return equalsNumbers(i, other)
}

func (i Int) GetType() *Type {
	return TYPE.Int
}

func (i Int) Native() interface{} {
	return i.I
}

func (i Int) Hash(env *Env) (uint32, error) {
	h := getHash()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i.I))
	h.Write(b)
	return h.Sum32(), nil
}

func (i Int) Compare(env *Env, other Object) (int, error) {
	os, err := other.GetType().ToString(env, false)
	if err != nil {
		return 0, err
	}

	n, err := AssertNumber(env, other, "Cannot compare Int and "+os)
	if err != nil {
		return 0, err
	}
	return CompareNumbers(i, n), nil
}

func (b Boolean) ToString(env *Env, escape bool) (string, error) {
	return fmt.Sprintf("%t", b.B), nil
}

func (b Boolean) Equals(env *Env, other interface{}) bool {
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

func (b Boolean) Hash(env *Env) (uint32, error) {
	h := getHash()
	var bs = make([]byte, 1)
	if b.B {
		bs[0] = 1
	} else {
		bs[0] = 0
	}
	h.Write(bs)
	return h.Sum32(), nil
}

func (b Boolean) Compare(env *Env, other Object) (int, error) {
	os, err := other.GetType().ToString(env, false)
	if err != nil {
		return 0, err
	}

	b2, err := AssertBoolean(env, other, "Cannot compare Boolean and "+os)
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

func (t Time) ToString(env *Env, escape bool) (string, error) {
	return t.T.String(), nil
}

func (t Time) Equals(env *Env, other interface{}) bool {
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

func (t Time) Hash(env *Env) (uint32, error) {
	return hashGobEncoder(t.T)
}

func (t Time) Compare(env *Env, other Object) (int, error) {
	os, err := other.GetType().ToString(env, false)
	if err != nil {
		return 0, err
	}

	t2, err := AssertTime(env, other, "Cannot compare Time and "+os)
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

func (k Keyword) ToString(env *Env, escape bool) (string, error) {
	if k.ns != "" {
		return ":" + k.ns + "/" + k.name, nil
	}
	return ":" + k.name, nil
}

func (k Keyword) String() string {
	if k.ns != "" {
		return ":" + k.ns + "/" + k.name
	}
	return ":" + k.name
}

func (k Keyword) RawString() string {
	if k.ns != "" {
		return k.ns + "/" + k.name
	}
	return k.name
}

func (k Keyword) Name() string {
	return k.name
}

func (k Keyword) Namespace() string {
	if k.ns != "" {
		return k.ns
	}
	return ""
}

func (k Keyword) Equals(env *Env, other interface{}) bool {
	switch other := other.(type) {
	case Keyword:
		return k.ns == other.ns && k.name == other.name
	default:
		return false
	}
}

func (s Keyword) Is(other Object) bool {
	switch other := other.(type) {
	case Keyword:
		return s.ns == other.ns && s.name == other.name
	default:
		return false
	}
}

func (k Keyword) GetType() *Type {
	return TYPE.Keyword
}

func (k Keyword) Hash(env *Env) (uint32, error) {
	return k.hash, nil
}

func (k Keyword) IsHash() uint32 {
	return k.hash
}

func (k Keyword) Compare(env *Env, other Object) (int, error) {
	os, err := other.GetType().ToString(env, false)
	if err != nil {
		return 0, err
	}

	k2, err := AssertKeyword(env, other, "Cannot compare Keyword and "+os)
	if err != nil {
		return 0, err
	}

	ks, err := k.ToString(env, false)
	if err != nil {
		return 0, err
	}
	k2s, err := k2.ToString(env, false)
	if err != nil {
		return 0, err
	}
	return strings.Compare(ks, k2s), nil
}

func (k Keyword) Call(env *Env, args []Object) (Object, error) {
	return getMap(env, k, args)
}

var _ Callable = Keyword{}

func MakeRegex(r *regexp.Regexp) *Regex {
	return &Regex{R: r}
}

func (rx *Regex) ToString(env *Env, escape bool) (string, error) {
	if escape {
		return "#\"" + rx.R.String() + "\"", nil
	}
	return rx.R.String(), nil
}

func (rx *Regex) Print(w io.Writer, printReadably bool) {
	fmt.Fprint(w, "#\""+rx.R.String()+"\"")
}

func (rx *Regex) Equals(env *Env, other interface{}) bool {
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

func (rx *Regex) Hash(env *Env) (uint32, error) {
	return HashPtr(uintptr(unsafe.Pointer(rx.R))), nil
}

func (s Symbol) ToString(env *Env, escape bool) (string, error) {
	if s.ns != "" {
		return s.ns + "/" + s.name, nil
	}
	return s.name, nil
}

func (s Symbol) String() string {
	if s.ns != "" {
		return s.ns + "/" + s.name
	}
	return s.name
}

func (s Symbol) Name() string {
	return s.name
}

func (s Symbol) Namespace() string {
	if s.ns != "" {
		return s.ns
	}
	return ""
}

func (s Symbol) Equals(env *Env, other interface{}) bool {
	switch other := other.(type) {
	case Symbol:
		return s.ns == other.ns && s.name == other.name
	default:
		return false
	}
}

func (s Symbol) Is(other Object) bool {
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

func (s Symbol) Hash(env *Env) (uint32, error) {
	return s.IsHash(), nil
}

func (s Symbol) IsHash() uint32 {
	return hashSymbol(s.ns, s.name) + 0x9e3779b9
}

func (s Symbol) Compare(env *Env, other Object) (int, error) {
	os, err := other.GetType().ToString(env, false)
	if err != nil {
		return 0, err
	}

	s2, err := AssertSymbol(env, other, "Cannot compare Symbol and "+os)
	if err != nil {
		return 0, err
	}

	ks, err := s.ToString(env, false)
	if err != nil {
		return 0, err
	}

	k2s, err := s2.ToString(env, false)
	if err != nil {
		return 0, err
	}
	return strings.Compare(ks, k2s), nil
}

func (s Symbol) Call(env *Env, args []Object) (Object, error) {
	return getMap(env, s, args)
}

var _ Callable = Symbol{}

func (s String) ToString(env *Env, escape bool) (string, error) {
	if escape {
		return escapeString(s.S), nil
	}
	return s.S, nil
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

func (s String) Equals(env *Env, other interface{}) bool {
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

func (s String) Hash(env *Env) (uint32, error) {
	h := getHash()
	h.Write([]byte(s.S))
	return h.Sum32(), nil
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

func (s String) Nth(env *Env, i int) (Object, error) {
	if i < 0 {
		return nil, env.NewError(fmt.Sprintf("Negative index: %d", i))
	}
	j := 0
	var r rune

	for j, r = range s.S {
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
	for j, r := range s.S {
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

func IsInstance(env *Env, t *Type, obj Object) bool {
	if obj.Equals(env, NIL) {
		return false
	}
	return IsEqualOrImplements(t, obj.GetType())
}

var specialSymbols = make(map[string]struct{})

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
		if obj.ns != "" {
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
		res.AddEqu(criticalKeywords.arglist, arglists)
	}
	res.AddEqu(criticalKeywords.doc, String{S: docstring})
	res.AddEqu(criticalKeywords.added, String{S: added})
	return res
}

func (p Position) String() string {
	if p.filename == "" {
		return "<unknown>:-1"
	}

	if p.startLine != p.endLine {
		return fmt.Sprintf("%s:%d-%d", p.Filename(), p.startLine, p.endLine)
	} else {
		return fmt.Sprintf("%s:%d", p.Filename(), p.startLine)
	}
}

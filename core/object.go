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
	"reflect"
	"regexp"
	"strconv"
	"time"
	"unsafe"
)

// interfaces
type (
	Equality interface {
		Equals(env *Env, other interface{}) bool
	}
	Conjable interface {
		Conj(env *Env, obj any) (Conjable, error)
	}
	Counted interface {
		Count() int
	}
	Error interface {
		error
		Message() any
	}
	Meta interface {
		GetMeta() Map
		WithMeta(*Env, Map) (any, error)
	}
	Ref interface {
		AlterMeta(env *Env, fn *Fn, args []any) (Map, error)
		ResetMeta(m Map) Map
	}
	Sequential interface {
		sequential()
	}
	Comparable interface {
		Compare(env *Env, other any) (int, error)
	}
	Comparator interface {
		Compare(env *Env, a, b any) (int, error)
	}
	Indexed interface {
		Nth(env *Env, i int) (any, error)
		TryNth(env *Env, i int, d any) (any, error)
	}
	IndexCounted interface {
		Indexed
		Counted
	}
	Stack interface {
		Peek(env *Env) (any, error)
		Pop(env *Env) (Stack, error)
	}
	Gettable interface {
		Get(env *Env, key any) (bool, any, error)
	}
	Associative interface {
		Conjable
		Gettable
		EntryAt(env *Env, key any) (*Vector, error)
		Assoc(env *Env, key, val any) (Associative, error)
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
		Counted
		Seqable
		Empty() Collection
	}
	Deref interface {
		Deref(env *Env) (any, error)
	}
	Native interface {
		Native() interface{}
	}
	KVReduce interface {
		kvreduce(env *Env, c Callable, init any) (any, error)
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
		value any
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
	Double struct {
		InfoHolder
		D float64
	}
	BigFloat struct {
		InfoHolder
		b big.Float
	}
	Ratio struct {
		InfoHolder
		r big.Rat
	}
	Boolean bool
	Regex   struct {
		InfoHolder
		R *regexp.Regexp
	}
	Time struct {
		InfoHolder
		T time.Time
	}
	RecurBindings []any
	Delay         struct {
		fn    Callable
		value any
	}
	SortableSlice struct {
		env *Env
		s   []any
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
	_ Conjable = &HashMap{}
	_ Conjable = &Vector{}
	_ Conjable = NIL
	_ Conjable = &ArrayMap{}
	_ Conjable = &MapSet{}
	_ Conjable = &List{}

	_ Counted = &Vector{}
	_ Counted = &List{}
	_ Counted = NIL
	_ Counted = GoString("")
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
	_ Meta = Symbol(nil)
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
	_ Comparable = Char(nil)
	_ Comparable = Double{}
	_ Comparable = Int(0)
	_ Comparable = Boolean(true)
	_ Comparable = Time{}
	_ Comparable = Keyword(nil)
	_ Comparable = Symbol(nil)
	_ Comparable = &Vector{}

	_ Comparator = &Fn{}
	_ Comparator = Proc{}

	_ Indexed = GoString("")
	_ Indexed = &Vector{}

	_ Stack = &Vector{}
	_ Stack = &List{}

	_ Gettable = &HashMap{}
	_ Gettable = NIL
	_ Gettable = &Vector{}
	_ Gettable = &ArrayMap{}
	_ Gettable = &MapSet{}

	_ Associative = NIL
	_ Associative = &Vector{}
	_ Associative = &HashMap{}
	_ Associative = &ArrayMap{}

	_ Reversible = &Vector{}

	_ Named = Keyword(nil)
	_ Named = Symbol(nil)

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

type HasInfo interface {
	GetInfo() *ObjectInfo
	WithInfo(info *ObjectInfo) any
}

type ReadObject interface {
	HasInfo
}

func SetInfo(obj any, info *ObjectInfo) any {
	if hi, ok := obj.(HasInfo); ok {
		return hi.WithInfo(info)
	}

	return obj
}

func GetInfo(obj any) *ObjectInfo {
	if hi, ok := obj.(HasInfo); ok {
		return hi.GetInfo()
	}

	return nil
}

func GetMeta(obj any) Map {
	if m, ok := obj.(interface{ GetMeta() Map }); ok {
		return m.GetMeta()
	}

	return nil
}

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

func CheckArity(env *Env, args []any, min int, max int) error {
	n := len(args)
	if n < min || n > max {
		return ReturnArityMinMax(env, n, min, max)
	}
	return nil
}

func getMap(env *Env, k any, args []any) (any, error) {
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

func HashPtr[T any](val *T) uint32 {
	ptr := uintptr(unsafe.Pointer(val))
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
	v, err := ToString(env, a.value)
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
	return HashPtr(a), nil
}

func (a *Atom) WithInfo(info *ObjectInfo) any {
	return a
}

func (a *Atom) WithMeta(env *Env, meta Map) (any, error) {
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

func (a *Atom) AlterMeta(env *Env, fn *Fn, args []any) (Map, error) {
	return AlterMeta(env, &a.MetaHolder, fn, args)
}

func (a *Atom) Deref(env *Env) (any, error) {
	return a.value, nil
}

func (d *Delay) ToString(env *Env, escape bool) (string, error) {
	return "#Delay[]", nil
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
	return HashPtr(d), nil
}

func (d *Delay) WithInfo(info *ObjectInfo) any {
	return d
}

func (d *Delay) Force(env *Env) (any, error) {
	if d.value == nil {
		val, err := d.fn.Call(env, []any{})
		if err != nil {
			return nil, err
		}
		d.value = val
	}
	return d.value, nil
}

func (d *Delay) Deref(env *Env) (any, error) {
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
	return HashPtr(t), nil
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

func compare(env *Env, c Callable, a, b any) (int, error) {
	val, err := c.Call(env, []any{a, b})
	if err != nil {
		return 0, err
	}

	switch r := val.(type) {
	case Boolean:
		if r {
			return -1, nil
		}

		v, err := c.Call(env, []any{b, a})
		if err != nil {
			return 0, err
		}

		b, err := AssertBoolean(env, v, "")
		if err != nil {
			return 0, err
		}

		if b {
			return 1, nil
		}
		return 0, nil
	default:
		a, err := AssertNumber(env, r, "Function is not a comparator since it returned a non-integer value")
		if err != nil {
			return 0, err
		}

		return a.Int().I(), nil
	}
}

func (b Boolean) GetInfo() *ObjectInfo {
	return nil
}

func (i InfoHolder) GetInfo() *ObjectInfo {
	return i.info
}

func (m MetaHolder) GetMeta() Map {
	return m.meta
}

func (m *MetaHolder) ClearMeta() {
	m.meta = nil
}

func ClearMeta(obj any) {
	if cm, ok := obj.(interface{ ClearMeta() }); ok {
		cm.ClearMeta()
	}
}

func AlterMeta(env *Env, m *MetaHolder, fn *Fn, args []any) (Map, error) {
	meta := m.meta
	if meta == nil {
		meta = NIL
	}

	fargs := append([]any{meta}, args...)

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

func (rat *Ratio) Compare(env *Env, other any) (int, error) {
	n, err := AssertNumber(env, other, "Cannot compare Ratio and "+TypeName(other))
	if err != nil {
		return 0, err
	}

	return CompareNumbers(rat, n), nil
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

func (bf *BigFloat) Compare(env *Env, other any) (int, error) {
	n, err := AssertNumber(env, other, "Cannot compare BigFloat and "+TypeName(other))
	if err != nil {
		return 0, err
	}

	return CompareNumbers(bf, n), nil
}

func MakeBoolean(b bool) Boolean {
	return Boolean(b)
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

func (d Double) Compare(env *Env, other any) (int, error) {
	n, err := AssertNumber(env, other, "Cannot compare Double and "+TypeName(other))
	if err != nil {
		return 0, err
	}
	return CompareNumbers(d, n), nil
}

func (b Boolean) ToString(env *Env, escape bool) (string, error) {
	return fmt.Sprintf("%t", b), nil
}

func (b Boolean) Equals(env *Env, other interface{}) bool {
	switch other := other.(type) {
	case Boolean:
		return b == other
	default:
		return false
	}
}

func (b Boolean) GetType() *Type {
	return TYPE.Boolean
}

func (b Boolean) Native() interface{} {
	return b
}

func (b Boolean) Hash(env *Env) (uint32, error) {
	h := getHash()
	var bs = make([]byte, 1)
	if b {
		bs[0] = 1
	} else {
		bs[0] = 0
	}
	h.Write(bs)
	return h.Sum32(), nil
}

func (b Boolean) Compare(env *Env, other any) (int, error) {
	b2, err := AssertBoolean(env, other, "Cannot compare Boolean and "+TypeName(other))
	if err != nil {
		return 0, err
	}
	if b == b2 {
		return 0, nil
	}
	if b {
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

func (t Time) Compare(env *Env, other any) (int, error) {
	t2, err := AssertTime(env, other, "Cannot compare Time and "+TypeName(other))
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
	return HashPtr(rx), nil
}

func MakeStringVector(ss []string) *Vector {
	res := EmptyVector()
	for _, s := range ss {
		res, _ = res.Conjoin(MakeString(s))
	}
	return res
}

func IsVector(obj any) bool {
	switch obj.(type) {
	case *Vector:
		return true
	default:
		return false
	}
}

func IsSeq(obj any) bool {
	switch obj.(type) {
	case Seq:
		return true
	default:
		return false
	}
}

func (x *Type) WithInfo(info *ObjectInfo) any {
	return x
}

func (x RecurBindings) WithInfo(info *ObjectInfo) any {
	return x
}

func IsEqualOrImplements(abstractType HasReflectType, concreteType HasReflectType) bool {
	at := abstractType.ReflectType()
	ct := concreteType.ReflectType()

	if at.Kind() == reflect.Interface {
		return ct.Implements(at)
	} else {
		return ct == at
	}
}

func IsInstance(env *Env, t *Type, obj any) bool {
	if Equals(env, obj, NIL) {
		return false
	}
	if hrt, ok := GetType(obj).(HasReflectType); ok {
		return IsEqualOrImplements(t, hrt)
	}

	return false
}

var specialSymbols = make(map[string]struct{})

func init() {
	specialSymbols[criticalSymbols._if.Name()] = struct{}{}
	specialSymbols[criticalSymbols.quote.Name()] = struct{}{}
	specialSymbols[criticalSymbols.fn_.Name()] = struct{}{}
	specialSymbols[criticalSymbols.let_.Name()] = struct{}{}
	specialSymbols[criticalSymbols.letfn_.Name()] = struct{}{}
	specialSymbols[criticalSymbols.loop_.Name()] = struct{}{}
	specialSymbols[criticalSymbols.recur.Name()] = struct{}{}
	specialSymbols[criticalSymbols.setMacro_.Name()] = struct{}{}
	specialSymbols[criticalSymbols.def.Name()] = struct{}{}
	specialSymbols[criticalSymbols.defLinter.Name()] = struct{}{}
	specialSymbols[criticalSymbols._var.Name()] = struct{}{}
	specialSymbols[criticalSymbols.do.Name()] = struct{}{}
	specialSymbols[criticalSymbols.throw.Name()] = struct{}{}
	specialSymbols[criticalSymbols.try.Name()] = struct{}{}
	specialSymbols[criticalSymbols.catch.Name()] = struct{}{}
	specialSymbols[criticalSymbols.finally.Name()] = struct{}{}
}

func IsSpecialSymbol(obj any) bool {
	switch obj := obj.(type) {
	case Symbol:
		if obj.Namespace() != "" {
			return false
		}

		_, found := specialSymbols[obj.Name()]
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
	res.AddEqu(criticalKeywords.doc, MakeString(docstring))
	res.AddEqu(criticalKeywords.added, MakeString(added))
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

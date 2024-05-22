package core

import (
	"math/big"
	"reflect"

	"github.com/lab47/lace/pkg/pkgreflect"
)

type CharImpl struct {
	ChFn       func() rune
	CompareFn  func(*Env, any) (int, error)
	charTypeFn func() string
}

func (s *CharImpl) Ch() rune {
	return s.ChFn()
}
func (s *CharImpl) Compare(a0 *Env, a1 any) (int, error) {
	return s.CompareFn(a0, a1)
}
func (s *CharImpl) charType() string {
	return s.charTypeFn()
}

type IntegerImpl struct {
	I64Fn         func() int64
	integerTypeFn func() string
}

func (s *IntegerImpl) I64() int64 {
	return s.I64Fn()
}
func (s *IntegerImpl) integerType() string {
	return s.integerTypeFn()
}

type KeywordImpl struct {
	CallFn        func(*Env, []any) (any, error)
	CompareFn     func(*Env, any) (int, error)
	GetInfoFn     func() *ObjectInfo
	IsFn          func(any) bool
	IsHashFn      func() uint32
	NameFn        func() string
	NamespaceFn   func() string
	RawStringFn   func() string
	StringFn      func() string
	WithInfoFn    func(*ObjectInfo) any
	keywordTypeFn func() string
}

func (s *KeywordImpl) Call(a0 *Env, a1 []any) (any, error) {
	return s.CallFn(a0, a1)
}
func (s *KeywordImpl) Compare(a0 *Env, a1 any) (int, error) {
	return s.CompareFn(a0, a1)
}
func (s *KeywordImpl) GetInfo() *ObjectInfo {
	return s.GetInfoFn()
}
func (s *KeywordImpl) Is(a0 any) bool {
	return s.IsFn(a0)
}
func (s *KeywordImpl) IsHash() uint32 {
	return s.IsHashFn()
}
func (s *KeywordImpl) Name() string {
	return s.NameFn()
}
func (s *KeywordImpl) Namespace() string {
	return s.NamespaceFn()
}
func (s *KeywordImpl) RawString() string {
	return s.RawStringFn()
}
func (s *KeywordImpl) String() string {
	return s.StringFn()
}
func (s *KeywordImpl) WithInfo(a0 *ObjectInfo) any {
	return s.WithInfoFn(a0)
}
func (s *KeywordImpl) keywordType() string {
	return s.keywordTypeFn()
}

type MapImpl struct {
	AssocFn   func(*Env, any, any) (Associative, error)
	ConjFn    func(*Env, any) (Conjable, error)
	CountFn   func() int
	EntryAtFn func(*Env, any) (*Vector, error)
	GetFn     func(*Env, any) (bool, any, error)
	GetEquFn  func(Equ) (bool, any)
	IterFn    func() MapIterator
	KeysFn    func() Seq
	MergeFn   func(*Env, Map) (Map, error)
	SeqFn     func() Seq
	ValsFn    func() Seq
	WithoutFn func(*Env, any) (Map, error)
}

func (s *MapImpl) Assoc(a0 *Env, a1 any, a2 any) (Associative, error) {
	return s.AssocFn(a0, a1, a2)
}
func (s *MapImpl) Conj(a0 *Env, a1 any) (Conjable, error) {
	return s.ConjFn(a0, a1)
}
func (s *MapImpl) Count() int {
	return s.CountFn()
}
func (s *MapImpl) EntryAt(a0 *Env, a1 any) (*Vector, error) {
	return s.EntryAtFn(a0, a1)
}
func (s *MapImpl) Get(a0 *Env, a1 any) (bool, any, error) {
	return s.GetFn(a0, a1)
}
func (s *MapImpl) GetEqu(a0 Equ) (bool, any) {
	return s.GetEquFn(a0)
}
func (s *MapImpl) Iter() MapIterator {
	return s.IterFn()
}
func (s *MapImpl) Keys() Seq {
	return s.KeysFn()
}
func (s *MapImpl) Merge(a0 *Env, a1 Map) (Map, error) {
	return s.MergeFn(a0, a1)
}
func (s *MapImpl) Seq() Seq {
	return s.SeqFn()
}
func (s *MapImpl) Vals() Seq {
	return s.ValsFn()
}
func (s *MapImpl) Without(a0 *Env, a1 any) (Map, error) {
	return s.WithoutFn(a0, a1)
}

type NilImpl struct {
	AssocFn    func(*Env, any, any) (Associative, error)
	ConjFn     func(*Env, any) (Conjable, error)
	ConsFn     func(any) Seq
	CountFn    func() int
	EntryAtFn  func(*Env, any) (*Vector, error)
	FirstFn    func(*Env) (any, error)
	GetFn      func(*Env, any) (bool, any, error)
	GetEquFn   func(Equ) (bool, any)
	GetInfoFn  func() *ObjectInfo
	IsEmptyFn  func(*Env) (bool, error)
	IterFn     func() MapIterator
	KeysFn     func() Seq
	MergeFn    func(*Env, Map) (Map, error)
	RestFn     func(*Env) (Seq, error)
	SeqFn      func() Seq
	ValsFn     func() Seq
	WithInfoFn func(*ObjectInfo) any
	WithoutFn  func(*Env, any) (Map, error)
	nilTypeFn  func() string
}

func (s *NilImpl) Assoc(a0 *Env, a1 any, a2 any) (Associative, error) {
	return s.AssocFn(a0, a1, a2)
}
func (s *NilImpl) Conj(a0 *Env, a1 any) (Conjable, error) {
	return s.ConjFn(a0, a1)
}
func (s *NilImpl) Cons(a0 any) Seq {
	return s.ConsFn(a0)
}
func (s *NilImpl) Count() int {
	return s.CountFn()
}
func (s *NilImpl) EntryAt(a0 *Env, a1 any) (*Vector, error) {
	return s.EntryAtFn(a0, a1)
}
func (s *NilImpl) First(a0 *Env) (any, error) {
	return s.FirstFn(a0)
}
func (s *NilImpl) Get(a0 *Env, a1 any) (bool, any, error) {
	return s.GetFn(a0, a1)
}
func (s *NilImpl) GetEqu(a0 Equ) (bool, any) {
	return s.GetEquFn(a0)
}
func (s *NilImpl) GetInfo() *ObjectInfo {
	return s.GetInfoFn()
}
func (s *NilImpl) IsEmpty(a0 *Env) (bool, error) {
	return s.IsEmptyFn(a0)
}
func (s *NilImpl) Iter() MapIterator {
	return s.IterFn()
}
func (s *NilImpl) Keys() Seq {
	return s.KeysFn()
}
func (s *NilImpl) Merge(a0 *Env, a1 Map) (Map, error) {
	return s.MergeFn(a0, a1)
}
func (s *NilImpl) Rest(a0 *Env) (Seq, error) {
	return s.RestFn(a0)
}
func (s *NilImpl) Seq() Seq {
	return s.SeqFn()
}
func (s *NilImpl) Vals() Seq {
	return s.ValsFn()
}
func (s *NilImpl) WithInfo(a0 *ObjectInfo) any {
	return s.WithInfoFn(a0)
}
func (s *NilImpl) Without(a0 *Env, a1 any) (Map, error) {
	return s.WithoutFn(a0, a1)
}
func (s *NilImpl) nilType() string {
	return s.nilTypeFn()
}

type NumberImpl struct {
	BigFloatFn     func() *big.Float
	BigIntFn       func() *big.Int
	DoubleFn       func() Double
	IntFn          func() Int
	NativeNumberFn func() any
	RatioFn        func() *big.Rat
}

func (s *NumberImpl) BigFloat() *big.Float {
	return s.BigFloatFn()
}
func (s *NumberImpl) BigInt() *big.Int {
	return s.BigIntFn()
}
func (s *NumberImpl) Double() Double {
	return s.DoubleFn()
}
func (s *NumberImpl) Int() Int {
	return s.IntFn()
}
func (s *NumberImpl) NativeNumber() any {
	return s.NativeNumberFn()
}
func (s *NumberImpl) Ratio() *big.Rat {
	return s.RatioFn()
}

type AssociativeImpl struct {
	AssocFn   func(*Env, any, any) (Associative, error)
	ConjFn    func(*Env, any) (Conjable, error)
	EntryAtFn func(*Env, any) (*Vector, error)
	GetFn     func(*Env, any) (bool, any, error)
}

func (s *AssociativeImpl) Assoc(a0 *Env, a1 any, a2 any) (Associative, error) {
	return s.AssocFn(a0, a1, a2)
}
func (s *AssociativeImpl) Conj(a0 *Env, a1 any) (Conjable, error) {
	return s.ConjFn(a0, a1)
}
func (s *AssociativeImpl) EntryAt(a0 *Env, a1 any) (*Vector, error) {
	return s.EntryAtFn(a0, a1)
}
func (s *AssociativeImpl) Get(a0 *Env, a1 any) (bool, any, error) {
	return s.GetFn(a0, a1)
}

type CollectionImpl struct {
	CountFn func() int
	EmptyFn func() Collection
	SeqFn   func() Seq
}

func (s *CollectionImpl) Count() int {
	return s.CountFn()
}
func (s *CollectionImpl) Empty() Collection {
	return s.EmptyFn()
}
func (s *CollectionImpl) Seq() Seq {
	return s.SeqFn()
}

type ComparableImpl struct {
	CompareFn func(*Env, any) (int, error)
}

func (s *ComparableImpl) Compare(a0 *Env, a1 any) (int, error) {
	return s.CompareFn(a0, a1)
}

type ComparatorImpl struct {
	CompareFn func(*Env, any, any) (int, error)
}

func (s *ComparatorImpl) Compare(a0 *Env, a1 any, a2 any) (int, error) {
	return s.CompareFn(a0, a1, a2)
}

type CountedImpl struct {
	CountFn func() int
}

func (s *CountedImpl) Count() int {
	return s.CountFn()
}

type DerefImpl struct {
	DerefFn func(*Env) (any, error)
}

func (s *DerefImpl) Deref(a0 *Env) (any, error) {
	return s.DerefFn(a0)
}

type ErrorImpl struct {
	ErrorFn   func() string
	MessageFn func() any
}

func (s *ErrorImpl) Error() string {
	return s.ErrorFn()
}
func (s *ErrorImpl) Message() any {
	return s.MessageFn()
}

type GettableImpl struct {
	GetFn func(*Env, any) (bool, any, error)
}

func (s *GettableImpl) Get(a0 *Env, a1 any) (bool, any, error) {
	return s.GetFn(a0, a1)
}

type IndexedImpl struct {
	NthFn    func(*Env, int) (any, error)
	TryNthFn func(*Env, int, any) (any, error)
}

func (s *IndexedImpl) Nth(a0 *Env, a1 int) (any, error) {
	return s.NthFn(a0, a1)
}
func (s *IndexedImpl) TryNth(a0 *Env, a1 int, a2 any) (any, error) {
	return s.TryNthFn(a0, a1, a2)
}

type KVReduceImpl struct {
	kvreduceFn func(*Env, Callable, any) (any, error)
}

func (s *KVReduceImpl) kvreduce(a0 *Env, a1 Callable, a2 any) (any, error) {
	return s.kvreduceFn(a0, a1, a2)
}

type MetaImpl struct {
	GetMetaFn  func() Map
	WithMetaFn func(*Env, Map) (any, error)
}

func (s *MetaImpl) GetMeta() Map {
	return s.GetMetaFn()
}
func (s *MetaImpl) WithMeta(a0 *Env, a1 Map) (any, error) {
	return s.WithMetaFn(a0, a1)
}

type NamedImpl struct {
	NameFn      func() string
	NamespaceFn func() string
}

func (s *NamedImpl) Name() string {
	return s.NameFn()
}
func (s *NamedImpl) Namespace() string {
	return s.NamespaceFn()
}

type PendingImpl struct {
	IsRealizedFn func() bool
}

func (s *PendingImpl) IsRealized() bool {
	return s.IsRealizedFn()
}

type RefImpl struct {
	AlterMetaFn func(*Env, *Fn, []any) (Map, error)
	ResetMetaFn func(Map) Map
}

func (s *RefImpl) AlterMeta(a0 *Env, a1 *Fn, a2 []any) (Map, error) {
	return s.AlterMetaFn(a0, a1, a2)
}
func (s *RefImpl) ResetMeta(a0 Map) Map {
	return s.ResetMetaFn(a0)
}

type ReversibleImpl struct {
	RseqFn func() Seq
}

func (s *ReversibleImpl) Rseq() Seq {
	return s.RseqFn()
}

type SequentialImpl struct {
	sequentialFn func()
}

func (s *SequentialImpl) sequential() {
	s.sequentialFn()
}

type StackImpl struct {
	PeekFn func(*Env) (any, error)
	PopFn  func(*Env) (Stack, error)
}

func (s *StackImpl) Peek(a0 *Env) (any, error) {
	return s.PeekFn(a0)
}
func (s *StackImpl) Pop(a0 *Env) (Stack, error) {
	return s.PopFn(a0)
}

type CallableImpl struct {
	CallFn func(*Env, []any) (any, error)
}

func (s *CallableImpl) Call(a0 *Env, a1 []any) (any, error) {
	return s.CallFn(a0, a1)
}

type SeqImpl struct {
	ConsFn    func(any) Seq
	FirstFn   func(*Env) (any, error)
	IsEmptyFn func(*Env) (bool, error)
	RestFn    func(*Env) (Seq, error)
	SeqFn     func() Seq
}

func (s *SeqImpl) Cons(a0 any) Seq {
	return s.ConsFn(a0)
}
func (s *SeqImpl) First(a0 *Env) (any, error) {
	return s.FirstFn(a0)
}
func (s *SeqImpl) IsEmpty(a0 *Env) (bool, error) {
	return s.IsEmptyFn(a0)
}
func (s *SeqImpl) Rest(a0 *Env) (Seq, error) {
	return s.RestFn(a0)
}
func (s *SeqImpl) Seq() Seq {
	return s.SeqFn()
}

type SeqableImpl struct {
	SeqFn func() Seq
}

func (s *SeqableImpl) Seq() Seq {
	return s.SeqFn()
}

type SetImpl struct {
	ConjFn    func(*Env, any) (Conjable, error)
	DisjoinFn func(*Env, any) (Set, error)
	GetFn     func(*Env, any) (bool, any, error)
	HasFn     func(Equ) bool
	SetIterFn func() SetIter
}

func (s *SetImpl) Conj(a0 *Env, a1 any) (Conjable, error) {
	return s.ConjFn(a0, a1)
}
func (s *SetImpl) Disjoin(a0 *Env, a1 any) (Set, error) {
	return s.DisjoinFn(a0, a1)
}
func (s *SetImpl) Get(a0 *Env, a1 any) (bool, any, error) {
	return s.GetFn(a0, a1)
}
func (s *SetImpl) Has(a0 Equ) bool {
	return s.HasFn(a0)
}
func (s *SetImpl) SetIter() SetIter {
	return s.SetIterFn()
}

type StringImpl struct {
	AppendToFn func(String) String
	CountFn    func() int
	NthFn      func(*Env, int) (any, error)
	SFn        func() string
	SeqFn      func() Seq
	TryNthFn   func(*Env, int, any) (any, error)
}

func (s *StringImpl) AppendTo(a0 String) String {
	return s.AppendToFn(a0)
}
func (s *StringImpl) Count() int {
	return s.CountFn()
}
func (s *StringImpl) Nth(a0 *Env, a1 int) (any, error) {
	return s.NthFn(a0, a1)
}
func (s *StringImpl) S() string {
	return s.SFn()
}
func (s *StringImpl) Seq() Seq {
	return s.SeqFn()
}
func (s *StringImpl) TryNth(a0 *Env, a1 int, a2 any) (any, error) {
	return s.TryNthFn(a0, a1, a2)
}

type SymbolImpl struct {
	CompareFn    func(*Env, any) (int, error)
	GetMetaFn    func() Map
	IsFn         func(any) bool
	IsHashFn     func() uint32
	NameFn       func() string
	NamespaceFn  func() string
	StringFn     func() string
	WithMetaFn   func(*Env, Map) (any, error)
	symbolTypeFn func() string
}

func (s *SymbolImpl) Compare(a0 *Env, a1 any) (int, error) {
	return s.CompareFn(a0, a1)
}
func (s *SymbolImpl) GetMeta() Map {
	return s.GetMetaFn()
}
func (s *SymbolImpl) Is(a0 any) bool {
	return s.IsFn(a0)
}
func (s *SymbolImpl) IsHash() uint32 {
	return s.IsHashFn()
}
func (s *SymbolImpl) Name() string {
	return s.NameFn()
}
func (s *SymbolImpl) Namespace() string {
	return s.NamespaceFn()
}
func (s *SymbolImpl) String() string {
	return s.StringFn()
}
func (s *SymbolImpl) WithMeta(a0 *Env, a1 Map) (any, error) {
	return s.WithMetaFn(a0, a1)
}
func (s *SymbolImpl) symbolType() string {
	return s.symbolTypeFn()
}

func init() {
	ArrayMap_methods := map[string]pkgreflect.Func{}
	BufferedReader_methods := map[string]pkgreflect.Func{}
	Char_methods := map[string]pkgreflect.Func{}
	EvalError_methods := map[string]pkgreflect.Func{}
	Fn_methods := map[string]pkgreflect.Func{}
	HashMap_methods := map[string]pkgreflect.Func{}
	BigInt_methods := map[string]pkgreflect.Func{}
	Int_methods := map[string]pkgreflect.Func{}
	Integer_methods := map[string]pkgreflect.Func{}
	Keyword_methods := map[string]pkgreflect.Func{}
	List_methods := map[string]pkgreflect.Func{}
	Map_methods := map[string]pkgreflect.Func{}
	Nil_methods := map[string]pkgreflect.Func{}
	Namespace_methods := map[string]pkgreflect.Func{}
	Number_methods := map[string]pkgreflect.Func{}
	Associative_methods := map[string]pkgreflect.Func{}
	Atom_methods := map[string]pkgreflect.Func{}
	BigFloat_methods := map[string]pkgreflect.Func{}
	Boolean_methods := map[string]pkgreflect.Func{}
	Collection_methods := map[string]pkgreflect.Func{}
	Comparable_methods := map[string]pkgreflect.Func{}
	Comparator_methods := map[string]pkgreflect.Func{}
	Counted_methods := map[string]pkgreflect.Func{}
	Delay_methods := map[string]pkgreflect.Func{}
	Deref_methods := map[string]pkgreflect.Func{}
	Double_methods := map[string]pkgreflect.Func{}
	Error_methods := map[string]pkgreflect.Func{}
	Gettable_methods := map[string]pkgreflect.Func{}
	Indexed_methods := map[string]pkgreflect.Func{}
	KVReduce_methods := map[string]pkgreflect.Func{}
	Meta_methods := map[string]pkgreflect.Func{}
	Named_methods := map[string]pkgreflect.Func{}
	Pending_methods := map[string]pkgreflect.Func{}
	Ratio_methods := map[string]pkgreflect.Func{}
	Ref_methods := map[string]pkgreflect.Func{}
	Regex_methods := map[string]pkgreflect.Func{}
	Reversible_methods := map[string]pkgreflect.Func{}
	Sequential_methods := map[string]pkgreflect.Func{}
	Stack_methods := map[string]pkgreflect.Func{}
	Time_methods := map[string]pkgreflect.Func{}
	Callable_methods := map[string]pkgreflect.Func{}
	Seq_methods := map[string]pkgreflect.Func{}
	Seqable_methods := map[string]pkgreflect.Func{}
	MapSet_methods := map[string]pkgreflect.Func{}
	Set_methods := map[string]pkgreflect.Func{}
	String_methods := map[string]pkgreflect.Func{}
	Symbol_methods := map[string]pkgreflect.Func{}
	Type_methods := map[string]pkgreflect.Func{}
	Var_methods := map[string]pkgreflect.Func{}
	Vector_methods := map[string]pkgreflect.Func{}
	pkgreflect.AddPackage("lace.lang", &pkgreflect.Package{
		Doc: "",
		Types: map[string]pkgreflect.Type{
			"ArrayMap":        {Doc: "A Map implementation that uses a simple array. Very efficient for small maps.", Value: reflect.TypeOf((*ArrayMap)(nil)).Elem(), Methods: ArrayMap_methods},
			"BufferedReader":  {Doc: "A value that can return data that has been buffered.", Value: reflect.TypeOf((*BufferedReader)(nil)).Elem(), Methods: BufferedReader_methods},
			"Char":            {Doc: "A single unicode rune.", Value: reflect.TypeOf((*Char)(nil)).Elem(), Methods: Char_methods},
			"EvalError":       {Doc: "The standard error thrown when evalution of a program encounters an error.", Value: reflect.TypeOf((*EvalError)(nil)).Elem(), Methods: EvalError_methods},
			"Fn":              {Doc: "A value that contains code and can be called to run that code.", Value: reflect.TypeOf((*Fn)(nil)).Elem(), Methods: Fn_methods},
			"HashMap":         {Doc: "A Map implementation that can store a large number of values efficiently.", Value: reflect.TypeOf((*HashMap)(nil)).Elem(), Methods: HashMap_methods},
			"BigInt":          {Doc: "An integer value that can be so large, it's hard to understand it.", Value: reflect.TypeOf((*BigInt)(nil)).Elem(), Methods: BigInt_methods},
			"Int":             {Doc: "The host int value", Value: reflect.TypeOf((*Int)(nil)).Elem(), Methods: Int_methods},
			"Integer":         {Doc: "The Common integer type (can be Int or BigInt)", Value: reflect.TypeOf((*Integer)(nil)).Elem(), Methods: Integer_methods},
			"Keyword":         {Doc: "A value that is just a name who's meaning is itself. It sounds meta,\nI know. It's just a name, usually used as a key in a association.\nMaybe you want to think of it as a short, compact, namespace'd string?\nThat's fine.", Value: reflect.TypeOf((*Keyword)(nil)).Elem(), Methods: Keyword_methods},
			"List":            {Doc: "The standard lisp persistent list.", Value: reflect.TypeOf((*List)(nil)).Elem(), Methods: List_methods},
			"Map":             {Doc: "A collection that contains the full complement of functionality for\ndealing with having objects that are stored by association with a key.", Value: reflect.TypeOf((*Map)(nil)).Elem(), Methods: Map_methods},
			"Nil":             {Doc: "The nothing type.", Value: reflect.TypeOf((*Nil)(nil)).Elem(), Methods: Nil_methods},
			"Namespace":       {Doc: "A namespace is a named collection and the foundation of clojure/lace.\nA namespace holds vars and functions run in the context of a namespace.", Value: reflect.TypeOf((*Namespace)(nil)).Elem(), Methods: Namespace_methods},
			"Number":          {Doc: "A number is any kind of sequence of numerals. This includes normal,\nintegers of any size, floating point values of any size, as well as\nratios (like 1/3).", Value: reflect.TypeOf((*Number)(nil)).Elem(), Methods: Number_methods},
			"Associative":     {Doc: "When a collection can store and retrieve values by an associated key.", Value: reflect.TypeOf((*Associative)(nil)).Elem(), Methods: Associative_methods},
			"Atom":            {Doc: "A value that contains space for a single other value that can be swapped\nin. Ie an atom is an atomic value.", Value: reflect.TypeOf((*Atom)(nil)).Elem(), Methods: Atom_methods},
			"BigFloat":        {Doc: "A floating point value that can be any size, practically.", Value: reflect.TypeOf((*BigFloat)(nil)).Elem(), Methods: BigFloat_methods},
			"Boolean":         {Doc: "It's true, or it's false. Never both.", Value: reflect.TypeOf((*Boolean)(nil)).Elem(), Methods: Boolean_methods},
			"Collection":      {Doc: "When a value has the standard collection interfaces.", Value: reflect.TypeOf((*Collection)(nil)).Elem(), Methods: Collection_methods},
			"Comparable":      {Doc: "When a value can report if it's less, same, or bigger than another value.", Value: reflect.TypeOf((*Comparable)(nil)).Elem(), Methods: Comparable_methods},
			"Comparator":      {Doc: "When a value can report if two values are are less, same, or bigger than another value.", Value: reflect.TypeOf((*Comparator)(nil)).Elem(), Methods: Comparator_methods},
			"Counted":         {Doc: "Values that can report how many they contain.", Value: reflect.TypeOf((*Counted)(nil)).Elem(), Methods: Counted_methods},
			"Delay":           {Doc: "A value that runs a function to produce a new value.\nIe, it 'delays' running the function until the value is needed.", Value: reflect.TypeOf((*Delay)(nil)).Elem(), Methods: Delay_methods},
			"Deref":           {Doc: "When a value can be dereference to return another value.", Value: reflect.TypeOf((*Deref)(nil)).Elem(), Methods: Deref_methods},
			"Double":          {Doc: "A floating point value who's size is constraint to the host's\nlargest floating point value.", Value: reflect.TypeOf((*Double)(nil)).Elem(), Methods: Double_methods},
			"Error":           {Doc: "When things how wrong, these show up.", Value: reflect.TypeOf((*Error)(nil)).Elem(), Methods: Error_methods},
			"Gettable":        {Doc: "When a collection can return a value by it's key.", Value: reflect.TypeOf((*Gettable)(nil)).Elem(), Methods: Gettable_methods},
			"Indexed":         {Doc: "When a collection can return a contained value an a fixed offset.", Value: reflect.TypeOf((*Indexed)(nil)).Elem(), Methods: Indexed_methods},
			"KVReduce":        {Doc: "When a value can orchestrate reducing itself value a callable.", Value: reflect.TypeOf((*KVReduce)(nil)).Elem(), Methods: KVReduce_methods},
			"Meta":            {Doc: "When a value can contain additional information.", Value: reflect.TypeOf((*Meta)(nil)).Elem(), Methods: Meta_methods},
			"Named":           {Doc: "When a value has a name.", Value: reflect.TypeOf((*Named)(nil)).Elem(), Methods: Named_methods},
			"Pending":         {Doc: "When a value can report if it has a pending operation.", Value: reflect.TypeOf((*Pending)(nil)).Elem(), Methods: Pending_methods},
			"Ratio":           {Doc: "A value that presents the division of two integers.", Value: reflect.TypeOf((*Ratio)(nil)).Elem(), Methods: Ratio_methods},
			"Ref":             {Doc: "When a value can change it's metadata.", Value: reflect.TypeOf((*Ref)(nil)).Elem(), Methods: Ref_methods},
			"Regex":           {Doc: "A value that contains a pre-compiled regular expression.", Value: reflect.TypeOf((*Regex)(nil)).Elem(), Methods: Regex_methods},
			"Reversible":      {Doc: "When a collection can return a sequence that returns values in reverse order.", Value: reflect.TypeOf((*Reversible)(nil)).Elem(), Methods: Reversible_methods},
			"Sequential":      {Doc: "When a value is sequential.", Value: reflect.TypeOf((*Sequential)(nil)).Elem(), Methods: Sequential_methods},
			"Stack":           {Doc: "When a container can add and remove values efficiently without regard for position.", Value: reflect.TypeOf((*Stack)(nil)).Elem(), Methods: Stack_methods},
			"Time":            {Doc: "A value that represents a point in time.", Value: reflect.TypeOf((*Time)(nil)).Elem(), Methods: Time_methods},
			"Callable":        {Doc: "When a value can be called and passed arguments, like a Function.", Value: reflect.TypeOf((*Callable)(nil)).Elem(), Methods: Callable_methods},
			"Seq":             {Doc: "The Seq interface, provides a sequence of values.", Value: reflect.TypeOf((*Seq)(nil)).Elem(), Methods: Seq_methods},
			"Seqable":         {Doc: "When a value can be converted into a Seq.", Value: reflect.TypeOf((*Seqable)(nil)).Elem(), Methods: Seqable_methods},
			"MapSet":          {Doc: "A Set implementation that uses a Map.", Value: reflect.TypeOf((*MapSet)(nil)).Elem(), Methods: MapSet_methods},
			"Set":             {Doc: "A collection that can store unique values.", Value: reflect.TypeOf((*Set)(nil)).Elem(), Methods: Set_methods},
			"String":          {Doc: "A sequence of bytes, usually containing utf-8.", Value: reflect.TypeOf((*String)(nil)).Elem(), Methods: String_methods},
			"Symbol":          {Doc: "A value that represents a value stored elsewhere by name,\nsuch as in a namespace or a local variable.", Value: reflect.TypeOf((*Symbol)(nil)).Elem(), Methods: Symbol_methods},
			"Type":            {Doc: "A value that describes a set of values.", Value: reflect.TypeOf((*Type)(nil)).Elem(), Methods: Type_methods},
			"Var":             {Doc: "A value that holds another value and can be changed. Ie, a Var is a variable.", Value: reflect.TypeOf((*Var)(nil)).Elem(), Methods: Var_methods},
			"Vector":          {Doc: "A collection that stores it's values at fixed integer offsets efficiently.", Value: reflect.TypeOf((*Vector)(nil)).Elem(), Methods: Vector_methods},
			"CharImpl":        {Doc: `Struct version of interface Char for implementation`, Value: reflect.TypeFor[CharImpl]()},
			"IntegerImpl":     {Doc: `Struct version of interface Integer for implementation`, Value: reflect.TypeFor[IntegerImpl]()},
			"KeywordImpl":     {Doc: `Struct version of interface Keyword for implementation`, Value: reflect.TypeFor[KeywordImpl]()},
			"MapImpl":         {Doc: `Struct version of interface Map for implementation`, Value: reflect.TypeFor[MapImpl]()},
			"NilImpl":         {Doc: `Struct version of interface Nil for implementation`, Value: reflect.TypeFor[NilImpl]()},
			"NumberImpl":      {Doc: `Struct version of interface Number for implementation`, Value: reflect.TypeFor[NumberImpl]()},
			"AssociativeImpl": {Doc: `Struct version of interface Associative for implementation`, Value: reflect.TypeFor[AssociativeImpl]()},
			"CollectionImpl":  {Doc: `Struct version of interface Collection for implementation`, Value: reflect.TypeFor[CollectionImpl]()},
			"ComparableImpl":  {Doc: `Struct version of interface Comparable for implementation`, Value: reflect.TypeFor[ComparableImpl]()},
			"ComparatorImpl":  {Doc: `Struct version of interface Comparator for implementation`, Value: reflect.TypeFor[ComparatorImpl]()},
			"CountedImpl":     {Doc: `Struct version of interface Counted for implementation`, Value: reflect.TypeFor[CountedImpl]()},
			"DerefImpl":       {Doc: `Struct version of interface Deref for implementation`, Value: reflect.TypeFor[DerefImpl]()},
			"ErrorImpl":       {Doc: `Struct version of interface Error for implementation`, Value: reflect.TypeFor[ErrorImpl]()},
			"GettableImpl":    {Doc: `Struct version of interface Gettable for implementation`, Value: reflect.TypeFor[GettableImpl]()},
			"IndexedImpl":     {Doc: `Struct version of interface Indexed for implementation`, Value: reflect.TypeFor[IndexedImpl]()},
			"KVReduceImpl":    {Doc: `Struct version of interface KVReduce for implementation`, Value: reflect.TypeFor[KVReduceImpl]()},
			"MetaImpl":        {Doc: `Struct version of interface Meta for implementation`, Value: reflect.TypeFor[MetaImpl]()},
			"NamedImpl":       {Doc: `Struct version of interface Named for implementation`, Value: reflect.TypeFor[NamedImpl]()},
			"PendingImpl":     {Doc: `Struct version of interface Pending for implementation`, Value: reflect.TypeFor[PendingImpl]()},
			"RefImpl":         {Doc: `Struct version of interface Ref for implementation`, Value: reflect.TypeFor[RefImpl]()},
			"ReversibleImpl":  {Doc: `Struct version of interface Reversible for implementation`, Value: reflect.TypeFor[ReversibleImpl]()},
			"SequentialImpl":  {Doc: `Struct version of interface Sequential for implementation`, Value: reflect.TypeFor[SequentialImpl]()},
			"StackImpl":       {Doc: `Struct version of interface Stack for implementation`, Value: reflect.TypeFor[StackImpl]()},
			"CallableImpl":    {Doc: `Struct version of interface Callable for implementation`, Value: reflect.TypeFor[CallableImpl]()},
			"SeqImpl":         {Doc: `Struct version of interface Seq for implementation`, Value: reflect.TypeFor[SeqImpl]()},
			"SeqableImpl":     {Doc: `Struct version of interface Seqable for implementation`, Value: reflect.TypeFor[SeqableImpl]()},
			"SetImpl":         {Doc: `Struct version of interface Set for implementation`, Value: reflect.TypeFor[SetImpl]()},
			"StringImpl":      {Doc: `Struct version of interface String for implementation`, Value: reflect.TypeFor[StringImpl]()},
			"SymbolImpl":      {Doc: `Struct version of interface Symbol for implementation`, Value: reflect.TypeFor[SymbolImpl]()},
		},

		Functions: map[string]pkgreflect.FuncValue{
			"CombineToString": {Doc: "Combine many values into a single string.", Args: []pkgreflect.Arg{{Name: "args", Tag: "[]any"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc2_2(CombineToString))},

			"ConcatSimple": {Doc: "Concatinate N sequences together", Args: []pkgreflect.Arg{{Name: "args", Tag: "[]any"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc2_2(ConcatSimple))},

			"Conj": {Doc: "Create a new Sequence by combine the value with the collection.", Args: []pkgreflect.Arg{{Name: "col", Tag: "any"}, {Name: "val", Tag: "any"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc3_2(Conj))},

			"Cons": {Doc: "Add an element to a Seq value, returning a new Seq", Args: []pkgreflect.Arg{{Name: "val", Tag: "any"}, {Name: "seq", Tag: "Seqable"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc3_2(Cons))},

			"ConvertToSeq": {Doc: "Convert the given value to a Seq", Args: []pkgreflect.Arg{{Name: "s", Tag: "Seqable"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc2_2(ConvertToSeq))},

			"Equals": {Doc: "Compare two values returning a boolean if they are equal or not", Args: []pkgreflect.Arg{{Name: "a", Tag: "any"}, {Name: "b", Tag: "any"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc3_2(EqualsValues))},

			"First": {Doc: "Return the first element in a Seq", Args: []pkgreflect.Arg{{Name: "s", Tag: "Seqable"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc2_2(First))},

			"LoadLibFromPath": {Doc: "Attempt to load a given lib from a given path.", Args: []pkgreflect.Arg{{Name: "libnamev", Tag: "Symbol"}, {Name: "pathnamev", Tag: "String"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc3_2(LoadLibFromPath))},

			"MakeList": {Doc: "Create a new lace List from the given arguments", Args: []pkgreflect.Arg{{Name: "args", Tag: "[]any"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc2_2(MakeList))},

			"NewFuture": {Doc: "NewFuture creates a new Future value and schedules the future\nto be run. Deref'ing the Future will retrieve the value (potentially\nwaiting if the value is not yet ready)", Args: []pkgreflect.Arg{{Name: "call", Tag: "Callable"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc2_2(NewFuture))},

			"Next": {Doc: "Return elements other than the first one in a Seq", Args: []pkgreflect.Arg{{Name: "s", Tag: "Seqable"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc2_2(Next))},

			"PushBindings": {Doc: "Add given bindings to the set of current Var bindings, returning\nthe original set.", Args: []pkgreflect.Arg{{Name: "assoc", Tag: "Map"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc2_2(PushBindings))},

			"Rest": {Doc: "Return all elements of a seq except for the first one.", Args: []pkgreflect.Arg{{Name: "s", Tag: "Seqable"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc2_2(Rest))},

			"SetBindings": {Doc: "Reset the local var bindings to the given value.", Args: []pkgreflect.Arg{{Name: "assoc", Tag: "Associative"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc2_2(SetBindings))},

			"StartGoRoutine": {Doc: "StartGoRoutine runs the given callable in a new goroutine, returning a channel\nthat can be used to retrieve the return value.", Args: []pkgreflect.Arg{{Name: "callable", Tag: "Callable"}}, Tag: "any", Value: reflect.ValueOf(WrapToProc2_2(StartGoRoutine))},
		},

		Variables: map[string]pkgreflect.Value{},

		Consts: map[string]pkgreflect.Value{},
	})
}

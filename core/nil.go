package core

type Nil interface {
	any
	Conjable
	Counted
	Gettable
	Associative
	HasInfo
	Seq
	Map

	nilType() string
}

var NIL = TinyNil{}

var _ Map = HeavyNil{}

type TinyNil struct{}

func (x TinyNil) nilType() string { return "tiny" }

func (x TinyNil) String() string {
	return "nil"
}

func (x TinyNil) WithInfo(info *ObjectInfo) any {
	return HeavyNil{info: info}
}

func (x TinyNil) GetInfo() *ObjectInfo {
	return nil
}

func (n TinyNil) ToString(env *Env, escape bool) (string, error) {
	return "nil", nil
}

func (n TinyNil) Equals(env *Env, other interface{}) bool {
	switch other.(type) {
	case Nil:
		return true
	default:
		return false
	}
}

func (n TinyNil) GetEqu(key Equ) (bool, any) {
	return false, NIL
}

func (n TinyNil) GetType() *Type {
	return TYPE.Nil
}

func (n TinyNil) Hash(env *Env) (uint32, error) {
	return 0, nil
}

func (n TinyNil) Seq() Seq {
	return n
}

func (n TinyNil) First(env *Env) (any, error) {
	return NIL, nil
}

func (n TinyNil) Rest(env *Env) (Seq, error) {
	return NIL, nil
}

func (n TinyNil) IsEmpty(env *Env) (bool, error) {
	return true, nil
}

func (n TinyNil) Cons(obj any) Seq {
	return NewListFrom(obj)
}

func (n TinyNil) Conj(env *Env, obj any) (Conjable, error) {
	return NewListFrom(obj), nil
}

func (n TinyNil) Without(env *Env, key any) (Map, error) {
	return n, nil
}

func (n TinyNil) Count() int {
	return 0
}

func (n TinyNil) Iter() MapIterator {
	return emptyMapIterator
}

func (n TinyNil) Merge(env *Env, other Map) (Map, error) {
	return other, nil
}

func (n TinyNil) Assoc(env *Env, key, value any) (Associative, error) {
	return EmptyArrayMap().Assoc(env, key, value)
}

func (n TinyNil) EntryAt(env *Env, key any) (*Vector, error) {
	return nil, nil
}

func (n TinyNil) Get(env *Env, key any) (bool, any, error) {
	return false, NIL, nil
}

func (n TinyNil) Disjoin(env *Env, key any) (Set, error) {
	return n, nil
}

func (n TinyNil) SetIter() SetIter {
	return emptySetIterator
}

func (n TinyNil) Has(key Equ) bool {
	return false
}

func (n TinyNil) Keys() Seq {
	return NIL
}

func (n TinyNil) Vals() Seq {
	return NIL
}

type HeavyNil struct {
	info *ObjectInfo
}

func (HeavyNil) nilType() string { return "heavy" }

func (x HeavyNil) String() string {
	return "nil"
}

func (x HeavyNil) WithInfo(info *ObjectInfo) any {
	x.info = info
	return x
}

func (x HeavyNil) GetInfo() *ObjectInfo {
	return x.info
}

func (n HeavyNil) ToString(env *Env, escape bool) (string, error) {
	return "nil", nil
}

func (n HeavyNil) Equals(env *Env, other interface{}) bool {
	switch other.(type) {
	case Nil:
		return true
	default:
		return false
	}
}

func (n HeavyNil) GetEqu(key Equ) (bool, any) {
	return false, NIL
}

func (n HeavyNil) GetType() *Type {
	return TYPE.Nil
}

func (n HeavyNil) Hash(env *Env) (uint32, error) {
	return 0, nil
}

func (n HeavyNil) Seq() Seq {
	return n
}

func (n HeavyNil) First(env *Env) (any, error) {
	return NIL, nil
}

func (n HeavyNil) Rest(env *Env) (Seq, error) {
	return NIL, nil
}

func (n HeavyNil) IsEmpty(env *Env) (bool, error) {
	return true, nil
}

func (n HeavyNil) Cons(obj any) Seq {
	return NewListFrom(obj)
}

func (n HeavyNil) Conj(env *Env, obj any) (Conjable, error) {
	return NewListFrom(obj), nil
}

func (n HeavyNil) Without(env *Env, key any) (Map, error) {
	return n, nil
}

func (n HeavyNil) Count() int {
	return 0
}

func (n HeavyNil) Iter() MapIterator {
	return emptyMapIterator
}

func (n HeavyNil) Merge(env *Env, other Map) (Map, error) {
	return other, nil
}

func (n HeavyNil) Assoc(env *Env, key, value any) (Associative, error) {
	return EmptyArrayMap().Assoc(env, key, value)
}

func (n HeavyNil) EntryAt(env *Env, key any) (*Vector, error) {
	return nil, nil
}

func (n HeavyNil) Get(env *Env, key any) (bool, any, error) {
	return false, NIL, nil
}

func (n HeavyNil) Disjoin(env *Env, key any) (Set, error) {
	return n, nil
}

func (n HeavyNil) SetIter() SetIter {
	return emptySetIterator
}

func (n HeavyNil) Has(key Equ) bool {
	return false
}

func (n HeavyNil) Keys() Seq {
	return NIL
}

func (n HeavyNil) Vals() Seq {
	return NIL
}

package core

import (
	"bytes"
	"fmt"
	"io"
)

type (
	Vector struct {
		InfoHolder
		MetaHolder
		root  []interface{}
		tail  []interface{}
		count int
		shift uint
	}
	VectorSeq struct {
		InfoHolder
		MetaHolder
		vector *Vector
		index  int
	}
	VectorRSeq struct {
		InfoHolder
		MetaHolder
		vector *Vector
		index  int
	}
)

var empty_node []interface{} = make([]interface{}, 32)

func (v *Vector) WithMeta(env *Env, meta Map) (Object, error) {
	res := *v
	m, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return &res, nil
}

func clone(s []interface{}) []interface{} {
	result := make([]interface{}, len(s), cap(s))
	copy(result, s)
	return result
}

func (v *Vector) tailoff() int {
	if v.count < 32 {
		return 0
	}
	return ((v.count - 1) >> 5) << 5
}

func (v *Vector) arrayFor(i int) []interface{} {
	if i >= v.count || i < 0 {
		panic(StubNewError(fmt.Sprintf("Index %d is out of bounds [0..%d]", i, v.count-1)))
	}
	if i >= v.tailoff() {
		return v.tail
	}
	node := v.root
	for level := v.shift; level > 0; level -= 5 {
		node = node[(i>>level)&0x01F].([]interface{})
	}
	return node
}

func (v *Vector) at(i int) Object {
	return v.arrayFor(i)[i&0x01F].(Object)
}

func newPath(level uint, node []interface{}) []interface{} {
	if level == 0 {
		return node
	}
	result := make([]interface{}, 32)
	result[0] = newPath(level-5, node)
	return result
}

func (v *Vector) pushTail(level uint, parent []interface{}, tailNode []interface{}) []interface{} {
	subidx := ((v.count - 1) >> level) & 0x01F
	result := clone(parent)
	var nodeToInsert []interface{}
	if level == 5 {
		nodeToInsert = tailNode
	} else {
		if parent[subidx] != nil {
			nodeToInsert = v.pushTail(level-5, parent[subidx].([]interface{}), tailNode)
		} else {
			nodeToInsert = newPath(level-5, tailNode)
		}
	}
	result[subidx] = nodeToInsert
	return result
}

func (v *Vector) Conjoin(obj Object) (*Vector, error) {
	var newTail []interface{}
	if v.count-v.tailoff() < 32 {
		newTail = append(clone(v.tail), obj)
		return &Vector{count: v.count + 1, shift: v.shift, root: v.root, tail: newTail}, nil
	}
	var newRoot []interface{}
	newShift := v.shift
	if (v.count >> 5) > (1 << v.shift) {
		newRoot = make([]interface{}, 32)
		newRoot[0] = v.root
		newRoot[1] = newPath(v.shift, v.tail)
		newShift += 5
	} else {
		newRoot = v.pushTail(v.shift, v.root, v.tail)
	}
	newTail = make([]interface{}, 1, 32)
	newTail[0] = obj
	return &Vector{count: v.count + 1, shift: newShift, root: newRoot, tail: newTail}, nil
}

func (v *Vector) ToString(env *Env, escape bool) (string, error) {
	var b bytes.Buffer
	b.WriteRune('[')
	if v.count > 0 {
		for i := 0; i < v.count-1; i++ {
			s, err := ToString(env, v.at(i))
			if err != nil {
				return "", err
			}
			b.WriteString(s)
			b.WriteRune(' ')
		}
		s, err := ToString(env, v.at(v.count-1))
		if err != nil {
			return "", err
		}
		b.WriteString(s)
	}
	b.WriteRune(']')
	return b.String(), nil
}

func (v *Vector) Equals(env *Env, other interface{}) bool {
	if v == other {
		return true
	}
	return IsSeqEqual(env, v.Seq(), other)
}

func (v *Vector) GetType() *Type {
	return TYPE.Vector
}

func (v *Vector) Hash(env *Env) (uint32, error) {
	return hashOrdered(env, v.Seq())
}

func (seq *VectorSeq) Seq() Seq {
	return seq
}

func (vseq *VectorSeq) Equals(env *Env, other interface{}) bool {
	return IsSeqEqual(env, vseq, other)
}

func (vseq *VectorSeq) ToString(env *Env, escape bool) (string, error) {
	return SeqToString(env, vseq, escape)
}

func (seq *VectorSeq) Pprint(env *Env, w io.Writer, indent int) (int, error) {
	return pprintSeq(env, seq, w, indent)
}

func (vseq *VectorSeq) WithMeta(env *Env, meta Map) (Object, error) {
	res := *vseq
	m, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return &res, nil
}

func (vseq *VectorSeq) GetType() *Type {
	return TYPE.VectorSeq
}

func (vseq *VectorSeq) Hash(env *Env) (uint32, error) {
	return hashOrdered(env, vseq)
}

func (vseq *VectorSeq) First(env *Env) (Object, error) {
	if vseq.index < vseq.vector.count {
		return vseq.vector.at(vseq.index), nil
	}
	return NIL, nil
}

func (vseq *VectorSeq) Rest(env *Env) (Seq, error) {
	if vseq.index+1 < vseq.vector.count {
		return &VectorSeq{vector: vseq.vector, index: vseq.index + 1}, nil
	}
	return EmptyList, nil
}

func (vseq *VectorSeq) IsEmpty(env *Env) (bool, error) {
	return vseq.index >= vseq.vector.count, nil
}

func (vseq *VectorSeq) Cons(obj Object) Seq {
	return &ConsSeq{first: obj, rest: vseq}
}

func (vseq *VectorSeq) sequential() {}

func (seq *VectorRSeq) Seq() Seq {
	return seq
}

func (vseq *VectorRSeq) Equals(env *Env, other interface{}) bool {
	return IsSeqEqual(env, vseq, other)
}

func (vseq *VectorRSeq) ToString(env *Env, escape bool) (string, error) {
	return SeqToString(env, vseq, escape)
}

func (seq *VectorRSeq) Pprint(env *Env, w io.Writer, indent int) (int, error) {
	return pprintSeq(env, seq, w, indent)
}

func (vseq *VectorRSeq) WithMeta(env *Env, meta Map) (Object, error) {
	res := *vseq
	m, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return &res, nil
}

func (vseq *VectorRSeq) GetType() *Type {
	return TYPE.VectorRSeq
}

func (vseq *VectorRSeq) Hash(env *Env) (uint32, error) {
	return hashOrdered(env, vseq)
}

func (vseq *VectorRSeq) First(env *Env) (Object, error) {
	if vseq.index >= 0 {
		return vseq.vector.at(vseq.index), nil
	}
	return NIL, nil
}

func (vseq *VectorRSeq) Rest(env *Env) (Seq, error) {
	if vseq.index-1 >= 0 {
		return &VectorRSeq{vector: vseq.vector, index: vseq.index - 1}, nil
	}
	return EmptyList, nil
}

func (vseq *VectorRSeq) IsEmpty(env *Env) (bool, error) {
	return vseq.index < 0, nil
}

func (vseq *VectorRSeq) Cons(obj Object) Seq {
	return &ConsSeq{first: obj, rest: vseq}
}

func (vseq *VectorRSeq) sequential() {}

func (v *Vector) Seq() Seq {
	return &VectorSeq{vector: v, index: 0}
}

func (v *Vector) Conj(env *Env, obj Object) (Conjable, error) {
	return v.Conjoin(obj)
}

func (v *Vector) Count() int {
	return v.count
}

func (v *Vector) Nth(env *Env, i int) (Object, error) {
	return v.at(i), nil
}

func (v *Vector) TryNth(env *Env, i int, d Object) (Object, error) {
	if i < 0 || i >= v.count {
		return d, nil
	}
	return v.at(i), nil
}

func (v *Vector) sequential() {}

func (v *Vector) Compare(env *Env, other Object) (int, error) {
	v2, err := AssertVector(env, other, "Cannot compare Vector and "+TypeName(other))
	if err != nil {
		return 0, err
	}
	if v.Count() > v2.Count() {
		return 1, nil
	}
	if v.Count() < v2.Count() {
		return -1, nil
	}
	for i := 0; i < v.Count(); i++ {
		v, err := AssertComparable(env, v.at(i), "")
		if err != nil {
			return 0, err
		}
		c, err := v.Compare(env, v2.at(i))
		if err != nil {
			return 0, err
		}
		if c != 0 {
			return c, nil
		}
	}
	return 0, nil
}

func (v *Vector) Peek(env *Env) (Object, error) {
	if v.count > 0 {
		return v.Nth(env, v.count-1)
	}
	return NIL, nil
}

func (v *Vector) popTail(level uint, node []interface{}) []interface{} {
	subidx := ((v.count - 2) >> level) & 0x01F
	if level > 5 {
		newChild := v.popTail(level-5, node[subidx].([]interface{}))
		if newChild == nil && subidx == 0 {
			return nil
		} else {
			ret := clone(node)
			ret[subidx] = newChild
			return ret
		}
	} else if subidx == 0 {
		return nil
	} else {
		ret := clone(node)
		ret[subidx] = nil
		return ret
	}
}

func (v *Vector) Pop(env *Env) (Stack, error) {
	if v.count == 0 {
		return nil, env.NewError("Can't pop empty vector")
	}
	if v.count == 1 {
		return EmptyVectorWithMeta(v.meta), nil
	}
	if v.count-v.tailoff() > 1 {
		newTail := clone(v.tail)[0 : len(v.tail)-1]
		res := &Vector{count: v.count - 1, shift: v.shift, root: v.root, tail: newTail}
		res.meta = v.meta
		return res, nil
	}
	newTail := v.arrayFor(v.count - 2)
	newRoot := v.popTail(v.shift, v.root)
	newShift := v.shift
	if newRoot == nil {
		newRoot = empty_node
	}
	if v.shift > 5 && newRoot[1] == nil {
		newRoot = newRoot[0].([]interface{})
		newShift -= 5
	}
	res := &Vector{count: v.count - 1, shift: newShift, root: newRoot, tail: newTail}
	res.meta = v.meta
	return res, nil
}

func (v *Vector) Get(env *Env, key Object) (bool, Object, error) {
	switch key := key.(type) {
	case Int:
		if key.I() >= 0 && key.I() < v.count {
			return true, v.at(key.I()), nil
		}
	}
	return false, nil, nil
}

func (v *Vector) EntryAt(env *Env, key Object) (*Vector, error) {
	ok, val, err := v.Get(env, key)
	if err != nil {
		return nil, err
	}

	if ok {
		return NewVectorFrom(key, val), nil
	}
	return nil, nil
}

func doAssoc(level uint, node []interface{}, i int, val Object) []interface{} {
	ret := clone(node)
	if level == 0 {
		ret[i&0x01f] = val
	} else {
		subidx := (i >> level) & 0x01f
		ret[subidx] = doAssoc(level-5, node[subidx].([]interface{}), i, val)
	}
	return ret
}

func (v *Vector) assocN(i int, val Object) (*Vector, error) {
	if i < 0 || i > v.count {
		return nil, StubNewError((fmt.Sprintf("Index %d is out of bounds [0..%d]", i, v.count)))
	}
	if i == v.count {
		return v.Conjoin(val)
	}
	if i < v.tailoff() {
		res := &Vector{count: v.count, shift: v.shift, root: doAssoc(v.shift, v.root, i, val), tail: v.tail}
		res.meta = v.meta
		return res, nil
	}
	newTail := clone(v.tail)
	newTail[i&0x01f] = val
	res := &Vector{count: v.count, shift: v.shift, root: v.root, tail: newTail}
	res.meta = v.meta
	return res, nil
}

func assertInteger(obj Object) (int, error) {
	var i int
	switch obj := obj.(type) {
	case Int:
		i = obj.I()
	case *BigInt:
		i = obj.Int().I()
	default:
		return 0, StubNewError("Key must be integer")
	}
	return i, nil
}

func (v *Vector) Assoc(env *Env, key, val Object) (Associative, error) {
	i, err := assertInteger(key)
	if err != nil {
		return nil, err
	}
	return v.assocN(i, val)
}

func (v *Vector) Rseq() Seq {
	return &VectorRSeq{vector: v, index: v.count - 1}
}

func (v *Vector) Call(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	i, err := assertInteger(args[0])
	if err != nil {
		return nil, err
	}
	return v.at(i), nil
}

var _ Callable = (*Vector)(nil)

func EmptyVector() *Vector {
	return &Vector{
		count: 0,
		shift: 5,
		root:  empty_node,
		tail:  make([]interface{}, 0, 32),
	}
}

func EmptyVectorWithMeta(m Map) *Vector {
	v := &Vector{
		count: 0,
		shift: 5,
		root:  empty_node,
		tail:  make([]interface{}, 0, 32),
	}
	v.meta = m

	return v
}

func NewVectorFrom(objs ...Object) *Vector {
	res := EmptyVector()
	for i := 0; i < len(objs); i++ {
		res, _ = res.Conjoin(objs[i])
	}
	return res
}

func NewVectorFromSeq(env *Env, seq Seq) (*Vector, error) {
	res := EmptyVector()
	for {
		ok, err := seq.IsEmpty(env)
		if err != nil {
			return nil, err
		}

		if ok {
			break
		}

		v, err := seq.First(env)
		if err != nil {
			return nil, err
		}
		res, _ = res.Conjoin(v)
		seq, err = seq.Rest(env)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (v *Vector) Empty() Collection {
	return EmptyVector()
}

func (v *Vector) kvreduce(env *Env, c Callable, init Object) (Object, error) {
	res := init
	for i := 0; i < v.Count(); i++ {
		o, err := v.Nth(env, i)
		if err != nil {
			return nil, err
		}
		res, err = c.Call(env, []Object{res, MakeInt(i), o})
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (v *Vector) Pprint(env *Env, w io.Writer, indent int) (int, error) {
	ind := indent + 1
	var err error
	fmt.Fprint(w, "[")
	if v.count > 0 {
		for i := 0; i < v.count-1; i++ {
			_, err = pprintObject(env, v.at(i), indent+1, w)
			if err != nil {
				return 0, err
			}
			fmt.Fprint(w, "\n")
			err = writeIndent(w, indent+1)
			if err != nil {
				return 0, err
			}
		}
		ind, err = pprintObject(env, v.at(v.count-1), indent+1, w)
		if err != nil {
			return 0, err
		}
	}
	fmt.Fprint(w, "]")
	return ind + 1, nil
}

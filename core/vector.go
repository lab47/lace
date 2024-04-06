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

func (v *Vector) WithMeta(meta Map) (Object, error) {
	res := *v
	m, err := SafeMerge(res.meta, meta)
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

func (v *Vector) ToString(escape bool) string {
	var b bytes.Buffer
	b.WriteRune('[')
	if v.count > 0 {
		for i := 0; i < v.count-1; i++ {
			b.WriteString(v.at(i).ToString(escape))
			b.WriteRune(' ')
		}
		b.WriteString(v.at(v.count - 1).ToString(escape))
	}
	b.WriteRune(']')
	return b.String()
}

func (v *Vector) Equals(other interface{}) bool {
	if v == other {
		return true
	}
	return IsSeqEqual(v.Seq(), other)
}

func (v *Vector) GetType() *Type {
	return TYPE.Vector
}

func (v *Vector) Hash() uint32 {
	return hashOrdered(v.Seq())
}

func (seq *VectorSeq) Seq() Seq {
	return seq
}

func (vseq *VectorSeq) Equals(other interface{}) bool {
	return IsSeqEqual(vseq, other)
}

func (vseq *VectorSeq) ToString(escape bool) string {
	return SeqToString(vseq, escape)
}

func (seq *VectorSeq) Pprint(w io.Writer, indent int) int {
	return pprintSeq(seq, w, indent)
}

func (vseq *VectorSeq) WithMeta(meta Map) (Object, error) {
	res := *vseq
	m, err := SafeMerge(res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return &res, nil
}

func (vseq *VectorSeq) GetType() *Type {
	return TYPE.VectorSeq
}

func (vseq *VectorSeq) Hash() uint32 {
	return hashOrdered(vseq)
}

func (vseq *VectorSeq) First() Object {
	if vseq.index < vseq.vector.count {
		return vseq.vector.at(vseq.index)
	}
	return NIL
}

func (vseq *VectorSeq) Rest() Seq {
	if vseq.index+1 < vseq.vector.count {
		return &VectorSeq{vector: vseq.vector, index: vseq.index + 1}
	}
	return EmptyList
}

func (vseq *VectorSeq) IsEmpty() bool {
	return vseq.index >= vseq.vector.count
}

func (vseq *VectorSeq) Cons(obj Object) Seq {
	return &ConsSeq{first: obj, rest: vseq}
}

func (vseq *VectorSeq) sequential() {}

func (seq *VectorRSeq) Seq() Seq {
	return seq
}

func (vseq *VectorRSeq) Equals(other interface{}) bool {
	return IsSeqEqual(vseq, other)
}

func (vseq *VectorRSeq) ToString(escape bool) string {
	return SeqToString(vseq, escape)
}

func (seq *VectorRSeq) Pprint(w io.Writer, indent int) int {
	return pprintSeq(seq, w, indent)
}

func (vseq *VectorRSeq) WithMeta(meta Map) (Object, error) {
	res := *vseq
	m, err := SafeMerge(res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return &res, nil
}

func (vseq *VectorRSeq) GetType() *Type {
	return TYPE.VectorRSeq
}

func (vseq *VectorRSeq) Hash() uint32 {
	return hashOrdered(vseq)
}

func (vseq *VectorRSeq) First() Object {
	if vseq.index >= 0 {
		return vseq.vector.at(vseq.index)
	}
	return NIL
}

func (vseq *VectorRSeq) Rest() Seq {
	if vseq.index-1 >= 0 {
		return &VectorRSeq{vector: vseq.vector, index: vseq.index - 1}
	}
	return EmptyList
}

func (vseq *VectorRSeq) IsEmpty() bool {
	return vseq.index < 0
}

func (vseq *VectorRSeq) Cons(obj Object) Seq {
	return &ConsSeq{first: obj, rest: vseq}
}

func (vseq *VectorRSeq) sequential() {}

func (v *Vector) Seq() Seq {
	return &VectorSeq{vector: v, index: 0}
}

func (v *Vector) Conj(obj Object) (Conjable, error) {
	return v.Conjoin(obj)
}

func (v *Vector) Count() int {
	return v.count
}

func (v *Vector) Nth(i int) Object {
	return v.at(i)
}

func (v *Vector) TryNth(i int, d Object) Object {
	if i < 0 || i >= v.count {
		return d
	}
	return v.at(i)
}

func (v *Vector) sequential() {}

func (v *Vector) Compare(env *Env, other Object) (int, error) {
	v2, err := AssertVector(env, other, "Cannot compare Vector and "+other.GetType().ToString(false))
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

func (v *Vector) Peek() Object {
	if v.count > 0 {
		return v.Nth(v.count - 1)
	}
	return NIL
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

func (v *Vector) Pop() Stack {
	if v.count == 0 {
		panic(StubNewError("Can't pop empty vector"))
	}
	if v.count == 1 {
		return EmptyVectorWithMeta(v.meta).(Stack)
	}
	if v.count-v.tailoff() > 1 {
		newTail := clone(v.tail)[0 : len(v.tail)-1]
		res := &Vector{count: v.count - 1, shift: v.shift, root: v.root, tail: newTail}
		res.meta = v.meta
		return res
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
	return res
}

func (v *Vector) Get(key Object) (bool, Object) {
	switch key := key.(type) {
	case Int:
		if key.I >= 0 && key.I < v.count {
			return true, v.at(key.I)
		}
	}
	return false, nil
}

func (v *Vector) EntryAt(key Object) (*Vector, error) {
	ok, val := v.Get(key)
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
		i = obj.I
	case *BigInt:
		i = obj.Int().I
	default:
		return 0, StubNewError("Key must be integer")
	}
	return i, nil
}

func (v *Vector) Assoc(key, val Object) (Associative, error) {
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
	CheckArity(env, args, 1, 1)
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

func EmptyVectorWithMeta(m Map) Object {
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

func NewVectorFromSeq(seq Seq) *Vector {
	res := EmptyVector()
	for !seq.IsEmpty() {
		res, _ = res.Conjoin(seq.First())
		seq = seq.Rest()
	}
	return res
}

func (v *Vector) Empty() Collection {
	return EmptyVector()
}

func (v *Vector) kvreduce(env *Env, c Callable, init Object) (Object, error) {
	res := init
	var err error
	for i := 0; i < v.Count(); i++ {
		res, err = c.Call(env, []Object{res, Int{I: i}, v.Nth(i)})
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (v *Vector) Pprint(w io.Writer, indent int) int {
	ind := indent + 1
	fmt.Fprint(w, "[")
	if v.count > 0 {
		for i := 0; i < v.count-1; i++ {
			pprintObject(v.at(i), indent+1, w)
			fmt.Fprint(w, "\n")
			writeIndent(w, indent+1)
		}
		ind = pprintObject(v.at(v.count-1), indent+1, w)
	}
	fmt.Fprint(w, "]")
	return ind + 1
}

package core

import "io"

type (
	Box struct {
		val interface{}
	}
	Node interface {
		assoc(env *Env, shift uint, hash uint32, key Object, val Object, addedLeaf *Box) (Node, error)
		without(env *Env, shift uint, hash uint32, key Object) Node
		find(env *Env, shift uint, hash uint32, key Object) *Pair
		findEqu(shift uint, hash uint32, key Equ) *Pair
		nodeSeq() Seq
		iter() MapIterator
	}
	HashMap struct {
		InfoHolder
		MetaHolder
		count int
		root  Node
	}
	BitmapIndexedNode struct {
		bitmap int
		array  []any
	}
	HashCollisionNode struct {
		hash  uint32
		count int
		array []any
	}
	ArrayNode struct {
		count int
		array []Node
	}
	NodeSeq struct {
		InfoHolder
		MetaHolder
		array []any
		i     int
		s     Seq
	}
	ArrayNodeSeq struct {
		InfoHolder
		MetaHolder
		nodes []Node
		i     int
		s     Seq
	}
	NodeIterator struct {
		array     []any
		i         int
		nextEntry *Pair
		nextIter  MapIterator
	}
	ArrayNodeIterator struct {
		array      []Node
		i          int
		nestedIter MapIterator
	}
)

var (
	emptyIndexedNode = &BitmapIndexedNode{}
	EmptyHashMap     = &HashMap{}
)

func (iter *ArrayNodeIterator) HasNext() bool {
	for {
		if iter.nestedIter != nil {
			if iter.nestedIter.HasNext() {
				return true
			} else {
				iter.nestedIter = nil
			}
		}
		if iter.i < len(iter.array) {
			node := iter.array[iter.i]
			iter.i++
			if node != nil {
				iter.nestedIter = node.iter()
			}
		} else {
			return false
		}
	}
}

func (iter *ArrayNodeIterator) Next() *Pair {
	if iter.HasNext() {
		return iter.nestedIter.Next()
	}
	panic(newIteratorError())
}

func (iter *NodeIterator) advance() bool {
	for iter.i < len(iter.array) {
		key := iter.array[iter.i]
		nodeOrVal := iter.array[iter.i+1]
		iter.i += 2
		if key != nil {
			iter.nextEntry = &Pair{Key: key.(Object), Value: nodeOrVal.(Object)}
			return true
		} else if nodeOrVal != nil {
			iter1 := nodeOrVal.(Node).iter()
			if iter1 != nil && iter1.HasNext() {
				iter.nextIter = iter1
				return true
			}
		}
	}
	return false
}

func (iter *NodeIterator) HasNext() bool {
	if iter.nextEntry != nil || iter.nextIter != nil {
		return true
	}
	return iter.advance()
}

func (iter *NodeIterator) Next() *Pair {
	ret := iter.nextEntry
	if ret != nil {
		iter.nextEntry = nil
		return ret
	} else if iter.nextIter != nil {
		ret := iter.nextIter.Next()
		if !iter.nextIter.HasNext() {
			iter.nextIter = nil
		}
		return ret
	} else if iter.advance() {
		return iter.Next()
	}
	panic(newIteratorError())
}

func newArrayNodeSeq(nodes []Node, i int, s Seq) Seq {
	if s != nil {
		return &ArrayNodeSeq{
			nodes: nodes,
			i:     i,
			s:     s,
		}
	}
	for j := i; j < len(nodes); j++ {
		if nodes[j] != nil {
			ns := nodes[j].nodeSeq()
			if ns != nil {
				return &ArrayNodeSeq{
					nodes: nodes,
					i:     j + 1,
					s:     ns,
				}
			}
		}
	}
	return nil
}

func (s *ArrayNodeSeq) WithMeta(env *Env, meta Map) (Object, error) {
	res := *s
	m, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return &res, nil
}

func (s *ArrayNodeSeq) Seq() Seq {
	return s
}

func (s *ArrayNodeSeq) Equals(env *Env, other interface{}) bool {
	return IsSeqEqual(env, s, other)
}

func (s *ArrayNodeSeq) ToString(env *Env, escape bool) (string, error) {
	return SeqToString(env, s, escape)
}

func (seq *ArrayNodeSeq) Pprint(env *Env, w io.Writer, indent int) (int, error) {
	return pprintSeq(env, seq, w, indent)
}

func (s *ArrayNodeSeq) GetType() *Type {
	return TYPE.ArrayNodeSeq
}

func (s *ArrayNodeSeq) Hash(env *Env) (uint32, error) {
	return hashOrdered(env, s)
}

func (s *ArrayNodeSeq) First(env *Env) (Object, error) {
	return s.s.First(env)
}

func (s *ArrayNodeSeq) Rest(env *Env) Seq {
	next := s.s.Rest(env)
	if next.IsEmpty(env) {
		next = nil
	}
	res := newArrayNodeSeq(s.nodes, s.i, next)
	if res == nil {
		return EmptyList
	}
	return res
}

func (s *ArrayNodeSeq) IsEmpty(env *Env) bool {
	if s.s != nil {
		return s.s.IsEmpty(env)
	}
	return false
}

func (s *ArrayNodeSeq) Cons(obj Object) Seq {
	return &ConsSeq{first: obj, rest: s}
}

func (s *ArrayNodeSeq) sequential() {}

func newNodeSeq(array []any, i int, s Seq) Seq {
	if s != nil {
		return &NodeSeq{
			array: array,
			i:     i,
			s:     s,
		}
	}
	for j := i; j < len(array); j += 2 {
		if array[j] != nil {
			return &NodeSeq{
				array: array,
				i:     j,
			}
		}
		switch node := array[j+1].(type) {
		case Node:
			nodeSeq := node.nodeSeq()
			if nodeSeq != nil {
				return &NodeSeq{
					array: array,
					i:     j + 2,
					s:     nodeSeq,
				}
			}
		}
	}
	return nil
}

func (s *NodeSeq) WithMeta(env *Env, meta Map) (Object, error) {
	res := *s
	m, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}

	s.meta = m
	return &res, nil
}

func (s *NodeSeq) Seq() Seq {
	return s
}

func (s *NodeSeq) Equals(env *Env, other interface{}) bool {
	return IsSeqEqual(env, s, other)
}

func (s *NodeSeq) ToString(env *Env, escape bool) (string, error) {
	return SeqToString(env, s, escape)
}

func (seq *NodeSeq) Pprint(env *Env, w io.Writer, indent int) (int, error) {
	return pprintSeq(env, seq, w, indent)
}

func (s *NodeSeq) GetType() *Type {
	return TYPE.NodeSeq
}

func (s *NodeSeq) Hash(env *Env) (uint32, error) {
	return hashOrdered(env, s)
}

func (s *NodeSeq) First(env *Env) (Object, error) {
	if s.s != nil {
		return s.s.First(env)
	}
	return NewVectorFrom(s.array[s.i].(Object), s.array[s.i+1].(Object)), nil
}

func (s *NodeSeq) Rest(env *Env) Seq {
	var res Seq
	if s.s != nil {
		next := s.s.Rest(env)
		if next.IsEmpty(env) {
			next = nil
		}
		res = newNodeSeq(s.array, s.i, next)
	} else {
		res = newNodeSeq(s.array, s.i+2, nil)
	}
	if res == nil {
		return EmptyList
	}
	return res
}

func (s *NodeSeq) IsEmpty(env *Env) bool {
	if s.s != nil {
		return s.s.IsEmpty(env)
	}
	return false
}

func (s *NodeSeq) Cons(obj Object) Seq {
	return &ConsSeq{first: obj, rest: s}
}

func (s *NodeSeq) sequential() {}

func (n *ArrayNode) iter() MapIterator {
	return &ArrayNodeIterator{
		array: n.array,
	}
}

func (n *ArrayNode) assoc(env *Env, shift uint, hash uint32, key Object, val Object, addedLeaf *Box) (Node, error) {
	idx := mask(hash, shift)
	node := n.array[idx]
	if node == nil {
		nn, err := emptyIndexedNode.assoc(env, shift+5, hash, key, val, addedLeaf)
		if err != nil {
			return nil, err
		}
		return &ArrayNode{
			count: n.count + 1,
			array: cloneAndSetNode(n.array, int(idx), nn),
		}, nil
	}
	nn, err := node.assoc(env, shift+5, hash, key, val, addedLeaf)
	if err != nil {
		return nil, err
	}
	if nn == node {
		return n, nil
	}
	return &ArrayNode{
		count: n.count,
		array: cloneAndSetNode(n.array, int(idx), nn),
	}, nil
}

func (n *ArrayNode) without(env *Env, shift uint, hash uint32, key Object) Node {
	idx := mask(hash, shift)
	node := n.array[idx]
	if node == nil {
		return n
	}
	nn := node.without(env, shift+5, hash, key)
	if nn == node {
		return n
	}
	if nn == nil {
		if n.count <= 8 {
			return n.pack(uint(idx))
		}
		return &ArrayNode{
			count: n.count - 1,
			array: cloneAndSetNode(n.array, int(idx), nn),
		}
	} else {
		return &ArrayNode{
			count: n.count,
			array: cloneAndSetNode(n.array, int(idx), nn),
		}
	}
}

func (n *ArrayNode) find(env *Env, shift uint, hash uint32, key Object) *Pair {
	idx := mask(hash, shift)
	node := n.array[idx]
	if node == nil {
		return nil
	}
	return node.find(env, shift+5, hash, key)
}

func (n *ArrayNode) findEqu(shift uint, hash uint32, key Equ) *Pair {
	idx := mask(hash, shift)
	node := n.array[idx]
	if node == nil {
		return nil
	}
	return node.findEqu(shift+5, hash, key)
}

func (n *ArrayNode) nodeSeq() Seq {
	return newArrayNodeSeq(n.array, 0, nil)
}

func (n *ArrayNode) pack(idx uint) Node {
	newArray := make([]any, 2*(n.count-1))
	j := 1
	bitmap := 0
	var i uint
	for i = 0; i < idx; i++ {
		if n.array[i] != nil {
			newArray[j] = n.array[i]
			bitmap |= 1 << i
			j += 2
		}
	}
	for i = idx + 1; i < uint(len(n.array)); i++ {
		if n.array[i] != nil {
			newArray[j] = n.array[i]
			bitmap |= 1 << i
			j += 2
		}
	}
	return &BitmapIndexedNode{
		bitmap: bitmap,
		array:  newArray,
	}
}

func (n *HashCollisionNode) findIndex(env *Env, key Object) int {
	for i := 0; i < 2*n.count; i += 2 {
		if key.Equals(env, n.array[i]) {
			return i
		}
	}
	return -1
}

func (n *HashCollisionNode) findIndexEqu(key Equ) int {
	for i := 0; i < 2*n.count; i += 2 {
		k, ok := n.array[i].(Object)
		if !ok {
			continue
		}

		if key.Is(k) {
			return i
		}
	}
	return -1
}

func (n *HashCollisionNode) iter() MapIterator {
	return &NodeIterator{
		array: n.array,
	}
}

func (n *HashCollisionNode) assoc(env *Env, shift uint, hash uint32, key Object, val Object, addedLeaf *Box) (Node, error) {
	if hash == n.hash {
		idx := n.findIndex(env, key)
		if idx != -1 {
			if n.array[idx+1] == val {
				return n, nil
			}
			return &HashCollisionNode{
				hash:  hash,
				count: n.count,
				array: cloneAndSet(n.array, idx+1, val),
			}, nil
		}
		newArray := make([]interface{}, 2*(n.count+1))
		for i := 0; i < 2*n.count; i++ {
			newArray[i] = n.array[i]
		}
		newArray[2*n.count] = key
		newArray[2*n.count+1] = val
		addedLeaf.val = addedLeaf
		return &HashCollisionNode{
			hash:  hash,
			count: n.count + 1,
			array: newArray,
		}, nil
	}

	return (&BitmapIndexedNode{
		bitmap: bitpos(n.hash, shift),
		array:  []interface{}{nil, n},
	}).assoc(env, shift, hash, key, val, addedLeaf)
}

func (n *HashCollisionNode) without(env *Env, shift uint, hash uint32, key Object) Node {
	idx := n.findIndex(env, key)
	if idx == -1 {
		return n
	}
	if n.count == 1 {
		return nil
	}
	return &HashCollisionNode{
		hash:  hash,
		count: n.count - 1,
		array: removePair(n.array, idx/2),
	}
}

func (n *HashCollisionNode) find(env *Env, shift uint, hash uint32, key Object) *Pair {
	idx := n.findIndex(env, key)
	if idx == -1 {
		return nil
	}
	return &Pair{
		Key:   n.array[idx].(Object),
		Value: n.array[idx+1].(Object),
	}
}

func (n *HashCollisionNode) findEqu(shift uint, hash uint32, key Equ) *Pair {
	idx := n.findIndexEqu(key)
	if idx == -1 {
		return nil
	}
	return &Pair{
		Key:   n.array[idx].(Object),
		Value: n.array[idx+1].(Object),
	}
}

func (n *HashCollisionNode) nodeSeq() Seq {
	return newNodeSeq(n.array, 0, nil)
}

func bitCount(n int) int {
	var count int
	for n != 0 {
		count++
		n &= n - 1
	}
	return count
}

func mask(hash uint32, shift uint) uint32 {
	return (hash >> shift) & 0x01f
}

func bitpos(hash uint32, shift uint) int {
	return 1 << mask(hash, shift)
}

func cloneAndSet(array []interface{}, i int, a interface{}) []interface{} {
	res := clone(array)
	res[i] = a
	return res
}

func cloneAndSet2(array []interface{}, i int, a interface{}, j int, b interface{}) []interface{} {
	res := clone(array)
	res[i] = a
	res[j] = b
	return res
}

func cloneAndSetNode(array []Node, i int, a Node) []Node {
	res := make([]Node, len(array), cap(array))
	copy(res, array)
	res[i] = a
	return res
}

func createNode(env *Env, shift uint, key1 Object, val1 Object, key2hash uint32, key2 Object, val2 Object) (Node, error) {
	key1hash, err := key1.Hash(env)
	if err != nil {
		return nil, err
	}
	if key1hash == key2hash {
		return &HashCollisionNode{
			hash:  key1hash,
			count: 2,
			array: []interface{}{key1, val1, key2, val2},
		}, nil
	}
	addedLeaf := &Box{}
	n, err := emptyIndexedNode.assoc(env, shift, key1hash, key1, val1, addedLeaf)
	if err != nil {
		return nil, err
	}

	return n.assoc(env, shift, key2hash, key2, val2, addedLeaf)
}

func removePair(array []any, n int) []interface{} {
	newArray := make([]interface{}, len(array)-2)
	for i := 0; i < 2*n; i++ {
		newArray[i] = array[i]
	}
	for i := 2 * (n + 1); i < len(array); i++ {
		newArray[i-2] = array[i]
	}
	return newArray
}

func (b *BitmapIndexedNode) index(bit int) int {
	return bitCount(b.bitmap & (bit - 1))
}

func (b *BitmapIndexedNode) iter() MapIterator {
	return &NodeIterator{
		array: b.array,
	}
}

func (b *BitmapIndexedNode) assoc(env *Env, shift uint, hash uint32, key Object, val Object, addedLeaf *Box) (Node, error) {
	bit := bitpos(hash, shift)
	idx := b.index(bit)
	if b.bitmap&bit != 0 {
		keyOrNull := b.array[2*idx]
		valOrNode := b.array[2*idx+1]
		if keyOrNull == nil {
			n, err := valOrNode.(Node).assoc(env, shift+5, hash, key, val, addedLeaf)
			if err != nil {
				return nil, err
			}
			if n == valOrNode {
				return b, nil
			}
			return &BitmapIndexedNode{
				bitmap: b.bitmap,
				array:  cloneAndSet(b.array, 2*idx+1, n),
			}, nil
		}
		if key.Equals(env, keyOrNull) {
			if val == valOrNode {
				return b, nil
			}
			return &BitmapIndexedNode{
				bitmap: b.bitmap,
				array:  cloneAndSet(b.array, 2*idx+1, val),
			}, nil
		}
		addedLeaf.val = addedLeaf
		nn, err := createNode(env, shift+5, keyOrNull.(Object), valOrNode.(Object), hash, key, val)
		if err != nil {
			return nil, err
		}

		return &BitmapIndexedNode{
			bitmap: b.bitmap,
			array:  cloneAndSet2(b.array, 2*idx, nil, 2*idx+1, nn),
		}, nil
	} else {
		n := bitCount(b.bitmap)
		var err error
		if n >= 16 {
			nodes := make([]Node, 32)
			jdx := mask(hash, shift)
			nodes[jdx], err = emptyIndexedNode.assoc(env, shift+5, hash, key, val, addedLeaf)
			if err != nil {
				return nil, err
			}
			j := 0
			var i uint
			for i = 0; i < 32; i++ {
				if (b.bitmap>>i)&1 != 0 {
					if b.array[j] == nil {
						nodes[i] = b.array[j+1].(Node)
					} else {
						h, err := b.array[j].(Object).Hash(env)
						if err != nil {
							return nil, err
						}
						nodes[i], err = emptyIndexedNode.assoc(env, shift+5, h, b.array[j].(Object), b.array[j+1].(Object), addedLeaf)
						if err != nil {
							return nil, err
						}
					}
					j += 2
				}
			}
			return &ArrayNode{
				count: n + 1,
				array: nodes,
			}, nil
		} else {
			newArray := make([]interface{}, 2*(n+1))
			for i := 0; i < 2*idx; i++ {
				newArray[i] = b.array[i]
			}
			newArray[2*idx] = key
			addedLeaf.val = addedLeaf
			newArray[2*idx+1] = val
			for i := 2 * idx; i < 2*n; i++ {
				newArray[i+2] = b.array[i]
			}
			return &BitmapIndexedNode{
				bitmap: b.bitmap | bit,
				array:  newArray,
			}, nil
		}
	}
}

func (b *BitmapIndexedNode) without(env *Env, shift uint, hash uint32, key Object) Node {
	bit := bitpos(hash, shift)
	if (b.bitmap & bit) == 0 {
		return b
	}
	idx := b.index(bit)
	keyOrNull := b.array[2*idx]
	valOrNode := b.array[2*idx+1]
	if keyOrNull == nil {
		n := valOrNode.(Node).without(env, shift+5, hash, key)
		if n == valOrNode {
			return b
		}
		if n != nil {
			return &BitmapIndexedNode{
				bitmap: b.bitmap,
				array:  cloneAndSet(b.array, 2*idx+1, n),
			}
		}
		if b.bitmap == bit {
			return nil
		}
		return &BitmapIndexedNode{
			bitmap: b.bitmap ^ bit,
			array:  removePair(b.array, idx),
		}
	}
	if key.Equals(env, keyOrNull) {
		return &BitmapIndexedNode{
			bitmap: b.bitmap ^ bit,
			array:  removePair(b.array, idx),
		}
	}
	return b
}

func (b *BitmapIndexedNode) find(env *Env, shift uint, hash uint32, key Object) *Pair {
	bit := bitpos(hash, shift)
	if (b.bitmap & bit) == 0 {
		return nil
	}
	idx := b.index(bit)
	keyOrNull := b.array[2*idx]
	valOrNode := b.array[2*idx+1]
	if keyOrNull == nil {
		return valOrNode.(Node).find(env, shift+5, hash, key)
	}
	if key.Equals(env, keyOrNull) {
		return &Pair{
			Key:   keyOrNull.(Object),
			Value: valOrNode.(Object),
		}
	}
	return nil
}

func (b *BitmapIndexedNode) findEqu(shift uint, hash uint32, key Equ) *Pair {
	bit := bitpos(hash, shift)
	if (b.bitmap & bit) == 0 {
		return nil
	}
	idx := b.index(bit)
	keyOrNull := b.array[2*idx]
	valOrNode := b.array[2*idx+1]
	if keyOrNull == nil {
		return valOrNode.(Node).findEqu(shift+5, hash, key)
	}
	obj, ok := keyOrNull.(Object)
	if ok && key.Is(obj) {
		return &Pair{
			Key:   keyOrNull.(Object),
			Value: valOrNode.(Object),
		}
	}
	return nil
}

func (b *BitmapIndexedNode) nodeSeq() Seq {
	return newNodeSeq(b.array, 0, nil)
}

func (m *HashMap) WithMeta(env *Env, meta Map) (Object, error) {
	res := *m
	v, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = v
	return &res, nil
}

func (m *HashMap) ToString(env *Env, escape bool) (string, error) {
	return mapToString(env, m, escape)
}

func (m *HashMap) Equals(env *Env, other interface{}) bool {
	return mapEquals(env, m, other)
}

func (m *HashMap) GetType() *Type {
	return TYPE.HashMap
}

func (m *HashMap) Hash(env *Env) (uint32, error) {
	return hashUnordered(env, m.Seq(), 1)
}

func (m *HashMap) Seq() Seq {
	if m.root != nil {
		s := m.root.nodeSeq()
		if s != nil {
			return s
		}
	}
	return EmptyList
}

func (m *HashMap) Count() int {
	return m.count
}

func (m *HashMap) containsKey(env *Env, key Object) bool {
	if m.root != nil {
		h, err := key.Hash(env)
		if err != nil {
			return false
		}

		return m.root.find(env, 0, h, key) != nil
	} else {
		return false
	}
}

func (m *HashMap) Assoc(env *Env, key, val Object) (Associative, error) {
	addedLeaf := &Box{}
	var newroot, t Node
	if m.root == nil {
		t = emptyIndexedNode
	} else {
		t = m.root
	}
	h, err := key.Hash(env)
	if err != nil {
		return nil, err
	}
	newroot, err = t.assoc(env, 0, h, key, val, addedLeaf)
	if err != nil {
		return nil, err
	}
	if newroot == m.root {
		return m, nil
	}
	newcount := m.count
	if addedLeaf.val != nil {
		newcount = m.count + 1
	}
	res := &HashMap{
		count: newcount,
		root:  newroot,
	}
	res.meta = m.meta
	return res, nil
}

func (m *HashMap) EntryAt(env *Env, key Object) (*Vector, error) {
	if m.root != nil {
		h, err := key.Hash(env)
		if err != nil {
			return nil, err
		}
		p := m.root.find(env, 0, h, key)
		if p != nil {
			return NewVectorFrom(p.Key, p.Value), nil
		}
	}
	return nil, nil
}

func (m *HashMap) Get(env *Env, key Object) (bool, Object, error) {
	if m.root != nil {
		h, err := key.Hash(env)
		if err != nil {
			return false, nil, err
		}
		if res := m.root.find(env, 0, h, key); res != nil {
			return true, res.Value, nil
		}
	}
	return false, nil, nil
}

func (m *HashMap) GetEqu(key Equ) (bool, Object) {
	if m.root != nil {
		if res := m.root.findEqu(0, key.IsHash(), key); res != nil {
			return true, res.Value
		}
	}
	return false, nil
}

func (m *HashMap) Conj(env *Env, obj Object) (Conjable, error) {
	return mapConj(env, m, obj)
}

func (m *HashMap) Iter() MapIterator {
	if m.root == nil {
		return emptyMapIterator
	}
	return m.root.iter()
}

func (m *HashMap) Keys() Seq {
	return &MappingSeq{
		seq: m.Seq(),
		fn: func(env *Env, obj Object) (Object, error) {
			var v *Vector
			if err := Cast(env, obj, &v); err != nil {
				return nil, err
			}
			return v.Nth(env, 0)
		},
	}
}

func (m *HashMap) Vals() Seq {
	return &MappingSeq{
		seq: m.Seq(),
		fn: func(env *Env, obj Object) (Object, error) {
			var v *Vector
			if err := Cast(env, obj, &v); err != nil {
				return nil, err
			}
			return v.Nth(env, 1)
		},
	}
}

func (m *HashMap) Merge(env *Env, other Map) (Map, error) {
	if other.Count() == 0 {
		return m, nil
	}
	if m.Count() == 0 {
		return other, nil
	}
	var res Associative = m
	var err error
	for iter := other.Iter(); iter.HasNext(); {
		p := iter.Next()
		res, err = res.Assoc(env, p.Key, p.Value)
		if err != nil {
			return nil, err
		}
	}
	var mm Map
	if err := Cast(env, res, &mm); err != nil {
		return nil, err
	}

	return mm, nil
}

func (m *HashMap) Without(env *Env, key Object) (Map, error) {
	if m.root == nil {
		return m, nil
	}

	h, err := key.Hash(env)
	if err != nil {
		return nil, err
	}

	newroot := m.root.without(env, 0, h, key)
	if newroot == m.root {
		return m, nil
	}
	res := &HashMap{
		count: m.count - 1,
		root:  newroot,
	}
	res.meta = m.meta
	return res, nil
}

func (m *HashMap) Call(env *Env, args []Object) (Object, error) {
	return callMap(env, m, args)
}

var _ Callable = (*HashMap)(nil)

func NewHashMap(env *Env, keyvals ...Object) (*HashMap, error) {
	var res Associative = EmptyHashMap
	var err error
	for i := 0; i < len(keyvals); i += 2 {
		res, err = res.Assoc(env, keyvals[i], keyvals[i+1])
		if err != nil {
			return nil, err
		}
	}
	var hm *HashMap
	if err := Cast(env, res, &hm); err != nil {
		return nil, err
	}
	return hm, nil
}

func (m *HashMap) Empty() Collection {
	return EmptyHashMap
}

func (m *HashMap) Pprint(env *Env, w io.Writer, indent int) (int, error) {
	return pprintMap(env, m, w, indent)
}

package marshal

import (
	"errors"
	"math/big"
	"reflect"

	"github.com/fxamacker/cbor/v2"
	"github.com/lab47/lace/core"
	"github.com/oklog/ulid/v2"
)

type Ref struct {
	Index int `json:"ref" cbor:"1,keyasint"`
}

type Symbol struct {
	Value string `json:"value" cbor:"1,keyasint"`
}

type Keyword struct {
	Value string `json:"value" cbor:"1,keyasint"`
}

type Char struct {
	Value string `json:"value" cbor:"1,keyasint"`
}

type BigFloat struct {
	Value string `json:"value" cbor:"1,keyasint"`
}

type Ratio struct {
	X *big.Int `json:"x" cbor:"1,keyasint"`
	Y *big.Int `json:"y" cbor:"2,keyasint"`
}

type Regex struct {
	Value string `json:"value" cbor:"1,keyasint"`
}

type ForeignRef struct {
	Ref string `json:"ref" cbor:"1,keyasint"`
}

type Collection struct {
	Type   string            `json:"type" cbor:"1,keyasint"`
	Values []cbor.RawMessage `json:"values" cbor:"2,keyasint"`
	Ref    int               `json:"ref" cbor:"3,keyasint"`
}

type Map struct {
	Values []cbor.RawMessage `json:"map" cbor:"1,keyasint"`
	Ref    int               `json:"ref" cbor:"2,keyasint"`
}

var (
	encoder cbor.EncMode
	decoder cbor.DecMode
)

func init() {
	to := cbor.TagOptions{
		DecTag: cbor.DecTagRequired,
		EncTag: cbor.EncTagRequired,
	}

	ts := cbor.NewTagSet()
	ts.Add(to, reflect.TypeFor[Symbol](), 470)
	ts.Add(to, reflect.TypeFor[Keyword](), 471)
	ts.Add(to, reflect.TypeFor[Char](), 472)
	ts.Add(to, reflect.TypeFor[BigFloat](), 473)
	ts.Add(to, reflect.TypeFor[Ratio](), 474)
	ts.Add(to, reflect.TypeFor[Regex](), 475)
	ts.Add(to, reflect.TypeFor[ForeignRef](), 476)
	ts.Add(to, reflect.TypeFor[Collection](), 477)
	ts.Add(to, reflect.TypeFor[Ref](), 480)

	m, err := cbor.EncOptions{}.EncModeWithTags(ts)
	if err != nil {
		panic(err)
	}

	encoder = m

	n, err := cbor.DecOptions{
		IntDec: cbor.IntDecConvertSignedOrBigInt,
	}.DecModeWithTags(ts)
	if err != nil {
		panic(err)
	}

	decoder = n
}

type marshalState struct {
	env   *core.Env
	refs  map[core.Object]int
	frefs map[string]core.Object
}

func (m *marshalState) newFref(val core.Object) ForeignRef {
	u, err := ulid.New(ulid.Now(), ulid.DefaultEntropy())
	if err != nil {
		panic(err)
	}

	str := u.String()

	m.frefs[str] = val

	return ForeignRef{Ref: str}
}

func (m *marshalState) newRef(val core.Object) Ref {
	id := len(m.refs)

	m.refs[val] = id

	return Ref{id}
}

func (m *marshalState) encode(val any) ([]byte, error) {
	return encoder.Marshal(val)
}

func (m *marshalState) encodeSeq(typ string, sq core.Seqable) ([]byte, error) {
	s := sq.Seq()

	ref := m.newRef(s)

	elems, err := core.ToSlice(m.env, s)
	if err != nil {
		return nil, err
	}

	col := Collection{
		Type: typ,
		Ref:  ref.Index,
	}

	for _, el := range elems {
		d, err := m.Marshal(el)
		if err != nil {
			return nil, err
		}

		col.Values = append(col.Values, cbor.RawMessage(d))
	}

	return m.encode(col)
}

func (m *marshalState) encodeSet(ma core.Set) ([]byte, error) {
	ref := m.newRef(ma)

	col := Collection{
		Type: "set",
		Ref:  ref.Index,
	}

	ret := map[any]any{}
	iter := ma.SetIter()
	for iter.HasNext(m.env) {
		p, err := iter.Next(m.env)
		if err != nil {
			return nil, err
		}

		k, err := m.Marshal(p)
		if err != nil {
			return nil, err
		}

		col.Values = append(col.Values, cbor.RawMessage(k))
	}

	return m.encode(ret)
}

func (m *marshalState) encodeMap(ma core.Map) ([]byte, error) {
	ref := m.newRef(ma)

	col := Collection{
		Type: "map",
		Ref:  ref.Index,
	}

	iter := ma.Iter()
	for iter.HasNext() {
		p := iter.Next()
		k, err := m.Marshal(p.Key)
		if err != nil {
			return nil, err
		}

		v, err := m.Marshal(p.Value)
		if err != nil {
			return nil, err
		}

		col.Values = append(col.Values, cbor.RawMessage(k), cbor.RawMessage(v))
	}

	return m.encode(col)
}

func (m *marshalState) Marshal(obj core.Object) ([]byte, error) {
	if idx, ok := m.refs[obj]; ok {
		return m.encode(Ref{Index: idx})
	}

	switch sv := obj.(type) {
	case core.Symbol:
		return m.encode(Symbol{Value: sv.String()})
	case core.Keyword:
		return m.encode(Keyword{Value: sv.RawString()})
	case core.Int:
		return m.encode(sv.I())
	case core.Double:
		return m.encode(sv.D)
	case *core.BigInt:
		return m.encode(sv.BigInt())
	case *core.BigFloat:
		return m.encode(BigFloat{Value: sv.BigFloat().String()})
	case *core.Ratio:
		r := sv.Ratio()

		return m.encode(Ratio{
			X: r.Num(),
			Y: r.Denom(),
		})
	case core.Char:
		return m.encode(Char{Value: string(sv.Ch())})
	case core.Boolean:
		return m.encode(bool(sv))
	case core.Nil:
		return m.encode(nil)
	case core.String:
		return m.encode(sv.S())
	case *core.Regex:
		return m.encode(Regex{Value: sv.R.String()})
	case core.Time:
		return m.encode(sv.T)
	case *core.Var:
		return m.encode(m.newFref(sv))
	case *core.Fn:
		return m.encode(m.newFref(sv))
	case *core.List:
		return m.encodeSeq("list", sv)
	case *core.Vector:
		return m.encodeSeq("vector", sv)
	case core.Set:
		return m.encodeSet(sv)
	case core.Map:
		return m.encodeMap(sv)
	case core.Seqable:
		return m.encodeSeq("list", sv)
	default:
		return nil, ErrUnsupported
	}
}

var ErrUnsupported = errors.New("value can not be marshaled")

func Marshal(obj core.Object) ([]byte, error) {
	var ms marshalState
	ms.refs = make(map[core.Object]int)
	ms.frefs = make(map[string]core.Object)

	return ms.Marshal(obj)
}

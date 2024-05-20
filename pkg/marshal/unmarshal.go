package marshal

import (
	"fmt"
	"math/big"
	"regexp"
	"time"

	"github.com/lab47/lace/core"
)

type unmarshalState struct {
	env   *core.Env
	refs  map[int]any
	frefs map[string]any
}

func Unmarshal(env *core.Env, data []byte) (any, error) {
	var ms unmarshalState

	ms.env = env
	ms.refs = make(map[int]any)

	return ms.unmarshal(data)
}

func (s *unmarshalState) unmarshalCol(col Collection) (any, error) {
	var (
		elems []any
		err   error
	)

	for _, e := range col.Values {
		obj, err := s.unmarshal(e)
		if err != nil {
			return nil, err
		}

		elems = append(elems, obj)
	}

	var ret any

	switch col.Type {
	case "list":
		ret = core.NewListFrom(elems...)
	case "vector":
		ret = core.NewVectorFrom(elems...)
	case "set":
		seq := core.NewListFrom(elems...)
		s, err := core.NewSetFromSeq(s.env, seq)
		if err != nil {
			return nil, err
		}

		ret = s
	case "map":
		ret, err = core.NewHashMap(s.env, elems...)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown collection type: %s", col.Type)
	}

	s.refs[col.Ref] = ret

	return ret, nil
}

func (s *unmarshalState) unmarshal(data []byte) (any, error) {
	var out any

	err := decoder.Unmarshal(data, &out)
	if err != nil {
		return nil, err
	}

	switch sv := out.(type) {
	case Ref:
		return s.refs[sv.Index], nil
	case Symbol:
		return core.MakeSymbol(sv.Value), nil
	case Keyword:
		return core.MakeKeyword(sv.Value), nil
	case int64:
		return core.MakeInt(int(sv)), nil
	case float64:
		return core.MakeDouble(sv), nil
	case bool:
		return core.MakeBoolean(sv), nil
	case *big.Int:
		return core.MakeBigIntFrom(sv), nil
	case BigFloat:
		var bf big.Float
		if _, ok := bf.SetPrec(256).SetString(sv.Value); !ok {
			return nil, err
		}
		return core.MakeBigFloatFrom(&bf), nil
	case Ratio:
		return core.MakeRatio(sv.X, sv.Y), nil
	case nil:
		return core.NIL, nil
	case string:
		return core.MakeString(sv), nil
	case time.Time:
		return core.MakeTime(sv), nil
	case Regex:
		re, err := regexp.Compile(sv.Value)
		if err != nil {
			return nil, err
		}

		return core.MakeRegex(re), nil
	case Collection:
		return s.unmarshalCol(sv)
	default:
		return nil, fmt.Errorf("unsupported type: %T", sv)
	}
}

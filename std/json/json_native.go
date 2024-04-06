package json

import (
	"encoding/json"
	"fmt"

	. "github.com/lab47/lace/core"
)

func fromObject(obj Object) interface{} {
	switch obj := obj.(type) {
	case Keyword:
		return obj.ToString(false)[1:]
	case Boolean:
		return obj.B
	case Number:
		return obj.Double().D
	case Nil:
		return nil
	case *Vector:
		cnt := obj.Count()
		res := make([]interface{}, cnt)
		for i := 0; i < cnt; i++ {
			res[i] = fromObject(obj.Nth(i))
		}
		return res
	case Map:
		res := make(map[string]interface{})
		for iter := obj.Iter(); iter.HasNext(); {
			p := iter.Next()
			var k string
			switch p.Key.(type) {
			case Keyword:
				k = p.Key.ToString(false)[1:]
			default:
				k = p.Key.ToString(false)
			}
			res[k] = fromObject(p.Value)
		}
		return res
	default:
		return obj.ToString(false)
	}
}

func toObject(v interface{}, keywordize bool) (Object, error) {
	switch v := v.(type) {
	case string:
		return MakeString(v), nil
	case float64:
		return Double{D: v}, nil
	case bool:
		return Boolean{B: v}, nil
	case nil:
		return NIL, nil
	case []interface{}:
		res := EmptyVector()
		for _, v := range v {
			o, err := toObject(v, keywordize)
			if err != nil {
				return nil, err
			}
			res, err = res.Conjoin(o)
			if err != nil {
				return nil, err
			}
		}
		return res, nil
	case map[string]interface{}:
		res := EmptyArrayMap()
		for k, v := range v {
			var key Object
			if keywordize {
				key = MakeKeyword(k)
			} else {
				key = MakeString(k)
			}
			o, err := toObject(v, keywordize)
			if err != nil {
				return nil, err
			}
			res.Add(key, o)
		}
		return res, nil
	default:
		return nil, StubNewError(fmt.Sprintf("Unknown json value: %v", v))
	}
}

func readString(s string, opts Map) (Object, error) {
	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return nil, StubNewError("Invalid json: " + err.Error())
	}
	var keywordize bool
	if opts != nil {
		if ok, v := opts.Get(MakeKeyword("keywords?")); ok {
			keywordize = ToBool(v)
		}
	}
	return toObject(v, keywordize)
}

func writeString(obj Object) (String, error) {
	res, err := json.Marshal(fromObject(obj))
	if err != nil {
		return String{}, StubNewError("Cannot encode value to json: " + err.Error())
	}
	return String{S: string(res)}, nil
}

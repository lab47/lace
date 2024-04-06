package yaml

import (
	"fmt"

	"gopkg.in/yaml.v2"

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

func toObject(v interface{}) (Object, error) {
	switch v := v.(type) {
	case string:
		return MakeString(v), nil
	case float64:
		return Double{D: v}, nil
	case int:
		return Int{I: v}, nil
	case bool:
		return Boolean{B: v}, nil
	case nil:
		return NIL, nil
	case []interface{}:
		res := EmptyVector()
		for _, v := range v {
			o, err := toObject(v)
			if err != nil {
				return nil, err
			}
			res, err = res.Conjoin(o)
			if err != nil {
				return nil, err
			}
		}
		return res, nil
	case map[interface{}]interface{}:
		res := EmptyArrayMap()
		for k, v := range v {
			ko, err := toObject(k)
			if err != nil {
				return nil, err
			}
			vo, err := toObject(v)
			if err != nil {
				return nil, err
			}
			res.Add(ko, vo)
		}
		return res, nil
	default:
		panic(StubNewError(fmt.Sprintf("Unknown yaml value: %v", v)))
	}
}

func readString(s string) (Object, error) {
	var v interface{}
	if err := yaml.Unmarshal([]byte(s), &v); err != nil {
		panic(StubNewError("Invalid yaml: " + err.Error()))
	}
	return toObject(v)
}

func writeString(obj Object) (String, error) {
	res, err := yaml.Marshal(fromObject(obj))
	if err != nil {
		return String{}, StubNewError("Cannot encode value to yaml: " + err.Error())
	}
	return String{S: string(res)}, nil
}

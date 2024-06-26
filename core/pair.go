package core

import "fmt"

type (
	NamedPair struct {
		Key   any
		Value any
	}
)

var _ any = &NamedPair{}

func MakeNamedPair(key any, obj any) *NamedPair {
	return &NamedPair{
		Key:   key,
		Value: obj,
	}
}

func (p *NamedPair) ToString(env *Env, escape bool) (string, error) {
	var (
		ks  string
		err error
	)

	if p.Key == nil {
		ks = "nil"
	} else {
		ks, err = ToString(env, p.Key)
		if err != nil {
			return "", err
		}
	}

	vs, err := ToString(env, p.Value)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"#object[NamedPair Key=%s Val=%s]", ks, vs,
	), nil
}

func (p *NamedPair) Equals(env *Env, other interface{}) bool {
	if e, ok := other.(*NamedPair); ok {
		if !Equals(env, p.Key, e.Key) {
			return false
		}

		return Equals(env, p.Value, e.Value)
	} else {
		return false
	}
}

func (p *NamedPair) Hash(env *Env) (uint32, error) {
	kh, err := HashValue(env, p.Key)
	if err != nil {
		return 0, err
	}

	vh, err := HashValue(env, p.Value)
	if err != nil {
		return 0, err
	}

	return kh ^ vh, nil
}

func (p *NamedPair) WithInfo(info *ObjectInfo) any {
	return p
}

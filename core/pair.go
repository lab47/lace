package core

import "fmt"

type (
	NamedPair struct {
		Key   Object
		Value Object
	}
)

var _ Object = &NamedPair{}

func MakeNamedPair(key Object, obj Object) *NamedPair {
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
		ks, err = p.Key.ToString(env, escape)
		if err != nil {
			return "", err
		}
	}

	vs, err := p.Value.ToString(env, escape)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"#object[NamedPair Key=%s Val=%s]", ks, vs,
	), nil
}

func (p *NamedPair) Equals(env *Env, other interface{}) bool {
	if e, ok := other.(*NamedPair); ok {
		if !p.Key.Equals(env, e.Key) {
			return false
		}

		return p.Value.Equals(env, e.Value)
	} else {
		return false
	}
}

func (p *NamedPair) GetInfo() *ObjectInfo {
	return nil
}

func (p *NamedPair) GetType() *Type {
	return TYPE.NamedPair
}

func (p *NamedPair) Hash(env *Env) (uint32, error) {
	kh, err := p.Key.Hash(env)
	if err != nil {
		return 0, err
	}

	vh, err := p.Value.Hash(env)
	if err != nil {
		return 0, err
	}

	return kh ^ vh, nil
}

func (p *NamedPair) WithInfo(info *ObjectInfo) Object {
	return p
}
package core

import "strings"

type Keyword interface {
	Object
	Equ
	Comparable
	Callable
	HasInfo

	Name() string
	Namespace() string
	String() string
	RawString() string

	keywordType() string
}

const KeywordHashMask uint32 = 0x7334c790

func MakeKeyword(nsname string) Keyword {
	index := strings.IndexRune(nsname, '/')
	if index == -1 || nsname == "/" {
		name := nsname
		return NewKeyword("", name)
	}
	ns := nsname[0:index]
	name := nsname[index+1:]
	return NewKeyword(ns, name)
}

func NewKeyword(ns, name string) Keyword {
	if ns == "" {
		return TinyKeyword(name)
	}
	return &HeavyKeyword{
		ns:   ns,
		name: name,
	}
}

type HeavyKeyword struct {
	InfoHolder
	ns   string
	name string
}

func (*HeavyKeyword) keywordType() string { return "heavy" }

func (k *HeavyKeyword) WithInfo(info *ObjectInfo) Object {
	k.info = info
	return k
}

func (k *HeavyKeyword) ToString(env *Env, escape bool) (string, error) {
	if k.ns != "" {
		return ":" + k.ns + "/" + k.name, nil
	}
	return ":" + k.name, nil
}

func (k *HeavyKeyword) String() string {
	if k.ns != "" {
		return ":" + k.ns + "/" + k.name
	}
	return ":" + k.name
}

func (k *HeavyKeyword) RawString() string {
	if k.ns != "" {
		return k.ns + "/" + k.name
	}
	return k.name
}

func (k *HeavyKeyword) Name() string {
	return k.name
}

func (k *HeavyKeyword) Namespace() string {
	if k.ns != "" {
		return k.ns
	}
	return ""
}

func (k *HeavyKeyword) Equals(env *Env, other interface{}) bool {
	switch other := other.(type) {
	case Keyword:
		return k.ns == other.Namespace() && k.name == other.Name()
	default:
		return false
	}
}

func (k *HeavyKeyword) Is(other Object) bool {
	switch other := other.(type) {
	case Keyword:
		return k.ns == other.Namespace() && k.name == other.Name()
	default:
		return false
	}
}

func (k *HeavyKeyword) GetType() *Type {
	return TYPE.Keyword
}

func (k *HeavyKeyword) Hash(env *Env) (uint32, error) {
	return hashSymbol(k.ns, k.name) ^ KeywordHashMask, nil
}

func (k *HeavyKeyword) IsHash() uint32 {
	return hashSymbol(k.ns, k.name) ^ KeywordHashMask
}

func (k *HeavyKeyword) Compare(env *Env, other Object) (int, error) {
	k2, err := AssertKeyword(env, other, "Cannot compare Keyword and "+TypeName(other))
	if err != nil {
		return 0, err
	}

	ks, err := ToString(env, k)
	if err != nil {
		return 0, err
	}
	k2s, err := ToString(env, k2)
	if err != nil {
		return 0, err
	}
	return strings.Compare(ks, k2s), nil
}

func (k *HeavyKeyword) Call(env *Env, args []Object) (Object, error) {
	return getMap(env, k, args)
}

var _ Callable = &HeavyKeyword{}

type TinyKeyword string

func (TinyKeyword) keywordType() string { return "tiny" }

func (TinyKeyword) GetInfo() *ObjectInfo {
	return nil
}

func (k TinyKeyword) WithInfo(info *ObjectInfo) Object {
	hk := &HeavyKeyword{
		name: string(k),
	}
	hk.info = info
	return hk
}

func (k TinyKeyword) ToString(env *Env, escape bool) (string, error) {
	return ":" + string(k), nil
}

func (k TinyKeyword) String() string {
	return ":" + string(k)
}

func (k TinyKeyword) RawString() string {
	return string(k)
}

func (k TinyKeyword) Name() string {
	return string(k)
}

func (k TinyKeyword) Namespace() string {
	return ""
}

func (k TinyKeyword) Equals(env *Env, other interface{}) bool {
	switch other := other.(type) {
	case Keyword:
		return other.Namespace() == "" && string(k) == other.Name()
	default:
		return false
	}
}

func (k TinyKeyword) Is(other Object) bool {
	switch other := other.(type) {
	case Keyword:
		return other.Namespace() == "" && string(k) == other.Name()
	default:
		return false
	}
}

func (k TinyKeyword) GetType() *Type {
	return TYPE.Keyword
}

func (k TinyKeyword) Hash(env *Env) (uint32, error) {
	return hashSymbol("", string(k)) ^ KeywordHashMask, nil
}

func (k TinyKeyword) IsHash() uint32 {
	return hashSymbol("", string(k)) ^ KeywordHashMask
}

func (k TinyKeyword) Compare(env *Env, other Object) (int, error) {
	k2, err := AssertKeyword(env, other, "Cannot compare Keyword and "+TypeName(other))
	if err != nil {
		return 0, err
	}

	ks, err := ToString(env, k)
	if err != nil {
		return 0, err
	}
	k2s, err := ToString(env, k2)
	if err != nil {
		return 0, err
	}
	return strings.Compare(ks, k2s), nil
}

func (k TinyKeyword) Call(env *Env, args []Object) (Object, error) {
	return getMap(env, k, args)
}

var _ Callable = TinyKeyword("")

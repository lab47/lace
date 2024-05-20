package core

import (
	"strings"
)

type HeavySymbol struct {
	InfoHolder
	MetaHolder
	ns   string
	name string
	hash uint32
}

type Symbol interface {
	Object
	Equ
	Meta
	Comparable

	Name() string
	Namespace() string
	String() string

	symbolType() string
}

func (h *HeavySymbol) symbolType() string { return "heavy" }
func (h *LightSymbol) symbolType() string { return "light" }

func SymbolSetInfo(sym Symbol, info *ObjectInfo) Symbol {
	if hs, ok := sym.(*HeavySymbol); ok {
		hs.info = info
		return hs
	}

	return sym
}

func MakeTaggedSymbol(nsname string, tag Symbol) Symbol {
	var sym HeavySymbol

	index := strings.IndexRune(nsname, '/')
	if index == -1 || nsname == "/" {
		sym = HeavySymbol{
			name: nsname,
		}
	} else {
		sym = HeavySymbol{
			ns:   nsname[0:index],
			name: nsname[index+1:],
		}
	}

	m := EmptyArrayMap()
	m.AddEqu(criticalKeywords.tag, tag)

	sym.meta = m

	return &sym
}

func hashSymbol(ns, name string) uint32 {
	h := getHash()
	if ns != "" {
		h.Write([]byte(ns))
	}
	h.Write([]byte("/" + name))
	return h.Sum32()
}

func AssembleSymbol(ns, name string) Symbol {
	if ns == "" {
		return TinySymbol(name)
	}
	return &LightSymbol{
		ns:   ns,
		name: name,
	}
}

func MakeSymbol(nsname string) Symbol {
	index := strings.IndexRune(nsname, '/')
	if index == -1 || nsname == "/" {
		return TinySymbol(nsname)
	}
	return &LightSymbol{
		ns:   nsname[0:index],
		name: nsname[index+1:],
	}
}

func MakeHeavySymbol(nsname string) *HeavySymbol {
	index := strings.IndexRune(nsname, '/')
	if index == -1 || nsname == "/" {
		return &HeavySymbol{
			name: nsname,
		}
	}
	return &HeavySymbol{
		ns:   nsname[0:index],
		name: nsname[index+1:],
	}
}

func MakeSymbolWithMeta(nsname string, m Map) Symbol {
	index := strings.IndexRune(nsname, '/')
	var sym HeavySymbol
	if index == -1 || nsname == "/" {
		sym = HeavySymbol{
			name: nsname,
		}
	} else {
		sym = HeavySymbol{
			ns:   nsname[0:index],
			name: nsname[index+1:],
		}
	}

	sym.meta = m

	return &sym
}

func IsSymbol(obj Object) bool {
	switch obj.(type) {
	case Symbol:
		return true
	default:
		return false
	}
}

func (sym *HeavySymbol) WithMeta(env *Env, meta Map) (Object, error) {
	res := sym
	m, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return res, nil
}

func (s *HeavySymbol) ToString(env *Env, escape bool) (string, error) {
	if s.ns != "" {
		return s.ns + "/" + s.name, nil
	}
	return s.name, nil
}

func (s *HeavySymbol) String() string {
	if s.ns != "" {
		return s.ns + "/" + s.name
	}
	return s.name
}

func (s *HeavySymbol) Name() string {
	return s.name
}

func (s *HeavySymbol) Namespace() string {
	if s.ns != "" {
		return s.ns
	}
	return ""
}

func (s *HeavySymbol) Equals(env *Env, other interface{}) bool {
	switch other := other.(type) {
	case Symbol:
		return s.ns == other.Namespace() && s.name == other.Name()
	default:
		return false
	}
}

func (s *HeavySymbol) Is(other Object) bool {
	switch other := other.(type) {
	case Symbol:
		return s.ns == other.Namespace() && s.name == other.Name()
	default:
		return false
	}
}

func (s *HeavySymbol) GetType() *Type {
	return TYPE.Symbol
}

func (s *HeavySymbol) Hash(env *Env) (uint32, error) {
	return s.IsHash(), nil
}

func (s *HeavySymbol) IsHash() uint32 {
	return hashSymbol(s.ns, s.name) + 0x9e3779b9
}

func (s *HeavySymbol) Compare(env *Env, other Object) (int, error) {
	s2, err := AssertSymbol(env, other, "Cannot compare Symbol and "+TypeName(other))
	if err != nil {
		return 0, err
	}

	ks, err := ToString(env, s)
	if err != nil {
		return 0, err
	}

	k2s, err := ToString(env, s2)
	if err != nil {
		return 0, err
	}
	return strings.Compare(ks, k2s), nil
}

func (s *HeavySymbol) Call(env *Env, args []Object) (Object, error) {
	return getMap(env, s, args)
}

var _ Callable = &HeavySymbol{}

type LightSymbol struct {
	ns, name string
}

var _ Symbol = &LightSymbol{}

func (sym *LightSymbol) WithMeta(env *Env, meta Map) (Object, error) {
	res := &HeavySymbol{
		ns:   sym.ns,
		name: sym.name,
	}
	m, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return res, nil
}

func (sym *LightSymbol) GetMeta() Map {
	return nil
}

func (sym *LightSymbol) WithInfo(info *ObjectInfo) Object {
	res := &HeavySymbol{
		ns:   sym.ns,
		name: sym.name,
	}
	res.info = info

	return res
}

func (sym *LightSymbol) GetInfo() *ObjectInfo {
	return nil
}

func (s *LightSymbol) ToString(env *Env, escape bool) (string, error) {
	if s.ns != "" {
		return s.ns + "/" + s.name, nil
	}
	return s.name, nil
}

func (s *LightSymbol) String() string {
	if s.ns != "" {
		return s.ns + "/" + s.name
	}
	return s.name
}

func (s *LightSymbol) Name() string {
	return s.name
}

func (s *LightSymbol) Namespace() string {
	if s.ns != "" {
		return s.ns
	}
	return ""
}

func (s *LightSymbol) Equals(env *Env, other interface{}) bool {
	switch other := other.(type) {
	case Symbol:
		return s.ns == other.Namespace() && s.name == other.Name()
	default:
		return false
	}
}

func (s *LightSymbol) Is(other Object) bool {
	switch other := other.(type) {
	case Symbol:
		return s.ns == other.Namespace() && s.name == other.Name()
	default:
		return false
	}
}

func (s *LightSymbol) GetType() *Type {
	return TYPE.Symbol
}

func (s *LightSymbol) Hash(env *Env) (uint32, error) {
	return s.IsHash(), nil
}

func (s *LightSymbol) IsHash() uint32 {
	return hashSymbol(s.ns, s.name) + 0x9e3779b9
}

func (s *LightSymbol) Compare(env *Env, other Object) (int, error) {
	s2, err := AssertSymbol(env, other, "Cannot compare Symbol and "+TypeName(other))
	if err != nil {
		return 0, err
	}

	ks, err := ToString(env, s)
	if err != nil {
		return 0, err
	}

	k2s, err := ToString(env, s2)
	if err != nil {
		return 0, err
	}
	return strings.Compare(ks, k2s), nil
}

func (s *LightSymbol) Call(env *Env, args []Object) (Object, error) {
	return getMap(env, s, args)
}

type TinySymbol string

var _ Symbol = TinySymbol("")

func (TinySymbol) symbolType() string { return "tiny" }

func (sym TinySymbol) WithMeta(env *Env, meta Map) (Object, error) {
	res := &HeavySymbol{
		ns:   "",
		name: string(sym),
	}
	m, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return res, nil
}

func (sym TinySymbol) GetMeta() Map {
	return nil
}

func (sym TinySymbol) WithInfo(info *ObjectInfo) Object {
	res := &HeavySymbol{
		ns:   "",
		name: string(sym),
	}
	res.info = info

	return res
}

func (sym TinySymbol) GetInfo() *ObjectInfo {
	return nil
}

func (s TinySymbol) ToString(env *Env, escape bool) (string, error) {
	return string(s), nil
}

func (s TinySymbol) String() string {
	return string(s)
}

func (s TinySymbol) Name() string {
	return string(s)
}

func (s TinySymbol) Namespace() string {
	return ""
}

func (s TinySymbol) Equals(env *Env, other interface{}) bool {
	switch other := other.(type) {
	case Symbol:
		return other.Namespace() == "" && string(s) == other.Name()
	default:
		return false
	}
}

func (s TinySymbol) Is(other Object) bool {
	switch other := other.(type) {
	case Symbol:
		return other.Namespace() == "" && string(s) == other.Name()
	default:
		return false
	}
}

func (s TinySymbol) GetType() *Type {
	return TYPE.Symbol
}

func (s TinySymbol) Hash(env *Env) (uint32, error) {
	return s.IsHash(), nil
}

func (s TinySymbol) IsHash() uint32 {
	return hashSymbol("", string(s)) + 0x9e3779b9
}

func (s TinySymbol) Compare(env *Env, other Object) (int, error) {
	s2, err := AssertSymbol(env, other, "Cannot compare Symbol and "+TypeName(other))
	if err != nil {
		return 0, err
	}

	ks, err := ToString(env, s)
	if err != nil {
		return 0, err
	}

	k2s, err := ToString(env, s2)
	if err != nil {
		return 0, err
	}
	return strings.Compare(ks, k2s), nil
}

func (s TinySymbol) Call(env *Env, args []Object) (Object, error) {
	return getMap(env, s, args)
}

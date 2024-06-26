package core

const MAX_RUNE = int(^uint32(0) >> 1)
const MIN_RUNE = -MAX_RUNE - 1

// A single unicode rune.
//
//lace:export
type Char interface {
	any
	Comparable

	Ch() rune

	charType() string
}

func NewChar(ch rune) Char {
	return TinyChar(ch)
}

type HeavyChar struct {
	InfoHolder
	ch rune
}

func (x *HeavyChar) charType() string { return "heavy" }

func (x *HeavyChar) Ch() rune {
	return x.ch
}

func (x *HeavyChar) WithInfo(info *ObjectInfo) any {
	x.info = info
	return x
}

func (c *HeavyChar) ToString(env *Env, escape bool) (string, error) {
	if escape {
		return escapeRune(c.Ch()), nil
	}
	return string(c.Ch()), nil
}

func (c *HeavyChar) Equals(env *Env, other interface{}) bool {
	switch other := other.(type) {
	case Char:
		return c.Ch() == other.Ch()
	default:
		return false
	}
}

func (c *HeavyChar) Native() interface{} {
	return c.Ch
}

func (c *HeavyChar) Hash(env *Env) (uint32, error) {
	h := getHash()
	h.Write([]byte(string(c.Ch())))
	return h.Sum32(), nil
}

func (c *HeavyChar) Compare(env *Env, other any) (int, error) {
	c2, err := AssertChar(env, other, "Cannot compare Char and "+TypeName(other))
	if err != nil {
		return 0, err
	}
	if c.Ch() < c2.Ch() {
		return -1, nil
	}
	if c2.Ch() < c.Ch() {
		return 1, nil
	}
	return 0, nil
}

type TinyChar rune

func (x TinyChar) charType() string { return "tiny" }

func (x TinyChar) Ch() rune {
	return rune(x)
}

func (x TinyChar) WithInfo(info *ObjectInfo) any {
	r := &HeavyChar{ch: rune(x)}
	r.info = info
	return r
}

func (c TinyChar) ToString(env *Env, escape bool) (string, error) {
	if escape || ToBool(env.printReadably.GetStatic()) {
		return escapeRune(c.Ch()), nil
	}
	return string(c.Ch()), nil
}

func (c TinyChar) Equals(env *Env, other interface{}) bool {
	switch other := other.(type) {
	case Char:
		return c.Ch() == other.Ch()
	default:
		return false
	}
}

func (c TinyChar) Native() interface{} {
	return c.Ch
}

func (c TinyChar) Hash(env *Env) (uint32, error) {
	h := getHash()
	h.Write([]byte(string(c.Ch())))
	return h.Sum32(), nil
}

func (c TinyChar) Compare(env *Env, other any) (int, error) {
	c2, err := AssertChar(env, other, "Cannot compare Char and "+TypeName(other))
	if err != nil {
		return 0, err
	}
	if c.Ch() < c2.Ch() {
		return -1, nil
	}
	if c2.Ch() < c.Ch() {
		return 1, nil
	}
	return 0, nil
}

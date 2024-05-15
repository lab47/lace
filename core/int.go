package core

import (
	"encoding/binary"
	"fmt"
)

func (i Int) ToString(env *Env, escape bool) (string, error) {
	return fmt.Sprintf("%d", i.Int()), nil
}

func MakeInt(i int) Int {
	return Int(i)
}

func (i Int) I() int {
	return int(i)
}

func (i Int) GetInfo() *ObjectInfo {
	return nil
}

func (i Int) Equals(env *Env, other interface{}) bool {
	return equalsNumbers(i, other)
}

func (i Int) GetType() *Type {
	return TYPE.Int
}

func (i Int) Native() interface{} {
	return i.Int()
}

func (i Int) Hash(env *Env) (uint32, error) {
	h := getHash()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i.Int()))
	h.Write(b)
	return h.Sum32(), nil
}

func (i Int) Compare(env *Env, other Object) (int, error) {
	os, err := other.GetType().ToString(env, false)
	if err != nil {
		return 0, err
	}

	n, err := AssertNumber(env, other, "Cannot compare Int and "+os)
	if err != nil {
		return 0, err
	}
	return CompareNumbers(i, n), nil
}

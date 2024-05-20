package core

import (
	"encoding/binary"
	"math/big"
	"strconv"
)

// The Common integer type (can be Int or BigInt)
//
//lace:export
type Integer interface {
	integerType() string
	I64() int64
}

// The host int value
//
//lace:export
type Int int

var _ Integer = Int(0)

func (Int) integerType() string { return "int" }

func (i Int) I64() int64 {
	return int64(i)
}

func (i Int) ToString(env *Env, escape bool) (string, error) {
	return strconv.Itoa(i.I()), nil
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

func (i Int) WithInfo(info *ObjectInfo) any {
	var bi BigInt
	bi.info = info
	bi.b.SetInt64(int64(i))

	return &bi
}

func (i Int) Equals(env *Env, other interface{}) bool {
	return equalsNumbers(i, other)
}

/*
func (i Int) GetType() *Type {
	return TYPE.Int
}
*/

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

func (i Int) Compare(env *Env, other any) (int, error) {
	n, err := AssertNumber(env, other, "Cannot compare Int and "+TypeName(other))
	if err != nil {
		return 0, err
	}
	return CompareNumbers(i, n), nil
}

type BigInt struct {
	InfoHolder
	b big.Int
}

var _ Integer = &BigInt{}

func (*BigInt) integerType() string { return "bigint" }

func (bi *BigInt) I64() int64 {
	return bi.b.Int64()
}

func MakeBigInt(bi int64) *BigInt {
	return &BigInt{b: *big.NewInt(bi)}
}

func MakeBigIntFrom(bi *big.Int) *BigInt {
	return &BigInt{b: *bi}
}

func (bi *BigInt) ToString(env *Env, escape bool) (string, error) {
	return bi.b.String(), nil
}

func (bi *BigInt) Equals(env *Env, other interface{}) bool {
	return equalsNumbers(bi, other)
}

func (bi *BigInt) GetType() *Type {
	return TYPE.BigInt
}

func (bi *BigInt) Hash(env *Env) (uint32, error) {
	return hashGobEncoder(&bi.b)
}

func (bi *BigInt) Compare(env *Env, other any) (int, error) {
	n, err := AssertNumber(env, other, "Cannot compare BigInt and "+TypeName(other))
	if err != nil {
		return 0, err
	}
	return CompareNumbers(bi, n), nil
}

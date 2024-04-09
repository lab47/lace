package core

import (
	"bufio"
	"io"
	"unsafe"
)

type (
	BufferedReader struct {
		*bufio.Reader
		hash uint32
	}
)

var _ Object = &BufferedReader{}

func MakeBufferedReader(rd io.Reader) *BufferedReader {
	res := &BufferedReader{bufio.NewReader(rd), 0}
	res.hash = HashPtr(uintptr(unsafe.Pointer(res)))
	return res
}

func (br *BufferedReader) ToString(env *Env, escape bool) (string, error) {
	return "#object[BufferedReader]", nil
}

func (br *BufferedReader) Equals(env *Env, other interface{}) bool {
	return br == other
}

func (br *BufferedReader) GetInfo() *ObjectInfo {
	return nil
}

func (br *BufferedReader) GetType() *Type {
	return TYPE.BufferedReader
}

func (br *BufferedReader) Hash(env *Env) (uint32, error) {
	return br.hash, nil
}

func (br *BufferedReader) WithInfo(info *ObjectInfo) Object {
	return br
}

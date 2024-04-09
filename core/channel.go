package core

import (
	"unsafe"
)

type (
	FutureResult struct {
		value Object
		err   Error
	}
	Channel struct {
		ch       chan FutureResult
		isClosed bool
		hash     uint32
	}
)

var _ Object = &Channel{}

func MakeFutureResult(value Object, err Error) FutureResult {
	return FutureResult{value: value, err: err}
}

func (ch *Channel) ToString(env *Env, escape bool) (string, error) {
	return "#object[Channel]", nil
}

func (ch *Channel) Equals(env *Env, other interface{}) bool {
	return ch == other
}

func (ch *Channel) GetInfo() *ObjectInfo {
	return nil
}

func (ch *Channel) GetType() *Type {
	return TYPE.Channel
}

func (ch *Channel) Hash(env *Env) (uint32, error) {
	return ch.hash, nil
}

func (ch *Channel) WithInfo(info *ObjectInfo) Object {
	return ch
}

func MakeChannel(ch chan FutureResult) *Channel {
	res := &Channel{ch: ch, hash: 0}
	res.hash = HashPtr(uintptr(unsafe.Pointer(res)))
	return res
}

func ExtractChannel(env *Env, args []Object, index int) (*Channel, error) {
	return EnsureChannel(env, args, index)
}

func (ch *Channel) Close() {
	if !ch.isClosed {
		close(ch.ch)
		ch.isClosed = true
	}
}

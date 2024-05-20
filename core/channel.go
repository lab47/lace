package core

type (
	FutureResult struct {
		value any
		err   Error
	}
	Channel struct {
		ch       chan FutureResult
		isClosed bool
		hash     uint32
	}
)

var _ any = &Channel{}

func MakeFutureResult(value any, err Error) FutureResult {
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

func (ch *Channel) WithInfo(info *ObjectInfo) any {
	return ch
}

func MakeChannel(ch chan FutureResult) *Channel {
	res := &Channel{ch: ch, hash: 0}
	res.hash = HashPtr(res)
	return res
}

func ExtractChannel(env *Env, args []any, index int) (*Channel, error) {
	return EnsureChannel(env, args, index)
}

func (ch *Channel) Close() {
	if !ch.isClosed {
		close(ch.ch)
		ch.isClosed = true
	}
}

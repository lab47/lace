package rpc

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/lab47/lace/core"
	"github.com/lab47/lace/pkg/marshal"
	"github.com/mr-tron/base58"
)

type Request struct {
	Endpoint   string
	Method     string
	Arguments  core.Seq
	RequestId  string
	NoResponse bool
}

type marshaledRequest struct {
	Endpoint     string          `json:"endpoint" cbor:"1,keyasint"`
	Arguments    cbor.RawMessage `json:"arguments" cbor:"2,keyasint"`
	RequestId    string          `json:"request-id" cbor:"3,keyasint"`
	NeedResponse bool            `json:"need-response" cbor:"4,keyasint"`
	Method       string          `json:"method" cbor:"5,keyasint"`
}

func (r *Request) Marshal() ([]byte, error) {
	mr := marshaledRequest{
		Endpoint:     r.Endpoint,
		Method:       r.Method,
		RequestId:    r.RequestId,
		NeedResponse: r.NoResponse,
	}

	data, err := marshal.Marshal(r.Arguments)
	if err != nil {
		return nil, err
	}

	mr.Arguments = cbor.RawMessage(data)

	return cbor.Marshal(mr)
}

func (r *Request) Unmarshal(env *core.Env, data []byte) error {
	var mr marshaledRequest

	err := cbor.Unmarshal(data, &mr)
	if err != nil {
		return err
	}

	args, err := marshal.Unmarshal(env, mr.Arguments)
	if err != nil {
		return err
	}

	seq, ok := args.(core.Seq)
	if !ok {
		return fmt.Errorf("bad arguments, not seq")
	}

	r.Endpoint = mr.Endpoint
	r.Method = mr.Method
	r.RequestId = mr.RequestId
	r.NoResponse = mr.NeedResponse
	r.Arguments = seq

	return nil
}

type Response struct {
	Value core.Object
}

type marshalResponse struct {
	Value     cbor.RawMessage `json:"value" cbor:"1,keyasint"`
	RequestId string          `json:"request-id" cbor:"2,keyasint"`
}

func (r *Request) MarshalResponse(val core.Object) ([]byte, error) {
	d, err := marshal.Marshal(val)
	if err != nil {
		return nil, err
	}

	mr := marshalResponse{
		Value:     cbor.RawMessage(d),
		RequestId: r.RequestId,
	}

	return cbor.Marshal(mr)
}

func (r *Response) Unmarshal(env *core.Env, data []byte) error {
	var mr marshalResponse

	err := cbor.Unmarshal(data, &mr)
	if err != nil {
		return err
	}

	val, err := marshal.Unmarshal(env, mr.Value)
	if err != nil {
		return err
	}

	r.Value = val

	return nil
}

type Sender interface {
	Exchange(req *Request) (*Response, error)
}

type RPC interface {
	Request() *Request
	Respond(val core.Object) error
}

type Listener interface {
	Accept(context.Context) (RPC, error)
}

type Capabilities struct {
	Endpoint string              `json:"endpoint" cbor:"1,keyasint"`
	Method   string              `json:"method" cbor:"2,keyasint"`
	Tags     map[string][]string `json:"tags" cbor:"3,keyasint"`
}

var capEncoder cbor.EncMode

func init() {
	m, err := cbor.EncOptions{
		Sort: cbor.SortBytewiseLexical,
	}.EncMode()
	if err != nil {
		panic(err)
	}

	capEncoder = m
}

// Calculate a secure identifier for the capabilites based on it's
// content.
func (c *Capabilities) Id() (string, error) {
	data, err := capEncoder.Marshal(c)
	if err != nil {
		return "", err
	}

	h := sha256.New()
	h.Write(data)

	return base58.Encode(h.Sum(nil)), nil
}

type Advertisement interface {
	Clear()
}

type Advertiser interface {
	Advertise(c *Capabilities) (Advertisement, error)
}

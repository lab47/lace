package rpc

import (
	"context"
	"testing"
	"time"

	"github.com/lab47/lablog/logger"
	"github.com/lab47/lace/core"
	"github.com/stretchr/testify/require"
)

func TestRPC(t *testing.T) {
	t.Run("can exchange messages on a bus", func(t *testing.T) {
		r := require.New(t)

		log := logger.New(logger.Trace)
		env := &core.Env{}

		b, err := StartBus(log, env)
		r.NoError(err)

		defer b.Close()

		c1, err := b.Connect(context.TODO())
		r.NoError(err)

		l, err := c1.Listen("test.service")
		r.NoError(err)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var rpc RPC

		go func() {
			rpc, err = l.Accept(ctx)
			if err != nil {
				panic("bad accept")
			}

			rpc.Respond(core.MakeSymbol("ok"))
		}()

		c2, err := b.Connect(context.TODO())
		r.NoError(err)

		args := core.NewListFrom(core.MakeSymbol("arg1"))
		req := &Request{
			Endpoint:  "test.service",
			Method:    "foo",
			Arguments: args,
		}

		resp, err := c2.Exchange(req)
		r.NoError(err)

		r.Equal(core.MakeSymbol("ok"), resp.Value)

		rreq := rpc.Request()
		r.Equal("foo", rreq.Method)
		r.Equal("test.service", rreq.Endpoint)
		r.Equal(args, rreq.Arguments)
	})

	t.Run("returns an error if there are no listeners", func(t *testing.T) {
		r := require.New(t)

		log := logger.New(logger.Trace)
		env := &core.Env{}

		b, err := StartBus(log, env)
		r.NoError(err)

		defer b.Close()

		c2, err := b.Connect(context.TODO())
		r.NoError(err)

		args := core.NewListFrom(core.MakeSymbol("arg1"))
		req := &Request{
			Endpoint:  "test.service",
			Method:    "foo",
			Arguments: args,
		}

		_, err = c2.Exchange(req)
		r.Error(err)
	})

	t.Run("can advertise capabilities", func(t *testing.T) {
		r := require.New(t)

		log := logger.New(logger.Trace)
		env := &core.Env{}

		b, err := StartBus(log, env)
		r.NoError(err)

		defer b.Close()

		c1, err := b.Connect(context.TODO())
		r.NoError(err)

		c2, err := b.Connect(context.TODO())
		r.NoError(err)

		c2cap := &Capabilities{
			Endpoint: "test.service.c2",
			Method:   "has-fun",
			Tags: map[string][]string{
				"env": {"test"},
			},
		}

		c2.Advertise(c2cap)

		time.Sleep(100 * time.Millisecond)

		known := c1.BrowseCapabilities()

		r.Len(known, 1)

		r.Equal(c2cap, known[0])
	})
}

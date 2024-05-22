package rpc

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/lab47/lablog/logger"
	"github.com/lab47/lace/core"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

func InitBusConnection(env *core.Env) (*BusConnection, error) {
	logNs, err := env.InitNamespace(core.MakeSymbol("lace.log"))
	if err != nil {
		return nil, err
	}
	defLog := logNs.Resolve("*logger*")

	var log logger.Logger

	if defLog != nil {
		if err := core.Cast(env, defLog.GetStatic(), &log); err != nil {
			return nil, err
		}
	}

	if log == nil {
		log = logger.New(logger.Info)
	}

	ns := env.EnsureNamespace(core.MakeSymbol("lace.rpc"))

	vr := ns.Resolve("*connection*")
	if vr != nil {
		if conn, ok := vr.GetStatic().(*BusConnection); ok {
			return conn, nil
		}
	}

	var url string

	urlVr := ns.Resolve("url")
	if urlVr != nil {
		core.Cast(env, urlVr.GetStatic(), &url)
	}

	if url == "" {
		url = os.Getenv("LACE_RPC_URL")
	}

	if url == "" {
		bus, err := StartBus(log, env)
		if err != nil {
			return nil, err
		}

		ns.InternVar(env, "*bus*", bus, nil)

		conn, err := bus.Connect(env.Context)
		if err != nil {
			return nil, err
		}

		ns.InternVar(env, "*connection*", conn, nil)
		return conn, nil
	}

	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	bc := &BusConnection{
		env: env,
		log: log,
		c:   nc,

		caps: make(map[string]*busCap),
		seen: make(map[string]*busCap),
	}

	ns.InternVar(env, "*connection*", bc, nil)
	return bc, nil
}

type Bus struct {
	env *core.Env
	log logger.Logger
	s   *server.Server
}

func StartBus(log logger.Logger, env *core.Env) (*Bus, error) {
	var opts server.Options
	opts.DontListen = true

	s := server.New(&opts)

	go s.Start()

	if !s.ReadyForConnections(10 * time.Second) {
		return nil, fmt.Errorf("bus unable to start")
	}

	b := &Bus{
		env: env,
		log: log,
		s:   s,
	}

	return b, nil
}

func (b *Bus) Close() {
	b.s.Shutdown()
	b.s.WaitForShutdown()
}

func (b *Bus) ServeUnix(path string) error {
	l, err := net.Listen("unix", path)
	if err != nil {
		return err
	}

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}

		go b.serve(c)
	}
}

func (b *Bus) serve(c net.Conn) {
	defer c.Close()

	ip, err := b.s.InProcessConn()
	if err != nil {
		b.log.Error("error starting in process nat connection", "error", err)
		return
	}

	defer ip.Close()

	io.Copy(c, ip)
	go io.Copy(ip, c)
}

type busCap struct {
	Capabilities
	lastTime time.Time
}

type BusConnection struct {
	env *core.Env
	log logger.Logger
	c   *nats.Conn

	advertMu sync.Mutex
	caps     map[string]*busCap
	seen     map[string]*busCap
}

func (b *Bus) Connect(ctx context.Context) (*BusConnection, error) {
	c, err := nats.Connect("", nats.InProcessServer(b.s))
	if err != nil {
		return nil, err
	}

	bc := &BusConnection{
		c:    c,
		log:  b.log,
		env:  b.env,
		caps: make(map[string]*busCap),
		seen: make(map[string]*busCap),
	}

	err = bc.watchCaps(ctx)
	if err != nil {
		c.Close()
		return nil, err
	}

	go bc.broadcastCaps(ctx)

	return bc, nil
}

func (b *BusConnection) BrowseCapabilities() []*Capabilities {
	b.advertMu.Lock()
	defer b.advertMu.Unlock()

	var cps []*Capabilities

	for _, bc := range b.seen {
		cps = append(cps, &bc.Capabilities)
	}

	return cps
}

func (b *BusConnection) watchCaps(ctx context.Context) error {
	ch := make(chan *nats.Msg)
	sub, err := b.c.ChanSubscribe("lace.capabilities", ch)
	if err != nil {
		return err
	}

	go func() {
		defer sub.Unsubscribe()

		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-ch:
				b.processAdvert(msg.Data)
			}
		}
	}()

	return nil
}

func (b *BusConnection) processAdvert(data []byte) {
	var cp Capabilities

	err := cbor.Unmarshal(data, &cp)
	if err != nil {
		b.log.Error("error decoding advert", "error", err)
		return
	}

	id, err := cp.Id()
	if err != nil {
		b.log.Error("error calculating id for capabilities", "error", err)
		return
	}

	b.advertMu.Lock()
	defer b.advertMu.Unlock()

	b.seen[id] = &busCap{
		lastTime:     time.Now(),
		Capabilities: cp,
	}
}

var (
	broadcastThresh = time.Minute
	expireThresh    = 5 * time.Minute
)

func (b *BusConnection) broadcastCaps(ctx context.Context) {
	t := time.NewTimer(5 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			b.advertMu.Lock()

			var toBroadcast []*busCap

			for _, bc := range b.caps {
				if time.Since(bc.lastTime) > broadcastThresh {
					bc.lastTime = time.Now()
					toBroadcast = append(toBroadcast, bc)
				}
			}

			go b.broadcastAdverts(toBroadcast)

			var toDelete []string
			for id, bc := range b.seen {
				if time.Since(bc.lastTime) > expireThresh {
					toDelete = append(toDelete, id)
				}
			}

			for _, id := range toDelete {
				delete(b.seen, id)
			}

			b.advertMu.Unlock()
		}
	}
}

func (b *BusConnection) broadcastAdverts(toBroadcast []*busCap) {
	for _, c := range toBroadcast {
		data, err := cbor.Marshal(c.Capabilities)
		if err != nil {
			b.log.Error("error marshaling capabilities", "error", err)
			continue
		}
		err = b.c.Publish("lace.capabilities", data)
		if err != nil {
			b.log.Error("error publish capabilitie", "error", err)
		}
	}
}

func (b *BusConnection) Exchange(req *Request) (*Response, error) {
	if req.Arguments == nil {
		req.Arguments = core.NIL
	}

	data, err := req.Marshal()
	if err != nil {
		return nil, err
	}

	if req.NoResponse {
		return nil, b.c.Publish(req.Endpoint, data)
	}

	msg, err := b.c.Request(req.Endpoint, data, 60*time.Second)
	if err != nil {
		return nil, err
	}

	var r Response

	err = r.Unmarshal(b.env, msg.Data)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

type BusListener struct {
	b   *BusConnection
	sub *nats.Subscription
	ch  chan *nats.Msg
}

func (b *BusConnection) Listen(endpoint string) (Listener, error) {
	ch := make(chan *nats.Msg, 1)
	sub, err := b.c.ChanQueueSubscribe(endpoint, endpoint, ch)
	if err != nil {
		return nil, err
	}

	l := &BusListener{
		b:   b,
		sub: sub,
		ch:  ch,
	}

	return l, nil
}

type BusRPC struct {
	msg *nats.Msg
	req *Request
}

func (b *BusListener) Accept(ctx context.Context) (RPC, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case msg := <-b.ch:
		var req Request

		err := req.Unmarshal(b.b.env, msg.Data)
		if err != nil {
			return nil, err
		}

		req.RequestId = msg.Reply

		return &BusRPC{
			msg: msg,
			req: &req,
		}, nil
	}
}

func (r *BusRPC) Request() *Request {
	return r.req
}

func (r *BusRPC) Respond(val any) error {
	if r.req.NoResponse {
		return nil
	}

	data, err := r.req.MarshalResponse(val)
	if err != nil {
		return err
	}

	return r.msg.Respond(data)
}

func (c *BusConnection) Advertise(cp *Capabilities) (Advertisement, error) {
	id, err := cp.Id()
	if err != nil {
		return nil, err
	}

	c.advertMu.Lock()
	defer c.advertMu.Unlock()

	bc := &busCap{
		Capabilities: *cp,
	}

	c.caps[id] = bc

	go c.broadcastAdverts([]*busCap{bc})

	return &BusAdvert{
		c:  c,
		id: id,
	}, nil
}

type BusAdvert struct {
	c  *BusConnection
	id string
}

func (a *BusAdvert) Clear() {
	a.c.advertMu.Lock()
	defer a.c.advertMu.Unlock()

	delete(a.c.caps, a.id)
}

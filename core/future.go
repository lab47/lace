package core

import (
	"fmt"
	"sync"
)

type IPending interface {
	IsRealized() bool
}

type Future struct {
	call     Callable
	ch       chan struct{}
	mu       sync.Mutex
	realized bool
	err      error
	value    any
}

var _ IPending = (*Future)(nil)

// NewFuture creates a new Future value and schedules the future
// to be run. Deref'ing the Future will retrieve the value (potentially
// waiting if the value is not yet ready)
//
//lace:export
func NewFuture(env *Env, call Callable) (*Future, error) {
	f := &Future{
		call: call,
		ch:   make(chan struct{}, 1),
	}

	err := f.Schedule(env)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (f *Future) Schedule(env *Env) error {
	go func() {
		env := env.Child()

		obj, err := f.call.Call(env, nil)
		f.mu.Lock()
		defer f.mu.Unlock()

		f.realized = true
		f.value = obj
		f.err = err

		close(f.ch)
	}()

	return nil
}

func (f *Future) IsRealized() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.realized
}

func (f *Future) Deref(env *Env) (any, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	for !f.realized {
		select {
		case <-env.Context.Done():
			return nil, env.Context.Err()
		case <-f.ch:
			if !f.realized {
				return nil, fmt.Errorf("dead future without value")
			}
		}
	}

	return f.value, nil
}

var _ Deref = (*Future)(nil)

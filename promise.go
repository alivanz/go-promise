package promise

import (
	"sync/atomic"
	"unsafe"
)

type Promise[T any] struct {
	success Pending[T]
	fail    Pending[error]
	final   PendingNoValue
	state   int32
}

type State int32

const (
	StatePending State = iota
	StateResolved
	StateError
)

func NewPromise[T any]() *Promise[T] {
	return &Promise[T]{}
}

func (prom *Promise[T]) Then(f func(v T)) *Promise[T] {
	prom.success.Then(f)
	return prom
}

func (prom *Promise[T]) Catch(f func(err error)) *Promise[T] {
	prom.fail.Then(f)
	return prom
}

func (prom *Promise[T]) Finally(f func()) *Promise[T] {
	if prom.final.then(unsafe.Pointer(&f)) {
		f()
	}
	return prom
}

func (prom *Promise[T]) Resolve(v T) *Promise[T] {
	// make sure this runs once
	if atomic.CompareAndSwapInt32(&prom.state, int32(StatePending), int32(StateResolved)) {
		defer prom.final.Resolve()
		prom.success.Resolve(v)
	}
	return prom
}

func (prom *Promise[T]) Error(err error) *Promise[T] {
	// make sure this runs once
	if atomic.CompareAndSwapInt32(&prom.state, int32(StatePending), int32(StateError)) {
		defer prom.final.Resolve()
		prom.fail.Resolve(err)
	}
	return prom
}

func (prom *Promise[T]) State() State {
	return State(atomic.LoadInt32(&prom.state))
}

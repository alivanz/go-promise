package promise

import (
	"sync/atomic"
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
	prom.final.Then(f)
	return prom
}

func (prom *Promise[T]) Resolve(v T) {
	// make sure this runs once
	if atomic.CompareAndSwapInt32(&prom.state, int32(StatePending), int32(StateResolved)) {
		defer prom.final.Resolve()
		prom.success.Resolve(v)
	}
}

func (prom *Promise[T]) Error(err error) {
	// make sure this runs once
	if atomic.CompareAndSwapInt32(&prom.state, int32(StatePending), int32(StateError)) {
		defer prom.final.Resolve()
		prom.fail.Resolve(err)
	}
}

func (prom *Promise[T]) State() State {
	return State(atomic.LoadInt32(&prom.state))
}

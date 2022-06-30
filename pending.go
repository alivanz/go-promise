package promise

import (
	"sync/atomic"
)

type PendingNoValue struct {
	f     *func()
	v     int32
	state int32
}

type Pending[T any] struct {
	f     *func(T)
	v     *T
	state int32
}

const (
	funcSet  = 1
	valueSet = 2
	bothSet  = 3
)

func (pending *PendingNoValue) Then(f func()) {
	if !cas(&pending.f, nil, &f) {
		return
	}
	if atomic.AddInt32(&pending.state, funcSet) != bothSet {
		return
	}
	f()
}

func (pending *PendingNoValue) Resolve() {
	if !atomic.CompareAndSwapInt32(&pending.v, 0, 1) {
		return
	}
	if atomic.AddInt32(&pending.state, valueSet) != bothSet {
		return
	}
	pf := load(&pending.f)
	(*pf)()
}

func (pending *Pending[T]) Then(f func(T)) {
	if !cas(&pending.f, nil, &f) {
		return
	}
	if atomic.AddInt32(&pending.state, funcSet) != bothSet {
		return
	}
	pv := load(&pending.v)
	f(*pv)
}

func (pending *Pending[T]) Resolve(v T) {
	if !cas(&pending.v, nil, &v) {
		return
	}
	if atomic.AddInt32(&pending.state, valueSet) != bothSet {
		return
	}
	pf := load(&pending.f)
	(*pf)(v)
}

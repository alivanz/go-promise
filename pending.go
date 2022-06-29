package promise

import (
	"sync/atomic"
	"unsafe"
)

type PendingNoValue struct {
	f     unsafe.Pointer
	state int32
}

type Pending[T any] struct {
	PendingNoValue
	v unsafe.Pointer
}

const (
	funcSet  = 1
	valueSet = 2
	bothSet  = 3
)

// then set callback, returns true if resolved
func (pending *PendingNoValue) then(f unsafe.Pointer) bool {
	// make sure f only set once
	if !atomic.CompareAndSwapPointer(&pending.f, nil, f) {
		return false
	}
	// mark f as done
	state := atomic.AddInt32(&pending.state, funcSet)
	return state == bothSet
}

// resolve returns the callback function if set
func (pending *PendingNoValue) resolve() unsafe.Pointer {
	// mark v as done
	state := atomic.AddInt32(&pending.state, valueSet)
	if state != bothSet {
		return nil
	}
	return atomic.LoadPointer(&pending.f)
}

func (pending *PendingNoValue) Then(f func()) {
	if !pending.then(unsafe.Pointer(&f)) {
		return
	}
	f()
}

func (pending *PendingNoValue) Resolve() {
	p := pending.resolve()
	if p == nil {
		return
	}
	pf := (*func())(p)
	(*pf)()
}

func (pending *Pending[T]) Then(f func(T)) {
	if !pending.then(unsafe.Pointer(&f)) {
		return
	}
	// get v
	p := atomic.LoadPointer(&pending.v)
	pv := (*T)(p)
	f(*pv)
}

func (pending *Pending[T]) Resolve(v T) {
	// make sure v only set once
	if !atomic.CompareAndSwapPointer(&pending.v, nil, unsafe.Pointer(&v)) {
		return
	}
	// get f
	p := pending.resolve()
	if p == nil {
		return
	}
	pf := (*func(T))(p)
	(*pf)(v)
}

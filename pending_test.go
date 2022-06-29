package promise

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestPendingThenFirst(t *testing.T) {
	var u Pending[int64]
	var count int32
	u.Then(func(arg1 int64) {
		count += 1
	})
	u.Resolve(10)
	if count != 1 {
		t.Fail()
	}
}

func TestPendingResolveFirst(t *testing.T) {
	var u Pending[int64]
	var count int32
	u.Resolve(10)
	u.Then(func(arg1 int64) {
		count += 1
	})
	if count != 1 {
		t.Fail()
	}
}

func TestPendingRace(t *testing.T) {
	var u Pending[int64]
	var wg sync.WaitGroup
	var count int32
	wg.Add(2)
	go func() {
		defer wg.Done()
		u.Then(func(arg1 int64) {
			atomic.AddInt32(&count, 1)
			if arg1 != 10 {
				t.Fail()
			}
		})
		u.Then(func(arg1 int64) {
			t.Fail()
		})
	}()
	go func() {
		defer wg.Done()
		u.Resolve(10)
		u.Resolve(11)
		u.Resolve(12)
	}()
	wg.Wait()
	if count != 1 {
		t.Fail()
	}
}

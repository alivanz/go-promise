package promise

import "testing"

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

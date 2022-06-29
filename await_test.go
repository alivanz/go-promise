package promise

import (
	"fmt"
	"testing"
)

func TestAwaitSuccess(t *testing.T) {
	prom := NewPromise[int]()
	go prom.Resolve(123)
	if v, err := prom.Await(); err != nil {
		t.Fatal(err)
	} else if v != 123 {
		t.Fatal(v)
	}
}

func TestAwaitError(t *testing.T) {
	prom := NewPromise[int]()
	go prom.Error(fmt.Errorf("test"))
	if v, err := prom.Await(); err == nil {
		t.Fatalf("unexpected success: %v", v)
	} else if err.Error() != "test" {
		t.Fatal(err)
	}
}

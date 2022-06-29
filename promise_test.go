package promise

import (
	"fmt"
	"testing"
	"time"
)

func TestResolve(t *testing.T) {
	prom := NewPromise[int]()
	prom.Resolve(123)
	i, err := prom.Await()
	if err != nil {
		t.Fatal(err)
	}
	if i != 123 {
		t.Fail()
	}
}

func TestError(t *testing.T) {
	prom := NewPromise[int]()
	prom.Error(fmt.Errorf("test"))
	_, err := prom.Await()
	if err == nil {
		t.Fail()
	} else if err.Error() != "test" {
		t.Fail()
	}
}

func TestResolveDouble(t *testing.T) {
	prom := NewPromise[int]()
	prom.Resolve(123)
	prom.Error(fmt.Errorf("test"))
	prom.Resolve(456)
	i, err := prom.Await()
	if err != nil {
		t.Fatal(err)
	}
	if i != 123 {
		t.Fail()
	}
}

func TestErrorDouble(t *testing.T) {
	prom := NewPromise[int]()
	prom.Error(fmt.Errorf("test"))
	prom.Resolve(123)
	prom.Error(fmt.Errorf("test2"))
	_, err := prom.Await()
	if err == nil {
		t.Fail()
	} else if err.Error() != "test" {
		t.Fail()
	}
}

func TestResolveThread(t *testing.T) {
	t.Parallel()
	prom := NewPromise[int]()
	go func() {
		time.Sleep(100 * time.Millisecond)
		prom.Resolve(123)
	}()
	i, err := prom.Await()
	if err != nil {
		t.Fatal(err)
	}
	if i != 123 {
		t.Fail()
	}
}

func TestErrorThread(t *testing.T) {
	t.Parallel()
	prom := NewPromise[int]()
	go func() {
		time.Sleep(100 * time.Millisecond)
		prom.Error(fmt.Errorf("test"))
	}()
	_, err := prom.Await()
	if err == nil {
		t.Fail()
	} else if err.Error() != "test" {
		t.Fail()
	}
}

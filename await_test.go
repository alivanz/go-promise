package promise

import (
	"fmt"
	"testing"
)

func TestResolve(t *testing.T) {
	prom, await := NewAwait[int]()
	prom.Resolve(123)
	i, err := await.Wait()
	if err != nil {
		t.Fatal(err)
	}
	if i != 123 {
		t.Fail()
	}
}

func TestError(t *testing.T) {
	prom, await := NewAwait[int]()
	prom.Error(fmt.Errorf("test"))
	_, err := await.Wait()
	if err == nil {
		t.Fail()
	} else if err.Error() != "test" {
		t.Fail()
	}
}

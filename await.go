package promise

import (
	"sync"
)

func (prom *Promise[T]) Await() (T, error) {
	var wg sync.WaitGroup
	wg.Add(1)
	var ret T
	var err error
	prom.Then(func(v T) {
		ret = v
	}).Catch(func(e error) {
		err = e
	}).Finally(func() {
		wg.Done()
	})
	wg.Wait()
	return ret, err
}

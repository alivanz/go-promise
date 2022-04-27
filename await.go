package promise

import (
	"sync"
)

type IAwait[T any] interface {
	Wait() (T, error)
}

type Await[T any] struct {
	v   T
	err error
	wg  sync.WaitGroup
}

func NewAwait[T any]() (Promise[T], *Await[T]) {
	prom := NewSimple[T]()
	await := NewAwaitFromPromise(prom)
	return prom, await
}

func NewAwaitFromPromise[T any](prom Promise[T]) *Await[T] {
	await := &Await[T]{}
	await.wg.Add(1)
	prom.Then(await.onResolve)
	prom.Catch(await.onError)
	prom.Finally(await.onFinal)
	return await
}

func (await *Await[T]) onResolve(v T) {
	await.v = v
}

func (await *Await[T]) onError(err error) {
	await.err = err
}

func (await *Await[T]) onFinal() {
	await.wg.Done()
}

func (await *Await[T]) Wait() (T, error) {
	await.wg.Wait()
	return await.v, await.err
}

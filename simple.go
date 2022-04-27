package promise

type Simple[T any] struct {
	then    func(T)
	catch   func(error)
	finally func()
}

func (prom *Simple[T]) Then(f func(v T)) Promise[T] {
	prom.then = f
	return prom
}

func (prom *Simple[T]) Catch(f func(err error)) Promise[T] {
	prom.catch = f
	return prom
}

func (prom *Simple[T]) Finally(f func()) Promise[T] {
	prom.finally = f
	return prom
}

func (prom *Simple[T]) Resolve(v T) {
	if prom.then != nil {
		prom.then(v)
	}
	if prom.finally != nil {
		prom.finally()
	}
}

func (prom *Simple[T]) Error(err error) {
	if prom.catch != nil {
		prom.catch(err)
	}
	if prom.finally != nil {
		prom.finally()
	}
}

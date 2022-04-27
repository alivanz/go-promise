package promise

type Promise[T any] interface {
	Then(func(T)) Promise[T]
	Catch(func(error)) Promise[T]
	Finally(func()) Promise[T]
	Resolve(v T)
	Error(err error)
}

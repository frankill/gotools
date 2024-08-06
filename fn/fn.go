package fn

type Fun[U, T any] func(x ...U) T

type Fn[F Fun[U, T], U, T any] struct {
	param []U
	fun   F
}

func New[F Fun[U, T], U, T any](f F) *Fn[F, U, T] {
	return &Fn[F, U, T]{fun: f}
}

func (f *Fn[F, U, T]) Partial(x ...U) *Fn[F, U, T] {
	f.param = x
	return f
}

func (f *Fn[F, U, T]) Call(x ...U) T {
	return f.fun(append(f.param, x...)...)
}

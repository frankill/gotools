package gotools

func Cmp[T Ordered](x, y T) int {
	if x > y {
		return 1
	} else if x < y {
		return -1
	}
	return 0
}

func Gt[T Ordered](x, y T) bool {
	return x > y
}

func Eq[T Comparable](x, y T) bool {
	return x == y
}

func Lte[T Ordered](x, y T) bool {
	return x <= y
}

func Lt[T Ordered](x, y T) bool {
	return x < y
}

func Gte[T Ordered](x, y T) bool {
	return x >= y
}

func NotEq[T Comparable](x, y T) bool {
	return x != y
}

func Identity[T any](x T) T {
	return x
}

var (
	ASCInt  = Lt[int]
	DESCInt = Gt[int]
)

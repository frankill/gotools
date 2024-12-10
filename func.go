package gotools

func Cmp[T Ordered](x, y T) int {
	if x > y {
		return 1
	} else if x < y {
		return -1
	}
	return 0
}

func Add[T Number](x, y T) T {
	return x + y
}

func Sub[T Number](x, y T) T {
	return x - y
}

func Mul[T Number](x, y T) T {
	return x * y
}

func Div[T Number](x, y T) T {
	return x / y
}

func Mod[T Integer](x, y T) T {
	return x % y
}

func Min[T Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func Max[T Ordered](x, y T) T {
	if x > y {
		return x
	}
	return y
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

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

// Ifelse 根据给定的布尔条件 condition，选择返回 trueVal 或 falseVal。
// 如果 condition 为 true，则返回 trueVal；否则返回 falseVal。
// 参数:
//
//   - condition - 用于判断的布尔条件。
//   - trueVal - 当 condition 为 true 时返回的值。
//   - falseVal - 当 condition 为 false 时返回的值。
//
// 返回:
//
//   - 根据 condition 的结果返回 trueVal 或 falseVal。
func Ifelse[T any](condition bool, trueVal T, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

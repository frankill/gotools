package operation

import (
	"github.com/frankill/gotools"
	"github.com/frankill/gotools/array"
	"github.com/frankill/gotools/fn"
)

// All 逐元素执行函数。
// 参数:
//
//	fun: 用于逐元素执行的函数。
//	arr: 变长参数列表，每个参数都是一个切片。
//
// 返回:
//
//	如果所有元素都满足条件，则返回 true，否则返回 false。
func All[T any](f func(T) bool, arr []T) bool {

	for _, v := range arr {
		if !f(v) {
			return false
		}
	}
	return true
}

// Any 逐元素执行函数。
// 参数:
//
//	fun: 用于逐元素执行的函数。
//	arr: 变长参数列表，每个参数都是一个切片。
//
// 返回:
//
//	如果至少有一个元素满足条件，则返回 true，否则返回 false。
func Any[T any](f func(T) bool, arr []T) bool {
	for _, v := range arr {
		if f(v) {
			return true
		}
	}
	return false
}

// ForEach 逐元素执行函数。
// 参数:
//
//	fun: 用于逐元素执行的函数。
//	arr: 变长参数列表，每个参数都是一个切片。
//
// 返回:
//
//	无返回值。
func ForEach[S ~[]T, T any](fun func(x ...T), arr ...S) {

	if len(arr) == 0 {
		return
	}

	la := array.Map(func(x ...S) int { return len(x[0]) }, arr)
	lm := array.Max(la)

	for i := 0; i < lm; i++ {
		parm := make([]T, len(arr))
		for j := 0; j < len(arr); j++ {
			parm[j] = arr[j][i%la[j]]
		}
		fun(parm...)
	}

}

// Add 逐元素相加
// 参数:
//
//	a: 一个切片
//	b: 一个切片
//
// 返回:
//
//	一个切片
func Add[S ~[]T, T gotools.Number](a, b S) []T {
	return Operator(func(x ...T) T { return x[0] + x[1] }, a, b)
}

// Sub 逐元素相减
// 参数:
//
//	a: 一个切片
//	b: 一个切片
//
// 返回:
//
//	一个切片
func Sub[S ~[]T, T gotools.Number](a, b S) []T {
	return Operator(func(x ...T) T { return x[0] - x[1] }, a, b)
}

// Mul 逐元素相乘
// 参数:
//
//	a: 一个切片
//	b: 一个切片
//
// 返回:
//
//	一个切片
func Mul[S ~[]T, T gotools.Number](a, b S) []T {
	return Operator(func(x ...T) T { return x[0] * x[1] }, a, b)
}

// Div 逐元素相除
// 参数:
//
//	a: 一个切片
//	b: 一个切片
//
// 返回:
//
//	一个切片
func Div[S ~[]T, T gotools.Number, R float64](a, b S) []R {
	return Operator(func(x ...T) R { return R(x[0]) / R(x[1]) }, a, b)
}

// Mod 逐元素取模
// 参数:
//
//	a: 一个切片
//	b: 一个切片
//
// 返回:
//
//	一个切片
func Mod[S ~[]T, T gotools.Integer](a, b S) []int {
	return fn.Lapply2(func(x, y T) int { return int(x % y) }, a, b)
}

// Lte 逐元素比较大小
// 参数:
//
//	a: 一个切片
//	b: 一个切片
//
// 返回:
//
//	一个切片
func Lte[S ~[]T, T gotools.Ordered](a, b S) []bool {
	return Operator(func(x ...T) bool { return x[0] <= x[1] }, a, b)
}

// Gte 逐元素比较大小
// 参数:
//
//	a: 一个切片
//	b: 一个切片
//
// 返回:
//
//	一个切片
func Gte[S ~[]T, T gotools.Ordered](a, b S) []bool {
	return Operator(func(x ...T) bool { return x[0] >= x[1] }, a, b)
}

// Lt 逐元素比较
// 参数:
//
//	a: 一个切片
//	b: 一个切片
//
// 返回:
//
//	一个切片
func Lt[S ~[]T, T gotools.Ordered](a, b S) []bool {
	return Operator(func(x ...T) bool { return x[0] < x[1] }, a, b)
}

// Gt 逐元素比较大小
// 参数:
//
//	a: 一个切片
//	b: 一个切片
//
// 返回:
//
//	一个切片
func Gt[S ~[]T, T gotools.Ordered](a, b S) []bool {
	return Operator(func(x ...T) bool { return x[0] > x[1] }, a, b)
}

// Eq 逐元素比较
// 参数:
//
//	a: 一个切片
//	b: 一个切片
//
// 返回:
//
//	一个切片
func Eq[S ~[]T, T gotools.Comparable](a, b S) []bool {
	return Operator(func(x ...T) bool { return x[0] == x[1] }, a, b)
}

// Neq 逐元素比较是否不相等
// 参数:
//
//	a: 一个切片
//	b: 一个切片
//
// 返回:
//
//	一个切片
func Neq[S ~[]T, T gotools.Comparable](a, b S) []bool {
	return Operator(func(x ...T) bool { return x[0] != x[1] }, a, b)
}

// Not 逐元素取反
// 参数:
//
//	a: 一个切片
//
// 返回:
//
//	一个切片
func Not(a []bool) []bool {
	return fn.Lapply(func(x bool) bool { return !x }, a)
}

// Operator 对一组切片应用指定的函数，每个切片元素按位置组合后作为函数的参数。
// 此泛型函数接受一个变长参数列表 `arr`，其中每个参数都是类型 S 的切片（S 必须是切片类型），
// 和一个函数 `fun`，该函数接受变长参数列表 x...T 并返回类型 U 的结果。
// 函数返回一个类型为 []U 的切片，其中包含应用给定函数 `fun` 到所有切片元素组合上的结果。
//
// 参数:
//
//   - fun: 一个函数，接受变长参数列表 x...T，并返回类型 U 的结果。
//   - arr: 变长参数列表，每个参数都是类型 S 的切片，S 必须是切片类型。
//
// 返回:
//
//   - 类型为 []U 的切片，包含应用给定函数 `fun` 到所有切片元素组合上的结果。
func Operator[S ~[]T, T, U any](fun func(x ...T) U, arr ...S) []U {

	if len(arr) == 0 {
		return make([]U, 0)
	}

	la := array.Map(func(x ...S) int { return len(x[0]) }, arr)
	lm := array.Max(la)

	res := make([]U, lm)

	for i := 0; i < lm; i++ {
		parm := make([]T, len(arr))
		for j := 0; j < len(arr); j++ {
			parm[j] = arr[j][i%la[j]]
		}
		res[i] = fun(parm...)
	}

	return res
}

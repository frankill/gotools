package operation

import (
	"cmp"

	"github.com/frankill/gotools/array"
)

// And 对多个布尔值进行逻辑与运算。
// 参数:
//
//	arr: 变长参数列表，每个参数都是布尔类型。
//
// 返回:
//
//	如果所有参数都为 true，则返回 true；否则返回 false。
func And[T bool](arr ...T) T {

	if len(arr) == 0 {
		return false
	}

	for _, v := range arr {

		if !v {
			return false
		}

	}
	return true

}

// Or 对多个布尔值进行逻辑或运算。
// 参数:
//
//	arr: 变长参数列表，每个参数都是布尔类型。
//
// 返回:
//
//	如果至少有一个参数为 true，则返回 true；否则返回 false。
func Or[T bool](arr ...T) T {
	if len(arr) == 0 {
		return false
	}
	for _, v := range arr {
		if v {
			return true
		}
	}
	return false
}

// Not 对单个布尔值进行逻辑非运算。
// 参数:
//
//	v: 布尔类型值。
//
// 返回:
//
//	如果 v 为 true，则返回 false；否则返回 true。
func Not[T bool](v T) T {
	return !v
}

// GT 检查有序类型参数列表是否满足严格递减顺序。
// 参数:
//
//	arr: 变长参数列表，每个参数都实现了 cmp.Ordered 接口。
//
// 返回:
//
//	如果参数列表满足严格递减顺序，则返回 true；否则返回 false。
func GT[T cmp.Ordered](arr ...T) bool {

	la := len(arr)
	if la == 0 {
		return true
	}

	for i := 1; i < la; i++ {
		if arr[i-1] < arr[i] {
			return false
		}
	}

	return true

}

// GTE 检查有序类型参数列表是否满足非严格递减顺序。
// 参数:
//
//	arr: 变长参数列表，每个参数都实现了 cmp.Ordered 接口。
//
// 返回:
//
//	如果参数列表满足非严格递减顺序，则返回 true；否则返回 false。
func GTE[T cmp.Ordered](arr ...T) bool {

	la := len(arr)
	if la == 0 {
		return true
	}

	for i := 1; i < la; i++ {
		if arr[i-1] <= arr[i] {
			return false
		}
	}

	return true

}

// LT 检查有序类型参数列表是否满足严格递增顺序。
// 参数:
//
//	arr: 变长参数列表，每个参数都实现了 cmp.Ordered 接口。
//
// 返回:
//
//	如果参数列表满足严格递增顺序，则返回 true；否则返回 false。
func LT[T cmp.Ordered](arr ...T) bool {

	la := len(arr)
	if la == 0 {
		return true
	}

	for i := 1; i < la; i++ {
		if arr[i-1] > arr[i] {
			return false
		}
	}

	return true

}

// LTE 检查有序类型参数列表是否满足非严格递增顺序。
// 参数:
//
//	arr: 变长参数列表，每个参数都实现了 cmp.Ordered 接口。
//
// 返回:
//
//	如果参数列表满足非严格递增顺序，则返回 true；否则返回 false。
func LTE[T cmp.Ordered](arr ...T) bool {

	la := len(arr)
	if la == 0 {
		return true
	}

	for i := 1; i < la; i++ {
		if arr[i-1] >= arr[i] {
			return false
		}
	}

	return true

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

// OperatorMap 对一组切片应用指定的函数，每个切片元素按位置组合后作为函数的参数。
// 此泛型函数接受一个变长参数列表 `arr`，其中每个参数都是类型 S 的切片（S 必须是切片类型），
// 和一个函数 `fun`，该函数接受变长参数列表 x...T 并返回类型 U 的结果。
// 函数返回一个类型为 []U 的切片，其中包含应用给定函数 `fun` 到所有切片元素组合上的结果。
//
// 参数:
//
//	fun: 一个函数，接受变长参数列表 x...T，并返回类型 U 的结果。
//	arr: 变长参数列表，每个参数都是类型 S 的切片，S 必须是切片类型。
//
// 返回:
//
//	类型为 []U 的切片，包含应用给定函数 `fun` 到所有切片元素组合上的结果。
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

// Cartesian函数生成多个切片的笛卡尔积。
//
// 它接受类型为`[]T`的可变参数`arr`，表示要组合的切片。
// 类型`T`表示切片中元素的类型。
//
// 函数返回类型为`[][]T`的切片，表示输入切片的笛卡尔积。
// 例如:Cartesian([][]int{{1, 2}, {3, 4}}) = [][]int{{1, 3}, {1, 4}, {2, 3}, {2, 4}}
func Cartesian[S []T, T any](arr ...S) [][]T {
	if len(arr) == 0 {
		return [][]T{}
	}

	colNum := len(arr)
	rowNum := int(array.Product(array.Map(func(x ...S) int { return len(x[0]) }, arr)))

	res := make([][]T, colNum)
	res[0] = arr[0]

	for i := 1; i < colNum; i++ {
		res[i] = array.Rep(arr[i], len(res[i-1]), true)
	}

	for i := 0; i < colNum; i++ {
		if n := rowNum / len(res[i]); n > 1 {
			res[i] = array.Rep(res[i], n, false)
		}

	}

	return array.Zip(res...)
}

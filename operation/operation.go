package operation

import (
	"github.com/frankill/gotools/array"
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

// Operator 对一组切片应用指定的函数，每个切片元素按位置组合后作为函数的参数。
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

	rowNum := 1
	for _, a := range arr {
		rowNum *= len(a)
	}

	res := make([][]T, rowNum)
	indices := make([]int, len(arr))

	for i := range res {
		row := make([]T, len(arr))
		for j, a := range arr {
			row[j] = a[indices[j]]
		}
		res[i] = row
		// 增加索引
		for j := len(arr) - 1; j >= 0; j-- {
			indices[j]++
			if indices[j] < len(arr[j]) {
				break
			}
			indices[j] = 0
		}
	}

	return res
}

// func Cartesian[S []T, T any](arr ...S) [][]T {
// 	if len(arr) == 0 {
// 		return [][]T{}
// 	}

// 	colNum := len(arr)
// 	rowNum := int(array.Product(array.Map(func(x ...S) int { return len(x[0]) }, arr)))

// 	res := make([][]T, colNum)

// 	for i := 0; i < colNum; i++ {
// 		res[i] = make([]T, rowNum)
// 	}

// 	res[0] = arr[0]

// 	for i := 1; i < colNum; i++ {
// 		copy(res[i], array.Rep(arr[i], len(res[i-1]), true))
// 	}

// 	for i := 0; i < colNum; i++ {
// 		if n := rowNum / len(res[i]); n > 1 {
// 			copy(res[i], array.Rep(res[i], n, false))
// 		}

// 	}

// 	return array.Zip(res...)
// }

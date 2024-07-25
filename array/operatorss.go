package array

import (
	"cmp"
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

// ArrayForEach 逐元素执行函数。
// 参数:
//
//	fun: 用于逐元素执行的函数。
//	arr: 变长参数列表，每个参数都是一个切片。
//
// 返回:
//
//	无返回值。
func ForEach[S ~[]T, T any](fun func(x ...T), arr ...S) {
	if len(arr) == 0 || len(arr[0]) == 0 {
		return
	}
	l := len(arr[0])
	f := len(arr)
	param := make([]T, f)

	for i := 0; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i]
		}
		fun(param...)
	}

}

// ArrayAnd 对两个布尔切片进行逐元素逻辑与操作。
// 参数:
//
//	x: 第一个布尔切片。
//	y: 第二个布尔切片。
//
// 返回:
//
//	类型为 []bool 的切片，表示 x 和 y 中对应元素的逻辑与结果。
func ArrayAnd[S ~[]T, T bool](x, y S) []T {
	return Operator(And[T], x, y)
}

// ArrayOr 对两个布尔切片进行逐元素逻辑或操作。
// 参数:
//
//	x: 第一个布尔切片。
//	y: 第二个布尔切片。
//
// 返回:
//
//	类型为 []bool 的切片，表示 x 和 y 中对应元素的逻辑或结果。
func ArrayOr[S ~[]T, T bool](x, y S) []T {
	return Operator(Or[T], x, y)
}

// ArrayAdd 对两个数值或字符串切片进行逐元素相加操作。
// 参数:
//
//	x: 第一个切片，可以是数字或字符串类型。
//	y: 第二个切片，可以是数字或字符串类型。
//
// 返回:
//
//	类型为 []T 的切片，表示 x 和 y 中对应元素的相加结果。
func ArrayAdd[S ~[]T, T Number | string](x, y S) []T {

	return Operator(func(x ...T) T { return x[0] + x[1] }, x, y)
}

// ArraySub 对两个数值切片进行逐元素相减操作。
// 参数:
//
//	x: 第一个数值切片。
//	y: 第二个数值切片。
//
// 返回:
//
//	类型为 []T 的切片，表示 x 和 y 中对应元素的相减结果。
func ArraySub[S ~[]T, T Number](x, y S) []T {

	return Operator(func(x ...T) T { return x[0] - x[1] }, x, y)
}

// ArrayMultiply 对两个数值切片进行逐元素相乘操作。
// 参数:
//
//	x: 第一个数值切片。
//	y: 第二个数值切片。
//
// 返回:
//
//	类型为 []T 的切片，表示 x 和 y 中对应元素的相乘结果。
func ArrayMultiply[S ~[]T, T Number](x, y S) []T {

	return Operator(func(x ...T) T { return x[0] * x[1] }, x, y)
}

// ArrayDivide 对两个数值切片进行逐元素相除操作。
// 参数:
//
//	x: 第一个数值切片。
//	y: 第二个数值切片。
//
// 返回:
//
//	类型为 []T 的切片，表示 x 和 y 中对应元素的相除结果。
func ArrayDivide[S ~[]T, T Number](x, y S) []T {

	return Operator(func(x ...T) T { return x[0] / x[1] }, x, y)
}

// ArrayGt 比较两个有序类型的切片中对应位置的元素，判断第一个切片中的元素是否大于第二个切片中的元素。
// 参数:
//
//	x: 第一个切片，其元素将与 y 中的元素进行比较。
//	y: 第二个切片，作为比较的目标。
//
// 返回:
//
//	类型为 []bool 的切片，表示 x 中相应位置的元素是否大于 y 的对应元素。
func ArrayGt[S ~[]T, T cmp.Ordered](x, y S) []bool {

	return Operator(func(x ...T) bool { return x[0] > x[1] }, x, y)
}

// ArrayLt 比较两个有序类型的切片中对应位置的元素，判断第一个切片中的元素是否小于第二个切片中的元素。
// 参数:
//
//	x: 第一个切片，其元素将与 y 中的元素进行比较。
//	y: 第二个切片，作为比较的目标。
//
// 返回:
//
//	类型为 []bool 的切片，表示 x 中相应位置的元素是否小于 y 的对应元素。
func ArrayLt[S ~[]T, T cmp.Ordered](x, y S) []bool {

	return Operator(func(x ...T) bool { return x[0] < x[1] }, x, y)
}

// ArrayGte 比较两个有序类型的切片中对应位置的元素，判断第一个切片中的元素是否大于等于第二个切片中的元素。
// 参数:
//
//	x: 第一个切片，其元素将与 y 中的元素进行比较。
//	y: 第二个切片，作为比较的目标。
//
// 返回:
//
//	类型为 []bool 的切片，表示 x 中相应位置的元素是否大于等于 y 的对应元素。
func ArrayGte[S ~[]T, T cmp.Ordered](x, y S) []bool {

	return Operator(func(x ...T) bool { return x[0] >= x[1] }, x, y)
}

// ArrayLte 比较两个切片中对应位置的元素，判断第一个切片中的元素是否小于等于第二个切片中的元素。
// 此泛型函数接受两个切片参数 x 和 y，其中切片元素必须实现 cmp.Ordered 接口，
// 表示它们可以进行比较排序。
//
// 参数:
//
//	x: 第一个切片，其元素将与 y 中的元素进行比较。
//	y: 第二个切片，作为比较的目标。
//
// 返回:
//
//	类型为 []bool 的切片，其中每个元素表示 x 中相应位置的元素是否小于等于 y 的对应元素。
func ArrayLte[S ~[]T, T cmp.Ordered](x, y S) []bool {

	return Operator(func(x ...T) bool { return x[0] <= x[1] }, x, y)
}

// ArrayBetween 比较三个切片中对应位置的元素，判断第一个切片中的元素是否在第二个和第三个切片元素之间（包括等于）。
// 此泛型函数接受三个切片参数 x, y, z，其中切片元素必须实现 cmp.Ordered 接口，
// 表示它们可以进行比较排序。
//
// 参数:
//
//	x: 第一个切片，其元素将与 y 和 z 中的元素进行比较。
//	y: 第二个切片，作为比较的下界。
//	z: 第三个切片，作为比较的上界。
//
// 返回:
//
//	类型为 []bool 的切片，其中每个元素表示 x 中相应位置的元素是否在 y 和 z 的对应元素之间（包括等于）。
func ArrayBetween[S ~[]T, T cmp.Ordered](x, y, z S) []bool {

	return Operator(func(x ...T) bool { return x[0] >= x[1] && x[0] <= x[2] }, x, y, z)
}

// Operator 对多个切片中的元素应用指定的函数，根据切片中元素的位置进行组合。
// 此泛型函数接受一个变长参数列表 `arr`，其中每个参数都是类型 S 的切片（S 必须是切片类型），
// 和一个函数 `fun`，该函数接受变长参数列表 x...T 并返回类型 T 的结果。
// 函数返回一个类型为 []T 的切片，其中包含应用给定函数 `fun` 到所有切片元素组合上的结果。
//
// 参数:
//
//	fun: 一个函数，接受变长参数列表 x...T，并返回类型 T 的结果。
//	arr: 变长参数列表，每个参数都是类型 S 的切片，S 必须是切片类型。
//
// 返回:
//
//	类型为 []T 的切片，包含应用给定函数 `fun` 到所有切片元素组合上的结果。
// func Operator[S ~[]T, T any](fun func(x ...T) T, arr ...S) []T {

// 	if len(arr) == 0 {
// 		return make([]T, 0)
// 	}

// 	la := ArrayMap(func(x ...S) int { return len(x[0]) }, arr)
// 	lm := ArrayMax(la)

// 	res := make([]T, lm)

// 	for i := 0; i < lm; i++ {
// 		parm := make([]T, len(arr))
// 		for j := 0; j < len(arr); j++ {
// 			parm[j] = arr[j][i%la[j]]
// 		}
// 		res[i] = fun(parm...)
// 	}

// 	return res
// }

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

	la := ArrayMap(func(x ...S) int { return len(x[0]) }, arr)
	lm := ArrayMax(la)

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
	rowNum := int(ArrayProduct(ArrayMap(func(x ...S) int { return len(x[0]) }, arr)))

	res := make([][]T, colNum)
	res[0] = arr[0]

	for i := 1; i < colNum; i++ {
		res[i] = ArrayRep(arr[i], len(res[i-1]), true)
	}

	for i := 0; i < colNum; i++ {
		if n := rowNum / len(res[i]); n > 1 {
			res[i] = ArrayRep(res[i], n, false)
		}

	}

	return ArrayZip(res...)
}

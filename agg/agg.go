package agg

import (
	"github.com/frankill/gotools"
	"github.com/frankill/gotools/array"
)

// Asccumulate 是一个高阶函数，用于对切片的元素进行累积计算。
// 它接受一个函数 fun 作为累积规则，一个默认值 default_v 作为累积的初始值，
// 和一个变长参数 arr，其中每个元素是一个切片。函数通过遍历所有切片的元素，
// 并应用累积函数 fun，最终返回累积的结果。
//
// 参数:
//
//	fun: 一个函数，定义了如何将一个累积值和当前元素结合起来产生新的累积值。
//	default_v: 累积的初始值。
//	arr: 变长参数，每个元素是一个切片，切片的元素类型为 T。
//
// 返回值:
//
//	U: 累积操作的结果类型，由 fun 函数定义。
//
// 使用示例:
//
//	func main() {
//	    gotools.Numbers := [][]int{{1, 2, 3}, {4, 5, 6}}
//	    sum := AsccumulateMap(func(x, y int) int { return x + y }, 0, gotools.Numbers)
//	    fmt.Println(sum) // 输出: 21，因为 0 + 1 + 2 + 3 + 4 + 5 + 6 = 21
//	}
func Asccumulate[S ~[]T, T, U any](fun func(x U, y T) U, default_v U, arr ...S) U {

	for _, v := range arr {
		for _, vv := range v {
			default_v = fun(default_v, vv)
		}
	}
	return default_v
}

// AsccumulateEach 对类型为 S（元素为 T 类型）的多个切片应用累计函数，并返回累计结果。
//
// 参数:
//
//	fun: 一个函数，定义了如何将一个累积值和当前元素结合起来产生新的累积值。
//	default_v: 累积的初始值。
//	arr: 变长参数，每个元素是一个切片，切片的元素类型为 T。
//
// 返回值:
//   - 一个 U 类型的累积结果。
//
// 由于此函数对多个切片元素应用函数，因此要求传递的切片必须具有相同的长度
func AsccumulateEach[S ~[]T, T, U any](fun func(x U, y ...T) U, default_v U, arr ...S) U {

	num := len(arr[0])
	pnum := len(arr)

	for i := 0; i < num; i++ {
		param := make([]T, pnum)

		for j := 0; j < pnum; j++ {
			param[j] = arr[j][i]
		}
		default_v = fun(default_v, param...)
	}

	return default_v
}

// ASum 计算多个切片中所有子切片元素的和。
// 该函数接受一个变长参数 slice，其中每个元素是一个切片，切片的元素类型必须实现 gotools.Number 接口。
// 函数返回所有子切片元素的总和。
// 使用 ArrayMap 和 ArraySum 函数进行嵌套操作，首先对每个子切片求和，然后对所有子切片的和求和。
// 参数:
//
//	slice ...[]T: 变长参数，其中每个元素是一个切片，切片的元素类型实现了 gotools.Number 接口。
//
// 返回值:
//
//	T: 所有子切片元素的总和，类型与切片的元素类型相同。
func ASum[S ~[]T, T gotools.Number](slice ...S) T {

	return Asccumulate(func(x, y T) T { return x + y }, T(0), slice...)
}

func Acount[S ~[]T, T any](slice ...S) int {
	return array.Sum(array.Map(func(x ...S) int { return len(x[0]) }, slice))
}

func ADistinct[S ~[]T, T comparable](slice ...S) int {
	return len(array.Unique(slice...))
}

// AMin 寻找多个切片中的最小元素。
// 该函数接受一个变长参数，其中每个参数是一个T类型的切片。
// 它通过先将每个切片的最小值找到，然后再从这些最小值中找出最终的最小值。
// 这样做的目的是为了在多个切片中找到全局最小值，而不是仅仅在一个切片中找到最小值。
// 参数:
//
//	slice ...[]T: 变长参数，每个参数是一个T类型的切片。
//
// 返回值:
//
//	T: 所有切片中的最小元素。
func AMin[S ~[]T, T gotools.Ordered](slice ...S) T {
	return array.Min(array.Map(func(x ...S) T { return array.Min(x[0]) }, slice))
}

// AMax 函数接收多个切片的切片作为输入，返回这些切片中的最大值。
// 它首先将每个切片中的最大值找到，然后在这些最大值中找出最终的最大值。
// 这个函数利用了泛型 T，使得它可以适用于任何实现了 gotools.Ordered 接口的类型。
// 参数:
//
//	slice ...[]T: 一个变长参数，包含了多个 T 类型的切片。
//
// 返回值:
//
//	T: 所有输入切片中的最大值。
func AMax[S ~[]T, T gotools.Ordered](slice ...S) T {
	return array.Max(array.Map(func(x ...S) T { return array.Max(x[0]) }, slice))
}

/*
AConcat 函数用于拼接多个切片。

参数:
- slice: 变长参数，每一项都是一个类型为 T 的切片，可以接受不同切片并将其合并。

返回值:
- 一个新切片，包含了输入的所有切片中的元素。新切片的类型与输入切片元素类型相同。

功能说明:
此函数通过泛型 T 支持任意数据类型的切片拼接，它内部调用 ArrayConcat 函数来完成实际的拼接操作。
*/
func AConcat[S ~[]T, T any](slice ...S) []T {

	return array.Concat(slice...)

}

// AMaxif 根据提供的条件函数从多个切片中找出满足条件的最大元素。
//
// 参数:
// - fun: 一个函数，接受可变数量的T类型参数并返回一个布尔值，用于判断元素是否满足条件。
// - slice: 可变数量的切片参数，所有切片的元素类型需为T，且T需要实现Ordered接口。
//
// 返回值:
// - T: 满足条件的所有输入切片元素中的最大值。
func AMaxif[S ~[]T, T gotools.Ordered](fun func(x ...T) bool, slice ...S) T {

	a := array.Map(func(x ...S) S { return array.Filter(fun, x[0]) }, slice)

	return AMax(a...)
}

// AMinif 是一个泛型函数，用于从一个切片中找到满足特定条件的最小元素。
// 它接受一个函数 fun 作为条件判断，和一个或多个切片 slice 作为待检查的集合。
// 函数 fun 用于测试切片中的元素是否满足某种条件，返回一个布尔值。
// AMinif 返回满足条件的最小元素，前提是切片中的元素类型必须实现了 gotools.Ordered 接口。
//
// 参数:
//
//	fun: 一个函数，接受一个或多个 T 类型的参数，并返回一个布尔值。
//	     该函数用于判断元素是否满足某种条件。
//	slice: 一个或多个切片，它们的元素类型必须是 T，并且 T 必须实现了 gotools.Ordered 接口。
//
// 返回值:
//
//	返回满足条件的最小元素，类型为 T。
func AMinif[S ~[]T, T gotools.Ordered](fun func(x ...T) bool, slice ...S) T {

	a := array.Map(func(x ...S) S { return array.Filter(fun, x[0]) }, slice)

	return AMin(a...)
}

// AargMax 函数在给定的值数组（val）中找到对应的键（arg）的最大值。
// 当值数组（val）为空时，返回键数组（arg）的第一个元素。
// 当键数组（arg）为空时，返回默认值（res）。
// 值数组（val）中的元素需要实现gotools.Ordered接口，键和值可以是任意类型。
//
// 参数:
//   - arg: D 类型，键数组，对应值数组（val）中的索引。
//   - val: S 类型，值数组，根据其元素大小进行比较。
//
// 返回:
//   - U 类型，找到的最大值对应的键
func AargMax[D ~[]U, S ~[]T, T gotools.Ordered, U any](arg D, val S) U {

	if len(val) == 0 {
		return arg[0]
	}
	var res U
	if len(arg) == 0 {
		return res
	}

	if len(arg) != len(val) {
		val = val[0:len(arg)]
	}

	index := array.FindMax(val)

	if index == -1 {
		return res
	}

	res = arg[index]
	return res

}

// AargMin 查找与目标序列中最小元素对应的数据源序列中的元素。
//
// 参数:
//   - arg(D): 数据源序列，类型为切片，元素类型为 U。
//   - val(S): 目标序列，类型为切片，元素类型需满足有序比较（gotools.Ordered），用于寻找最小值。
//
// 返回值:
//   - U: 与目标序列中最小元素对应的数据源序列中的元素。如果目标序列为空，返回数据源序列的第一个元素。
//     如果数据源序列为空，返回 U 类型的零值。
//
// 注意:
//   - 函数利用泛型支持多种数据类型的操作，D 和 S 分别约束为切片类型，T 需实现 Ordered 接口，
//     而 U 可以为任何类型。
func AargMin[D ~[]U, S ~[]T, T gotools.Ordered, U any](arg D, val S) U {

	if len(val) == 0 {
		return arg[0]
	}
	var res U
	if len(arg) == 0 {
		return res
	}

	if len(arg) != len(val) {
		val = val[0:len(arg)]
	}

	index := array.FindMin(val)

	if index == -1 {
		return res
	}

	res = arg[index]
	return res

}

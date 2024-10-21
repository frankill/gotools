package array

// Sort 对类型为 S（元素类型为 T）的切片进行自定义排序。
//
// 参数:
//   - fun: 一个比较函数，接受两个 T 类型的参数并返回一个布尔值，指示是否需要交换这两个参数的位置。
//     当 `fun(x, y)` 返回 `true` 时，在排序过程中 `x` 应该排在 `y` 之前。
//   - arr: 要排序的切片 S。
//
// 返回值:
//   - 返回一个新的 S 类型切片，其中的元素根据提供的比较函数 `fun` 进行排序。
func Sort[S ~[]T, T any](fun func(x, y T) bool, arr S) []T {

	la := len(arr)
	if la == 0 {
		return []T{}
	}
	if la == 1 {
		return []T{arr[0]}
	}

	res := make([]T, la)
	copy(res, arr)

	for i := 1; i < la; {
		if i == 0 {
			i = 2
		}
		if i >= la {
			break
		}
		if fun(res[i], res[i-1]) {
			res[i], res[i-1] = res[i-1], res[i]
			i--
			continue
		}
		i++

	}

	return res

}

// SortQuick 使用快速排序的思路（尽管实现并不完全正确）尝试对类型为 S（元素类型为 T）的切片进行原地排序。
// 注意：当前实现并非标准快速排序算法，更像是简化冒泡排序变种。
//
// 参数:
//    - arr: 要排序的切片 S，其中元素类型 T 必须是可比较的（实现 gotools.Ordered 接口）。
// func SortQuick[S ~[]T, T gotools.Ordered](arr S) []T {

// 	return ArraySort(func(x, y T) bool { return x < y }, arr)

// }

// SortL 对类型为 S（元素类型为 T）的切片进行原地排序，依据提供的比较函数 `fun`。
//
// 参数:
//   - fun: 自定义比较函数，接受两个 T 类型的参数并返回一个布尔值。
//     当 `fun(x, y)` 返回 `true`，则在排序时 `x` 应位于 `y` 之前。
//   - arr: 要排序的本地切片 S，函数会直接修改传入的切片。
func SortL[S ~[]T, T any](fun func(x, y T) bool, arr S) {
	la := len(arr)
	if la == 0 {
		return
	}
	if la == 1 {
		return
	}

	for i := 1; i < la; {
		if i == 0 {
			i = 2
		}
		if i >= la {
			break
		}
		if fun(arr[i], arr[i-1]) {
			arr[i], arr[i-1] = arr[i-1], arr[i]
			i--
			continue
		}
		i++

	}

}

// func SortByQ[D ~[]U, S ~[]T, T any, U gotools.Ordered](arr S, order D) (S, D) {

// 	return ArraySortBy(func(x, y U) bool { return x < y }, arr, order)

// }

// SortTwo 通过自定义比较函数和排序顺序对数组进行排序。
// 参数:
//
//   - fun: 自定义比较函数，接受两个 T 类型的参数并返回一个布尔值。
//   - arr: 需要排序的数组，类型为 S。
//   - order: 与 arr 相关的数组，类型为 D，其元素类型与 fun 的参数类型相同。
//
// 返回值:
//
//   - 返回值是排序后的数组。
func SortBy[D ~[]U, S ~[]T, T, U any](fun func(x, y U) bool, arr S, order D) (S, D) {

	la := len(arr)
	if la == 0 {
		return S{}, D{}
	}
	if la == 1 {
		return arr, order
	}

	res := make(S, la)
	tmp := make(D, la)
	copy(res, arr)
	copy(tmp, order)

	for i := 1; i < la; {
		if i == 0 {
			i = 2
		}
		if i >= la {
			break
		}
		if fun(tmp[i], tmp[i-1]) {
			tmp[i], tmp[i-1] = tmp[i-1], tmp[i]
			res[i], res[i-1] = res[i-1], res[i]
			i--
			continue
		}
		i++

	}

	return res, tmp
}

// func SortByLQ[D ~[]U, S ~[]T, T any, U gotools.Ordered](arr S, order D) {
// 	ArraySortByL(func(current, before U) bool { return current < before }, arr, order)
// }

// SortTwoLocal 是一个泛型排序函数，用于对两个关联数组进行排序。
// 它接受一个比较函数、一个需要排序的数组和一个与之相关的数组，
// 根据比较函数的规则对相关数组进行排序。
// 比较函数接收两个元素并返回一个布尔值，表示第一个元素是否应排在第二个元素之前。
// 注意：这个函数会直接修改输入的数组。
//
// 参数:
//   - fun: 比较函数，用于确定元素的排序顺序。
//   - arr: 需要排序的数组，类型为 S。
//   - order: 与 arr 相关的数组，类型为 D，其元素类型与 fun 的参数类型相同。
//
// 返回:
//   - 无返回值，直接修改输入的 arr 和 order 数组。
func SortByL[D ~[]U, S ~[]T, T, U any](fun func(x, y U) bool, arr S, order D) {
	la := len(arr)
	if la == 0 {
		return
	}
	if la == 1 {
		return
	}

	for i := 1; i < la; {
		if i == 0 {
			i = 2
		}
		if i >= la {
			break
		}
		if fun(order[i], order[i-1]) {
			order[i], order[i-1] = order[i-1], order[i]
			arr[i], arr[i-1] = arr[i-1], arr[i]
			i--
			continue
		}
		i++

	}

}

package array

import "github.com/frankill/gotools"

// SortFun 通过比较函数对切片进行排序
//
// 参数:
//   - f: 一个函数，接受两个 T 类型的变长参数并返回布尔值，决定是否交换两个元素的位置。true 表示交换，false 表示不交换。
//   - arr: 一个切片，表示要进行排序的切片。
//
// 返回:
//   - 一个新的切片，表示排序后的结果。
func SortFun[S ~[]T, T any](f func(x T, y T) bool, arr S) []T {

	res := make([]T, len(arr))

	copy(res, arr)

	SortFunLocal(f, res)

	return res

}

// SortFunLocal 通过比较函数对切片进行排序，修改原切片
//
// 参数:
//   - f: 一个函数，接受两个 T 类型的变长参数并返回布尔值，决定是否交换两个元素的位置。true 表示交换，false 表示不交换。
//   - arr: 一个切片，表示要进行排序的切片。
//
// 无返回值
func SortFunLocal[S ~[]T, T any](f func(x T, y T) bool, arr S) {
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
		if f(arr[i], arr[i-1]) {
			arr[i], arr[i-1] = arr[i-1], arr[i]
			i--
			continue
		}
		i++

	}

}

// OrderFun 通过比较函数对切片进行排序
//
// 参数:
//   - f: 一个函数，接受两个 T 类型的变长参数并返回布尔值，决定是否交换两个元素的位置。true 表示交换，false 表示不交换。
//   - arr: 一个切片，表示要进行排序的切片。
//   - order: 一个切片，表示排序后的顺序。
//
// 返回:
//   - 一个新的切片，表示排序后的结果。
func OrderFun[D ~[]U, S ~[]T, T any, U any](f func(x, y U) bool, arr S, order D) (S, D) {

	if len(arr) != len(order) {
		panic("length of arr and order must be equal")
	}

	res := make([]T, len(arr))
	copy(res, arr)

	resOrder := make([]U, len(order))
	copy(resOrder, order)

	OrderFunLocal(f, res, resOrder)

	return res, resOrder

}

// OrderFunLocal 通过比较函数对切片进行排序，修改原切片
//
// 参数:
//   - f: 一个函数，接受两个 T 类型的变长参数并返回布尔值，决定是否交换两个元素的位置。true 表示交换，false 表示不交换。
//   - arr: 一个切片，表示要进行排序的切片。
//   - order: 一个切片，表示排序后的顺序。
//
// 无返回值
func OrderFunLocal[D ~[]U, S ~[]T, T any, U any](f func(x, y U) bool, arr S, order D) {

	if len(arr) != len(order) {
		panic("length of arr and order must be equal")
	}

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
		if f(order[i], order[i-1]) {
			order[i], order[i-1] = order[i-1], order[i]
			arr[i], arr[i-1] = arr[i-1], arr[i]
			i--
			continue
		}
		i++

	}

}

// Sort 对切片进行排序, 切片类型必须实现 Ordered 接口
//
// 参数:
//   - arr: 一个切片，表示要进行排序的切片。
//   - descending: 一个布尔值，表示是否降序排序。
//
// 返回:
//   - 一个新的切片，表示排序后的结果。
func Sort[S ~[]T, T gotools.Ordered](arr S, descending bool) []T {

	res := make([]T, len(arr))

	copy(res, arr)

	if descending {
		SortR(res)
	} else {
		SortL(res)
	}

	return res

}

// SortL 对切片进行升序排序, 切片类型必须实现 Ordered 接口
//
// 参数:
//   - arr: 一个切片，表示要进行排序的切片。
//
// 无返回值
func SortL[S ~[]T, T gotools.Ordered](arr S) {
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
		if arr[i] < arr[i-1] {
			arr[i], arr[i-1] = arr[i-1], arr[i]
			i--
			continue
		}
		i++

	}

}

// SortR 对切片进行降序排序, 切片类型必须实现 Ordered 接口
//
// 参数:
//   - arr: 一个切片，表示要进行排序的切片。
//
// 无返回值
func SortR[S ~[]T, T gotools.Ordered](arr S) {
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
		if arr[i] > arr[i-1] {
			arr[i], arr[i-1] = arr[i-1], arr[i]
			i--
			continue
		}
		i++

	}

}

// Order 对切片进行排序，order切片类型必须实现 Ordered 接口
//
// 参数:
//   - arr: 一个切片，表示要进行排序的切片。
//   - order: 一个切片，表示排序后的顺序。
//   - descending: 一个布尔值，表示是否降序排序。
//
// 返回:
//   - 一个新的切片，表示排序后的结果。
//   - 一个新的切片，表示排序后的顺序。
func Order[D ~[]U, S ~[]T, T any, U gotools.Ordered](arr S, order D, descending bool) (S, D) {

	if len(arr) != len(order) {
		panic("length of arr and order must be equal")
	}

	res := make([]T, len(arr))
	copy(res, arr)

	resOrder := make([]U, len(order))
	copy(resOrder, order)

	if descending {
		OrderR(res, resOrder)
	} else {
		OrderL(res, resOrder)
	}

	return res, resOrder

}

// OrderL 对切片进行升序排序，order切片类型必须实现 Ordered 接口
//
// 参数:
//   - arr: 一个切片，表示要进行排序的切片。
//   - order: 一个切片，表示排序后的顺序。
//
// 无返回值
func OrderL[D ~[]U, S ~[]T, T any, U gotools.Ordered](arr S, order D) {

	if len(arr) != len(order) {
		panic("length of arr and order must be equal")
	}

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
		if order[i] < order[i-1] {
			order[i], order[i-1] = order[i-1], order[i]
			arr[i], arr[i-1] = arr[i-1], arr[i]
			i--
			continue
		}
		i++

	}

}

// OrderR 对切片进行降序排序，order切片类型必须实现 Ordered 接口
//
// 参数:
//   - arr: 一个切片，表示要进行排序的切片。
//   - order: 一个切片，表示排序后的顺序。
//
// 无返回值
func OrderR[D ~[]U, S ~[]T, T any, U gotools.Ordered](arr S, order D) {

	if len(arr) != len(order) {
		panic("length of arr and order must be equal")
	}

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
		if order[i] > order[i-1] {
			order[i], order[i-1] = order[i-1], order[i]
			arr[i], arr[i-1] = arr[i-1], arr[i]
			i--
			continue
		}
		i++

	}

}

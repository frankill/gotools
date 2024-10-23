package iter

import "github.com/frankill/gotools"

// UniqueCount 去重, 要求传入的ch必须是排序过的
// 参数:
//   - f: 一个函数，接受两个类型为 T 的值，返回一个布尔值，表示是否相等。
//   - ch: 一个通道，用于接收数据。
//
// 返回:
//   - 去重后的数量
func UniqueCount[T any](f func(x, y T) bool) func(ch chan T) int {

	return func(ch chan T) int {

		count := 0

		prev, ok := <-ch
		if !ok {
			return count
		} else {
			count++
		}

		for v := range ch {
			if !f(prev, v) {
				count++
				prev = v
			}
		}
		return count
	}
}

// Count 计算通道中元素的数量
// 参数:
//   - ch: 一个通道，用于接收数据。
//
// 返回:
//   - int: 通道中元素的数量。
func Count[T any](ch chan T) int {
	count := 0
	for range ch {
		count++
	}
	return count
}

// DistinctCount 去重计数
// 参数:
//   - ch: 一个通道，用于接收数据。
//
// 返回:
//   - int: 通道中去重后元素的数量。
func DistinctCount[T comparable](ch chan T) int {
	count := 0
	seen := make(map[T]struct{})
	for v := range ch {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			count++
		}
	}
	return count
}

// Reduce 对通道中的数据进行归约操作，将通道中的所有数据合并成一个值。
// 参数:
//   - f: 一个函数，接受两个参数：一个类型为 U 的累加器和一个类型为 T 的当前值，返回一个类型为 U 的新累加器。
//   - init: 初始值，类型为 U，用作归约操作的起始累加器。
//   - ch: 一个通道，通道中的每个值是类型为 T 的数据，将被用于归约操作。
//
// 返回:
//   - 返回归约操作后的结果，类型为 U。
//
// 函数功能:
//   - 从输入通道 ch 中读取数据，并将每个数据应用归约函数 f，使用 init 作为初始值。
//   - 在处理完所有数据后，返回最终的累加器值。
func Reduce[T, U any](f func(x U, y T) U, init U) func(ch chan T) U {

	return func(ch chan T) U {
		for v := range ch {
			init = f(init, v)
		}
		return init
	}

}

// Find 从通道中查找第一个满足条件的值。
// 参数:
//   - f: 一个函数，接受一个类型为 T 的值，返回一个布尔值，表示该值是否满足条件。
//   - ch: 一个通道，通道中的值是类型为 T 的数据。
//
// 返回:
//   - 第一个满足条件的值。如果通道中的值都不满足条件，返回类型 T 的零值。
//
// 函数功能:
//   - 从输入通道 ch 中读取数据，应用函数 f 来检查每个值是否满足条件。
//   - 如果找到第一个满足条件的值，则返回该值。
//   - 如果通道关闭且没有找到满足条件的值，则返回类型 T 的零值。
func First[T any](f func(x T) bool) func(ch chan T) T {

	return func(ch chan T) T {

		var result T
		var ok bool
		for v := range ch {
			if !ok && f(v) {
				ok = true
				result = v
			} else {
				continue
			}
		}

		return result
	}
}

// Last 从通道中查找最后一个满足条件的值。
// 参数:
//   - f: 一个函数，接受一个类型为 T 的值，返回一个布尔值，表示该值是否满足条件。
//   - ch: 一个通道，通道中的值是类型为 T 的数据。
//
// 返回:
//   - 最后一个满足条件的值。如果通道中的值都不满足条件，返回类型 T 的零值。
//
// 函数功能:
//   - 从输入通道 ch 中读取数据，应用函数 f 来检查每个值是否满足条件。
//   - 如果找到最后一个满足条件的值，则返回该值。
//   - 如果通道关闭且没有找到满足条件的值，则返回类型 T 的零值。
func Last[T any](f func(x T) bool) func(ch chan T) T {
	return func(ch chan T) T {

		var result T
		for v := range ch {
			if f(v) {
				result = v
			}
		}
		return result
	}
}

// Collect 从通道中收集所有值并返回一个切片。
// 参数:
//   - ch: 一个通道，通道中的值是类型为 T 的数据。
//
// 返回:
//   - 一个切片，其中包含从通道中读取的所有值。
//
// 函数功能:
//   - 从输入通道 ch 中读取所有值，并将它们按顺序添加到一个切片中。
//   - 当通道关闭时，函数返回包含所有读取值的切片。
func Collect[T any](ch chan T) []T {
	var result []T
	for v := range ch {
		result = append(result, v)
	}
	return result
}

func CollectFun[T any, U gotools.Comparable](f func(x T) U) func(ch chan T) []U {

	return func(ch chan T) []U {

		var result []U
		for v := range ch {
			result = append(result, f(v))
		}
		return result

	}

}

// Walk 对通道中的数据进行遍历操作。
// 参数:
//   - f: 一个函数，接受类型为 T 的输入。
//   - ch: 一个通道，通道中的每个值是类型为 T 的数据。
//
// 函数功能:
//   - 从输入通道 ch 中读取数据，对每个数据应用函数 f。
func Walk[T any](f func(x T)) func(ch chan T) {

	return func(ch chan T) {

		for v := range ch {
			f(v)
		}
	}

}

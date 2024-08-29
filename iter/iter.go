package iter

import (
	"cmp"
	"sync"

	"github.com/frankill/gotools/array"
)

var (
	BufferSize = 100
)

func Identity[T any](x T) T {
	return x
}

// Map 将通道中的每个元素应用函数 f，并将结果发送到一个新的通道。
// 参数:
//   - f: 一个函数，接受类型为 T 的输入，返回类型为 U 的结果。
//   - ch: 一个通道，通道中的每个值是类型为 T 的数据。
//
// 返回:
//   - 一个通道，通道中的值是类型为 U 的数据，每个值是应用函数 f 之后的结果。
//
// 函数功能:
//   - 从输入通道 ch 中读取数据，将每个数据应用函数 f，然后将结果写入新的通道 ch_。
//   - 使用一个 goroutine 执行这些操作，并在完成后关闭通道 ch_。
func Map[T any, U any](f func(x T) U) func(ch chan T) chan U {

	return func(ch chan T) chan U {

		ch_ := make(chan U, BufferSize)

		go func() {
			defer close(ch_)
			for v := range ch {
				ch_ <- f(v)
			}
		}()

		return ch_
	}
}

// FlatMap 返回一个处理步骤函数，该函数接受一个输入通道（chan T），
// 对通道中的每个元素应用函数 f，并将结果展平到一个新的输出通道（chan U）中。
// 函数 f 接受类型 T 的元素并返回类型 U 的切片。函数返回的处理步骤函数创建一个新的通道
// 用于发送展平后的结果。返回的输出通道在处理完成后会被关闭。
func FlatMap[T any, U any](f func(x T) []U) func(ch chan T) chan U {
	return func(ch chan T) chan U {
		ch_ := make(chan U, BufferSize)
		go func() {
			defer close(ch_)
			for v := range ch {
				for _, value := range f(v) {
					ch_ <- value
				}
			}
		}()
		return ch_
	}
}

// Distinct 返回一个处理步骤函数，该函数接受一个输入通道（chan T），
// 从中移除重复的元素，并将唯一的元素发送到一个新的输出通道（chan T）中。
// 函数内部使用一个映射（map）来跟踪已经见过的元素，以确保每个元素只出现一次。
// 返回的输出通道在处理完成后会被关闭。
func Distinct[T comparable](ch chan T) chan T {
	ch_ := make(chan T, BufferSize)
	set := make(map[T]struct{})
	go func() {
		defer close(ch_)
		for v := range ch {
			if _, ok := set[v]; !ok {
				set[v] = struct{}{}
				ch_ <- v
			}
		}
	}()
	return ch_
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

// Filter 过滤通道中的数据，只将符合条件的数据发送到新的通道。
// 参数:
//   - f: 一个函数，接受类型为 T 的输入，并返回布尔值。如果返回 true，则将该数据发送到新通道；如果返回 false，则忽略该数据。
//   - ch: 一个通道，通道中的每个值是类型为 T 的数据。
//
// 返回:
//   - 一个通道，通道中的值是类型为 T 的数据，只包含符合函数 f 条件的数据。
//
// 函数功能:
//   - 从输入通道 ch 中读取数据，对每个数据应用函数 f。如果函数 f 返回 true，则将数据写入新的通道 ch_。
//   - 使用一个 goroutine 执行这些操作，并在完成后关闭通道 ch_。
func Filter[T any](f func(x T) bool) func(ch chan T) chan T {

	return func(ch chan T) chan T {

		ch_ := make(chan T, BufferSize)

		go func() {
			defer close(ch_)

			for v := range ch {
				if f(v) {
					ch_ <- v
				}
			}
		}()
		return ch_
	}
}

// Sequence 生成一个整数序列，并将其发送到通道。
// 参数:
//   - start: 序列的起始值，整数类型。
//   - end: 序列的结束值，整数类型。序列生成会在达到该值时停止（不包括此值）。
//   - step: 序列的步长，整数类型。可以为正数或负数。
//
// 返回:
//   - 一个通道，通道中的值是整数类型，表示生成的序列。
//
// 函数功能:
//   - 从 start 开始，根据步长 step 生成整数序列，直到达到 end 为止（不包括 end）。
//   - 使用一个 goroutine 执行这些操作，并在完成后关闭通道 ch。
func Sequence(start, end, step int) chan int {

	ch := make(chan int, BufferSize)

	go func() {
		defer close(ch)
		for i := start; (step > 0 && i < end) || (step < 0 && i > end); i += step {
			ch <- i
		}

	}()
	return ch
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

// Scanl 对通道中的数据进行扫描操作，生成一个累加序列，并将结果发送到新的通道。
// 参数:
//   - f: 一个函数，接受两个参数：一个类型为 U 的累加器和一个类型为 T 的当前值，返回一个类型为 U 的新累加器。
//   - init: 初始值，类型为 U，用作扫描操作的起始累加器。
//   - ch: 一个通道，通道中的每个值是类型为 T 的数据，将用于扫描操作。
//
// 返回:
//   - 一个通道，通道中的值是类型为 U 的数据，表示扫描操作的累加序列。
//
// 函数功能:
//   - 从输入通道 ch 中读取数据，并将每个数据应用扫描函数 f，使用 init 作为初始值。
//   - 将每个步骤的累加器值写入新的通道 ch_。
//   - 使用一个 goroutine 执行这些操作，并在完成后关闭通道 ch_。
func Scanl[T, U any](f func(x U, y T) U, init U) func(ch chan T) chan U {

	return func(ch chan T) chan U {

		ch_ := make(chan U, BufferSize)

		go func() {
			defer close(ch_)

			for v := range ch {
				init = f(init, v)
				ch_ <- init
			}

		}()

		return ch_
	}
}

// Zip 结合两个通道中的值，生成一个新的通道。
// 参数:
//   - f: 一个函数，接受两个参数：一个类型为 T 的值和一个类型为 U 的值，返回一个类型为 V 的值。
//   - ch1: 一个通道，通道中的值是类型为 T 的数据。
//   - ch2: 一个通道，通道中的值是类型为 U 的数据。
//
// 返回:
//   - 一个通道，通道中的值是类型为 V 的数据，表示将 ch1 和 ch2 中的值通过函数 f 组合后的结果。
//
// 函数功能:
//   - 从通道 ch1 和 ch2 中读取数据，并使用函数 f 将这两个值组合成一个新值。
//   - 将生成的新值发送到新的通道 ch 中。
//   - 当两个通道都关闭时，停止处理并关闭通道 ch。
func Zip[T any, U any, V any](f func(x T, y U) V) func(ch1 chan T, ch2 chan U) chan V {

	return func(ch1 chan T, ch2 chan U) chan V {
		ch := make(chan V, BufferSize)

		go func() {
			defer close(ch)
			for {
				v1, ok1 := <-ch1
				v2, ok2 := <-ch2

				if !ok1 && !ok2 {
					return
				}

				ch <- f(v1, v2)
			}
		}()

		return ch
	}

}

// Partition 根据给定的条件函数将通道中的值分为两个通道。
// 参数:
//   - f: 一个函数，接受一个类型为 T 的值，返回一个布尔值，表示该值是否满足条件。
//   - ch: 一个通道，通道中的值是类型为 T 的数据。
//
// 返回:
//   - ch1: 一个通道，通道中的值是类型为 T 的数据，表示满足条件的值。
//   - ch2: 一个通道，通道中的值是类型为 T 的数据，表示不满足条件的值。
//
// 函数功能:
//   - 从输入通道 ch 中读取数据，根据函数 f 的结果将每个值分配到两个新的通道 ch1 和 ch2。
//   - 当输入通道 ch 关闭时，关闭新的通道 ch1 和 ch2。
func Partition[T any](f func(x T) bool) func(ch chan T) (chan T, chan T) {

	return func(ch chan T) (chan T, chan T) {
		ch1 := make(chan T, BufferSize)
		ch2 := make(chan T, BufferSize)

		go func() {

			defer close(ch1)
			defer close(ch2)

			for v := range ch {
				if f(v) {
					ch1 <- v
				} else {
					ch2 <- v
				}
			}
		}()
		return ch1, ch2
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
		for v := range ch {
			if f(v) {
				return v
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

func CollectFun[T any, U comparable](f func(x T) U) func(ch chan T) []U {

	return func(ch chan T) []U {

		var result []U
		for v := range ch {
			result = append(result, f(v))
		}
		return result

	}

}

// TakeWhile 从通道中读取值，直到遇到第一个不满足条件的值。
// 参数:
//   - f: 一个函数，接受一个类型为 T 的值，返回一个布尔值，表示该值是否满足条件。
//   - ch: 一个通道，通道中的值是类型为 T 的数据。
//
// 返回:
//   - 一个新通道，其中包含从输入通道中读取的所有满足条件的值。
//
// 函数功能:
//   - 从输入通道 ch 中读取值，应用函数 f 来检查每个值是否满足条件。
//   - 将所有满足条件的值写入新通道 ch_。
//   - 当遇到第一个不满足条件的值时，停止读取并关闭新通道 ch_。
//   - 新通道 ch_ 只包含在遇到第一个不满足条件的值之前的所有值
func TakeWhile[T any](f func(x T) bool) func(ch chan T) chan T {

	return func(ch chan T) chan T {
		ch_ := make(chan T, BufferSize)

		go func() {
			defer close(ch_)
			for v := range ch {
				if f(v) {
					ch_ <- v
				} else {
					return
				}
			}
		}()
		return ch_
	}
}

// DropWhile 从通道中读取值，直到遇到第一个不满足条件的值，之后将所有后续值传递到新通道。
// 参数:
//   - f: 一个函数，接受一个类型为 T 的值，返回一个布尔值，表示该值是否满足条件。
//   - ch: 一个通道，通道中的值是类型为 T 的数据。
//
// 返回:
//   - 一个新通道，其中包含从第一个不满足条件的值之后的所有值。
//
// 函数功能:
//   - 从输入通道 ch 中读取值，应用函数 f 来检查每个值是否满足条件。
//   - 跳过所有满足条件的值，直到遇到第一个不满足条件的值。
//   - 从第一个不满足条件的值开始，将所有后续值写入新通道 ch_。
//   - 新通道 ch_ 包含从第一个不满足条件的值之后的所有值。
func DropWhile[T any](f func(x T) bool) func(ch chan T) chan T {
	return func(ch chan T) chan T {
		ch_ := make(chan T, BufferSize)

		num := 0
		go func() {
			defer close(ch_)
			for v := range ch {

				if num == 0 {
					if f(v) {
						continue
					}
					num = 1
				}

				if num == 1 {
					ch_ <- v
				}

			}
		}()
		return ch_
	}
}

// Merge 将多个通道合并为一个通道。
//
// 参数:
//
//	chs: 一个或多个需要合并的通道，这些通道的数据类型必须相同。
//
// 返回值:
//
//	一个通道，该通道将接收来自所有输入通道的数据。
//
// 说明:
//
//	Merge 函数会启动多个 goroutine，每个 goroutine 从一个输入通道中读取数据
//	并将数据写入到返回的通道中。所有输入通道的数据都将被合并到这个返回的通道中。
//	当所有输入通道的数据都被读取完毕后，返回的通道将会被关闭。
func Union[T any](chs ...chan T) chan T {

	out := make(chan T, BufferSize)

	var wg sync.WaitGroup

	for _, ch := range chs {
		wg.Add(1)
		go func(c chan T) {
			defer wg.Done()
			for v := range c {
				out <- v
			}
		}(ch)
	}

	go func() {
		defer close(out)
		wg.Wait()
	}()

	return out
}

func Intersection[T comparable](ch1 chan T, ch2 chan T) chan T {

	ch := make(chan T, BufferSize)

	go func() {

		m := array.ArrayToMap(Collect(ch2))

		for v := range ch1 {

			if _, ok := m[v]; ok {
				ch <- v
			}

		}

	}()

	return ch

}

func Subtract[T comparable](ch1 chan T, ch2 chan T) chan T {

	ch := make(chan T, BufferSize)

	go func() {

		m := array.ArrayToMap(Collect(ch2))

		for v := range ch1 {

			if _, ok := m[v]; !ok {
				ch <- v
			}

		}

	}()

	return ch
}

func Cartesian[T any](ch1 chan T, ch2 chan T) chan array.Pair[T, T] {

	ch := make(chan array.Pair[T, T], BufferSize)

	dd := Collect(ch2)

	var wg sync.WaitGroup
	wg.Add(2)

	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()
			for v1 := range ch1 {
				for _, v2 := range dd {
					ch <- array.PairOf(v1, v2)
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch

}

func GroupBy[U cmp.Ordered, T any, K comparable](by chan K, data chan T, order chan U) chan array.Pair[K, []T] {
	ch := make(chan array.Pair[K, []T], BufferSize)

	go func() {
		defer close(ch)

		// 收集所有通道的数据
		keys := Collect(by)
		values := Collect(data)
		orders := Collect(order)

		// 使用 GroupByOrder 分组数据
		group := array.GroupByOrder(keys, values, orders)

		// 将分组结果发送到结果通道
		for k, v := range group {
			ch <- array.Pair[K, []T]{
				First:  k,
				Second: v,
			}
		}
	}()

	return ch
}

// Window 滚动窗口函数
// 参数:
//   - windowSize: 窗口的大小。
//   - ch: 输入通道，包含要处理的数据。
//
// 返回:
//   - 输出通道，其中包含滚动窗口的数据切片。
func Window[T any](windowSize int, ch chan T) chan []T {
	window := make([]T, 0, windowSize)
	out := make(chan []T, BufferSize)

	go func() {
		defer close(out)
		for v := range ch {
			window = append(window, v)
			if len(window) == windowSize {
				// 发送当前窗口的副本到通道
				out <- window
				// 重新创建窗口切片
				window = make([]T, 0, windowSize)
			}
		}
		// 处理剩余的窗口数据
		if len(window) > 0 {
			out <- window
		}
	}()

	return out
}

// Split 将通道中的数据分组
// 参数:
//   - fn: 一个函数，接受一个类型为 T 的值，返回一个整数值，表示该值属于哪个分组。
//   - num: 分组的数量。
// 返回:
//   - 一个函数，接受一个通道，返回一个包含 num 个通道的切片。

func Split[T any](fn func(T) int, num int) func(ch chan T) []chan T {
	// 创建一个包含 num 个通道的切片
	a := make([]chan T, num)
	for i := 0; i < num; i++ {
		a[i] = make(chan T, BufferSize)
	}

	return func(ch chan T) []chan T {
		go func() {
			defer func() {
				for _, v := range a {
					close(v)
				}
			}()

			for v := range ch {
				// 确保 fn(v) 不超出通道切片的范围
				index := fn(v)
				if index >= 0 && index < num {
					a[index] <- v
				}
			}
		}()

		return a
	}
}

// Maps 将多个通道中的数据进行映射
// 参数:
//   - fn: 一个函数，接受一个类型为 T 的值 ，类型为int的通道位置，返回一个类型为 U 的值。
//   - cs: 一个包含多个通道的切片。
//
// 返回:
//   - 一个通道，用于接收映射后的数据。
func Maps[T any, U any](fn func(ch chan T, num int) U) func(cs ...chan T) chan U {
	return func(cs ...chan T) chan U {
		out := make(chan U, len(cs)) // 创建接收结果的通道
		var wg sync.WaitGroup

		go func() {
			defer close(out)
			for loc, c := range cs {
				wg.Add(1)
				go func(index int, ch chan T) {
					defer wg.Done()
					// 将 fn 结果发送到 out 通道
					out <- fn(ch, index)
				}(loc, c)
			}
			wg.Wait()
		}()

		return out
	}
}

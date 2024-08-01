package iter

import (
	"cmp"
	"log"
	"sync"

	"github.com/frankill/gotools/array"
)

var (
	BufferSize = 10
)

// Parallel 允许并行运行多个函数，并控制同时运行的最大协程数
type Parallel struct {
	arr []func()       // 存储要执行的函数的切片
	wg  sync.WaitGroup // 用于同步 goroutine 的 WaitGroup
	num int            // 同时运行的最大协程数
	sem chan struct{}  // 信号量，用于限制同时运行的协程数量
}

// NewParallel 创建一个新的 Parallel 实例，指定最大并发数
func NewParallel(maxConcurrency int) *Parallel {
	return &Parallel{
		num: maxConcurrency,
		sem: make(chan struct{}, maxConcurrency),
	}
}

// Add 添加一个新的函数到并行执行列表
func (p *Parallel) Add(f func()) {
	p.arr = append(p.arr, f)
}

// Compute 并行执行所有添加的函数，并等待它们完成
func (p *Parallel) Compute() {
	for _, f := range p.arr {
		p.sem <- struct{}{}
		p.wg.Add(1)
		go func(f func()) {
			defer p.wg.Done()
			defer func() { <-p.sem }()
			f()
		}(f)
	}

	p.wg.Wait()
}

// Pipeline 表示一系列处理步骤，可以对数据进行处理
type Pipeline[T any] struct {
	steps []func(chan T) chan T // 表示处理步骤的函数切片
}

// NewPipeline 创建一个新的 Pipeline 实例，支持任意数据类型
func NewPipeline[T any]() *Pipeline[T] {
	return &Pipeline[T]{}
}

// AddStep 添加一个新的处理步骤到管道中
func (p *Pipeline[T]) AddStep(f func(chan T) chan T) {
	p.steps = append(p.steps, f)
}

// Compute 对输入通道应用管道中的所有处理步骤，并返回输出通道
func (p *Pipeline[T]) Compute(input chan T) chan T {
	ch := input
	for _, step := range p.steps {
		ch = step(ch) // 对通道应用每一个步骤
	}
	return ch // 返回最终的输出通道
}

// FromArray 将输入切片 `a` 中的每个元素应用函数 `f`，并返回一个包含结果的通道。
// 参数:
//   - f: 一个函数，接受切片中的元素 `T` 类型，并返回 `U` 类型的结果。
//   - a: 一个 `T` 类型的切片，将对其每个元素应用函数 `f`。
//
// 返回:
//   - 一个 `U` 类型的通道，通道中的值是对切片 `a` 中的每个元素应用函数 `f` 的结果。
func FromArray[T any, U any](f func(x T) U, a []T) chan U {

	ch := make(chan U, BufferSize)

	go func() {
		defer close(ch)
		for _, v := range a {
			ch <- f(v)
		}

	}()
	return ch

}

// FromMap 将一个映射转换为一个通道，每个通道中的元素是映射中的键值对。
// 参数:
//   - m: 一个映射，键类型为 K，值类型为 V。
//
// 返回:
//   - 一个通道，每个通道中的元素是包含一个键值对的切片，类型为 []array.Pair[K, V]。
//
// 函数功能:
//   - 遍历映射 m，将每个键值对包装在一个切片中，然后将这些切片逐个发送到通道 ch 中。
//   - 每个通道中的元素都是一个包含单个键值对的切片。
//   - 当所有键值对都被发送到通道后，关闭通道。
func FromMap[K comparable, V any](m map[K]V) chan array.Pair[K, V] {

	ch := make(chan array.Pair[K, V], BufferSize)

	go func() {
		defer close(ch)
		for k, v := range m {
			ch <- array.Pair[K, V]{
				First:  k,
				Second: v,
			}
		}
	}()

	return ch
}

// FromCsv 从指定的 CSV 文件路径读取数据，并将其以切片的形式发送到通道。
// 参数:
//   - path: CSV 文件的路径，字符串类型。
//
// 返回:
//   - 一个通道，通道中的值是读取的 CSV 文件中的每一行数据，每一行数据是一个字符串切片（[]string）。
func FromCsv(path string) chan []string {

	ch := make(chan []string, BufferSize)

	go func() {
		err := array.ReadFromCsvSliceChannel(path, ch)
		if err != nil {
			log.Println(err)
		}
	}()
	return ch
}

func FromExcel(path string, sheet string) chan []string {

	ch := make(chan []string, BufferSize)

	go func() {
		err := array.ReadFromExcelSliceChannel(path, sheet, ch)
		if err != nil {
			log.Println(err)
		}
	}()
	return ch
}

// ToCsv 将通道中的数据写入指定的 CSV 文件。
// 参数:
//   - path: CSV 文件的路径，字符串类型。
//   - ch: 一个通道，通道中的每个值是一个字符串切片（[]string），表示 CSV 文件中的一行数据。
//
// 函数功能:
//   - 从通道中读取数据，并将数据写入指定的 CSV 文件。
//   - 使用一个 goroutine 执行写入操作，并通过 stop 通道同步写入完成。
func ToCsv(path string, ch chan []string) {

	stop := make(chan struct{})

	go func() {
		err := array.WriteToCSVStringSliceChannel(ch, stop, path)
		if err != nil {
			log.Println(err)
		}
	}()

	<-stop

}

func ToExcel(path string, sheet string, ch chan []string) {

	stop := make(chan struct{})

	go func() {
		err := array.WriteToExcelStringSliceChannel(ch, stop, path, sheet)
		if err != nil {
			log.Println(err)
		}
	}()

	<-stop

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
func Map[T any, U any](f func(x T) U, ch chan T) chan U {

	ch_ := make(chan U, BufferSize)

	go func() {
		defer close(ch_)
		for v := range ch {
			ch_ <- f(v)
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
func Walk[T any, U any](f func(x T) U, ch chan T) {

	for v := range ch {
		f(v)
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
func Filter[T any](f func(x T) bool, ch chan T) chan T {

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
func Reduce[T, U any](f func(x U, y T) U, init U, ch chan T) U {

	for v := range ch {
		init = f(init, v)
	}
	return init

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
func Scanl[T, U any](f func(x U, y T) U, init U, ch chan T) chan U {

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
func Zip[T any, U any, V any](f func(x T, y U) V, ch1 chan T, ch2 chan U) chan V {
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
func Partition[T any](f func(x T) bool, ch chan T) (chan T, chan T) {
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
func Find[T any](f func(x T) bool, ch chan T) T {

	for v := range ch {
		if f(v) {
			return v
		}
	}
	var result T
	return result
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
func TakeWhile[T any](f func(x T) bool, ch chan T) chan T {
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
func DropWhile[T any](f func(x T) bool, ch chan T) chan T {
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
func Merge[T any](chs ...chan T) chan T {

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
	out := make(chan []T, 3)

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

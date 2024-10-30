package iter

import (
	"sync"

	"github.com/frankill/gotools"
	"github.com/frankill/gotools/array"
)

func More[T, U, R any](f func(x T, y U) R) func(ch1 chan T, ch2 chan U) chan R {

	return func(ch1 chan T, ch2 chan U) chan R {
		ch_ := make(chan R, bufferSize) // 使用合理的缓冲区大小

		go func() {

			defer close(ch_)

			for {

				t, ok1 := <-ch1
				u, ok2 := <-ch2

				if !ok1 && !ok2 {
					return
				}

				if !ok1 || !ok2 {
					continue
				}

				ch_ <- f(t, u)
			}
		}()

		return ch_
	}
}

// Subtract 返回两个通道的差集
// 参数:
//   - f: 一个函数，接受两个类型为 T 的值，返回一个布尔值，表示比较大小 0 相等，-1 小于，1 大于。
//   - ch1: 一个通道，用于接收数据。必须排序
//   - ch2: 一个通道，用于接收数据。必须排序
//
// 返回:
//   - 一个通道，用于接收排序后的数据。
func Subtract[T any](f func(x, y T) int) func(ch1 chan T, ch2 chan T) chan T {
	return func(ch1 chan T, ch2 chan T) chan T {
		ch_ := make(chan T, bufferSize)

		go func() {
			defer close(ch_)

			var v1, v2 T
			var ok1, ok2 bool

			// Get first value from both channels
			v1, ok1 = <-ch1
			v2, ok2 = <-ch2

			for ok1 {
				// Compare values from ch1 and ch2
				switch {
				case !ok2: // ch2 is exhausted
					ch_ <- v1
					v1, ok1 = <-ch1
				case f(v1, v2) < 0: // v1 < v2
					ch_ <- v1
					v1, ok1 = <-ch1
				case f(v1, v2) > 0: // v1 > v2
					v2, ok2 = <-ch2
				default: // v1 == v2
					v1, ok1 = <-ch1
				}
			}
		}()

		return ch_
	}
}

// Intersect 返回两个通道的交集
// 参数:
//   - f: 一个函数，接受两个类型为 T 的值，返回一个布尔值，表示比较大小 0 相等，-1 小于，1 大于。
//   - ch1: 一个通道，用于接收数据。必须排序
//   - ch2: 一个通道，用于接收数据。必须排序
//
// 返回:
//   - 一个通道，用于接收排序后的数据。
func Intersect[T any](f func(x, y T) int) func(ch1, ch2 chan T) chan T {
	return func(ch1, ch2 chan T) chan T {
		ch_ := make(chan T, bufferSize)

		go func() {
			defer close(ch_)

			var v1, v2 T
			var ok1, ok2 bool

			// Get first value from both channels
			v1, ok1 = <-ch1
			v2, ok2 = <-ch2

			for ok1 {
				// Compare values from ch1 and ch2
				switch {
				case !ok2: // ch2 is exhausted
					return
				case f(v1, v2) < 0: // v1 < v2
					v1, ok1 = <-ch1
				case f(v1, v2) > 0: // v1 > v2
					v2, ok2 = <-ch2
				default: // v1 == v2
					ch_ <- v1
					v1, ok1 = <-ch1
				}
			}
		}()

		return ch_
	}
}

// InterS 返回两个通道的交集
// 参数:
//   - ch1: 一个通道，通道中的值是类型为 T 的数据。
//   - ch2: 一个通道，通道中的值是类型为 T 的数据。
//
// 返回:
//   - 一个通道，通道中的值是类型为 T 的数据，表示两个通道的交集。
//
// 注意:
// 由于需要对收集第二个通道的数据，因此可以将较少数据的通道传递给第二个通道。
// 如果第二个通道数据很多，要考虑内存占用问题。
func InterS[T gotools.Comparable](ch1 chan T, ch2 chan T) chan T {

	ch_ := make(chan T, bufferSize)
	m := array.ToMap(Collect(ch2))

	var wg sync.WaitGroup
	wg.Add(parallerNum)

	go func() {
		defer close(ch_)
		defer wg.Wait()
	}()

	for num := 0; num < parallerNum; num++ {

		go func() {

			defer wg.Done()

			for v := range ch1 {

				if _, ok := m[v]; ok {
					ch_ <- v
				}

			}

		}()
	}

	return ch_

}

// func toMapInt[T gotools.Comparable](arr []T) map[T]int {

// 	result := make(map[T]int, len(arr))

// 	for i := range arr {
// 		result[arr[i]] = 1
// 	}
// 	return result
// }

// SubS 返回两个通道的差集
// 参数:
//   - ch1: 一个通道，通道中的值是类型为 T 的数据。
//   - ch2: 一个通道，通道中的值是类型为 T 的数据。
//
// 返回:
//   - 一个通道，通道中的值是类型为 T 的数据，表示两个通道的差集。
//
// 注意:
// 由于需要对收集第二个通道的数据，因此可以将较少数据的通道传递给第二个通道。
// 如果第二个通道数据很多，要考虑内存占用问题。
func SubS[T gotools.Comparable](ch1 chan T, ch2 chan T) chan T {

	ch_ := make(chan T, bufferSize)
	m := array.ToMap(Collect(ch2))

	var wg sync.WaitGroup

	wg.Add(parallerNum)

	go func() {
		defer close(ch_)
		defer wg.Wait()
	}()

	for num := 0; num < parallerNum; num++ {

		go func() {

			defer wg.Done()

			for v := range ch1 {

				if _, ok := m[v]; !ok {
					ch_ <- v
				}

			}

		}()
	}

	return ch_
}

// InnerJoin 按照某个函数进行连接
// 参数:
//   - f: 一个函数，接受两个类型为 T 和 U 的值，返回一个类型为 R 的值。
//   - f1: 一个函数，用于比较大小， -1 表示小于， 0 表示相等， 1 表示大于。,通道数据必须升序，当通道数据降序，函数结果数字要相反
//   - ch1: 一个通道，用于接收第一个数据
//   - ch2: 一个通道，用于接收第二个数据
//
// 返回:
//   - 一个通道，用于接收连接后的数据。
func InnerJoin[T any, U any, R any](f func(x T, y U) R, f1 func(x T, y U) int) func(ch1 chan T, ch2 chan U) chan R {

	return func(ch1 chan T, ch2 chan U) chan R {
		ch_ := make(chan R, bufferSize) // 使用合理的缓冲区大小

		go func() {
			defer close(ch_)

			var t T
			var u U
			var tok, uok, ok1, ok2 = true, true, true, true

			us := make([]U, 0)

			for {

				if tok || !ok2 {
					t, ok1 = <-ch1
					tok = ok1
				}

				if uok || !ok1 {
					u, ok2 = <-ch2
					uok = ok2
				}

				if !ok1 && !ok2 {
					break
				}

				if tok {
					if len(us) > 0 {
						if f1(t, us[0]) == 0 {
							for _, v := range us {
								ch_ <- f(t, v)
							}

							uok = false
							tok = true

							continue
						} else {
							us = us[:0]
						}

					}

				}

				if !ok2 || !ok1 {
					continue
				}

				switch f1(t, u) {
				case -1:
					uok = false
					tok = true

				case 1:
					tok = false
					uok = true

				case 0:
					tok = false
					uok = true
					ch_ <- f(t, u)
					us = append(us, u)

				}

			}

		}()

		return ch_
	}
}

// LeftJoin 按照某个函数进行连接
// 参数:
//   - f: 一个函数，接受两个类型为 T 和 U 的值，返回一个类型为 R 的值。
//   - f1: 一个函数，用于比较大小， -1 表示小于， 0 表示相等， 1 表示大于,通道数据必须升序，当通道数据降序，函数结果数字要相反
//   - ch1: 一个通道，用于接收第一个数据
//   - ch2: 一个通道，用于接收第二个数据
//
// 返回:
//   - 一个通道，用于接收连接后的数据。
func LeftJoin[T any, U any, R any](f func(x T, y U) R, f1 func(x T, y U) int) func(ch1 chan T, ch2 chan U) chan R {

	return func(ch1 chan T, ch2 chan U) chan R {
		ch_ := make(chan R, bufferSize) // 使用合理的缓冲区大小

		go func() {
			defer close(ch_)

			var t T
			var u U
			var um U
			var tok, uok, ok1, ok2 = true, true, true, true

			us := make([]U, 0)

			for {

				if tok || !ok2 {
					t, ok1 = <-ch1
					tok = ok1
				}

				if uok || !ok1 {
					u, ok2 = <-ch2
					uok = ok2
					if !ok2 && !tok && ok1 {
						tok = true
					}
				}

				if !ok1 && !ok2 {
					break
				}

				if tok {
					if len(us) > 0 {
						if f1(t, us[0]) == 0 {
							for _, v := range us {
								ch_ <- f(t, v)
							}

							uok = false
							tok = true

							continue
						} else {
							us = us[:0]
						}

					}

					if !ok2 {
						ch_ <- f(t, u)
					}

				}

				if !ok2 || !ok1 {
					continue
				}

				switch f1(t, u) {
				case -1:
					uok = false
					tok = true

					if len(us) == 0 {
						ch_ <- f(t, um)
					}

				case 1:
					tok = false
					uok = true

				case 0:
					tok = false
					uok = true
					ch_ <- f(t, u)
					us = append(us, u)

				}

			}

		}()

		return ch_
	}
}

// Merge 合并排序多个通道，按照func(x, y T) bool进行排序
// 参数:
//   - f: 一个函数，接受两个类型为 T 和 U 的值，返回一个布尔值，表示是否满足排序条件。
//     当 `fun(x, y)` 返回 `true`，则在排序时 `x` 应位于 `y` 之前。
//   - cs: 一个包含多个通道的切片。
//
// 返回:
//   - 一个通道，用于接收排序后的数据。
func Merge[T any](f func(x, y T) bool) func(cs ...chan T) chan T {

	return func(cs ...chan T) chan T {

		mins := make([]T, len(cs))
		sortedCh := make(chan T, bufferSize)

		go func() {

			defer close(sortedCh)

			for i, ch := range cs {
				item, ok := <-ch
				if ok {
					mins[i] = item
				}
			}
			for {

				minIndex := -1

				for i := range mins {

					if minIndex == -1 {
						minIndex = i
					}

					if !f(mins[minIndex], mins[i]) {
						minIndex = i
					}
				}

				if minIndex == -1 {

					break
				}

				sortedCh <- mins[minIndex]

				item, ok := <-cs[minIndex]

				if ok {
					mins[minIndex] = item
				} else {

					mins = append(mins[:minIndex], mins[minIndex+1:]...)
					cs = append(cs[:minIndex], cs[minIndex+1:]...)

				}
			}
		}()

		return sortedCh

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
		ch := make(chan V, bufferSize)

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

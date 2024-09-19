package slider

import "fmt"

// Slide 滑动窗口计算
// 参数:
//   - f: 一个函数，接受类型为 []T 的输入，返回类型为 U 的结果。
//   - data: 一个类型为 T 的切片。
//   - before: 一个整数，表示滑动窗口的前面部分的长度。
//   - after: 一个整数，表示滑动窗口的后面部分的长度。
//   - defaultValue: 一个类型为 T 的值，表示滑动窗口中的默认值。
//
// 返回:
//   - 一个类型为 U 的切片，表示滑动窗口计算后的结果。
func Slide[S ~[]T, U, T any](f func(x []T) U, data S, before int, after int, defaultValue T) []U {
	l := len(data)
	wl := before + after + 1
	if wl > l {
		return []U{} // 如果窗口大小大于数据长度，返回空结果
	}
	windows := make([]T, wl)
	index := wl / 2
	result := make([]U, l)

	for i := 0; i < l; i++ {

		windows[index] = data[i]

		for j := 1; j <= before; j++ {
			if i-j >= 0 {
				windows[index-j] = data[i-j]
			} else {
				windows[index-j] = defaultValue
			}

		}

		for j := 1; j <= after; j++ {
			if i+j < l {
				windows[index+j] = data[i+j]
			} else {
				windows[index+j] = defaultValue
			}
		}
		fmt.Println(windows)
		result[i] = f(windows)

	}

	return result
}

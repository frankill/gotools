package slider

import (
	"strings"
	"time"

	"github.com/frankill/gotools"
	"github.com/frankill/gotools/array"
)

// Sum 对数据进行滑动窗口求和。
// 参数:
//   - before: 一个整数，表示滑动窗口的前面部分的长度。
//   - after: 一个整数，表示滑动窗口的后面部分的长度。
//   - defaultValue: 一个类型为 T 的值，表示滑动窗口中的默认值。
//   - data: 一个类型为 T 的切片。
//
// 返回值:
//   - 一个类型为 T 的切片，表示滑动窗口计算后的结果。
func Sum[S ~[]T, T gotools.Number](before int, after int, defaultValue T, data S) []T {

	return Slide(func(x []T) T {
		return array.Sum(x)
	}, before, after, defaultValue, data)

}

// Max 对数据进行滑动窗口求最大值。
// 参数:
//   - before: 一个整数，表示滑动窗口的前面部分的长度。
//   - after: 一个整数，表示滑动窗口的后面部分的长度。
//   - defaultValue: 一个类型为 T 的值，表示滑动窗口中的默认值。
//   - data: 一个类型为 T 的切片。
//
// 返回值:
//   - 一个类型为 T 的切片，表示滑动窗口计算后的结果。
func Max[S ~[]T, T gotools.Number](before int, after int, defaultValue T, data S) []T {

	return Slide(func(x []T) T {
		return array.Max(x)
	}, before, after, defaultValue, data)

}

// Min 对数据进行滑动窗口求最小值。
// 参数:
//   - before: 一个整数，表示滑动窗口的前面部分的长度。
//   - after: 一个整数，表示滑动窗口的后面部分的长度。
//   - defaultValue: 一个类型为 T 的值，表示滑动窗口中的默认值。
//   - data: 一个类型为 T 的切片。
//
// 返回值:
//   - 一个类型为 T 的切片，表示滑动窗口计算后的结果。
func Min[S ~[]T, T gotools.Number](before int, after int, defaultValue T, data S) []T {

	return Slide(func(x []T) T {
		return array.Min(x)
	}, before, after, defaultValue, data)

}

// Mean 对数据进行滑动窗口求平均值。
// 参数:
//   - before: 一个整数，表示滑动窗口的前面部分的长度。
//   - after: 一个整数，表示滑动窗口的后面部分的长度。
//   - defaultValue: 一个类型为 T 的值，表示滑动窗口中的默认值。
//   - data: 一个类型为 T 的切片。
//
// 返回值:
//   - 一个类型为 T 的切片，表示滑动窗口计算后的结果。
func Mean[S ~[]T, T gotools.Number](before int, after int, defaultValue T, data S) []T {

	return Slide(func(x []T) T {
		return array.Mean(x)
	}, before, after, defaultValue, data)

}

//	Paste 拼接滑动窗口字符串
//
// 参数:
//   - before: 一个整数，表示滑动窗口的前面部分的长度。
//   - after: 一个整数，表示滑动窗口的后面部分的长度。
//   - defaultValue: 一个类型为 string 的值，表示滑动窗口中的默认值。
//   - data: 一个类型为 string 的切片。
//
// 返回值:
//   - 一个类型为 string 的切片，表示滑动窗口计算后的结果。
func Paste(before int, after int, defaultValue string, data []string) []string {

	return Slide(func(x []string) string {
		return strings.Join(x, "")
	}, before, after, defaultValue, data)

}

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
func Slide[S ~[]T, U, T any](f func(x []T) U, before int, after int, defaultValue T, data S) []U {
	l := len(data)
	wl := before + after + 1
	if before+1 > l || after+1 > l {
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

		result[i] = f(windows)

	}

	return result
}

// Pslide 多输入切片滑动窗口计算
// 参数:
//   - f: 一个函数，接受类型为 []T 的输入，返回类型为 U 的结果。
//   - data: 一个类型为 T 的切片。
//   - before: 一个整数，表示滑动窗口的前面部分的长度。
//   - after: 一个整数，表示滑动窗口的后面部分的长度。
//   - defaultValue: 一个类型为 T 的值，表示滑动窗口中的默认值。
//
// 返回:
//   - 一个类型为 U 的切片，表示滑动窗口计算后的结果。
//
// 注意:
// 多个切片会分别计算 并最终合并相同索引的元素到一个切片元素中
func Pslide[S ~[]T, U, T any](f func(x []T) U, before int, after int, defaultValue T, data ...S) [][]U {

	if len(data) == 0 {
		return [][]U{}
	}

	res := make([][]U, len(data))

	ch_ := make([]chan []U, len(data))

	for i := 0; i < len(data); i++ {

		ch_[i] = make(chan []U, 1)

		go func(i int) {
			ch_[i] <- Slide(f, before, after, defaultValue, data[i])
		}(i)

	}

	for i := 0; i < len(data); i++ {

		res[i] = <-ch_[i]
	}

	return array.Zip(res...)

}

// SlideIndex 基于时间切片的滑动窗口计算
// 参数:
//   - f: 一个函数，接受类型为 []T 的输入，返回类型为 U 的结果。
//   - data: 一个类型为 T 的切片。
//   - before: 一个整数，表示滑动窗口的前面部分的天数。
//   - after: 一个整数，表示滑动窗口的后面部分的天数。
//   - index : 一个类型是time的切片，用于计算滑动窗口的长度
//
// 返回:
//   - 一个类型为 U 的切片，表示滑动窗口计算后的结果。
func SlideIndex[S ~[]T, T, U any](f func(x []T) U, before, after int, index []time.Time, data S) []U {
	l := len(data)

	result := make([]U, l)
	windows := make([]T, 0)
	for i := 0; i < l; i++ {
		windows := windows[:0]
		for j := 1; i-j >= 0; j++ {
			sub := index[i].Sub(index[i-j])
			if sub.Hours() <= float64(before*24) && sub.Hours() >= 0.0 {
				windows = append(windows, data[i-j])
			}

		}

		windows = append(windows, data[i])

		for j := 1; i+j < l; j++ {
			sub := index[i+j].Sub(index[i])
			if sub.Hours() <= float64(after*24) && sub.Hours() >= 0.0 {
				windows = append(windows, data[i+j])
			}
		}

		result[i] = f(windows)

	}

	return result
}

// PslideIndex 多输入切片滑动窗口计算
// 参数:
//   - f: 一个函数，接受类型为 []T 的输入，返回类型为 U 的结果。
//   - data: 一个类型为 T 的切片。
//   - before: 一个整数，表示滑动窗口的前面部分的长度。
//   - after: 一个整数，表示滑动窗口的后面部分的长度。
//
// 返回:
//   - 一个类型为 U 的切片，表示滑动窗口计算后的结果。
//
// 注意:
// 多个切片会分别计算 并最终合并相同索引的元素到一个切片元素中
func PslideIndex[S ~[]T, U, T any](f func(x []T) U, before int, after int, index []time.Time, data ...S) [][]U {

	if len(data) == 0 {
		return [][]U{}
	}

	res := make([][]U, len(data))

	ch_ := make([]chan []U, len(data))

	for i := 0; i < len(data); i++ {

		ch_[i] = make(chan []U, 1)

		go func(i int) {
			ch_[i] <- SlideIndex(f, before, after, index, data[i])
		}(i)

	}

	for i := 0; i < len(data); i++ {

		res[i] = <-ch_[i]
	}

	return array.Zip(res...)

}

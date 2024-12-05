package slider

import (
	"strings"
	"time"

	"github.com/frankill/gotools"
	"github.com/frankill/gotools/array"
	"github.com/frankill/gotools/op"
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

	return Slider(func(x ...T) T {
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

	return Slider(func(x ...T) T {
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

	return Slider(func(x ...T) T {
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

	return Slider(func(x ...T) T {
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

	return Slider(func(x ...string) string {
		return strings.Join(x, "")
	}, before, after, defaultValue, data)

}

// Slider 滑动窗口计算
// 参数:
//   - f: 一个函数，接受类型为 ...T 的输入，返回类型为 U 的结果。
//   - data: 一个类型为 T 的切片。
//   - before: 一个整数，表示滑动窗口的前面部分的长度。
//   - after: 一个整数，表示滑动窗口的后面部分的长度。
//   - defaultValue: 一个类型为 T 的值，表示滑动窗口中的默认值。
//
// 返回:
//   - 一个类型为 U 的切片，表示滑动窗口计算后的结果。
//
// 注意:
//
//   - 如果函数f 返回x 需要使用 array.Copy(x) 复制, 否则结果可能不正确
//
//   - a := array.RandomSample(array.Seq(1, 10, 1), 20, true)
//
//   - slider.Slide(func(x []int) []int { return array.Copy(x) }, 1, 1, 0, a)
func Slider[S ~[]T, U, T any](f func(x ...T) U, before int, after int, defaultValue T, data S) []U {
	l := len(data)
	cap := before + after + 1

	windows := make([]T, 0, cap)
	result := make([]U, 0, l)
	num := 0

	for i := 0; i < l; i++ {

	llp:
		if num < before {
			windows = append(windows, defaultValue)
			num++
			goto llp
		}

		windows = append(windows, data[i])
		num++

		if num == cap {
			result = append(result, f(windows...))
			windows = shift(windows)
			num--
		}

	}

	for i := num; i < cap-1; i++ {
		windows = append(windows, defaultValue)
	}

	r := min(after, l)

	if num > 0 {
		for i := 0; i < r; i++ {
			windows = append(windows, defaultValue)
			result = append(result, f(windows...))
			windows = shift(windows)
		}
	}
	return result
}

func shift[T any](x []T) []T {

	copy(x, x[1:])

	return x[:len(x)-1]

}

// Pslide 多输入切片滑动窗口计算
// 参数:
//   - f: 一个函数，接受类型为 ...T 的输入，返回类型为 U 的结果。
//   - data: 一个类型为 T 的切片。
//   - before: 一个整数，表示滑动窗口的前面部分的长度。
//   - after: 一个整数，表示滑动窗口的后面部分的长度。
//   - defaultValue: 一个类型为 T 的值，表示滑动窗口中的默认值。
//
// 返回:
//   - 一个类型为 U 的切片，表示滑动窗口计算后的结果。
//
// 注意:
//
//   - 多个切片会分别计算 并最终合并相同索引的元素到一个切片元素中
//
//   - 如果函数f 返回x 需要使用 array.Copy(x) 复制, 否则结果可能不正确
func Pslide[S ~[]T, U, T any](f func(x ...T) U, before int, after int, defaultValue T, data ...S) [][]U {

	if len(data) == 0 {
		return [][]U{}
	}

	res := make([][]U, len(data))

	ch_ := make([]chan []U, len(data))

	for i := 0; i < len(data); i++ {

		ch_[i] = make(chan []U, 1)

		go func(i int) {
			ch_[i] <- Slider(f, before, after, defaultValue, data[i])
		}(i)

	}

	for i := 0; i < len(data); i++ {

		res[i] = <-ch_[i]
	}

	return op.Zip(res...)

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
//   - 如果函数f 返回x 需要使用 array.Copy(x) 复制, 否则结果可能不正确
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
//   - 多个切片会分别计算 并最终合并相同索引的元素到一个切片元素中
//   - 如果函数f 返回x 需要使用 array.Copy(x) 复制, 否则结果可能不正确
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

	return op.Zip(res...)

}

package group

import (
	"github.com/frankill/gotools"
	"github.com/frankill/gotools/array"
	"github.com/frankill/gotools/maps"
)

func WindowFun[B ~[]U, C ~[]S, U gotools.Comparable, S gotools.Ordered](by B, order C) func(func(by B, data C) []S) []S {

	return func(fn func(by B, data C) []S) []S {
		return fn(by, order)
	}

}

// SequenceMatch 函数用于生成一个根据指定序列和排序规则来检查元素匹配情况的函数。
// 参数：
//   - by: 实际的排序规则依据序列，类型为 B。
//   - data: 需要进行匹配操作的数据序列，类型为 D。
//   - order: 指定的排序顺序，类型为 O。
//
// 返回值：
//   - 该函数返回一个高阶函数，该高阶函数接受一个与 D 类型相同的序列作为输入，
//     进一步返回一个函数，该函数接受一个整数切片作为索引，最终产出一个 map，
//     其键为 T 类型，值为布尔类型，表示 data 中相应位置的元素是否在输入序列中按照给定排序规则出现。
//
// 示例 :
//   - a, b, c := []int{1, 1, 1, 2, 2, 2}, []int{1, 2, 3, 4, 5, 6}, []int{1, 1, 1, 1, 1, 1}
//   - SequenceMatch(a, b, c)([]int{1, 2, 3})([]int{1, 2}) = map[1:true 2:false]
func SequenceMatch[B ~[]T, D ~[]U, O ~[]S, T gotools.Comparable, U gotools.Comparable, S gotools.Ordered](by B, data D, order O) func([]U) func([]int) map[T]bool {

	group := ByOrder(by, data, order)

	return func(eventID []U) func([]int) map[T]bool {

		group = maps.ApplyValue(func(k T, v []U) []U {
			return array.Filter(func(x U) bool { return array.Has(eventID, x) }, v)
		}, group)

		return func(modeID []int) map[T]bool {

			modeID = array.Map(func(x int) int { return x - 1 }, modeID)
			eID := array.Map(func(x int) U { return eventID[x] }, modeID)
			return maps.ApplyValue(func(k T, v []U) bool {
				ok, _ := array.HasSequence(v, eID)
				return ok
			}, group)

		}

	}

}

// SequenceMatch 函数用于生成一个根据指定序列和排序规则返回匹配次数。
// 参数说明：
//  - by: 实际的排序规则依据序列，类型为 B。
//  - data: 需要进行匹配操作的数据序列，类型为 D。
//  - order: 指定的排序顺序，类型为 O。

// 返回值：
//  - 该函数返回一个高阶函数，该高阶函数接受一个与 D 类型相同的序列作为输入，
//    进一步返回一个函数，该函数接受一个整数切片作为索引，最终产出一个 map，
//    其键为 T 类型，值为int类型，表示 data 中相应位置的元素是否在输入序列中按照给定排序规则出现次数。

// 示例 :
//   - a, b, c := []int{1, 1, 1}, []string{"a", "a", "b"}, []int{}
//   - SequenceCount(a, b, c)([]string{"a"})([]int{1}) = map[1:2]
func SequenceCount[B ~[]T, D ~[]U, O ~[]S, T gotools.Comparable, U gotools.Comparable, S gotools.Ordered](by B, data D, order O) func([]U) func([]int) map[T]int {

	group := ByOrder(by, data, order)

	return func(eventID []U) func([]int) map[T]int {

		group = maps.ApplyValue(func(k T, v []U) []U {
			return array.Filter(func(x U) bool { return array.Has(eventID, x) }, v)
		}, group)

		return func(modeID []int) map[T]int {

			modeID = array.Map(func(x int) int { return x - 1 }, modeID)
			eID := array.Map(func(x int) U { return eventID[x] }, modeID)
			return maps.ApplyValue(func(k T, v []U) int {
				return array.ArrSequenceCount(v, eID)
			}, group)

		}

	}

}

// WindowFunnel 创建一个基于窗口的漏斗分析函数，该函数接受一组事件ID，
// 并根据不同的模式返回每个组内满足条件的事件序列的最大计数。
//
// Parameters:
//
//   - by: 用于分组的键值切片。通常这些键值代表了数据的不同维度，如用户ID或产品类别。
//   - data: 包含实际数据的切片，每个元素对应于一个事件或记录。
//   - order: 用于排序的有序切片，确保数据按照时间或其它有意义的顺序排列。
//
// Returns:
//
//   - 一个函数，该函数接收事件ID切片作为参数，并进一步返回一个函数。
//     这个进一步返回的函数接受模式字符串并返回一个映射，其中键是分组键，
//     值是在该组内满足给定模式的事件序列的最大计数。
//
// Modes:
//
//   - "strict_order": 严格顺序模式，计算事件序列在数据中严格按照事件ID顺序出现的最大次数，意外的事件会中断。
//   - "strict_dedup": 严格去重模式，计算事件序列在数据中出现的最大次数，重复的事件会中断。
//   - "strict_increase": 严格递增模式，计算事件序列在数据中按顺序出现的最大次数，重复，意外事件不会中断。
func WindowFunnel[B ~[]T, D ~[]U, O ~[]S, T gotools.Comparable, U gotools.Comparable, S gotools.Ordered](by B, data D, order O) func([]U) func(mode string) map[T]int {

	group := ByOrder(by, data, order)

	return func(eventID []U) func(mode string) map[T]int {

		exist := array.ToMap(eventID)

		return func(mode string) map[T]int {

			return maps.ApplyValue(func(k T, v []U) int {

				var num int

				switch mode {
				case "strict_order":
					num = array.HasOrderMaxCount(v, eventID, exist)
				case "strict_dedup":
					num = array.HasDupMaxCount(v, eventID)
				case "strict_increase":
					num = array.HasIncreaseMaxCount(v, eventID)
				}

				return num

			}, group)

		}

	}
}

type Retentionfun[T any] func(x []T) bool

// Retention 函数根据提供的条件函数列表，生成一个新的闭包函数，该闭包函数能够处理分类标识和数据集，
// 并返回一个映射，其中键为分类标识，值为一个布尔值切片，表示每个类别在不同条件下的筛选结果。
//
// 参数:
//   - fun: 一系列的条件函数，每个函数接收一个数据切片并返回一个布尔值。
//     这些函数将被用来决定数据项是否满足特定条件。
//
// 返回:
//   - 一个闭包函数，它需要两个参数：
//     by: 一个分类标识的切片，用于区分数据集中的不同组。
//     data: 数据集切片，与分类标识一一对应，用于条件检查。
//   - 闭包函数返回一个映射，除第一个条件外，条件成对应用：如果第一个和第二个为真，则第二个结果为真，如果第一个和第三个为真，则第三个结果为真，等等
func Retention[B ~[]U, C ~[]S, U gotools.Comparable, S any](by B, data C) func(fun ...Retentionfun[S]) map[U][]bool {

	group := By(by, data)

	return func(fun ...Retentionfun[S]) map[U][]bool {

		firstFun := fun[0]
		fun = fun[1:]

		value := make(map[U][]bool, len(group))

		for k, v := range group {

			tmp := make([]bool, 0, len(fun)-1)
			ts := firstFun(v)
			tmp = append(tmp, ts)

			for _, f := range fun {
				tmp = append(tmp, ts && f(v))
			}

			value[k] = tmp

		}

		return value
	}

}

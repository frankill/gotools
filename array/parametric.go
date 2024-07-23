package array

import "cmp"

/*
SequenceMatch 函数用于生成一个根据指定序列和排序规则来检查元素匹配情况的函数。

参数说明：
- B: 类型参数，表示用于确定排序规则的切片类型，元素类型为 T，要求可比较。
- D: 类型参数，表示数据集的切片类型，元素类型为 U，要求可比较。
- O: 类型参数，表示排序顺序的切片类型，元素类型为 S，要求为有序类型。
- by: 实际的排序规则依据序列，类型为 B。
- data: 需要进行匹配操作的数据序列，类型为 D。
- order: 指定的排序顺序，类型为 O。

返回值：
该函数返回一个高阶函数，该高阶函数接受一个与 D 类型相同的序列作为输入，
进一步返回一个函数，该函数接受一个整数切片作为索引，最终产出一个 map，
其键为 T 类型，值为布尔类型，表示 data 中相应位置的元素是否在输入序列中按照给定排序规则出现。

示例 :
a, b, c := []int{1, 1, 1, 2, 2, 2}, []int{1, 2, 3, 4, 5, 6}, []int{1, 1, 1, 1, 1, 1}
SequenceMatch(a, b, c)([]int{1, 2, 3})([]int{1, 2}) = map[1:true 2:false]
*/
func SequenceMatch[B ~[]T, D ~[]U, O ~[]S, T comparable, U comparable, S cmp.Ordered](by B, data D, order O) func([]U) func([]int) map[T]bool {

	group := GroupByOrder(by, data, order)

	return func(eventID []U) func([]int) map[T]bool {

		group = MapApplyValue(func(k T, v []U) []U {
			return ArrayFilter(func(x ...U) bool { return ArrayHas(eventID, x[0]) }, v)
		}, group)

		return func(modeID []int) map[T]bool {

			modeID = ArrayMap(func(x ...int) int { return x[0] - 1 }, modeID)
			eID := ArrayMap(func(x ...int) U { return eventID[x[0]] }, modeID)
			return MapApplyValue(func(k T, v []U) bool {
				ok, _ := ArrayHasSequence(v, eID)
				return ok
			}, group)

		}

	}

}

/*
SequenceMatch 函数用于生成一个根据指定序列和排序规则返回匹配次数。

参数说明：
- B: 类型参数，表示用于确定排序规则的切片类型，元素类型为 T，要求可比较。
- D: 类型参数，表示数据集的切片类型，元素类型为 U，要求可比较。
- O: 类型参数，表示排序顺序的切片类型，元素类型为 S，要求为有序类型。
- by: 实际的排序规则依据序列，类型为 B。
- data: 需要进行匹配操作的数据序列，类型为 D。
- order: 指定的排序顺序，类型为 O。

返回值：
该函数返回一个高阶函数，该高阶函数接受一个与 D 类型相同的序列作为输入，
进一步返回一个函数，该函数接受一个整数切片作为索引，最终产出一个 map，
其键为 T 类型，值为int类型，表示 data 中相应位置的元素是否在输入序列中按照给定排序规则出现次数。

示例 :
a, b, c := []int{1, 1, 1}, []string{"a", "a", "b"}, []int{}
SequenceCount(a, b, c)([]string{"a"})([]int{1}) = map[1:2]
*/
func SequenceCount[B ~[]T, D ~[]U, O ~[]S, T comparable, U comparable, S cmp.Ordered](by B, data D, order O) func([]U) func([]int) map[T]int {

	group := GroupByOrder(by, data, order)

	return func(eventID []U) func([]int) map[T]int {

		group = MapApplyValue(func(k T, v []U) []U {
			return ArrayFilter(func(x ...U) bool { return ArrayHas(eventID, x[0]) }, v)
		}, group)

		return func(modeID []int) map[T]int {

			modeID = ArrayMap(func(x ...int) int { return x[0] - 1 }, modeID)
			eID := ArrayMap(func(x ...int) U { return eventID[x[0]] }, modeID)
			return MapApplyValue(func(k T, v []U) int {
				return ArraySequenceCount(v, eID)
			}, group)

		}

	}

}

// WindowFunnel 创建一个基于窗口的漏斗分析函数，该函数接受一组事件ID，
// 并根据不同的模式返回每个组内满足条件的事件序列的最大计数。
//
// Parameters:
//
//	by: 用于分组的键值切片。通常这些键值代表了数据的不同维度，如用户ID或产品类别。
//	data: 包含实际数据的切片，每个元素对应于一个事件或记录。
//	order: 用于排序的有序切片，确保数据按照时间或其它有意义的顺序排列。
//
// Returns:
//
//	一个函数，该函数接收事件ID切片作为参数，并进一步返回一个函数。
//	这个进一步返回的函数接受模式字符串并返回一个映射，其中键是分组键，
//	值是在该组内满足给定模式的事件序列的最大计数。
//
// Modes:
//
//	"strict_order": 严格顺序模式，计算事件序列在数据中严格按照事件ID顺序出现的最大次数，意外的事件会中断。
//	"strict_dedup": 严格去重模式，计算事件序列在数据中出现的最大次数，重复的事件会中断。
//	"strict_increase": 严格递增模式，计算事件序列在数据中按顺序出现的最大次数，重复，意外事件不会中断。
func WindowFunnel[B ~[]T, D ~[]U, O ~[]S, T comparable, U comparable, S cmp.Ordered](by B, data D, order O) func([]U) func(mode string) map[T]int {

	group := GroupByOrder(by, data, order)

	return func(eventID []U) func(mode string) map[T]int {

		exist := ArrayToMap(eventID)

		return func(mode string) map[T]int {

			return MapApplyValue(func(k T, v []U) int {

				var num int

				switch mode {
				case "strict_order":
					num = HasOrderMaxCount(v, eventID, exist)
				case "strict_dedup":
					num = HasDupMaxCount(v, eventID)
				case "strict_increase":
					num = HasIncreaseMaxCount(v, eventID)
				}

				return num

			}, group)

		}

	}
}

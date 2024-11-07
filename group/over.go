package group

import (
	"github.com/frankill/gotools"
	"github.com/frankill/gotools/array"
)

// RowNumber 根据指定的分组键(by)和排序序列(order)，为列表中的每一项元素分配一个行号。
// 分配规则确保相同分组内的元素根据order排序后获得连续递增的行号。
// 参数:
//
//   - by []U: 用于分组的依据，U类型需支持比较操作。
//   - order []S: 指定每个元素在组内的排序顺序，S类型需实现gotools.Ordered接口。
//
// 返回值:
//
//   - []int: 包含为列表中每个元素分配的行号的切片，反映元素在排序后的相对位置。
func RowNumber[B ~[]U, C ~[]S, U gotools.Comparable, S gotools.Ordered](by B, order C) []int {
	id := array.Seq(0, len(by), 1)

	group := ByOrder(by, id, order)

	sortid := make([]int, 0, len(by))
	numberid := make([]int, 0, len(by))

	for _, v := range group {
		sortid = append(sortid, v...)
		numberid = append(numberid, array.Seq(0, len(v), 1)...)
	}

	array.SortByL(gotools.ASCInt, numberid, sortid)

	return numberid
}

func RowNumberDesc[B ~[]U, C ~[]S, U gotools.Comparable, S gotools.Ordered](by B, order C) []int {
	id := array.Seq(0, len(by), 1)

	group := ByOrderDesc(by, id, order)

	sortid := make([]int, 0, len(by))
	numberid := make([]int, 0, len(by))

	for _, v := range group {
		sortid = append(sortid, v...)
		numberid = append(numberid, array.Seq(0, len(v), 1)...)
	}

	array.SortByL(gotools.ASCInt, numberid, sortid)

	return numberid
}

// MaxValue 返回一个 S 类型数组，其中包含输入数据中按指定字段分组后的最大值。
// 参数：
//   - data: C 类型的数组，其中 C 是 S 类型元素的切片，S 必须实现 gotools.Ordered 接口。
//   - by: B 类型的数组，表示用于分组的字段，B 是 U 类型元素的切片，U 可以比较。
//
// 返回：
//   - []S: 包含按原始顺序排列的最大值的有序数组。
//
// 泛型约束：
//   - B 和 C 分别是数据和分组依据的类型，它们需要满足相应的类型约束。
func MaxValue[B ~[]U, C ~[]S, U gotools.Comparable, S gotools.Ordered](by B, data C) []S {

	id := array.Seq(0, len(data), 1)
	group := PairBy(by, data, id)

	value := make([]S, 0, len(data))
	numberid := make([]int, 0, len(data))

	for _, v := range group {
		value = append(value, array.Repeat(array.Max(v.First), len(v.First))...)
		numberid = append(numberid, v.Second...)
	}

	array.SortByL(gotools.ASCInt, value, numberid)

	return value
}

// MinValue 返回一个 S 类型数组，其中包含输入数据中按指定字段分组后的最小值。
// 参数：
//   - data: C 类型的数组，其中 C 是 S 类型元素的切片，S 必须实现 gotools.Ordered 接口。
//   - by: B 类型的数组，表示用于分组的字段，B 是 U 类型元素的切片，U 可以比较。
//
// 返回：
//   - []S: 包含按原始顺序排列的最小值的有序数组。
//
// 泛型约束：
//   - B 和 C 分别是数据和分组依据的类型，它们需要满足相应的类型约束。
func MinValue[B ~[]U, C ~[]S, U gotools.Comparable, S gotools.Ordered](by B, data C) []S {

	id := array.Seq(0, len(data), 1)
	group := PairBy(by, data, id)

	value := make([]S, 0, len(data))
	numberid := make([]int, 0, len(data))

	for _, v := range group {
		value = append(value, array.Repeat(array.Min(v.First), len(v.First))...)
		numberid = append(numberid, v.Second...)
	}

	array.SortByL(gotools.ASCInt, value, numberid)

	return value
}

// FirstValue 根据提供的排序键by，对数据集data中的元素进行分组并提取每个分组的第一个元素值。
// 数据集data为任意类型C的切片，排序键by为可比较类型U的切片。U必须实现gotools.Comparable接口，
// 而S可以是任何类型。函数返回一个切片，包含每个分组的第一个元素值。
//
// 参数:
//   - data: 待处理的数据集，元素类型为S。
//   - by: 用于对data中元素进行分组的键，元素类型为U，需可比较。
//
// 返回值:
//   - []S: 一个切片，包含根据by分组后每个组的第一个元素值
func FirstValue[B ~[]U, C ~[]S, U gotools.Comparable, S any](by B, data C) []S {

	id := array.Seq(0, len(data), 1)
	group := PairBy(by, data, id)

	value := make([]S, 0, len(data))
	numberid := make([]int, 0, len(data))

	for _, v := range group {
		value = append(value, array.Repeat(v.First[0], len(v.First))...)
		numberid = append(numberid, v.Second...)
	}

	array.SortByL(gotools.ASCInt, value, numberid)

	return value
}

// LastValue 根据指定的键值（by）对数据（data）进行分组，
// 并返回每个分组的最后一个元素组成的切片。
//   - data: 输入的数据切片，类型为 C，其中 C 可以是任何切片类型。
//   - by: 分组依据的键值切片，类型为 B，其中 B 也是切片类型，且其元素可比较。
//
// 返回值:
//
//   - 返回一个 S 类型的切片，包含了按照分组规则提取的每个分组的最后一个元素。
func LastValue[B ~[]U, C ~[]S, U gotools.Comparable, S any](by B, data C) []S {

	id := array.Seq(0, len(data), 1)
	group := PairBy(by, data, id)

	value := make([]S, 0, len(data))
	numberid := make([]int, 0, len(data))

	for _, v := range group {
		value = append(value, array.Repeat(v.First[len(v.First)-1], len(v.First))...)
		numberid = append(numberid, v.Second...)
	}

	array.SortByL(gotools.ASCInt, value, numberid)

	return value
}

package group

import (
	"github.com/frankill/gotools"
	"github.com/frankill/gotools/array"
	"github.com/frankill/gotools/pair"
)

// Count2 对数据进行分组计数。
//
// 此函数接收两个参数：
//   - by：B 类型，代表分组的依据，B 是一个切片类型，其元素类型 U 必须可比较（实现 gotools.Comparable 接口）。
//   - data：C 类型，表示待分组的数据集合，C 同样是切片类型，其元素类型 S 需要支持排序（实现 gotools.Ordered 接口）。
//
// 返回值：map[U]int
// 函数返回一个 map[U]int，其中键 U 是分组的依据值，值 int 表示该组在 data 中出现的次数。
func Count2[B ~[]U, C ~[]S, U gotools.Comparable, S gotools.Ordered](by B, data C) map[U]int {

	group := By(by, data)

	res := make(map[U]int, len(group))
	for k, v := range group {
		res[k] = len(v)
	}
	return res

}

// Distinct 对数据进行分组并计算每个组的唯一元素数量。
//
// 参数：
//   - by[B ~[]U]：一个切片，其元素作为分组的依据。
//   - data[C ~[]S]：另一个切片，是需要根据 by 中的元素进行分组的数据。
//
// 返回值：
//   - map[U]int：一个映射，键为 by 中的分组依据值，值为对应组内唯一元素的数量。
func Distinct[B ~[]U, C ~[]S, U gotools.Comparable, S gotools.Ordered](by B, data C) map[U]int {

	group := By(by, data)

	res := make(map[U]int, len(group))

	for k, v := range group {
		res[k] = len(array.Distinct(v))
	}
	return res
}

// Max 根据指定的分组键对数据进行分组处理，并计算每个组中的最大值。
// 参数说明：
//   - by(B): 分组依据的键列表，类型为切片，元素需可比较。
//   - data(C): 待分组的数据列表，类型为切片，元素需要能排序。
//
// 返回值：
//   - 返回一个 map，key 为分组键 by 中的元素值，value 为对应组中的最大值。
//
// 类型约束：
//   - B: 类型为切片，元素类型为 U，要求 U 类型是可比较的。
//   - C: 类型为切片，元素类型为 S，要求 S 类型实现了 Ordered 接口，表明元素可排序。
//   - U: 可比较类型，用于分组键的元素类型。
//   - S: Ordered 类型，数据元素类型，需支持排序操作以便找出最大值。
func Max[B ~[]U, C ~[]S, U gotools.Comparable, S gotools.Ordered](by B, data C) map[U]S {

	group := By(by, data)

	res := make(map[U]S, len(group))

	for k, v := range group {
		res[k] = array.Max(v)
	}
	return res
}

// Min 根据提供的分组条件和数据集，对数据进行分组，并计算每个组的最小值。
// 参数说明：
//   - by(B): 分组的依据，类型为切片U的别名。
//   - data(C): 待分组和计算最小值的数据集，类型为切片S的别名。
//
// 返回值：
//   - 返回一个从类型U到类型S的映射，键为分组的依据，值为对应组的最小值。
//
// 类型参数：
//   - B: 类型约束为切片U的别名，表示分组键的序列类型。
//   - C: 类型约束为切片S的别名，表示数据项的序列类型。
//   - U: 可比较类型，用作分组的键。
//   - S: 有序类型，实现了gotools.Ordered接口，用于确定最小值。
//
// 注意：
//   - 确保U类型的元素可以相互比较。
//   - S类型的元素需要支持排序操作。
func Min[B ~[]U, C ~[]S, U gotools.Comparable, S gotools.Ordered](by B, data C) map[U]S {

	group := By(by, data)

	res := make(map[U]S, len(group))

	for k, v := range group {
		res[k] = array.Min(v)
	}
	return res
}

// Sum 对数据进行分组求和。
//
// 参数:
//   - by: B 类型，表示分组依据的键值列表，B 应为切片类型。
//   - data: C 类型，表示需要分组并求和的数据列表，C 也应为切片类型。
//
// 类型参数:
//   - B: 类型约束为切片 U，表示分组键的切片类型。
//   - C: 类型约束为切片 S，表示数值数据的切片类型。
//   - U: 可比较类型，作为分组的键。
//   - S: 数值类型，要求实现 Number 接口，以便进行求和操作。
//
// 返回值:
//   - map[U]S: 返回一个以 U 类型为键，S 类型为值的映射，表示每个分组的键与该组数据的总和。
func Sum[B ~[]U, C ~[]S, U gotools.Comparable, S gotools.Number](by B, data C) map[U]S {

	group := By(by, data)

	res := make(map[U]S, len(group))

	for k, v := range group {
		res[k] = array.Sum(v)
	}
	return res

}

// Pair 根据一个关键数组将两个数据数组分组。
//   - by: 用于分组的键值数组，类型为 B，需要与 frist 的元素类型 T 可比较。
//   - frist: 需要被分组的第一个数据数组，类型为 D，元素类型为 U。
//   - second: 需要被分组的第二个数据数组，类型为 O，元素类型为 S。
//
// 返回值：
//   - 一个映射，键为 T 类型（与 by 数组相同），值为 Pair 类型，包含两个分组后的数组（类型为 []U 和 []S）。
//     此函数利用泛型，可适用于不同类型的数组。如果输入的 frist、by 或 second 为空，将直接返回空映射。
func Pair[D ~[]U, B ~[]T, O ~[]S, T gotools.Comparable, S, U any](by B, frist D, second O) map[T]pair.Pair[[]U, []S] {

	res := map[T]pair.Pair[[]U, []S]{}

	if len(frist) == 0 {
		return res
	}

	if len(second) == 0 {
		return res
	}

	if len(by) == 0 {
		var t T
		res[t] = pair.Pair[[]U, []S]{
			First:  frist,
			Second: second,
		}
		return res
	}

	for i := 0; i < len(frist); i++ {
		res[by[i]] = pair.Pair[[]U, []S]{
			First:  append(res[by[i]].First, frist[i]),
			Second: append(res[by[i]].Second, second[i]),
		}
	}

	return res
}

// By 将数据切片 D 按照键值切片 B 进行分组，返回一个映射。
// 键 T 必须是可比较的类型，D 和 B 分别是元素类型为 U 和 T 的切片，U 可以为任意类型。
// 如果输入数据切片 data 或键值切片 by 为空，则直接返回空映射。
//
// 参数:
//   - data: 要进行分组的元素数据切片。
//   - by: 与 data 中元素对应的分组键值切片。
//
// 返回值:
//   - map[T][]U: 一个映射，键为 T 类型的分组标识，值为对应分组内的 U 类型元素切片。
func By[D ~[]U, B ~[]T, T gotools.Comparable, U any](by B, data D) map[T][]U {
	res := map[T][]U{}

	if len(data) == 0 {
		return res
	}
	if len(by) == 0 {
		var t T
		res[t] = data
		return res
	}

	for i := 0; i < len(data); i++ {
		res[by[i]] = append(res[by[i]], data[i])
	}
	return res
}

// ByFn 根据给定的函数进行分组计算。
//
// 参数:
//   - by: 与 data 中元素对应的分组键值切片。
//   - data: 要进行分组的元素数据切片。
//   - fn: 用于计算分组数据的函数。
//
// 返回值:
//   - map[T][]U: 一个映射，键为 T 类型的分组标识，值为对应分组内的 U 类型元素切片。
func ByFn[B ~[]T, D ~[]U, T gotools.Comparable, R, U any](fn func([]U) R, by B, data D) map[T]R {

	group := By(by, data)

	res := make(map[T]R, len(group))

	for k, v := range group {
		res[k] = fn(v)
	}
	return res

}

// ByFun2 根据给定的函数进行分组计算。
//
// 参数:
//   - fn : 用于计算分组数据的函数。
//   - fn1: 用于计算分组数据的函数。
//   - data: 要进行分组的元素数据切片。
//
// 返回值:
//   - map[T]R: 一个映射，键为 T 类型的分组标识，值为对应分组内的 R 类型元素切片。
func ByFun2[D ~[]S, T gotools.Comparable, S, R any](fn func(S) T, fn1 func([]S) R, data D) map[T]R {

	by := array.Map(fn, data)

	return ByFn(fn1, by, data)

}

//	Index 函数根据输入的切片 by 对数据进行分组，返回一个映射，映射的键是 by 中的唯一元素，值是对应元素在原切片中的索引集合。
//
// by: 需要分组的切片
//
// 返回值:
//   - map[T][]int: 包含分组信息的映射
//
// 示例:
//
//	by := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
//	grouped := Group(by)
//
//	fmt.Println(grouped)
//
//	输出:
//	map[1:[0] 2:[1 2] 3:[3 4 5] 4:[6 7 8 9]]
func Index[B ~[]T, T gotools.Comparable](by B) map[T][]int {
	res := map[T][]int{}

	if len(by) == 0 {
		return res
	}

	data := array.Seq(0, len(by), 1)

	for i := 0; i < len(by); i++ {
		res[by[i]] = append(res[by[i]], data[i])
	}
	return res
}

func IndexAsc[B ~[]T, O ~[]S, T gotools.Comparable, S gotools.Ordered](by B, order O) map[T][]int {

	data := array.Seq(0, len(by), 1)

	if len(order) == 0 {
		return By(by, data)
	}

	group := Pair(by, data, order)

	res := make(map[T][]int, len(group))

	for k, v := range group {
		array.SortByL(gotools.ASCGeneric, v.First, v.Second)
		res[k] = v.First
	}

	return res
}

func IndexDesc[B ~[]T, O ~[]S, T gotools.Comparable, S gotools.Ordered](by B, order O) map[T][]int {

	data := array.Seq(0, len(by), 1)

	if len(order) == 0 {
		return By(by, data)
	}

	group := Pair(by, data, order)

	res := make(map[T][]int, len(group))

	for k, v := range group {
		array.SortByL(gotools.DESCGeneric, v.First, v.Second)
		res[k] = v.First
	}

	return res
}

// Count 函数根据输入的切片 by 对数据进行分组，返回一个映射，映射的键是 by 中的唯一元素，值是对应元素在原切片中出现的次数。
//
// by: 需要分组的切片
//
// 返回值:
//   - map[T]int: 包含分组信息的映射
func Count[B ~[]T, T gotools.Comparable](by B) map[T]int {

	if len(by) == 0 {
		return map[T]int{}
	}

	res := map[T]int{}
	for i := 0; i < len(by); i++ {
		res[by[i]]++
	}
	return res
}

// Order 根据给定的排序规则对数组进行分组操作。
//
// 参数说明：
//   - data(D): 待分组的数据，类型为 D，要求 D 是 U 类型的切片。
//   - by(B): 分组依据的键值数组，类型为 B，要求 B 是 T 类型的切片，其中 T 需要可比较。
//   - order(O): 排序依据的序列，类型为 O，要求 O 是 S 类型的切片，S 必须实现了有序接口Ordered。
//
// 返回值：
//   - map[T][]U: 返回一个字典，键为 T 类型，值为对应分组的 U 类型元素组成的切片。
//
// 注意：
//   - 当 order 为空时，直接调用 GroupArray 进行分组。
//   - 若 order 提供了排序依据，函数首先依据 by 和 order 进行分组及排序，然后将排序后的数据作为结果值。
func Order[D ~[]U, B ~[]T, O ~[]S, T gotools.Comparable, S gotools.Ordered, U any](by B, data D, order O) map[T][]U {

	if len(order) == 0 {
		return By(by, data)
	}

	group := Pair(by, data, order)

	res := make(map[T][]U, len(group))

	for k, v := range group {
		array.SortByL(gotools.ASCGeneric, v.First, v.Second)
		res[k] = v.First
	}

	return res
}

// OrderDesc Order倒序版
func OrderDesc[D ~[]U, B ~[]T, O ~[]S, T gotools.Comparable, S gotools.Ordered, U any](by B, data D, order O) map[T][]U {

	if len(order) == 0 {
		return By(by, data)
	}

	group := Pair(by, data, order)

	res := make(map[T][]U, len(group))

	for k, v := range group {
		array.SortByL(gotools.DESCGeneric, v.First, v.Second)
		res[k] = v.First
	}

	return res
}

package group

import (
	"github.com/frankill/gotools"
	"github.com/frankill/gotools/array"
	"github.com/frankill/gotools/maps"
	"github.com/frankill/gotools/pair"
)

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

// CountBy 对数据进行分组计数。
//
// 此函数接收两个参数：
//   - by：B 类型，代表分组的依据，B 是一个切片类型，其元素类型 U 必须可比较（实现 gotools.Comparable 接口）。
//   - data：C 类型，表示待分组的数据集合，C 同样是切片类型，其元素类型 S 需要支持排序（实现 gotools.Ordered 接口）。
//
// 返回值：map[U]int
// 函数返回一个 map[U]int，其中键 U 是分组的依据值，值 int 表示该组在 data 中出现的次数。
func CountBy[B ~[]U, C ~[]S, U gotools.Comparable, S gotools.Ordered](by B, data C) map[U]int {

	group := By(by, data)

	res := make(map[U]int, len(group))
	for k, v := range group {
		res[k] = len(v)
	}
	return res

}

// DistinctBy 对数据进行分组并计算每个组的唯一元素数量。
//
// 参数：
//   - by[B ~[]U]：一个切片，其元素作为分组的依据。
//   - data[C ~[]S]：另一个切片，是需要根据 by 中的元素进行分组的数据。
//
// 返回值：
//   - map[U]int：一个映射，键为 by 中的分组依据值，值为对应组内唯一元素的数量。
func DistinctBy[B ~[]U, C ~[]S, U gotools.Comparable, S gotools.Ordered](by B, data C) map[U]int {

	group := By(by, data)

	res := make(map[U]int, len(group))

	for k, v := range group {
		res[k] = len(array.Distinct(v))
	}
	return res
}

// MaxBy 根据指定的分组键对数据进行分组处理，并计算每个组中的最大值。
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
func MaxBy[B ~[]U, C ~[]S, U gotools.Comparable, S gotools.Ordered](by B, data C) map[U]S {

	group := By(by, data)

	res := make(map[U]S, len(group))

	for k, v := range group {
		res[k] = array.Max(v)
	}
	return res
}

// MinBy 根据提供的分组条件和数据集，对数据进行分组，并计算每个组的最小值。
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
func MinBy[B ~[]U, C ~[]S, U gotools.Comparable, S gotools.Ordered](by B, data C) map[U]S {

	group := By(by, data)

	res := make(map[U]S, len(group))

	for k, v := range group {
		res[k] = array.Min(v)
	}
	return res
}

// SumBy 对数据进行分组求和。
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
func SumBy[B ~[]U, C ~[]S, U gotools.Comparable, S gotools.Number](by B, data C) map[U]S {

	group := By(by, data)

	res := make(map[U]S, len(group))

	for k, v := range group {
		res[k] = array.Sum(v)
	}
	return res
}

// PairBy 根据一个关键数组将两个数据数组分组。
//   - by: 用于分组的键值数组，类型为 B，需要与 frist 的元素类型 T 可比较。
//   - frist: 需要被分组的第一个数据数组，类型为 D，元素类型为 U。
//   - second: 需要被分组的第二个数据数组，类型为 O，元素类型为 S。
//
// 返回值：
//   - 一个映射，键为 T 类型（与 by 数组相同），值为 Pair 类型，包含两个分组后的数组（类型为 []U 和 []S）。
//     此函数利用泛型，可适用于不同类型的数组。如果输入的 frist、by 或 second 为空，将直接返回空映射。
func PairBy[D ~[]U, B ~[]T, O ~[]S, T gotools.Comparable, S, U any](by B, frist D, second O) map[T]pair.Pair[[]U, []S] {

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

// ByFunc 根据给定的函数进行分组。
//
// 参数:
//   - f: 用于分组的函数。
//   - data: 要进行分组的元素数据切片。
//
// 返回值:
//   - map[T][]U: 一个映射，键为 T 类型的分组标识，值为对应分组内的 U 类型元素切片。
func ByFunc[D ~[]U, T gotools.Comparable, U any](f func(U) T, data D) map[T][]U {

	if len(data) == 0 {
		return map[T][]U{}
	}

	res := map[T][]U{}

	for i := 0; i < len(data); i++ {
		res[f(data[i])] = append(res[f(data[i])], data[i])
	}
	return res
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

// ArrayByOrder 根据给定的排序规则对数组进行分组操作。
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
func ByOrder[D ~[]U, B ~[]T, O ~[]S, T gotools.Comparable, S gotools.Ordered, U any](by B, data D, order O) map[T][]U {

	if len(order) == 0 {
		return By(by, data)
	}

	group := PairBy(by, data, order)

	res := make(map[T][]U, len(group))

	for k, v := range group {
		array.SortByL(gotools.ASCGeneric, v.First, v.Second)
		res[k] = v.First
	}

	return res
}

// ArrayByOrderDesc GroupArrayByOrder倒序版
func ByOrderDesc[D ~[]U, B ~[]T, O ~[]S, T gotools.Comparable, S gotools.Ordered, U any](by B, data D, order O) map[T][]U {

	if len(order) == 0 {
		return By(by, data)
	}

	group := PairBy(by, data, order)

	res := make(map[T][]U, len(group))

	for k, v := range group {
		array.SortByL(gotools.DESCGeneric, v.First, v.Second)
		res[k] = v.First
	}

	return res
}

func WindowFun[B ~[]U, C ~[]S, U gotools.Comparable, S gotools.Ordered](by B, order C) func(func(by B, data C) []S) []S {

	return func(fn func(by B, data C) []S) []S {
		return fn(by, order)
	}

}

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

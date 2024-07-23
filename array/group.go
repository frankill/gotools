package array

import (
	"cmp"
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
func Retention[B ~[]U, C ~[]S, U comparable, S any](by B, data C) func(fun ...Retentionfun[S]) map[U][]bool {

	group := GroupData(by, data)

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

// GroupCount 对数据进行分组计数。
//
// 此函数接收两个参数：
//   - by：B 类型，代表分组的依据，B 是一个切片类型，其元素类型 U 必须可比较（实现 comparable 接口）。
//   - data：C 类型，表示待分组的数据集合，C 同样是切片类型，其元素类型 S 需要支持排序（实现 cmp.Ordered 接口）。
//
// 返回值：map[U]int
// 函数返回一个 map[U]int，其中键 U 是分组的依据值，值 int 表示该组在 data 中出现的次数。
func GroupCount[B ~[]U, C ~[]S, U comparable, S cmp.Ordered](by B, data C) map[U]int {

	group := GroupData(by, data)

	res := make(map[U]int, len(group))
	for k, v := range group {
		res[k] = len(v)
	}
	return res

}

// GroupDistinct 对数据进行分组并计算每个组的唯一元素数量。
//
// 参数：
// by[B ~[]U]：一个切片，其元素作为分组的依据。
// data[C ~[]S]：另一个切片，是需要根据 by 中的元素进行分组的数据。
// U：类型参数，要求是可比较的类型，用于约束 by 切片的元素类型。
// S：类型参数，要求实现了 ordered 接口，用于约束 data 切片的元素类型。
//
// 返回值：
// map[U]int：一个映射，键为 by 中的分组依据值，值为对应组内唯一元素的数量。
func GroupDistinct[B ~[]U, C ~[]S, U comparable, S cmp.Ordered](by B, data C) map[U]int {

	group := GroupData(by, data)

	res := make(map[U]int, len(group))

	for k, v := range group {
		res[k] = len(ArrayDistinct(v))
	}
	return res
}

type Groupfun[T any] func(x []T) T
type GroupFilterfun[K comparable, T any] func(x K, y []T) []bool

// GroupGenerate 是一个函数，用于根据给定的分组条件和数据，生成一个函数，该函数接受一个 Groupfun 类型的参数，并返回一个映射。
// 该函数将根据给定的分组条件将数据分组，并对每个分组应用给定的函数，并将结果存储在映射中。
//
// 参数：
//   - by: 分组条件，类型为 B，必须是可比较类型的切片。
//   - data: 数据，类型为 C，必须是可比较类型的切片。
//
// 返回值：
//
//	一个函数，该函数接受一个 Groupfun 类型的参数，并返回一个映射。
//	该映射的键类型为 U，值类型为 S。
func GroupGenerate[B ~[]U, C ~[]S, U comparable, S any](by B, data C) func(fun Groupfun[S]) map[U]S {

	group := GroupData(by, data)

	return func(fun Groupfun[S]) map[U]S {

		res := make(map[U]S, len(group))
		for k, v := range group {
			res[k] = fun(v)
		}
		return res

	}
}

// GroupGenerateFilter 根据指定的分组规则（by）和数据集（data），
// 生成一个过滤器函数。此过滤器函数允许用户自定义过滤逻辑（通过 GroupFilterfun），
// 并返回根据该逻辑处理后的分组,数据。
//
// 参数:
//   - by: 分组依据的元素切片，类型为B，约束为可比较的切片。
//   - data: 需要分组和过滤的原始数据切片，类型为C。
//
// 返回值:
//   - 一个函数，接受一个GroupFilterfun类型的参数用于定制过滤条件，
//     并返回过滤后的分组，数据 （类型为 []U 和 []S）。
func GroupGenerateFilter[B ~[]U, C ~[]S, U comparable, S any](by B, data C) func(fun GroupFilterfun[U, S]) ([]U, []S) {

	id := ArraySeq(0, len(data), 1)
	group := GroupPair(by, data, id)

	return func(fun GroupFilterfun[U, S]) ([]U, []S) {

		value := make([]S, 0, len(data))
		valueby := make([]U, 0, len(by))

		numberid := make([]int, 0, len(data))

		for k, v := range group {

			farr := fun(k, v.First)

			vv := ArrayFilter(func(x ...any) bool { return x[1].(bool) }, ArrayToAny(v.First), ArrayToAny(farr))

			value = append(value, ArrayToGeneric[S](vv)...)

			valueby = append(valueby, Rep(k, len(vv))...)

			numberid = append(numberid, ArrayFilter(func(x ...int) bool { return x[1] > 0 }, v.Second, ArrayToint(farr))...)

		}

		numberids := ArrayCopy(numberid)

		ArraySortByL(ASCInt, value, numberid)
		ArraySortByL(ASCInt, valueby, numberids)

		return valueby, value

	}
}

/*
GroupMax 根据指定的分组键对数据进行分组处理，并计算每个组中的最大值。

参数说明：
- by(B): 分组依据的键列表，类型为切片，元素需可比较。
- data(C): 待分组的数据列表，类型为切片，元素需要能排序。

返回值：
- 返回一个 map，key 为分组键 by 中的元素值，value 为对应组中的最大值。

类型约束：
- B: 类型为切片，元素类型为 U，要求 U 类型是可比较的。
- C: 类型为切片，元素类型为 S，要求 S 类型实现了 Ordered 接口，表明元素可排序。
- U: 可比较类型，用于分组键的元素类型。
- S: Ordered 类型，数据元素类型，需支持排序操作以便找出最大值。
*/
func GroupMax[B ~[]U, C ~[]S, U comparable, S cmp.Ordered](by B, data C) map[U]S {

	group := GroupData(by, data)

	res := make(map[U]S, len(group))

	for k, v := range group {
		res[k] = ArrayMax(v)
	}
	return res
}

/*
GroupMin 根据提供的分组条件和数据集，对数据进行分组，并计算每个组的最小值。

参数说明：
- by(B): 分组的依据，类型为切片U的别名。
- data(C): 待分组和计算最小值的数据集，类型为切片S的别名。

返回值：
- 返回一个从类型U到类型S的映射，键为分组的依据，值为对应组的最小值。

类型参数：
- B: 类型约束为切片U的别名，表示分组键的序列类型。
- C: 类型约束为切片S的别名，表示数据项的序列类型。
- U: 可比较类型，用作分组的键。
- S: 有序类型，实现了cmp.Ordered接口，用于确定最小值。

注意：
- 确保U类型的元素可以相互比较。
- S类型的元素需要支持排序操作。
*/
func GroupMin[B ~[]U, C ~[]S, U comparable, S cmp.Ordered](by B, data C) map[U]S {

	group := GroupData(by, data)

	res := make(map[U]S, len(group))

	for k, v := range group {
		res[k] = ArrayMin(v)
	}
	return res
}

// GroupSum 对数据进行分组求和。
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
func GroupSum[B ~[]U, C ~[]S, U comparable, S Number](by B, data C) map[U]S {

	group := GroupData(by, data)

	res := make(map[U]S, len(group))

	for k, v := range group {
		res[k] = ArraySum(v)
	}
	return res
}

/*
GroupArrayPair 根据一个关键数组将两个数据数组分组。
- by: 用于分组的键值数组，类型为 B，需要与 frist 的元素类型 T 可比较。
- frist: 需要被分组的第一个数据数组，类型为 D，元素类型为 U。
- second: 需要被分组的第二个数据数组，类型为 O，元素类型为 S。
返回值：
- 一个映射，键为 T 类型（与 by 数组相同），值为 Pair 类型，包含两个分组后的数组（类型为 []U 和 []S）。

此函数利用泛型，可适用于不同类型的数组。如果输入的 frist、by 或 second 为空，将直接返回空映射。
*/
func GroupPair[D ~[]U, B ~[]T, O ~[]S, T comparable, S, U any](by B, frist D, second O) map[T]Pair[[]U, []S] {

	res := map[T]Pair[[]U, []S]{}

	if len(frist) == 0 {
		return res
	}

	if len(second) == 0 {
		return res
	}

	if len(by) == 0 {
		var t T
		res[t] = Pair[[]U, []S]{
			First:  frist,
			Second: second,
		}
		return res
	}

	for i := 0; i < len(frist); i++ {
		res[by[i]] = Pair[[]U, []S]{
			First:  append(res[by[i]].First, frist[i]),
			Second: append(res[by[i]].Second, second[i]),
		}
	}

	return res
}

// GroupArray 将数据切片 D 按照键值切片 B 进行分组，返回一个映射。
// 键 T 必须是可比较的类型，D 和 B 分别是元素类型为 U 和 T 的切片，U 可以为任意类型。
// 如果输入数据切片 data 或键值切片 by 为空，则直接返回空映射。
//
// 参数:
// - data: 要进行分组的元素数据切片。
// - by: 与 data 中元素对应的分组键值切片。
//
// 返回值:
// - map[T][]U: 一个映射，键为 T 类型的分组标识，值为对应分组内的 U 类型元素切片。
func GroupData[D ~[]U, B ~[]T, T comparable, U any](by B, data D) map[T][]U {
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

// Group 函数根据输入的切片 by 对数据进行分组，返回一个映射，映射的键是 by 中的唯一元素，值是对应元素在原切片中的索引集合。
//
// by: 需要分组的切片
//
// 返回值:
// - map[T][]int: 包含分组信息的映射
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
func GroupLocation[B ~[]T, T comparable](by B) map[T][]int {
	res := map[T][]int{}

	if len(by) == 0 {
		return res
	}

	data := ArraySeq(0, len(by), 1)

	for i := 0; i < len(by); i++ {
		res[by[i]] = append(res[by[i]], data[i])
	}
	return res
}

// GroupArrayByOrder 根据给定的排序规则对数组进行分组操作。
//
// 参数说明：
// - data(D): 待分组的数据，类型为 D，要求 D 是 U 类型的切片。
// - by(B): 分组依据的键值数组，类型为 B，要求 B 是 T 类型的切片，其中 T 需要可比较。
// - order(O): 排序依据的序列，类型为 O，要求 O 是 S 类型的切片，S 必须实现了有序接口Ordered。
//
// 返回值：
// - map[T][]U: 返回一个字典，键为 T 类型，值为对应分组的 U 类型元素组成的切片。
//
// 注意：
// - 当 order 为空时，直接调用 GroupArray 进行分组。
// - 若 order 提供了排序依据，函数首先依据 by 和 order 进行分组及排序，然后将排序后的数据作为结果值。
func GroupByOrder[D ~[]U, B ~[]T, O ~[]S, T comparable, S cmp.Ordered, U any](by B, data D, order O) map[T][]U {

	if len(order) == 0 {
		return GroupData(by, data)
	}

	group := GroupPair(by, data, order)

	res := make(map[T][]U, len(group))

	for k, v := range group {
		ArraySortByL(ASCGeneric, v.First, v.Second)
		res[k] = v.First
	}

	return res
}

// GroupArrayByOrderDesc GroupArrayByOrder倒序版
func GroupByOrderDesc[D ~[]U, B ~[]T, O ~[]S, T comparable, S cmp.Ordered, U any](by B, data D, order O) map[T][]U {

	if len(order) == 0 {
		return GroupData(by, data)
	}

	group := GroupPair(by, data, order)

	res := make(map[T][]U, len(group))

	for k, v := range group {
		ArraySortByL(DESCGeneric, v.First, v.Second)
		res[k] = v.First
	}

	return res
}

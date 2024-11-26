package array

import "github.com/frankill/gotools"

var (
	LETTERS = []string{
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	}

	Letters = []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	}
)

// Which 函数根据提供的判断函数fun，遍历多个切片数组arr的对应位置元素，判断这些元素是否满足fun定义的条件，
// 并记录满足条件的位置（以1表示）和不满足条件的位置（以0表示）到结果切片中返回。
//
// 参数:
//   - fun: 一个通用型函数，接受任意数量的T类型参数并返回一个布尔值，用于判断一组元素是否满足特定条件。
//   - arr: 可变数量的切片数组，每个切片包含相同数量的T类型元素，这些切片将被并行遍历以应用fun函数。
//
// 返回值:
//   - 一个int类型的切片，表示在对应位置上元素是否满足fun条件的结果。满足条件的位置记为1，不满足则记为0。
//
// 注意:
//
//   - 所有输入的切片arr必须具有相同的长度。
//
//   - 使用了Go泛型[T any]，允许该函数接受任何类型的参数和切片。
func Which[T any](fun func(x ...T) bool, arr ...[]T) []int {

	if len(arr) == 0 {
		return []int{}
	}

	if len(arr[0]) == 0 {
		return []int{}
	}

	l := len(arr[0])
	f := len(arr)

	param := make([]T, f)

	result := make([]int, 0, len(arr[0]))

	for i := 0; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i]
		}
		if fun(param...) {
			result = append(result, 1)
		} else {
			result = append(result, 0)
		}
	}
	return result
}

func Toint(arr []bool) []int {

	result := make([]int, len(arr))

	for i, v := range arr {
		if v {
			result[i] = 1
		} else {
			result[i] = 0
		}
	}
	return result
}

func ToGeneric[T any](arr []any) []T {
	result := make([]T, len(arr))

	for i, v := range arr {
		result[i] = v.(T)
	}
	return result
}

// ToAny 将类型为 S（元素类型为 T）的切片转换为任意类型切片。
// 要求 T 类型实现 `any` 接口。
//
// 参数:
//   - arr: 要转换的切片 S。
//
// 返回值:
//   - 返回一个新的 []any 类型切片，包含输入的所有元素，保持原有的顺序。
func ToAny[S ~[]T, T any](arr S) []any {

	res := make([]any, 0, len(arr))

	for i := 0; i < len(arr); i++ {
		res = append(res, arr[i])
	}
	return res
}

// ToZero 将类型为 S（元素类型为 T）的切片转换为 map[T]int 类型。
// 要求 T 类型实现 `Comparable` 接口。
//
// 参数:
//   - arr: 要转换的切片 S。
//
// 返回值:
//   - 返回一个 map[T]int 类型的映射， 包含输入的所有元素，键为元素，值为 0。
func ToZero[S ~[]T, T gotools.Comparable](arr S) map[T]int {

	res := make(map[T]int)

	for i := 0; i < len(arr); i++ {
		res[arr[i]] = 0
	}
	return res
}

// ArrayCopy 创建给定切片的一个副本。
// 它返回一个新的切片，其中包含与输入切片相同的元素。
// 参数:
//
//   - arr - 要复制的切片。
//
// 返回:
//
//   - 一个新的切片，它是输入切片的副本。
func Copy[T any](arr []T) []T {
	return CopyWithNum(arr, len(arr))
}

// ArrayCopyWithNum 复制给定切片并创建一个新的切片，新切片的长度为 num。
// 如果 num 大于原切片长度，则重复原切片中的元素以填充新切片。
// 参数:
//
//   - arr - 要复制的切片。
//   - num - 新切片的长度。
//
// 返回:
//
//   - 一个新的切片，其长度为 num，且元素为 arr 的重复。
func CopyWithNum[T any](arr []T, num int) []T {

	res := make([]T, num)
	copy(res, arr)

	la := len(arr)
	if la < num {

		for i, j := la, 0; i < num; i, j = i+1, j+1 {
			if j%la == 0 {
				j = 0
			}
			res[i] = arr[j]
		}

	}

	return res
}

// Transform 将切片 x 中的元素根据 array_from 和 array_to 的映射关系转换为新的类型 T，
// 并返回转换后的切片。如果 x 中的元素在映射中不存在，则使用 default_value 填充。
// 参数:
//
//   - x - 待转换的切片。
//   - array_from - 映射的键值切片。
//   - array_to - 映射的值切片。
//   - default_value - 当 x 中的元素在映射中不存在时使用的默认值。
//
// 返回:
//
//	一个新切片，其中元素为转换后的值。
func Transform[S ~[]F, D ~[]T, F gotools.Comparable, T any](x S, array_from S, array_to D, default_value T) []T {

	if len(x) == 0 {
		return []T{}
	}

	if len(array_from) == 0 || len(array_to) == 0 {
		return []T{}
	}

	la := len(x)

	res := make([]T, la)

	dict := ToMap2(array_from, array_to)

	for i := 0; i < la; i++ {
		res[i] = default_value
		if v, ok := dict[x[i]]; ok {
			res[i] = v
		}
	}
	return res

}

// ToMap2 根据两个数组创建一个映射。
// 它接受两个参数：key和value，它们分别是K和V类型的数组。
// 函数返回一个map[T]S类型的映射，其中T和S是K和V数组元素的类型。
// 该函数的目的是通过索引匹配将key数组的元素作为映射的键，value数组的元素作为对应的值。
// 参数:
//
//   - key K: 用于映射的键的数组。
//   - value V: 与键对应的值的数组。
//
// 返回值:
//
//   - map[T]S: 一个映射，其中键来自key数组，值来自value数组。
//
// 注意:
//
//   - K和V的类型参数必须是数组类型，且K的元素类型必须是可比较的。
//   - 这个函数假设key和value数组的长度是相同的，以便进行索引匹配。
func ToMap2[K ~[]T, V ~[]S, T gotools.Comparable, S any](key K, value V) map[T]S {

	dict := make(map[T]S, len(key))

	for i := range key {

		dict[key[i]] = value[i]

	}

	return dict

}

// HasOrderMaxCount 检查在给定的映射存在性验证下，切片 source 中有多少个元素可以按照切片 match 的顺序出现。
// 在source中出现了match中未出现的元素，停止计算返回
// 这个函数返回的是 source 中与 match 顺序一致的最长连续子序列的长度。
//
// 参数:
//
//   - a: 类型为 S 的切片，S 是一个泛型切片，其元素类型为 T。
//   - b: 类型同样为 S 的切片，作为参考顺序。
//   - exist: 一个映射，键为 T 类型，用于快速判断元素是否存在于某个集合中。
//
// 返回值:
//
//   - 返回一个整数，表示 source 中与 match 顺序一致的最长连续子序列的长度。
func HasOrderMaxCount[S ~[]T, T gotools.Comparable](source, match S, exist map[T]struct{}) int {

	count := 0
	j := 0 // 用于追踪 b 切片中的当前元素位置

	n := len(match)
	for i := range source {

		if _, ok := exist[source[i]]; !ok {
			break
		}

		if j < n && source[i] == match[j] {
			count++
			j++
		} else if j >= n {

			break
		}

	}

	return count

}

// HasDupMaxCount 计算源切片中与匹配切片元素相等的元素个数，
// 并在遇到第一个重复元素时停止计算。
// 这个函数使用泛型，S 表示切片的类型，而 T 表示切片元素的类型，
// 其中 T 必须是可比较的。
//
// 参数：
//
//   - source: 待检查的源切片。
//   - match: 用于匹配的参考切片。
//
// 返回值：
//
//   - 返回源切片中与 match 中元素相等的元素个数。
//   - 如果在计数过程中发现源切片中有重复元素，则提前终止并返回当前计数值。
func HasDupMaxCount[S ~[]T, T gotools.Comparable](source, match S) int {
	count := 0
	j := 0 // 用于追踪 b 切片中的当前元素位置

	n := len(source) - 1
	nn := len(match)

	for i := range source {

		if j < nn && source[i] == match[j] {
			count++
			j++
		} else if j >= nn {

			break
		}

		if i < n && source[i] == source[i+1] {
			break
		}

	}

	return count
}

// HasIncreaseMaxCount 计算源切片中与匹配切片元素相等的连续元素个数。
// 此函数使用泛型，其中 S 表示切片的类型，T 表示切片元素的类型，
// 要求 T 必须是可比较的类型。
//
// 参数：
//
//   - source: 待检查的源切片。
//   - match: 用于匹配的参考切片。
//
// 返回值：
//
//   - 返回源切片中与 match 中元素相等的连续元素个数。
//   - 当 match 中的所有元素都已在 source 中找到匹配或 source 中不再有匹配元素时，停止计数并返回结果。
func HasIncreaseMaxCount[S ~[]T, T gotools.Comparable](source, match S) int {
	count := 0
	j := 0 // 用于追踪 b 切片中的当前元素位置

	n := len(match)
	for i := range source {

		if j < n && source[i] == match[j] {
			count++
			j++
		} else if j >= n {
			break
		}

	}

	return count
}

// Repeat 创建一个包含相同元素的切片。
// 它接受一个类型为 T 的元素和一个整数 n，返回一个长度为 n 的切片，其中所有元素都与 x 相同。
// 这个函数的目的是为了方便地初始化一个切片，当切片的所有元素都相同时，可以避免重复的初始化代码。
// 参数:
//
//   - x: 类型为 T 的元素，切片中的每个元素都将复制这个值。
//   - n: 切片的长度，指定切片将包含多少个元素。
//
// 返回值:
//
//   - 一个类型为 []T 的切片，长度为 n，其中所有元素都等于 x。
func Repeat[T any](x T, n int) []T {
	result := make([]T, n)
	for i := 0; i < n; i++ {
		result[i] = x
	}
	return result
}

// Cycle 用于实现循环队列。
type Cycle[T any] struct {
	index int
	len   int
	data  []T
}

// Index 返回循环队列的当前索引。
func (c *Cycle[T]) Index() int {
	return c.index
}

// Len 返回循环队列的长度。
func (c *Cycle[T]) Len() int {
	return c.len
}

// 创建一个新的循环队列。
func NewCycle[T any](length int) *Cycle[T] {
	data := make([]T, length)
	return &Cycle[T]{
		index: 0,
		len:   len(data),
		data:  data,
	}
}

// Push 将元素 d 添加到循环队列的index位置,并更新index
func (c *Cycle[T]) Push(d T) {

	if c.index == c.len {
		c.index = 0
	}

	c.data[c.index%c.len] = d
	c.index++
}

// Pop 从循环队列的index位置取出元素
func (c *Cycle[T]) Pop() T {
	return c.data[c.index%c.len]
}

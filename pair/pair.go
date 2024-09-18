package pair

import "github.com/frankill/gotools"

type Pair[T, U any] struct {
	First  T
	Second U
}

// Of 创建一个包含两种类型元素的Pair。
// 这个函数展示了Go的泛型特性，允许创建泛型对，其中每个元素可以是任何类型。
// 参数:
//
//	first - Pair中的第一个元素，类型为T。
//	second - Pair中的第二个元素，类型为U。
//
// 返回值:
//
//	返回一个Pair[T, U]实例，其中T和U是泛型类型，可以是任何类型。
//
// 使用场景:
//
//	当需要将两个不相关或不同类型的数据作为一个单元返回时，可以使用PairOf函数。
func Of[T, U any](first T, second U) Pair[T, U] {
	return Pair[T, U]{
		First:  first,
		Second: second,
	}
}

// FromMap 将一个映射转换为一个包含键值对的切片。
// 参数:
//   - data: 一个映射，键类型为 K，值类型为 V。
//
// 返回:
//   - 一个切片，其中包含所有映射中的键值对，每个键值对表示为 Pair[K, V]。
//
// 函数功能:
//   - 遍历映射 data，将每个键值对转换为 Pair[K, V] 类型，并将其添加到切片 pairs 中。
//   - 返回包含所有键值对的切片 pairs。
func FromMap[K gotools.Comparable, V any](data map[K]V) []Pair[K, V] {
	pairs := make([]Pair[K, V], 0, len(data))
	for k, v := range data {
		pairs = append(pairs, Pair[K, V]{First: k, Second: v})
	}
	return pairs
}

// FromArray 根据两个数组创建一个包含对应元素对的数组。
// 参数 first 和 second 分别代表两个源数组。
// 返回值是一个 Pair[T, V] 类型的数组，其中每个 Pair 包含来自 first 和 second 的对应元素。
// 如果任一输入数组为空，则返回一个空的 Pair 数组。
// 该函数允许泛型使用，可以适用于任何类型的数组。
//
// 注意：Pair[T, V] 必须是一个已定义的结构体，且包含 First 和 Second 两个字段。
func FromArray[F ~[]T, S ~[]V, T, V any](first F, second S) []Pair[T, V] {

	if len(first) == 0 || len(second) == 0 {
		return []Pair[T, V]{}
	}

	pairs := make([]Pair[T, V], 0, len(first))

	for i := range first {
		pairs = append(pairs, Pair[T, V]{First: first[i], Second: second[i]})
	}

	return pairs

}

// Firsts 提取一个包含配对的第一个元素的切片。
// 该函数接受一个类型为D的切片，其中D是Pair[K, V]类型的切片的约束。
// 它返回一个K类型的切片，包含了输入切片中每个配对的第一个元素。
// 参数:
//
//	data D - 一个包含Pair[K, V]类型元素的切片。
//
// 返回值:
//
//	[]K - 一个切片，包含输入切片中每个配对的第一个元素。
func Firsts[D ~[]Pair[K, V], K any, V any](data D) []K {

	if len(data) == 0 {
		return []K{}
	}

	res := make([]K, 0, len(data))

	for i := range data {
		res = append(res, data[i].First)
	}
	return res
}

// Seconds 提取一个包含配对数据的切片中所有第二个元素，返回一个只包含这些第二个元素的切片。
// 该函数适用于任何实现了配对（Pair）接口的切片类型D，其中每个配对包含键（K）和值（V）。
// 参数:
//
//	data D: 一个包含配对数据的切片，每个配对包含一个键和一个值。
//
// 返回值:
//
//	[]V: 一个切片，包含输入切片中所有配对的第二个元素（值）。
//
// 如果输入切片为空，函数将返回一个空切片，以避免处理空切片时的潜在错误。
func Seconds[D ~[]Pair[K, V], K any, V any](data D) []V {
	ld := len(data)
	if ld == 0 {
		return []V{}
	}

	res := make([]V, 0, ld)

	for i := range data {
		res = append(res, data[i].Second)
	}
	return res
}

func ToMap[D ~[]Pair[K, V], K gotools.Comparable, V any](data D) map[K][]V {
	res := make(map[K][]V, len(data))
	for _, p := range data {
		if res[p.First] == nil {
			res[p.First] = []V{p.Second}
		} else {
			res[p.First] = append(res[p.First], p.Second)
		}
	}
	return res
}

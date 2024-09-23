package maps

import (
	"github.com/frankill/gotools"
	"github.com/frankill/gotools/array"
	"github.com/frankill/gotools/pair"
)

// OfIndex 从给定的数据中创建一个映射，其中键是函数 f 的返回值，值是每个元素在数据中的位置。
// 该函数允许泛型使用，确保键类型K是可比较的，值类型V可以是任何类型。
// 参数:
//   - f: 一个函数，它接受每个元素作为输入，返回一个可比较的键。
//   - data: 一个切片，它包含要处理的数据。
//
// 返回值:
//   - 一个映射，其中键是函数 f 的返回值，值是每个元素在数据中的位置。
func OfIndex[K gotools.Comparable, V any](f func(V) K, data []V) map[K][]int {

	result := map[K][]int{}

	for i := range data {
		result[f(data[i])] = append(result[f(data[i])], i)
	}

	return result

}

// OfCount 从给定的数据中创建一个映射，其中键是函数 f 的返回值，值是每个元素在数据中出现的次数。
// 该函数允许泛型使用，确保键类型K是可比较的，值类型V可以是任何类型。
// 参数:
//   - f: 一个函数，它接受每个元素作为输入，返回一个可比较的键。
//   - data: 一个切片，它包含要处理的数据。
//
// 返回值:
//   - 一个映射，其中键是函数 f 的返回值，值是每个元素在数据中出现的次数。
func OfCount[K gotools.Comparable, V any](f func(V) K, data []V) map[K]int {

	result := map[K]int{}

	for i := range data {
		result[f(data[i])]++
	}

	return result

}

// Keys 返回给定映射的所有键的切片。
// 该函数允许泛型使用，确保键类型K是可比较的，值类型V可以是任何类型。
// 参数:
//
//	m: 一个映射，它的键将被收集到切片中。
//
// 返回值:
//
//	一个包含映射所有键的切片。
//
// 使用场景:
//
//	当需要单独处理或传递映射的键集合时，此函数非常有用。
func Keys[K gotools.Comparable, V any](m ...map[K]V) []K {

	if len(m) == 0 {
		return make([]K, 0)
	}

	num := array.Sum(array.Map(func(x ...map[K]V) int { return len(x[0]) }, m))

	keys := make([]K, 0, num)

	for _, v := range m {
		for k := range v {
			keys = append(keys, k)
		}
	}

	return keys
}

// MpaValues 将map的值转换为切片返回。
// 该函数接受一个类型为[K gotools.Comparable, V any]的map作为参数，
// 其中K代表map的键类型，必须实现gotools.Comparable接口，
// V代表map的值类型，可以是任何类型。
// 函数返回一个类型为[]V的切片，包含了map中所有值的副本。
// 这个函数的存在是为了在不改变原map值的情况下，
// 提供一个方便的方式来使用或处理map的值集合。
//
// 参数:
//
//	m: 一个类型为map[K]V的map，其中K是可比较的，V可以是任何类型。
//
// 返回值:
//
//	一个类型为[]V的切片，包含了map m中所有值的副本。
func Values[K gotools.Comparable, V any](m ...map[K]V) []V {

	if len(m) == 0 {
		return make([]V, 0)
	}

	num := array.Sum(array.Map(func(x ...map[K]V) int { return len(x[0]) }, m))

	values := make([]V, 0, num)

	for _, v := range m {
		for k := range v {
			values = append(values, v[k])
		}
	}

	return values
}

// Filter 过滤给定映射中的元素，仅保留满足特定条件的键值对。
// 参数 f 是一个函数，用于测试每个键值对是否满足条件。如果满足，则返回 true，否则返回 false。
// 参数 m 是要过滤的映射。
// 返回值是一个新的映射，仅包含满足条件的键值对。
// Filter 的类型参数 K 和 V 分别表示映射的键类型和值类型，其中 K 必须实现 gotools.Comparable 接口，以支持作为映射的键。
func Filter[K gotools.Comparable, V any](f func(K, V) bool, m map[K]V) map[K]V {
	filtered := make(map[K]V)
	for k, v := range m {
		if f(k, v) {
			filtered[k] = v
		}
	}
	return filtered
}

// FilterArr 根据给定的数组过滤给定映射中的元素，仅保留满足特定条件的键值对。
// 参数:
//
//	m: 要过滤的映射。
//
// arr : 要过滤的数组。
// 返回值:
//
//	一个新的映射，其中键都是在 arr 中存在的键，

func FilterArr[S ~[]K, K gotools.Comparable, V any](m map[K]V, arr S) map[K]V {
	filtered := make(map[K]V)
	for k, v := range m {

		if array.Has(arr, k) {
			filtered[k] = v
		}

	}
	return filtered
}

// ApplyValue 对给定映射中的每个元素应用函数 f，并返回一个新的映射，
// 其中包含原始键和应用函数后得到的新值。
// 参数:
//
//	f: 一个函数，它接受映射中的键和值作为参数，并返回一个新值。
//	m: 要应用函数的映射。
//
// 返回值:
//
//	一个新的映射，其中每个键都对应于原始映射中的键，
//	但其值是通过应用函数 f 得到的新值。
//
// 使用场景:
//
//	当需要对映射中的所有值进行某种转换时，此函数非常有用。
//	它允许在不改变原始映射的情况下，创建一个包含转换后值的新映射。
func ApplyValue[K gotools.Comparable, V, U any](f func(K, V) U, m map[K]V) map[K]U {
	applied := make(map[K]U, len(m))

	for k, v := range m {
		applied[k] = f(k, v)
	}
	return applied

}

// ApplyKey 对给定映射中的每个键值对应用函数 f，并返回一个新的映射，
// 其中键是应用函数后的结果，值保持不变。
// 参数:
//
//	f: 一个函数，接受映射的键和值作为参数，返回一个新的键。
//	m: 待处理的映射。
//
// 返回值:
//
//	一个新的映射，其中键是应用函数 f 后的结果，值是从原始映射中对应的值。
//
// 使用场景:
//
//	当需要根据原始映射的键和值计算出新的键，并保持值不变时，可以使用这个函数。
func ApplyKey[K, U gotools.Comparable, V any](f func(K, V) U, m map[K]V) map[U]V {
	applied := make(map[U]V, len(m))

	for k, v := range m {
		applied[f(k, v)] = v
	}
	return applied
}

// ApplyBoth 接受一个函数和一个映射，应用函数到映射的每个键值对上，
// 并返回一个新的映射，其中新映射的键是原映射的键和值通过函数处理后的结果，
// 值是原映射的值通过函数处理后的结果。
// 参数:
//
//	f: 一个函数，它接受原映射的键和值作为参数，返回新的键和值。
//	m: 原映射，它的键值对将被应用函数处理。
//
// 返回值:
//
//	一个新的映射，其中键是原映射的键和值通过函数处理后的结果，值是原映射的值通过函数处理后的结果。
func ApplyBoth[K, U gotools.Comparable, V, S any](f func(K, V) (U, S), m map[K]V) map[U]S {
	applied := make(map[U]S, len(m))

	for k, v := range m {
		k1, v1 := f(k, v)
		applied[k1] = v1
	}
	return applied
}

// Apply 应用给定的函数 f 到 map m 的每个键值对上，返回应用函数后的结果集合。
// 参数:
//
//	f: 一个函数，它接受 map 的键和值作为参数，返回一个任意类型的值。
//	m: 一个 map，它的每个键值对都会被函数 f 处理。
//
// 返回值:
//
//	一个切片，包含对 map 中每个键值对应用函数 f 后的结果。
//
// Apply 的目的是提供一个通用的方式，来对 map 的所有元素执行某种操作，而无需直接修改原 map。
func Apply[K gotools.Comparable, V, U any](f func(K, V) U, m map[K]V) []U {

	applied := make([]U, 0, len(m))
	for k, v := range m {
		applied = append(applied, f(k, v))
	}

	return applied
}

// FromArray2 根据一个函数和一个数组创建一个映射。
// 它接受两个参数：f和value，其中f是一个函数，它接受一个S类型的参数，返回一个T类型的值。
// 该函数的目的是通过索引匹配将value数组的元素作为映射的键，通过调用f函数将value数组的元素作为对应的值。
// 参数:
//
//	f: 一个函数，它接受一个S类型的参数，返回一个T类型的值。
//	value: 用于映射的值的数组。
//
// 返回值:
//
//	map[T]S: 一个映射，其中键来自value数组，值来自f函数的返回值。
//

func FromArray2[V ~[]S, T gotools.Comparable, S any](f func(x S) T, value V) map[T]S {

	dict := make(map[T]S, len(value))

	for i := range value {

		dict[f(value[i])] = value[i]

	}

	return dict

}

// ToPairs 将一个映射(map)转换为一个由Pair组成的切片(slice)。
// 这个函数接受一个类型为[K]V的映射作为输入，其中K和V可以是任何可比较和任意类型的值。
// 函数返回一个由Pair[K, V]组成的切片，每个Pair包含原始映射中的一对键值对。
// 这个函数的目的是为了提供一种将map结构转换为更方便处理的切片形式的方法。
// 参数:
//
//	m: 输入的映射(map)，其中K是键类型，V是值类型。
//
// 返回值:
//
//	一个由Pair[K, V]组成的切片(slice)，每个Pair包含映射中的一对键值对。
func ToArrayPairs[K gotools.Comparable, V any](m map[K]V) []pair.Pair[K, V] {

	pairs := make([]pair.Pair[K, V], 0, len(m))

	for k, v := range m {
		pairs = append(pairs, pair.Pair[K, V]{First: k, Second: v})
	}

	return pairs
}

// FromPairs 根据提供的转换函数和数据数组，创建并返回一个映射。
// 这个函数接受一个函数参数 fun，该函数用于将数据数组中的元素转换为映射的键，
// 同时接受一个数据数组 data，该数组的元素类型是 Pair，其中 First 字段作为映射的键，
// Second 字段作为映射的值。函数返回一个映射，其中键是通过 fun 函数转换得到的，
// 值是对应 Pair 中的 Second 字段。
// 参数:
//
//	fun - 一个函数，用于将数据数组中的元素转换为映射的键。
//	data - 一个数组，其元素类型是 Pair，用于生成映射。
//
// 返回值:
//
//	一个映射，其中键是通过 fun 函数转换得到的，值是对应 Pair 中的 Second 字段。
func FromPairs[K ~[]pair.Pair[V, S], T gotools.Comparable, V, S any](fun func(V) T, data K) map[T]S {

	dict := make(map[T]S, len(data))

	for i := range data {
		dict[fun(data[i].First)] = data[i].Second
	}

	return dict

}

// FromArrayWithFun 根据提供的函数和两个数组，创建一个映射。
// 它接受一个函数，该函数应用于键数组的每个元素，并使用该函数的返回值作为映射的键，
// 使用值数组的相应元素作为映射的值。这允许快速查找与特定键相关联的值。
//
// 参数:
//
//	fun - 一个函数，它接受一个类型为S的参数，并返回一个可比较的类型T的结果。
//	key - 一个类型为K的数组，其元素将被用于生成映射的键。
//	value - 一个类型为V的数组，其元素将被用作生成映射的值。
//
// 返回值:
//
//	一个映射，其中键是通过应用fun函数到key数组的元素得到的，值是对应于键在value数组中的元素。
//
// 注意:
//
//	K和V的类型参数必须是数组类型，而T必须是可比较的类型，S可以是任何类型。
//	这个函数假设key和value数组的长度是相同的，以便在生成映射时保持键-值对的一致性。
func FromArrayWithFun[A ~[]K, B ~[]V, T gotools.Comparable, K, V any](fun func(K) T, key A, value B) map[T]V {

	dict := make(map[T]V, len(key))

	for i := range key {

		dict[fun(key[i])] = value[i]

	}

	return dict

}

// PopulateSeries 根据给定的键和值的序列，以及一个最大范围，创建并返回一个映射。
// 键和值的序列用于初始化映射，随后映射将被扩展到指定的最大范围。
// 参数:
//
//	key K: 键的序列，类型为切片。
//	value V: 值的序列，类型为切片，必须与键的序列长度相同。
//	max int: 映射的最大范围。
//
// 返回值:
//
//	map[T]S: 创建的映射，其中T是键的类型，S是值的类型。
//
// 注意：该函数的类型参数使用了通用约束，要求K为切片类型，V也为切片类型，T和S分别是切片元素的类型。
func PopulateSeries[K ~[]T, V ~[]S, T gotools.Number, S any](key K, value V, max int) map[T]S {

	var v S

	dict := make(map[T]S, len(key))

	for i := range key {

		dict[key[i]] = value[i]

	}

	start := int(key[0])
	for i := start; i < start+max; i++ {

		if _, ok := dict[T(i)]; ok {
			continue
		}

		dict[T(i)] = v

	}

	return dict

}

// Contains 检查一个map是否包含指定的一组键，并返回这些键对应的值。
// 如果map不包含某个键，则对应值为该类型的零值。
//
// 参数:
//
//	m: 要检查的map。
//	key: 一个可变长参数，包含要检查的键。
//
// 返回值:
//
//	一个切片，包含所有指定键在map中的对应值。
//	如果某个键不存在于map中，则对应位置的值为提供的默认值。
func Contains[M map[K]V, A []V, K gotools.Comparable, V any](m M, default_value V, key ...K) A {

	lk := len(key)
	res := make(A, 0, lk)

	if lk == 0 {
		return res
	}

	for i := range key {

		if v, ok := m[key[i]]; !ok {
			res = append(res, default_value)
			continue
		} else {
			res = append(res, v)
		}

	}

	return res

}

// Remove 从给定的映射中移除指定的键。
// 该函数接受一个映射和一个或多个键作为参数，然后从映射中移除这些键及其对应的值。
// 这样做的目的是为了减少映射的大小或删除不再需要的条目。
// 参数:
//
//	M: 映射的类型，它定义了键和值的类型。
//	K: 映射中键的类型，必须是可比较的。
//	V: 映射中值的类型。
//	m: 要操作的映射。
//	key: 要移除的键的切片。可以同时移除多个键。
func Remove[M map[K]V, K gotools.Comparable, V any](m M, key ...K) {

	for i := range key {

		delete(m, key[i])

	}

}

// Concat 合并多个映射表为一个。
// 该函数接受一个或多个类型为[K]V的映射表作为参数，返回一个合并后的映射表。
// 参数 m1 是一个变长参数，表示要合并的映射表序列。
// 返回值是一个映射表，包含了所有输入映射表的键值对。
// 如果没有提供任何映射表作为参数，函数将返回 空map
func Concat[K gotools.Comparable, V any](m1 ...map[K]V) map[K]V {

	if len(m1) == 0 {
		return make(map[K]V)
	}

	num := array.Map(func(x ...map[K]V) int { return len(x[0]) }, m1)

	res := make(map[K]V, array.Sum(num))

	for _, value := range m1 {

		for k, v := range value {
			res[k] = v
		}
	}

	return res
}

// Exists 检查给定的映射中是否存在满足特定条件的键值对。
// 参数 f 是一个函数，用于定义条件，它接受映射的键和值作为参数，并返回一个布尔值。
// 参数 m 是要检查的映射。
// 返回值表示是否找到满足条件的键值对。
// Exists 的泛型设计允许它适用于任何类型的映射，只要键是可比较的。
func Exists[K gotools.Comparable, V any](f func(K, V) bool, m map[K]V) bool {

	for k, v := range m {
		if f(k, v) {
			return true
		}
	}
	return false
}

// All 对给定映射中的所有键值对应用函数 f，并检查是否所有应用都返回 true。
// 如果函数 f 对任何键值对返回 false，则立即返回 false，表示不是所有应用都成功。
// 参数:
//
//	f: 一个函数，接受映射的键和值作为参数，并返回一个布尔值。
//	    这个函数用于对每个键值对进行某种测试或处理。
//	m: 要处理的映射，其键类型和值类型分别为 K 和 V。
//
// 返回值:
//
//	如果函数 f 对映射中的所有键值对都返回 true，则返回 true；
//	如果函数 f 对任何键值对返回 false，则返回 false。
func All[K gotools.Comparable, V any](f func(K, V) bool, m map[K]V) bool {
	for k, v := range m {
		if !f(k, v) {
			return false
		}
	}
	return true
}

// ToPairsArray2 将一个带有多个值的映射表转换为键与值的二维数组对。
// 此函数接收一个映射表 m，其中键 K 是可比较的任意类型，值 S 是 V 类型元素的切片。
// 函数返回一个 Pair 结构体实例，其中 First 是 K 类型的切片，Second 是 V 类型的切片，
// 这些切片包含了原始映射表中所有的键和值，按照映射表中的键值对顺序排列。
//
// 参数:
//
//	m: 输入的映射表，键为 K 类型，值为 V 类型元素的切片 S。
//
// 返回:
//
//	Pair[[]K, []V]: 包含两个切片的 Pair 结构体，First 切片包含映射表的所有键，Second 切片包含映射表的所有值。
func ToPairsArray2[K gotools.Comparable, S ~[]V, V any](m map[K]S) pair.Pair[[]K, []V] {

	pairs := pair.Pair[[]K, []V]{}

	num := 0
	for _, v := range m {
		num += len(v)
	}

	pairs.First = make([]K, 0, num)
	pairs.Second = make([]V, 0, num)

	for k, v := range m {
		pairs.First = append(pairs.First, array.Repeat(k, len(v))...)
		pairs.Second = append(pairs.Second, v...)
	}

	return pairs

}

// ToPairsArray 将一个映射表转换为键与值的数组对。
// 此函数接收一个映射表 m，其中键 K 是可比较的任意类型，值 V 也是任意类型。
// 函数返回一个 Pair 结构体实例，其中 First 是 K 类型的切片，Second 是 V 类型的切片，
// 这些切片包含了原始映射表中所有的键和值，按照映射表中的键值对顺序排列。
//
// 参数:
//
//	m: 输入的映射表，键为 K 类型，值为 V 类型。
//
// 返回:
//
//	Pair[[]K, []V]: 包含两个切片的 Pair 结构体，First 切片包含映射表的所有键，Second 切片包含映射表的所有值。
func ToPairsArray[K gotools.Comparable, V any](m map[K]V) pair.Pair[[]K, []V] {

	pairs := pair.Pair[[]K, []V]{}

	pairs.First = make([]K, 0, len(m))
	pairs.Second = make([]V, 0, len(m))

	for k, v := range m {
		pairs.First = append(pairs.First, k)
		pairs.Second = append(pairs.Second, v)
	}

	return pairs

}

// Merge 合并多个映射表为一个。
// K: 映射表的键的类型，必须是可比较的。
// V: 映射表的值的类型。
// 参数 data 是一个变长参数，包含多个映射表。
// 返回值是一个新的映射表，包含所有输入映射表的键值对，相同的键只会出现一次。
func Merge[K gotools.Comparable, V any](data ...map[K]V) map[K]V {

	if len(data) == 0 {
		return make(map[K]V)
	}

	num := array.Unique(Keys(data...))

	res := make(map[K]V, len(num))

	for k, v := range data[0] {
		res[k] = v
	}

	for _, m := range data[1:] {
		for k, v := range m {
			if _, ok := res[k]; !ok {
				res[k] = v
			}
		}
	}
	return res
}

// Intersect 接受多个 map[K]V 类型的参数，返回一个 Pair[[][]K, []V] 类型的值。
// 该函数遍历输入映射的键的笛卡尔积。
// 对于每个键的组合，它创建一个空的 [][]U 类型的临时切片。
// 然后，函数从每个映射中获取相应的值，并将它们追加到临时切片中。
// 如果交集非空，则函数将键的组合和交集追加到结果的 First 和 Second 字段中。
// 如果没有找到交集，函数返回一个空的 Pair。
//
// Parameters:
//   - data: 输入映射的可变参数，类型为 map[K]V。
//
// Returns:
//   - res: 返回一个 Pair[[][]K, []V] 类型的值，
//     First 字段包含键的组合，Second 字段包含交集。
//
// Example:
//
//	data1 := map[string][]int{"a": {1, 2, 3}, "b": {5, 6, 7, 8, 9}, "c": {4}}
//	data2 := map[string][]int{"a": {3, 4, 5}, "b": {7}, "c": {8, 9, 10}}
//
//	result := MapIntersect(data1, data2)
//
//	// Output:
//	// res.First = [
//	//   ["a" "a"],
//	//   ["b" "a"],
//	//   ["c" "a"],
//	//   ["b" "b"],
//	//   ["b" "c"],
//	// ]
//	// res.Second = [
//	//   [3],
//	//   [5],
//	//   [4],
//	//   [7],
//	//   [8, 9],
//	// ]
func Intersect[V []U, K, U gotools.Comparable](data ...map[K]V) pair.Pair[[][]K, []V] {

	if len(data) == 0 {
		return pair.Pair[[][]K, []V]{}
	}

	res := pair.Pair[[][]K, []V]{}
	rf := array.Cartesian(array.Map(func(x ...map[K]V) []K { return Keys(x[0]) }, data)...)

	for i := range rf {
		tmp := make([][]U, 0)
		for k, v := range rf[i] {
			tmp = append(tmp, data[k][v])
		}

		if dd := array.Intersect(tmp...); len(dd) > 0 {
			res.First = append(res.First, rf[i])
			res.Second = append(res.Second, dd)
		}

	}

	return res

}

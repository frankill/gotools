package array

import (
	"cmp"
)

// MatchZero 执行精确匹配，在 'lookup_array' 中查找与 'lookup_value' 相同的元素。
// 参数:
//
//	lookup_value - 类型为 A 的切片，其中包含要匹配的值。
//	lookup_array - 类型为 A 的切片，其中包含可能的匹配项。
//
// 返回:
//
//	一个整数切片，表示在 'lookup_array' 中找到的每个 'lookup_value' 元素的第一个匹配项的索引，如果没有找到匹配项则返回 -1。
func MatchZero[A ~[]T, T comparable](lookup_value A, lookup_array A) []int {

	return ArrayMap(func(x ...T) int {
		return ArrayFindFirst(func(y ...T) bool {
			return y[0] == x[0]
		}, lookup_array)
	}, lookup_value)

}

// MatchOne 执行匹配，对于 'lookup_value' 中的每个元素，在排序后的 'lookup_array' 中找到第一个大于等于它的元素。
// 参数:
//
//	lookup_value - 类型为 A 的切片，其中包含要匹配的值。
//	lookup_array - 类型为 A 的切片，其中包含可能的匹配项。
//
// 返回:
//
//	一个整数切片，表示在 'lookup_array' 中找到的每个 'lookup_value' 元素的匹配项的索引，如果没有找到匹配项则返回 -1。
func MatchOne[A ~[]T, T cmp.Ordered](lookup_value A, lookup_array A) []int {

	ll := len(lookup_value)
	la := len(lookup_array)

	id := ArraySeq(0, la, 1)
	res := Rep(-1, ll)

	id, lookup_array = ArraySortBy(func(x, y T) bool {
		return x < y
	}, id, lookup_array)

	for i := 0; i < ll; i++ {

		for j := 0; j < la; j++ {

			if lookup_value[i] <= lookup_array[j] {

				res[i] = id[j]

				break

			}

		}

	}

	return res
}

// MatchMinusOne 执行匹配，对于 'lookup_value' 中的每个元素，在逆序排序后的 'lookup_array' 中找到第一个小于等于它的元素。
// 参数:
//
//	lookup_value - 类型为 A 的切片，其中包含要匹配的值。
//	lookup_array - 类型为 A 的切片，其中包含可能的匹配项。
//
// 返回:
//
//	一个整数切片，表示在 'lookup_array' 中找到的每个 'lookup_value' 元素的匹配项的索引，如果没有找到匹配项则返回 -1。
func MatchMinusOne[A ~[]T, T cmp.Ordered](lookup_value A, lookup_array A) []int {

	ll := len(lookup_value)
	la := len(lookup_array)

	id := ArraySeq(0, la, 1)
	res := Rep(-1, ll)

	id, lookup_array = ArraySortBy(func(x, y T) bool {
		return x > y
	}, id, lookup_array)

	for i := 0; i < ll; i++ {

		for j := 0; j < la; j++ {

			if lookup_value[i] >= lookup_array[j] {

				res[i] = id[j]

				break

			}

		}

	}

	return res

}

// Xlookup 根据 lookup_value 在 lookup_array 中查找匹配项，并返回对应在 lookup_result 中的元素。
// 该函数利用泛型实现了类型安全的查找操作，适用于任何实现了 comparable 接口的类型。
// 参数:
//
//	lookup_value - 需要查找的值。
//	lookup_array - 用于查找的数组，必须包含与 lookup_value 相同类型的元素。
//	lookup_result - 用于返回结果的数组，类型为 U，它应该包含与 lookup_array 中元素相关联的值。
//
// 返回值:
//
//	返回一个数组，包含与 lookup_value 在 lookup_array 中找到的元素相对应的值。
//
// 使用泛型允许函数适用于多种类型，提高了代码的复用性和灵活性。
func Xlookup[S ~[]T, T comparable, R ~[]U, U any](lookup_value S, lookup_array S, lookup_result R) []U {

	index := MatchZero(lookup_value, lookup_array)

	return ArrayChoose(index, lookup_result)

}

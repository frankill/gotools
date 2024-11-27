package operation

import (
	"github.com/frankill/gotools"
	"github.com/frankill/gotools/array"
)

// ForEach 逐元素执行函数。
// 参数:
//
//	fun: 用于逐元素执行的函数。
//	arr: 变长参数列表，每个参数都是一个切片。
//
// 返回:
//
//	无返回值。
func ForEach[S ~[]T, T any](fun func(x ...T), arr ...S) {

	if len(arr) == 0 {
		return
	}

	la := array.Map(func(x S) int { return len(x) }, arr)
	lm := array.Max(la)

	parm := make([]T, len(arr))
	for i := 0; i < lm; i++ {
		for j := 0; j < len(arr); j++ {
			parm[j] = arr[j][i%la[j]]
		}
		fun(parm...)
	}

}

// Add 逐元素相加
// 参数:
//
//	a: 一个切片
//	b: 一个切片
//
// 返回:
//
//	一个切片
func Add[S ~[]T, T gotools.Number](a, b S) []T {
	return array.Map2(func(x, y T) T { return x + y }, a, b)
}

// Sub 逐元素相减
// 参数:
//
//	a: 一个切片
//	b: 一个切片
//
// 返回:
//
//	一个切片
func Sub[S ~[]T, T gotools.Number](a, b S) []T {
	return array.Map2(func(x, y T) T { return x - y }, a, b)
}

// Mul 逐元素相乘
// 参数:
//
//	a: 一个切片
//	b: 一个切片
//
// 返回:
//
//	一个切片
func Mul[S ~[]T, T gotools.Number](a, b S) []T {
	return array.Map2(func(x, y T) T { return x * y }, a, b)
}

// Div 逐元素相除
// 参数:
//
//	a: 一个切片
//	b: 一个切片
//
// 返回:
//
//	一个切片
func Div[S ~[]T, T gotools.Number, R float64](a, b S) []R {
	return array.Map2(func(x, y T) R { return R(x) / R(y) }, a, b)
}

// Mod 逐元素取模
// 参数:
//
//	a: 一个切片
//	b: 一个切片
//
// 返回:
//
//	一个切片
func Mod[S ~[]T, T gotools.Integer](a, b S) []int {
	return array.Map2(func(x, y T) int { return int(x) % int(y) }, a, b)
}

// Lte 逐元素比较大小
// 参数:
//
//	a: 一个切片
//	b: 一个切片
//
// 返回:
//
//	一个切片
func Lte[S ~[]T, T gotools.Ordered](a, b S) []bool {
	return array.Map2(func(x, y T) bool { return x <= y }, a, b)
}

// Gte 逐元素比较大小
// 参数:
//
//	a: 一个切片
//	b: 一个切片
//
// 返回:
//
//	一个切片
func Gte[S ~[]T, T gotools.Ordered](a, b S) []bool {
	return array.Map2(func(x, y T) bool { return x >= y }, a, b)
}

// Lt 逐元素比较
// 参数:
//
//	a: 一个切片
//	b: 一个切片
//
// 返回:
//
//	一个切片
func Lt[S ~[]T, T gotools.Ordered](a, b S) []bool {
	return array.Map2(func(x, y T) bool { return x < y }, a, b)
}

// Gt 逐元素比较大小
// 参数:
//
//	a: 一个切片
//	b: 一个切片
//
// 返回:
//
//	一个切片
func Gt[S ~[]T, T gotools.Ordered](a, b S) []bool {
	return array.Map2(func(x, y T) bool { return x > y }, a, b)
}

// Eq 逐元素比较
// 参数:
//
//	a: 一个切片
//	b: 一个切片
//
// 返回:
//
//	一个切片
func Eq[S ~[]T, T gotools.Comparable](a, b S) []bool {
	return array.Map2(func(x, y T) bool { return x == y }, a, b)
}

// Neq 逐元素比较是否不相等
// 参数:
//
//	a: 一个切片
//	b: 一个切片
//
// 返回:
//
//	一个切片
func Neq[S ~[]T, T gotools.Comparable](a, b S) []bool {
	return array.Map2(func(x, y T) bool { return x != y }, a, b)
}

// Not 逐元素取反
// 参数:
//
//	a: 一个切片
//
// 返回:
//
//	一个切片
func Not(a []bool) []bool {
	return array.Map(func(x bool) bool { return !x }, a)
}

// Map 对一组切片应用指定的函数，每个切片元素按位置组合后作为函数的参数。
// 此泛型函数接受一个变长参数列表 `arr`，其中每个参数都是类型 S 的切片（S 必须是切片类型），
// 和一个函数 `fun`，该函数接受变长参数列表 x...T 并返回类型 U 的结果。
// 函数返回一个类型为 []U 的切片，其中包含应用给定函数 `fun` 到所有切片元素组合上的结果。
//
// 参数:
//
//   - fun: 一个函数，接受变长参数列表 x...T，并返回类型 U 的结果。
//   - arr: 变长参数列表，每个参数都是类型 S 的切片，S 必须是切片类型。
//
// 返回:
//
//   - 类型为 []U 的切片，包含应用给定函数 `fun` 到所有切片元素组合上的结果。
func Map[S ~[]T, T, U any](fun func(x ...T) U, arr ...S) []U {

	if len(arr) == 0 {
		return make([]U, 0)
	}

	la := array.Map(func(x S) int { return len(x) }, arr)
	lm := array.Max(la)

	res := make([]U, lm)

	parm := make([]T, len(arr))
	for i := 0; i < lm; i++ {
		for j := 0; j < len(arr); j++ {
			parm[j] = arr[j][i%la[j]]
		}
		res[i] = fun(parm...)
	}

	return res
}

func Reduce[S ~[]T, T, U any](fun func(x U, y ...T) U, init U, arr ...S) U {

	result := init

	if len(arr) == 0 {
		return result
	}

	la := array.Map(func(x S) int { return len(x) }, arr)
	lm := array.Max(la)

	parm := make([]T, len(arr))

	for i := 0; i < lm; i++ {
		for j := 0; j < len(arr); j++ {
			parm[j] = arr[j][i%la[j]]
		}
		result = fun(result, parm...)
	}
	return result
}

func ReduceR[S ~[]T, T, U any](fun func(x U, y ...T) U, init U, arr ...S) U {

	result := init

	if len(arr) == 0 {
		return result
	}

	la := array.Map(func(x S) int { return len(x) }, arr)
	lm := array.Max(la)

	parm := make([]T, len(arr))

	for i := lm - 1; i >= 0; i-- {
		for j := 0; j < len(arr); j++ {
			parm[j] = arr[j][i%la[j]]
		}
		result = fun(result, parm...)
	}

	return result
}

// FindLast 查找类型为 S（元素类型为 T）的切片数组中最后一个使条件函数 `fun` 返回 true 的元素组合所在的索引位置。
//
// 参数:
//   - fun: 一个函数，接受 T 类型的变长参数并返回布尔值，用于测试一组元素是否满足条件。
//   - arr: 变长参数，每个元素为类型为 S 的切片（T 类型的切片） ,切片不一致会循环补齐
//
// 返回值:
//   - 第一个满足条件的元素组合在原数组中的起始索引位置。
//     如果没有找到满足条件的组合，则返回   - 1。
//     如果输入切片数组为空，则直接返回   - 1。
func FindLast[S ~[]T, T any](fun func(x ...T) bool, arr ...S) int {

	result := -1

	if len(arr) == 0 {
		return result
	}

	la := array.Map(func(x S) int { return len(x) }, arr)
	lm := array.Max(la)

	parm := make([]T, len(arr))
	for i := 0; i < lm; i++ {
		for j := 0; j < len(arr); j++ {
			parm[j] = arr[j][i%la[j]]
		}

		if fun(parm...) {
			result = i
		}
	}

	return result
}

// FindFirst 查找类型为 S（元素类型为 T）的切片数组中第一个使条件函数 `fun` 返回 true 的元素组合所在的索引位置。
//
// 参数:
//   - fun: 一个函数，接受 T 类型的变长参数并返回布尔值，用于测试一组元素是否满足条件。
//   - arr: 变长参数，每个元素为类型为 S 的切片（T 类型的切片） ,切片不一致会循环补齐
//
// 返回值:
//   - 第一个满足条件的元素组合在原数组中的起始索引位置。
//     如果没有找到满足条件的组合，则返回   - 1。
//     如果输入切片数组为空，则直接返回   - 1。
func FindFirst[S ~[]T, T any](fun func(x ...T) bool, arr ...S) int {

	result := -1

	if len(arr) == 0 {
		return result
	}

	la := array.Map(func(x S) int { return len(x) }, arr)
	lm := array.Max(la)

	parm := make([]T, len(arr))
	for i := 0; i < lm; i++ {
		for j := 0; j < len(arr); j++ {
			parm[j] = arr[j][i%la[j]]
		}

		if fun(parm...) {
			result = i
			break
		}
	}

	return result
}

// Last 查找类型为 S（元素类型为 T）的切片数组中最后一个使条件函数 `fun` 返回 true 的元素组合，并返回该组合的第一个元素。
//
// 参数:
//   - fun: 一个函数，接受 T 类型的变长参数并返回布尔值，用于测试一组元素是否满足条件。
//   - arr: 变长参数，每个元素为类型为 S 的切片（T 类型的切片），切片不一致会循环补齐
//
// 返回值:
//   - 最后一个满足条件的元素组合中的第一个元素。
//     如果没有找到满足条件的组合，则返回 T 类型的零值。
//     如果输入切片数组为空，则直接返回 T 类型的零值。
func Last[S ~[]T, T any](fun func(x ...T) bool, arr ...S) T {

	var result T

	if len(arr) == 0 {
		return result
	}

	la := array.Map(func(x S) int { return len(x) }, arr)
	lm := array.Max(la)

	parm := make([]T, len(arr))
	for i := 0; i < lm; i++ {
		for j := 0; j < len(arr); j++ {
			parm[j] = arr[j][i%la[j]]
		}

		if fun(parm...) {
			result = parm[0]
		}
	}

	return result
}

// First 查找类型为 S（元素类型为 T）的切片数组中第一个使条件函数 `fun` 返回 true 的元素组合，并返回该组合的第一个元素。
//
// 参数:
//   - fun: 一个函数，接受 T 类型的变长参数并返回布尔值，用于测试一组元素是否满足条件。
//   - arr: 变长参数，每个元素为类型为 S 的切片（T 类型的切片），所有切片长度需一致。
//
// 返回值:
//   - 第一个满足条件的元素组合中的第一个元素。
//     如果没有找到满足条件的组合，则返回 T 类型的零值。
//     如果输入切片数组为空，则直接返回 T 类型的零值。
//
// 注意:
//   - 所有输入切片的长度必须相等，否则函数的行为未定义。
func First[S ~[]T, T any](fun func(x ...T) bool, arr ...S) T {

	var result T

	if len(arr) == 0 {
		return result
	}

	la := array.Map(func(x S) int { return len(x) }, arr)
	lm := array.Max(la)

	parm := make([]T, len(arr))
	for i := 0; i < lm; i++ {
		for j := 0; j < len(arr); j++ {
			parm[j] = arr[j][i%la[j]]
		}

		if fun(parm...) {
			result = parm[0]
			break
		}
	}

	return result
}

// All 检查类型为 S（元素类型为 T）的切片数组中所有元素组合是否都满足提供的条件函数。
//
// 参数:
//   - fun: 一个函数，接受 T 类型的变长参数并返回布尔值，用于测试一组元素是否满足条件。
//   - arr: 变长参数，每个元素为类型为 S 的切片（T 类型的切片） ，切片长度不一致会循环补齐
//
// 返回值:
//   - 如果所有元素组合均使得 `fun` 返回 true，则返回 true；只要有一个不满足则返回 false。
//     如果输入切片数组为空，则直接返回 false
func All[S ~[]T, T any](fun func(x ...T) bool, arr ...S) bool {

	if len(arr) == 0 {
		return false
	}

	la := array.Map(func(x S) int { return len(x) }, arr)
	lm := array.Max(la)

	parm := make([]T, len(arr))
	for i := 0; i < lm; i++ {
		for j := 0; j < len(arr); j++ {
			parm[j] = arr[j][i%la[j]]
		}

		if !fun(parm...) {
			return false
		}
	}

	return true

}

// Any 检查类型为 S（元素类型为 T）的切片数组中是否有任一元素组合满足提供的条件函数。
//
// 参数:
//   - fun: 一个函数，接受 T 类型的变长参数并返回布尔值，用于测试一组元素是否满足条件。
//   - arr: 变长参数，每个元素为类型为 S 的切片（T 类型的切片），长度不一致的切片会被循环补齐。
//
// 返回值:
//   - 如果至少有一个元素组合使得 `fun` 返回 true，则返回 true；否则返回 false。
//     如果输入切片数组为空，则直接返回 false。
func Any[S ~[]T, T any](fun func(x ...T) bool, arr ...S) bool {

	if len(arr) == 0 {
		return false
	}

	la := array.Map(func(x S) int { return len(x) }, arr)
	lm := array.Max(la)

	parm := make([]T, len(arr))
	for i := 0; i < lm; i++ {
		for j := 0; j < len(arr); j++ {
			parm[j] = arr[j][i%la[j]]
		}

		if fun(parm...) {
			return true
		}
	}

	return false

}

// ReverseSplit 根据提供的条件函数反向地将输入切片 S（类型为 T）分割成多个子切片，并返回这些子切片组成的切片。
// 与 `ArraySplit` 不同，此函数在条件满足的位置进行切割，并且包含切割点的元素在下一个子切片中。
//
// 参数:
//   - fun: 一个函数，接受 T 类型的变长参数并返回布尔值，决定是否在当前位置进行切分。
//   - arr: 变长参数，每个元素为类型为 S 的切片（T 类型的切片），长度不一致的切片会被循环补齐
//
// 返回值:
//   - 一个由 S 类型子切片组成的切片，每个子切片代表原切片中满足分割条件的相邻部分。
//     区别在于，当条件满足时，该元素会包含在后续的子切片中，而非当前子切片的结尾。
//     若输入为空或首切片为空，则返回空 S 类型切片的切片
//
// 注意:
//   - 数组将在元素的右侧进行拆分。
func ReverseSplit[S ~[]T, T any](fun func(x ...T) bool, arr ...S) [][]T {

	if len(arr) == 0 || len(arr[0]) == 0 {
		return [][]T{}
	}

	l := len(arr[0])
	f := len(arr)
	la := array.Map(func(x S) int { return len(x) }, arr)
	param := make([]T, f)
	result := make([][]T, 0)

	num := 0

	for i := 0; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i%la[j]]
		}
		if fun(param...) {
			result = append(result, arr[0][num:i+1])
			num = i + 1
		}
	}

	if num < l && num >= 0 {
		result = append(result, arr[0][num:])
	}

	return result
}

// Split 根据提供的条件函数将输入切片 S（类型为 T）分割成多个子切片，并返回这些子切片组成的切片。
// 条件函数 `fun` 应用于输入切片的每个元素，当 `fun` 返回 true 时，会在该位置切割切片。
//
// 参数:
//   - fun: 一个函数，接受 T 类型的变长参数并返回布尔值，决定是否在当前位置进行切分。
//   - arr: 变长参数，每个元素为类型为 S 的切片（T 类型的切片）， 长度不一致的切片会被循环补齐
//
// 返回值:
//   - 一个由 S 类型子切片组成的切片，每个子切片代表原切片中满足分割条件相邻的部分。
//     若输入为空或首切片为空，则返回空 S 类型切片的切片。
//
// 注意:
//   - 数组将在元素的左侧进行拆分。
//   - 数组不会在第一个元素之前被分割。
func Split[S ~[]T, T any](fun func(x ...T) bool, arr ...S) [][]T {
	if len(arr) == 0 || len(arr[0]) == 0 {
		return [][]T{}
	}

	l := len(arr[0])
	f := len(arr)
	la := array.Map(func(x S) int { return len(x) }, arr)
	param := make([]T, f)
	result := make([][]T, 0)

	num := 0

	for i := 1; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i%la[j]]
		}
		if fun(param...) {
			result = append(result, arr[0][num:i])
			num = i
		}
	}

	if num < l && num >= 0 {
		result = append(result, arr[0][num:])
	}

	return result
}

// ReverseFill 根据提供的条件函数反向填充新切片。它从最后一个元素开始向前遍历，
// 对于每个索引位置，使用条件函数 `fun` 应用于对应位置的元素，决定该位置的值。
// 如果 `fun` 返回 false，则新切片中的该位置元素取自后一个索引的值（即更靠近末尾的值）；
// 如果 `fun` 返回 true，则取自当前索引的值。
//
// 参数:
//   - fun: 一个函数，接受 T 类型的变长参数并返回布尔值，决定是否采用当前或下一个索引的值。
//   - arr: 变长参数，每个元素为类型为 S 的切片（T 类型的切片）， 长度不一致的切片会被循环补齐
//
// 返回值:
//   - 一个新的切片 S，其中元素根据 `fun` 的判断结果从前一个或当前索引的值中选取。
//     若输入为空或首切片为空，则返回空 T 类型切片。
func ReverseFill[S ~[]T, T any](fun func(x ...T) bool, arr ...S) []T {

	if len(arr) == 0 || len(arr[0]) == 0 {
		return []T{}
	}

	if len(arr[0]) == 1 {
		return append(S{}, arr[0][0])
	}

	l := len(arr[0])
	f := len(arr)
	la := array.Map(func(x S) int { return len(x) }, arr)
	param := make([]T, f)
	result := make([]T, l)

	result[l-1] = arr[0][l-1]

	for i := l - 2; i >= 0; i-- {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i%la[j]]
		}
		if !fun(param...) {
			result[i] = result[i+1]
		} else {
			result[i] = arr[0][i]
		}
	}

	return result
}

// Fill 根据提供的条件函数填充新切片。对于每个索引位置，如果条件函数应用于对应位置的元素返回 false，
// 则新切片中的该位置元素取自前一个索引位置的首个切片的元素；否则，取自当前索引位置的首个切片的元素。
//
// 参数:
//   - fun: 一个函数，接受 T 类型的变长参数并返回布尔值，决定是否采用当前索引的值。
//   - arr: 变长参数，每个元素为类型为 S 的切片（T 类型的切片）， 长度不一致的切片会被循环补齐
//
// 返回值:
//   - 一个新的切片 S，其中元素根据 `fun` 的判断结果从输入切片的相应位置或前一位置选取。
//     若输入为空或首切片为空，则返回空 T 类型切片。
func Fill[S ~[]T, T any](fun func(x ...T) bool, arr ...S) []T {
	if len(arr) == 0 || len(arr[0]) == 0 {
		return []T{}
	}

	if len(arr[0]) == 1 {
		return append(S{}, arr[0][0])
	}

	l := len(arr[0])
	f := len(arr)
	la := array.Map(func(x S) int { return len(x) }, arr)
	param := make([]T, f)
	result := make([]T, l)

	result[0] = arr[0][0]

	for i := 1; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i%la[j]]
		}
		if !fun(param...) {
			result[i] = result[i-1]
		} else {
			result[i] = arr[0][i]
		}
	}
	return result
}

// Filter 根据提供的函数过滤多个同结构切片（类型为 T 的切片）的元素。
// 它将每个切片的对应元素作为参数传递给函数 `fun`，并仅保留 `fun` 返回真值时的首个切片中的元素。
//
// 参数:
//   - fun: 一个函数，接受 T 类型的变长参数，返回布尔值，指示是否保留当前元素。
//   - arr: 变长参数，每个元素为 T 类型的切片, 长度不一致的切片会被循环补齐
//
// 返回值:
//   - 一个新的切片 S，包含根据 `fun` 筛选后的元素。若输入为空或首切片为空，则返回空切片 S。
//
// 。
func Filter[S ~[]T, T any](fun func(x ...T) bool, arr ...S) []T {
	if len(arr) == 0 || len(arr[0]) == 0 {
		return []T{}
	}
	l := len(arr[0])
	f := len(arr)
	la := array.Map(func(x S) int { return len(x) }, arr)
	param := make([]T, f)
	result := make([]T, 0)

	for i := 0; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i%la[j]]
		}
		if fun(param...) {
			result = append(result, arr[0][i])
		}
	}
	return result
}

// FilterIndex 根据提供的函数过滤多个同结构切片（类型为 T 的切片）的元素。
// 它将每个切片的对应元素作为参数传递给函数 `fun`，并仅保留 `fun` 返回真值时的位置索引。
//
// 参数:
//   - fun: 一个函数，接受 T 类型的变长参数，返回布尔值，指示是否保留当前元素。
//   - arr: 变长参数，每个元素为 T 类型的切片, 长度不一致的切片会被循环补齐
//
// 返回值:
//   - 一个新的切片 S，包含根据 `fun` 筛选后的元素的索引。若输入为空或首切片为空，则返回空切片 []int。
func FilterIndex[S ~[]T, T any](fun func(x ...T) bool, arr ...S) []int {
	if len(arr) == 0 || len(arr[0]) == 0 {
		return []int{}
	}
	l := len(arr[0])
	f := len(arr)
	la := array.Map(func(x S) int { return len(x) }, arr)
	param := make([]T, f)
	result := make([]int, 0, l)

	for i := 0; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i%la[j]]
		}
		if fun(param...) {
			result = append(result, i)
		}
	}
	return result
}

// FlatMap 对多个同结构切片（S，类型为 T 的切片）应用一个函数，将每个切片的对应元素作为参数传递：
// 并收集返回值形成一个新的 U 类型切片序列。
//
// 参数:
//   - fun: 一个函数，接受 T 类型的变长参数，并返回 U 类型的结果。
//   - arr: 变长参数，每个元素为类型为 S 的切片， 长度不一致的切片会被循环补齐
//
// 返回值:
//   - 一个 U 类型的切片，其元素为对输入切片每相同索引位置的元素应用 `fun` 函数后的结果。
//     若输入为空或首切片为空，则返回空切片。
func FlatMap[S ~[]T, T any, U any](fun func(x ...T) []U, arr ...S) []U {

	if len(arr) == 0 || len(arr[0]) == 0 {
		return []U{}
	}
	l := len(arr[0])
	f := len(arr)
	la := array.Map(func(x S) int { return len(x) }, arr)
	param := make([]T, f)
	result := make([]U, 0)

	for i := 0; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i%la[j]]
		}
		result = append(result, fun(param...)...)
	}
	return result
}

// Fold 对类型为 S 的多维数组（其中 S 是 T 类型的切片）执行折叠操作，生成一个 U 类型的结果切片。
// 它接收三个参数：
//   - fun：一个函数，接受变长参数 T 的切片并返回一个 U 类型的值，用于单个维度的聚合操作。
//   - acc：一个累积函数，接受两个 U 类型的参数并返回一个 U 类型的值，用于将相邻结果累积。
//   - arr：变长参数，表示待折叠的多维数组，数组的每个元素也是类型为 S 的切片, 长度不一致的切片会被循环补齐
//
// 返回值：
//   - 一个 U 类型的切片，表示经过聚合和累积操作后的结果序列。
//
// 示例用途：
// 可以用来计算多维数组中各维度对应位置的元素经过特定运算后的序列，如多序列的逐元素加法、乘法等。
func Fold[S ~[]T, T, U any](fun func(x ...T) U, acc func(x, y U) U, arr ...S) []U {
	if len(arr) == 0 || len(arr[0]) == 0 {
		return []U{}
	}

	l := len(arr[0])
	f := len(arr)
	la := array.Map(func(x S) int { return len(x) }, arr)
	param := make(S, f)
	result := make([]U, l)

	for i := 0; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i%la[j]]
		}
		result[i] = fun(param...)

		if i > 0 {
			result[i] = acc(result[i-1], result[i])
		}

	}

	return result
}

// Intersect 计算多个切片（类型为 []T，元素类型 T 可比较）的交集。
// 参数:
//   - fun: 一个函数，接受 T 类型的变长参数并返回 U 类型的结果，用于比较元素。
//   - arr: 变长参数，每个元素为类型为 []T 的切片
//
// 返回值:
//   - 一个 []T 类型的切片，其元素为所有输入切片的交集元素。顺序与输入切片元素第一次出现的顺序相同。
//     若输入为空，则返回空切片。
func Intersect[S ~[]T, T any, U gotools.Comparable](f func(x T) U, arr ...S) []T {

	if len(arr) == 0 {
		return make([]T, 0)
	}

	// 如果只有一个切片，则直接返回它
	if len(arr) == 1 {
		return arr[0]
	}

	// 使用第一个切片作为基数来收集交集元素
	intersectionMap := make(map[U]*counts, len(arr[0]))
	for index, item := range arr[0] {

		t := f(item)

		if v, ok := intersectionMap[t]; !ok {
			intersectionMap[t] = &counts{
				Key:   1,
				Count: 1,
				index: index,
			}
		} else {
			v.Count++
			v.Key++

		}
	}

	for k := range arr[1:] {

		for _, item := range arr[1:][k] {

			t := f(item)
			if v, ok := intersectionMap[t]; ok {
				v.Count++
			}
		}

	}

	res := make([]T, 0)

	for _, v := range intersectionMap {

		if v.Count > v.Key {
			for i := 0; i < v.Key; i++ {
				res = append(res, arr[0][v.index])
			}
		}
	}

	return res
}

// Subtract 计算多个切片（类型为 []T，元素类型 T 可比较）的差集。
//
// 参数:
//   - arr: 变长参数，每个参数为一个待求差集的切片。
//
// 返回值:
//   - 一个新的 []T 类型的切片，包含所有输入切片中不在其它输入切片中出现的元素，且元素顺序与它们在第一个切片中出现的顺序一致。
func Subtract[S ~[]T, T any, U gotools.Comparable](f func(x T) U, arr ...S) []T {

	if len(arr) == 0 {
		return make([]T, 0)
	}

	// 如果只有一个切片，则直接返回它
	if len(arr) == 1 {
		return arr[0]
	}

	// 使用第一个切片作为基数来收集交集元素
	intersectionMap := make(map[U]*counts, len(arr[0]))
	for index, item := range arr[0] {

		t := f(item)

		if v, ok := intersectionMap[t]; !ok {
			intersectionMap[t] = &counts{
				Key:   1,
				Count: 1,
				index: index,
			}
		} else {
			v.Count++
			v.Key++

		}
	}

	for k := range arr[1:] {

		for _, item := range arr[1:][k] {

			t := f(item)
			if v, ok := intersectionMap[t]; ok {
				v.Count++
			}
		}

	}

	res := make([]T, 0)

	for _, v := range intersectionMap {

		if v.Count == v.Key {
			for i := 0; i < v.Key; i++ {
				res = append(res, arr[0][v.index])
			}
		}
	}

	return res
}

// InterS 计算多个切片（类型为 []T，元素类型 T 可比较）的交集。
//
// 参数:
//   - arr: 变长参数，每个参数为一个待求交集的切片。
//
// 返回值:
//   - 一个新的 []T 类型的切片，包含所有输入切片中共有的元素，且元素顺序与它们在第一个切片中出现的顺序一致。
//     如果没有交集或输入为空，则返回一个空切片。
//
// 注意: T 必须实现 gotools.Comparable 接口，允许元素之间的比较操作。
func InterS[S ~[]T, T gotools.Comparable](arr ...S) []T {

	if len(arr) == 0 {
		return make([]T, 0)
	}

	// 如果只有一个切片，则直接返回它
	if len(arr) == 1 {
		return arr[0]
	}

	// 使用第一个切片作为基数来收集交集元素
	intersectionMap := make(map[T]*count, len(arr[0]))
	for _, item := range arr[0] {
		if v, ok := intersectionMap[item]; !ok {
			intersectionMap[item] = &count{
				Key:   1,
				Count: 1,
			}
		} else {
			v.Count++
			v.Key++

		}
	}

	for k := range arr[1:] {

		for _, item := range arr[1:][k] {
			if v, ok := intersectionMap[item]; ok {
				v.Count++
			}
		}

	}

	res := make([]T, 0)

	for k, v := range intersectionMap {

		if v.Count > v.Key {
			for i := 0; i < v.Key; i++ {
				res = append(res, k)
			}
		}
	}

	return res
}

// SubS 计算多个切片（类型为 []T，元素类型 T 可比较）的交集。
//
// 参数:
//   - arr: 变长参数，每个参数为一个待求交集的切片。
//
// 返回值:
//   - 一个新的 []T 类型的切片，包含所有输入切片中共有的元素，且元素顺序与它们在第一个切片中出现的顺序一致。
//     如果没有交集或输入为空，则返回一个空切片。
//
// 注意: T 必须实现 gotools.Comparable 接口，允许元素之间的比较操作。
func SubS[S ~[]T, T gotools.Comparable](arr ...S) []T {
	if len(arr) == 0 {
		return make([]T, 0)
	}

	// 如果只有一个切片，则直接返回它
	if len(arr) == 1 {
		return arr[0]
	}

	// 使用第一个切片作为基数来收集交集元素
	intersectionMap := make(map[T]*count, len(arr[0]))
	for _, item := range arr[0] {
		if v, ok := intersectionMap[item]; !ok {
			intersectionMap[item] = &count{
				Key:   1,
				Count: 1,
			}
		} else {
			v.Count++
			v.Key++

		}
	}

	for k := range arr[1:] {

		for _, item := range arr[1:][k] {
			if v, ok := intersectionMap[item]; ok {
				v.Count++
			}
		}

	}

	res := make([]T, 0)

	for k, v := range intersectionMap {

		if v.Count == v.Key {
			for i := 0; i < v.Key; i++ {
				res = append(res, k)
			}
		}
	}

	return res
}

type count struct {
	Key   int
	Count int
}

type counts struct {
	Key   int
	Count int
	index int
}

package array

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/frankill/gotools"
)

var (
	ASCInt  = func(x, y int) bool { return x < y }
	DESCInt = func(x, y int) bool { return x < y }
)

func ASCGeneric[T gotools.Ordered](x, y T) bool {
	return x < y
}

func DESCGeneric[T gotools.Ordered](x, y T) bool {
	return x > y
}

// FromAny函数接受任意类型的输入参数，并将它们转换为切片。
//
// 参数:
//
//	input: 可变长度参数，表示要转换为切片的输入值
//
// 返回:
//
//	一个包含输入值的切片
func FromAny[T any](input ...T) []T {
	return input
}

// Rep 生成一个新的切片，该切片由输入切片 `x` 的元素重复 `n` 次组成。
//
// 参数:
// - x: 类型为 `S` 的输入切片，包含类型为 `T` 的元素。
// - n: 每个 `x` 中的元素应重复的次数。
// - sort: 一个布尔标志，指示生成的切片是否应进行排序。
//
// 返回:
// - 一个类型为 `[]T` 的新切片，包含重复的 `x` 元素。
func Rep[S ~[]T, T any](x S, n int, sort bool) []T {
	// Create an empty slice with a capacity equal to `n` times the length of `x`.
	result := make([]T, 0, n*len(x))

	// If `sort` is set to `true`, repeat each element of `x` `n` times and append them to the result slice.
	if sort {
		for i := range x {
			for j := 0; j < n; j++ {
				result = append(result, x[i])
			}
		}
		return result
	}

	// If `sort` is set to `false`, repeat the entire `x` slice `n` times and append it to the result slice.
	for i := 0; i < n; i++ {
		result = append(result, x...)
	}

	return result
}

// Repeated 创建一个包含相同重复切片的切片，该切片的长度等于n*len(x)。
//
// 参数:
//
//	x: 类型为 S 的切片，将复制到结果切片中。
//	n: 指定结果切片的长度，每个元素将复制 `x` `n` 次。
//
// 返回值:
//
//	一个类型为 []T 的切片，长度为 n*len(x)，其中每个元素都是 `x` 的副本。
func Repeated[S ~[]T, T any](x S, n int) []T {
	result := make([]T, 0, n*len(x))
	for i := 0; i < n; i++ {
		result = append(result, x...)
	}
	return result
}

// Seq 创建一个整数数组，该数组基于起始值、结束值和步长。
// 它返回一个包含从start到end（但不包括end）的整数的切片，每个相邻元素之间的差值为step。
// 如果step为0，表示没有有效的步长，返回一个空切片。
// 如果step为正数且start大于end，或者step为负数且start小于end，也返回一个空切片，因为在这种情况下无法生成有效的序列。
// 参数:
//
//	start - 序列的起始整数值。
//	end - 序列的结束整数值，不包括在结果中。
//	step - 序列中相邻元素之间的步长。
//
// 返回值:
//
//	一个整数切片，包含从start到end（但不包括end）的整数，每个相邻元素之间的差值为step
func Seq(start, end, step int) []int {
	if step == 0 {
		return []int{}
	}

	// 新增逻辑：直接根据step的正负判断生成序列的方向，忽略原逻辑中对start和end大小关系的限制
	size := int(math.Abs(float64(end-start)))/int(math.Abs(float64(step))) + 1

	result := make([]int, 0, size)
	for i := start; (step > 0 && i < end) || (step < 0 && i > end); i += step {
		result = append(result, i)
	}
	return result
}

// RandomSample 从输入切片 S（类型为 T）中随机抽取指定数量的元素，返回一个新的切片 S。
//
// 参数:
// - input: 类型为 S 的切片，从中进行随机抽样。
// - samples: 需要抽取的样本数量，必须小于输入切片的长度。
//
// 返回值:
//   - 一个新的切片 S，包含从输入切片中随机选取的 `samples` 个元素。
//     如果 `samples` 大于等于输入切片的长度，则直接返回原切片。
//
// 此函数使用当前时间作为随机数生成器的种子，确保每次调用都能得到不同的随机结果。
func RandomSample[S ~[]T, T any](input S, samples int, replace bool) S {

	li := len(input)

	if !replace {
		if samples > li {
			return input
		}
	}

	source := rand.NewSource(time.Now().UnixNano())

	r := rand.New(source)

	result := make(S, 0, samples)
	mapping := make(map[int]struct{}, samples)

	for {
		if len(result) >= samples {
			break
		}

		j := r.Intn(li)

		if !replace {
			if _, ok := mapping[j]; ok {
				continue
			}
		}
		mapping[j] = struct{}{}
		result = append(result, input[j])

	}

	return result
}

// Shif 对输入的切片 S（类型为 T）进行循环移位操作。
//
// 参数:
// - arr: 类型为 S 的切片，将被移位。
// - n: 移位步数，正数表示向右移位，负数表示向左移位。
//
// 返回值:
//   - 新的 T 类型切片，为移位操作后得到的结果。
//     如果输入切片为空，则直接返回原切片。

func Shif[S ~[]T, T any](arr S, n int) []T {

	la := len(arr)

	if la == 0 {
		return arr
	}

	if la <= n {
		return make([]T, la)
	}

	ns := int(math.Abs(float64(n)))
	t := make([]T, ns)

	if n > 0 {
		return append(t, arr[0:la-2]...)
	}

	return append(arr[ns:], t...)

}

// Rotate 对输入的切片 S（类型为 T）进行循环旋转操作。
//
// 参数:
// - arr: 类型为 S 的切片，需要进行旋转操作。
// - n: 旋转步数，正数表示向右旋转，负数表示向左旋转。
//
// 返回值:
//   - 新的 T 类型切片，为旋转操作后得到的结果。
//     如果输入切片为空，则直接返回原切片。
func Rotate[S ~[]T, T any](arr S, n int) []T {

	la := len(arr)

	if la == 0 {
		return arr
	}

	n = n % la

	ns := int(math.Abs(float64(n)))

	if n > 0 {
		return append(arr[ns:], arr[0:ns]...)
	}

	return append(arr[la-2:], arr[0:la-2]...)
}

// Product 计算类型为 S（元素为 gotools.Number 类型）的切片的所有元素的乘积，并返回结果为 float64 类型。
//
// 参数:
// - arr: 类型为 S 的切片，其元素需要实现 gotools.Number 接口（通常包括 int、float32、float64 等数字类型）。
//
// 返回值:
//   - 所有切片元素相乘的结果，转换为 float64 类型返回。
//     如果切片为空，则返回 1.0，遵循数学中乘积的空集定义。
func Product[S ~[]T, T gotools.Number](arr S) float64 {
	var result float64 = 1.0
	for _, v := range arr {
		result *= float64(v)
	}
	return result
}

// CumFun 对类型为 S（元素为 gotools.Number 类型）的切片应用累计函数，并返回累计结果的新切片 S。
//
// 参数:
// - arr: 类型为 S 的切片，元素应支持加减乘除等操作（实现 gotools.Number 接口）。
// - fun: 一个二元函数，接受两个 T 类型的参数并返回一个 T 类型的结果，用于定义累积操作的方式（如加法、乘法等）。
//
// 返回值:
//   - 一个新的切片 S，其中每个元素是原切片元素依序应用 `fun` 函数累积计算的结果。
//     第一个元素与原切片的第一个元素相同，之后每个元素都是前一个累积值与下一个数组元素的函数运算结果。
//
// 示例用途:
// - 使用加法函数时，可计算累积和。
// - 使用乘法函数时，可计算累积积。
func CumFun[S ~[]T, T gotools.Number](fun func(T, T) T, arr S) S {
	la := len(arr)
	result := make(S, la)
	result[0] = arr[0]
	for i := 1; i < la; i++ {
		result[i] = fun(result[i-1], arr[i])
	}
	return result
}

// CumSum 计算类型为 S（元素为 gotools.Number 类型）的切片的累积和，并返回结果为同类型的新切片 S。
//
// 参数:
// - arr: 类型为 S 的切片，元素应为可以相加的数字（实现 gotools.Number 接口）。
//
// 返回值:
//   - 一个新的切片 S，其中每个元素是从原切片开始到当前位置的所有元素的和。
//     例如，给定切片 [1, 2, 3, 4]，返回的累积和切片将是 [1, 3, 6, 10]。
func CumSum[S ~[]T, T gotools.Number](arr S) S {
	return CumFun(func(a, b T) T { return a + b }, arr)
}

// CumDiff 计算类型为 S（元素类型为 T）的切片中元素的累积差分，并返回一个新的切片。
// 要求 T 类型实现 gotools.Number 接口，支持减法运算。
//
// 参数:
// - arr: 输入的切片 S，元素为可以进行减法运算的数值类型。
//
// 返回值:
//   - 返回一个新的 S 类型切片，其中第 i 个元素是原切片中从第 0 项到第 i 项的累积差分值。
//     即新切片的第 i 项等于原切片第 i 项减去第 i-1 项的结果，首项特殊处理（通常为原切片的首项）。
func CumDiff[S ~[]T, T gotools.Number](arr S) S {
	return CumFun(func(a, b T) T { return b - a }, arr)
}

// CumProd 计算类型为 S（元素类型为 T）的切片中累积乘积序列。
//
// 参数:
// - arr: 类型为 S 的切片，元素应为支持乘法运算的数字类型（实现 gotools.Number 接口）。
//
// 返回值:
// - 一个新的 S 类型切片，其中每个元素是原始切片中从开始到当前位置（含当前位置）的所有元素的累积乘积。
func CumProd[S ~[]T, T gotools.Number](arr S) S {
	return CumFun(func(a, b T) T { return a * b }, arr)
}

// CumMax 计算类型为 S（元素类型为 T）的切片中累积最大值序列。
//
// 参数:
// - arr: 类型为 S 的切片，元素应为数字类型（实现 gotools.Number 接口）。
//
// 返回值:
// - 一个新的 S 类型切片，其中每个元素是原始切片中从开始到当前位置（含当前位置）的所有元素的累积最大值。
func CumMax[S ~[]T, T gotools.Number](arr S) S {
	return CumFun(func(a, b T) T { return max(a, b) }, arr)
}

// CumMin 计算类型为 S（元素类型为 T）的切片中累积最小值序列。
//
// 参数:
// - arr: 类型为 S 的切片，元素应为数字类型（实现 gotools.Number 接口）。
//
// 返回值:
// - 一个新的 S 类型切片，其中每个元素是原始切片中从开始到当前位置（含当前位置）的所有元素的累积最小值。
func CumMin[S ~[]T, T gotools.Number](arr S) S {
	return CumFun(func(a, b T) T { return min(a, b) }, arr)
}

// Mean 计算类型为 S（元素类型为 T）的切片中所有元素的平均值。
//
// 参数:
// - arr: 类型为 S 的切片，元素应为可以相加和除法运算的数字类型（实现 gotools.Number 接口）。
//
// 返回值:
// - 切片中所有元素的平均值。
func Mean[S ~[]T, T gotools.Number](arr S) T {
	la := len(arr)

	if la == 0 {
		return T(0)
	}
	return Sum(arr) / T(la)
}

// Len 计算类型为 S（元素类型为 T）的多个切片中元素的总数。
// 参数:
// - arr: 类型为 S 的切片
//
// 返回值:
// - 切片中元素的总数。
func Len[S ~[]T, T any](arr ...S) int {

	if len(arr) == 0 {
		return 0
	}

	var num int

	for _, v := range arr {
		num += len(v)
	}
	return num

}

// Sum 计算类型为 S（元素类型为 T）的切片中所有元素的总和。
//
// 参数:
// - arr: 类型为 S 的切片，元素应为可以相加的数字类型（实现 gotools.Number 接口）。
//
// 返回值:
//   - 切片中所有元素的和。
//     如果切片为空，则返回 T 类型的零值（对于数值类型通常是 0）。
func Sum[S ~[]T, T gotools.Number](arr ...S) T {
	var result T
	for _, v := range arr {
		for _, vv := range v {
			result += vv
		}
	}
	return result
}

// Max 返回类型为 S（元素类型为 T）的切片中的最大元素。
//
// 参数:
// - arr: 类型为 S 的切片，元素需要可比较（实现 gotools.Ordered 接口）。
//
// 返回值:
//   - 切片中的最大元素。
//     如果切片为空，则返回 T 类型的默认值（这可能是未定义的行为，具体取决于 T 的类型）。
func Max[S ~[]T, T gotools.Ordered](arr ...S) T {
	var result T
	if len(arr) == 0 {
		return result
	}
	if len(arr[0]) == 0 {
		return result
	}

	result = arr[0][0]
	for _, v := range arr {
		for _, vv := range v {
			result = max(result, vv)
		}
	}
	return result
}

// Min 返回类型为 S（元素类型为 T）的切片中的最小元素。
//
// 参数:
// - arr: 类型为 S 的切片，元素需要可比较（实现 gotools.Ordered 接口）。
//
// 返回值:
//   - 切片中的最小元素。
//     如果切片为空，则返回 T 类型的默认值（这可能是未定义的行为，具体取决于 T 的类型）。
func Min[S ~[]T, T gotools.Ordered](arr ...S) T {
	var result T
	if len(arr) == 0 {
		return result
	}
	if len(arr[0]) == 0 {
		return result
	}

	result = arr[0][0]
	for _, v := range arr {
		for _, vv := range v {
			result = min(result, vv)
		}
	}
	return result
}

// FindMin 查找类型为 S（元素类型为 T）的切片中最小元素的索引位置。
//
// 参数:
// - arr: 类型为 S 的切片，元素必须是可比较的（实现 gotools.Ordered 接口）。
//
// 返回值:
// - 返回切片中最小元素的索引。如果切片为空，则行为未定义（可能返回 0，具体取决于编译器和运行环境）。
func FindMin[S ~[]T, T gotools.Ordered](arr S) int {

	if len(arr) == 0 {
		return -1
	}

	var index int
	var value T

	value = arr[0]
	for k, v := range arr {
		if v < value {
			value = v
			index = k
		}
	}

	return index
}

// FindMax 查找类型为 S（元素类型为 T）的切片中最大元素的索引位置。
//
// 参数:
// - arr: 类型为 S 的切片，元素必须是可比较的（实现 gotools.Ordered 接口）。
//
// 返回值:
// - 返回切片中最大元素的索引。如果切片为空，则行为未定义（可能返回 0，具体取决于编译器和运行环境）。
//
// 注意: 此函数假定切片非空，并且切片中的元素能够相互比较以确定大小关系。
func FindMax[S ~[]T, T gotools.Ordered](arr S) int {

	if len(arr) == 0 {
		return -1
	}

	var index int
	var value T

	// 初始化索引和最大值，默认为切片的第一个元素
	index = 0
	value = arr[0]

	// 遍历切片以寻找最大值及其索引
	for k, v := range arr {
		if v > value {
			value = v
			index = k
		}
	}

	return index
}

// FindLast 查找类型为 S（元素类型为 T）的切片数组中最后一个使条件函数 `fun` 返回 true 的元素组合所在的索引位置。
//
// 参数:
// - fun: 一个函数，接受 T 类型的变长参数并返回布尔值，用于测试一组元素是否满足条件。
// - arr: 变长参数，每个元素为类型为 S 的切片（T 类型的切片），所有切片长度需一致。
//
// 返回值:
//   - 最后一个满足条件的元素组合在原数组中的起始索引位置。
//     如果没有找到满足条件的组合，则返回 -1。
//     如果输入切片数组为空，则直接返回 -1。
//
// 注意:
// - 所有输入切片的长度必须相等，否则函数的行为未定义。
func FindLast[S ~[]T, T any](fun func(x ...T) bool, arr ...S) int {

	result := -1

	if len(arr) == 0 {
		return result
	}

	l := len(arr[0])
	f := len(arr)
	param := make([]T, f)

	for i := 0; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i]
		}
		if fun(param...) {
			result = i
		}
	}
	return result
}

// FindFirst 查找类型为 S（元素类型为 T）的切片数组中第一个使条件函数 `fun` 返回 true 的元素组合所在的索引位置。
//
// 参数:
// - fun: 一个函数，接受 T 类型的变长参数并返回布尔值，用于测试一组元素是否满足条件。
// - arr: 变长参数，每个元素为类型为 S 的切片（T 类型的切片），所有切片长度需一致。
//
// 返回值:
//   - 第一个满足条件的元素组合在原数组中的起始索引位置。
//     如果没有找到满足条件的组合，则返回 -1。
//     如果输入切片数组为空，则直接返回 -1。
//
// 注意:
// - 所有输入切片的长度必须相等，否则函数的行为未定义。
func FindFirst[S ~[]T, T any](fun func(x ...T) bool, arr ...S) int {

	result := -1

	if len(arr) == 0 {
		return result
	}

	l := len(arr[0])
	f := len(arr)
	param := make(S, f)

	for i := 0; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i]
		}
		if fun(param...) {
			result = i
			break
		}
	}
	return result
}

// Last 查找类型为 S（元素类型为 T）的切片数组中最后一个使条件函数 `fun` 返回 true 的元素组合，并返回该组合的第一个元素。
//
// 参数:
// - fun: 一个函数，接受 T 类型的变长参数并返回布尔值，用于测试一组元素是否满足条件。
// - arr: 变长参数，每个元素为类型为 S 的切片（T 类型的切片），所有切片长度需一致。
//
// 返回值:
//   - 最后一个满足条件的元素组合中的第一个元素。
//     如果没有找到满足条件的组合，则返回 T 类型的零值。
//     如果输入切片数组为空，则直接返回 T 类型的零值。
//
// 注意:
// - 所有输入切片的长度必须相等，否则函数的行为未定义。
func Last[S ~[]T, T any](fun func(x ...T) bool, arr ...S) T {

	var result T

	if len(arr) == 0 {
		return result
	}

	l := len(arr[0])
	f := len(arr)
	param := make(S, f)

	for i := 0; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i]
		}
		if fun(param...) {
			result = param[0]
		}
	}
	return result
}

// First 查找类型为 S（元素类型为 T）的切片数组中第一个使条件函数 `fun` 返回 true 的元素组合，并返回该组合的第一个元素。
//
// 参数:
// - fun: 一个函数，接受 T 类型的变长参数并返回布尔值，用于测试一组元素是否满足条件。
// - arr: 变长参数，每个元素为类型为 S 的切片（T 类型的切片），所有切片长度需一致。
//
// 返回值:
//   - 第一个满足条件的元素组合中的第一个元素。
//     如果没有找到满足条件的组合，则返回 T 类型的零值。
//     如果输入切片数组为空，则直接返回 T 类型的零值。
//
// 注意:
// - 所有输入切片的长度必须相等，否则函数的行为未定义。
func First[S ~[]T, T any](fun func(x ...T) bool, arr ...S) T {

	var result T

	if len(arr) == 0 {
		return result
	}

	l := len(arr[0])
	f := len(arr)
	param := make(S, f)

	for i := 0; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i]
		}
		if fun(param...) {
			result = param[0]
			break
		}
	}
	return result
}

// All 检查类型为 S（元素类型为 T）的切片数组中所有元素组合是否都满足提供的条件函数。
//
// 参数:
// - fun: 一个函数，接受 T 类型的变长参数并返回布尔值，用于测试一组元素是否满足条件。
// - arr: 变长参数，每个元素为类型为 S 的切片（T 类型的切片），所有切片长度需一致。
//
// 返回值:
//   - 如果所有元素组合均使得 `fun` 返回 true，则返回 true；只要有一个不满足则返回 false。
//     如果输入切片数组为空，则直接返回 false
//
// 注意:
// - 所有输入切片的长度必须相等，否则函数的行为未定义。
func All[S ~[]T, T any](fun func(x ...T) bool, arr ...S) bool {

	if len(arr) == 0 {
		return false
	}

	l := len(arr[0])
	f := len(arr)
	param := make(S, f)

	for i := 0; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i]
		}
		if !fun(param...) {
			return false
		}
	}

	return true

}

// Any 检查类型为 S（元素类型为 T）的切片数组中是否有任一元素组合满足提供的条件函数。
//
// 参数:
// - fun: 一个函数，接受 T 类型的变长参数并返回布尔值，用于测试一组元素是否满足条件。
// - arr: 变长参数，每个元素为类型为 S 的切片（T 类型的切片），所有切片长度需一致。
//
// 返回值:
//   - 如果至少有一个元素组合使得 `fun` 返回 true，则返回 true；否则返回 false。
//     如果输入切片数组为空，则直接返回 false。
//
// 注意:
// - 所有输入切片的长度必须相等，否则函数的行为未定义。
func Any[S ~[]T, T any](fun func(x ...T) bool, arr ...S) bool {

	if len(arr) == 0 {
		return false
	}

	l := len(arr[0])
	f := len(arr)
	param := make([]T, f)

	for i := 0; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i]
		}
		if fun(param...) {
			return true
		}
	}

	return false

}

// ReverseSplit 根据提供的条件函数反向地将输入切片 S（类型为 T）分割成多个子切片，并返回这些子切片组成的切片。
// 与 `ArraySplit` 不同，此函数在条件满足的位置进行切割，并且包含切割点的元素在下一个子切片中。
//
// 参数:
// - fun: 一个函数，接受 T 类型的变长参数并返回布尔值，决定是否在当前位置进行切分。
// - arr: 变长参数，每个元素为类型为 S 的切片（T 类型的切片），所有切片长度需一致。
//
// 返回值:
//   - 一个由 S 类型子切片组成的切片，每个子切片代表原切片中满足分割条件的相邻部分。
//     区别在于，当条件满足时，该元素会包含在后续的子切片中，而非当前子切片的结尾。
//     若输入为空或首切片为空，则返回空 S 类型切片的切片
//
// 注意:
// - 所有输入切片的长度必须相等，否则函数的行为未定义。
// - 数组将在元素的右侧进行拆分。
func ReverseSplit[S ~[]T, T any](fun func(x ...T) bool, arr ...S) []S {

	if len(arr) == 0 || len(arr[0]) == 0 {
		return []S{}
	}

	l := len(arr[0])
	f := len(arr)
	param := make([]T, f)
	result := make([]S, 0)

	num := 0

	for i := 0; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i]
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
// - fun: 一个函数，接受 T 类型的变长参数并返回布尔值，决定是否在当前位置进行切分。
// - arr: 变长参数，每个元素为类型为 S 的切片（T 类型的切片），所有切片长度需一致。
//
// 返回值:
//   - 一个由 S 类型子切片组成的切片，每个子切片代表原切片中满足分割条件相邻的部分。
//     若输入为空或首切片为空，则返回空 S 类型切片的切片。
//
// 注意:
// - 所有输入切片的长度必须相等，否则函数的行为未定义。
// - 数组将在元素的左侧进行拆分。
// - 数组不会在第一个元素之前被分割。
func Split[S ~[]T, T any](fun func(x ...T) bool, arr ...S) []S {
	if len(arr) == 0 || len(arr[0]) == 0 {
		return []S{}
	}

	l := len(arr[0])
	f := len(arr)
	param := make([]T, f)
	result := make([]S, 0)

	num := 0

	for i := 1; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i]
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
// - fun: 一个函数，接受 T 类型的变长参数并返回布尔值，决定是否采用当前或下一个索引的值。
// - arr: 变长参数，每个元素为类型为 S 的切片（T 类型的切片），所有切片长度需一致。
//
// 返回值:
//   - 一个新的切片 S，其中元素根据 `fun` 的判断结果从前一个或当前索引的值中选取。
//     若输入为空或首切片为空，则返回空 T 类型切片。
//
// 注意:
// - 所有输入切片的长度必须相等，否则函数的行为未定义。
func ReverseFill[S ~[]T, T any](fun func(x ...T) bool, arr ...S) S {

	if len(arr) == 0 || len(arr[0]) == 0 {
		return S{}
	}

	if len(arr[0]) == 1 {
		return append(S{}, arr[0][0])
	}

	l := len(arr[0])
	f := len(arr)
	param := make([]T, f)
	result := make([]T, l)

	result[l-1] = arr[0][l-1]

	for i := l - 1; i >= 0; i-- {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i]
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
// - fun: 一个函数，接受 T 类型的变长参数并返回布尔值，决定是否采用当前索引的值。
// - arr: 变长参数，每个元素为类型为 S 的切片（T 类型的切片），所有切片长度需一致。
//
// 返回值:
//   - 一个新的切片 S，其中元素根据 `fun` 的判断结果从输入切片的相应位置或前一位置选取。
//     若输入为空或首切片为空，则返回空 T 类型切片。
//
// 注意:
// - 所有输入切片的长度必须相等，否则函数的行为未定义。
func Fill[S ~[]T, T any](fun func(x ...T) bool, arr ...S) S {
	if len(arr) == 0 || len(arr[0]) == 0 {
		return []T{}
	}

	if len(arr[0]) == 1 {
		return append(S{}, arr[0][0])
	}

	l := len(arr[0])
	f := len(arr)
	param := make([]T, f)
	result := make([]T, l)

	result[0] = arr[0][0]

	for i := 1; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i]
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
// - fun: 一个函数，接受 T 类型的变长参数，返回布尔值，指示是否保留当前元素。
// - arr: 变长参数，每个元素为 T 类型的切片，所有切片长度需一致。
//
// 返回值:
// - 一个新的切片 S，包含根据 `fun` 筛选后的元素。若输入为空或首切片为空，则返回空切片 S。
// 注意:
// - 所有输入切片的长度必须相等，否则函数的行为未定义。
func Filter[S ~[]T, T any](fun func(x ...T) bool, arr ...S) []T {
	if len(arr) == 0 || len(arr[0]) == 0 {
		return S{}
	}
	l := len(arr[0])
	f := len(arr)
	param := make([]T, f)
	result := make([]T, 0)

	for i := 0; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i]
		}
		if fun(param...) {
			result = append(result, arr[0][i])
		}
	}
	return result
}

// Filter2 根据提供的函数过滤多个同结构切片（类型为 T 的切片）的元素。
// 它将每个切片的对应元素作为参数传递给函数 `fun`，并仅保留 `fun`的结果切片 。
//
// 参数:
// - fun: 一个函数，接受 T 类型的变长参数，返回布尔值，指示是否保留当前元素。
// - arr: 变长参数，每个元素为 T 类型的切片，所有切片长度需一致。
//
// 返回值:
// - 一个新的切片 S，包含根据 `fun` 执行后的结果。若输入为空或首切片为空，则返回空切片 []bool。
// 注意:
// - 所有输入切片的长度必须相等，否则函数的行为未定义。
func Filter2[S ~[]T, T any](fun func(x ...T) bool, arr ...S) []bool {

	return Map(fun, arr...)
	// if len(arr) == 0 || len(arr[0]) == 0 {
	// 	return make([]bool, 0)
	// }
	// l := len(arr[0])
	// f := len(arr)
	// param := make([]T, f)
	// result := make([]bool, 0)

	// for i := 0; i < l; i++ {
	// 	for j := 0; j < f; j++ {
	// 		param[j] = arr[j][i]
	// 	}
	// 	result = append(result, fun(param...))
	// }
	// return result
}

// FilterIndex 根据提供的函数过滤多个同结构切片（类型为 T 的切片）的元素。
// 它将每个切片的对应元素作为参数传递给函数 `fun`，并仅保留 `fun` 返回真值时的位置索引。
//
// 参数:
// - fun: 一个函数，接受 T 类型的变长参数，返回布尔值，指示是否保留当前元素。
// - arr: 变长参数，每个元素为 T 类型的切片，所有切片长度需一致。
//
// 返回值:
// - 一个新的切片 S，包含根据 `fun` 筛选后的元素的索引。若输入为空或首切片为空，则返回空切片 []int。
// 注意:
// - 所有输入切片的长度必须相等，否则函数的行为未定义。
func FilterIndex[S ~[]T, T any](fun func(x ...T) bool, arr ...S) []int {
	if len(arr) == 0 || len(arr[0]) == 0 {
		return []int{}
	}
	l := len(arr[0])
	f := len(arr)
	param := make([]T, f)
	result := make([]int, 0, l)

	for i := 0; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i]
		}
		if fun(param...) {
			result = append(result, i)
		}
	}
	return result
}

// Map 对多个同结构切片（S，类型为 T 的切片）应用一个函数，将每个切片的对应元素作为参数传递，
// 并收集返回值形成一个新的 U 类型切片序列。
//
// 参数:
// - fun: 一个函数，接受 T 类型的变长参数，并返回 U 类型的结果。
// - arr: 变长参数，每个元素为类型为 S 的切片，所有切片长度需一致。
//
// 返回值:
//   - 一个 U 类型的切片，其元素为对输入切片每相同索引位置的元素应用 `fun` 函数后的结果。
//     若输入为空或首切片为空，则返回空切片。
//
// 注意:
// - 所有输入切片的长度必须相等，否则函数的行为未定义。
func Map[S ~[]T, T any, U any](fun func(x ...T) U, arr ...S) []U {
	if len(arr) == 0 || len(arr[0]) == 0 {
		return []U{}
	}
	l := len(arr[0])
	f := len(arr)
	param := make([]T, f)
	result := make([]U, l)

	for i := 0; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i]
		}
		result[i] = fun(param...)
	}
	return result
}

// FlatMap 对多个同结构切片（S，类型为 T 的切片）应用一个函数，将每个切片的对应元素作为参数传递：
// 并收集返回值形成一个新的 U 类型切片序列。
//
// 参数:
// - fun: 一个函数，接受 T 类型的变长参数，并返回 U 类型的结果。
// - arr: 变长参数，每个元素为类型为 S 的切片，所有切片长度需一致。
//
// 返回值:
//   - 一个 U 类型的切片，其元素为对输入切片每相同索引位置的元素应用 `fun` 函数后的结果。
//     若输入为空或首切片为空，则返回空切片。
//
// 注意:
// - 所有输入切片的长度必须相等，否则函数的行为未定义。
func FlatMap[S ~[]T, T any, U any](fun func(x ...T) []U, arr ...S) []U {

	if len(arr) == 0 || len(arr[0]) == 0 {
		return []U{}
	}
	l := len(arr[0])
	f := len(arr)
	param := make([]T, f)
	result := make([]U, 0)

	for i := 0; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i]
		}
		result = append(result, fun(param...)...)
	}
	return result
}

// Zip 将多个同长度的切片 S（类型为 T 的切片）按索引位置组合成新的切片 S 序列。
// 每个新切片包含的是原始切片在该索引位置上的元素。
//
// 参数:
// - arr: 变长参数，每个元素都是类型为 S 的切片，所有输入切片的长度必须相等。
//
// 返回值:
//   - 一个新切片，其中每个元素是由原始切片在相同索引处的元素组成的 S 类型切片。
//     如果输入为空或首个切片为空，则返回空切片。
//
// 注意:
// - 所有输入切片的长度必须相等，否则函数的行为未定义。
func Zip[S ~[]T, T any](arr ...S) [][]T {
	if len(arr) == 0 || len(arr[0]) == 0 {
		return [][]T{}
	}
	l := len(arr[0])
	f := len(arr)
	param := make([]T, f)
	result := make([][]T, l)

	for i := 0; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i]
		}
		result[i] = param
		param = make([]T, f)
	}
	return result
}

// Compact 移除给定切片 S 中连续重复的元素，其中 S 是泛型类型 T 的切片，且 T 必须实现了 gotools.Ordered 接口。
//
// 参数:
// - arr: 类型为 S 的切片，可能包含连续重复的元素。
//
// 返回值:
// - 一个新的切片，其中连续重复的元素已被移除。
//
// 此函数遍历输入切片，仅将与前一个元素不同的元素添加到结果切片中，
// 从而实现连续重复元素的紧凑化处理。如果输入切片为空，则返回同类型的空切片。
func Compact[S ~[]T, T gotools.Comparable](arr S) S {
	if len(arr) == 0 {
		return S{}
	}

	l := len(arr)
	result := make([]T, 0)
	result = append(result, arr[0])

	for i := 1; i < l; i++ {
		if arr[i] != arr[i-1] {
			result = append(result, arr[i])
		}
	}
	return result
}

// Compact 移除给定切片 S 中连续重复的元素，其中 S 是泛型类型 T 的切片 ，且 T 为任意类型。
//
// 参数:
// - arr: 类型为 S 的切片，可能包含连续重复的元素。
// - fun: 用于比较两个元素是否相等的函数，类型为 func(x, y T) bool，其中 T 为任意类型。 相等则视为连续重复元素，不会追加到结果中
//
// 返回值:
// - 一个新的切片，其中连续重复的元素已被移除。
//
// 此函数遍历输入切片，仅将与前一个元素不同的元素添加到结果切片中，
// 从而实现连续重复元素的紧凑化处理。如果输入切片为空，则返回同类型的空切片。
func CompactAny[S ~[]T, T any](fun func(x, y T) bool, arr S) S {
	if len(arr) == 0 {
		return S{}
	}

	l := len(arr)
	result := make([]T, 0)
	result = append(result, arr[0])

	for i := 1; i < l; i++ {
		if fun(arr[i], arr[i-1]) {
			result = append(result, arr[i])
		}
	}
	return result
}

// Reverse 反转给定的切片 S，其中 S 是泛型类型 T 的切片。
//
// 参数:
// - arr: 类型为 S 的切片，需要被反转。
//
// 返回值:
// - 反转后的切片 S，与输入切片类型相同。
//
// 此函数通过双指针技巧实现切片元素的原地交换，达到反转的目的。
// 当输入切片为空时，会直接返回同类型的空切片。
func Reverse[S ~[]T, T any](arr S) S {

	la := len(arr)
	if la == 0 {
		return S{}
	}
	if la == 1 {
		return S{arr[0]}
	}
	res := make(S, la)

	for i, j := 0, la-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = arr[j], arr[i]
	}

	if la%2 == 1 {
		res[la/2] = arr[la/2]
	}

	return res
}

// Fold 对类型为 S 的多维数组（其中 S 是 T 类型的切片）执行折叠操作，生成一个 U 类型的结果切片。
// 它接收三个参数：
// - fun：一个函数，接受变长参数 T 的切片并返回一个 U 类型的值，用于单个维度的聚合操作。
// - acc：一个累积函数，接受两个 U 类型的参数并返回一个 U 类型的值，用于将相邻结果累积。
// - arr：变长参数，表示待折叠的多维数组，数组的每个元素也是类型为 S 的切片。
//
// 返回值：
// - 一个 U 类型的切片，表示经过聚合和累积操作后的结果序列。
//
// 示例用途：
// 可以用来计算多维数组中各维度对应位置的元素经过特定运算后的序列，如多序列的逐元素加法、乘法等。
// 注意:
// - 所有输入切片的长度必须相等，否则函数的行为未定义。
func Fold[S ~[]T, T, U any](fun func(x ...T) U, acc func(x, y U) U, arr ...S) []U {
	if len(arr) == 0 || len(arr[0]) == 0 {
		return []U{}
	}

	l := len(arr[0])
	f := len(arr)
	param := make(S, f)
	result := make([]U, l)

	for i := 0; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i]
		}
		result[i] = fun(param...)

		if i > 0 {
			result[i] = acc(result[i-1], result[i])
		}

	}

	return result
}

// Reduce 对切片中的元素应用一个累积函数，并返回该函数处理后的结果。
// 此函数泛型，适用于任何类型的切片和累积函数，只要累积函数的输入输出类型与切片元素类型兼容。
//
// Parameters:
// arr - 需要被归约处理的切片。类型为 S，其中 S 是 T 类型元素的切片。
// fun - 累积函数，接收两个参数：累积结果（类型为 U）和切片中的当前元素（类型为 T），
// dvalue - 初始累积值，类型为 U。
//
//	并返回一个新的累积结果（同样为 U 类型）。此函数定义了如何将单个元素累积到整体结果中。
//
// Returns:
// result - 应用累积函数后得到的最终结果。类型为 U，累积过程的起始值为 U 的零值。
//
//	如果切片为空，则直接返回 U 的零值。
//
// 注意事项:
// - 累积函数 `fun` 应确保对于所有可能的输入都是正确的，并且应当处理好任何潜在的边界条件或错误情况。
// - 若 `result` 是引用类型（如切片、map），其初始零值可能影响结果的预期。确保理解并适当处理此类情况。
func Reduce[S ~[]T, T, U any](fun func(x U, y T) U, dvalue U, arr S) U {

	result := dvalue

	if len(arr) == 0 {
		return result
	}

	for i := 0; i < len(arr); i++ {
		result = fun(result, arr[i])
	}

	return result
}

// Intersect 计算多个切片（类型为 []T，元素类型 T 可比较）的交集。
//
// 参数:
// - arr: 变长参数，每个参数为一个待求交集的切片。
//
// 返回值:
//   - 一个新的 []T 类型的切片，包含所有输入切片中共有的元素，且元素顺序与它们在第一个切片中出现的顺序一致。
//     如果没有交集或输入为空，则返回一个空切片。
//
// 注意: T 必须实现 gotools.Comparable 接口，允许元素之间的比较操作。
func Intersect[S ~[]T, T gotools.Comparable](arr ...S) []T {
	if len(arr) == 0 {
		return make([]T, 0)
	}

	// 如果只有一个切片，则直接返回它
	if len(arr) == 1 {
		return arr[0]
	}

	nums := Map(func(x ...S) int { return len(x[0]) }, arr)

	index := FindMin(nums)

	// 使用第一个切片作为基数来收集交集元素
	intersectionMap := make(map[T][]int, len(arr[index]))
	for _, item := range arr[index] {
		intersectionMap[item] = []int{index}
	}

	// 遍历剩余的切片，仅保留交集中的元素
	for k, otherSlice := range arr {

		if k == index {
			continue
		}

		for _, item := range otherSlice {
			if _, ok := intersectionMap[item]; ok {
				intersectionMap[item] = append(intersectionMap[item], k)
			}
		}

	}

	result := mapfitler(func(k T, v []int) bool { return len(Distinct(v)) == len(arr) }, intersectionMap)

	return mapKeys(result)
}

func mapKeys[K gotools.Comparable, V any](m ...map[K]V) []K {

	if len(m) == 0 {
		return make([]K, 0)
	}

	num := Sum(Map(func(x ...map[K]V) int { return len(x[0]) }, m))

	keys := make([]K, 0, num)

	for _, v := range m {
		for k := range v {
			keys = append(keys, k)
		}
	}

	return keys
}

func mapfitler[K gotools.Comparable, V any](f func(K, V) bool, m map[K]V) map[K]V {
	filtered := make(map[K]V)
	for k, v := range m {
		if f(k, v) {
			filtered[k] = v
		}
	}
	return filtered
}

// EnumerateDense 为输入的数组中每个元素生成一个索引列表，其中的值对应该元素在数组中首次出现的位置。
// 参数:
//
//	arr: 类型为 S 的数组，S 必须是类似切片的类型且其元素类型 T 可比较。
//
// 返回值:
//
//	一个整数切片，长度与输入数组相同，其中的值表示对应元素在数组中首次出现的索引。
func EnumerateDense[S ~[]T, T gotools.Comparable](arr S) []int {

	la := len(arr)
	firstIndexMap := make(map[T]int, la)

	result := make([]int, la)

	for i, v := range arr {
		if _, exists := firstIndexMap[v]; !exists {

			firstIndexMap[v] = i
		}

		result[i] = firstIndexMap[v]
	}

	return result
}

// func SortByQ[D ~[]U, S ~[]T, T any, U gotools.Ordered](arr S, order D) (S, D) {

// 	return ArraySortBy(func(x, y U) bool { return x < y }, arr, order)

// }

// SortTwo 通过自定义比较函数和排序顺序对数组进行排序。
// 参数 fun 是用于比较两个元素大小的函数，返回 true 表示第一个元素小于第二个元素。
// 参数 arr 是需要被排序的数组。
// 参数 order 是排序顺序数组，用于定义排序规则。
// 返回值是排序后的数组。
func SortBy[D ~[]U, S ~[]T, T, U any](fun func(x, y U) bool, arr S, order D) (S, D) {

	la := len(arr)
	if la == 0 {
		return S{}, D{}
	}
	if la == 1 {
		return arr, order
	}

	res := make(S, la)
	tmp := make(D, la)
	copy(res, arr)
	copy(tmp, order)

	for i := 1; i < la; {
		if i == 0 {
			i = 2
		}
		if i >= la {
			break
		}
		if fun(tmp[i], tmp[i-1]) {
			tmp[i], tmp[i-1] = tmp[i-1], tmp[i]
			res[i], res[i-1] = res[i-1], res[i]
			i--
			continue
		}
		i++

	}

	return res, tmp
}

// func SortByLQ[D ~[]U, S ~[]T, T any, U gotools.Ordered](arr S, order D) {
// 	ArraySortByL(func(current, before U) bool { return current < before }, arr, order)
// }

// SortTwoLocal 是一个泛型排序函数，用于对两个关联数组进行排序。
// 它接受一个比较函数、一个需要排序的数组和一个与之相关的数组，
// 根据比较函数的规则对相关数组进行排序。
// 比较函数接收两个元素并返回一个布尔值，表示第一个元素是否应排在第二个元素之前。
// 注意：这个函数会直接修改输入的数组。
//
// 参数:
//   - fun: 比较函数，用于确定元素的排序顺序。
//   - arr: 需要排序的数组，类型为 S。
//   - order: 与 arr 相关的数组，类型为 D，其元素类型与 fun 的参数类型相同。
//
// 返回:
//   - 无返回值，直接修改输入的 arr 和 order 数组。
func SortByL[D ~[]U, S ~[]T, T, U any](fun func(x, y U) bool, arr S, order D) {
	la := len(arr)
	if la == 0 {
		return
	}
	if la == 1 {
		return
	}

	for i := 1; i < la; {
		if i == 0 {
			i = 2
		}
		if i >= la {
			break
		}
		if fun(order[i], order[i-1]) {
			order[i], order[i-1] = order[i-1], order[i]
			arr[i], arr[i-1] = arr[i-1], arr[i]
			i--
			continue
		}
		i++

	}

}

// Sort 对类型为 S（元素类型为 T）的切片进行自定义排序。
//
// 参数:
//   - fun: 一个比较函数，接受两个 T 类型的参数并返回一个布尔值，指示是否需要交换这两个参数的位置。
//     当 `fun(x, y)` 返回 `true` 时，在排序过程中 `x` 应该排在 `y` 之前。
//   - arr: 要排序的切片 S。
//
// 返回值:
// - 返回一个新的 S 类型切片，其中的元素根据提供的比较函数 `fun` 进行排序。
func Sort[S ~[]T, T any](fun func(x, y T) bool, arr S) S {

	la := len(arr)
	if la == 0 {
		return S{}
	}
	if la == 1 {
		return arr
	}

	res := make(S, la)
	copy(res, arr)

	for i := 1; i < la; {
		if i == 0 {
			i = 2
		}
		if i >= la {
			break
		}
		if fun(res[i], res[i-1]) {
			res[i], res[i-1] = res[i-1], res[i]
			i--
			continue
		}
		i++

	}

	return res

}

// SortQuick 使用快速排序的思路（尽管实现并不完全正确）尝试对类型为 S（元素类型为 T）的切片进行原地排序。
// 注意：当前实现并非标准快速排序算法，更像是简化冒泡排序变种。
//
// 参数:
// - arr: 要排序的切片 S，其中元素类型 T 必须是可比较的（实现 gotools.Ordered 接口）。
// func SortQuick[S ~[]T, T gotools.Ordered](arr S) S {

// 	return ArraySort(func(x, y T) bool { return x < y }, arr)

// }

// SortLocal 对类型为 S（元素类型为 T）的切片进行原地排序，依据提供的比较函数 `fun`。
//
// 参数:
//   - fun: 自定义比较函数，接受两个 T 类型的参数并返回一个布尔值。
//     当 `fun(x, y)` 返回 `true`，则在排序时 `x` 应位于 `y` 之前。
//   - arr: 要排序的本地切片 S，函数会直接修改传入的切片。
func SortLocal[S ~[]T, T any](fun func(x, y T) bool, arr S) {
	la := len(arr)
	if la == 0 {
		return
	}
	if la == 1 {
		return
	}

	for i := 1; i < la; {
		if i == 0 {
			i = 2
		}
		if i >= la {
			break
		}
		if fun(arr[i], arr[i-1]) {
			arr[i], arr[i-1] = arr[i-1], arr[i]
			i--
			continue
		}
		i++

	}

}

// Choose 根据提供的索引数组重新排列给定的泛型数组元素。
// 参数:
//
//	arr: 需要被重新排序的原数组，类型为泛型数组 S。
//	index: 一个整数索引数组，指示原数组元素在结果数组中的新位置。
//
// 返回值:
//
//	返回一个与原数组相同类型的数组 S，其元素按照 index 指定的新顺序排列。
//
// 注意:

// 索引数组中的 -1 值会被跳过，不会影响结果数组的构建。
func Choose[S ~[]T, T any](index []int, arr S) []T {

	if len(arr) < Max(index) {
		log.Println("index length not equal arr length return arr ")
		return arr
	}

	res := make([]T, len(index))

	for k, v := range index {
		if v == -1 {
			continue
		}
		res[k] = arr[v]
	}

	return res

}

// Distinct 移除类型为 S（元素类型为 T）的切片中的重复元素，并返回一个新的无重复元素的切片。
// 要求 T 类型实现 gotools.Ordered 接口，以便进行比较操作。
//
// 参数:
// - arr: 输入的切片 S，可能包含重复元素。
//
// 返回值:
// - 返回一个新的 S 类型切片，其中重复的元素已被移除，剩余元素按升序排列。
func Distinct[S ~[]T, T gotools.Ordered](arr S) S {

	if len(arr) == 0 {
		return S{}
	}
	if len(arr) == 1 {
		return arr
	}

	res := Sort(func(x, y T) bool { return x < y }, arr)

	return Compact(res)
}

// Difference 计算类型为 S（元素类型为 T）的切片中相邻元素的差值，并返回一个新的切片。
// 要求 T 类型实现 gotools.Number 接口，支持减法运算。
//
// 参数:
// - arr: 输入的切片 S，元素为可以进行减法运算的数值类型，且长度至少为 1。
//
// 返回值:
//   - 返回一个新的 S 类型切片，其中第 i 项是原切片中第 i 项与第 i-1 项的差值。
//     第一项默认为原切片的第一项，之后的每一项都是后一项减去前一项的结果。
//
// 示例:
//
//	输入: []int{5, 2, 9, 1}
//	输出: []int{5, -3, 7, -8}
func Difference[S ~[]T, T gotools.Number](arr S) S {

	la := len(arr)
	if la == 0 {
		return S{}
	}

	res := make(S, 0, la)

	res = append(res, arr[0])

	for i := 1; i < la; i++ {
		res = append(res, arr[i]-arr[i-1])
	}
	return res
}

// Count 计算满足特定条件的元素数量，这些元素来自多个具有相同长度的切片。
// 条件由提供的函数 `fun` 定义，该函数接受与输入切片数量相等的参数并返回一个布尔值。
//
// 参数:
//   - fun: 一个 variadic 函数，接受与输入切片数量相同的 T 类型参数，并返回一个布尔值。
//     当给定的元素组合满足某种条件时，应返回 `true`。
//   - arr: 变长参数，每个参数为一个 S 类型的切片（元素类型为 T）。所有切片必须具有相同的长度。
//
// 返回值:
// - 返回一个整数，表示在所有切片中满足 `fun` 条件的元素组合的数量。
//
// 注意:
// - 所有输入切片的长度必须相等，否则函数的行为未定义。
// - 如果没有切片被提供，或者提供的切片为空，则函数返回 0。
// 注意:
// - 所有输入切片的长度必须相等，否则函数的行为未定义。
func Count[S ~[]T, T any](fun func(x ...T) bool, arr ...S) int {

	if len(arr) == 0 {
		return 0
	}

	if len(arr[0]) == 0 {
		return 0
	}

	l := len(arr[0])
	f := len(arr)
	param := make([]T, f)

	num := 0
	for i := 0; i < l; i++ {
		for j := 0; j < f; j++ {
			param[j] = arr[j][i]
		}
		if fun(param...) {
			num++
		}
	}
	return num
}

// Has 检查类型为 S（元素类型为 T）的切片中是否包含指定的元素。
// 要求 T 类型实现 gotools.Comparable 接口，以便进行相等性比较。
//
// 参数:
// - arr: 要检查的切片 S。
// - elem: 要搜索的元素 T。
//
// 返回值:
// - 如果切片 `arr` 包含元素 `elem`，则返回 `true`；否则返回 `false`。
func Has[S ~[]T, T gotools.Comparable](arr S, elem T) bool {

	la := len(arr)
	if la == 0 {
		return false
	}

	for i := 0; i < la; i++ {
		if arr[i] == elem {
			return true
		}
	}
	return false
}

// HasAny 检查类型为 S（元素类型为 T）的切片是否包含至少一个指定的元素。
// 要求 T 类型实现 gotools.Comparable 接口，允许元素之间的相等性比较。
//
// 参数:
// - arr: 要检查的切片 S。
// - elems: 可变数量参数，表示要查找的一个或多个元素 T。
//
// 返回值:
// - 如果切片 `arr` 中包含 `elems` 中的至少一个元素，则返回 `true`；否则返回 `false`。
func HasAny[S ~[]T, T gotools.Comparable](arr S, elems ...T) bool {
	if len(arr) == 0 || len(elems) == 0 {
		return false
	}

	return len(Intersect(arr, elems)) > 0
}

// HasAll 检查类型为 S（元素类型为 T）的切片是否包含指定的所有元素。
// 要求 T 类型实现 gotools.Comparable 接口，允许元素之间的相等性比较。
//
// 参数:
// - arr: 要检查的切片 S。
// - elems: 可变数量参数，表示要查找的所有元素 T 组成的集合。
//
// 返回值:
// - 如果切片 `arr` 包含 `elems` 中的所有元素，则返回 `true`；否则返回 `false`。
func HasAll[S ~[]T, T gotools.Comparable](arr S, elems ...T) bool {
	if len(arr) == 0 || len(elems) == 0 {
		return false
	}

	res := Intersect(arr, elems)

	return len(res) == len(ToMap(elems))
}

// Concat 合并任意数量的同类型切片为一个新的切片。
//
// 参数:
// - arr: 可变数量的参数，每个参数为一个类型为 T 的切片，所有切片将被连接。
//
// 返回值:
// - 返回一个新的 []T 类型切片，包含输入的所有切片中的元素，保持原有的顺序。
func Concat[S ~[]T, T any](arr ...S) []T {

	if len(arr) == 0 {
		return []T{}
	}

	num := Map(func(x ...S) int { return len(x[0]) }, arr)

	res := make([]T, 0, Sum(num))

	for i := 0; i < len(num); i++ {
		res = append(res, arr[i]...)
	}
	return res
}

/*
ArrayHasSequence 检查数组arr1中是否包含数组arr2作为连续子序列。

参数：
- arr1 (A): 可能包含子序列的数组，类型A为切片的约束类型。
- arr2 (A): 需要查找的连续子序列，类型与arr1相同。

返回值：
- bool: 如果arr2是arr1中的一个连续子序列，则返回true，否则返回false。

此函数利用类型约束[A ~[]T, T gotools.Comparable]确保传入的参数为切片类型且元素可比较。
通过遍历arr1并逐一比对arr2的所有元素来判断子序列是否存在。
*/
func HasSequence[A ~[]T, T gotools.Comparable](arr1 A, arr2 A) (bool, int) {

	l1 := len(arr1)
	l2 := len(arr2)

	l := l1 - l2 + 1
	for i := 0; i < l; i++ {
		match := true
		for j := 0; j < l2; j++ {
			if arr1[i+j] != arr2[j] {
				match = false
				break
			}

		}
		if match {
			return true, i + l2 - 1
		}
	}
	return false, 0
}

/*
SequenceCount 计算一个数组中特定序列出现的次数。

参数:
- arr1[A]: 被搜索的数组，A 类型为切片，元素类型为 T。
- arr2[A]: 需要计数的序列，类型与 arr1 相同。

返回值:
- int: arr1 中 arr2 序列出现的次数。

注意:
- A 和 T 使用类型参数，要求 T 类型的元素可比较。
*/
func ArrSequenceCount[A ~[]T, T gotools.Comparable](arr1 A, arr2 A) int {

	count := 0
	num := len(arr1) - len(arr2) + 1
	for i := 0; i < num; i++ {

		if ok, index := HasSequence(arr1, arr2); ok {
			count++
			arr1 = arr1[index+1:]
		}

	}
	return count

}

// ToMap 将一个切片转换为一个映射(map)，其中切片中的每个元素作为键(key)，
// 值(value)是一个空结构体(struct{})。这个函数的目的是为了创建一个唯一的键集合，
// 由于结构体不占用存储空间，因此这种方式可以有效地表示一个集合。
// 参数:
//
//	arr []K: 一个包含待转换为映射键的元素的切片。
//
// 返回值:
//
//	map[K]struct{}: 一个映射，其中每个切片元素作为键，值是一个空结构体。
//
// 使用场景:
//
//	当需要从切片中快速查找某个元素是否存在时，可以将切片转换为映射，利用映射的O(1)查找复杂度。
//	该函数非常有用，因为它可以快速地创建一个唯一的键集合，从而节省内存空间。
func ToMap[K gotools.Comparable](arr []K) map[K]struct{} {
	m := make(map[K]struct{}, len(arr))
	for _, v := range arr {
		m[v] = struct{}{}
	}
	return m
}

// Unique 去重数组元素
// 参数:
//   - arr: 一个切片，表示要进行去重的切片。
//
// 返回:
//   - []T: 一个新的切片，表示去重后的切片。
func Unique[S ~[]T, T gotools.Comparable](arr ...S) []T {
	seen := make(map[T]struct{})
	var result []T

	for _, vv := range arr {
		for _, v := range vv {
			if _, exists := seen[v]; !exists {
				seen[v] = struct{}{}
				result = append(result, v)
			}
		}
	}

	return result
}

// Merge 通过比较函数对两个切片进行合并，返回一个新的切片, 要求传入切片必须排序
// 参数:
//   - f: 一个函数，接受两个类型为 T 的值，返回一个布尔值
//     当 `fun(x, y)` 返回 `true`，则在排序时 `x` 应位于 `y` 之前。
//   - arr: 一个切片，表示要进行合并的切片。
//
// 返回:
//   - []T: 一个新的切片，表示合并后的切片。
func Merge[S ~[]T, T any](f func(x T, y T) bool, arr ...S) []T {

	length := Reduce(func(x int, y S) int {
		return x + len(y)
	}, 0, arr)

	result := make([]T, 0, length)

	mins := make([]T, len(arr))
	index := make([]int, len(arr))

	for i := 0; i < len(arr); i++ {
		mins[i] = arr[i][0]
		index[i] = 0
	}

	for {

		minIndex := -1

		for i := range mins {

			if minIndex == -1 {
				minIndex = i
			}

			if !f(mins[minIndex], mins[i]) {
				minIndex = i
			}
		}

		if minIndex == -1 {
			break
		}

		result = append(result, mins[minIndex])

		index[minIndex] = index[minIndex] + 1

		if len(arr[minIndex]) > index[minIndex] {
			mins[minIndex] = arr[minIndex][index[minIndex]]
		} else {
			mins = append(mins[:minIndex], mins[minIndex+1:]...)
			arr = append(arr[:minIndex], arr[minIndex+1:]...)
			index = append(index[:minIndex], index[minIndex+1:]...)
		}

	}

	return result
}

// Chunk 切片分片
// 参数:
//   - size: 一个整数，表示每个切片的大小。
//   - arr: 一个切片，表示要进行分片的切片。
//
// 返回:
//   - [][]T: 一个二维切片，表示分片后的切片。
//
// 例如：ArrayChunk(2, [1, 2, 3, 4, 5, 6, 7, 8, 9])) // => [[1, 2], [3, 4], [5, 6], [7, 8], [9]]
func Chunk[S ~[]T, T any](size int, arr S) [][]T {

	if size <= 0 {
		return [][]T{}
	}

	result := make([][]T, 0, len(arr)/size+1)
	for i := 0; i < len(arr); i += size {
		result = append(result, arr[i:i+size])
	}
	return result
}

// Cartesian函数生成多个切片的笛卡尔积。
//
// 它接受类型为`[]T`的可变参数`arr`，表示要组合的切片。
// 类型`T`表示切片中元素的类型。
//
// 函数返回类型为`[][]T`的切片，表示输入切片的笛卡尔积。
// 例如:Cartesian([][]int{{1, 2}, {3, 4}}) = [][]int{{1, 3}, {1, 4}, {2, 3}, {2, 4}}
func Cartesian[S []T, T any](arr ...S) [][]T {
	if len(arr) == 0 {
		return [][]T{}
	}

	rowNum := 1
	for _, a := range arr {
		rowNum *= len(a)
	}

	res := make([][]T, rowNum)
	indices := make([]int, len(arr))

	for i := range res {
		row := make([]T, len(arr))
		for j, a := range arr {
			row[j] = a[indices[j]]
		}
		res[i] = row
		// 增加索引
		for j := len(arr) - 1; j >= 0; j-- {
			indices[j]++
			if indices[j] < len(arr[j]) {
				break
			}
			indices[j] = 0
		}
	}

	return res
}

// func Cartesian[S []T, T any](arr ...S) [][]T {
// 	if len(arr) == 0 {
// 		return [][]T{}
// 	}

// 	colNum := len(arr)
// 	rowNum := int(array.Product(array.Map(func(x ...S) int { return len(x[0]) }, arr)))

// 	res := make([][]T, colNum)

// 	for i := 0; i < colNum; i++ {
// 		res[i] = make([]T, rowNum)
// 	}

// 	res[0] = arr[0]

// 	for i := 1; i < colNum; i++ {
// 		copy(res[i], array.Rep(arr[i], len(res[i-1]), true))
// 	}

// 	for i := 0; i < colNum; i++ {
// 		if n := rowNum / len(res[i]); n > 1 {
// 			copy(res[i], array.Rep(res[i], n, false))
// 		}

// 	}

// 	return array.Zip(res...)
// }

// Flatten 从多维数组中拆分成单维数组
// 参数:
//   - arr: 一个多维数组，表示要拆分的多维数组。
//
// 返回:
//   - []T: 一个单维数组，表示拆分后的单维数组。
func Flatten[S ~[]T, T any](arr ...S) []T {

	if len(arr) == 0 {
		return []T{}
	}

	res := make([]T, 0, Len(arr))

	for i := 0; i < len(arr); i++ {
		res = append(res, arr[i]...)
	}

	return res
}

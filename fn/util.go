package fn

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"strconv"

	"github.com/zentures/cityhash"
)

// Md5 返回一个字符串的md5
// 参数：string
// 返回：string
func Md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}

// Crc32 返回一个字符串的crc32
// 参数：string
// 返回：string
func Crc32(str string) string {
	data := []byte(str)
	return fmt.Sprint(crc32.ChecksumIEEE(data))
}

func sToInt(str string) int {
	int64Val, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return int64Val
}

// StrToInt 将字符串转换为int, 首先计算md5, 然后计算crc32, 然后取余
// 参数:
//
//	str - 字符串
//	num - int范围
//
// 返回:
//
//	转换后的int
func StrToInt(str string, num int) int {
	if num <= 0 {
		return 0 // Handle invalid range
	}

	md5Res := Md5(str)
	crc32IntRes := sToInt(Crc32(md5Res))

	return crc32IntRes % num
}

func CityHash64(str string) uint64 {
	return cityhash.CityHash64([]byte(str), uint32(len([]byte(str))))
}

func CityHash32(str string) uint32 {
	return cityhash.CityHash32([]byte(str), uint32(len([]byte(str))))
}

func Println[T any](data T) {
	fmt.Println(data)
}

// Lapply 用于将一个函数应用于一个切片的每个元素。
// 参数:
//
//		f - 一个函数，接受一个类型为 T 的参数，返回一个类型为 U 的值。
//	 s - 一个包含 T 类型元素的切片。
//
// 返回:
//
//	一个包含 U 类型元素的切片。
func Lapply[S ~[]T, T any, U any](f func(x T) U, s S) []U {
	res := make([]U, len(s))

	for i := 0; i < len(s); i++ {
		res[i] = f(s[i])
	}

	return res
}

func Lapply2[S ~[]T, T any, U ~[]V, V any, R any](f func(x T, y V) R, s S, v U) []R {

	res := make([]R, len(s))

	for i := 0; i < len(s); i++ {
		res[i] = f(s[i], v[i])
	}

	return res
}

// Ifelse 根据给定的布尔条件 condition，选择返回 trueVal 或 falseVal。
// 如果 condition 为 true，则返回 trueVal；否则返回 falseVal。
// 参数:
//
//	condition - 用于判断的布尔条件。
//	trueVal - 当 condition 为 true 时返回的值。
//	falseVal - 当 condition 为 false 时返回的值。
//
// 返回:
//
//	根据 condition 的结果返回 trueVal 或 falseVal。
func Ifelse[T any](condition bool, trueVal T, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

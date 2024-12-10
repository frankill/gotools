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

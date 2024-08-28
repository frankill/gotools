package fn

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"strconv"
)

// Md5 returns the MD5 hash of the given string in hexadecimal format.
func Md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}

// Crc32 returns the CRC32 checksum of the given string in decimal format.
func Crc32(str string) string {
	data := []byte(str)
	return fmt.Sprint(crc32.ChecksumIEEE(data))
}

// StringToInt converts a string to an integer. Returns 0 if conversion fails.
func StringToInt(str string) int {
	int64Val, err := strconv.Atoi(str)
	if err != nil {
		return 0 // Default value or handle the error as needed
	}
	return int64Val
}

// StrToInt converts a string to an integer in the range [0, num). It first computes
// the MD5 hash of the string, then computes the CRC32 checksum of the hash,
// and finally takes modulo num to ensure the result is in the specified range.
func StrToInt(str string, num int) int {
	if num <= 0 {
		return 0 // Handle invalid range
	}

	md5Res := Md5(str)
	crc32IntRes := StringToInt(Crc32(md5Res))

	return crc32IntRes % num
}

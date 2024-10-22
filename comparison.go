package gotools

import "fmt"

func GtFloat64(x, y float64) bool {
	return x > y
}

func GteFloat64(x, y float64) bool {
	return x >= y
}

func LtFloat64(x, y float64) bool {
	return x < y
}

func LteFloat64(x, y float64) bool {
	return x <= y
}

func EqFloat64(x, y float64) bool {
	return x == y
}

func NotEqFloat64(x, y float64) bool {
	return x != y
}

func GtFloat32(x, y float32) bool {
	return x > y
}

func GteFloat32(x, y float32) bool {
	return x >= y
}

func LtFloat32(x, y float32) bool {
	return x < y
}

func LteFloat32(x, y float32) bool {
	return x <= y
}

func EqFloat32(x, y float32) bool {
	return x == y
}

func NotEqFloat32(x, y float32) bool {
	return x != y
}

func LtInt(x, y int) bool {
	return x < y
}

func EqInt(x, y int) bool {
	return x == y
}

func GtInt(x, y int) bool {
	return x > y
}

func GteInt(x, y int) bool {
	return x >= y
}

func LteInt(x, y int) bool {
	return x <= y
}

func EqStr(x, y string) bool {
	return x == y
}

func NotEqInt(x, y string) bool {
	return x != y
}

func Gtstr(x, y string) bool {
	return x > y
}

func GteStr(x, y string) bool {
	return x >= y
}

func LtStr(x, y string) bool {
	return x < y
}

func LteStr(x, y string) bool {
	return x <= y
}

func NotEqStr(x, y string) bool {
	return x != y
}

// Println 打印数据
// 参数:
//   - x: 数据
func Println[T any](x T) {
	fmt.Println(x)
}

package gotools

import "cmp"

type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

type Logical interface {
	bool
}

type Integer interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

type Float interface {
	float32 | float64
}

type String interface {
	string
}

type Comparable interface {
	comparable
}

type Ordered interface {
	cmp.Ordered
}

func Identity[T any](x T) T {
	return x
}

var (
	ASCInt  = ASCGeneric[int]
	DESCInt = DESCGeneric[int]
)

func ASCGeneric[T Ordered](x, y T) bool {
	return x < y
}

func DESCGeneric[T Ordered](x, y T) bool {
	return x > y
}

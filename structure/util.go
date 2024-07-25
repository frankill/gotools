package structure

import "cmp"

func Compare[T cmp.Ordered](a, b T) bool {
	return a < b
}

func CompareDesc[T cmp.Ordered](a, b T) bool {
	return a > b
}

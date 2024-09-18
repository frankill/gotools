package structure

import "github.com/frankill/gotools"

func Compare[T gotools.Ordered](a, b T) bool {
	return a < b
}

func CompareDesc[T gotools.Ordered](a, b T) bool {
	return a > b
}

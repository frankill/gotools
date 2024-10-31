package structure

import (
	"github.com/frankill/gotools"
)

type (
	Set[T gotools.Comparable] struct {
		m map[T]struct{}
	}
)

func NewSet[T gotools.Comparable](data ...T) Set[T] {
	res := make(map[T]struct{}, len(data))

	for _, v := range data {
		res[v] = struct{}{}
	}
	return Set[T]{m: res}
}

func (s *Set[T]) Has(element T) bool {

	for k := range s.m {
		if k == element {
			return true
		}
	}
	return false
}

func (s *Set[T]) Push(element T) {

	s.m[element] = struct{}{}
}

func (s *Set[T]) Move(element T) {

	delete(s.m, element)
}

func (s *Set[T]) Foreach(fn func(e T)) {

	for k := range s.m {
		fn(k)
	}
}

func (s *Set[T]) ToArray() []T {

	res := make([]T, 0, len(s.m))
	for k := range s.m {
		res = append(res, k)
	}
	return res
}

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

func (s *Set[T]) Pop() (T, bool) {

	for k := range s.m {
		delete(s.m, k)
		return k, true
	}
	var a T
	return a, false
}

func (s *Set[T]) FromChan(ch chan T) {

	for v := range ch {
		s.Push(v)
	}
}

func (s *Set[T]) FromArr(data []T) {

	for _, v := range data {
		s.Push(v)
	}
}

func (s *Set[T]) Move(element T) {

	delete(s.m, element)
}

func (s *Set[T]) Foreach(fn func(e T)) {

	for k := range s.m {
		fn(k)
	}
}

func (s *Set[T]) Len() int {

	return len(s.m)
}

func (s *Set[T]) Clear() {

	s.m = make(map[T]struct{})
}

func (s *Set[T]) IsEmpty() bool {

	return len(s.m) == 0
}

func (s *Set[T]) ToChan() chan T {

	res := make(chan T, 10)
	go func() {
		defer close(res)
		for k := range s.m {
			res <- k
		}
	}()
	return res
}

func (s *Set[T]) ToArr() []T {

	res := make([]T, 0, len(s.m))
	for k := range s.m {
		res = append(res, k)
	}
	return res
}

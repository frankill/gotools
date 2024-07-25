package structure

import "sync"

type (
	Set[T comparable] struct {
		m     map[T]struct{}
		mutex sync.RWMutex
	}
)

func NewSet[T comparable](data ...T) Set[T] {
	res := make(map[T]struct{}, len(data))

	for _, v := range data {
		res[v] = struct{}{}
	}
	return Set[T]{m: res, mutex: sync.RWMutex{}}
}

func (s *Set[T]) Has(element T) bool {

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for k := range s.m {
		if k == element {
			return true
		}
	}
	return false
}

func (s *Set[T]) Push(element T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.m[element] = struct{}{}
}

func (s *Set[T]) Move(element T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.m, element)
}

func (s *Set[T]) Foreach(fn func(e T)) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for k := range s.m {
		fn(k)
	}
}

func (s *Set[T]) ToArray() []T {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	res := make([]T, 0, len(s.m))
	for k := range s.m {
		res = append(res, k)
	}
	return res
}

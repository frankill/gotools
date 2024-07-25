package structure

import "sync"

// Queue Stack
type (
	Stack[T any] struct {
		top   *node[T]
		num   int64
		mutex sync.RWMutex
	}

	node[T any] struct {
		value T
		next  *node[T]
	}
)

func NewStack[T any](data ...T) *Stack[T] {

	s := &Stack[T]{nil, 0, sync.RWMutex{}}

	for _, v := range data {
		s.Push(v)
	}

	return s
}

func (s *Stack[T]) Len() int64 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.num
}

func (s *Stack[T]) Peek() (T, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.num == 0 {
		var a T
		return a, false
	}
	return s.top.value, true
}

func (s *Stack[T]) Pop() (T, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.num == 0 {
		var a T
		return a, false
	}

	r := s.top
	s.top = r.next
	s.num--
	return r.value, true
}

func (s *Stack[T]) Push(data T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.top = &node[T]{data, s.top}
	s.num++
}

func (s *Stack[T]) ToArray() []T {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	res := make([]T, 0, s.num)
	for current := s.top; current != nil; current = current.next {
		res = append(res, current.value)
	}
	return res
}

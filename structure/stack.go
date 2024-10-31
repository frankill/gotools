package structure

// Queue Stack
type (
	Stack[T any] struct {
		top *node[T]
		num int64
	}

	node[T any] struct {
		value T
		next  *node[T]
	}
)

func NewStack[T any](data ...T) *Stack[T] {

	s := &Stack[T]{nil, 0}

	for _, v := range data {
		s.Push(v)
	}

	return s
}

func (s *Stack[T]) Len() int64 {

	return s.num
}

func (s *Stack[T]) Peek() (T, bool) {

	if s.num == 0 {
		var a T
		return a, false
	}
	return s.top.value, true
}

func (s *Stack[T]) Pop() (T, bool) {

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

	s.top = &node[T]{data, s.top}
	s.num++
}

func (s *Stack[T]) FromArr(data []T) {

	for _, v := range data {
		s.Push(v)
	}
}

func (s *Stack[T]) FromChan(data chan T) {

	for v := range data {
		s.Push(v)
	}
}

func (s *Stack[T]) Clear() {

	s.top = nil
	s.num = 0
}

func (s *Stack[T]) IsEmpty() bool {

	return s.num == 0
}

func (s *Stack[T]) ToArr() []T {

	res := make([]T, 0, s.num)
	for current := s.top; current != nil; current = current.next {
		res = append(res, current.value)
	}
	return res
}

func (s *Stack[T]) ToChan() chan T {

	res := make(chan T, 10)
	go func() {

		defer close(res)
		for current := s.top; current != nil; current = current.next {
			res <- current.value
		}

	}()

	return res
}

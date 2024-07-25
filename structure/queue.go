package structure

import "sync"

type (
	Queue[T any] struct {
		start, end *node[T]
		num        int64
		mutex      sync.RWMutex
	}
)

func NewQueue[T any](data T) *Queue[T] {

	n := node[T]{data, nil}

	return &Queue[T]{&n, &n, 1, sync.RWMutex{}}
}

func (q *Queue[T]) Len() int64 {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return q.num
}

func (q *Queue[T]) Push(data T) {

	q.mutex.Lock()
	defer q.mutex.Unlock()

	n := &node[T]{data, nil}

	if q.num == 0 {
		q.start = n
	} else {
		q.end.next = n
	}

	q.end = n
	q.num++

}

func (q *Queue[T]) Pop() (T, bool) {

	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.num == 0 {
		var a T
		return a, false
	}

	n := q.start

	if q.num == 1 {
		q.start = nil
		q.end = nil
	} else {
		q.start = q.start.next
	}
	q.num--
	return n.value, true
}
func (q *Queue[T]) Peek() (T, bool) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	if q.num == 0 {
		var a T
		return a, false
	}
	return q.start.value, true
}

func (q *Queue[T]) ToArray() []T {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	res := make([]T, q.num)
	i := 0
	for current := q.start; current != nil; current = current.next {
		res[i] = current.value
		i++
	}
	return res
}

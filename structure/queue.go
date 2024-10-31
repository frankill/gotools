package structure

type (
	Queue[T any] struct {
		start, end *node[T]
		num        int64
	}
)

func NewQueue[T any](data ...T) *Queue[T] {

	n := &Queue[T]{nil, nil, 0}

	for _, v := range data {
		n.Push(v)
	}

	return n
}

func (q *Queue[T]) Len() int64 {

	return q.num
}

func (q *Queue[T]) Push(data T) {

	n := &node[T]{data, nil}

	if q.num == 0 {
		q.start = n
	} else {
		q.end.next = n
	}

	q.end = n
	q.num++

}

func (q *Queue[T]) FromArr(data []T) {

	for _, v := range data {
		q.Push(v)
	}

}

func (q *Queue[T]) FromChan(ch chan T) {

	for v := range ch {
		q.Push(v)
	}
}

func (q *Queue[T]) Empty() bool {

	return q.num == 0
}

func (q *Queue[T]) Clear() {

	q.start = nil
	q.end = nil
	q.num = 0
}

func (q *Queue[T]) Pop() (T, bool) {

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

	if q.num == 0 {
		var a T
		return a, false
	}
	return q.start.value, true
}

func (q *Queue[T]) ToArr() []T {

	res := make([]T, q.num)
	i := 0
	for current := q.start; current != nil; current = current.next {
		res[i] = current.value
		i++
	}
	return res
}

func (q *Queue[T]) ToChan() chan T {

	res := make(chan T, 10)
	go func() {
		defer close(res)
		for current := q.start; current != nil; current = current.next {
			res <- current.value
		}
	}()
	return res
}

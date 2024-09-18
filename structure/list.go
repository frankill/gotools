package structure

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/frankill/gotools"
)

// skip list
type (
	List[T gotools.Comparable] struct {
		root *skipnode[T]
		num  int
		gen  *rand.Rand
		pre  float64
		fun  Slf[T]
	}
	skipnode[T gotools.Comparable] struct {
		value T
		next  []*skipnode[T]
	}
	Slf[T gotools.Comparable] func(a, b T) bool
)

func NewList[T gotools.Comparable](pre float64, fun Slf[T]) *List[T] {
	gen := rand.New(rand.NewSource(time.Now().UnixNano()))
	var a T
	// 初始化根节点为类型零值
	n := &skipnode[T]{a, make([]*skipnode[T], 0)}
	res := &List[T]{n, 0, gen, pre, fun}

	return res
}

func (s *List[T]) Len() int {
	return s.num
}

func (s *List[T]) Get(value T) (T, bool) {
	var a T
	if s.num == 0 {
		return a, false
	}

	cur := s.root
	// 开始于根节点数组指针中最高值
	for i := len(cur.next) - 1; i >= 0; i-- {
		for s.fun(cur.next[i].value, value) {
			cur = cur.next[i]
		}
	}
	cur = cur.next[0]

	if cur.value == value {
		return cur.value, true
	}

	return a, false
}

func (s *List[T]) Foreach(value T) {
	cur := s.root
	h := len(cur.next)

	count := 0
	data := make([]T, 0)
	for i := h - 1; i >= 0; i-- {
		for cur.next[i] != nil && s.fun(cur.next[i].value, value) {
			count++
			data = append(data, cur.next[i].value)
			cur = cur.next[i]
		}
	}
	fmt.Println("gotools.Comparable times", count)
	fmt.Println("search node data", data)
}

func (s *List[T]) Levels() {
	cur := s.root
	h := len(cur.next)

	for i := h - 1; i >= 0; i-- {
		cur := s.root
		data := make([]T, 0)
		for cur.next[i] != nil {
			data = append(data, cur.next[i].value)
			cur = cur.next[i]
		}
		fmt.Println(i, "==>", data)
	}
	fmt.Println()
}

func (s *List[T]) Push(data T) {
	// 获取当前要插入节点不同层级的前置节点
	prev := s.getPrevious(data)

	// 去除重复输入的值
	if len(prev) > 0 && prev[0].next[0] != nil && prev[0].next[0].value == data {
		return
	}

	h := len(s.root.next)
	// 获取当前需要插入节点的层级数量
	nh := s.pickHeight()
	n := &skipnode[T]{data, make([]*skipnode[T], nh)}

	// 高于根节点层级 直接添加至根节点中
	if nh > h {
		s.root.next = append(s.root.next, n)
	}

	// 插入当前节点
	for i := 0; i < h && i < nh; i++ {
		n.next[i] = prev[i].next[i]
		prev[i].next[i] = n
	}
	// 更新计数器
	s.num++
}

func (s *List[T]) Pop(data T) (T, bool) {

	prev := s.getPrevious(data)
	h := len(prev)

	if len(prev) == 0 {
		var a T
		return a, false
	}

	cur := prev[0].next[0]

	if cur == nil || cur.value != data {
		return data, false
	}

	for i := 0; i < h && i < len(cur.next); i++ {
		prev[i].next[i] = cur.next[i]
	}

	for i := h - 1; i >= 0; i-- {
		if s.root.next[i] == nil {
			s.root.next = s.root.next[:i]
		} else {
			break
		}
	}

	return cur.value, true

}

func (s *List[T]) getPrevious(value T) []*skipnode[T] {
	cur := s.root
	h := len(cur.next)
	nodes := make([]*skipnode[T], h)
	for i := h - 1; i >= 0; i-- {
		for cur.next[i] != nil && s.fun(cur.next[i].value, value) {
			cur = cur.next[i]
		}
		nodes[i] = cur
	}
	return nodes
}

func (s *List[T]) pickHeight() int {
	h := 1
	// 遍历多次，获取需要产生的层级数量
	for s.gen.Float64() > s.pre {
		h++
	}
	if h > len(s.root.next) {
		return len(s.root.next) + 1
	}
	return h
}

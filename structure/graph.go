package structure

import (
	"fmt"
	"strings"
)

type Graph[T any, U comparable] struct {
	ednum     int
	vernum    int
	vertexs   []*Vertex[T]
	match     func(v, d T) bool
	index     map[U]int
	transform func(v T) U
}

type Vertex[T any] struct {
	data T
	e    *Edge
}

type Edge struct {
	loc  int // 边指向的元素位置索引
	next *Edge
}

func (v *Vertex[T]) SetNext(next *Edge) bool {

	if v.e == nil {
		v.e = next
		return true
	}

	for e := v.e; e != nil; e = e.next {
		if e.loc == next.loc {
			return false
		}
	}
	v.e.next = next
	return true
}
func (e *Edge) Next() *Edge {
	return e.next
}

func (g *Graph[T, U]) GetEdge(v T) *Edge {

	if index, ok := g.index[g.transform(v)]; ok {
		return g.vertexs[index].e
	}

	return nil
}

func (g *Graph[T, U]) EageToElement(e *Edge) []T {

	var r []T
	for e != nil {
		r = append(r, g.vertexs[e.loc].data)
		e = e.next
	}
	return r
}

func (g *Graph[T, U]) Getloc() int {

	return g.vernum

}

func (g *Graph[T, U]) GetIndex(data T) int {

	if index, ok := g.index[g.transform(data)]; ok {
		return index
	}

	return -1

}
func (g *Graph[T, U]) DrawGraph() string {
	var result strings.Builder

	result.WriteString("Graph Structure:\n")
	result.WriteString("Vertices: ")

	for i, v := range g.vertexs {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(fmt.Sprintf("%v", v.data))
	}

	result.WriteString("\nEdges:\n")

	for _, v := range g.vertexs {
		result.WriteString(fmt.Sprintf("Vertex %v: ", v.data))
		for e := v.e; e != nil; e = e.next {
			result.WriteString(fmt.Sprintf("%v => %v; ", v.data, g.vertexs[e.loc].data))
		}
		result.WriteString("\n")
	}

	return result.String()
}
func (g *Graph[T, U]) Addedge(ds, de T) bool {

	vs := g.GetIndex(ds)
	ve := g.GetIndex(de)

	if vs == -1 || ve == -1 {
		return false
	}

	if g.vertexs[vs].SetNext(newEdge(ve)) {
		g.ednum++
	}

	// if g.vertexs[ve].SetNext(newEdge(vs)) {
	// 	g.ednum++
	// }

	return true

}

func (g *Graph[T, U]) AddVer(v T) {

	if _, ok := g.index[g.transform(v)]; ok {
		return
	}

	g.vertexs = append(g.vertexs, newVertex(v))
	g.vernum++

	g.index[g.transform(v)] = g.vernum - 1

}

func newEdge(loc int) *Edge {
	return &Edge{loc: loc, next: nil}
}

func newVertex[T any](data T) *Vertex[T] {
	return &Vertex[T]{data: data, e: nil}
}

func NewGraph[T any, U comparable](fn func(v, d T) bool, size int, index func(v T) U) *Graph[T, U] {

	r := new(Graph[T, U])
	r.match = fn
	r.vertexs = make([]*Vertex[T], 0, size)
	r.index = make(map[U]int)
	r.transform = index
	return r

}

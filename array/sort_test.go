package array_test

import (
	"reflect"
	"testing"

	"github.com/frankill/gotools/array"
)

func TestSortFun(t *testing.T) {
	arr := []int{3, 2, 1, 4, 5}
	expected := []int{1, 2, 3, 4, 5}
	actual := array.SortFun(func(x, y int) bool { return x < y }, arr)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestSortFunLocal(t *testing.T) {
	arr := []int{3, 2, 1, 4, 5}
	expected := []int{1, 2, 3, 4, 5}
	array.SortFun2(func(x, y int) bool { return x < y }, arr)
	if !reflect.DeepEqual(arr, expected) {
		t.Errorf("expected %v, got %v", expected, arr)
	}
}

func TestOrderFun(t *testing.T) {
	arr := []int{3, 2, 1, 4, 5}
	order := []int{1, 2, 3, 4, 5}
	expected := []int{3, 2, 1, 4, 5}
	actual, _ := array.OrderFun(func(x, y int) bool { return x < y }, arr, order)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestOrderFunLocal(t *testing.T) {
	arr := []int{3, 2, 1, 4, 5}
	order := []int{1, 2, 3, 4, 5}
	expected := []int{3, 2, 1, 4, 5}
	array.OrderFun2(func(x, y int) bool { return x < y }, arr, order)
	if !reflect.DeepEqual(arr, expected) {
		t.Errorf("expected %v, got %v", expected, arr)
	}
}

func TestSort(t *testing.T) {
	arr := []int{3, 2, 1, 4, 5}
	expected := []int{1, 2, 3, 4, 5}
	actual := array.Sort(arr, false)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestSortL(t *testing.T) {
	arr := []int{3, 2, 1, 4, 5}
	expected := []int{1, 2, 3, 4, 5}
	array.SortL(arr)
	if !reflect.DeepEqual(arr, expected) {
		t.Errorf("expected %v, got %v", expected, arr)
	}
}

func TestSortR(t *testing.T) {
	arr := []int{3, 2, 1, 4, 5}
	expected := []int{5, 4, 3, 2, 1}
	array.SortR(arr)
	if !reflect.DeepEqual(arr, expected) {
		t.Errorf("expected %v, got %v", expected, arr)
	}
}

func TestOrder(t *testing.T) {
	arr := []int{3, 2, 1, 4, 5}
	order := []int{1, 2, 3, 4, 5}
	expected := []int{3, 2, 1, 4, 5}
	actual, _ := array.Order(arr, order, false)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestOrderL(t *testing.T) {
	arr := []int{3, 2, 1, 4, 5}
	order := []int{1, 2, 3, 4, 5}
	expected := []int{3, 2, 1, 4, 5}
	array.OrderL(arr, order)
	if !reflect.DeepEqual(arr, expected) {
		t.Errorf("expected %v, got %v", expected, arr)
	}
}

func TestOrderR(t *testing.T) {
	arr := []int{3, 2, 1, 4, 5}
	order := []int{1, 2, 3, 4, 5}
	expected := []int{5, 4, 1, 2, 3}
	array.OrderR(arr, order)
	if !reflect.DeepEqual(arr, expected) {
		t.Errorf("expected %v, got %v", expected, arr)
	}
}

package structure

import (
	"fmt"
	"testing"
)

func TestSetOperations(t *testing.T) {
	// 创建一个新的集合
	set := NewSet(1, 2, 3)

	// 测试 Has 方法
	if !set.Has(2) {
		t.Errorf("Expected set to have element 2, but it doesn't")
	}

	// 测试 Push 方法
	set.Push(4)

	expected := []int{1, 2, 3, 4}
	for _, e := range expected {
		if !set.Has(e) {
			t.Errorf("Expected set to have element %d after Push, but it doesn't", e)
		}
	}

	// 测试 Move 方法
	set.Move(2)
	if set.Has(2) {
		t.Errorf("Expected set to not have element 2 after Move, but it does")
	}

	// 测试 Foreach 方法
	fmt.Println("Foreach Print:")
	set.Foreach(func(e int) {
		fmt.Println(e)
	})

	// 可以继续添加其他测试场景...
}

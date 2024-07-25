package structure

import (
	"testing"

	"github.com/frankill/gotools/array"
)

func TestListOperations(t *testing.T) {
	// 创建一个新的跳表
	list := NewList(0.5, Compare[int])

	// 测试 Push 方法

	num := array.ArraySeq(1, 101, 1)

	array.ForEach(func(x ...int) {
		list.Push(x[0])
	}, num)

	// 验证 Len 方法
	expectedLen := len(num)
	if list.Len() != expectedLen {
		t.Errorf("Expected length %d, got %d", expectedLen, list.Len())
	}

	// 测试 Get 方法
	val, found := list.Get(6)
	if !found || val != 6 {
		t.Errorf("Expected to find value 6, got %v (found: %v)", val, found)
	}

	// 测试 Pop 方法
	poppedVal, ok := list.Pop(6)
	if !ok || poppedVal != 6 {
		t.Errorf("Expected to pop value 6, got %v (success: %v)", poppedVal, ok)
	}

	// 测试 Get 方法，确认已删除
	_, found = list.Get(6)
	if found {
		t.Errorf("Expected value 6 to be deleted, but found it")
	}

	// 添加更多测试场景...

	// 打印跳表各层数据，仅供查看
	list.Levels()
}

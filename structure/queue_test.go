package structure_test

import (
	"testing"

	"github.com/frankill/gotools/structure"
)

func TestQueueOperations(t *testing.T) {
	// 创建一个新的队列，并进行操作测试
	queue := structure.NewQueue(1)

	// 测试 Push 操作
	queue.Push(2)
	queue.Push(3)

	// 测试 Len 方法
	if queue.Len() != 3 {
		t.Errorf("Expected queue length to be 3, got %d", queue.Len())
	}

	// 测试 Peek 方法
	if front, ok := queue.Peek(); !ok || front != 1 {
		t.Errorf("Expected Peek() to return value 1, got %v", front)
	}

	// 测试 Pop 方法
	if front, ok := queue.Pop(); !ok || front != 1 {
		t.Errorf("Expected Pop() to return value 1, got %v", front)
	}
	if queue.Len() != 2 {
		t.Errorf("Expected queue length after Pop() to be 2, got %d", queue.Len())
	}

	// 测试 ToArray 方法
	arr := queue.ToArr()
	expected := []int{2, 3}
	for i := 0; i < len(expected); i++ {
		if arr[i] != expected[i] {
			t.Errorf("Expected element at index %d to be %d, got %d", i, expected[i], arr[i])
		}
	}
}

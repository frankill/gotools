package structure_test

import (
	"reflect"
	"testing"

	"github.com/frankill/gotools/structure"
)

func TestStackOperations(t *testing.T) {
	// Create a new stack with some initial data
	stack := structure.NewStack(1, 2, 3)

	// Test push method
	stack.Push(4)

	if top, ok := stack.Pop(); ok && top != 4 {
		t.Errorf("Expected top element after push to be 4, got %v", top)
	}

	// Test pop method
	if top, ok := stack.Pop(); ok && top != 3 {
		t.Errorf("Expected top element after pop to be 3, got %v", top)
	}

	// Test To_arr method
	expectedArr := []int{2, 1}
	arr := stack.ToArray()

	if !reflect.DeepEqual(arr, expectedArr) {
		t.Errorf("Expected stack elements in array form to be %v, got %v", expectedArr, arr)
	}

	// Test pop on empty stack
	stack.Pop()
	stack.Pop()
	if _, ok := stack.Pop(); ok {
		t.Errorf("Expected pop on empty stack to return false, but got true")
	}
}

package array_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/frankill/gotools/array"
)

func TestSequenceMatch(t *testing.T) {
	tests := []struct {
		name  string
		data  []int
		by    []int
		order []int
		input []int
		id    []int
		want  map[int]bool
	}{
		{
			name:  "Test1",
			data:  []int{1, 2, 3, 4, 5},
			by:    []int{1, 1, 2, 2, 2},
			order: []int{1, 1, 1, 1, 1},
			input: []int{3, 4, 6},
			id:    []int{1, 3},
			want:  map[int]bool{1: false, 2: false},
		},
		{
			name:  "Test1",
			data:  []int{1, 2, 3, 4, 5},
			by:    []int{1, 1, 2, 2, 2},
			order: []int{1, 1, 1, 1, 1},
			input: []int{3, 6},
			id:    []int{1, 2},
			want:  map[int]bool{1: false, 2: false},
		},
		// 添加更多测试用例
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 调用待测试的函数
			got := array.SequenceMatch(tt.by, tt.data, tt.order)(tt.input)(tt.id)

			// 比较结果
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupArrayByOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSequenceCount(t *testing.T) {
	tests := []struct {
		name  string
		data  []int
		by    []int
		order []int
		input []int
		id    []int
		want  map[int]int
	}{
		{
			name:  "Test1",
			data:  []int{1, 2, 1, 2, 3, 1, 2},
			by:    []int{1, 1, 2, 2, 2, 2, 2},
			order: []int{1, 1, 1, 1, 1, 1, 1},
			input: []int{1, 2, 3, 4, 5},
			id:    []int{1, 2},
			want:  map[int]int{1: 1, 2: 2},
		},
		{
			name:  "Test1",
			data:  []int{1, 3, 1, 2, 3, 1, 2},
			by:    []int{1, 1, 2, 2, 2, 2, 2},
			order: []int{1, 1, 1, 1, 1, 1, 1},
			input: []int{1, 2, 3, 4, 5},
			id:    []int{1, 3},
			want:  map[int]int{1: 1, 2: 0},
		},
		// 添加更多测试用例
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 调用待测试的函数
			got := array.SequenceCount(tt.by, tt.data, tt.order)(tt.input)(tt.id)

			// 比较结果
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupArrayByOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

type Event struct {
	UserID    string
	EventType string
	Timestamp int64
}

// TestWindowFunnel 是一个测试用例，用于验证 WindowFunnel 函数的正确性。
// 它使用预定义的数据集和事件序列，针对不同的模式检查结果是否符合预期。
func TestWindowFunnel(t *testing.T) {
	// 测试数据准备
	var events = []Event{
		{"user1", "view_product", 1},
		{"user1", "aaa", 2},
		{"user1", "add_to_cart", 3},
		{"user2", "view_product", 1},
		{"user2", "view_product", 2},
		{"user2", "add_to_cart", 3},
		{"user3", "view_product", 1},
		{"user3", "add_to_cart", 2},
		{"user3", "add_to_cart", 3},
		{"user3", "purchase", 4},
	}

	// 提取用户ID、事件类型和时间戳
	var userIDs = make([]string, len(events))
	var eventTypes = make([]string, len(events))
	var timestamps = make([]int64, len(events))

	for i, e := range events {
		userIDs[i] = e.UserID
		eventTypes[i] = e.EventType
		timestamps[i] = e.Timestamp
	}

	// 定义漏斗的事件序列
	funnelSteps := []string{"view_product", "add_to_cart", "purchase"}

	// 创建漏斗分析函数
	funnel := array.WindowFunnel(userIDs, eventTypes, timestamps)

	// 预期结果
	expectedStrictOrder := map[string]int{"user1": 1, "user2": 2, "user3": 3}
	expectedStrictDedup := map[string]int{"user1": 2, "user2": 1, "user3": 2}
	expectedStrictIncrease := map[string]int{"user1": 2, "user2": 2, "user3": 3}

	// 测试 strict_order 模式
	resultStrictOrder := funnel(funnelSteps)("strict_order")
	checkResult(t, resultStrictOrder, expectedStrictOrder, "strict_order")

	// 测试 strict_dedup 模式
	resultStrictDedup := funnel(funnelSteps)("strict_dedup")
	checkResult(t, resultStrictDedup, expectedStrictDedup, "strict_dedup")

	// 测试 strict_increase 模式
	resultStrictIncrease := funnel(funnelSteps)("strict_increase")
	checkResult(t, resultStrictIncrease, expectedStrictIncrease, "strict_increase")
}

// checkResult 是一个辅助函数，用于比较测试结果和预期结果，并报告差异。
func checkResult(t *testing.T, result, expected map[string]int, mode string) {
	t.Helper()
	for k, v := range expected {
		if r, ok := result[k]; !ok || r != v {
			t.Errorf("For mode %s, expected %v for key %s but got %v", mode, v, k, r)
		}
	}
	fmt.Printf("Mode %s passed.\n", mode)
}

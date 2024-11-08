package group_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/frankill/gotools/array"
	"github.com/frankill/gotools/group"
	"github.com/frankill/gotools/pair"
)

// 测试用例函数 TestRetention 用于验证 Retention 函数的功能是否正确。
// 本测试通过提供一组预定义的分类标识、数据集以及条件函数，来检查 Retention 函数生成的闭包是否能正确返回预期的筛选结果映射。
func TestRetention(t *testing.T) {
	// 定义分类标识
	by := []string{"b", "d", "a", "a", "a", "b"}

	// 定义数据集，与分类标识一一对应
	order := []int{1, 2, 3, 1, 4, 4}

	// 定义条件函数
	f1 := func(x []int) bool { return array.Has(x, 1) } // 检查数据集中是否包含数字1
	f2 := func(x []int) bool { return array.Has(x, 2) } // 检查数据集中是否包含数字2
	f3 := func(x []int) bool { return array.Has(x, 4) } // 检查数据集中是否包含数字4

	// 生成基于条件函数的闭包
	reten := group.Retention(by, order)

	// 预期结果映射
	expected := map[string][]bool{
		"a": {true, false, true},   // 类别'a'在第一个条件(f1)下为真，因此f2结果也为真（因为f1为真），f3同样处理
		"b": {true, false, true},   // 同上逻辑
		"d": {false, false, false}, // 类别'd'在所有条件下的结果都为假
	}

	// 执行闭包函数并获取实际结果
	result := reten(f1, f2, f3)

	// 比较预期结果与实际结果
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Test failed. Expected: %v, got: %v", expected, result)
	} else {
		t.Logf("Test passed. Result: %v", result)
	}
}

func TestGroupCount(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		by       []int
		data     []string
		expected map[int]int
	}{
		{
			name:     "Test 1",
			by:       []int{1, 1, 2, 2},
			data:     []string{"a", "b", "c", "d"},
			expected: map[int]int{1: 2, 2: 2},
		},

		// Add more test cases as needed
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function under test
			result := group.Count2(tt.by, tt.data)

			// Check if the result matches the expected value
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
			}
			for k, v := range tt.expected {
				if result[k] != v {
					t.Errorf("Expected count %d for key %v, got %d", v, k, result[k])
				}
			}
		})
	}
}

func TestGroupDistinct(t *testing.T) {
	// Define test cases
	tests := []struct {
		by     []int       // Input for 'by' parameter
		data   []string    // Input for 'data' parameter
		expect map[int]int // Expected result
	}{
		{
			by:     []int{1, 2, 1, 3, 3},
			data:   []string{"a", "b", "c", "d", "e"},
			expect: map[int]int{1: 2, 2: 1, 3: 2},
		},
		{
			by:     []int{1, 1, 1, 1},
			data:   []string{"a", "b", "d", "d"},
			expect: map[int]int{1: 3},
		},
		// Add more test cases here
	}

	// Run test cases
	for _, tt := range tests {
		result := group.Distinct(tt.by, tt.data)

		// Check if the result matches the expected result
		if len(result) != len(tt.expect) {
			t.Errorf("GroupDistinct() got %v, want %v", len(result), len(tt.expect))
		}
		for k, v := range tt.expect {
			if result[k] != v {
				t.Errorf("GroupDistinct() got %v for key %v, want %v", result[k], k, v)
			}
		}
	}
}

func TestGroupMax(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		by       []int
		data     []string
		expected map[int]string
	}{
		{
			name:     "Test 1",
			by:       []int{1, 1, 2, 2},
			data:     []string{"a", "b", "c", "d"},
			expected: map[int]string{1: "b", 2: "d"},
		},
		{
			name:     "Test 2",
			by:       []int{1, 2, 3},
			data:     []string{"a", "b", "c"},
			expected: map[int]string{1: "a", 2: "b", 3: "c"},
		},
		// Add more test cases if needed
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function
			result := group.Max(tt.by, tt.data)

			// Check if the result is equal to the expected result
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, but got %d", len(tt.expected), len(result))
			}
			for k, v := range tt.expected {
				if result[k] != v {
					t.Errorf("Expected %s for key %d, but got %s", v, k, result[k])
				}
			}
		})
	}
}

func TestGroupMin(t *testing.T) {
	tests := []struct {
		name string
		by   []int
		data []string
		want map[int]string
	}{
		{
			name: "Test 1",
			by:   []int{1, 1, 2, 2},
			data: []string{"apple", "banana", "cherry", "date"},
			want: map[int]string{1: "apple", 2: "cherry"},
		},
		{
			name: "Test 2",
			by:   []int{1, 2, 3, 4},
			data: []string{"apple", "banana", "cherry", "date"},
			want: map[int]string{1: "apple", 2: "banana", 3: "cherry", 4: "date"},
		},
		// Add more test cases if needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := group.Min(tt.by, tt.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupMin() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestGroupSum(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		by       []int
		data     []int
		expected map[int]int
	}{
		{
			name:     "Test 1",
			by:       []int{1, 2, 1, 3},
			data:     []int{10, 20, 10, 30},
			expected: map[int]int{1: 20, 2: 20, 3: 30},
		},

		// Add more test cases if needed
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function under test
			result := group.Sum(tt.by, tt.data)

			// Check if the result matches the expected value
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
			for k, v := range tt.expected {
				if result[k] != v {
					t.Errorf("Expected %v, got %v for key %v", v, result[k], k)
				}
			}
		})
	}
}
func TestGroupArrayPair(t *testing.T) {
	tests := []struct {
		name   string
		by     []int
		first  []string
		second []float64
		want   map[int]pair.Pair[[]string, []float64]
	}{
		{
			name:   "Test 1",
			by:     []int{1, 1, 2, 2},
			first:  []string{"a", "b", "c", "d"},
			second: []float64{1.1, 1.2, 2.1, 2.2},
			want: map[int]pair.Pair[[]string, []float64]{
				1: {
					First:  []string{"a", "b"},
					Second: []float64{1.1, 1.2},
				},
				2: {
					First:  []string{"c", "d"},
					Second: []float64{2.1, 2.2},
				},
			},
		},
		// Add more test cases if needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := group.Pair(tt.by, tt.first, tt.second); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupArrayPair() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestGroupArray(t *testing.T) {
	// Test case 1
	data1 := []int{1, 2, 3}
	by1 := []int{0, 0, 1}
	expected1 := map[int][]int{0: {1, 2}, 1: {3}}
	result1 := group.By(by1, data1)
	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Test case 1 failed. Expected %v, got %v", expected1, result1)
	}

	// Test case 3
	data3 := [][]float64{}
	by3 := []float64{}
	expected3 := map[float64][][]float64{}
	result3 := group.By(by3, data3)
	if !reflect.DeepEqual(result3, expected3) {
		t.Errorf("Test case 3 failed. Expected %v, got %v", expected3, result3)
	}

	// Test case 4
	data4 := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	by4 := []int{0, 0, 2}
	expected4 := map[int][][]int{0: {{1, 2, 3}, {4, 5, 6}}, 2: {{7, 8, 9}}}
	result4 := group.By(by4, data4)
	if !reflect.DeepEqual(result4, expected4) {
		t.Errorf("Test case 4 failed. Expected %v, got %v", expected4, result4)
	}
}

func TestGroupArrayByOrder(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name  string
		data  []int
		by    []int
		order []int
		want  map[int][]int
	}{
		{
			name:  "Test1",
			data:  []int{1, 2, 3, 4, 5},
			by:    []int{1, 1, 2, 2, 2},
			order: []int{2, 1, 3, 4, 5},
			want:  map[int][]int{1: {2, 1}, 2: {3, 4, 5}},
		},
		// 添加更多测试用例
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 调用待测试的函数
			got := group.ByOrder(tt.by, tt.data, tt.order)

			// 比较结果
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupArrayByOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRowNumber(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		by       []int
		order    []int
		expected []int
	}{
		{
			name:     "Test1",
			by:       []int{1, 2, 3, 1},
			order:    []int{4, 3, 2, 1},
			expected: []int{1, 0, 0, 0},
		},
		{
			name:     "Test2",
			by:       []int{5, 5, 2, 2, 3},
			order:    []int{1, 2, 1, 2, 1},
			expected: []int{0, 1, 0, 1, 0},
		},
		{
			name:     "Test3",
			by:       []int{},
			order:    []int{},
			expected: []int{},
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function under test
			result := group.RowNumber(tt.by, tt.order)

			// Check if the result matches the expected value
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("RowNumber(%v, %v) = %v, expected %v", tt.by, tt.order, result, tt.expected)
			}
		})
	}
}

func TestMaxValue(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		data     []int
		by       []int
		expected []int
	}{
		{
			name:     "Test 1",
			data:     []int{1, 2, 3, 4, 5},
			by:       []int{},
			expected: []int{5, 5, 5, 5, 5},
		},
		{
			name:     "Test 2",
			data:     []int{10, 20, 30, 40, 50},
			by:       []int{1, 1, 2, 2, 2},
			expected: []int{20, 20, 50, 50, 50},
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function
			result := group.MaxValue(tt.by, tt.data)

			// Check if the result is correct
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MaxValue(%v, %v) = %v, expected %v", tt.data, tt.by, result, tt.expected)
			}
		})
	}
}

func TestMinValue(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name string
		data []int
		by   []bool
		want []int
	}{
		{
			name: "Test 1",
			data: []int{4, 2, 5, 1, 3},
			by:   []bool{true, true, true, false, false},
			want: []int{2, 2, 2, 1, 1},
		},
		{
			name: "Test 2",
			data: []int{4, 2, 5, 1, 3},
			by:   []bool{},
			want: []int{1, 1, 1, 1, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 调用函数得到结果
			got := group.MinValue(tt.by, tt.data)

			// 检查结果是否正确
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MinValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLastValue(t *testing.T) {
	tests := []struct {
		name string
		data []int
		by   []int
		want []int
	}{
		{
			name: "Test 1",
			data: []int{1, 2, 3, 4, 5},
			by:   []int{},
			want: []int{5, 5, 5, 5, 5},
		},

		// Add more test cases if needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := group.LastValue(tt.by, tt.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LastValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFirstValue(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		data     []int
		by       []int
		expected []int
	}{
		{
			name:     "Test1",
			data:     []int{1, 2, 3, 4, 5},
			by:       []int{1, 1, 2, 2, 3},
			expected: []int{1, 1, 3, 3, 5},
		},

		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function under test
			result := group.FirstValue(tt.by, tt.data)

			// Check if the result matches the expected value
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FirstValue(%v, %v) = %v, expected %v", tt.data, tt.by, result, tt.expected)
			}
		})
	}
}

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
			got := group.SequenceMatch(tt.by, tt.data, tt.order)(tt.input)(tt.id)

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
			got := group.SequenceCount(tt.by, tt.data, tt.order)(tt.input)(tt.id)

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
	funnel := group.WindowFunnel(userIDs, eventTypes, timestamps)

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

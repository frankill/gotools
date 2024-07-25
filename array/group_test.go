package array

import (
	"reflect"
	"testing"
)

// 测试用例函数 TestRetention 用于验证 Retention 函数的功能是否正确。
// 本测试通过提供一组预定义的分类标识、数据集以及条件函数，来检查 Retention 函数生成的闭包是否能正确返回预期的筛选结果映射。
func TestRetention(t *testing.T) {
	// 定义分类标识
	by := []string{"b", "d", "a", "a", "a", "b"}

	// 定义数据集，与分类标识一一对应
	order := []int{1, 2, 3, 1, 4, 4}

	// 定义条件函数
	f1 := func(x []int) bool { return ArrayHas(x, 1) } // 检查数据集中是否包含数字1
	f2 := func(x []int) bool { return ArrayHas(x, 2) } // 检查数据集中是否包含数字2
	f3 := func(x []int) bool { return ArrayHas(x, 4) } // 检查数据集中是否包含数字4

	// 生成基于条件函数的闭包
	reten := Retention(by, order)

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

func TestGroupGenerateFilter(t *testing.T) {
	// Test case 1: Empty data slice

	f := func(x int, y []int) []bool { return ArrayMap(func(x ...int) bool { return x[0]%2 == 0 }, y) }

	data := []int{}
	by := []int{1, 2, 3}
	expected := []int{}
	result, _ := GroupGenerateFilter(by, data)(f)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Test case 2: Grouped and filtered data
	data = []int{1, 2, 3, 4, 5, 6}
	by = []int{1, 2, 3, 1, 2, 3}
	expected = []int{2, 1, 3}
	result, _ = GroupGenerateFilter(by, data)(f)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

}

func TestGroupGenerate(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name string
		by   []int
		data []string
		fun  Groupfun[string]
		want map[int]string
	}{
		{
			name: "Group by first character",
			by:   []int{0, 2, 2, 1},
			data: []string{"apple", "banana", "cherry", "date"},
			fun:  func(s []string) string { return s[0] },
			want: map[int]string{0: "apple", 1: "date", 2: "banana"},
		},
		// TODO: 添加更多测试用例
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GroupApply(tt.by, tt.data)(tt.fun)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupGenerate() = %v, want %v", got, tt.want)
			}
		})
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
			result := GroupCount(tt.by, tt.data)

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
		result := GroupDistinct(tt.by, tt.data)

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
			result := GroupMax(tt.by, tt.data)

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
			if got := GroupMin(tt.by, tt.data); !reflect.DeepEqual(got, tt.want) {
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
			result := GroupSum(tt.by, tt.data)

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
		want   map[int]Pair[[]string, []float64]
	}{
		{
			name:   "Test 1",
			by:     []int{1, 1, 2, 2},
			first:  []string{"a", "b", "c", "d"},
			second: []float64{1.1, 1.2, 2.1, 2.2},
			want: map[int]Pair[[]string, []float64]{
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
			if got := GroupPair(tt.by, tt.first, tt.second); !reflect.DeepEqual(got, tt.want) {
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
	result1 := GroupData(by1, data1)
	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Test case 1 failed. Expected %v, got %v", expected1, result1)
	}

	// Test case 3
	data3 := [][]float64{}
	by3 := []float64{}
	expected3 := map[float64][][]float64{}
	result3 := GroupData(by3, data3)
	if !reflect.DeepEqual(result3, expected3) {
		t.Errorf("Test case 3 failed. Expected %v, got %v", expected3, result3)
	}

	// Test case 4
	data4 := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	by4 := []int{0, 0, 2}
	expected4 := map[int][][]int{0: {{1, 2, 3}, {4, 5, 6}}, 2: {{7, 8, 9}}}
	result4 := GroupData(by4, data4)
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
			got := GroupByOrder(tt.by, tt.data, tt.order)

			// 比较结果
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupArrayByOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

package array

import (
	"reflect"
	"testing"
)

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
			result := RowNumber(tt.by, tt.order)

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
			result := MaxValue(tt.by, tt.data)

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
			got := MinValue(tt.by, tt.data)

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
			if got := LastValue(tt.by, tt.data); !reflect.DeepEqual(got, tt.want) {
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
			result := FirstValue(tt.by, tt.data)

			// Check if the result matches the expected value
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FirstValue(%v, %v) = %v, expected %v", tt.data, tt.by, result, tt.expected)
			}
		})
	}
}

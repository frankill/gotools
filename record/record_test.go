package record_test

import (
	"fmt"
	"testing"

	"github.com/frankill/gotools/array"
)

// TestMapIntersect tests the MapIntersect function
func TestMapIntersect(t *testing.T) {
	// Define some sample input maps for testing
	// input1 := map[string][]int{"a": {1, 2, 3}, "b": {4, 5, 7, 8}}
	// // input2 := map[string][]int{"c": {2, 3, 4}, "d": {5, 6, 7, 8}}
	// // input3 := map[string][]int{"a": {3, 4, 5}, "b": {4, 5, 6}}

	// // Call the MapIntersect function with the sample input maps
	// result := MapIntersect(input1)

	// Define the expected outputf
	// fmt.Println(result)
	// fmt.Println(ArrayZip(result.First...))

	df := [][]byte{[]byte("abc"), []byte("test"), []byte("中心")}

	a := array.ArrayConcat(df...)
	fmt.Println(string(a[0:3]), a, df)

}

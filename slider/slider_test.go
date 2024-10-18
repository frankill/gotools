package slider_test

import (
	"testing"
	"time"

	"github.com/frankill/gotools"
	"github.com/frankill/gotools/array"
	"github.com/frankill/gotools/date"
	"github.com/frankill/gotools/operation"
	"github.com/frankill/gotools/slider"
)

func TestSum(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	sum := slider.Sum(1, 1, 0, data)
	expectedSum := []float64{3, 6, 9, 12, 9}

	if array.All(gotools.Identity[bool], operation.Eq(sum, expectedSum)) {
		t.Log("测试通过")
	} else {
		t.Error("测试失败")
	}
}

func TestMax(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	max := slider.Max(1, 1, 0, data)
	expectedMax := []float64{2, 3, 4, 5, 5}
	if array.All(gotools.Identity[bool], operation.Eq(max, expectedMax)) {
		t.Log("测试通过")
	} else {
		t.Error("测试失败")
	}
}

func TestMin(t *testing.T) {
	data := []float64{5, 4, 3, 2, 1}
	min := slider.Min(1, 1, 10, data)
	expectedMin := []float64{4, 3, 2, 1, 1}
	if array.All(gotools.Identity[bool], operation.Eq(min, expectedMin)) {
		t.Log("测试通过")
	} else {
		t.Error("测试失败")
	}
}

func TestMean(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	mean := slider.Mean(1, 1, 0, data)
	expectedMean := []float64{1, 2, 3, 4, 3}
	if array.All(gotools.Identity[bool], operation.Eq(mean, expectedMean)) {
		t.Log("测试通过")
	} else {
		t.Error("测试失败")
	}
}

func TestPaste(t *testing.T) {
	data := []string{"a", "b", "c", "d", "e"}
	paste := slider.Paste(1, 1, "", data)
	expectedPaste := []string{"ab", "abc", "bcd", "cde", "de"}
	if array.All(gotools.Identity[bool], operation.Eq(paste, expectedPaste)) {
		t.Log("测试通过")
	} else {
		t.Error("测试失败")
	}
}

func TestSlideIndex(t *testing.T) {
	index := date.Days(time.Now(), 1, 2, 3)
	data := []float64{1, 2, 3}
	slideIndex := slider.SlideIndex(func(x []float64) float64 {
		return array.Sum(x)
	}, 1, 1, index, data)
	// 需要根据实际情况编写期望的结果
	expectedSlideIndex := []float64{3, 6, 5}
	if array.All(gotools.Identity[bool], operation.Eq(slideIndex, expectedSlideIndex)) {
		t.Log("测试通过")
	} else {
		t.Error("测试失败")
	}
}

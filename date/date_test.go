package date_test

import (
	"testing"
	"time"

	"github.com/frankill/gotools/date"
)

func TestFloorYear(t *testing.T) {
	currentTime := time.Now()
	flooredTime := date.FloorYear(currentTime)
	// 验证时间是否被向下取整到年份的第一天
	if flooredTime.Year() != currentTime.Year() || flooredTime.Month() != time.January || flooredTime.Day() != 1 {
		t.Errorf("Expected floored time to be the first day of the year, got %v", flooredTime)
	}
}

func TestCeilingYear(t *testing.T) {
	currentTime := time.Now()
	ceilingTime := date.CeilingYear(currentTime)
	// 验证时间是否被向上取整到年份的最后一天
	if ceilingTime.Year() != currentTime.Year() || ceilingTime.Month() != time.December || ceilingTime.Day() != 31 {
		t.Errorf("Expected ceiling time to be the last day of the year, got %v", ceilingTime)
	}
}

func TestToUnix(t *testing.T) {
	dateStr := "2023-09-17 15:04:05"
	expectedTime, err := time.Parse("2006-01-02 15:04:05", dateStr)
	if err != nil {
		t.Fatal(err)
	}
	unixTimestamps := date.ToUnix([]string{dateStr}...)
	if len(unixTimestamps) != 1 {
		t.Fatalf("Expected one Unix timestamp, got %v", unixTimestamps)
	}
	if unixTimestamps[0] != expectedTime.Unix() {
		t.Errorf("Expected Unix timestamp %d, got %d", expectedTime.Unix(), unixTimestamps[0])
	}
}

func TestToStr(t *testing.T) {
	currentTime := time.Now()
	strTime := date.ToStr(currentTime)[0]
	expectedFormat := "2006-01-02 15:04:05"
	_, err := time.Parse(expectedFormat, strTime)
	if err != nil {
		t.Errorf("ToStr function failed to format time correctly: %v", err)
	}
}

func TestSub(t *testing.T) {
	t1 := time.Now()
	t2 := t1.Add(48 * time.Hour)
	days := date.Sub(t2, t1)
	if days != 2 {
		t.Errorf("Expected 2 days difference, got %f", days)
	}
}

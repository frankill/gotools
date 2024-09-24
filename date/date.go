package date

import (
	"fmt"
	"strings"
	"time"

	"github.com/frankill/gotools/fn"
)

// Floor 将时间向下取整到指定的时间单位
//
// 参数:
//   - t: 要向下取整的时间对象。
//   - unit: 字符串，指定要向下取整到的时间单位，如 "year", "month", "day" 等。
//   - weekStart: 整数，指定一周的开始是星期几，0 表示星期日，1 表示星期一，依此类推。
//
// 返回:
//   - 一个 time.Time 对象，表示向下取整后的时间。
func Floor(t time.Time, unit string, weekStart int) time.Time {
	switch strings.ToLower(unit) {
	case "years", "year", "Y":
		return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
	case "months", "month", "M":
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	case "days", "day", "D":
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	case "hours", "hour", "h":
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
	case "minutes", "minute", "m":
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
	case "seconds", "second", "s":
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
	case "weeks", "week", "w":
		// Calculate the number of days to subtract based on the week start day
		days := int(t.Weekday()) - weekStart
		if days < 0 {
			days += 7
		}
		return t.AddDate(0, 0, -days)
	default:
		panic(fmt.Sprintf("unsupported unit: %s", unit))
	}
}

// Floor 将时间向上取整到指定的时间单位
//
// 参数:
//   - t: 要向上取整的时间对象。
//   - unit: 字符串，指定要向下取整到的时间单位，如 "year", "month", "day" 等。
//   - weekStart: 整数，指定一周的开始是星期几，0 表示星期日，1 表示星期一，依此类推。
//
// 返回:
//   - 一个 time.Time 对象，表示向上取整后的时间。
func Ceiling(t time.Time, unit string, weekStart int) time.Time {
	switch strings.ToLower(unit) {
	case "years", "year":
		// 获取年份的最后一天
		year := t.Year()
		return time.Date(year, 12, 31, 23, 59, 59, 0, t.Location())
	case "months", "month":
		// 获取月份的最后一天
		month := t.Month()
		year := t.Year()
		// 计算该月份有多少天
		lastDayOfMonth := getLastDayOfMonth(year, month)
		return time.Date(year, month, lastDayOfMonth, 23, 59, 59, 0, t.Location())
	case "days", "day":
		// Days are handled at the current day, so we do nothing.
		return t
	case "hours", "hour", "h":
		if t.Minute() == 0 && t.Second() == 0 && t.Nanosecond() == 0 {
			return t
		}
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour()+1, 0, 0, 0, t.Location())
	case "minutes", "minute", "m":
		if t.Second() == 0 && t.Nanosecond() == 0 {
			return t
		}
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute()+1, 0, 0, t.Location())
	case "seconds", "second", "s":
		if t.Nanosecond() == 0 {
			return t
		}
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()+1, 0, t.Location())
	case "weeks", "week", "w":
		// Calculate the number of days to add based on the week start day
		days := (7 - int(t.Weekday()) + weekStart) % 7
		return t.AddDate(0, 0, days)
	default:
		panic(fmt.Sprintf("unsupported unit: %s", unit))
	}
}

// getLastDayOfMonth 返回指定年月的最后一天
func getLastDayOfMonth(year int, month time.Month) int {
	// 通过将日期设置为下个月的第0天来获取这个月的最后一天
	lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)
	return lastDay.Day()
}

// Zone 将时间转换为指定时区的时间
//
// 参数:
//   - d: 要转换的时间对象。
//   - zone: 字符串，指定要转换的时区。空字符串使用"Asia/Shanghai"
//
// 返回:
//   - 一个 time.Time 对象，表示转换后的时间。
func Zone(d time.Time, zone string) time.Time {

	if zone == "" {
		zone = "Asia/Shanghai"
	}
	_, err := time.LoadLocation(zone)
	if err != nil {
		return time.Time{}
	}

	return d.In(time.FixedZone(zone, 0))
}

// ToUnix 将日期字符串转换为 Unix 时间戳
//
// 参数:
//   - d: 要转换的日期字符串。"2006-01-02 15:04:05"
//
// 返回:
//   - 一个 int64，表示转换后的 Unix 时间戳。
func ToUnix(d string) int64 {
	t, _ := time.Parse("2006-01-02 15:04:05", d)
	return t.Unix()
}

// ToUnixMilli 将日期字符串转换为 Unix 时间戳
//
// 参数:
//   - d: 要转换的日期字符串。"2006-01-02 15:04:05"
//
// 返回:
//   - 一个 int64，表示转换后的 Unix 时间戳。
func ToUnixMilli(d string) int64 {
	t, _ := time.Parse("2006-01-02 15:04:05", d)
	return t.UnixMilli()
}

// ToUnixNano 将日期字符串转换为 Unix 时间戳
//
// 参数:
//   - d: 要转换的日期字符串。"2006-01-02 15:04:05"
//
// 返回:
//   - 一个 int64，表示转换后的 Unix 时间戳。
func ToUnixNano(d string) int64 {
	t, _ := time.Parse("2006-01-02 15:04:05", d)
	return t.UnixNano()
}

// ToTime 将日期字符串转换为 time.Time 对象
//
// 参数:
//   - d: 要转换的日期字符串。"2006-01-02 15:04:05"
//
// 返回:
//   - 一个 time.Time 对象，表示转换后的时间。
func ToTime(d string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", d)
	return t
}

// UnixToTime 将 Unix 时间戳转换为 time.Time 对象
//
// 参数:
//   - unix: 要转换的 Unix 时间戳。
//
// 返回:
//   - 一个 time.Time 对象，表示转换后的时间。
func UnixToTime(unix int64) time.Time {
	return time.Unix(unix, 0)
}

func Days(t time.Time, d ...int) []time.Time {

	return fn.Lapply(func(x int) time.Time {
		return t.AddDate(0, 0, x)
	}, d)
}

func Months(t time.Time, d ...int) []time.Time {
	return fn.Lapply(func(x int) time.Time {
		return t.AddDate(0, x, 0)
	}, d)
}

func Years(t time.Time, d ...int) []time.Time {
	return fn.Lapply(func(x int) time.Time {
		return t.AddDate(x, 0, 0)
	}, d)
}

func YMD(d string) time.Time {
	t, _ := time.Parse("2006-01-02", d)
	return t
}

func YMDHMS(d string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", d)
	return t
}

// Sub 计算两个时间相差多少天
func Sub(d1, d2 time.Time) int {
	return int(d1.Sub(d2).Hours() / 24)
}

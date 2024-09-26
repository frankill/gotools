package date

import (
	"fmt"
	"strings"
	"time"

	"github.com/frankill/gotools/fn"
)

// FloorMonth 将日期向下取整到月份
//
// 参数:
//   - t: 要向下取整的日期对象。
//
// 返回:
//   - 一个 time.Time 对象，表示向下取整后的日期。
func FloorMonth(t time.Time) time.Time {
	return floor(t, "month", 0)
}

// FloorYear 将日期向下取整到年份
//
// 参数:
//   - t: 要向下取整的日期对象。
//
// 返回:
//   - 一个 time.Time 对象，表示向下取整后的日期。
func FloorYear(t time.Time) time.Time {
	return floor(t, "year", 0)
}

// Floor 将时间向下取整到指定的时间单位
//
// 参数:
//   - unit: 字符串，指定要向下取整到的时间单位，如 "year", "month", "day" 等。
//   - weekStart: 整数，指定一周的开始是星期几，0 表示星期日，1 表示星期一，依此类推。
//   - t: 要向下取整的时间对象。
//
// 返回:
//   - 一个 time.Time 对象，表示向下取整后的时间。
func Floor(unit string, weekStart int) func(d ...time.Time) []time.Time {

	return func(d ...time.Time) []time.Time {

		if len(d) == 0 {
			return []time.Time{}
		}

		var res []time.Time
		for _, t := range d {
			res = append(res, floor(t, unit, weekStart))
		}
		return res
	}
}

func floor(t time.Time, unit string, weekStart int) time.Time {
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

// CeilingMonth 将日期向上取整到月份
//
// 参数:
//   - t: 要向上取整的日期对象。
//
// 返回:
//   - 一个 time.Time 对象，表示向上取整后的日期。
func CeilingMonth(t time.Time) time.Time {
	return ceiling(t, "month", 0)
}

// CeilingYear 将日期向上取整到年份
//
// 参数:
//   - t: 要向上取整的日期对象。
//
// 返回:
//   - 一个 time.Time 对象，表示向上取整后的日期。
func CeilingYear(t time.Time) time.Time {
	return ceiling(t, "year", 0)
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
func Ceiling(unit string, weekStart int) func(d ...time.Time) []time.Time {

	return func(d ...time.Time) []time.Time {

		if len(d) == 0 {
			return []time.Time{}
		}

		var res []time.Time
		for _, t := range d {
			res = append(res, ceiling(t, unit, weekStart))
		}
		return res
	}
}

func ceiling(t time.Time, unit string, weekStart int) time.Time {
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
//   - 一个 []time.Time 对象，表示转换后的时间。
func Zone(zone string) func(d ...time.Time) []time.Time {

	return func(d ...time.Time) []time.Time {

		return fn.Lapply(func(x time.Time) time.Time {
			if zone == "" {
				zone = "Asia/Shanghai"
			}
			return x.In(time.FixedZone(zone, 0))
		}, d)
	}

}

// ToUnix 将日期字符串转换为 Unix 时间戳
//
// 参数:
//   - d: 要转换的日期字符串。"2006-01-02 15:04:05"
//
// 返回:
//   - 一个 []int64，表示转换后的 Unix 时间戳。
func ToUnix(d ...string) []int64 {

	return fn.Lapply(func(x string) int64 {
		t, _ := time.Parse("2006-01-02 15:04:05", x)
		return t.Unix()
	}, d)
}

// ToUnix2 将日期转换为 Unix 时间戳
//
// 参数:
//   - d: 要转换的日期对象。
//
// 返回:
//   - 一个 []int64，表示转换后的 Unix 时间戳。
func ToUnix2(d ...time.Time) []int64 {
	return fn.Lapply(func(x time.Time) int64 {
		return x.Unix()
	}, d)
}

// ToUnixMilli 将日期字符串转换为 Unix 时间戳
//
// 参数:
//   - d: 要转换的日期字符串。"2006-01-02 15:04:05"
//
// 返回:
//   - 一个 []int64，表示转换后的 Unix 时间戳。
func ToUnixMilli(d ...string) []int64 {

	return fn.Lapply(func(x string) int64 {
		t, _ := time.Parse("2006-01-02 15:04:05", x)
		return t.UnixMilli()
	}, d)
}

// ToUnixNano 将日期字符串转换为 Unix 时间戳
//
// 参数:
//   - d: 要转换的日期字符串。"2006-01-02 15:04:05"
//
// 返回:
//   - 一个 []int64，表示转换后的 Unix 时间戳。
func ToUnixNano(d ...string) []int64 {

	return fn.Lapply(func(x string) int64 {
		t, _ := time.Parse("2006-01-02 15:04:05", x)
		return t.UnixNano()
	}, d)
}

// ToTime 将日期字符串转换为 time.Time 对象
//
// 参数:
//   - d: 要转换的日期字符串。"2006-01-02 15:04:05"
//
// 返回:
//   - 一个 []time.Time 对象，表示转换后的时间。
func ToTime(d ...string) []time.Time {

	return fn.Lapply(func(x string) time.Time {
		t, _ := time.Parse("2006-01-02 15:04:05", x)
		return t
	}, d)
}

// ToStr 将日期转换为 "2006-01-02 15:04:05" 格式的字符串
//
// 参数:
//   - d: 要转换的日期对象。
//
// 返回:
//   - 一个 []string，表示转换后的字符串。
func ToStr(d ...time.Time) []string {

	return fn.Lapply(func(x time.Time) string {
		return x.Format("2006-01-02 15:04:05")
	}, d)
}

// ToYMD 将日期字符串转换为 "20060102" 格式的字符串
//
// 参数:
//   - d: 要转换的日期字符串。"2006-01-02 15:04:05"
//
// 返回:
//   - 一个 []string，表示转换后的字符串。
func ToYMD(d ...time.Time) []string {

	return fn.Lapply(func(x time.Time) string {
		return x.Format("20060102")
	}, d)
}

// UnixToStr 将 Unix 时间戳转换为日期字符串
//
// 参数:
//   - unix: 要转换的 Unix 时间戳。
//
// 返回:
//   - 一个 []string，表示转换后的日期字符串。
func UnixToStr(unix ...int64) []string {

	return fn.Lapply(func(x int64) string {
		return time.Unix(x, 0).Format("2006-01-02 15:04:05")
	}, unix)
}

// UnixToTime 将 Unix 时间戳转换为 time.Time 对象
//
// 参数:
//   - unix: 要转换的 Unix 时间戳。
//
// 返回:
//   - 一个 []time.Time 对象，表示转换后的时间。
func UnixToTime(unix ...int64) []time.Time {

	return fn.Lapply(func(x int64) time.Time {
		return time.Unix(x, 0)
	}, unix)
}

// Days 生成指定天数的日期列表
//
// 参数:
//   - d: 要生成的天数。
//   - t: 要生成日期列表的起始日期。
//
// 返回:
//   - 一个 []time.Time 对象，表示生成的日期列表。
func Days(t time.Time, d ...int) []time.Time {

	return fn.Lapply(func(x int) time.Time {
		return t.AddDate(0, 0, x)
	}, d)
}

// Months 生成指定月数的日期列表
//
// 参数:
//   - d: 要生成的月数。
//   - t: 要生成日期列表的起始日期。
//
// 返回:
//   - 一个 []time.Time 对象，表示生成的日期列表。
func Months(t time.Time, d ...int) []time.Time {
	return fn.Lapply(func(x int) time.Time {
		return t.AddDate(0, x, 0)
	}, d)
}

// Years 生成指定年数的日期列表
//
// 参数:
//   - d: 要生成的年数。
//   - t: 要生成日期列表的起始日期。
//
// 返回:
//   - 一个 []time.Time 对象，表示生成的日期列表。
func Years(t time.Time, d ...int) []time.Time {
	return fn.Lapply(func(x int) time.Time {
		return t.AddDate(x, 0, 0)
	}, d)
}

// YMD 将日期字符串转换为 "2006-01-02" 格式的字符串
//
// 参数:
//   - d: 要转换的日期字符串。"2006-01-02 15:04:05"
//
// 返回:
//   - 一个 []time.Time 对象，表示转换后的时间。
func YMD(d ...string) []time.Time {

	return fn.Lapply(func(x string) time.Time {
		t, _ := time.Parse("2006-01-02", x)
		return t
	}, d)
}

// YMDHMS 将日期字符串转换为 "2006-01-02 15:04:05" 格式的字符串
//
// 参数:
//   - d: 要转换的日期字符串。"2006-01-02 15:04:05"
//
// 返回:
//   - 一个 []time.Time 对象，表示转换后的时间。
func YMDHMS(d ...string) []time.Time {

	return fn.Lapply(func(x string) time.Time {
		t, _ := time.Parse("2006-01-02 15:04:05", x)
		return t
	}, d)
}

// Sub 计算两个时间相差多少天
//
// 参数:
//   - d1: 第一个时间
//   - d2: 第二个时间
//
// 返回:
//   - int
func Sub(d1, d2 time.Time) float64 {
	return d1.Sub(d2).Hours() / 24
}

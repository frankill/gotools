package lubridate

import "time"

func MakeDate(year int, month int, day int) time.Time {

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).Add(-8 * time.Hour)

}

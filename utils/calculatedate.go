package utils

import "time"

// GetFirstDateOfMonth 获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

// GetLastDateOfMonth 获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

// GetFirstDateOfYear 获取传入的时间所在年份的第一天
func GetFirstDateOfYear(d time.Time) time.Time {
	return time.Date(d.Year(), time.January, 1, 0, 0, 0, 0, d.Location())
}

// GetLastDateOfYear 获取传入的时间所在年份的最后一天
func GetLastDateOfYear(d time.Time) time.Time {
	return GetFirstDateOfYear(d).AddDate(1, 0, -1)
}

// GetZeroTime 获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

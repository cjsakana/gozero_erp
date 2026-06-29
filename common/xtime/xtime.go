package xtime

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// EndOfDay 获取给定日期当天的最后一刻
func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, t.Location())
}

// MinutesToDecimalHours 分钟转小时（保留2位小数）
func MinutesToDecimalHours(minutes int) float64 {
	return float64(minutes) / 60.0
}

// TimeToDecimal HH:MM 格式的分钟转小时（保留2位小数）
func TimeToDecimal(timeStr string) (float64, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid time format")
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	return float64(hours) + float64(minutes)/60.0, nil
}

// GetTimestampOfToday 获取当天指定时间的时间戳
// hour: 小时（0-23）, minute: 分钟（0-59）, second: 秒（0-59）, isMilli: 是否返回毫秒级
func GetTimestampOfToday(hour, minute, second int, isMilli bool) int64 {
	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, 0, now.Location())

	if isMilli {
		return t.UnixMilli() // 毫秒级时间戳
	}
	return t.Unix() // 秒级时间戳
}

// GetStandardClockInTime 获取当天的标准上班时间（9:00）的时间戳
func GetStandardClockInTime(timestamp int64, loc *time.Location) int64 {
	t := time.Unix(timestamp, 0).In(loc)
	return time.Date(t.Year(), t.Month(), t.Day(), 9, 0, 0, 0, loc).Unix()
}

// GetStandardClockOutTime 获取当天的标准下班时间（18:00）的时间戳
func GetStandardClockOutTime(timestamp int64, loc *time.Location) int64 {
	t := time.Unix(timestamp, 0).In(loc)
	return time.Date(t.Year(), t.Month(), t.Day(), 18, 0, 0, 0, loc).Unix()
}

// IsLate 判断是否迟到（基于int64时间戳）
func IsLate(clockInTimestamp int64, loc *time.Location) bool {
	standardClockIn := GetStandardClockInTime(clockInTimestamp, loc)
	return clockInTimestamp > (standardClockIn + 10*60) // 10分钟迟到阈值
}

// IsEarly 判断是否早退（基于int64时间戳）
func IsEarly(clockOutTimestamp int64, loc *time.Location) bool {
	standardClockOut := GetStandardClockOutTime(clockOutTimestamp, loc)
	return clockOutTimestamp < (standardClockOut - 10*60) // 10分钟早退阈值
}

// IsLateShanghai 判断是否迟到（基于int64时间戳）上海时区
func IsLateShanghai(clockInTimestamp int64) bool {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return IsLate(clockInTimestamp, loc)
}

// IsEarlyShanghai 判断是否早退（基于int64时间戳）上海时区
func IsEarlyShanghai(clockOutTimestamp int64) bool {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return IsEarly(clockOutTimestamp, loc)
}

func IsZeroTime(t time.Time) bool {
	return t.IsZero() || t.Unix() == 0
}

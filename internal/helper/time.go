package helper

import "time"

func FormatTime(time time.Time) string {
	return time.Format("2006-01-02")
}

func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse("2006-01-02", timeStr)
}

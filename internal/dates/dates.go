package dates

import (
	"shifts-go/internal/helper"
	"time"
)

func ThisWeek() time.Time {
	today := Today()
	offset := (int(today.Weekday()) - int(time.Friday) + 7) % 7
	return today.AddDate(0, 0, -offset)
}
func NextWeek() time.Time { return ThisWeek().AddDate(0, 0, 7) }
func LastWeek() time.Time { return ThisWeek().AddDate(0, 0, -7) }

func Today() time.Time {
	now := time.Now()
	y, m, d := now.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, now.Location())
}

func Tomorrow() time.Time  { return Today().AddDate(0, 0, 1) }
func Yesterday() time.Time { return Today().AddDate(0, 0, -1) }

func Friday() time.Time    { return ThisWeek() }
func Saturday() time.Time  { return ThisWeek().AddDate(0, 0, 1) }
func Sunday() time.Time    { return ThisWeek().AddDate(0, 0, 2) }
func Monday() time.Time    { return ThisWeek().AddDate(0, 0, 3) }
func Tuesday() time.Time   { return ThisWeek().AddDate(0, 0, 4) }
func Wednesday() time.Time { return ThisWeek().AddDate(0, 0, 5) }
func Thursday() time.Time  { return ThisWeek().AddDate(0, 0, 6) }

func NFriday() time.Time    { return NextWeek() }
func NSaturday() time.Time  { return NextWeek().AddDate(0, 0, 1) }
func NSunday() time.Time    { return NextWeek().AddDate(0, 0, 2) }
func NMonday() time.Time    { return NextWeek().AddDate(0, 0, 3) }
func NTuesday() time.Time   { return NextWeek().AddDate(0, 0, 4) }
func NWednesday() time.Time { return NextWeek().AddDate(0, 0, 5) }
func NThursday() time.Time  { return NextWeek().AddDate(0, 0, 6) }

func LThursday() time.Time  { return ThisWeek().AddDate(0, 0, -1) }
func LWednesday() time.Time { return ThisWeek().AddDate(0, 0, -2) }
func LTuesday() time.Time   { return ThisWeek().AddDate(0, 0, -3) }
func LMonday() time.Time    { return ThisWeek().AddDate(0, 0, -4) }
func LSunday() time.Time    { return ThisWeek().AddDate(0, 0, -5) }
func LSaturday() time.Time  { return ThisWeek().AddDate(0, 0, -6) }
func LFriday() time.Time    { return ThisWeek().AddDate(0, 0, -7) }

func MakeIntoFriday(t time.Time) time.Time {
	offset := (int(t.Weekday()) - int(time.Friday) + 7) % 7
	return t.AddDate(0, 0, offset)
}

func ParseArg(dateString string) (time.Time, error) {
	return helper.ParseTime(dateString)
}

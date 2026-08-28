package dates

import (
	"strings"
	"time"
)

var periodAliases = map[string]func() time.Time{
	"thisweek": ThisWeek,
	"tweek":    ThisWeek,
	"tw":       ThisWeek,

	"today": Today,
	"tod":   Today,
	"td":    Today,

	"tomorrow": Tomorrow,
	"tom":      Tomorrow,

	"yesterday": Yesterday,
	"yes":       Yesterday,

	"friday": Friday,
	"fri":    Friday,
	"tfri":   Friday,

	"saturday": Saturday,
	"sat":      Saturday,
	"tsat":     Saturday,

	"sunday": Sunday,
	"sun":    Sunday,
	"tsun":   Sunday,

	"monday": Monday,
	"mon":    Monday,
	"tmon":   Monday,

	"tuesday": Tuesday,
	"tue":     Tuesday,
	"ttue":    Tuesday,

	"wednesday": Wednesday,
	"wed":       Wednesday,
	"twed":      Wednesday,

	"thursday": Thursday,
	"thu":      Thursday,
	"tthu":     Thursday,

	// next week
	"nextweek":   NextWeek,
	"nweek":      NextWeek,
	"nw":         NextWeek,
	"nfriday":    NFriday,
	"nfri":       NFriday,
	"nsaturday":  NSaturday,
	"nsat":       NSaturday,
	"nsunday":    NSunday,
	"nsun":       NSunday,
	"nmonday":    NMonday,
	"nmon":       NMonday,
	"ntuesday":   NTuesday,
	"ntue":       NTuesday,
	"nwednesday": NWednesday,
	"nwed":       NWednesday,
	"nthursday":  NThursday,
	"nthu":       NThursday,

	// last week
	"lastweek":   LastWeek,
	"lweek":      LastWeek,
	"lw":         LastWeek,
	"lthursday":  LThursday,
	"lthu":       LThursday,
	"lwednesday": LWednesday,
	"lwed":       LWednesday,
	"ltuesday":   LTuesday,
	"ltue":       LTuesday,
	"lmonday":    LMonday,
	"lmon":       LMonday,
	"lsunday":    LSunday,
	"lsun":       LSunday,
	"lsaturday":  LSaturday,
	"lsat":       LSaturday,
	"lfriday":    LFriday,
	"lfri":       LFriday,
}

func DefinePeriod(period string) (time.Time, error) {
	period = strings.ToLower(strings.TrimSpace(period))
	if fn, ok := periodAliases[period]; ok {
		return fn(), nil
	}

	return ParseArg(period)
}

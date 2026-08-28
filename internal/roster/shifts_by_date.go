package roster

import (
	"fmt"
	"shifts-go/internal/helper"
	"sort"
	"time"
)

func ShiftsByDate(t time.Time, emps []Employee) []EmployeeShifts {
	time := helper.FormatTime(t)
	fmt.Println(time)
	var res []EmployeeShifts

	for i := range emps {
		e := &emps[i]
		for date, shifts := range e.Shifts {
			if date != time {
				continue
			}

			var validShifts []Shift

			for _, s := range shifts {
				if s.Duration.DecimalDuration == 0 {
					continue
				}

				validShifts = append(validShifts, s)

				if len(validShifts) > 0 {
					res = append(res, EmployeeShifts{
						Employee: e, Shifts: validShifts,
					})
				}
			}
		}

	}
	sort.Slice(res, func(a, b int) bool {
		return res[a].Shifts[0].StartTime.OrderableTime < res[b].Shifts[0].StartTime.OrderableTime
	})
	return res
}

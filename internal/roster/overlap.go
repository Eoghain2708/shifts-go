package roster

import "fmt"

type ShiftOverlap struct {
	Date         string  `json:"date"`
	ShiftOne     string  `json:"shiftOne"`
	ShiftOneName string  `json:"shiftOneName"`
	ShiftTwo     string  `json:"shiftTwo"`
	ShiftTwoName string  `json:"shiftTwoName"`
	Overlap      float64 `json:"overlap"`
}

func FindShiftsInCommon(emp1, emp2 *Employee) (result []ShiftOverlap, message string, err error) {
	if len(emp1.Shifts) == 0 && len(emp2.Shifts) == 0 {
		return nil, "This rota has not been finished", nil
	}

	if len(emp1.Shifts) == 0 {
		return nil, fmt.Sprintf("%s has no shifts scheduled yet", emp1.Name), nil
	}

	if len(emp2.Shifts) == 0 {
		return nil, fmt.Sprintf("%s has no shifts scheduled yet", emp2.Name), nil
	}

	var res []ShiftOverlap
	for date, emp1Shifts := range emp1.Shifts {
		emp2Shifts, exists := emp2.Shifts[date]
		if !exists {
			continue
		}

		for i := range emp1Shifts {
			for j := range emp2Shifts {
				shift1 := &emp1Shifts[i]
				shift2 := &emp2Shifts[j]
				found, overlap := findOverlap(shift1, shift2)
				if !found {
					continue
				}

				res = append(res, ShiftOverlap{
					Date:         date,
					ShiftOne:     shift1.ShiftText.Time12hr,
					ShiftOneName: emp1.Name,
					ShiftTwo:     shift2.ShiftText.Time12hr,
					ShiftTwoName: emp2.Name,
					Overlap:      overlap,
				})
			}
		}
	}

	if len(res) == 0 {
		return nil, "These employees will not see each other this week", nil
	}

	return res, "", nil
}

func FindMostSeen(emp *Employee, emps []*Employee) (map[string][]ShiftOverlap, error) {
	var res = make(map[string][]ShiftOverlap)
	for _, otherEmp := range emps {
		if otherEmp.Name == emp.Name {
			continue
		}

		commonShifts, msg, err := FindShiftsInCommon(emp, otherEmp)
		if err != nil {
			return nil, err
		}

		if len(msg) > 0 {
			fmt.Println(msg)
			return nil, nil
		}

		res[otherEmp.Name] = commonShifts
	}

	return res, nil
}

func findOverlap(shift1, shift2 *Shift) (found bool, overlap float64) {
	s1Start := shift1.StartTime.OrderableTime
	s1End := shift1.EndTime.OrderableTime
	s2Start := shift2.StartTime.OrderableTime
	s2End := shift2.EndTime.OrderableTime

	if s1Start >= s2End || s2Start >= s1End {
		return false, 0
	}

	var start float64
	var end float64

	start = max(s1Start, s2Start)
	end = min(s1End, s2End)

	if start-end == 0 {
		return false, 0
	}

	return true, end - start
}

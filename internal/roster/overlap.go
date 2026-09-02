package roster

import (
	"errors"
	"fmt"
	"os"
	"shifts-go/cli/ui"
	"shifts-go/internal/helper"
	"sort"
	"strconv"
	"strings"

	"github.com/aquasecurity/table"
)

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
		return []ShiftOverlap{}, "This rota has not been finished", nil
	}

	if len(emp1.Shifts) == 0 {
		return []ShiftOverlap{}, fmt.Sprintf("%s has no shifts scheduled yet", emp1.Name), nil
	}

	if len(emp2.Shifts) == 0 {
		return []ShiftOverlap{}, fmt.Sprintf("%s has no shifts scheduled yet", emp2.Name), nil
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

func (so *ShiftOverlap) Print() {
	fmt.Printf("%s\n", ui.BoldLightCyan.Render(helper.FormatShiftDate(so.Date)))
	helper.Divide()
	fmt.Printf("%s: %s\n", ui.BoldLightGreen.Render(so.ShiftOneName), ui.BoldWhite.Render(so.ShiftOne))
	fmt.Printf("%s: %s\n", ui.BoldLightGreen.Render(so.ShiftTwoName), ui.BoldWhite.Render(so.ShiftTwo))
	overlap := fmt.Sprintf("%.2f hours", so.Overlap)
	fmt.Printf("%s: %s\n", ui.BoldLightGreen.Render("Overlap"), ui.BoldWhite.Render(overlap))
	helper.Divide()
}

type EmployeeOverlapResponse struct {
	Employee *Employee      `json:"employee"`
	Overlaps []ShiftOverlap `json:"overlaps"`
}

func FindMostSeen(emp *Employee, emps []Employee) ([]EmployeeOverlapResponse, error) {
	var res = []EmployeeOverlapResponse{}
	for i := range emps {
		otherEmp := &emps[i]
		if otherEmp.ID == emp.ID {
			continue
		}

		commonShifts, _, err := FindShiftsInCommon(emp, otherEmp)

		if err != nil {
			return nil, errors.New("Cannot find shifts in common")
		}

		eos := EmployeeOverlapResponse{Employee: otherEmp, Overlaps: commonShifts}
		res = append(res, eos)
	}

	sort.Slice(res, func(a, b int) bool {
		var hoursA, hoursB float64
		for _, overlap := range res[a].Overlaps {
			hoursA += overlap.Overlap
		}

		for _, overlap := range res[b].Overlaps {
			hoursB += overlap.Overlap
		}

		return hoursA > hoursB
	})
	return res, nil
}

func PrintMostSeenTable(eos []EmployeeOverlapResponse) {
	t := table.New(os.Stdout)
	t.SetHeaders("Employee", "Frequency", "Total hrs seen", "Dates")

	for _, response := range eos {
		var datesBuilder strings.Builder
		var hours float64

		for _, overlap := range response.Overlaps {
			hours += overlap.Overlap
			fmt.Fprintf(&datesBuilder, "%s\n", helper.FormatShiftDate(overlap.Date))
		}

		row := []string{
			ui.BoldLightCyan.Render(response.Employee.Name),
			renderShiftFrequencyColour(len(response.Overlaps)),
			strconv.FormatFloat(hours, 'f', 2, 64),
			datesBuilder.String(),
		}

		t.AddRow(row...)
	}

	t.Render()
}

// Render number in different colour depending on frequency. >= 3 green, 1-2 yellow, 0 red
func renderShiftFrequencyColour(sf int) string {
	if sf >= 3 {
		return ui.BoldLightGreen.Render(strconv.Itoa(sf))
	} else if sf == 1 || sf == 2 {
		return ui.BoldLightYellow.Render(strconv.Itoa(sf))
	} else {
		return ui.Red.Render(strconv.Itoa(sf))
	}
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

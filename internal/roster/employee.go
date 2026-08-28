package roster

import (
	"errors"
	"fmt"
	"shifts-go/cli/ui"
	"strings"

	"github.com/MarkusZoppelt/fuzzymatch"
)

const GenManagerCode = "1489343527823711206"
const SupervisorCode = "1489339758941806044"
const DutyManagerCode = "1489343542827589598"
const GeneralStaffCode = "1537027076506875352"

const EqualOver21Wage int32 = 1271
const Under21Wage int32 = 1085

type Printable interface {
	Print()
}

type Employee struct {
	JobCode   string             `json:"defaultJob"`
	Name      string             `json:"displayName"`
	Age       float32            `json:"age"`
	Shifts    map[string][]Shift `json:"shifts"`
	StartDate string             `json:"startDate"`
}

func (e *Employee) Print() {
	fmt.Printf("%s: %s", ui.BoldLightGreen.Render("Name"), e.Name)
	fmt.Printf("%s: %v", ui.BoldLightGreen.Render("Age"), e.Age)
	fmt.Printf("%s: %s", ui.BoldLightGreen.Render("Is a"), e.JobRole())
	fmt.Printf("%s: %s", ui.BoldLightGreen.Render("Start"), e.StartDate)
}

func (emp *Employee) HourlyWagePence() int32 {
	var wage int32

	if emp.Age < 21 {
		wage = 1085
	} else {
		wage = 1271
	}

	if emp.JobCode != GeneralStaffCode {
		wage += 40
	}

	return wage
}

func (emp *Employee) JobRole() string {
	if emp.JobCode == GenManagerCode {
		return "General Manager"
	}
	if emp.JobCode == SupervisorCode {
		return "Supervisor"
	}
	if emp.JobCode == DutyManagerCode {
		return "Duty Manager"
	}
	if emp.JobCode == GeneralStaffCode {
		return "General Staff"
	}
	return "N/A"
}

type EmployeeResponse struct {
	StartDate string     `json:"startDate"`
	EndDate   string     `json:"endDate"`
	Employees []Employee `json:"employees"`
}

func FuzzyFindEmployee(emps []Employee, empName string) (*Employee, error) {
	possibleMatches := []string{}
	for _, emp := range emps {
		possibleMatches = append(possibleMatches, emp.Name)
	}

	match := fuzzymatch.SuggestClosestMatch(empName, possibleMatches, 40)

	for i := range emps {
		if emps[i].Name == match {
			return &emps[i], nil
		}
	}

	return nil, errors.New("employee not found")
}

func FindEmployee(emps []Employee, empName string) (*Employee, error) {
	for _, emp := range emps {
		if strings.TrimSpace(strings.ToLower(emp.Name)) == strings.ToLower(strings.TrimSpace(empName)) {
			return &emp, nil
		}
	}

	return nil, errors.New("employee not found")
}

func FindShifts(emps []Employee, empName string) (map[string][]Shift, error) {
	emp, err := FuzzyFindEmployee(emps, empName)
	if err != nil {
		return nil, err
	}

	return emp.Shifts, nil
}

type EmployeeShifts struct {
	Employee *Employee
	Shifts   []Shift
}

func (es *EmployeeShifts) Print() {
	fmt.Printf("%s: %v\n", ui.BoldGreen.Render("Name"), ui.BoldLightCyan.Render(es.Employee.Name))
	for _, s := range es.Shifts {
		fmt.Printf("%s: %v\n", ui.BoldGreen.Render("Shift"), s.ShiftText.Time12hr)
		fmt.Printf("%s: %v\n", ui.BoldGreen.Render("Hours"), s.Duration.DecimalDuration)
		fmt.Printf("%s: %v\n", ui.BoldGreen.Render("Pay"), (float64(s.Payment(es.Employee)) / 100))
		fmt.Println(ui.BoldWhite.Render(strings.Repeat("-", 40)))
	}
}

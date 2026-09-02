package roster

import (
	"errors"
	"fmt"
	"shifts-go/cli/ui"
	"shifts-go/internal/helper"
	"strings"

	"github.com/hbollon/go-edlib"
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
	ID        string             `json:"id"`
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
	empName = helper.NormaliseName(empName)

	possibleMatches := []string{}

	for i := range emps {
		name := helper.NormaliseName(emps[i].Name)
		if name == empName {
			return &emps[i], nil
		}

		possibleMatches = append(possibleMatches, name)
	}

	res, err := edlib.FuzzySearchThreshold(empName, possibleMatches, 0.18, edlib.JaroWinkler)
	if err != nil {
		return nil, fmt.Errorf("Cannot find employee %s", empName)
	}

	for i := range emps {
		if res == helper.NormaliseName(emps[i].Name) {
			return &emps[i], nil
		}
	}

	return nil, fmt.Errorf("employee matching %s not found", empName)
}

func FindEmployee(emps []Employee, empName string) (*Employee, error) {
	for _, emp := range emps {
		if helper.NormaliseName(emp.Name) == helper.NormaliseName(empName) {
			return &emp, nil
		}
	}

	return nil, errors.New("employee not found")
}

func FindShifts(emps []Employee, e *Employee) (map[string][]Shift, error) {
	return e.Shifts, nil
}

type EmployeeShifts struct {
	Employee *Employee `json:"employee"`
	Shifts   []Shift   `json:"shifts"`
}

func (es *EmployeeShifts) Print() {
	fmt.Printf("%-6s %v\n", ui.BoldGreen.Render("Name:"), ui.BoldLightCyan.Render(es.Employee.Name))

	for _, s := range es.Shifts {

		fmt.Printf("%-6s %v\n", ui.BoldGreen.Render("Shift:"), ui.BoldWhite.Render(s.ShiftText.Time12hr))

		decimalDuration := fmt.Sprintf("%.2f", s.Duration.DecimalDuration)
		fmt.Printf("%-6s %v\n", ui.BoldGreen.Render("Hours:"), ui.BoldWhite.Render(decimalDuration))

		pay := fmt.Sprintf("%.2f", float64(s.Payment(es.Employee))/100)
		fmt.Printf("%-6s %s%s\n", ui.BoldGreen.Render("Pay:"), ui.BoldWhite.Render("£"), ui.BoldWhite.Render(pay))

		fmt.Println(ui.BoldWhite.Render(strings.Repeat("-", 40)))
	}
}

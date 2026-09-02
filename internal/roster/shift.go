package roster

import (
	"fmt"
	"shifts-go/cli/ui"
	"shifts-go/internal/helper"
	"time"
)

type Shift struct {
	Role      JobRole      `json:"job"`
	ShiftText ShiftText    `json:"shiftText"`
	StartTime StartEndTime `json:"startTime"`
	EndTime   StartEndTime `json:"endTime"`
	Duration  NetDuration  `json:"netDuration"`
	ShiftDate string       `json:"shiftDate"`
}

type ShiftText struct {
	Time12hr string `json:"time12Hr"`
	ToClose  bool   `json:"toClose"`
}

type StartEndTime struct {
	OrderableTime float64 `json:"orderableTime"`
	Display12     string  `json:"display12"`
}

type NetDuration struct {
	DecimalDuration float64 `json:"decimal"`
}

type JobRole struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (s *Shift) Payment(emp *Employee) int32 {
	hourlyWage := s.calculateHourlyWage(emp)
	return int32(float64(hourlyWage) * s.Duration.DecimalDuration)
}

func FormattedPayment(wage int32) string {
	return fmt.Sprintf("%2.f", float64(wage)/100)
}

func (s *Shift) calculateHourlyWage(emp *Employee) int32 {
	var wage int32
	if emp.Age > 21 {
		wage = EqualOver21Wage
	} else {
		wage = Under21Wage
	}

	if s.Role.ID != GeneralStaffCode {
		wage += 40
	}

	return wage
}

func (s *Shift) Print(e *Employee) {
	t, err := time.Parse("2006-01-02", s.ShiftDate)
	formatted := t.Format("Monday 2 January, 2006")
	if err != nil {
		fmt.Println(err)
		return
	}
	helper.Divide()
	fmt.Printf("%s: %-10s\n", ui.BoldLightGreen.Render("Day"), ui.BoldWhite.Render(formatted))
	fmt.Printf("%s: %-10s\n", ui.BoldLightGreen.Render("Shift"), ui.BoldWhite.Render(s.ShiftText.Time12hr))
	duration := fmt.Sprintf("%.2f hrs", s.Duration.DecimalDuration)
	fmt.Printf("%s: %-10s\n", ui.BoldLightGreen.Render("Duration"), ui.BoldWhite.Render(duration))
	payment := fmt.Sprintf("%.2f", float64(s.Payment(e))/100)
	fmt.Printf("%s: %-10v\n", ui.BoldLightGreen.Render("Pay"), ui.BoldWhite.Render("£"+payment))
	helper.Divide()
}

package roster

import (
	"fmt"
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

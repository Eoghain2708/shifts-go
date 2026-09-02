package roster

import (
	"encoding/json/v2"
	"fmt"
	"os"
	"shifts-go/cli/ui"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/aquasecurity/table"
)

type LifetimeOverlap struct {
	OtherEmployee    *Employee `json:"otherEmployee"`
	TimesSeen        int       `json:"timesSeen"`
	OverlappingHours float64   `json:"overlappingHours"`
}

type Report struct {
	Employee      *Employee          `json:"employee"`
	Overlaps      []*LifetimeOverlap `json:"lifetimeOverlaps"`
	TotalHours    float64            `json:"totalHours"`
	TotalPayPence int32              `json:"totalPayPence"`
}

type ProcessedWeek struct {
	Hours            float64
	PayPence         int32
	EmployeeOverlaps []EmployeeOverlapResponse
}

type WeekResult struct {
	Report *ProcessedWeek
	Err    error
}

func processWeek(empName string, date time.Time, c *Client) (*ProcessedWeek, error) {
	r, err := c.GetEmployees(date)
	if err != nil {
		return nil, err
	}

	weekEmp, err := FindEmployee(r.Employees, empName)
	if err != nil {
		return nil, nil
	}

	var hours float64
	var payPence int32

	empOverlaps, err := FindMostSeen(weekEmp, r.Employees)
	if err != nil {
		return nil, err
	}

	for _, shifts := range weekEmp.Shifts {
		for _, s := range shifts {
			hours += s.Duration.DecimalDuration
			payPence += s.Payment(weekEmp)
		}
	}

	return &ProcessedWeek{Hours: hours, PayPence: payPence, EmployeeOverlaps: empOverlaps}, nil
}

func GetCachedReport(filename string) (*Report, error) {
	contents, exists, err := ReadLifetime(filename)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("Cannot find report by name %v", filename)
	}
	var report Report
	err = json.Unmarshal(contents, &report)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func PrintReport(r *Report) {
	fmt.Printf("Report for %s\n", r.Employee.Name)

	for _, overlap := range r.Overlaps {
		fmt.Printf(
			"%s - %d (%.2f hours)\n",
			overlap.OtherEmployee.Name,
			overlap.TimesSeen,
			overlap.OverlappingHours,
		)
	}

	fmt.Printf("Total hours - %.2f\n", r.TotalHours)
	fmt.Printf("Total wages - %.2f\n", float64(r.TotalPayPence)/100)
}

func PrintReportTable(r *Report) {
	t := table.New(os.Stdout)
	t.AddHeaders("Employee", "Times Seen", "Hours Together")
	for _, overlap := range r.Overlaps {
		t.AddRow(
			ui.BoldLightCyan.Render(overlap.OtherEmployee.Name),
			ui.BoldWhite.Render(strconv.Itoa(overlap.TimesSeen)),
			ui.BoldWhite.Render(strconv.FormatFloat(overlap.OverlappingHours, 'f', 2, 64)))
	}
	pay := fmt.Sprintf("%.2f", float32(r.TotalPayPence)/100)
	t.AddFooters("", "Total Pay", ui.BoldLightGreen.Render(pay))
	hours := fmt.Sprintf("%.2f", r.TotalHours)
	t.AddFooters("", "Total Hours", ui.BoldWhite.Render(hours))

	t.Render()
}

func GenerateReport(empName string, from time.Time, c *Client, emps []Employee) (*Report, error) {
	emp, err := FindEmployee(emps, empName)
	if err != nil {
		return nil, fmt.Errorf("employee could not be found in this week. Please select a week where they were in")
	}

	to, err := time.Parse("2006-01-02", emp.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid employee start date: %w", err)
	}

	var lifeHours float64
	var lifePayPence int32

	lifetimeOverlaps := make(map[string]*LifetimeOverlap)
	weekCount := 0
	for d := from; !d.Before(to); d = d.AddDate(0, 0, -7) {
		weekCount++
	}

	results := make(chan WeekResult, weekCount)
	sem := make(chan struct{}, 5)
	var wg sync.WaitGroup

	for ; !from.Before(to); from = from.AddDate(0, 0, -7) {
		formatted := from.Format("January 2, 2006")
		fmt.Printf("Processing week: %s\n", formatted)
		date := from
		wg.Add(1)
		go func() {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			processedWeek, err := processWeek(empName, date, c)
			results <- WeekResult{
				Report: processedWeek,
				Err:    err,
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Println("Generating report table...")

	for result := range results {
		if result.Err != nil {
			return nil, result.Err
		}

		if result.Report == nil {
			continue
		}

		lifeHours += result.Report.Hours
		lifePayPence += result.Report.PayPence

		for _, eo := range result.Report.EmployeeOverlaps {
			id := eo.Employee.ID
			lifetimeOverlap, exists := lifetimeOverlaps[id]
			if !exists {
				lifetimeOverlap = &LifetimeOverlap{
					OtherEmployee: eo.Employee,
				}

				lifetimeOverlaps[id] = lifetimeOverlap
			}
			lifetimeOverlap.TimesSeen += len(eo.Overlaps)

			for _, overlap := range eo.Overlaps {
				lifetimeOverlap.OverlappingHours += overlap.Overlap
			}
		}
	}

	overlaps := make([]*LifetimeOverlap, 0, len(lifetimeOverlaps))
	for _, overlap := range lifetimeOverlaps {
		overlaps = append(overlaps, overlap)
	}

	sort.Slice(overlaps, func(a, b int) bool {
		if overlaps[a].OverlappingHours == overlaps[b].OverlappingHours {
			return overlaps[a].TimesSeen > overlaps[b].TimesSeen
		}
		return overlaps[a].OverlappingHours > overlaps[b].OverlappingHours
	})

	report := Report{
		Employee:      emp,
		Overlaps:      overlaps,
		TotalHours:    lifeHours,
		TotalPayPence: lifePayPence,
	}

	if err := WriteLifetimeReport(report.Employee.Name, report); err != nil {
		return nil, err
	}

	return &report, nil
}

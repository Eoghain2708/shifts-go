/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"shifts-go/cli/ui"
	"shifts-go/internal/dates"
	"shifts-go/internal/roster"

	"github.com/spf13/cobra"
)

// hoursCmd represents the hours command
var hoursCmd = &cobra.Command{
	Use:     "hours <EMPLOYEE_NAME> <DATE>",
	Aliases: []string{"hrs", "h"},
	Short:   "Get the hours of a single employee for the week",
	Long: `Get the hours, payment, and each shift of a single employee for the week 
	along with total pre-tax pay.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("Too few arguments. Include date and employee name")
		}
		date, err := dates.DefinePeriod(args[1])
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		r, err := client.GetEmployees(date)
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		emp, err := roster.FuzzyFindEmployee(r.Employees, args[0])
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		data, err := roster.FindShifts(r.Employees, emp)
		if err != nil {
			return fmt.Errorf("Cannot get employee: %v", err)
		}

		formatted := date.Format("Monday 2 January 2006")
		var hours float64
		var totalPay int32
		hourlyWageStr := fmt.Sprintf("£%.2f", float64(emp.HourlyWagePence())/100)
		fmt.Printf("%s: %s\n", ui.BoldWhite.Render("Shifts for week"), ui.BoldLightCyan.Render(formatted))
		fmt.Printf("%s - %s (%s)\n", ui.BoldLightCyan.Render(emp.Name), ui.BoldLightGreen.Render(emp.JobRole()), ui.BoldGreen.Render(hourlyWageStr+"/hr"))
		for _, shifts := range data {
			for _, shift := range shifts {
				hours += shift.Duration.DecimalDuration
				totalPay += shift.Payment(emp)
				shift.Print(emp)
			}
		}
		hoursStr := fmt.Sprintf("%.2f", hours)
		payStr := fmt.Sprintf("£%.2f", float64(totalPay)/100)
		fmt.Printf("%s: %s\n", ui.BoldLightGreen.Render("Total hours"), ui.BoldWhite.Render(hoursStr))
		fmt.Printf("%s: %s\n", ui.BoldLightGreen.Render("Total pay"), ui.BoldWhite.Render(payStr))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(hoursCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hoursCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hoursCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

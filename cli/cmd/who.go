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

// whoCmd represents the who command
var whoCmd = &cobra.Command{
	Use:   "who <date-string> in format YYYY-MM-DD or else a shorthand command",
	Short: "See who is in work on a given day",
	Long: `Receive a list of each person working on a given day, ordered by shift start time, and showing
	 Name, Shift, Total hours worked, and Total payment for the day.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		date, err := dates.DefinePeriod(args[0])
		if err != nil {
			return fmt.Errorf("Invalid date")
		}
		r, err := client.GetEmployees(date)
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		namesAndShifts := roster.ShiftsByDate(date, r.Employees)
		fmt.Printf("%s: %v\n", ui.BoldWhite.Render("Shifts for"), ui.BoldLightCyan.Render(date.Format("Monday 2 January 2006")))
		for _, shiftData := range namesAndShifts {
			shiftData.Print()
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(whoCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// whoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// whoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

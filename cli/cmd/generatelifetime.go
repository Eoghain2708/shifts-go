/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"shifts-go/internal/dates"
	"shifts-go/internal/helper"
	"shifts-go/internal/roster"

	"github.com/hbollon/go-edlib"
	"github.com/spf13/cobra"
)

// generatelifetimeCmd represents the generatelifetime command
var generatelifetimeCmd = &cobra.Command{
	Use:     "generatelifetime <EMPLOYEE_NAME> <DATE> where <DATE> represents the starting point counting back to employee's start date",
	Aliases: []string{"glifetime", "glife", "gl"},
	Short:   "Generate a lifetime report for an employee",
	Long: `Generate and cache a lifetime report for an employee including how many hours they worked total, how much 
	they were paid before tax (approximately), and how many times they saw and hours they spent with each employee`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("Too few arguments. Include date and employee name")
		}
		date, err := dates.DefinePeriod(args[1])
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		rota, err := client.GetEmployees(date)
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		var names []string
		for _, emp := range rota.Employees {
			names = append(names, helper.NormaliseName(emp.Name))
		}

		// using explicity fuzzysearch instead of wrapper method for modified strictness
		empName, err := edlib.FuzzySearchThreshold(args[0], names, 0.75, edlib.JaroWinkler)
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		report, err := roster.GenerateReportConcurrent(empName, date, client, rota.Employees)
		if err != nil {
			return errors.New("Error creating report")
		}

		roster.PrintReportTable(report)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generatelifetimeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generatelifetimeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generatelifetimeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

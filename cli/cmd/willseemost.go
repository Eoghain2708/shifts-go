/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"shifts-go/internal/dates"
	"shifts-go/internal/roster"

	"github.com/spf13/cobra"
)

// willseemostCmd represents the willseemost command
var willseemostCmd = &cobra.Command{
	Use:     "willseemost <EMPLOYEE_NAME> <DATE>",
	Aliases: []string{"wsm"},
	Short:   "See who a given employee will see most of in a given week",
	Long: `Receive a list in tabular form ordered by which how many hours the given employee will spend with each of their colleagues
	throughout the week, alongside the dates and frequencies.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("Invalid format. Call using: willseemost <EMPLOYEE_NAME> <DATE>")
		}
		date, err := dates.DefinePeriod(args[1])
		if err != nil {
			return fmt.Errorf("Invalid date")
		}

		r, err := client.GetEmployees(date)
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		emp, err := roster.FuzzyFindEmployee(r.Employees, args[0])
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		allOverlaps, err := roster.FindMostSeen(emp, r.Employees)
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		roster.PrintMostSeenTable(allOverlaps)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(willseemostCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// willseemostCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// willseemostCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

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

// willseeCmd represents the willsee command
var willseeCmd = &cobra.Command{
	Use:   "willsee <NAME_ONE> <NAME_TWO> <DATE>",
	Short: "Enter the names of two employees and receive an output of any overlapping shifts they will have throughout the week",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		date, err := dates.DefinePeriod(args[2])
		if err != nil {
			return fmt.Errorf("Invalid date")
		}

		r, err := client.GetEmployees(date)
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		emp1, err := roster.FuzzyFindEmployee(r.Employees, args[0])
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		emp2, err := roster.FuzzyFindEmployee(r.Employees, args[1])
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		commonShifts, message, err := roster.FindShiftsInCommon(emp1, emp2)
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		if len([]rune(message)) > 0 {
			fmt.Println(message)
		}

		for _, overlap := range commonShifts {
			overlap.Print()
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(willseeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// willseeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// willseeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

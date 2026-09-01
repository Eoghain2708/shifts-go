/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"shifts-go/internal/helper"
	"shifts-go/internal/roster"

	"github.com/hbollon/go-edlib"
	"github.com/spf13/cobra"
)

// lifetimeCmd represents the lifetime command
var lifetimeCmd = &cobra.Command{
	Use:     "lifetime <EMPLOYEE_NAME>",
	Aliases: []string{"life, l"},
	Short:   "Return a cached lifetime report for an employee if there is one saved",
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := os.ReadDir(roster.LifeTimeReportsDir())
		if err != nil {
			return err
		}

		var names []string
		for _, e := range entries {
			names = append(names, helper.NormaliseName(e.Name()))
		}

		empName, err := edlib.FuzzySearchThreshold(args[0], names, 0.75, edlib.JaroWinkler)
		if err != nil {
			return err
		}
		report, err := roster.GetCachedReport(empName)
		if err != nil {
			return err
		}

		roster.PrintReportTable(report)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(lifetimeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lifetimeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lifetimeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// listDoneCmd represents the listDone command
var listDoneCmd = &cobra.Command{
	Use:   "list-done",
	Short: "List completed tasks",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := queries.ListDoneTasks(ctx)
		cobra.CheckErr(err)
		displayTasksTable(tasks)
	},
}

func init() {
	rootCmd.AddCommand(listDoneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listDoneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listDoneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

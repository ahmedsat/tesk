/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// archivedCmd represents the archived command
var archivedCmd = &cobra.Command{
	Use:   "archived",
	Short: "List archived tasks",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := queries.ListArchivedTasks(ctx)
		cobra.CheckErr(err)
		displayTasksTable(tasks)
	},
}

func init() {
	rootCmd.AddCommand(archivedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// archivedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// archivedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

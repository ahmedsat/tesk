/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	sqlc "github.com/ahmedsat/tesk/sql"
	"github.com/spf13/cobra"
)

var listPage, listSize *int

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		offset := (*listPage - 1) * *listSize
		tasks, err := queries.ListTasksPage(ctx, sqlc.ListTasksPageParams{Limit: int64(*listSize), Offset: int64(offset)})
		cobra.CheckErr(err)
		displayTasksTable(tasks)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listPage = listCmd.Flags().IntP("page", "p", 1, "Page number")
	listSize = listCmd.Flags().IntP("size", "s", 10, "Page size")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

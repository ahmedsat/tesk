/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"slices"

	sqlc "github.com/ahmedsat/tesk/sql"
	"github.com/spf13/cobra"
)

func assert(cond bool) {
	if !cond {
		panic("assertion failed")
	}
}

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search tasks",
	Run: func(cmd *cobra.Command, args []string) {
		q := cmd.Flag("query").Value.String()
		assert(q != "")
		titleRes, err := queries.SearchTasksByTitle(ctx, sql.NullString{Valid: true, String: q})
		cobra.CheckErr(err)
		descRes, err := queries.SearchTasksByDescription(ctx, sql.NullString{Valid: true, String: q})
		cobra.CheckErr(err)
		all := append(titleRes, descRes...)
		slices.SortFunc(all, func(a, b sqlc.Task) int {
			return int(a.ID - b.ID)
		})
		displayTasksTable(slices.Compact(all))
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringP("query", "q", "", "Search term (required)")
	searchCmd.MarkFlagRequired("query")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"strconv"

	sqlc "github.com/ahmedsat/tesk/sql"
	"github.com/spf13/cobra"
)

func ternary[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

// update [id]Cmd represents the update [id] command
var updateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title, desc := cmd.Flag("title").Value.String(), cmd.Flag("description").Value.String()
		id, _ := strconv.ParseInt(args[0], 10, 64)
		oldTask, err := queries.GetTask(ctx, id)
		cobra.CheckErr(err)

		_, err = queries.UpdateTask(ctx, sqlc.UpdateTaskParams{
			ID:          id,
			Title:       ternary(title == "", oldTask.Title, title),
			Description: ternary(desc == "", oldTask.Description, sql.NullString{Valid: true, String: desc}),
		})
		cobra.CheckErr(err)
		fmt.Printf("✅ Updated task %d\n", id)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringP("title", "t", "", "Task title (required)")
	updateCmd.Flags().StringP("description", "d", "", "Task description")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// update [id]Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// update [id]Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

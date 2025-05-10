/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"

	sqlc "github.com/ahmedsat/tesk/sql"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new task",
	Run: func(cmd *cobra.Command, args []string) {
		title, desc := cmd.Flag("title"), cmd.Flag("description")
		fmt.Println(title.Value.String(), desc.Value.String(), desc.Changed)
		task, err := queries.CreateTask(ctx, sqlc.CreateTaskParams{
			Title:       title.Value.String(),
			Description: sql.NullString{Valid: desc.Changed, String: desc.Value.String()},
		})
		cobra.CheckErr(err)
		fmt.Printf("✅ Created task %d\n", task.ID)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// cmd.Flags().StringVarP(&title, "title", "t", "", "Task title (required)")
	// cmd.Flags().StringVarP(&desc, "description", "d", "", "Task description")

	createCmd.Flags().StringP("title", "t", "", "Task title (required)")
	createCmd.MarkFlagRequired("title")
	createCmd.Flags().StringP("description", "d", "", "Task description")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

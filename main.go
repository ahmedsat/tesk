package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	sqlc "github.com/ahmedsat/tesks/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var debug bool = true

//go:generate sqlc generate

// Embed migration files
//
//go:embed migrations
var migrations embed.FS

var (
	ctx     = context.Background()
	db      *sql.DB
	queries *sqlc.Queries
)

func main() {

	rootCmd := &cobra.Command{
		Use:   "tesks",
		Short: "Task management CLI",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setupDatabase()
			runMigrations()
			cleanupTasks()
		},
	}

	// Register subcommands
	rootCmd.AddCommand(
		newCreateCmd(),
		newListCmd(),
		newGetCmd(),
		newUpdateCmd(),
		newDeleteCmd(),
		newSearchCmd(),
		newDoneCmd(),
		newRestoreCmd(),
		newArchivedCmd(),
		newMarkDoneCmd(),
	)

	cobra.CheckErr(rootCmd.Execute())
}

func setupDatabase() {
	userHome := os.Getenv("HOME")
	if userHome == "" {
		exitErr("HOME environment variable not set")
	}

	dataDir := filepath.Join(userHome, ".local", "share", "tesks")
	cobra.CheckErr(os.MkdirAll(dataDir, os.ModePerm))

	if debug {
		dataDir = "data"
		err := os.MkdirAll(dataDir, os.ModePerm)
		cobra.CheckErr(err)
	}

	var err error
	db, err = sql.Open("sqlite3", filepath.Join(dataDir, "tesks.db"))
	cobra.CheckErr(err)
	queries = sqlc.New(db)

}

func runMigrations() {

	files, err := migrations.ReadDir("migrations")
	cobra.CheckErr(err)
	slices.SortFunc(files, func(a, b fs.DirEntry) int {
		return strings.Compare(a.Name(), b.Name())
	})
	for _, f := range files {

		b, err := migrations.ReadFile("migrations/" + f.Name())
		cobra.CheckErr(err)
		_, err = db.ExecContext(ctx, string(b))
		cobra.CheckErr(err)
	}

}

func cleanupTasks() {
	cobra.CheckErr(queries.DeleteOldTasks(ctx))
	cobra.CheckErr(queries.DeleteOlderTasks(ctx))

}

func exitErr(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

// ---- Commands ----
func newCreateCmd() *cobra.Command {
	var title, desc string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new task",
		Run: func(cmd *cobra.Command, args []string) {
			if title == "" {
				cmd.Usage()
				os.Exit(1)
			}
			task, err := queries.CreateTask(ctx, sqlc.CreateTaskParams{
				Title:       title,
				Description: sql.NullString{Valid: desc != "", String: desc},
			})
			cobra.CheckErr(err)
			fmt.Printf("✅ Created task %d\n", task.ID)
		},
	}
	cmd.Flags().StringVarP(&title, "title", "t", "", "Task title (required)")
	cmd.Flags().StringVarP(&desc, "description", "d", "", "Task description")
	return cmd
}

func newListCmd() *cobra.Command {
	var page, size int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List active tasks",
		Run: func(cmd *cobra.Command, args []string) {
			offset := (page - 1) * size
			tasks, err := queries.ListTasksPage(ctx, sqlc.ListTasksPageParams{Limit: int64(size), Offset: int64(offset)})
			cobra.CheckErr(err)
			displayTasksTable(tasks)
		},
	}
	cmd.Flags().IntVarP(&page, "page", "p", 1, "Page number")
	cmd.Flags().IntVarP(&size, "size", "s", 10, "Page size")
	return cmd
}

func newGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get [id]",
		Short: "Get task by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.ParseInt(args[0], 10, 64)
			task, err := queries.GetTask(ctx, id)
			cobra.CheckErr(err)
			displayTasksTable([]sqlc.Task{task})
		},
	}
}

func newUpdateCmd() *cobra.Command {
	var title, desc string
	cmd := &cobra.Command{
		Use:   "update [id]",
		Short: "Update a task",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.ParseInt(args[0], 10, 64)
			_, err := queries.UpdateTask(ctx, sqlc.UpdateTaskParams{
				ID:          id,
				Title:       title,
				Description: sql.NullString{Valid: desc != "", String: desc},
			})
			cobra.CheckErr(err)
			fmt.Printf("✅ Updated task %d\n", id)
		},
	}
	cmd.Flags().StringVarP(&title, "title", "t", "", "New title")
	cmd.Flags().StringVarP(&desc, "description", "d", "", "New description")
	return cmd
}

func newDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete [id]",
		Short: "Archive a task",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.ParseInt(args[0], 10, 64)
			_, err := queries.DeleteTask(ctx, id)
			cobra.CheckErr(err)
			fmt.Printf("🗑️ Archived task %d\n", id)
		},
	}
}

func newSearchCmd() *cobra.Command {
	var q string
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search tasks",
		Run: func(cmd *cobra.Command, args []string) {
			titleRes, err := queries.SearchTasksByTitle(ctx, sql.NullString{Valid: q != "", String: q})
			cobra.CheckErr(err)
			descRes, err := queries.SearchTasksByDescription(ctx, sql.NullString{Valid: q != "", String: q})
			cobra.CheckErr(err)
			all := append(titleRes, descRes...)
			displayTasksTable(all)
		},
	}
	cmd.Flags().StringVarP(&q, "query", "q", "", "Search term (required)")
	cmd.MarkFlagRequired("query")
	return cmd
}

func newDoneCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "done",
		Short: "List completed tasks",
		Run: func(cmd *cobra.Command, args []string) {
			tasks, err := queries.ListDoneTasks(ctx)
			cobra.CheckErr(err)
			displayTasksTable(tasks)
		},
	}
}

func newRestoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restore [id]",
		Short: "Restore a completed task",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.ParseInt(args[0], 10, 64)
			_, err := queries.RestoreTask(ctx, id)
			cobra.CheckErr(err)
			fmt.Printf("♻️ Restored task %d\n", id)
		},
	}
}

func newArchivedCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "archived",
		Short: "List archived tasks",
		Run: func(cmd *cobra.Command, args []string) {
			tasks, err := queries.ListArchivedTasks(ctx)
			cobra.CheckErr(err)
			displayTasksTable(tasks)
		},
	}
}

func newMarkDoneCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "markdone [id]",
		Short: "Mark a task as done",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.ParseInt(args[0], 10, 64)
			_, err := queries.MarkTaskDone(ctx, id)
			cobra.CheckErr(err)
			fmt.Printf("✔️ Task %d marked done\n", id)
		},
	}
}

// displayTasksTable prints tasks in a formatted table
func displayTasksTable(tasks []sqlc.Task) {
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	lessThanDay := []tablewriter.Colors{
		{tablewriter.FgGreenColor},
		{tablewriter.FgGreenColor},
		{tablewriter.FgGreenColor},
		{tablewriter.FgGreenColor},
	}

	lessThanWeek := []tablewriter.Colors{
		{tablewriter.FgYellowColor},
		{tablewriter.FgYellowColor},
		{tablewriter.FgYellowColor},
		{tablewriter.FgYellowColor},
	}

	moreThanWeek := []tablewriter.Colors{
		{tablewriter.FgRedColor},
		{tablewriter.FgRedColor},
		{tablewriter.FgRedColor},
		{tablewriter.FgRedColor},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Title", "Description", "Age"})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlueColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlueColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlueColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlueColor},
	)

	for _, task := range tasks {

		row := []string{
			strconv.FormatInt(task.ID, 10),
			task.Title,
			task.Description.String,
			time.Since(task.CreationDate).Truncate(time.Second).String(),
		}

		if task.CreationDate.After(time.Now().Add(-24 * time.Hour)) {
			table.Rich(row, lessThanDay)
		} else if task.CreationDate.After(time.Now().Add(-7 * 24 * time.Hour)) {
			table.Rich(row, lessThanWeek)
		} else {
			table.Rich(row, moreThanWeek)
		}

	}
	table.Render()
}

// formatTime formats a timestamp for human readability
func formatTimeDuration(t time.Duration) string {

	return t.String()

}

// formatBool returns a check mark or cross emoji for a boolean
func formatBool(b bool) string {
	if b {
		return "✅"
	}
	return "❌"
}

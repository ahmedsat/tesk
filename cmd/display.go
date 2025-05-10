package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	sqlc "github.com/ahmedsat/tesk/sql"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func displayTasksTable(tasks []sqlc.Task) {
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	t := table.NewWriter()
	t.SetTitle("tesk: the task manager for your terminal")

	style := t.Style()
	style.Box = table.StyleBoxRounded

	style.Color.Header = text.Colors{text.BgHiRed, text.FgBlack}
	style.Color.Row = text.Colors{text.BgBlue, text.FgWhite}
	style.Color.RowAlternate = text.Colors{text.BgCyan, text.FgBlack}

	// make border colors transparent
	style.Color.Border = text.Colors{}
	style.Color.Footer = text.Colors{}
	style.Color.Separator = text.Colors{}
	style.Color.IndexColumn = text.Colors{}

	style.Options.SeparateColumns = false

	style.Title = table.TitleOptions{
		Align:  text.AlignCenter,
		Colors: text.Colors{text.FgHiWhite},
		Format: text.FormatTitle,
	}

	t.AppendHeader(table.Row{"ID", "Title", "Description", "Age"})

	for _, task := range tasks {

		row := table.Row{
			strconv.FormatInt(task.ID, 10),
			task.Title,
			task.Description.String,
			// time.Since(task.CreationDate).Truncate(time.Second).String(),
			FormatDuration(time.Since(task.CreationDate)),
		}

		t.AppendRow(row)

	}

	t.SetOutputMirror(os.Stdout)
	t.Render()
}

func FormatDuration(d time.Duration) string {
	d = d.Round(time.Second)

	if d < time.Second {
		return "<1s"
	}

	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	// Format based on largest unit present
	if days > 0 {
		if hours > 0 {
			return fmt.Sprintf("%dd %dh", days, hours)
		}
		return fmt.Sprintf("%dd", days)
	}

	if hours > 0 {
		if minutes > 0 {
			return fmt.Sprintf("%dh %dm", hours, minutes)
		}
		return fmt.Sprintf("%dh", hours)
	}

	if minutes > 0 {
		if seconds > 0 {
			return fmt.Sprintf("%dm %ds", minutes, seconds)
		}
		return fmt.Sprintf("%dm", minutes)
	}

	return fmt.Sprintf("%ds", seconds)
}

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
			time.Since(task.CreationDate).Truncate(time.Second).String(),
		}

		t.AppendRow(row)

	}

	t.SetOutputMirror(os.Stdout)
	t.Render()
}

package ui

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

const (
	purple    = lipgloss.Color("#cffafe")
	gray      = lipgloss.Color("#fb7185")
	lightGray = lipgloss.Color("#22d3ee")
)

func PrintTaskList(rows [][]string) *table.Table {
	re := lipgloss.NewRenderer(os.Stdout)

	var (
		HeaderStyle      = re.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
		CellStyle        = re.NewStyle().Padding(0, 1)
		EvenRowStyle     = CellStyle.Foreground(gray)
		OddRowStyle      = CellStyle.Foreground(lightGray)
		DurationColStyle = CellStyle.Foreground(purple)
		BorderStyle      = lipgloss.NewStyle().Foreground(purple)
	)

	RowHeaders := []string{"ID", "Name", "Status", "Category", "Duration"}
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(BorderStyle).
		StyleFunc(func(row, col int) lipgloss.Style {
			var style lipgloss.Style

			switch {
			case row == 0:
				return HeaderStyle
			case col == 4:
				style = DurationColStyle
			case row%2 == 0:
				style = EvenRowStyle
			default:
				style = OddRowStyle
			}

			return style
		}).
		Headers(RowHeaders...).
		Rows(rows...)

	return t
}

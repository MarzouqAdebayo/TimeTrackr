package ui

import "github.com/charmbracelet/lipgloss"

func PrintTask(str string) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("228")).
		BorderBackground(lipgloss.Color("63")).
		Foreground(lipgloss.Color("228"))
	return style.Render(str)
}


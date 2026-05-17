package ui

import "github.com/charmbracelet/lipgloss"

// Theme defines the color palette.
type Theme struct {
	Primary    lipgloss.Color
	Success    lipgloss.Color
	Warning    lipgloss.Color
	Error      lipgloss.Color
	Muted      lipgloss.Color
	Background lipgloss.Color
	Foreground lipgloss.Color
	Border     lipgloss.Color
	Accent     lipgloss.Color
}

var DefaultTheme = Theme{
	Primary:    lipgloss.Color("#6C8EBF"),
	Success:    lipgloss.Color("#82B366"),
	Warning:    lipgloss.Color("#D6B656"),
	Error:      lipgloss.Color("#B85450"),
	Muted:      lipgloss.Color("#666666"),
	Background: lipgloss.Color("#1E1E2E"),
	Foreground: lipgloss.Color("#CDD6F4"),
	Border:     lipgloss.Color("#45475A"),
	Accent:     lipgloss.Color("#89B4FA"),
}

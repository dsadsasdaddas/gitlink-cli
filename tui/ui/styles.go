package ui

import "github.com/charmbracelet/lipgloss"

// Common styles used across components.
var (
	StyleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(DefaultTheme.Primary)

	StyleMuted = lipgloss.NewStyle().
			Foreground(DefaultTheme.Muted)

	StyleError = lipgloss.NewStyle().
			Foreground(DefaultTheme.Error)

	StyleSuccess = lipgloss.NewStyle().
			Foreground(DefaultTheme.Success)

	StyleWarning = lipgloss.NewStyle().
			Foreground(DefaultTheme.Warning)

	StyleAccent = lipgloss.NewStyle().
			Foreground(DefaultTheme.Accent)

	StyleBorder = lipgloss.RoundedBorder()

	StyleSelected = lipgloss.NewStyle().
			Background(lipgloss.Color("#45475A"))

	StyleCmd = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A6E3A1"))
)

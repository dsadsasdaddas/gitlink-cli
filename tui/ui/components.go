package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// TabDef defines a single tab.
type TabDef struct {
	Key   string
	Label string
}

// TabBarModel renders a horizontal tab bar.
type TabBarModel struct {
	tabs   []TabDef
	active int
}

func NewTabBar(tabs []TabDef) *TabBarModel {
	return &TabBarModel{tabs: tabs}
}

func (t *TabBarModel) View(active int) string {
	var items []string
	for i, tab := range t.tabs {
		label := fmt.Sprintf("[%s] %s", tab.Key, tab.Label)
		if i == active {
			label = StyleSelected.Render(" " + label + " ")
		} else {
			label = StyleMuted.Render(label)
		}
		items = append(items, label)
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, items...)
}

// StatusBarModel renders a bottom status bar.
type StatusBarModel struct {
	project    string
	message    string
	authStatus string
}

func NewStatusBar() *StatusBarModel {
	return &StatusBarModel{}
}

func (s *StatusBarModel) SetProject(p string) {
	s.project = p
}

func (s *StatusBarModel) SetMessage(msg string) {
	s.message = msg
}

func (s *StatusBarModel) SetAuth(status string) {
	s.authStatus = status
}

func (s *StatusBarModel) View() string {
	left := StyleMuted.Render(fmt.Sprintf(" Project: %s ", s.project))
	center := StyleCmd.Render(s.message)
	right := StyleMuted.Render(fmt.Sprintf(" Ctrl+A AI | Ctrl+Space Chat | Ctrl+P Cmd | ? Help "))

	width := 80 // will be overridden by layout
	leftW := lipgloss.Width(left)
	rightW := lipgloss.Width(right)
	centerW := lipgloss.Width(center)

	gapLeft := (width - leftW - centerW - rightW) / 2
	if gapLeft < 0 {
		gapLeft = 0
	}
	gapRight := width - leftW - centerW - rightW - gapLeft
	if gapRight < 0 {
		gapRight = 0
	}

	bg := lipgloss.NewStyle().Background(lipgloss.Color("#313244")).Width(width)
	return bg.Render(left + strings.Repeat(" ", gapLeft) + center + strings.Repeat(" ", gapRight) + right)
}

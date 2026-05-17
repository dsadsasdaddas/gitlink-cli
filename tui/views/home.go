package views

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gitlink-org/gitlink-cli/tui/store"
	"github.com/gitlink-org/gitlink-cli/tui/ui"
)

// HomeModel is the dashboard view.
type HomeModel struct {
	store  *store.Store
	width  int
	height int

	// Data
	project     string
	authStatus  string
	issuesCount string
	prsCount    string
	ciStatus    string
	releaseVer  string

	// Onboarding
	firstRun   bool
	onboardIdx int
}

func NewHomeModel(s *store.Store) *HomeModel {
	return &HomeModel{
		store:    s,
		firstRun: true,
	}
}

func (m *HomeModel) Init() tea.Cmd {
	return tea.Batch(
		m.store.LoadAuthStatus,
		m.store.LoadProjectContext,
	)
}

func (m *HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case store.AuthStatusMsg:
		if msg.LoggedIn {
			m.authStatus = "Logged in"
		} else {
			m.authStatus = "Not logged in (run gitlink-cli auth login)"
		}

	case store.ProjectContextMsg:
		if msg.Error != nil {
			m.project = "No GitLink project detected"
		} else {
			m.project = msg.Owner + "/" + msg.Repo
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.firstRun {
				m.firstRun = false
			}
		}
	}

	return m, nil
}

func (m *HomeModel) View() string {
	if m.firstRun {
		return m.onboardingView()
	}
	return m.dashboardView()
}

func (m *HomeModel) dashboardView() string {
	title := ui.StyleTitle.Render("🏠 Home Dashboard")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		fmt.Sprintf("  Project:    %s", ui.StyleAccent.Render(m.project)),
		fmt.Sprintf("  Auth:       %s", m.authStatus),
		"",
		ui.StyleMuted.Render("  Quick start:"),
		fmt.Sprintf("  %s  Switch views (1-7)", ui.StyleCmd.Render("[1]-[7]")),
		fmt.Sprintf("  %s  AI actions (context-aware)", ui.StyleCmd.Render("Ctrl+A")),
		fmt.Sprintf("  %s  Open AI chat panel", ui.StyleCmd.Render("Ctrl+Space")),
		fmt.Sprintf("  %s  Command palette", ui.StyleCmd.Render("Ctrl+P")),
		fmt.Sprintf("  %s  Help", ui.StyleCmd.Render("?")),
	)

	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(content)
}

func (m *HomeModel) onboardingView() string {
	welcome := lipgloss.NewStyle().
		Bold(true).
		Foreground(ui.DefaultTheme.Primary).
		Render("Welcome to GitLink TUI")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.DefaultTheme.Primary).
		Padding(1, 2).
		Width(60)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		welcome,
		"",
		lipgloss.NewStyle().Bold(true).Render("One key to remember:"),
		"",
		fmt.Sprintf("  %s → AI actions, anywhere, context-aware", ui.StyleCmd.Render("Ctrl+A")),
		"",
		lipgloss.NewStyle().Bold(true).Render("Also useful:"),
		fmt.Sprintf("  %s  AI Chat panel", ui.StyleCmd.Render("Ctrl+Space")),
		fmt.Sprintf("  %s  Command palette", ui.StyleCmd.Render("Ctrl+P")),
		fmt.Sprintf("  %s  Switch views", ui.StyleCmd.Render("[1]-[7]")),
		fmt.Sprintf("  %s  Help", ui.StyleCmd.Render("?")),
		"",
		fmt.Sprintf("  Project: %s", ui.StyleAccent.Render(m.project)),
		fmt.Sprintf("  Auth:    %s", m.authStatus),
		"",
		ui.StyleCmd.Render("  [Enter] Got it!"),
	)

	return box.Render(content)
}

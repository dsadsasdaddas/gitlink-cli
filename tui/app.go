package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gitlink-org/gitlink-cli/tui/navigation"
	"github.com/gitlink-org/gitlink-cli/tui/store"
	"github.com/gitlink-org/gitlink-cli/tui/ui"
	"github.com/gitlink-org/gitlink-cli/tui/views"
)

// View enum
type View int

const (
	ViewHome View = iota
	ViewCode
	ViewIssueBoard
	ViewIssueDetail
	ViewPRList
	ViewPRDetail
	ViewCIDashboard
	ViewCILog
	ViewReleaseList
	ViewReleaseCreate
	ViewPMDashboard
	ViewExecutionDashboard
	ViewSearch
	ViewSettings
)

// AppModel is the root Bubble Tea model.
type AppModel struct {
	// Navigation
	currentView View
	router      *navigation.Router

	// Sub-models for each view
	home           *views.HomeModel
	codeView       tea.Model
	issueBoard     tea.Model
	issueDetail    tea.Model
	prList         tea.Model
	prDetail       tea.Model
	ciDashboard    tea.Model
	ciLog          tea.Model
	releaseList    tea.Model
	releaseCreate  tea.Model
	pmDashboard    tea.Model
	execDashboard  tea.Model
	searchView     tea.Model
	settingsView   tea.Model

	// Layout
	width  int
	height int
	ready  bool

	// Panels
	statusBar *ui.StatusBarModel
	tabBar    *ui.TabBarModel
	aiPanel   tea.Model // nil = not created yet (Phase 3)
	aiOpen    bool

	// Shared state
	store *store.Store

	// Key bindings
	keyMap ui.KeyMap

	// Notifications
	notifications []string
}

// tabs defines the global tab bar.
var tabs = []ui.TabDef{
	{Key: "1", Label: "Home"},
	{Key: "2", Label: "Code"},
	{Key: "3", Label: "Issues"},
	{Key: "4", Label: "PRs"},
	{Key: "5", Label: "CI"},
	{Key: "6", Label: "Release"},
	{Key: "7", Label: "PM"},
}

func NewApp() *AppModel {
	s := store.New()

	statusBar := ui.NewStatusBar()
	tabBar := ui.NewTabBar(tabs)

	return &AppModel{
		currentView: ViewHome,
		router:      navigation.NewRouter(),
		store:       s,
		statusBar:   statusBar,
		tabBar:      tabBar,
		keyMap:      ui.DefaultKeyMap,
		home:        views.NewHomeModel(s),
	}
}

func (m *AppModel) Init() tea.Cmd {
	// Load auth status and current project info
	return tea.Batch(
		m.store.LoadAuthStatus,
		m.store.LoadProjectContext,
	)
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if !m.ready {
			m.ready = true
		}
		return m, nil

	case tea.KeyMsg:
		// Global keybindings take priority over view-specific ones
		cmd := m.handleGlobalKeys(msg)
		if cmd != nil {
			return m, cmd
		}

	case store.AuthStatusMsg:
		// Update statusbar with auth info
		_ = msg
	case store.ProjectContextMsg:
		// Update statusbar with project info
		m.statusBar.SetProject(msg.Owner + "/" + msg.Repo)
	}
	cmds = nil

	// Delegate to active view
	viewModel := m.currentViewModel()
	if viewModel != nil {
		newModel, cmd := viewModel.Update(msg)
		// If the model returned a different model (unlikely but handle)
		_ = newModel
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *AppModel) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Compute layout dimensions
	statusHeight := 1
	tabHeight := 1

	mainHeight := m.height - statusHeight - tabHeight
	if mainHeight < 0 {
		mainHeight = 0
	}

	mainWidth := m.width
	aiWidth := 0
	if m.aiOpen && m.aiPanel != nil {
		aiWidth = min(40, m.width*30/100)
		mainWidth = m.width - aiWidth
	}

	// Render active view
	var mainContent string
	viewModel := m.currentViewModel()
	if viewModel != nil {
		mainContent = viewModel.View()
	}
	mainBox := lipgloss.NewStyle().
		Width(mainWidth).
		Height(mainHeight).
		Render(mainContent)

	// AI panel placeholder
	var aiContent string
	if m.aiOpen {
		aiContent = lipgloss.NewStyle().
			Width(aiWidth).
			Height(mainHeight).
			BorderLeft(true).
			BorderStyle(lipgloss.NormalBorder()).
			Padding(0, 1).
			Render("AI Panel\n───\nCtrl+Space to toggle")
	}

	// Combine main + AI
	topSection := lipgloss.JoinHorizontal(lipgloss.Top, mainBox, aiContent)

	// Tab bar
	var tabBarStr string
	if !m.aiOpen || aiWidth == 0 {
		// Highlight active tab
		tabBarStr = m.tabBar.View(int(m.currentView))
	} else {
		tabBarStr = m.tabBar.View(int(m.currentView))
	}

	// Status bar
	statusStr := m.statusBar.View()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		topSection,
		tabBarStr,
		statusStr,
	)
}

func (m *AppModel) handleGlobalKeys(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+c":
		return tea.Quit

	case "ctrl+space":
		m.aiOpen = !m.aiOpen
		return nil

	case "ctrl+p":
		// Command palette (Phase 5, placeholder)
		m.notifications = append(m.notifications, "Command palette coming in Phase 5")
		return nil

	case "?":
		// Help overlay (placeholder)
		m.notifications = append(m.notifications, "Help: j/k navigate, Enter select, Esc back, Ctrl+A AI actions")
		return nil

	case "ctrl+a":
		// AI context menu (Phase 3, placeholder)
		m.notifications = append(m.notifications, fmt.Sprintf("AI menu for: %s", m.viewName()))
		return nil

	case "1", "2", "3", "4", "5", "6", "7":
		return m.switchToTab(msg.String())

	case "tab":
		// Cycle through views
		next := (int(m.currentView) + 1) % len(tabs)
		m.currentView = View(next)
		return nil
	}

	// Delegate to view
	return nil
}

func (m *AppModel) switchToTab(key string) tea.Cmd {
	switch key {
	case "1":
		m.currentView = ViewHome
	case "2":
		m.currentView = ViewCode
	case "3":
		m.currentView = ViewIssueBoard
	case "4":
		m.currentView = ViewPRList
	case "5":
		m.currentView = ViewCIDashboard
	case "6":
		m.currentView = ViewReleaseList
	case "7":
		m.currentView = ViewPMDashboard
	}
	return nil
}

func (m *AppModel) currentViewModel() tea.Model {
	switch m.currentView {
	case ViewHome:
		return m.home
	case ViewCode:
		return m.codeView
	case ViewIssueBoard:
		return m.issueBoard
	case ViewIssueDetail:
		return m.issueDetail
	case ViewPRList:
		return m.prList
	case ViewPRDetail:
		return m.prDetail
	case ViewCIDashboard:
		return m.ciDashboard
	case ViewCILog:
		return m.ciLog
	case ViewReleaseList:
		return m.releaseList
	case ViewReleaseCreate:
		return m.releaseCreate
	case ViewPMDashboard:
		return m.pmDashboard
	case ViewExecutionDashboard:
		return m.execDashboard
	case ViewSearch:
		return m.searchView
	case ViewSettings:
		return m.settingsView
	}
	return nil
}

func (m *AppModel) viewName() string {
	names := map[View]string{
		ViewHome:               "Home",
		ViewCode:               "Code",
		ViewIssueBoard:         "Issues",
		ViewIssueDetail:        "Issue Detail",
		ViewPRList:             "PRs",
		ViewPRDetail:           "PR Detail",
		ViewCIDashboard:        "CI",
		ViewCILog:              "CI Log",
		ViewReleaseList:        "Releases",
		ViewReleaseCreate:      "New Release",
		ViewPMDashboard:        "PM",
		ViewExecutionDashboard: "Execution",
		ViewSearch:             "Search",
		ViewSettings:           "Settings",
	}
	if n, ok := names[m.currentView]; ok {
		return n
	}
	return "Unknown"
}

func (m *AppModel) navigate(to View) {
	m.router.Push(int(m.currentView))
	m.currentView = to
}

func (m *AppModel) goBack() {
	if prev, ok := m.router.Pop(); ok {
		m.currentView = View(prev)
	}
}

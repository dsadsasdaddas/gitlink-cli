package store

import (
	"fmt"
	"sync"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/gitlink-org/gitlink-cli/internal/auth"
	"github.com/gitlink-org/gitlink-cli/internal/client"
	"github.com/gitlink-org/gitlink-cli/internal/config"
	"github.com/gitlink-org/gitlink-cli/internal/context"
)

// Store is the central state container.
type Store struct {
	mu sync.RWMutex

	// Client for GitLink API calls
	Client *client.Client

	// Config
	Config *config.Config

	// Auth
	Token    string
	LoggedIn bool
	User     *UserInfo

	// Project
	Owner string
	Repo  string

	// Subscribers
	subs map[string][]chan tea.Msg
}

// UserInfo holds basic authenticated user data.
type UserInfo struct {
	Login    string
	ID       int
	Username string
}

func New() *Store {
	cli, err := client.New()
	if err != nil {
		cli = nil
	}

	return &Store{
		Client: cli,
		subs:   make(map[string][]chan tea.Msg),
	}
}

// LoadAuthStatus checks whether the user is authenticated.
func (s *Store) LoadAuthStatus() tea.Msg {
	token, err := auth.LoadToken()
	if err != nil {
		return AuthStatusMsg{LoggedIn: false, Error: err}
	}

	s.mu.Lock()
	s.Token = token
	s.LoggedIn = token != ""
	s.mu.Unlock()

	return AuthStatusMsg{LoggedIn: token != ""}
}

// LoadProjectContext detects the current GitLink project from git remote.
func (s *Store) LoadProjectContext() tea.Msg {
	owner, repo, err := context.ResolveOwnerRepo("", "")
	if err != nil {
		return ProjectContextMsg{Error: err}
	}

	s.mu.Lock()
	s.Owner = owner
	s.Repo = repo
	s.mu.Unlock()

	return ProjectContextMsg{
		Owner: owner,
		Repo:  repo,
	}
}

// Dispatch sends an action to the store and publishes to subscribers.
func (s *Store) Dispatch(action string, data interface{}) {
	s.mu.RLock()
	subs := s.subs[action]
	s.mu.RUnlock()

	for _, ch := range subs {
		select {
		case ch <- StoreEvent{Action: action, Data: data}:
		default:
		}
	}
}

// Subscribe returns a channel for store events.
func (s *Store) Subscribe(action string) chan tea.Msg {
	ch := make(chan tea.Msg, 10)
	s.mu.Lock()
	s.subs[action] = append(s.subs[action], ch)
	s.mu.Unlock()
	return ch
}

// AuthStatusMsg is sent when auth status is loaded.
type AuthStatusMsg struct {
	LoggedIn bool
	Error    error
}

// ProjectContextMsg is sent when project context is loaded.
type ProjectContextMsg struct {
	Owner string
	Repo  string
	Error error
}

func (m ProjectContextMsg) String() string {
	if m.Error != nil {
		return fmt.Sprintf("No project (%v)", m.Error)
	}
	return m.Owner + "/" + m.Repo
}

// StoreEvent is a generic store event.
type StoreEvent struct {
	Action string
	Data   interface{}
}

func (e StoreEvent) String() string {
	return fmt.Sprintf("[%s]", e.Action)
}

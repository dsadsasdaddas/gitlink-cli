package ui

import "github.com/charmbracelet/bubbles/key"

// KeyMap holds all global keybindings.
type KeyMap struct {
	Quit     key.Binding
	Help     key.Binding
	AI       key.Binding
	AIPanel  key.Binding
	Commands key.Binding
	Refresh  key.Binding

	Up    key.Binding
	Down  key.Binding
	Back  key.Binding
	Enter key.Binding

	Create key.Binding
	Delete key.Binding
	Edit   key.Binding
	Filter key.Binding
}

var DefaultKeyMap = KeyMap{
	Quit:     key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("Ctrl+C", "Quit")),
	Help:     key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "Help")),
	AI:       key.NewBinding(key.WithKeys("ctrl+a"), key.WithHelp("Ctrl+A", "AI actions")),
	AIPanel:  key.NewBinding(key.WithKeys("ctrl+space"), key.WithHelp("Ctrl+Space", "AI Panel")),
	Commands: key.NewBinding(key.WithKeys("ctrl+p"), key.WithHelp("Ctrl+P", "Commands")),
	Refresh:  key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "Refresh")),

	Up:    key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("j/k", "Navigate")),
	Down:  key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("", "")),
	Back:  key.NewBinding(key.WithKeys("esc", "h"), key.WithHelp("Esc", "Back")),
	Enter: key.NewBinding(key.WithKeys("enter", "l"), key.WithHelp("Enter", "Select")),

	Create: key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "Create")),
	Delete: key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "Delete")),
	Edit:   key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "Edit")),
	Filter: key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "Filter")),
}

package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/gitlink-org/gitlink-cli/cmd"
	"github.com/gitlink-org/gitlink-cli/tui"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "tui" {
		// TUI mode
		f, err := tea.LogToFile("gitlink-tui.log", "")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to init log: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()

		app := tui.NewApp()
		program := tea.NewProgram(
			app,
			tea.WithAltScreen(),
			tea.WithMouseCellMotion(),
		)
		if _, err := program.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// CLI mode
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

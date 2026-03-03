package main

import (
	"fmt"
	"os"
	"runtime/debug"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oronbz/nag/internal/reminders"
	"github.com/oronbz/nag/internal/ui"
)

var version = "dev"

func getVersion() string {
	if version != "dev" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return version
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "help", "--help", "-h":
			printHelp()
			return
		case "version", "--version", "-v":
			fmt.Println("nag " + getVersion())
			return
		}
	}

	client, err := reminders.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to access Reminders.")
		fmt.Fprintln(os.Stderr, "Make sure you've granted Reminders access in System Settings > Privacy & Security > Reminders.")
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	model := ui.NewModel(client)
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Print(`nag — A terminal UI for Apple Reminders

Usage:
  nag              Launch the TUI
  nag help         Show this help message
  nag version      Show version

Note:
  On first run, macOS will prompt for Reminders access.
  You can manage this in System Settings > Privacy & Security > Reminders.

Keybindings (inside TUI):
  j/k, ↑/↓            Navigate lists
  g / G                Jump to top / bottom
  Ctrl-u / Ctrl-d      Page up / page down
  Enter                Select list
  Tab / Shift-Tab      Switch panel
  Space / x            Toggle reminder complete
  n                    New reminder / list (context-sensitive)
  e                    Edit reminder / list (context-sensitive)
  d                    Delete reminder / list (with confirmation)
  o                    Open in Apple Reminders
  c                    Toggle show completed
  r                    Refresh
  /                    Filter / search
  ?                    Show all keybindings
  Esc                  Close dialog / popup
  q, Ctrl-C            Quit
`)
}

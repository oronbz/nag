package ui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up             key.Binding
	Down           key.Binding
	Enter          key.Binding
	Tab            key.Binding
	ShiftTab       key.Binding
	Top            key.Binding
	Bottom         key.Binding
	PageDown       key.Binding
	PageUp         key.Binding
	ToggleComplete key.Binding
	NewReminder    key.Binding
	Edit           key.Binding
	Delete          key.Binding
	OpenInApp       key.Binding
	ShowCompleted   key.Binding
	Sort            key.Binding
	Refresh         key.Binding
	Filter         key.Binding
	Help           key.Binding
	Quit           key.Binding
	Escape         key.Binding
}

var Keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("⏎", "select"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("Tab", "next panel"),
	),
	ShiftTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("S-Tab", "prev panel"),
	),
	Top: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "top"),
	),
	Bottom: key.NewBinding(
		key.WithKeys("G"),
		key.WithHelp("G", "bottom"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("C-d", "page down"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("ctrl+u"),
		key.WithHelp("C-u", "page up"),
	),
	ToggleComplete: key.NewBinding(
		key.WithKeys(" ", "x"),
		key.WithHelp("Space/x", "toggle complete"),
	),
	NewReminder: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new reminder / list"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit reminder / list"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	OpenInApp: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "open in Reminders"),
	),
	ShowCompleted: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "toggle completed"),
	),
	Sort: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "cycle sort"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("Esc", "close"),
	),
}

package statusbar

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/oronbz/nag/internal/ui/styles"
)

type Panel int

const (
	PanelLists Panel = iota
	PanelReminders
)

type Model struct {
	panel   Panel
	errMsg  string
	infoMsg string
	loading string
}

func New() Model {
	return Model{}
}

func (m *Model) SetPanel(panel Panel) {
	m.panel = panel
}

func (m *Model) SetError(msg string) {
	m.errMsg = msg
}

func (m *Model) ClearError() {
	m.errMsg = ""
}

func (m *Model) SetLoading(msg string) {
	m.loading = msg
}

func (m *Model) ClearLoading() {
	m.loading = ""
}

func (m *Model) SetInfo(msg string) {
	m.infoMsg = msg
}

func (m *Model) ClearInfo() {
	m.infoMsg = ""
}

type hint struct {
	key  string
	desc string
}

var (
	brandStyle = lipgloss.NewStyle().
			Background(styles.Teal).
			Foreground(lipgloss.Color("#000000")).
			Bold(true).
			Padding(0, 1)
	keyStyle = lipgloss.NewStyle().Foreground(styles.Teal).Bold(true)
	dimStyle = lipgloss.NewStyle().Foreground(styles.DimGray)
	errStyle = lipgloss.NewStyle().Foreground(styles.Red)
	infoStyle = lipgloss.NewStyle().Foreground(styles.Green)

	listsHints = []hint{
		{"↑/k", "up"}, {"↓/j", "down"},
		{"Tab", "panel"}, {"/", "filter"}, {"?", "help"}, {"q", "quit"},
	}
	remindersHints = []hint{
		{"↑/k", "up"}, {"↓/j", "down"},
		{"Tab", "panel"}, {"Space", "toggle"}, {"n", "new"},
		{"d", "delete"}, {"c", "completed"}, {"/", "filter"}, {"?", "help"}, {"q", "quit"},
	}
)

func (m Model) View() string {
	brand := brandStyle.Render("NAG")

	var hints []hint
	switch m.panel {
	case PanelLists:
		hints = listsHints
	case PanelReminders:
		hints = remindersHints
	}

	var parts []string
	for _, h := range hints {
		parts = append(parts, keyStyle.Render(h.key)+" "+h.desc)
	}
	hintsStr := strings.Join(parts, "  ")

	var rightText string
	if m.errMsg != "" {
		rightText = "  " + errStyle.Render(m.errMsg)
	} else if m.infoMsg != "" {
		rightText = "  " + infoStyle.Render(m.infoMsg)
	} else if m.loading != "" {
		rightText = "  " + dimStyle.Render(m.loading)
	}

	return brand + " " + dimStyle.Render("│") + " " + hintsStr + rightText
}

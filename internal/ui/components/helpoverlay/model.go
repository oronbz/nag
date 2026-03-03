package helpoverlay

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oronbz/nag/internal/ui/styles"
)

type Model struct {
	viewport viewport.Model
	visible  bool
	width    int
	height   int
}

func New() Model {
	return Model{}
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	vw := width * 60 / 100
	vh := height * 70 / 100
	if vw < 50 {
		vw = 50
	}
	if vh < 20 {
		vh = 20
	}
	m.viewport = viewport.New(vw-4, vh-4)
	m.viewport.SetContent(helpContent())
}

func (m *Model) Toggle() {
	m.visible = !m.visible
	if m.visible {
		m.viewport.GotoTop()
	}
}

func (m Model) Visible() bool {
	return m.visible
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.visible {
		return m, nil
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if !m.visible {
		return ""
	}

	content := styles.HelpOverlayStyle.
		Width(m.viewport.Width + 4).
		Render(m.viewport.View())

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

func helpContent() string {
	k := styles.HelpKeyStyle
	d := styles.HelpDescStyle
	h := styles.HelpHeaderStyle

	var sb strings.Builder

	sb.WriteString(h.Render("Navigation"))
	sb.WriteString("\n")
	sb.WriteString(k.Render("j / ↓") + d.Render("Move down") + "\n")
	sb.WriteString(k.Render("k / ↑") + d.Render("Move up") + "\n")
	sb.WriteString(k.Render("Enter") + d.Render("Select list") + "\n")
	sb.WriteString(k.Render("Tab / S-Tab") + d.Render("Switch panel") + "\n")
	sb.WriteString(k.Render("g") + d.Render("Jump to top") + "\n")
	sb.WriteString(k.Render("G") + d.Render("Jump to bottom") + "\n")
	sb.WriteString(k.Render("Ctrl-d") + d.Render("Page down") + "\n")
	sb.WriteString(k.Render("Ctrl-u") + d.Render("Page up") + "\n")

	sb.WriteString("\n")
	sb.WriteString(h.Render("Actions"))
	sb.WriteString("\n")
	sb.WriteString(k.Render("Space / x") + d.Render("Toggle reminder complete") + "\n")
	sb.WriteString(k.Render("n") + d.Render("New reminder / list") + "\n")
	sb.WriteString(k.Render("e") + d.Render("Edit reminder / list") + "\n")
	sb.WriteString(k.Render("d") + d.Render("Delete reminder / list") + "\n")
	sb.WriteString(k.Render("o") + d.Render("Open in Reminders app") + "\n")
	sb.WriteString(k.Render("s") + d.Render("Cycle sort order") + "\n")
	sb.WriteString(k.Render("c") + d.Render("Toggle show completed") + "\n")
	sb.WriteString(k.Render("r") + d.Render("Refresh current view") + "\n")
	sb.WriteString(k.Render("/") + d.Render("Filter / search") + "\n")

	sb.WriteString("\n")
	sb.WriteString(h.Render("General"))
	sb.WriteString("\n")
	sb.WriteString(k.Render("?") + d.Render("Toggle this help") + "\n")
	sb.WriteString(k.Render("Esc") + d.Render("Close dialog / popup") + "\n")
	sb.WriteString(k.Render("q / Ctrl-C") + d.Render("Quit") + "\n")

	return sb.String()
}

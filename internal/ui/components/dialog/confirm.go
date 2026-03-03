package dialog

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oronbz/nag/internal/ui/styles"
)

type ConfirmAction int

const (
	ConfirmDelete ConfirmAction = iota
	ConfirmDeleteList
)

type ConfirmYesMsg struct{ Action ConfirmAction }
type ConfirmNoMsg struct{}

type ConfirmModel struct {
	message string
	action  ConfirmAction
	visible bool
	width   int
	height  int
}

func NewConfirm() ConfirmModel {
	return ConfirmModel{}
}

func (m *ConfirmModel) Show(message string, action ConfirmAction) {
	m.message = message
	m.action = action
	m.visible = true
}

func (m *ConfirmModel) Hide() {
	m.visible = false
}

func (m ConfirmModel) Visible() bool {
	return m.visible
}

func (m *ConfirmModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m ConfirmModel) Update(msg tea.Msg) (ConfirmModel, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y", "enter":
			action := m.action
			m.Hide()
			return m, func() tea.Msg { return ConfirmYesMsg{Action: action} }
		case "n", "N", "esc":
			m.Hide()
			return m, func() tea.Msg { return ConfirmNoMsg{} }
		}
	}

	return m, nil
}

func (m ConfirmModel) View() string {
	if !m.visible {
		return ""
	}

	title := styles.DialogTitleStyle.Render("Confirm")
	content := title + "\n\n" +
		m.message + "\n\n" +
		lipgloss.NewStyle().Foreground(styles.DimGray).Render("y: yes  n/Esc: no")

	dialog := styles.DialogStyle.Render(content)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		dialog,
	)
}

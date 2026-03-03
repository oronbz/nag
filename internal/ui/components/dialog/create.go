package dialog

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oronbz/nag/internal/reminders"
	"github.com/oronbz/nag/internal/ui/styles"
)

type CreateSubmitMsg struct {
	Input reminders.CreateReminderInput
}

type CreateModel struct {
	titleInput    textinput.Model
	dueDateInput  textinput.Model
	priorityInput textinput.Model
	focusIndex    int
	visible       bool
	listName      string
	width         int
	height        int
}

func NewCreate() CreateModel {
	ti := textinput.New()
	ti.Placeholder = "Buy groceries"
	ti.CharLimit = 256
	ti.Width = 44
	ti.Prompt = "Title:    "
	ti.PromptStyle = lipgloss.NewStyle().Foreground(styles.Teal)

	di := textinput.New()
	di.Placeholder = "today, tomorrow, 2025-03-15, 2025-03-15 14:30"
	di.CharLimit = 32
	di.Width = 44
	di.Prompt = "Due date: "
	di.PromptStyle = lipgloss.NewStyle().Foreground(styles.Teal)

	pi := textinput.New()
	pi.Placeholder = "none, low, medium, high"
	pi.CharLimit = 16
	pi.Width = 44
	pi.Prompt = "Priority: "
	pi.PromptStyle = lipgloss.NewStyle().Foreground(styles.Teal)

	return CreateModel{
		titleInput:    ti,
		dueDateInput:  di,
		priorityInput: pi,
	}
}

func (m *CreateModel) Show(listName string) {
	m.visible = true
	m.listName = listName
	m.focusIndex = 0
	m.titleInput.SetValue("")
	m.dueDateInput.SetValue("")
	m.priorityInput.SetValue("")
	m.titleInput.Focus()
	m.dueDateInput.Blur()
	m.priorityInput.Blur()
}

func (m *CreateModel) Hide() {
	m.visible = false
	m.titleInput.Blur()
	m.dueDateInput.Blur()
	m.priorityInput.Blur()
}

func (m CreateModel) Visible() bool {
	return m.visible
}

func (m *CreateModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m CreateModel) Update(msg tea.Msg) (CreateModel, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.Hide()
			return m, nil
		case "tab", "shift+tab":
			m.focusIndex = (m.focusIndex + 1) % 3
			m.titleInput.Blur()
			m.dueDateInput.Blur()
			m.priorityInput.Blur()
			switch m.focusIndex {
			case 0:
				m.titleInput.Focus()
			case 1:
				m.dueDateInput.Focus()
			case 2:
				m.priorityInput.Focus()
			}
			return m, nil
		case "enter":
			title := strings.TrimSpace(m.titleInput.Value())
			if title == "" {
				return m, nil
			}
			input := reminders.CreateReminderInput{
				Title:    title,
				ListName: m.listName,
				DueDate:  parseDueDate(m.dueDateInput.Value()),
				Priority: parsePriority(m.priorityInput.Value()),
			}
			m.Hide()
			return m, func() tea.Msg {
				return CreateSubmitMsg{Input: input}
			}
		}
	}

	var cmd tea.Cmd
	switch m.focusIndex {
	case 0:
		m.titleInput, cmd = m.titleInput.Update(msg)
	case 1:
		m.dueDateInput, cmd = m.dueDateInput.Update(msg)
	case 2:
		m.priorityInput, cmd = m.priorityInput.Update(msg)
	}
	return m, cmd
}

func (m CreateModel) View() string {
	if !m.visible {
		return ""
	}

	title := styles.DialogTitleStyle.Render("New Reminder")
	content := title + "\n\n" +
		m.titleInput.View() + "\n\n" +
		m.dueDateInput.View() + "\n\n" +
		m.priorityInput.View() + "\n\n" +
		lipgloss.NewStyle().Foreground(styles.DimGray).Render("Enter: create  Tab: next field  Esc: cancel")

	dialog := styles.DialogStyle.Render(content)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		dialog,
	)
}

func parseDueDate(s string) *time.Time {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return nil
	}

	now := time.Now()
	loc := now.Location()

	switch s {
	case "today":
		t := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, loc)
		return &t
	case "tomorrow":
		t := time.Date(now.Year(), now.Month(), now.Day()+1, 9, 0, 0, 0, loc)
		return &t
	}

	// Date + time: "2025-03-15 14:30"
	if parsed, err := time.ParseInLocation("2006-01-02 15:04", s, loc); err == nil {
		return &parsed
	}

	// Date only: "2025-03-15" (defaults to 9:00 AM)
	if parsed, err := time.ParseInLocation("2006-01-02", s, loc); err == nil {
		t := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 9, 0, 0, 0, loc)
		return &t
	}

	return nil
}

func parsePriority(s string) int {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "high", "h":
		return reminders.PriorityHigh
	case "medium", "med", "m":
		return reminders.PriorityMedium
	case "low", "l":
		return reminders.PriorityLow
	default:
		return reminders.PriorityNone
	}
}

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

type EditSubmitMsg struct {
	ID    string
	Input reminders.UpdateReminderInput
}

const fieldCount = 4

type CreateModel struct {
	titleInput    textinput.Model
	notesInput    textinput.Model
	dueDateInput  textinput.Model
	priorityInput textinput.Model
	focusIndex    int
	visible       bool
	listName      string
	editingID     string
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

	ni := textinput.New()
	ni.Placeholder = "Optional notes"
	ni.CharLimit = 1024
	ni.Width = 44
	ni.Prompt = "Notes:    "
	ni.PromptStyle = lipgloss.NewStyle().Foreground(styles.Teal)

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
		notesInput:    ni,
		dueDateInput:  di,
		priorityInput: pi,
	}
}

func (m *CreateModel) Show(listName string) {
	m.visible = true
	m.listName = listName
	m.editingID = ""
	m.focusIndex = 0
	m.titleInput.SetValue("")
	m.notesInput.SetValue("")
	m.dueDateInput.SetValue("")
	m.priorityInput.SetValue("")
	m.focusField(0)
}

func (m *CreateModel) ShowEdit(r reminders.Reminder) {
	m.visible = true
	m.editingID = r.ID
	m.listName = ""
	m.focusIndex = 0

	m.titleInput.SetValue(r.Title)
	m.notesInput.SetValue(r.Notes)

	if r.DueDate != nil {
		m.dueDateInput.SetValue(r.DueDate.Local().Format("2006-01-02 15:04"))
	} else {
		m.dueDateInput.SetValue("")
	}

	switch r.Priority {
	case reminders.PriorityHigh:
		m.priorityInput.SetValue("high")
	case reminders.PriorityMedium:
		m.priorityInput.SetValue("medium")
	case reminders.PriorityLow:
		m.priorityInput.SetValue("low")
	default:
		m.priorityInput.SetValue("")
	}

	m.focusField(0)
}

func (m *CreateModel) Hide() {
	m.visible = false
	m.titleInput.Blur()
	m.notesInput.Blur()
	m.dueDateInput.Blur()
	m.priorityInput.Blur()
}

func (m CreateModel) Visible() bool {
	return m.visible
}

func (m CreateModel) isEditing() bool {
	return m.editingID != ""
}

func (m *CreateModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *CreateModel) focusField(index int) {
	m.titleInput.Blur()
	m.notesInput.Blur()
	m.dueDateInput.Blur()
	m.priorityInput.Blur()
	switch index {
	case 0:
		m.titleInput.Focus()
	case 1:
		m.notesInput.Focus()
	case 2:
		m.dueDateInput.Focus()
	case 3:
		m.priorityInput.Focus()
	}
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
			m.focusIndex = (m.focusIndex + 1) % fieldCount
			m.focusField(m.focusIndex)
			return m, nil
		case "enter":
			title := strings.TrimSpace(m.titleInput.Value())
			if title == "" {
				return m, nil
			}
			m.Hide()
			if m.isEditing() {
				return m, m.buildEditCmd(title)
			}
			input := reminders.CreateReminderInput{
				Title:    title,
				ListName: m.listName,
				Notes:    strings.TrimSpace(m.notesInput.Value()),
				DueDate:  parseDueDate(m.dueDateInput.Value()),
				Priority: parsePriority(m.priorityInput.Value()),
			}
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
		m.notesInput, cmd = m.notesInput.Update(msg)
	case 2:
		m.dueDateInput, cmd = m.dueDateInput.Update(msg)
	case 3:
		m.priorityInput, cmd = m.priorityInput.Update(msg)
	}
	return m, cmd
}

func (m CreateModel) buildEditCmd(title string) tea.Cmd {
	id := m.editingID
	input := reminders.UpdateReminderInput{
		Title: &title,
	}

	notes := strings.TrimSpace(m.notesInput.Value())
	input.Notes = &notes

	dueDateStr := strings.TrimSpace(m.dueDateInput.Value())
	if dueDateStr == "" {
		input.ClearDueDate = true
	} else {
		input.DueDate = parseDueDate(dueDateStr)
	}

	priority := parsePriority(m.priorityInput.Value())
	input.Priority = &priority

	return func() tea.Msg {
		return EditSubmitMsg{ID: id, Input: input}
	}
}

func (m CreateModel) View() string {
	if !m.visible {
		return ""
	}

	dialogTitle := "New Reminder"
	footer := "Enter: create  Tab: next field  Esc: cancel"
	if m.isEditing() {
		dialogTitle = "Edit Reminder"
		footer = "Enter: save  Tab: next field  Esc: cancel"
	}

	title := styles.DialogTitleStyle.Render(dialogTitle)
	content := title + "\n\n" +
		m.titleInput.View() + "\n\n" +
		m.notesInput.View() + "\n\n" +
		m.dueDateInput.View() + "\n\n" +
		m.priorityInput.View() + "\n\n" +
		lipgloss.NewStyle().Foreground(styles.DimGray).Render(footer)

	dlg := styles.DialogStyle.Render(content)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		dlg,
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

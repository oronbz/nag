package reminderpanel

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oronbz/nag/internal/reminders"
	"github.com/oronbz/nag/internal/ui/styles"
)

type Item struct {
	Reminder reminders.Reminder
}

func (i Item) Title() string       { return i.Reminder.Title }
func (i Item) Description() string { return "" }
func (i Item) FilterValue() string { return i.Reminder.Title + " " + i.Reminder.Notes }

type Delegate struct {
	selectedStyle lipgloss.Style
	normalStyle   lipgloss.Style
}

func NewDelegate() Delegate {
	return Delegate{
		selectedStyle: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(styles.Teal).
			PaddingLeft(1),
		normalStyle: lipgloss.NewStyle().
			PaddingLeft(2),
	}
}

func (d Delegate) Height() int                             { return 2 }
func (d Delegate) Spacing() int                            { return 0 }
func (d Delegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d Delegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(Item)
	if !ok {
		return
	}

	r := item.Reminder

	// Line 1: checkbox + title + priority
	checkbox := styles.CheckboxIcon(r.Completed)
	var title string
	if r.Completed {
		title = styles.ReminderCompletedStyle.Render(r.Title)
	} else {
		title = styles.ReminderTitleStyle.Render(r.Title)
	}
	priority := styles.PriorityIcon(r.Priority)
	line1 := fmt.Sprintf("%s %s%s", checkbox, title, priority)

	// Line 2: due date + notes preview
	var parts []string
	if r.DueDate != nil {
		dueStr := formatDueDate(*r.DueDate)
		if isOverdue(*r.DueDate) && !r.Completed {
			parts = append(parts, styles.ReminderOverdueStyle.Render(dueStr))
		} else {
			parts = append(parts, styles.ReminderDimStyle.Render(dueStr))
		}
	}
	if r.Notes != "" {
		note := firstLine(r.Notes)
		if len(note) > 40 {
			note = note[:37] + "..."
		}
		parts = append(parts, styles.ReminderDimStyle.Render(note))
	}
	line2 := "  " + strings.Join(parts, "  ")

	style := d.normalStyle
	if index == m.Index() {
		style = d.selectedStyle
	}
	style = style.MaxWidth(m.Width())

	fmt.Fprint(w, style.Render(line1+"\n"+line2))
}

func formatDueDate(t time.Time) string {
	t = t.Local()
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	due := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, now.Location())

	days := int(due.Sub(today).Hours() / 24)
	var dateStr string
	switch {
	case days < -1:
		dateStr = fmt.Sprintf("%d days ago", -days)
	case days == -1:
		dateStr = "yesterday"
	case days == 0:
		dateStr = "today"
	case days == 1:
		dateStr = "tomorrow"
	case days < 7:
		dateStr = t.Format("Monday")
	default:
		dateStr = t.Format("Jan 2")
	}

	// Append time if set (skip midnight which means no time component)
	if t.Hour() != 0 || t.Minute() != 0 {
		dateStr += " " + t.Format("15:04")
	}

	return dateStr
}

func isOverdue(t time.Time) bool {
	t = t.Local()
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	due := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, now.Location())
	return due.Before(today)
}

func firstLine(s string) string {
	if i := strings.IndexAny(s, "\r\n"); i >= 0 {
		return s[:i]
	}
	return s
}

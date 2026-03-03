package listpanel

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oronbz/nag/internal/reminders"
	"github.com/oronbz/nag/internal/ui/styles"
)

type Item struct {
	List reminders.ReminderList
}

func (i Item) Title() string       { return i.List.Title }
func (i Item) Description() string { return fmt.Sprintf("%d reminders", i.List.Count) }
func (i Item) FilterValue() string {
	if i.List.Kind == reminders.ListSeparator {
		return ""
	}
	return i.List.Title
}

// Delegate renders list items with smart list styling and separators.
type Delegate struct{}

func (d Delegate) Height() int  { return 2 }
func (d Delegate) Spacing() int { return 0 }
func (d Delegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d Delegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(Item)
	if !ok {
		return
	}

	rl := item.List

	if rl.Kind == reminders.ListSeparator {
		sep := styles.ReminderDimStyle.Render(strings.Repeat("─", m.Width()))
		fmt.Fprint(w, "\n"+sep)
		return
	}

	selected := index == m.Index()

	var titleStyle, descStyle lipgloss.Style
	if selected {
		titleStyle = lipgloss.NewStyle().Foreground(styles.Teal).Bold(true)
		descStyle = lipgloss.NewStyle().Foreground(styles.DimGray)
	} else {
		titleStyle = lipgloss.NewStyle().Foreground(styles.White)
		descStyle = lipgloss.NewStyle().Foreground(styles.DimGray)
	}

	// Smart lists get an icon
	title := rl.Title
	if rl.Kind == reminders.ListSmart {
		switch rl.ID {
		case reminders.SmartListToday:
			title = "◉ " + title
		case reminders.SmartListScheduled:
			title = "▦ " + title
		}
	}

	line1 := titleStyle.Render(title)
	line2 := descStyle.Render(fmt.Sprintf("%d reminders", rl.Count))

	var style lipgloss.Style
	if selected {
		style = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(styles.Teal).
			PaddingLeft(1)
	} else {
		style = lipgloss.NewStyle().PaddingLeft(2)
	}

	style = style.MaxWidth(m.Width())
	fmt.Fprint(w, style.Render(line1+"\n"+line2))
}

type Model struct {
	list    list.Model
	focused bool
}

func New(width, height int) Model {
	l := list.New([]list.Item{}, Delegate{}, width, height)
	l.Title = "Lists"
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.Styles.Title = styles.TitleStyle
	l.DisableQuitKeybindings()

	return Model{list: l}
}

func (m *Model) SetSize(width, height int) {
	m.list.SetSize(width, height)
}

func (m *Model) SetFocused(focused bool) {
	m.focused = focused
}

func (m *Model) SetLists(lists []reminders.ReminderList) {
	prevIndex := m.list.Index()
	items := make([]list.Item, len(lists))
	for i, l := range lists {
		items[i] = Item{List: l}
	}
	m.list.SetItems(items)
	if prevIndex > 0 && prevIndex < len(items) {
		m.list.Select(prevIndex)
	}
}

func (m Model) SelectedList() (reminders.ReminderList, bool) {
	item, ok := m.list.SelectedItem().(Item)
	if !ok {
		return reminders.ReminderList{}, false
	}
	// Don't select separators
	if item.List.Kind == reminders.ListSeparator {
		return reminders.ReminderList{}, false
	}
	return item.List, true
}

func (m Model) Filtering() bool {
	return m.list.FilterState() == list.Filtering
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focused {
		if _, ok := msg.(tea.MouseMsg); !ok {
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.list.View()
}

package commands

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/oronbz/nag/internal/reminders"
	"github.com/oronbz/nag/internal/ui/messages"
)

func FetchReminders(client *reminders.Client, listName string, showCompleted bool) tea.Cmd {
	return func() tea.Msg {
		items, err := client.Reminders(listName, showCompleted)
		return messages.RemindersLoadedMsg{Reminders: items, Err: err}
	}
}

func FetchTodayReminders(client *reminders.Client, showCompleted bool) tea.Cmd {
	return func() tea.Msg {
		items, err := client.TodayReminders(showCompleted)
		return messages.RemindersLoadedMsg{Reminders: items, Err: err}
	}
}

func FetchScheduledReminders(client *reminders.Client, showCompleted bool) tea.Cmd {
	return func() tea.Msg {
		items, err := client.ScheduledReminders(showCompleted)
		return messages.RemindersLoadedMsg{Reminders: items, Err: err}
	}
}

func CreateReminder(client *reminders.Client, input reminders.CreateReminderInput) tea.Cmd {
	return func() tea.Msg {
		r, err := client.CreateReminder(input)
		return messages.ReminderCreatedMsg{Reminder: r, Err: err}
	}
}

func UpdateReminder(client *reminders.Client, id string, input reminders.UpdateReminderInput) tea.Cmd {
	return func() tea.Msg {
		r, err := client.UpdateReminder(id, input)
		return messages.ReminderUpdatedMsg{Reminder: r, Err: err}
	}
}

func ToggleComplete(client *reminders.Client, id string, currentlyCompleted bool) tea.Cmd {
	return func() tea.Msg {
		var r *reminders.Reminder
		var err error
		if currentlyCompleted {
			r, err = client.UncompleteReminder(id)
		} else {
			r, err = client.CompleteReminder(id)
		}
		return messages.ReminderCompletedMsg{Reminder: r, Err: err}
	}
}

func DeleteReminder(client *reminders.Client, id string) tea.Cmd {
	return func() tea.Msg {
		err := client.DeleteReminder(id)
		return messages.ReminderDeletedMsg{ID: id, Err: err}
	}
}

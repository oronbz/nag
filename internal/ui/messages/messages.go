package messages

import "github.com/oronbz/nag/internal/reminders"

type ListsLoadedMsg struct {
	Lists []reminders.ReminderList
	Err   error
}

type RemindersLoadedMsg struct {
	Reminders []reminders.Reminder
	Err       error
}

type ReminderCreatedMsg struct {
	Reminder *reminders.Reminder
	Err      error
}

type ReminderCompletedMsg struct {
	Reminder *reminders.Reminder
	Err      error
}

type ReminderDeletedMsg struct {
	ID  string
	Err error
}

type ClearInfoMsg struct{}

type TickMsg struct{}

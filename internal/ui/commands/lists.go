package commands

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/oronbz/nag/internal/reminders"
	"github.com/oronbz/nag/internal/ui/messages"
)

func FetchLists(client *reminders.Client) tea.Cmd {
	return func() tea.Msg {
		lists, err := client.ListsWithSmart()
		return messages.ListsLoadedMsg{Lists: lists, Err: err}
	}
}

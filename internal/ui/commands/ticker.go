package commands

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oronbz/nag/internal/ui/messages"
)

const RefreshInterval = 10 * time.Second

func AutoRefreshTick() tea.Cmd {
	return tea.Tick(RefreshInterval, func(time.Time) tea.Msg {
		return messages.TickMsg{}
	})
}

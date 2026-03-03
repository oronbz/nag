package reminders

import ekreminders "github.com/BRO3886/go-eventkit/reminders"

const (
	SmartListToday     = "__smart_today__"
	SmartListScheduled = "__smart_scheduled__"
)

func (c *Client) Lists() ([]ReminderList, error) {
	ekLists, err := c.ek.Lists()
	if err != nil {
		return nil, err
	}

	// Single query for all incomplete reminders, then count by list ID
	incompleteCounts := make(map[string]int)
	if all, err := c.ek.Reminders(ekreminders.WithCompleted(false)); err == nil {
		for _, r := range all {
			incompleteCounts[r.ListID]++
		}
	}

	lists := make([]ReminderList, len(ekLists))
	for i, l := range ekLists {
		lists[i] = ReminderList{
			ID:    l.ID,
			Title: l.Title,
			Count: incompleteCounts[l.ID],
		}
	}
	return lists, nil
}

// ListsWithSmart returns smart lists (Today, Scheduled) + separator + normal lists.
func (c *Client) ListsWithSmart() ([]ReminderList, error) {
	normal, err := c.Lists()
	if err != nil {
		return nil, err
	}

	todayItems, _ := c.TodayReminders(false)
	scheduledItems, _ := c.ScheduledReminders(false)

	smart := []ReminderList{
		{ID: SmartListToday, Title: "Today", Count: len(todayItems), Kind: ListSmart},
		{ID: SmartListScheduled, Title: "Scheduled", Count: len(scheduledItems), Kind: ListSmart},
		{Kind: ListSeparator},
	}

	return append(smart, normal...), nil
}

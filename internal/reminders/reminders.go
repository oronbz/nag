package reminders

import (
	"sort"
	"time"

	ekreminders "github.com/BRO3886/go-eventkit/reminders"
)

func (c *Client) Reminders(listName string, showCompleted bool) ([]Reminder, error) {
	opts := []ekreminders.ListOption{ekreminders.WithList(listName)}
	if !showCompleted {
		opts = append(opts, ekreminders.WithCompleted(false))
	}
	items, err := c.ek.Reminders(opts...)
	if err != nil {
		return nil, err
	}
	result := convertReminders(items)
	if showCompleted {
		sortReminders(result)
	}
	return result, nil
}

// TodayReminders returns reminders due today or overdue, across all lists.
func (c *Client) TodayReminders(showCompleted bool) ([]Reminder, error) {
	now := time.Now()
	loc := now.Location()
	endOfDay := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, loc)

	opts := []ekreminders.ListOption{
		ekreminders.WithDueBefore(endOfDay),
	}
	if !showCompleted {
		opts = append(opts, ekreminders.WithCompleted(false))
	}
	items, err := c.ek.Reminders(opts...)
	if err != nil {
		return nil, err
	}
	result := convertReminders(items)
	if showCompleted {
		sortReminders(result)
	}
	return result, nil
}

// ScheduledReminders returns all reminders that have a due date, across all lists.
func (c *Client) ScheduledReminders(showCompleted bool) ([]Reminder, error) {
	// Use a distant past date to get all reminders with any due date
	distantPast := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	opts := []ekreminders.ListOption{
		ekreminders.WithDueAfter(distantPast),
	}
	if !showCompleted {
		opts = append(opts, ekreminders.WithCompleted(false))
	}
	items, err := c.ek.Reminders(opts...)
	if err != nil {
		return nil, err
	}
	result := convertReminders(items)
	// Sort by due date for scheduled view
	sort.SliceStable(result, func(i, j int) bool {
		a, b := result[i], result[j]
		if a.Completed != b.Completed {
			return !a.Completed
		}
		if a.DueDate != nil && b.DueDate != nil {
			return a.DueDate.Before(*b.DueDate)
		}
		return a.DueDate != nil
	})
	return result, nil
}

// sortReminders preserves original order but moves completed items after incomplete.
func sortReminders(items []Reminder) {
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Completed != items[j].Completed {
			return !items[i].Completed
		}
		return false
	})
}

func (c *Client) CreateReminder(input CreateReminderInput) (*Reminder, error) {
	ekInput := ekreminders.CreateReminderInput{
		Title:    input.Title,
		ListName: input.ListName,
		Priority: ekreminders.Priority(input.Priority),
		Notes:    input.Notes,
	}
	if input.DueDate != nil {
		ekInput.DueDate = input.DueDate
	}
	result, err := c.ek.CreateReminder(ekInput)
	if err != nil {
		return nil, err
	}
	r := convertReminder(*result)
	return &r, nil
}

func (c *Client) UpdateReminder(id string, input UpdateReminderInput) (*Reminder, error) {
	ekInput := ekreminders.UpdateReminderInput{}
	if input.Title != nil {
		ekInput.Title = input.Title
	}
	if input.Notes != nil {
		ekInput.Notes = input.Notes
	}
	if input.ClearDueDate {
		ekInput.ClearDueDate = true
	} else if input.DueDate != nil {
		ekInput.DueDate = input.DueDate
	}
	if input.Priority != nil {
		p := ekreminders.Priority(*input.Priority)
		ekInput.Priority = &p
	}
	result, err := c.ek.UpdateReminder(id, ekInput)
	if err != nil {
		return nil, err
	}
	r := convertReminder(*result)
	return &r, nil
}

func (c *Client) CompleteReminder(id string) (*Reminder, error) {
	result, err := c.ek.CompleteReminder(id)
	if err != nil {
		return nil, err
	}
	r := convertReminder(*result)
	return &r, nil
}

func (c *Client) UncompleteReminder(id string) (*Reminder, error) {
	result, err := c.ek.UncompleteReminder(id)
	if err != nil {
		return nil, err
	}
	r := convertReminder(*result)
	return &r, nil
}

func (c *Client) DeleteReminder(id string) error {
	return c.ek.DeleteReminder(id)
}

func convertReminders(items []ekreminders.Reminder) []Reminder {
	result := make([]Reminder, len(items))
	for i, item := range items {
		result[i] = convertReminder(item)
	}
	return result
}

func convertReminder(r ekreminders.Reminder) Reminder {
	return Reminder{
		ID:             r.ID,
		Title:          r.Title,
		Notes:          r.Notes,
		ListID:         r.ListID,
		DueDate:        r.DueDate,
		Completed:      r.Completed,
		CompletionDate: r.CompletionDate,
		Priority:       int(r.Priority),
		CreatedAt:      r.CreatedAt,
		ModifiedAt:     r.ModifiedAt,
	}
}

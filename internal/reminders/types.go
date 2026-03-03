package reminders

import "time"

const (
	PriorityNone   = 0
	PriorityHigh   = 1
	PriorityMedium = 5
	PriorityLow    = 9
)

type ListKind int

const (
	ListNormal ListKind = iota
	ListSmart
	ListSeparator
)

type ReminderList struct {
	ID    string
	Title string
	Count int
	Kind  ListKind
}

type Reminder struct {
	ID             string
	Title          string
	Notes          string
	ListID         string
	DueDate        *time.Time
	Completed      bool
	CompletionDate *time.Time
	Priority       int
	CreatedAt      *time.Time
	ModifiedAt     *time.Time
}

type CreateReminderInput struct {
	Title    string
	ListName string
	DueDate  *time.Time
	Priority int
	Notes    string
}

type UpdateReminderInput struct {
	Title        *string
	Notes        *string
	DueDate      *time.Time
	ClearDueDate bool
	Priority     *int
}

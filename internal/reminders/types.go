package reminders

import "time"

const (
	PriorityNone   = 0
	PriorityHigh   = 1
	PriorityMedium = 5
	PriorityLow    = 9
)

type SortMode int

const (
	SortDefault SortMode = iota
	SortCreated
	SortDueDate
	SortTitle
)

func (s SortMode) Next() SortMode {
	return (s + 1) % 4
}

func (s SortMode) Label() string {
	switch s {
	case SortCreated:
		return "created"
	case SortDueDate:
		return "due date"
	case SortTitle:
		return "title"
	default:
		return "default"
	}
}

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

package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Colors — teal accent (Apple Reminders aesthetic)
	Teal    = lipgloss.Color("#5AC8FA")
	DimGray = lipgloss.Color("#555555")
	White   = lipgloss.Color("#FFFFFF")
	Green   = lipgloss.Color("#34C759")
	Red     = lipgloss.Color("#FF3B30")
	Orange  = lipgloss.Color("#FF9500")
	Yellow  = lipgloss.Color("#FFCC00")
	Blue    = lipgloss.Color("#007AFF")

	// Panel borders
	FocusedBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Teal)

	UnfocusedBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(DimGray)

	// Panel titles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Teal).
			Padding(0, 1)

	// Reminder checkbox icons
	CheckboxEmpty   = lipgloss.NewStyle().Foreground(DimGray).SetString("☐")
	CheckboxChecked = lipgloss.NewStyle().Foreground(Green).SetString("☑")

	// Priority indicators
	PriorityHighStyle   = lipgloss.NewStyle().Foreground(Red).SetString("!!!")
	PriorityMediumStyle = lipgloss.NewStyle().Foreground(Orange).SetString("!!")
	PriorityLowStyle    = lipgloss.NewStyle().Foreground(Yellow).SetString("!")

	// Reminder text
	ReminderTitleStyle     = lipgloss.NewStyle().Foreground(White)
	ReminderDimStyle       = lipgloss.NewStyle().Foreground(DimGray)
	ReminderCompletedStyle = lipgloss.NewStyle().Foreground(DimGray).Strikethrough(true)
	ReminderOverdueStyle   = lipgloss.NewStyle().Foreground(Red)

	// Error & info
	ErrorStyle = lipgloss.NewStyle().Foreground(Red)

	// Spinner
	SpinnerStyle = lipgloss.NewStyle().Foreground(Teal)

	// Help overlay
	HelpOverlayStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(Teal).
				Padding(1, 2)

	HelpKeyStyle    = lipgloss.NewStyle().Foreground(Teal).Bold(true).Width(16)
	HelpDescStyle   = lipgloss.NewStyle().Foreground(White)
	HelpHeaderStyle = lipgloss.NewStyle().Foreground(Teal).Bold(true).
			Underline(true).
			MarginBottom(1)

	// Dialog styles
	DialogStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Teal).
			Padding(1, 2).
			Width(55)

	DialogTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(Teal).
				MarginBottom(1)
)

func CheckboxIcon(completed bool) string {
	if completed {
		return CheckboxChecked.String()
	}
	return CheckboxEmpty.String()
}

func PriorityIcon(priority int) string {
	switch priority {
	case 1:
		return " " + PriorityHighStyle.String()
	case 5:
		return " " + PriorityMediumStyle.String()
	case 9:
		return " " + PriorityLowStyle.String()
	default:
		return ""
	}
}

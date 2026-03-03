package ui

const (
	statusBarHeight = 1
	borderSize      = 2
	minPanelWidth   = 20
)

type Layout struct {
	ListsWidth     int
	RemindersWidth int
	PanelHeight    int
	TotalWidth     int
	TotalHeight    int
}

func ComputeLayout(width, height int) Layout {
	l := Layout{
		TotalWidth:  width,
		TotalHeight: height,
	}

	usableHeight := height - statusBarHeight
	l.PanelHeight = usableHeight - borderSize
	if l.PanelHeight < 3 {
		l.PanelHeight = 3
	}

	// 30% / 70% split
	usableWidth := width - 4 // subtract borders for 2 panels
	l.ListsWidth = usableWidth * 30 / 100
	l.RemindersWidth = usableWidth - l.ListsWidth

	if l.ListsWidth < minPanelWidth {
		l.ListsWidth = minPanelWidth
	}
	if l.RemindersWidth < minPanelWidth {
		l.RemindersWidth = minPanelWidth
	}

	return l
}

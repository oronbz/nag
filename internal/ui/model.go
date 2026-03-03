package ui

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oronbz/nag/internal/reminders"
	"github.com/oronbz/nag/internal/ui/commands"
	"github.com/oronbz/nag/internal/ui/components/dialog"
	"github.com/oronbz/nag/internal/ui/components/helpoverlay"
	"github.com/oronbz/nag/internal/ui/components/listpanel"
	"github.com/oronbz/nag/internal/ui/components/reminderpanel"
	"github.com/oronbz/nag/internal/ui/components/statusbar"
	"github.com/oronbz/nag/internal/ui/messages"
	"github.com/oronbz/nag/internal/ui/styles"
)

type Panel int

const (
	PanelLists Panel = iota
	PanelReminders
)

type Model struct {
	client *reminders.Client

	listPanel     listpanel.Model
	reminderPanel reminderpanel.Model
	statusBar     statusbar.Model
	helpOverlay   helpoverlay.Model
	createDlg     dialog.CreateModel
	createListDlg dialog.CreateListModel
	confirmDlg    dialog.ConfirmModel
	spinner       spinner.Model

	focusedPanel  Panel
	selectedList  *reminders.ReminderList
	showCompleted bool
	layout        Layout
	pendingFocus  bool
	width         int
	height        int
	ready         bool
}

func NewModel(client *reminders.Client) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.SpinnerStyle

	lp := listpanel.New(30, 20)
	lp.SetFocused(true)

	return Model{
		client:        client,
		listPanel:     lp,
		reminderPanel: reminderpanel.New(50, 20),
		statusBar:     statusbar.New(),
		helpOverlay:   helpoverlay.New(),
		createDlg:     dialog.NewCreate(),
		createListDlg: dialog.NewCreateList(),
		confirmDlg:    dialog.NewConfirm(),
		spinner:       s,
		focusedPanel:  PanelLists,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		commands.FetchLists(m.client),
		commands.AutoRefreshTick(),
		m.spinner.Tick,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.resize()
		return m, nil

	case tea.MouseMsg:
		if m.createDlg.Visible() || m.createListDlg.Visible() || m.confirmDlg.Visible() || m.helpOverlay.Visible() {
			break
		}
		target, ok := m.panelForMouse(msg.X, msg.Y)
		if !ok {
			break
		}
		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
			m.setFocus(target)
			return m, nil
		}
		if msg.Button == tea.MouseButtonWheelUp || msg.Button == tea.MouseButtonWheelDown {
			var cmd tea.Cmd
			switch target {
			case PanelLists:
				m.listPanel, cmd = m.listPanel.Update(msg)
			case PanelReminders:
				m.reminderPanel, cmd = m.reminderPanel.Update(msg)
			}
			return m, cmd
		}

	case clearInfoMsg:
		m.statusBar.ClearInfo()
		return m, nil

	case delayedRefreshMsg:
		return m, m.fetchSelectedReminders()

	case openFailedMsg:
		m.statusBar.SetError("Open failed: " + msg.err.Error())
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case messages.ListsLoadedMsg:
		m.statusBar.ClearLoading()
		if msg.Err != nil {
			m.statusBar.SetError(msg.Err.Error())
			return m, nil
		}
		m.statusBar.ClearError()
		m.listPanel.SetLists(msg.Lists)
		m.resize()
		return m, nil

	case messages.RemindersLoadedMsg:
		m.statusBar.ClearLoading()
		if msg.Err != nil {
			m.statusBar.SetError(msg.Err.Error())
			return m, nil
		}
		m.statusBar.ClearError()
		m.reminderPanel.SetReminders(msg.Reminders)
		m.resize()
		if m.pendingFocus {
			m.pendingFocus = false
			m.cycleFocus(1)
		}
		return m, nil

	case messages.ReminderCreatedMsg:
		m.statusBar.ClearLoading()
		if msg.Err != nil {
			m.statusBar.SetError("Create failed: " + msg.Err.Error())
			return m, nil
		}
		m.statusBar.ClearError()
		m.statusBar.SetInfo("Reminder created")
		cmds = append(cmds, tea.Tick(2*time.Second, func(time.Time) tea.Msg { return clearInfoMsg{} }))
		cmds = append(cmds, commands.FetchLists(m.client))
		if cmd := m.fetchSelectedReminders(); cmd != nil {
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)

	case messages.ReminderCompletedMsg:
		if msg.Err != nil {
			m.statusBar.SetError("Toggle failed: " + msg.Err.Error())
			return m, nil
		}
		m.statusBar.ClearError()
		if msg.Reminder != nil {
			m.reminderPanel.UpdateReminder(*msg.Reminder)
		}
		// Refresh lists immediately for count update; delay reminder re-fetch
		// so the user sees the toggle and can undo with Space.
		return m, tea.Batch(
			commands.FetchLists(m.client),
			tea.Tick(2*time.Second, func(time.Time) tea.Msg { return delayedRefreshMsg{} }),
		)

	case messages.ReminderDeletedMsg:
		m.statusBar.ClearLoading()
		if msg.Err != nil {
			m.statusBar.SetError("Delete failed: " + msg.Err.Error())
			return m, nil
		}
		m.statusBar.ClearError()
		m.statusBar.SetInfo("Reminder deleted")
		cmds = append(cmds, tea.Tick(2*time.Second, func(time.Time) tea.Msg { return clearInfoMsg{} }))
		cmds = append(cmds, commands.FetchLists(m.client))
		if cmd := m.fetchSelectedReminders(); cmd != nil {
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)

	case messages.TickMsg:
		batch := []tea.Cmd{commands.AutoRefreshTick(), commands.FetchLists(m.client)}
		if cmd := m.fetchSelectedReminders(); cmd != nil {
			batch = append(batch, cmd)
		}
		return m, tea.Batch(batch...)

	case dialog.CreateSubmitMsg:
		m.statusBar.SetLoading("Creating reminder...")
		return m, commands.CreateReminder(m.client, msg.Input)

	case dialog.EditSubmitMsg:
		m.statusBar.SetLoading("Updating reminder...")
		return m, commands.UpdateReminder(m.client, msg.ID, msg.Input)

	case messages.ReminderUpdatedMsg:
		m.statusBar.ClearLoading()
		if msg.Err != nil {
			m.statusBar.SetError("Update failed: " + msg.Err.Error())
			return m, nil
		}
		m.statusBar.ClearError()
		m.statusBar.SetInfo("Reminder updated")
		cmds = append(cmds, tea.Tick(2*time.Second, func(time.Time) tea.Msg { return clearInfoMsg{} }))
		cmds = append(cmds, commands.FetchLists(m.client))
		if cmd := m.fetchSelectedReminders(); cmd != nil {
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)

	case dialog.CreateListSubmitMsg:
		m.statusBar.SetLoading("Creating list...")
		return m, commands.CreateList(m.client, msg.Title)

	case dialog.EditListSubmitMsg:
		m.statusBar.SetLoading("Updating list...")
		return m, commands.UpdateList(m.client, msg.ID, msg.Title)

	case messages.ListUpdatedMsg:
		m.statusBar.ClearLoading()
		if msg.Err != nil {
			m.statusBar.SetError("Update list failed: " + msg.Err.Error())
			return m, nil
		}
		m.statusBar.ClearError()
		m.statusBar.SetInfo("List updated")
		return m, tea.Batch(
			commands.FetchLists(m.client),
			tea.Tick(2*time.Second, func(time.Time) tea.Msg { return clearInfoMsg{} }),
		)

	case messages.ListCreatedMsg:
		m.statusBar.ClearLoading()
		if msg.Err != nil {
			m.statusBar.SetError("Create list failed: " + msg.Err.Error())
			return m, nil
		}
		m.statusBar.ClearError()
		m.statusBar.SetInfo("List created")
		return m, tea.Batch(
			commands.FetchLists(m.client),
			tea.Tick(2*time.Second, func(time.Time) tea.Msg { return clearInfoMsg{} }),
		)

	case messages.ListDeletedMsg:
		m.statusBar.ClearLoading()
		if msg.Err != nil {
			m.statusBar.SetError("Delete list failed: " + msg.Err.Error())
			return m, nil
		}
		m.statusBar.ClearError()
		m.statusBar.SetInfo("List deleted")
		m.selectedList = nil
		m.reminderPanel.SetReminders(nil)
		m.reminderPanel.SetTitle("")
		return m, tea.Batch(
			commands.FetchLists(m.client),
			tea.Tick(2*time.Second, func(time.Time) tea.Msg { return clearInfoMsg{} }),
		)

	case dialog.ConfirmYesMsg:
		switch msg.Action {
		case dialog.ConfirmDelete:
			if r, ok := m.reminderPanel.SelectedReminder(); ok {
				m.statusBar.SetLoading("Deleting reminder...")
				return m, commands.DeleteReminder(m.client, r.ID)
			}
		case dialog.ConfirmDeleteList:
			if list, ok := m.listPanel.SelectedList(); ok {
				m.statusBar.SetLoading("Deleting list...")
				return m, commands.DeleteList(m.client, list.ID)
			}
		}
		return m, nil

	case dialog.ConfirmNoMsg:
		return m, nil
	}

	// Route to dialogs/overlays first if visible
	if m.createDlg.Visible() {
		var cmd tea.Cmd
		m.createDlg, cmd = m.createDlg.Update(msg)
		return m, cmd
	}
	if m.createListDlg.Visible() {
		var cmd tea.Cmd
		m.createListDlg, cmd = m.createListDlg.Update(msg)
		return m, cmd
	}
	if m.confirmDlg.Visible() {
		var cmd tea.Cmd
		m.confirmDlg, cmd = m.confirmDlg.Update(msg)
		return m, cmd
	}
	if m.helpOverlay.Visible() {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if key.Matches(keyMsg, Keys.Help) || key.Matches(keyMsg, Keys.Escape) || key.Matches(keyMsg, Keys.Quit) {
				m.helpOverlay.Toggle()
				return m, nil
			}
		}
		var cmd tea.Cmd
		m.helpOverlay, cmd = m.helpOverlay.Update(msg)
		return m, cmd
	}

	// Global keys (skip when filtering)
	filtering := m.listPanel.Filtering() || m.reminderPanel.Filtering()
	if keyMsg, ok := msg.(tea.KeyMsg); ok && !filtering {
		switch {
		case key.Matches(keyMsg, Keys.Quit):
			return m, tea.Quit
		case key.Matches(keyMsg, Keys.Help):
			m.helpOverlay.Toggle()
			return m, nil
		case key.Matches(keyMsg, Keys.Tab):
			m.cycleFocus(1)
			return m, nil
		case key.Matches(keyMsg, Keys.ShiftTab):
			m.cycleFocus(-1)
			return m, nil
		case key.Matches(keyMsg, Keys.Enter):
			return m, m.handleEnter()
		case key.Matches(keyMsg, Keys.ToggleComplete):
			return m, m.handleToggleComplete()
		case key.Matches(keyMsg, Keys.NewReminder):
			m.handleNew()
			return m, nil
		case key.Matches(keyMsg, Keys.Edit):
			m.handleEdit()
			return m, nil
		case key.Matches(keyMsg, Keys.Delete):
			m.handleDelete()
			return m, nil
		case key.Matches(keyMsg, Keys.OpenInApp):
			return m, m.handleOpenInApp()
		case key.Matches(keyMsg, Keys.ShowCompleted):
			m.showCompleted = !m.showCompleted
			if m.showCompleted {
				m.statusBar.SetInfo("Showing completed")
			} else {
				m.statusBar.SetInfo("Hiding completed")
			}
			return m, tea.Batch(m.fetchSelectedReminders(), tea.Tick(2*time.Second, func(time.Time) tea.Msg { return clearInfoMsg{} }))
		case key.Matches(keyMsg, Keys.Refresh):
			return m, m.handleRefresh()
		}
	}

	// Route to focused panel
	switch m.focusedPanel {
	case PanelLists:
		var cmd tea.Cmd
		m.listPanel, cmd = m.listPanel.Update(msg)
		cmds = append(cmds, cmd)
	case PanelReminders:
		var cmd tea.Cmd
		m.reminderPanel, cmd = m.reminderPanel.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return m.spinner.View() + " Loading..."
	}

	if m.helpOverlay.Visible() {
		return m.helpOverlay.View()
	}
	if m.createDlg.Visible() {
		return m.createDlg.View()
	}
	if m.createListDlg.Visible() {
		return m.createListDlg.View()
	}
	if m.confirmDlg.Visible() {
		return m.confirmDlg.View()
	}

	l := m.layout
	listsPanel := m.renderPanel(m.listPanel.View(), l.ListsWidth, l.PanelHeight, m.focusedPanel == PanelLists)
	remindersPanel := m.renderPanel(m.reminderPanel.View(), l.RemindersWidth, l.PanelHeight, m.focusedPanel == PanelReminders)

	panels := lipgloss.JoinHorizontal(lipgloss.Top, listsPanel, remindersPanel)
	return lipgloss.JoinVertical(lipgloss.Left, panels, m.statusBar.View())
}

func (m Model) renderPanel(content string, width, height int, focused bool) string {
	style := styles.UnfocusedBorder
	if focused {
		style = styles.FocusedBorder
	}
	return style.Width(width).Height(height).MaxHeight(height + 2).Render(content)
}

func (m *Model) resize() {
	m.layout = ComputeLayout(m.width, m.height)
	m.listPanel.SetSize(m.layout.ListsWidth, m.layout.PanelHeight)
	m.reminderPanel.SetSize(m.layout.RemindersWidth, m.layout.PanelHeight)
	m.helpOverlay.SetSize(m.width, m.height)
	m.createDlg.SetSize(m.width, m.height)
	m.createListDlg.SetSize(m.width, m.height)
	m.confirmDlg.SetSize(m.width, m.height)
}

func (m *Model) setFocus(panel Panel) {
	if m.focusedPanel == panel {
		return
	}
	m.listPanel.SetFocused(false)
	m.reminderPanel.SetFocused(false)

	m.focusedPanel = panel
	switch panel {
	case PanelLists:
		m.listPanel.SetFocused(true)
		m.statusBar.SetPanel(statusbar.PanelLists)
	case PanelReminders:
		m.reminderPanel.SetFocused(true)
		m.statusBar.SetPanel(statusbar.PanelReminders)
	}
}

func (m Model) panelForMouse(x, _ int) (Panel, bool) {
	l := m.layout
	if x < l.ListsWidth+2 {
		return PanelLists, true
	}
	return PanelReminders, true
}

func (m *Model) cycleFocus(dir int) {
	m.listPanel.SetFocused(false)
	m.reminderPanel.SetFocused(false)

	m.focusedPanel = Panel((int(m.focusedPanel) + dir + 2) % 2)

	switch m.focusedPanel {
	case PanelLists:
		m.listPanel.SetFocused(true)
		m.statusBar.SetPanel(statusbar.PanelLists)
	case PanelReminders:
		m.reminderPanel.SetFocused(true)
		m.statusBar.SetPanel(statusbar.PanelReminders)
	}
}

func (m *Model) handleEnter() tea.Cmd {
	switch m.focusedPanel {
	case PanelLists:
		if list, ok := m.listPanel.SelectedList(); ok {
			m.selectedList = &list
			m.pendingFocus = true
			m.reminderPanel.SetTitle(list.Title)
			m.statusBar.SetLoading("Loading reminders...")
			return m.fetchSelectedReminders()
		}
	}
	return nil
}

func (m *Model) fetchSelectedReminders() tea.Cmd {
	if m.selectedList == nil {
		return nil
	}
	switch m.selectedList.ID {
	case reminders.SmartListToday:
		return commands.FetchTodayReminders(m.client, m.showCompleted)
	case reminders.SmartListScheduled:
		return commands.FetchScheduledReminders(m.client, m.showCompleted)
	default:
		return commands.FetchReminders(m.client, m.selectedList.Title, m.showCompleted)
	}
}

func (m *Model) handleToggleComplete() tea.Cmd {
	if m.focusedPanel != PanelReminders {
		return nil
	}
	if r, ok := m.reminderPanel.SelectedReminder(); ok {
		return commands.ToggleComplete(m.client, r.ID, r.Completed)
	}
	return nil
}

func (m *Model) handleNew() {
	switch m.focusedPanel {
	case PanelLists:
		m.createListDlg.Show()
	case PanelReminders:
		if m.selectedList == nil {
			m.statusBar.SetInfo("Select a list first")
			return
		}
		if m.selectedList.Kind == reminders.ListSmart {
			m.statusBar.SetInfo("Select a regular list to create reminders")
			return
		}
		m.createDlg.Show(m.selectedList.Title)
	}
}

func (m *Model) handleEdit() {
	switch m.focusedPanel {
	case PanelLists:
		list, ok := m.listPanel.SelectedList()
		if !ok {
			return
		}
		if list.Kind == reminders.ListSmart || list.Kind == reminders.ListSeparator {
			m.statusBar.SetInfo("Cannot edit smart lists")
			return
		}
		m.createListDlg.ShowEdit(list.ID, list.Title)
	case PanelReminders:
		if r, ok := m.reminderPanel.SelectedReminder(); ok {
			m.createDlg.ShowEdit(r)
		}
	}
}

func (m *Model) handleDelete() {
	switch m.focusedPanel {
	case PanelLists:
		list, ok := m.listPanel.SelectedList()
		if !ok {
			return
		}
		if list.Kind == reminders.ListSmart || list.Kind == reminders.ListSeparator {
			m.statusBar.SetInfo("Cannot delete smart lists")
			return
		}
		m.confirmDlg.Show(
			fmt.Sprintf("Delete list \"%s\" and all its reminders?", list.Title),
			dialog.ConfirmDeleteList,
		)
	case PanelReminders:
		if r, ok := m.reminderPanel.SelectedReminder(); ok {
			m.confirmDlg.Show(
				fmt.Sprintf("Delete \"%s\"?", r.Title),
				dialog.ConfirmDelete,
			)
		}
	}
}

func (m *Model) handleRefresh() tea.Cmd {
	m.statusBar.SetLoading("Refreshing...")
	batch := []tea.Cmd{commands.FetchLists(m.client)}
	if cmd := m.fetchSelectedReminders(); cmd != nil {
		batch = append(batch, cmd)
	}
	return tea.Batch(batch...)
}

func (m *Model) handleOpenInApp() tea.Cmd {
	if m.focusedPanel != PanelReminders {
		return nil
	}
	r, ok := m.reminderPanel.SelectedReminder()
	if !ok {
		return nil
	}
	return func() tea.Msg {
		err := exec.Command("open", "x-apple-reminderkit://REMCDReminder/"+r.ID).Run()
		if err != nil {
			return openFailedMsg{err}
		}
		return nil
	}
}

type clearInfoMsg struct{}
type delayedRefreshMsg struct{}
type openFailedMsg struct{ err error }

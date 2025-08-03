package src

import (
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//go:embed art.ascii
var ascii string

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type Model struct {
	width         int
	height        int
	list          list.Model
	db            *sql.DB
	keys          *listKeyMap
	creatingNew   bool
	editing       bool
	confirmDelete bool
	textInput     textinput.Model
	editingTaskID int
}

func NewModel(db *sql.DB) Model {
	records, err := LoadTasks(db)
	if err != nil {
		panic(fmt.Sprintf("failed to load tasks: %v", err))
	}

	var items []list.Item
	for _, r := range records {
		items = append(items, NewTaskItemFromRecord(r))
	}

	keys := newListKeyMap()
	delegate := NewTaskItemDelegate(keys)

	l := list.New(items, delegate, 0, 0)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.DisableQuitKeybindings()
	l.SetShowHelp(true)
	l.SetShowTitle(true)
	l.Title = "Tasks"
	l.Styles.Title = TitleStyle
	l.Styles.NoItems = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Height(0)

	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.newTask,
			keys.editTask,
			keys.deleteTask,
			keys.toggleTask,
			keys.quit,
		}
	}

	ti := textinput.New()
	ti.Placeholder = "New task title"
	ti.Focus()
	ti.CharLimit = 64
	ti.Width = 30

	return Model{
		list:          l,
		db:            db,
		keys:          keys,
		creatingNew:   false,
		editing:       false,
		confirmDelete: false,
		textInput:     ti,
		editingTaskID: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.list.SetSize(m.width/3, m.height-4)

	case tea.KeyMsg:
		if m.creatingNew {
			switch msg.String() {
			case "enter":
				title := m.textInput.Value()
				if title != "" {
					id, err := InsertTask(m.db, title)
					if err == nil {
						newTask := NewTaskItemFromRecord(TaskRecord{
							ID:      int(id),
							Title:   title,
							Elapsed: 0,
						})
						m.list.InsertItem(len(m.list.Items()), newTask)
						statusCmd := m.list.NewStatusMessage(StatusMessageStyle("Added task: " + title))
						cmds = append(cmds, statusCmd)
					}
				}
				m.creatingNew = false
				m.textInput.Reset()
			case "esc":
				m.creatingNew = false
				m.textInput.Reset()
			default:
				var cmd tea.Cmd
				m.textInput, cmd = m.textInput.Update(msg)
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		}

		if m.confirmDelete {
			switch msg.String() {
			case "y":
				if item, ok := m.list.SelectedItem().(TaskItem); ok {
					err := DeleteTask(m.db, item.id)
					if err == nil {
						m.list.RemoveItem(m.list.Index())
						statusCmd := m.list.NewStatusMessage(StatusMessageStyle("Deleted task: " + item.title))
						cmds = append(cmds, statusCmd)
					}
				}
				m.confirmDelete = false
				return m, tea.Batch(cmds...)
			case "n", "esc":
				m.confirmDelete = false
				return m, nil
			}
			return m, nil
		}

		if m.editing {
			switch msg.String() {
			case "enter":
				title := m.textInput.Value()
				if title != "" {
					err := UpdateTask(m.db, m.editingTaskID, title)
					if err == nil {
						if item, ok := m.list.SelectedItem().(TaskItem); ok {
							item.title = title
							m.list.SetItem(m.list.Index(), item)
							statusCmd := m.list.NewStatusMessage(StatusMessageStyle("Updated task: " + title))
							cmds = append(cmds, statusCmd)
						}
					}
				}
				m.editing = false
				m.editingTaskID = 0
				m.textInput.Reset()
			case "esc":
				m.editing = false
				m.editingTaskID = 0
				m.textInput.Reset()
			default:
				var cmd tea.Cmd
				m.textInput, cmd = m.textInput.Update(msg)
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		}

		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.newTask):
			m.creatingNew = true
			m.textInput.Focus()
			return m, nil

		case key.Matches(msg, m.keys.editTask):
			if item, ok := m.list.SelectedItem().(TaskItem); ok {
				m.editing = true
				m.editingTaskID = item.id
				m.textInput.SetValue(item.title)
				m.textInput.Focus()
			}
			return m, nil

		case key.Matches(msg, m.keys.deleteTask):
			if _, ok := m.list.SelectedItem().(TaskItem); ok {
				m.confirmDelete = true
			}
			return m, nil

		case key.Matches(msg, m.keys.toggleTask):
			if item, ok := m.list.SelectedItem().(TaskItem); ok {
				if item.timer == nil {
					item.timer = NewTimer()
				}
				item.timer.Toggle()
				m.list.SetItem(m.list.Index(), item)

				if !item.timer.IsRunning() {
					_ = UpdateTaskElapsed(m.db, item.id, item.timer.Elapsed().Nanoseconds())
					statusCmd := m.list.NewStatusMessage(StatusMessageStyle("Stopped timer for: " + item.title))
					cmds = append(cmds, statusCmd)
				} else {
					statusCmd := m.list.NewStatusMessage(StatusMessageStyle("Started timer for: " + item.title))
					cmds = append(cmds, statusCmd)
					cmds = append(cmds, tick())
				}
			}

		case key.Matches(msg, m.keys.quit):
			for i, listItem := range m.list.Items() {
				item := listItem.(TaskItem)
				if item.timer != nil && item.timer.IsRunning() {
					item.timer.Stop()
					_ = UpdateTaskElapsed(m.db, item.id, item.timer.Elapsed().Nanoseconds())
					m.list.SetItem(i, item)
				}
			}
			return m, tea.Quit
		}

	case tickMsg:
		for _, listItem := range m.list.Items() {
			item := listItem.(TaskItem)
			if item.timer != nil && item.timer.IsRunning() {
				cmds = append(cmds, tick())
				break
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	contentWidth := m.width - 6
	contentHeight := m.height - 3
	taskWidth := contentWidth / 3
	timerWidth := 2 * contentWidth / 3

	tasksView := BorderStyle.Width(taskWidth).Height(contentHeight).
		Render(m.list.View())
	timerContent := ""
	var item TaskItem
	if selected := m.list.SelectedItem(); selected != nil {
		item = selected.(TaskItem)
	}
	timerStyle := TimerStyle

	if len(m.list.Items()) == 0 {
		timerContent = LogoStyle.Render(ascii)
	} else {
		timeStr := "00:00:00"
		if item.timer != nil {
			timeStr = item.timer.String()
			if item.timer.IsRunning() {
				timerStyle = timerStyle.Foreground(successColor)
			} else if item.timer.Elapsed() > 0 {
				timerStyle = timerStyle.Foreground(pauseColor)
			} else {
				timerStyle = timerStyle.Foreground(errorColor)
			}
		}
		timerContent = timerStyle.Render(timerBox(timeStr))
	}

	timerView := BorderStyle.Width(timerWidth).Height(contentHeight).
		Align(lipgloss.Center, lipgloss.Center).
		Render(timerContent)

	mainView := lipgloss.JoinHorizontal(lipgloss.Top, tasksView, "  ", timerView)

	if m.creatingNew {
		popupWidth := 40
		popupHeight := 7
		popup := lipgloss.NewStyle().
			Width(popupWidth).
			Height(popupHeight).
			Border(lipgloss.RoundedBorder()).
			Align(lipgloss.Center, lipgloss.Center).
			Padding(1, 2).
			Render(fmt.Sprintf("New Task:\n\n%s\n\nesc to cancel", m.textInput.View()))

		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, popup)
	}

	if m.editing {
		popupWidth := 40
		popupHeight := 7
		popup := lipgloss.NewStyle().
			Width(popupWidth).
			Height(popupHeight).
			Border(lipgloss.RoundedBorder()).
			Align(lipgloss.Center, lipgloss.Center).
			Padding(1, 2).
			Render(fmt.Sprintf("Edit Task:\n\n%s\n\nesc to cancel", m.textInput.View()))

		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, popup)
	}

	if m.confirmDelete {
		if item, ok := m.list.SelectedItem().(TaskItem); ok {
			confirmMsg := fmt.Sprintf("Delete task '%s'?\n\n(y/n)\n\nesc to cancel", item.title)
			popup := lipgloss.NewStyle().
				Width(40).
				Height(7).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(errorColor).
				Align(lipgloss.Center, lipgloss.Center).
				Padding(1, 2).
				Render(confirmMsg)

			return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, popup)
		}
	}

	return mainView
}

func timerBox(s string) string {
	return lipgloss.JoinVertical(lipgloss.Center,
		"▀▀▀▀▀ ▀▀▀▀▀ ▀▀▀▀▀",
		s,
		"▄▄▄▄▄ ▄▄▄▄▄ ▄▄▄▄▄",
	)
}

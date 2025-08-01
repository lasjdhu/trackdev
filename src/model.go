package src

import (
	_ "embed"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//go:embed art.ascii
var ascii string

type tickMsg time.Time

type model struct {
	width       int
	height      int
	tasks       []TaskItem
	selectedIdx int
	showTimer   bool
}

func NewModel() model {
	return model{
		tasks: []TaskItem{
			NewTaskItem("Fix bug #123"),
			NewTaskItem("Write tests"),
			NewTaskItem("Code review"),
			NewTaskItem("Implement feature X"),
			NewTaskItem("Refactor Y"),
		},
		selectedIdx: 0,
		showTimer:   true,
	}
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
		case "j":
			if m.selectedIdx < len(m.tasks)-1 {
				m.selectedIdx++
			}
		case "v":
			m.showTimer = !m.showTimer
			return m, nil
		case " ":
			if m.selectedIdx >= 0 && m.selectedIdx < len(m.tasks) {
				if m.tasks[m.selectedIdx].timer == nil {
					m.tasks[m.selectedIdx].timer = NewTimer()
				}
				m.tasks[m.selectedIdx].timer.Toggle()
				if m.tasks[m.selectedIdx].timer.IsRunning() {
					cmds = append(cmds, tick())
				}
			}
			return m, tea.Batch(cmds...)
		}

	case tickMsg:
		needsTick := false
		for _, task := range m.tasks {
			if task.timer != nil && task.timer.IsRunning() {
				needsTick = true
				break
			}
		}
		if needsTick {
			cmds = append(cmds, tick())
		}
		return m, tea.Batch(cmds...)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	contentWidth := m.width - 2
	contentHeight := m.height - 4
	taskWidth := (contentWidth / 3) - 2
	timerWidth := ((2 * contentWidth) / 3) - 4

	var taskLines []string
	for i, task := range m.tasks {
		style := ItemStyle
		if i == m.selectedIdx {
			style = SelectedStyle
		}
		taskLines = append(taskLines, style.Render(task.title))
	}

	tasksContent := lipgloss.JoinVertical(lipgloss.Left, taskLines...)

	tasksView := BorderStyle.
		Width(taskWidth).
		Height(contentHeight).
		Render(tasksContent)

	var timerContent string
	if m.selectedIdx >= 0 && m.selectedIdx < len(m.tasks) {
		task := m.tasks[m.selectedIdx]
		timeStr := "00:00:00"
		timerStyle := TimerStyle

		if task.timer != nil {
			timeStr = task.timer.String()
			if task.timer.IsRunning() {
				timerStyle = timerStyle.Foreground(successColor)
			} else if task.timer.Elapsed() > 0 {
				timerStyle = timerStyle.Foreground(pauseColor)
			} else {
				timerStyle = timerStyle.Foreground(errorColor)
			}
		} else {
			timerStyle = timerStyle.Foreground(errorColor)
		}

		if m.showTimer {
			timerStr := timerWithBorder(timeStr)
			timerContent = timerStyle.Render(timerStr)
		} else {
			timerContent = LogoStyle.Render(getAsciiArt())
		}
	}

	timerView := BorderStyle.
		Width(timerWidth).
		Height(contentHeight).
		Align(lipgloss.Center, lipgloss.Center).
		Render(timerContent)

	tasksSection := lipgloss.JoinVertical(lipgloss.Left,
		TitleStyle.Render("Tasks"),
		tasksView,
	)

	timerSection := lipgloss.JoinVertical(lipgloss.Left,
		TitleStyle.Render("Timer"),
		timerView,
	)

	mainView := lipgloss.NewStyle().
		MarginLeft(1).
		MarginRight(1).
		Render(lipgloss.JoinHorizontal(lipgloss.Top, tasksSection, "  ", timerSection))

	legend := lipgloss.JoinHorizontal(lipgloss.Top,
		KeyStyle.Render(NavigateKeys)+" to navigate",
		"  ",
		KeyStyle.Render(PauseResumeKey)+" to toggle timer",
		"  ",
		KeyStyle.Render(ToggleViewKey)+" to toggle view",
		"  ",
		QuitStyle.Render(QuitKey)+" to quit")

	legendView := lipgloss.NewStyle().
		Width(m.width).
		Height(1).
		Align(lipgloss.Center, lipgloss.Center).
		Render(legend)

	return lipgloss.JoinVertical(lipgloss.Left, mainView, legendView)
}

func timerWithBorder(s string) string {
	lines := []string{
		"▀▀▀▀▀ ▀▀▀▀▀ ▀▀▀▀▀",
		s,
		"▄▄▄▄▄ ▄▄▄▄▄ ▄▄▄▄▄",
	}
	return lipgloss.JoinVertical(lipgloss.Center, lines...)
}

func getAsciiArt() string {
	return ascii
}


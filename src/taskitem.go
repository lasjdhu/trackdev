package src

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type TaskItem struct {
	title  string
	timer  *Timer
	active bool
}

func NewTaskItem(title string) TaskItem {
	return TaskItem{
		title: title,
		timer: NewTimer(),
	}
}

func (t TaskItem) FilterValue() string { return t.title }
func (t TaskItem) Title() string       { return t.title }
func (t TaskItem) Description() string {
	if t.timer == nil {
		return "Not started"
	}
	elapsed := t.timer.Elapsed().Round(time.Second)
	status := "Paused"
	if t.timer.IsRunning() {
		status = "Running"
	}
	return fmt.Sprintf("⏱️  %s - %s", status, elapsed)
}

type TaskList struct {
	list list.Model
}

func NewTaskList() TaskList {
	items := []list.Item{
		NewTaskItem("Fix bug #123"),
		NewTaskItem("Write tests"),
		NewTaskItem("Code review"),
		NewTaskItem("Implement feature X"),
		NewTaskItem("Refactor Y"),
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = SelectedStyle
	delegate.Styles.SelectedDesc = SelectedStyle

	const defaultWidth = 40
	l := list.New(items, delegate, defaultWidth, 14)
	l.Title = "Tasks"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.Styles.Title = TitleStyle
	l.Styles.TitleBar = ListStyle
	l.KeyMap.Quit.Unbind()

	return TaskList{
		list: l,
	}
}

func (t TaskList) Init() tea.Cmd {
	return nil
}

func (t TaskList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	t.list, cmd = t.list.Update(msg)
	return t, cmd
}

func (t TaskList) View() string {
	return t.list.View()
}

func (t *TaskList) GetSelectedTask() (TaskItem, bool) {
	item := t.list.SelectedItem()
	if item == nil {
		return TaskItem{}, false
	}
	if taskItem, ok := item.(TaskItem); ok {
		return taskItem, true
	}
	return TaskItem{}, false
}

func (t *TaskList) UpdateSelectedTask(item TaskItem) {
	t.list.SetItem(t.list.Index(), item)
}

func (t *TaskList) SetSize(width, height int) {
	t.list.SetSize(width, height)
}

func (t *TaskList) ForEachTask(fn func(TaskItem) bool) {
	items := t.list.Items()
	for i := 0; i < len(items); i++ {
		if item, ok := items[i].(TaskItem); ok {
			if !fn(item) {
				break
			}
		}
	}
}


package src

import (
	"fmt"
	"io"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type TaskItem struct {
	id    int
	title string
	timer *Timer
}

func NewTaskItemFromRecord(t TaskRecord) TaskItem {
	timer := NewTimer()
	timer.elapsed = time.Duration(t.Elapsed)
	return TaskItem{
		id:    t.ID,
		title: t.Title,
		timer: timer,
	}
}

func (t TaskItem) Title() string {
	return t.title
}

func (t TaskItem) Description() string {
	status := "🔴"
	if t.timer != nil {
		if t.timer.IsRunning() {
			status = "🟢"
		} else if t.timer.Elapsed() > 0 {
			status = "🟡"
		}
	}
	return status
}

func (t TaskItem) FilterValue() string {
	return t.title
}

type TaskItemDelegate struct {
	keys *listKeyMap
}

func NewTaskItemDelegate(keys *listKeyMap) TaskItemDelegate {
	return TaskItemDelegate{keys: keys}
}

func (d TaskItemDelegate) Height() int                         { return 1 }
func (d TaskItemDelegate) Spacing() int                        { return 1 }
func (d TaskItemDelegate) Update(tea.Msg, *list.Model) tea.Cmd { return nil }

func (d TaskItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	t, ok := item.(TaskItem)
	if !ok {
		return
	}

	style := ItemStyle
	if index == m.Index() {
		style = SelectedStyle
	}

	var spinner string
	if t.timer != nil && t.timer.IsRunning() {
		spinner = "⏱ "
	}

	line := fmt.Sprintf("%s%s %s", spinner, t.Description(), t.Title())
	_, err := fmt.Fprint(w, style.Width(m.Width()).Render(line))
	if err != nil {
		return
	}
}

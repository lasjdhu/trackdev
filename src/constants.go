package src

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

type listKeyMap struct {
	newTask    key.Binding
	editTask   key.Binding
	deleteTask key.Binding
	toggleTask key.Binding
	quit       key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		newTask: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new task"),
		),
		editTask: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit task"),
		),
		deleteTask: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete task"),
		),
		toggleTask: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle task timer"),
		),
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

var (
	errorColor   = lipgloss.Color("1")
	successColor = lipgloss.Color("2")
	pauseColor   = lipgloss.Color("3")
	primaryColor = lipgloss.Color("5")

	BorderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder())

	ItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	SelectedStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(primaryColor).
			Bold(true)

	TitleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1)

	TimerStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(1, 0).
			MarginTop(2)

	LogoStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Align(lipgloss.Center)

	StatusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

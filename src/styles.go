package src

import "github.com/charmbracelet/lipgloss"

var (
	errorColor     = lipgloss.Color("1")
	successColor   = lipgloss.Color("2")
	pauseColor     = lipgloss.Color("3")
	secondaryColor = lipgloss.Color("4")
	primaryColor   = lipgloss.Color("5")

	BorderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor)

	KeyStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)

	QuitStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	ItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	SelectedStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(primaryColor).
			Bold(true)

	TitleStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)

	TimerStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(1, 0).
			MarginTop(2)

	ListStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor)

	LogoStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Align(lipgloss.Center)
)

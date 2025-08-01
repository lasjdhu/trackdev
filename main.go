package main

import (
	"fmt"
	"os"

	"github.com/lasjdhu/trackdev/src"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := src.NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

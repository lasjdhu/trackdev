package main

import (
	"fmt"
	"os"

	"github.com/lasjdhu/trackdev/src"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	db, err := src.InitDB("tasks.sqlite")
	if err != nil {
		fmt.Println("Failed to open DB:", err)
		os.Exit(1)
	}

	m := src.NewModel(db)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

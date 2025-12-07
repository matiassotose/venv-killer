package main

import (
	"fmt"
	"os"
	"venv-killer/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	p := tea.NewProgram(ui.NewModel(root))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

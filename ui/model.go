package ui

import (
	"fmt"
	"venv-killer/deleter"
	"venv-killer/scanner"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type State int

const (
	StateScanning State = iota
	StateList
	StateDeleting
)

type Model struct {
	state         State
	venvs         []scanner.Venv
	cursor        int
	selected      map[int]struct{}
	spinner       spinner.Model
	progress      progress.Model
	deletedCount  int
	totalToDelete int
	err           error
	width         int
	height        int
	root          string
}

func NewModel(root string) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	p := progress.New(progress.WithDefaultGradient())

	return Model{
		state:    StateScanning,
		spinner:  s,
		progress: p,
		selected: make(map[int]struct{}),
		root:     root,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, scanCmd(m.root))
}

func scanCmd(root string) tea.Cmd {
	return func() tea.Msg {
		venvs, err := scanner.Scan(root)
		if err != nil {
			return err
		}
		return venvs
	}
}

type deleteMsg struct {
	index int
	err   error
}

func deleteCmd(path string, index int) tea.Cmd {
	return func() tea.Msg {
		err := deleter.Delete(path)
		return deleteMsg{index: index, err: err}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.state == StateDeleting {
			return m, nil // Ignore keys while deleting
		}
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.venvs)-1 {
				m.cursor++
			}
		case " ":
			if _, ok := m.selected[m.cursor]; ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "enter":
			if len(m.selected) > 0 {
				m.state = StateDeleting
				m.totalToDelete = len(m.selected)
				m.deletedCount = 0
				var cmds []tea.Cmd
				for i := range m.selected {
					cmds = append(cmds, deleteCmd(m.venvs[i].Path, i))
				}
				return m, tea.Batch(cmds...)
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = msg.Width - 4
	case []scanner.Venv:
		m.venvs = msg
		m.state = StateList
	case deleteMsg:
		m.deletedCount++
		cmd := m.progress.SetPercent(float64(m.deletedCount) / float64(m.totalToDelete))

		if m.deletedCount >= m.totalToDelete {
			// Done deleting
			var newVenvs []scanner.Venv
			for i, v := range m.venvs {
				if _, ok := m.selected[i]; !ok {
					newVenvs = append(newVenvs, v)
				}
			}
			m.venvs = newVenvs
			m.selected = make(map[int]struct{})
			m.cursor = 0
			m.state = StateList
			// Reset progress
			m.progress.SetPercent(0)
			return m, cmd
		}
		return m, cmd
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	case error:
		m.err = msg
		return m, tea.Quit
	}

	var cmd tea.Cmd
	if m.state == StateScanning {
		m.spinner, cmd = m.spinner.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}
	if m.state == StateScanning {
		return fmt.Sprintf("\n %s Scanning for venvs in %s...\n\n", m.spinner.View(), m.root)
	}

	if m.state == StateDeleting {
		return fmt.Sprintf("\n Deleting venvs...\n\n%s\n", m.progress.View())
	}

	if len(m.venvs) == 0 {
		return "No venvs found.\nPress q to quit.\n"
	}

	s := "Select venvs to delete (space to select, enter to delete):\n\n"

	for i, v := range m.venvs {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		// Format size
		size := fmt.Sprintf("%.2f MB", float64(v.Size)/1024/1024)

		style := lipgloss.NewStyle()
		if m.cursor == i {
			style = style.Foreground(lipgloss.Color("205")).Bold(true)
		}
		if _, ok := m.selected[i]; ok {
			style = style.Foreground(lipgloss.Color("196")) // Red for deletion
		}

		line := fmt.Sprintf("%s [%s] %s (%s)", cursor, checked, v.Path, size)
		s += style.Render(line) + "\n"
	}

	s += "\nPress q to quit.\n"
	return s
}

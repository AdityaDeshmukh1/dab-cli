package cmd 

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices []string
	cursor int
	selected map[int]struct{}
}

func initialModel() model {
	return model {
		choices: []string{"Search Songs", "Play a song", "Download a song", "Login", "Quit"},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

// Messages
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {


		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

			case"down", "j":
			if m.cursor < len(m.choices) - 1 {
				m.cursor++
			}

		case "enter":
			choice := m.choices[m.cursor]
			fmt.Println("You chose: ", choice)

			switch choice {
			case "Search Songs":
				// call search.go here
			case "Play a song":
				// call play.go here

			case "Download a song":
				// call download.go here

			case "Login ":
				// call login.go here

			case "Quit":
				return m, tea.Quit
			}

		}

	}
	return m, nil
}

func (m model) View() string {
	s := "What do you want to do?\n\n"

	for i ,choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nPress q to quit.\n"
	return s
}

func RunTUI() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

package downloadTUI

import (
	"strconv"

	"github.com/adityadeshmukh1/dab-cli/internal/download"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	step         int // 1 = input, 2 = result display
	input        string
	message      string
	shouldReturn bool
}

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)

func NewModel() Model {
	return Model{
		step: 1,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.step {
		case 1: // Input step
			return m.handleInput(msg)
		case 2: // Result display step
			switch msg.String() {
			case "enter", "esc":
				m.shouldReturn = true
			}
		}
	}
	return m, nil
}

func (m Model) handleInput(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyRunes:
		// Only allow numeric input
		for _, r := range msg.Runes {
			if r >= '0' && r <= '9' {
				m.input += string(r)
			}
		}
	case tea.KeyBackspace:
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
	case tea.KeyEnter:
		if m.input != "" {
			// Attempt download
			trackNum, err := strconv.Atoi(m.input)
			if err != nil {
				m.message = errorStyle.Render("Invalid track number")
			} else {
				if download.Download(trackNum) {
					m.message = successStyle.Render("Track downloaded successfully!")
				} else {
					m.message = errorStyle.Render("Failed to download track")
				}
			}
			m.step = 2
		}
	case tea.KeyEsc:
		m.shouldReturn = true
	}
	return m, nil
}

func (m Model) View() string {
	switch m.step {
	case 1: // Input step
		s := titleStyle.Render("Download a Song") + "\n\n"
		s += "Enter track number to download:\n\n"
		s += m.input + "\n\n"
		s += "Press Enter to download, Esc to go back."
		return s

	case 2: // Result display step
		s := titleStyle.Render("Download Result") + "\n\n"
		s += m.message + "\n\n"
		s += "Press Enter or Esc to continue."
		return s
	}

	return ""
}

func (m Model) ShouldReturn() bool {
	return m.shouldReturn
}

func (m *Model) Reset() {
	m.step = 1
	m.input = ""
	m.message = ""
	m.shouldReturn = false
}

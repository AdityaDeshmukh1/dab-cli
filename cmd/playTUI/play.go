package playTUI

import (
	"strconv"

	"github.com/adityadeshmukh1/dab-cli/internal/play"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	step         int // 1 = input track, 2 = select quality, 3 = result display
	input        string
	quality      string
	qualityIndex int
	message      string
	shouldReturn bool
}

var (
	titleStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true)
	itemStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	successStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
	errorStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)

var qualityOptions = []string{"low", "medium", "high"}

func NewModel() Model {
	return Model{
		step:         1,
		qualityIndex: 1, // Default to "medium"
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.step {
		case 1: // Track input step
			return m.handleTrackInput(msg)
		case 2: // Quality selection step
			return m.handleQualitySelection(msg)
		case 3: // Result display step
			switch msg.String() {
			case "enter", "esc":
				m.shouldReturn = true
			}
		}
	}
	return m, nil
}

func (m Model) handleTrackInput(msg tea.KeyMsg) (Model, tea.Cmd) {
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
			m.step = 2
		}
	case tea.KeyEsc:
		m.shouldReturn = true
	}
	return m, nil
}

func (m Model) handleQualitySelection(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.qualityIndex > 0 {
			m.qualityIndex--
		}
	case "down", "j":
		if m.qualityIndex < len(qualityOptions)-1 {
			m.qualityIndex++
		}
	case "enter":
		// Attempt to play
		trackNum, err := strconv.Atoi(m.input)
		if err != nil {
			m.message = errorStyle.Render("Invalid track number")
		} else {
			m.quality = qualityOptions[m.qualityIndex]
			err := play.Play(trackNum, m.quality)
			if err != nil {
				m.message = errorStyle.Render("Error playing track: " + err.Error())
			} else {
				m.message = successStyle.Render("Playing track with " + m.quality + " quality!")
			}
		}
		m.step = 3
	case "esc":
		m.step = 1 // Go back to track input
	}
	return m, nil
}

func (m Model) View() string {
	switch m.step {
	case 1: // Track input step
		s := titleStyle.Render("Play a Song") + "\n\n"
		s += "Enter track number to play:\n\n"
		s += m.input + "\n\n"
		s += "Press Enter to continue, Esc to go back."
		return s

	case 2: // Quality selection step
		s := titleStyle.Render("Select Audio Quality") + "\n\n"
		for i, quality := range qualityOptions {
			if i == m.qualityIndex {
				s += selectedItemStyle.Render("> "+quality) + "\n"
			} else {
				s += itemStyle.Render(quality) + "\n"
			}
		}
		s += "\nUse up/down to select, Enter to play, Esc to go back."
		return s

	case 3: // Result display step
		s := titleStyle.Render("Play Result") + "\n\n"
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
	m.quality = ""
	m.qualityIndex = 1 // Reset to "medium"
	m.message = ""
	m.shouldReturn = false
}

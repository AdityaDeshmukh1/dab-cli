package searchTUI

import (
	"fmt"

	"github.com/adityadeshmukh1/dab-cli/internal/download"
	"github.com/adityadeshmukh1/dab-cli/internal/play"
	"github.com/adityadeshmukh1/dab-cli/internal/search"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	// Search state
	step         int // 0 = not started, 1 = entering query, 2 = displaying results
	query        string
	results      []search.Track
	err          string
	searching    bool
	shouldReturn bool
	cursor       int

	// Search action submenu state
	actionOpen   bool // whether submenu (Play/Download) is open
	actionCursor int  // cursor for submenu (0=Play, 1=Download)

	// UI messages
	message string

	spinner spinner.Model
}

type searchResultsMsg struct {
	tracks []search.Track
	err    error
}

var (
	titleStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true)
	itemStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

func NewModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		step:    1, // Start at query input
		spinner: s,
	}
}

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Handle spinner tick messages
	case spinner.TickMsg:
		if m.searching {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil

	// Async search results
	case searchResultsMsg:
		m.searching = false
		if msg.err != nil {
			m.err = msg.err.Error()
			m.results = nil
		} else {
			m.results = msg.tracks
			m.err = ""
		}
		return m, nil

	case tea.KeyMsg:
		// Handle input based on current step
		switch m.step {
		case 1: // Query input
			return m.handleQueryInput(msg)
		case 2: // Results display
			return m.handleResultsInput(msg)
		}
	}

	return m, nil
}

func (m Model) handleQueryInput(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyRunes:
		for _, r := range msg.Runes {
			m.query += string(r)
		}
	case tea.KeySpace:
		m.query += " "
	case tea.KeyBackspace:
		if len(m.query) > 0 {
			m.query = m.query[:len(m.query)-1]
		}
	case tea.KeyEnter:
		if m.query != "" {
			// Start search
			m.cursor = 0
			m.results = nil
			m.err = ""
			m.step = 2
			m.searching = true

			return m, tea.Batch(
				m.spinner.Tick,
				m.doSearch(m.query),
			)
		}
	case tea.KeyEsc:
		m.shouldReturn = true
	}
	return m, nil
}

func (m Model) handleResultsInput(msg tea.KeyMsg) (Model, tea.Cmd) {
	if m.searching {
		return m, nil
	}

	if m.actionOpen {
		return m.handleActionInput(msg)
	}

	// Main search result navigation
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.results)-1 {
			m.cursor++
		}
	case "enter":
		if len(m.results) > 0 {
			m.actionOpen = true
			m.actionCursor = 0
		}
	case "esc":
		m.shouldReturn = true
	}
	return m, nil
}

func (m Model) handleActionInput(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.actionCursor > 0 {
			m.actionCursor--
		}
	case "down", "j":
		if m.actionCursor < 1 {
			m.actionCursor++
		}
	case "enter":
		selectedTrack := m.results[m.cursor]
		if m.actionCursor == 0 {
			// Play
			err := play.Play(m.cursor+1, "medium")
			if err != nil {
				m.message = fmt.Sprintf("Error playing track: %s", err.Error())
			} else {
				m.message = fmt.Sprintf("Playing: %s - %s", selectedTrack.Title, selectedTrack.Artist)
			}
		} else if m.actionCursor == 1 {
			// Download
			if download.Download(m.cursor + 1) {
				m.message = fmt.Sprintf("Track %d (%s) downloaded successfully!", m.cursor+1, selectedTrack.Title)
			} else {
				m.message = fmt.Sprintf("Failed to download track %d.", m.cursor+1)
			}
		}
		m.actionOpen = false
	case "esc":
		m.actionOpen = false
	}
	return m, nil
}

func (m Model) doSearch(query string) tea.Cmd {
	return func() tea.Msg {
		tracks, err := search.Search(query)
		return searchResultsMsg{tracks: tracks, err: err}
	}
}

func (m Model) View() string {
	switch m.step {
	case 1: // Query input
		s := titleStyle.Render("Search for a track:") + "\n\n"
		s += m.query + "\n\n"
		s += "Press Enter to search, Backspace to delete, Esc to go back."
		return s

	case 2: // Results display
		s := titleStyle.Render("Search Results:") + "\n\n"

		if m.searching {
			s += fmt.Sprintf("Searching for %q %s\n", m.query, m.spinner.View())
		} else if m.err != "" {
			s += fmt.Sprintf("[ERROR] %s\n", m.err)
		} else if len(m.results) == 0 {
			s += "No tracks found.\n"
		} else {
			for i, t := range m.results {
				if m.cursor == i {
					s += selectedItemStyle.Render(fmt.Sprintf("> %2d. %s - %s", i+1, t.Title, t.Artist)) + "\n"
					if m.actionOpen {
						actions := []string{"Play", "Download"}
						for j, act := range actions {
							prefix := "   "
							if m.actionCursor == j {
								prefix = " > "
								s += selectedItemStyle.Render(fmt.Sprintf("%s%s", prefix, act)) + "\n"
							} else {
								s += itemStyle.Render(fmt.Sprintf("%s%s", prefix, act)) + "\n"
							}
						}
					}
				} else {
					s += itemStyle.Render(fmt.Sprintf("%2d. %s - %s", i+1, t.Title, t.Artist)) + "\n"
				}
			}
		}

		if m.message != "" {
			s += "\n" + m.message + "\n"
		}

		s += "\nUse up/down to navigate, Enter to select, Esc to go back."
		return s
	}

	return ""
}

func (m Model) ShouldReturn() bool {
	return m.shouldReturn
}

func (m *Model) Reset() {
	m.step = 1
	m.query = ""
	m.results = nil
	m.err = ""
	m.searching = false
	m.shouldReturn = false
	m.cursor = 0
	m.actionOpen = false
	m.actionCursor = 0
	m.message = ""
}

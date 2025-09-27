package searchScreen

import (
	"github.com/adityadeshmukh1/dab-cli/internal/search"
	"github.com/adityadeshmukh1/dab-cli/ui/searchBarTUI"
	"github.com/adityadeshmukh1/dab-cli/ui/searchResultsTUI"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type searchMsg string

type searchResultsMsg struct {
	tracks []search.Track
	err    error
}

// Screen Model
type Model struct {
	searchBar     searchBarTUI.SearchBar
	searchResults searchResultsTUI.SearchResults
	focusBar      bool
	width         int
	height        int

	// state
	shouldReturn bool
	err          string
	query        string
}

// Constructor
func New(width, height int) Model {
	return Model{
		searchBar:     searchBarTUI.NewSearchBar(),
		searchResults: searchResultsTUI.NewSearchResults(width/3, height-5),
		focusBar:      true,
		width:         width,
		height:        height,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles all messages for the screen
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle key messages
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab": // switch focus
			m.focusBar = !m.focusBar

		case "enter":
			if m.focusBar {
				// Trigger async search
				m.query = m.searchBar.Value()
				if m.query != "" {
					return m, m.doSearch(m.query)
				}
			} else {
				// If results panel has focus, select item
				if item, ok := m.searchResults.SelectedItem(); ok {
					return m, func() tea.Msg { return searchMsg(item.Title) }
				}
			}

		case "esc":
			// Pressing esc anywhere sets shouldReturn
			m.shouldReturn = true
		}

	case searchResultsMsg:
		if msg.err != nil {
			m.err = msg.err.Error()
		} else {
			results := make([]searchResultsTUI.SearchItem, 0, len(msg.tracks))
			for _, t := range msg.tracks {
				results = append(results, searchResultsTUI.SearchItem{
					Title: t.Title,
					Desc:  t.Artist,
				})
			}
			m.searchResults.SetResults(results)
		}
	}

	// Delegate updates
	if m.focusBar {
		m.searchBar, cmd = m.searchBar.Update(msg)
		if m.searchBar.ShouldReturn() {
			m.shouldReturn = true
		}
	} else {
		m.searchResults, cmd = m.searchResults.Update(msg)
	}

	return m, cmd
}

func (m Model) doSearch(query string) tea.Cmd {
	return func() tea.Msg {
		tracks, err := search.Search(query)
		return searchResultsMsg{tracks: tracks, err: err}
	}
}

// View defines how the whole screen looks
func (m Model) View() string {
	layout := lipgloss.JoinVertical(
		lipgloss.Left,
		m.searchBar.View(),
		m.searchResults.View(),
	)

	if m.err != "" {
		layout += "\n[ERROR] " + m.err
	}
	return layout
}

// Helpers so searchScreen can manage its own navigation
func (m Model) ShouldReturn() bool {
	return m.shouldReturn
}

func (m *Model) Reset() {
	m.searchBar.Reset()
	m.searchResults = searchResultsTUI.NewSearchResults(m.width/3, m.height-5)
	m.focusBar = true
	m.shouldReturn = false
	m.err = ""
	m.query = ""
}

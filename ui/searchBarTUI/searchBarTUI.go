package searchBarTUI

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SearchBar struct {
	input        textinput.Model
	style        lipgloss.Style
	title        lipgloss.Style
	shouldReturn bool
}

// Constructor
func NewSearchBar() SearchBar {
	ti := textinput.New()
	ti.Placeholder = "Type to search..."
	ti.Focus()

	return SearchBar{
		input: ti,
		style: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("12")).
			Padding(1, 2),
		title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")).
			Bold(true),
	}
}

// Init satisfies tea.Model
func (s SearchBar) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles events for the search bar
func (s SearchBar) Update(msg tea.Msg) (SearchBar, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			s.shouldReturn = true
		}
	}

	s.input, cmd = s.input.Update(msg)
	return s, cmd
}

// View renders the search bar
func (s SearchBar) View() string {
	label := s.title.Render("Search")
	box := s.style.Render(s.input.View())
	return lipgloss.JoinVertical(lipgloss.Left, label, box)
}

// Value gets the current search text
func (s SearchBar) Value() string {
	return s.input.Value()
}

// ShouldReturn allows parent to know if user pressed esc
func (s SearchBar) ShouldReturn() bool {
	return s.shouldReturn
}

// Reset clears the search bar state
func (s *SearchBar) Reset() {
	s.input.SetValue("")
	s.shouldReturn = false
}

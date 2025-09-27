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
	width        int
}

// Constructor
func NewSearchBar() SearchBar {
	ti := textinput.New()
	ti.Placeholder = "Type to search..."
	ti.Width = 40      // Set explicit width for the input
	ti.CharLimit = 100 // Optional: set character limit
	ti.Focus()

	return SearchBar{
		input: ti,
		width: 50, // Total width for the component
		style: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("12")).
			Padding(1, 2).
			Width(44), // Adjust width to accommodate padding and borders
		title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")).
			Bold(true),
	}
}

// NewSearchBarWithWidth allows custom width
func NewSearchBarWithWidth(width int) SearchBar {
	ti := textinput.New()
	ti.Placeholder = "Type to search..."
	ti.Width = width - 8 // Account for padding and borders
	ti.CharLimit = 200
	ti.Focus()

	return SearchBar{
		input: ti,
		width: width,
		style: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("12")).
			Padding(1, 2).
			Width(width - 4), // Account for borders
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
			return s, nil
		case "enter":
			// Optional: handle enter key if needed
			// You might want to trigger search here
		}
	case tea.WindowSizeMsg:
		// Responsive width adjustment
		if msg.Width > 20 {
			newWidth := min(msg.Width-10, 80) // Max width of 80, with margins
			s.input.Width = newWidth - 8
			s.style = s.style.Width(newWidth - 4)
			s.width = newWidth
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
	s.input.Blur()
	s.input.Focus() // Refocus after reset
	s.shouldReturn = false
}

// SetPlaceholder allows changing placeholder text
func (s *SearchBar) SetPlaceholder(placeholder string) {
	s.input.Placeholder = placeholder
}

// SetFocus controls input focus
func (s *SearchBar) SetFocus(focused bool) {
	if focused {
		s.input.Focus()
	} else {
		s.input.Blur()
	}
}

// IsEmpty checks if the search bar is empty
func (s SearchBar) IsEmpty() bool {
	return s.input.Value() == ""
}

// Helper function for Go versions < 1.21
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

package searchResultsTUI

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Each search result item
type SearchItem struct {
	Title string
	Desc  string
}

func (i SearchItem) TitleView() string   { return i.Title }
func (i SearchItem) Description() string { return i.Desc }
func (i SearchItem) FilterValue() string { return i.Title } // for list's filtering

// Wrapper component
type SearchResults struct {
	list  list.Model
	style lipgloss.Style
}

// Constructor
func NewSearchResults(width, height int) SearchResults {
	items := []list.Item{} // empty initially

	// custom delegate
	delegate := list.NewDefaultDelegate()

	l := list.New(items, delegate, width, height)
	l.Title = "Results"
	l.SetShowHelp(false) // hide help bar
	l.SetFilteringEnabled(false)

	return SearchResults{
		list: l,
		style: lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("10")).
			Padding(1, 2).
			Width(width).
			Height(height + 2), // account for title/border
	}
}

// Update handles key events
func (s SearchResults) Update(msg tea.Msg) (SearchResults, tea.Cmd) {
	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

// View renders the results panel
func (s SearchResults) View() string {
	return s.style.Render(s.list.View())
}

// SetResults allows replacing results dynamically
func (s *SearchResults) SetResults(items []SearchItem) {
	newItems := make([]list.Item, len(items))
	for i, it := range items {
		newItems[i] = it
	}
	s.list.SetItems(newItems)
}

// SelectedItem returns the currently selected search item
func (s SearchResults) SelectedItem() (SearchItem, bool) {
	if i, ok := s.list.SelectedItem().(SearchItem); ok {
		return i, true
	}
	return SearchItem{}, false
}

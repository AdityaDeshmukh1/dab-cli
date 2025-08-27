package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/adityadeshmukh1/dab-cli/internal/download"
	"github.com/adityadeshmukh1/dab-cli/internal/login"
	"github.com/adityadeshmukh1/dab-cli/internal/play"
	"github.com/adityadeshmukh1/dab-cli/internal/search"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFCC00")).
			MarginTop(1).
			MarginBottom(1).
			MarginLeft(2)

	itemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			PaddingLeft(4)

	selectedItemStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#00FFAA")).
				PaddingLeft(2)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#777777")).
			MarginTop(1).
			MarginLeft(2)
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}

	// TUI login state
	email     string
	password  string
	loggedIn  bool
	errMsg    string
	loginStep int // 0 = not started, 1 = email, 2 = password

	// Search Song State
	searchStep   int // 0 = not started, 1 = entering query, 2 = displaying results
	searchQuery  string
	searchResult []search.Track
	searchErr    string

	// Search action submenu state
	searchActionOpen   bool // whether submenu (Play/Download) is open
	searchActionCursor int  // cursor for submenu (0=Play, 1=Download)

	// Download state
	downloadStep    int
	downloadInput   string
	downloadMessage string

	// Play Song State
	playStep    int
	playInput   string
	playErr     string
	playQuality string
}

func initialModel() model {
	return model{
		choices:  []string{"Search Songs", "Play a song", "Download a song", "Login", "Quit"},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		// Quit
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		// -------------------
		// LOGIN HANDLER
		// -------------------
		if m.loginStep > 0 {
			switch msg.Type {
			case tea.KeyRunes:
				if m.loginStep == 1 {
					m.email += string(msg.Runes)
				} else if m.loginStep == 2 {
					m.password += string(msg.Runes)
				}
			case tea.KeyBackspace:
				if m.loginStep == 1 && len(m.email) > 0 {
					m.email = m.email[:len(m.email)-1]
				} else if m.loginStep == 2 && len(m.password) > 0 {
					m.password = m.password[:len(m.password)-1]
				}
			case tea.KeyEnter:
				if m.loginStep == 1 {
					m.loginStep = 2
				} else if m.loginStep == 2 {
					err := login.Login(m.email, m.password)
					if err != nil {
						m.errMsg = err.Error()
						m.loggedIn = false
					} else {
						m.loggedIn = true
						m.errMsg = ""
					}
					m.loginStep = 0
					m.email = ""
					m.password = ""
				}
			}
			return m, nil
		}

		// -------------------
		// SEARCH HANDLER
		// -------------------
		if m.searchStep > 0 {
			switch msg.Type {
			case tea.KeyRunes:
				if !m.searchActionOpen {
					m.searchQuery += string(msg.Runes)
				}
			case tea.KeyBackspace:
				if !m.searchActionOpen && len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				}
			case tea.KeyEnter:
				if m.searchStep == 1 {
					tracks, err := search.Search(m.searchQuery)
					if err != nil {
						m.searchErr = err.Error()
						m.searchResult = nil
					} else {
						m.searchResult = tracks
						m.searchErr = ""
					}
					m.cursor = 0
					m.searchStep = 2
				}
			}

			if m.searchStep == 2 {
				if m.searchActionOpen {
					// Submenu navigation
					switch msg.String() {
					case "up", "k":
						if m.searchActionCursor > 0 {
							m.searchActionCursor--
						}
					case "down", "j":
						if m.searchActionCursor < 1 {
							m.searchActionCursor++
						}
					case "enter":
						selectedTrack := m.searchResult[m.cursor]
						if m.searchActionCursor == 0 {
							err := play.Play(m.cursor+1, "medium")
							if err != nil {
								m.playErr = err.Error()
							}
						} else if m.searchActionCursor == 1 {
							if download.Download(m.cursor + 1) {
								m.downloadMessage = fmt.Sprintf("Track %d (%s) downloaded successfully!", m.cursor+1, selectedTrack.Title)
							} else {
								m.downloadMessage = fmt.Sprintf("Failed to download track %d.", m.cursor+1)
							}
						}
						m.searchActionOpen = false
					case "esc":
						m.searchActionOpen = false
					}
				} else {
					// Main search result navigation
					switch msg.String() {
					case "up", "k":
						if m.cursor > 0 {
							m.cursor--
						}
					case "down", "j":
						if m.cursor < len(m.searchResult)-1 {
							m.cursor++
						}
					case "enter":
						m.searchActionOpen = true
						m.searchActionCursor = 0
					case "esc":
						// back to menu
						m.searchStep = 0
						m.searchQuery = ""
						m.searchResult = nil
					}
				}
			}
			return m, nil
		}
		// -------------------
		// MAIN MENU
		// -------------------
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			switch m.choices[m.cursor] {
			case "Search Songs":
				m.searchStep = 1
				m.searchQuery = ""
				m.searchResult = nil
			case "Play a song":
				m.playStep = 1
				m.playInput = ""
				m.playErr = ""
			case "Download a song":
				m.downloadStep = 1
				m.downloadInput = ""
				m.downloadMessage = ""
			case "Login":
				m.loginStep = 1
			case "Quit":
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	// -------------------
	// LOGIN VIEW
	// -------------------
	if m.loginStep > 0 {
		s := "Login to DAB\n\n"
		s += fmt.Sprintf("Email: %s\n", m.email)
		if m.loginStep == 2 {
			s += fmt.Sprintf("Password: %s\n", strings.Repeat("*", len(m.password)))
		}
		if m.errMsg != "" {
			s += fmt.Sprintf("\n[ERROR] %s\n", m.errMsg)
		}
		s += "\nPress Enter to continue, Backspace to delete."
		return s
	}

	// -------------------
	// SEARCH VIEW
	// -------------------
	if m.searchStep == 1 {
		s := "Search for a track:\n\n"
		s += m.searchQuery
		s += "\n\nPress Enter to search, Backspace to delete."
		return s
	}
	if m.searchStep == 2 {
		s := "Search Results:\n\n"
		if m.searchErr != "" {
			s += fmt.Sprintf("[ERROR] %s\n", m.searchErr)
		} else if len(m.searchResult) == 0 {
			s += "No tracks found.\n"
		} else {
			for i, t := range m.searchResult {
				if m.cursor == i {
					s += selectedItemStyle.Render(fmt.Sprintf("> %2d. %s - %s", i+1, t.Title, t.Artist)) + "\n"
					if m.searchActionOpen {
						actions := []string{"Play", "Download"}
						for j, act := range actions {
							prefix := "   "
							if m.searchActionCursor == j {
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
		s += "\nUse up/down to navigate, Enter to select, Esc to go back."
		return s
	}

	// -------------------
	// MAIN MENU
	// -------------------
	s := titleStyle.Render("What do you want to do?") + "\n\n"
	for i, choice := range m.choices {
		if m.cursor == i {
			s += selectedItemStyle.Render(fmt.Sprintf("> %s", choice)) + "\n"
		} else {
			s += itemStyle.Render(choice) + "\n"
		}
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


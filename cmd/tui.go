package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/adityadeshmukh1/dab-cli/internal/download"
	"github.com/adityadeshmukh1/dab-cli/internal/login"
	"github.com/adityadeshmukh1/dab-cli/internal/play"
	"github.com/adityadeshmukh1/dab-cli/internal/search"
	tea "github.com/charmbracelet/bubbletea"
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
	searchStep   int // 0 = not started, 1 = entering query, displaying results
	searchQuery  string
	searchResult []search.Track
	searchErr    string

	// Download state
	downloadStep    int
	downloadInput   string
	downloadMessage string

	// Play Song State
	playStep    int
	playInput   string // accumulate digits as string
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

		// Login Handler
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

		// Search Handler
		if m.searchStep > 0 {
			switch msg.Type {
			case tea.KeyRunes:
				m.searchQuery += string(msg.Runes)
			case tea.KeyBackspace:
				if len(m.searchQuery) > 0 {
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
					m.searchStep = 2
				} else if m.searchStep == 2 {
					m.searchStep = 0
					m.searchQuery = ""
					m.searchResult = nil
					m.searchErr = ""
				}
			}
			return m, nil
		}

		// Download Handler
		if m.downloadStep > 0 {
			switch msg.Type {
			case tea.KeyRunes:
				m.downloadInput += string(msg.Runes)
			case tea.KeyBackspace:
				if len(m.downloadInput) > 0 {
					m.downloadInput = m.downloadInput[:len(m.downloadInput)-1]
				}
			case tea.KeyEnter:
				trackNum, err := strconv.Atoi(m.downloadInput)
				if err != nil || trackNum <= 0 {
					m.downloadMessage = "Invalid number!"
				} else if download.Download(trackNum) {
					m.downloadMessage = fmt.Sprintf("Track %d downloaded successfully!", trackNum)
				} else {
					m.downloadMessage = fmt.Sprintf("Failed to download track %d.", trackNum)
				}
				m.downloadStep = 2
			default:
				// ignore other keys
			}
			return m, nil
		}

		// Play Handler
		if m.playStep > 0 {
			switch msg.String() {
			case "enter":
				trackNum, err := strconv.Atoi(m.playInput)
				if err != nil || trackNum <= 0 {
					m.playErr = "Invalid track number"
				} else {
					m.playErr = ""
					err := play.Play(trackNum, "medium") // default quality
					if err != nil {
						m.playErr = err.Error()
					}
				}
				m.playStep = 0
				m.playInput = ""
			case "backspace", "backspace2":
				if len(m.playInput) > 0 {
					m.playInput = m.playInput[:len(m.playInput)-1]
				}
			default:
				if len(msg.Runes) > 0 && msg.Runes[0] >= '0' && msg.Runes[0] <= '9' {
					m.playInput += string(msg.Runes)
				}
			}
			return m, nil
		}

		// Menu navigation
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
	// Login view
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

	// Search
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
				s += fmt.Sprintf("%2d. %s - %s (ID: %d)\n", i+1, t.Title, t.Artist, t.ID)
			}
		}
		s += "\nPress any key to return to menu."
		return s
	}

	// Download
	if m.downloadStep == 1 {
		s := "Enter track number to download:\n\n"
		s += m.downloadInput
		s += "\n\nPress Enter to download, Backspace to delete."
		return s
	}
	if m.downloadStep == 2 {
		s := m.downloadMessage + "\n\nPress any key to return to menu."
		return s
	}

	// Play
	if m.playStep == 1 {
		s := "Enter track number to play:\n\n"
		s += m.playInput
		if m.playErr != "" {
			s += fmt.Sprintf("\n[ERROR] %s", m.playErr)
		}
		s += "\n\nPress Enter to play, Backspace to delete."
		return s
	}

	// Main Menu
	s := "What do you want to do?\n\n"
	for i, choice := range m.choices {
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

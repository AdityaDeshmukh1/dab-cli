package ui

import (
	"fmt"
	"os"

	"github.com/adityadeshmukh1/dab-cli/ui/downloadTUI"
	"github.com/adityadeshmukh1/dab-cli/ui/loginTUI"
	"github.com/adityadeshmukh1/dab-cli/ui/playTUI"
	"github.com/adityadeshmukh1/dab-cli/ui/searchScreen" // <-- updated

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}

	// Component models
	login    loginTUI.Model
	search   searchScreen.Model // <-- updated
	download downloadTUI.Model
	play     playTUI.Model

	// Current state
	currentView string // "menu", "login", "search", "download", "play"

	spinner spinner.Model
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		choices:     []string{"Search Songs", "Play a Song", "Download a Song", "Login", "Quit"},
		selected:    make(map[int]struct{}),
		currentView: "menu",
		login:       loginTUI.NewModel(),
		// give search screen some width/height for layout
		search:   searchScreen.New(80, 20),
		download: downloadTUI.NewModel(),
		play:     playTUI.NewModel(),
		spinner:  s,
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global quit
		switch msg.String() {
		case "ctrl+c", "q":
			if m.currentView == "menu" {
				return m, tea.Quit
			}
		}

		// Route to appropriate handler
		switch m.currentView {
		case "login":
			var cmd tea.Cmd
			m.login, cmd = m.login.Update(msg)
			if m.login.ShouldReturn() {
				m.currentView = "menu"
				m.login.Reset()
			}
			return m, cmd

		case "search":
			var cmd tea.Cmd
			m.search, cmd = m.search.Update(msg)
			if m.search.ShouldReturn() {
				m.currentView = "menu"
				m.search.Reset()
			}
			return m, cmd

		case "download":
			var cmd tea.Cmd
			m.download, cmd = m.download.Update(msg)
			if m.download.ShouldReturn() {
				m.currentView = "menu"
				m.download.Reset()
			}
			return m, cmd

		case "play":
			var cmd tea.Cmd
			m.play, cmd = m.play.Update(msg)
			if m.play.ShouldReturn() {
				m.currentView = "menu"
				m.play.Reset()
			}
			return m, cmd

		case "menu":
			return m.handleMenuInput(msg)
		}
	}

	// Handle spinner updates for active components
	var cmd tea.Cmd
	switch m.currentView {
	case "search":
		m.search, cmd = m.search.Update(msg)
	}

	return m, cmd
}

func (m model) handleMenuInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
			m.currentView = "search"
			return m, m.search.Init()
		case "Play a Song":
			m.currentView = "play"
			return m, m.play.Init()
		case "Download a Song":
			m.currentView = "download"
			return m, m.download.Init()
		case "Login":
			m.currentView = "login"
			return m, m.login.Init()
		case "Quit":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	switch m.currentView {
	case "login":
		return m.login.View()
	case "search":
		return m.search.View()
	case "download":
		return m.download.View()
	case "play":
		return m.play.View()
	default:
		return m.menuView()
	}
}

func (m model) menuView() string {
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

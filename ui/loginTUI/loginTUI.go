package loginTUI

import (
	"fmt"
	"strings"

	"github.com/adityadeshmukh1/dab-cli/internal/login"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	email        string
	password     string
	loggedIn     bool
	errMsg       string
	loginStep    int // 1 = email, 2 = password, 3 = result display
	shouldReturn bool
}

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	inputStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
)

func NewModel() Model {
	return Model{
		loginStep: 1, // Start at email input
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.loginStep {
		case 1: // Email input
			return m.handleEmailInput(msg)
		case 2: // Password input
			return m.handlePasswordInput(msg)
		case 3: // Result display
			switch msg.String() {
			case "enter", "esc":
				m.shouldReturn = true
			}
		}
	}
	return m, nil
}

func (m Model) handleEmailInput(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyRunes:
		m.email += string(msg.Runes)
	case tea.KeySpace:
		// Don't allow spaces in email
		return m, nil
	case tea.KeyBackspace:
		if len(m.email) > 0 {
			m.email = m.email[:len(m.email)-1]
		}
	case tea.KeyEnter:
		if m.email != "" {
			m.loginStep = 2
		}
	case tea.KeyEsc:
		m.shouldReturn = true
	}
	return m, nil
}

func (m Model) handlePasswordInput(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyRunes:
		m.password += string(msg.Runes)
	case tea.KeySpace:
		m.password += " "
	case tea.KeyBackspace:
		if len(m.password) > 0 {
			m.password = m.password[:len(m.password)-1]
		}
	case tea.KeyEnter:
		if m.password != "" {
			// Attempt login
			err := login.Login(m.email, m.password)
			if err != nil {
				m.errMsg = "Invalid credentials"
				m.loggedIn = false
			} else {
				m.loggedIn = true
				m.errMsg = ""
			}
			m.loginStep = 3
		}
	case tea.KeyEsc:
		m.loginStep = 1 // Go back to email input
		m.password = "" // Clear password for security
	}
	return m, nil
}

func (m Model) View() string {
	switch m.loginStep {
	case 1: // Email input
		s := titleStyle.Render("Login to DAB") + "\n\n"
		s += "Email:\n"
		s += inputStyle.Render(m.email) + "\n\n"
		s += "Press Enter to continue, Esc to go back."
		return s

	case 2: // Password input
		s := titleStyle.Render("Login to DAB") + "\n\n"
		s += fmt.Sprintf("Email: %s\n", m.email)
		s += "Password:\n"
		s += inputStyle.Render(strings.Repeat("*", len(m.password))) + "\n\n"
		s += "Press Enter to login, Esc to go back to email."
		return s

	case 3: // Result display
		s := titleStyle.Render("Login Result") + "\n\n"
		if m.loggedIn {
			s += successStyle.Render("Successfully logged in!") + "\n\n"
			s += fmt.Sprintf("Welcome, %s!", m.email) + "\n\n"
		} else {
			s += errorStyle.Render("Login failed") + "\n\n"
			if m.errMsg != "" {
				s += errorStyle.Render(m.errMsg) + "\n\n"
			}
		}
		s += "Press Enter or Esc to continue."
		return s
	}

	return ""
}

func (m Model) ShouldReturn() bool {
	return m.shouldReturn
}

func (m *Model) Reset() {
	m.loginStep = 1
	m.email = ""
	m.password = ""
	m.loggedIn = false
	m.errMsg = ""
	m.shouldReturn = false
}

// Additional methods for checking login status
func (m Model) IsLoggedIn() bool {
	return m.loggedIn
}

func (m Model) GetEmail() string {
	if m.loggedIn {
		return m.email
	}
	return ""
}

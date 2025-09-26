package loginTUI

import (
	"fmt"
	"strings"

	"github.com/adityadeshmukh1/dab-cli/internal/login"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Email     string
	Password  string
	LoggedIn  bool
	ErrMsg    string
	LoginStep int // 0 = not started, 1 = email, 2 = password
}

func New() Model {
	return Model{}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.LoginStep > 0 {
			switch msg.Type {
			case tea.KeyRunes:
				if m.LoginStep == 1 {
					m.Email += string(msg.Runes)
				} else if m.LoginStep == 2 {
					m.Password += string(msg.Runes)
				}
			case tea.KeyBackspace:
				if m.LoginStep == 1 && len(m.Email) > 0 {
					m.Email = m.Email[:len(m.Email)-1]
				} else if m.LoginStep == 2 && len(m.Password) > 0 {
					m.Password = m.Password[:len(m.Password)-1]
				}
			case tea.KeyEnter:
				if m.LoginStep == 1 {
					m.LoginStep = 2
				} else if m.LoginStep == 2 {
					// Replace with your real auth logic
					err := login.Login(m.Email, m.Password)
					if err != nil {
						m.ErrMsg = "Invalid credentials"
						m.LoggedIn = false
					} else {
						m.LoggedIn = true
						m.ErrMsg = ""
					}
					m.LoginStep = 0
					m.Email = ""
					m.Password = ""
				}
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.LoginStep == 0 {
		return ""
	}

	s := "Login to DAB\n\n"
	s += fmt.Sprintf("Email: %s\n", m.Email)
	if m.LoginStep == 2 {
		s += fmt.Sprintf("Password: %s\n", strings.Repeat("*", len(m.Password)))
	}
	if m.ErrMsg != "" {
		s += fmt.Sprintf("\n[ERROR] %s\n", m.ErrMsg)
	}
	s += "\nPress Enter to continue, Backspace to delete."
	return s
}

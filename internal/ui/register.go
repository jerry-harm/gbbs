package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"

	"gbbs/internal/database"
)

type registerModel struct {
	focuseIndex int
	inputs      []textinput.Model
	cursorMode  cursor.Mode

	focusedStyle  lipgloss.Style
	blurredStyle  lipgloss.Style
	cursorStyle   lipgloss.Style
	noStyle       lipgloss.Style
	helpStyle     lipgloss.Style
	focusedButton string
	blurredButton string
}

func initRgsiterModel(s ssh.Session) registerModel {
	renderer := bubbletea.MakeRenderer(s)

	focusedStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	blurredStyle := renderer.NewStyle().Foreground(lipgloss.Color("7"))
	noStyle := renderer.NewStyle()

	focusedButton := focusedStyle.Render("[ Submit ]")
	blurredButton := fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))

	m := registerModel{
		inputs:        make([]textinput.Model, 3),
		focusedStyle:  focusedStyle,
		blurredStyle:  blurredStyle,
		cursorStyle:   focusedStyle,
		noStyle:       noStyle,
		helpStyle:     blurredStyle,
		focusedButton: focusedButton,
		blurredButton: blurredButton,
		cursorMode:    cursor.CursorStatic,
	}
	var t textinput.Model

	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = focusedStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Nickname"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Email"
			t.CharLimit = 64
		case 2:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

func (m registerModel) Init() tea.Cmd {
	return textinput.Blink
}
func (m registerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focuseIndex == len(m.inputs) {
				err := database.Register(m.inputs[0].Value(), m.inputs[2].Value(), m.inputs[1].Value())
				if err != nil {
					return m, func() tea.Msg { return err }
				}
				return m, func() tea.Msg { return showBBS }
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focuseIndex--
			} else {
				m.focuseIndex++
			}

			if m.focuseIndex > len(m.inputs) {
				m.focuseIndex = 0
			} else if m.focuseIndex < 0 {
				m.focuseIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focuseIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = m.focusedStyle
					m.inputs[i].TextStyle = m.focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = m.noStyle
				m.inputs[i].TextStyle = m.noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m registerModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m registerModel) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := m.blurredButton
	if m.focuseIndex == len(m.inputs) {
		button = m.focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", button)

	b.WriteString(m.helpStyle.Render("register here"))

	return b.String()
}

package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
)

type errorModel struct {
	newState state
	errorMsg string

	errorStyle lipgloss.Style
	tipStyle   lipgloss.Style
}

func initErrorModel(s ssh.Session) errorModel {
	renderer := bubbletea.MakeRenderer(s)
	m := errorModel{
		errorStyle: renderer.NewStyle().Foreground(lipgloss.Color("9")).Background(lipgloss.Color("12")),
		tipStyle:   renderer.NewStyle().Foreground(lipgloss.Color("8")),
	}
	return m
}

func (m errorModel) Init() tea.Cmd {
	return nil
}

func (m errorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, func() tea.Msg { return m.newState }
	case errMsg:
		m.errorMsg = msg.err
		m.newState = msg.newState
	}
	return m, nil
}

func (m errorModel) View() string {
	var b strings.Builder
	b.WriteString(m.errorStyle.Render(m.errorMsg))
	b.WriteRune('\n')
	b.WriteString(m.tipStyle.Render("press any key to continue"))
	return b.String()
}

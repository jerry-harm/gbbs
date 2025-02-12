package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/spf13/viper"
)

type infoModel struct {
	session ssh.Session

	width    int
	height   int
	term     string
	profile  string
	bg       string
	txtStyle lipgloss.Style
	tipStyle lipgloss.Style
}

func initInfoModel(s ssh.Session) infoModel {
	pty, _, _ := s.Pty()
	renderer := bubbletea.MakeRenderer(s)
	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	tipStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))
	bg := "light"
	if renderer.HasDarkBackground() {
		bg = "dark"
	}
	m := infoModel{
		session:  s,
		term:     pty.Term,
		profile:  renderer.ColorProfile().Name(),
		width:    pty.Window.Width,
		height:   pty.Window.Height,
		bg:       bg,
		txtStyle: txtStyle,
		tipStyle: tipStyle,
	}
	return m
}

func (m infoModel) Init() tea.Cmd {
	return nil
}

func (m infoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

	case tea.KeyMsg:
		// switch msg.Type {
		// case tea.KeyEnter:
		// 	return m, func() tea.Msg { return errMsg{newState: showInfo, err: "test here"} }
		// }
		if m.session.User() == viper.GetString("ssh.newuser") {
			return m, func() tea.Msg { return showRegister }
		} else {
			return m, func() tea.Msg { return showBBS }
		}
	}
	return m, nil
}

func (m infoModel) View() string {
	s := fmt.Sprintf("Your term is %s\nYour window size is %dx%d\nBackground: %s\nColor Profile: %s\nusername: %s\nip: %s",
		m.term, m.width, m.height, m.bg, m.profile, m.session.User(), m.session.RemoteAddr())
	return m.txtStyle.Render(s) + "\n\n" + m.tipStyle.Render("Press any key to going on\n")
}

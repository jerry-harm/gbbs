package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/spf13/viper"
)

type state int

const (
	showInfo state = iota
	showRegister
	showBBS
	showBoard
	showError
)

type (
	errMsg struct {
		newState state
		err      string
	}
	chatMsg struct {
		id   string
		text string
	}
)

var progs []*tea.Program

// msg between clients
func sendToAll(msg tea.Msg) {
	for _, p := range progs {
		go p.Send(msg)
	}
}

func ProgramHandler(s ssh.Session) *tea.Program {
	m := initModel(s)
	p := tea.NewProgram(m, append(bubbletea.MakeOptions(s), tea.WithAltScreen())...)
	progs = append(progs, p)
	return p
}

type model struct {
	state   state
	session ssh.Session
	info    tea.Model
	regiser tea.Model
	board   tea.Model
	err     tea.Model
}

func initModel(s ssh.Session) tea.Model {
	m := model{
		state:   showInfo,
		session: s,
		info:    initInfoModel(s),
		board:   initBoardModel(s),
		err:     initErrorModel(s),
		regiser: initRgsiterModel(s),
	}
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case errMsg:
		m.state = showError
		m.err, cmd = m.err.Update(msg)
		return m, cmd
	case state:
		if msg != showInfo && msg != showError && msg != showRegister {
			if m.session.User() == viper.GetString("ssh.newuser") {
				return m, func() tea.Msg { return errMsg{newState: showInfo, err: "you should not be here"} }
			}
		}
		if msg == showRegister && m.session.User() != viper.GetString("ssh.newuser") {
			return m, func() tea.Msg { return errMsg{newState: showInfo, err: "you should not be here"} }
		}
		m.state = msg
	}

	switch m.state {
	case showInfo:
		m.info, cmd = m.info.Update(msg)
	case showRegister:
		m.regiser, cmd = m.regiser.Update(msg)
	case showBoard:
		m.board, cmd = m.board.Update(msg)
	case showError:
		m.err, cmd = m.err.Update(msg)
	default:
		m.state = showInfo
	}
	if cmd == nil {
		return m, cmd
	}

	return m, cmd
}

func (m model) View() string {
	switch m.state {
	case showInfo:
		return m.info.View()
	case showRegister:
		return m.regiser.View()
	case showBoard:
		return m.board.View()
	case showError:
		return m.err.View()
	default:
		return m.info.View()
	}
}

package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"

	"gbbs/internal/database"
)

type boardModel struct {
	boardList []database.Board
	page      int
	input     textinput.Model
	session   ssh.Session

	messageStyle  lipgloss.Style
	userNameStyle lipgloss.Style
	dateStyle     lipgloss.Style
	titleStyle    lipgloss.Style
	blurredStyle  lipgloss.Style
	helpStyle     lipgloss.Style
}

func initBoardModel(s ssh.Session) boardModel {
	renderer := bubbletea.MakeRenderer(s)
	ti := textinput.New()
	ti.Placeholder = "hello to the world"
	ti.Focus()
	ti.CharLimit = 60
	ti.Width = 20

	m := boardModel{
		page:    0,
		input:   ti,
		session: s,

		messageStyle:  renderer.NewStyle().Foreground(lipgloss.Color("15")),
		userNameStyle: renderer.NewStyle().Bold(true).Foreground(lipgloss.Color("5")),
		dateStyle:     renderer.NewStyle().Foreground(lipgloss.Color("4")),
		titleStyle:    renderer.NewStyle().Foreground(lipgloss.Color("1")).Background(lipgloss.Color("2")),
		helpStyle:     renderer.NewStyle().Foreground(lipgloss.Color("7")),
	}
	database.DB.Order("created_at desc").Preload("User").Limit(10).Find(&m.boardList)
	return m
}

func (m boardModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m boardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	// up down change page
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// new board message
			var user database.User
			database.DB.Where("name = ?", m.session.User()).First(&user)
			user.Board(m.input.Value())
			return m, func() tea.Msg { return showBBS }
		case "esc":
			// not post
			return m, func() tea.Msg { return showBBS }
		case "up":
			m.page++
			database.DB.Order("created_at desc").Preload("User").Limit(10).Offset(m.page * 10).Find(&m.boardList)
			if len(m.boardList) == 0 {
				m.page--
				database.DB.Order("created_at desc").Preload("User").Limit(10).Offset(m.page * 10).Find(&m.boardList)
			}
		case "down":
			if m.page > 0 {
				m.page--
			}
			database.DB.Order("created_at desc").Preload("User").Limit(10).Offset(m.page * 10).Find(&m.boardList)
		}

	}
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m boardModel) View() string {
	var b strings.Builder
	b.WriteString(m.titleStyle.Render("Board"))
	b.WriteString("  ")
	b.WriteString(m.helpStyle.Render(fmt.Sprintf("page %d", m.page)))
	b.WriteRune('\n')
	for _, item := range m.boardList {
		b.WriteString(m.userNameStyle.Render(item.User.Name))
		b.WriteRune(' ')
		b.WriteString(m.dateStyle.Render(item.CreatedAt.Format("2006-01-02 15:04:05")))
		b.WriteRune('\n')
		b.WriteString(m.messageStyle.Render(item.Text))
		b.WriteRune('\n')
		b.WriteRune('\n')
	}

	b.WriteString(m.input.View())
	b.WriteRune('\n')
	b.WriteString(m.helpStyle.Render("press enter to post,press esc to not post,press up and down to view"))
	return b.String()
}

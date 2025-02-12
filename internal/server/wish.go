package server

import (
	"net"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/muesli/termenv"
	"github.com/spf13/viper"

	"gbbs/internal/database"
	"gbbs/internal/ui"
)

func GetServer() *ssh.Server {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(viper.GetString("host"), viper.GetString("ssh.port"))),
		wish.WithHostKeyPath(viper.GetString("ssh.keypath")),

		wish.WithPasswordAuth(func(ctx ssh.Context, password string) bool {
			if ctx.User() == viper.GetString("ssh.newuser") {
				return true // TODO regiser and invite password maybe?
			}
			return database.Login(ctx.User(), password)
		}),

		wish.WithMiddleware(
			bubbletea.MiddlewareWithProgramHandler(ui.ProgramHandler, termenv.ANSI256),
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}
	return s

}

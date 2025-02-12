package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/spf13/viper"

	"gbbs/config"
	"gbbs/internal/database"
	"gbbs/internal/server"
)

func main() {
	config.InitConfig()
	log.Info("Start!")

	viper.WatchConfig()
	database.InitDB()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	s := server.GetServer()
	log.Info("Starting SSH server", "host", viper.GetString("host"), "port", viper.GetString("ssh.port"))
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			quit <- nil
		}
	}()

	<-quit
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}

	log.Info("Shutdown...")
}

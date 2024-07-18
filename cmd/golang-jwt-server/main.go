package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatalistix/slogattr"

	"github.com/fatalistix/golang-jwt-server/internal/app"
	"github.com/fatalistix/golang-jwt-server/internal/config"
	"github.com/fatalistix/golang-jwt-server/internal/env"
	"github.com/fatalistix/golang-jwt-server/internal/log"
)

func main() {
	env.MustLoad()

	configPath := os.Getenv("CONFIG_PATH")

	config := config.MustLoad(configPath)

	log := log.MustSetup(config.Env)

	log.Debug("loaded config", slog.Any("config", config))

	application, err := app.New(log, config)
	if err != nil {
		log.Error("cannot create app", slogattr.Err(err))
		os.Exit(1)
	}

	go func() {
		if err := application.Run(); err != nil {
			log.Error("server stopped with error", slogattr.Err(err))
		}
	}()

	log.Info("application started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	if err := application.Stop(ctx); err != nil {
		log.Error("application stopped with error", slogattr.Err(err))

		os.Exit(1)
	}

	log.Info("application stopped")
}

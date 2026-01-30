package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/bekzat-kamen/startGO.git/internal/config"
	"github.com/bekzat-kamen/startGO.git/internal/pkg/logger"
	"github.com/bekzat-kamen/startGO.git/internal/server"
	"github.com/gin-gonic/gin"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %s", err.Error())
	}

	slogger := logger.New(cfg.LogLevel)
	slog.SetDefault(slogger)

	r := gin.New()

	srv := server.New(r, cfg.Port)
	err = srv.Run()
	if err != nil {
		slog.Error("failed to start server", "error", err.Error())
		os.Exit(1)
	}
}

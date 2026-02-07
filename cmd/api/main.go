package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/bekzat-kamen/startGO.git/internal/config"
	"github.com/bekzat-kamen/startGO.git/internal/handler"
	"github.com/bekzat-kamen/startGO.git/internal/pkg/logger"
	"github.com/bekzat-kamen/startGO.git/internal/repository"
	"github.com/bekzat-kamen/startGO.git/internal/server"
	"github.com/bekzat-kamen/startGO.git/internal/service"
	_ "github.com/gin-gonic/gin"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %s", err.Error())
	}

	slogger := logger.New(cfg.LogLevel)
	slog.SetDefault(slogger)

	// 	r := gin.New()
	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err.Error())
		return
	}
	
	courseRepo := repository.NewPsgCourseRepo(db)
	courseService := service.NewCourseService(courseRepo)
	h := handler.NewHandler(courseService)
	router, err := h.InitRoutes()

	srv := server.New(router, cfg.Port)
	err = srv.Run()
	if err != nil {
		slog.Error("failed to start server", "error", err.Error())
		os.Exit(1)
	}
}

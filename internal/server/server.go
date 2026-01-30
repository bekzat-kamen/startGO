package server

import (
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	myServer *http.Server
}

func New(router http.Handler, port string) *Server {
	return &Server{
		myServer: &http.Server{
			Addr:           ":" + port,
			Handler:        router,
			MaxHeaderBytes: http.DefaultMaxHeaderBytes,
			ReadTimeout:    40 * time.Second,
			WriteTimeout:   40 * time.Second,
		},
	}
}

func (s *Server) Run() error {
	slog.Info("server is started", "ADDR", s.myServer.Addr)
	return s.myServer.ListenAndServe()
}

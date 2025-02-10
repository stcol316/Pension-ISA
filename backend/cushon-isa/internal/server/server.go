package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/stcol316/cushon-isa/internal/database"
)

type Server struct {
	port int
	db   database.Service
}

func NewServer(listenAddr int, db database.Service) *http.Server {
	NewServer := &Server{
		port: 8080, //TODO: Use env var here
		db:   db,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

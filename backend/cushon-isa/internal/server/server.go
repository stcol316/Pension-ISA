package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/stcol316/cushon-isa/internal/customer"
	"github.com/stcol316/cushon-isa/internal/fund"
	"github.com/stcol316/cushon-isa/internal/investment"
)

type Server struct {
	port              int
	customerHandler   *customer.Handler
	fundHandler       *fund.Handler
	investmentHandler *investment.Handler
}

func NewServer(listenAddr int, ch *customer.Handler, fh *fund.Handler, ih *investment.Handler) *http.Server {
	NewServer := &Server{
		port:              8080, //TODO: Use env var here
		customerHandler:   ch,
		fundHandler:       fh,
		investmentHandler: ih,
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

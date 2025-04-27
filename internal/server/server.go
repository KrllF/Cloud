package server

import (
	"net"
	"net/http"

	"github.com/KrllF/Cloud/internal/config"
)

type (
	// Handler интерфейс хендлера
	Handler interface {
		Init() http.HandlerFunc
	}
	// HTTP структура сервера
	HTTP struct {
		httpServer *http.Server
	}
)

// NewServer конструктор сервера
func NewServer(conf config.Config, handler http.Handler) *HTTP {
	return &HTTP{&http.Server{
		Addr:              net.JoinHostPort(conf.HTTP_HOST, conf.HTTP_PORT),
		Handler:           handler,
		ReadTimeout:       conf.Read,
		WriteTimeout:      conf.Write,
		IdleTimeout:       conf.Idle,
		ReadHeaderTimeout: conf.ReadHeader,
	}}
}

// Run запустить HTTP сервер
func (s *HTTP) Run() error {
	return s.httpServer.ListenAndServe()
}

// Close закрыть HTTP сервер
func (s *HTTP) Close() {
	s.httpServer.Close()
}

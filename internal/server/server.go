package server

import (
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(bindAddr string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:         bindAddr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

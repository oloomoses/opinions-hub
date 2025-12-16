package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func New(addr string, handler http.Handler) *Server {
	svr := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{httpServer: svr}
}

func (s *Server) Start() {
	go func() {
		log.Printf("server listening on %s", s.httpServer.Addr)

		if err := s.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("listen error: %v", err)
		}
	}()
	// return s.httpServer.ListenAndServe()
}

func (s *Server) ShutDown(ctx context.Context) error {
	log.Println("shutting down server ....")
	return s.httpServer.Shutdown(ctx)
}

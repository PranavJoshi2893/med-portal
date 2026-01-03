package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PranavJoshi2893/boilerplate/internal/config"
)

type Server struct {
	httpServer *http.Server
	db         *sql.DB
}

func NewServer(cfg *config.Config, db *sql.DB, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:              cfg.ServerPort,
			Handler:           handler,
			ReadHeaderTimeout: time.Second * 5,
			ReadTimeout:       time.Second * 15,
			WriteTimeout:      time.Second * 30,
			IdleTimeout:       time.Second * 60,
		},
		db: db,
	}
}

func (s *Server) Run() error {
	errChan := make(chan error, 1)

	go func() {
		err := s.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return fmt.Errorf("failed to start server: %v", err)
	case <-quit:
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("forced shutdown: %v", err)
		}

		return nil
	}

}

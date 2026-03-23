package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/lstratta/flatpeak-take-home-task/config"
)

type Server struct {
  srv *http.Server
	conf config.Config 
}

func NewServer(c config.Config) *Server {
	m := NewServeMux()
	m.AddRoutes()

	return &Server{
		srv: &http.Server{
			Addr: fmt.Sprintf("%s:%s", c.Host, c.Port),
			Handler: m,
		},
		conf: c,
	}
}

func (s *Server) Start() error {
	serveErr := make(chan error, 1)
	defer close(serveErr)
	go func() {
	  log.Println("Starting server...")
		if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select{
	case err := <-serveErr:
		return err
	case <-stop:
		log.Println("Shutdown signal received")
	}
 
	if err := s.srv.Close(); err != nil {
		return err
	}
	log.Println("Server closed")
 	return nil
}

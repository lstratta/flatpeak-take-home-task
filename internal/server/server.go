package server

import (
	"fmt"
	"net/http"

	"github.com/lstratta/flatpeak-take-home-task/config"
)

type Server struct {
  srv *http.Server
	conf config.Config 
}

func NewServer(c config.Config, mux *serveMux) *Server {
	return &Server{
		srv: &http.Server{
			Addr: fmt.Sprintf("%s:%s", c.Host, c.Port),
			Handler: mux,
		},
		conf: c,
	}
}

func (s *Server) Start() error {

	return s.srv.ListenAndServe()
}

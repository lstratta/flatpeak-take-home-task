package main

import (
	"log"
	"net/http"

	"github.com/lstratta/flatpeak-take-home-task/config"
	"github.com/lstratta/flatpeak-take-home-task/internal/server"
)

func main() {
	c := config.New()
	s := server.NewServer(c)
	
	if err := s.Start(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

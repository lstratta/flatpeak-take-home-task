package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/lstratta/flatpeak-take-home-task/internal/neso"
)

func (s *serveMux) testHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d, err := neso.GetNesoData()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		b, err := json.Marshal(d)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
		if err != nil {
			log.Println("failed to write body")
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}

func (s *serveMux) healthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("service alive"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

	})
}

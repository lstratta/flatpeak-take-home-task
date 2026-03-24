package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/lstratta/flatpeak-take-home-task/internal/calculate"
	"github.com/lstratta/flatpeak-take-home-task/internal/neso"
)

func (s *serveMux) slotsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d, err := neso.GetNesoData()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		var dur time.Duration = 30

		lowSlots, err := calculate.FilterPeriodsByLowestIntensity(d.Data, dur)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}

		b, err := json.Marshal(lowSlots)
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

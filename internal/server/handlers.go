package server

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/lstratta/flatpeak-take-home-task/internal/calculate"
	"github.com/lstratta/flatpeak-take-home-task/internal/models"
	"github.com/lstratta/flatpeak-take-home-task/internal/neso"
)

func (s *serveMux) slotsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.RawPath)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		q := u.Query()
		durParam := q["duration"][0]
		isContinuousParam := q["continuous"][0]

		dur, err := strconv.Atoi(durParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		isContinuous, err := strconv.ParseBool(isContinuousParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		d, err := neso.GetNesoData()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		duration := time.Duration(dur)
		pArr, err := calculate.FilterPeriodsByLowestIntensity(d.Data, duration)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}

		var slot *models.Slot

		if isContinuous {
			fArr, err := calculate.FilterPeriodsByDuration(pArr, duration)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			slot, err = calculate.CalculateContinuousPeriodIntensity(fArr, duration)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		b, err := json.Marshal(&slot)
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

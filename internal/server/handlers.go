package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/lstratta/flatpeak-take-home-task/internal/calculate"
	"github.com/lstratta/flatpeak-take-home-task/internal/models"
	"github.com/lstratta/flatpeak-take-home-task/internal/neso"
)

func (s *serveMux) slotsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		durParam := q["duration"]
		isContinuousParam := q["continuous"]

		// if durParam is empty, use 30 default value
		if len(durParam) < 1 || len(isContinuousParam) < 1 {
			durParam = []string{"30"}
		}

		// if isContinuousParam is empty, use false default value
		if len(isContinuousParam) < 1 {
			isContinuousParam = []string{"false"}
		}

		dur := durParam[0]
		isContinuousString := isContinuousParam[0]

		isContinuous, err := strconv.ParseBool(isContinuousString)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("error converting continuous url param to bool")
			return
		}

		d, err := neso.GetNesoData()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		durInMinutes := dur + "m"
		duration, err := time.ParseDuration(durInMinutes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("error parsing duration")
			return
		}

		pArr, err := calculate.FilterPeriodsByLowestIntensity(d.Data, duration)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		var slots []models.Slot

		if isContinuous {
			slot, err := calculate.CalculateContinuousPeriodIntensity(pArr, duration)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("error calculating continuous period by duration: %v", err)
				return
			}
			slots = append(slots, *slot)
		} else {
			slots, err = calculate.CalculateWeightedAverage(pArr, duration)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("error calculating weighted average: %v", err)
				return
			}

		}

		b, err := json.Marshal(&slots)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error marshalling object to json: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
		if err != nil {
			log.Printf("failed to write body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *serveMux) healthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("service alive"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to write body: %v", err)
			return
		}

	})
}

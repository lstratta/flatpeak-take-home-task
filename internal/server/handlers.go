package server

import (
	"encoding/json"
	"fmt"
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
		if len(durParam) < 1 {
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
			log.Println("error converting continuous url param to bool")
			return
		}

		d, err := neso.GetNesoData()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("error getting data", err)
			return
		}

		durInMinutes := dur + "m"
		duration, err := time.ParseDuration(durInMinutes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("error parsing duration: ", err)
			return
		}

		if duration.Minutes() < 0 || duration.Minutes() > 1440 {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("duration out of range: ", duration.Minutes())
			return
		}

		slots, err := calculations(d, duration, isContinuous)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("error handling logic: ", err)
			return
		}

		data := struct {
			Data []models.Slot `json:"data"`
		}{
			Data: slots,
		}

		b, err := json.Marshal(&data)
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

func calculations(d *models.Data, duration time.Duration, isContinuous bool) ([]models.Slot, error) {
	c := calculate.NewCalculationService()
	slots := []models.Slot{}

	pArr, err := c.FilterPeriodsByLowestIntensity(d.Data, duration)
	if err != nil {
		return nil, fmt.Errorf("error calculating lowest intesity: %v", err)
	}

	// if isContinuous == true, find the average for all periods and return as one period
	// else, return all the number of lowest periods over the next 24 hours that fit
	// within the given time duration
	if isContinuous {
		slot, err := c.CalculateContinuousPeriodIntensity(pArr, duration)
		if err != nil {
			return nil, fmt.Errorf("error calculating continuous period by duration: %v", err)
		}
		slots = append(slots, *slot)
	} else {
		slots, err = c.CalculateWeightedAverage(pArr, duration)
		if err != nil {
			return nil, fmt.Errorf("error calculating weighted average: %v", err)
		}
	}

	return slots, nil
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

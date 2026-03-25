package server

import (
	"encoding/json"
	"fmt"
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
		q := r.URL.Query()

		// validate and convert query params
		durStr, isContinuousStr := validateSlotsQueryParams(q)

		isContinuous, err := strconv.ParseBool(isContinuousStr)
		if err != nil {
			log.Println("error converting continuous url param to bool")
			return
		}

		durInMinutes := durStr + "m"
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

		// fetch data from NESO
		d, err := neso.GetNesoData()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("error getting data", err)
			return
		}

		// execute calculations
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

func validateSlotsQueryParams(q url.Values) (dur string, isContinuous string) {
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

	dur = durParam[0]
	isContinuous = isContinuousParam[0]
	return dur, isContinuous
}

func calculations(d *models.Data, duration time.Duration, isContinuous bool) ([]models.Slot, error) {
	c := calculate.NewCalculationService()
	slots := []models.Slot{}

	data := d.Data

	// if isContinuous == true, find the average for all periods and return as one period
	// else, return all the number of lowest periods over the next 24 hours that fit
	// within the given time duration
	if isContinuous {
		p, err := c.FilterPeriodsByDuration(data, duration)
		if err != nil {
			return nil, fmt.Errorf("error filtering periods by duration: %v", err)
		}

		slot, err := c.CalculateContinuousPeriodIntensity(p, duration)
		if err != nil {
			return nil, fmt.Errorf("error calculating continuous period by duration: %v", err)
		}

		slots = append(slots, *slot)

	} else {
		pArr, err := c.FilterPeriodsByLowestIntensity(data, duration)
		if err != nil {
			return nil, fmt.Errorf("error calculating lowest intesity: %v", err)
		}
		s, err := c.CalculateWeightedAverage(pArr, duration)
		if err != nil {
			return nil, fmt.Errorf("error calculating weighted average: %v", err)
		}
		slots = append(slots, s...)
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

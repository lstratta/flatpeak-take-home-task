package calculate

import (
	"fmt"
	"time"

	"github.com/lstratta/flatpeak-take-home-task/internal/neso"
)

type Slot struct {
	ValidFrom time.Time `json:"valid_from"`
	ValidTo   time.Time `json:"valid_to"`
	Carbon    Carbon    `json:"carbon"`
}

type Carbon struct {
	Intensity int64 `json:"intensity"`
}

func FilterPeriodsByLowIntensity(p []neso.Period) ([]neso.Period, error) {
	var lowPeriods []neso.Period

	for i := range p {
		idx := p[i]
		if idx.Intensity.Index != "low" {
			continue
		}

		lowPeriods = append(lowPeriods, idx)
	}

	return lowPeriods, nil
}

func FilterPeriodsByDuration(p []neso.Period, duration time.Duration) ([]neso.Period, error) {
	var selectedPeriods []neso.Period
	startTime, err := formatTime(p[0].From)
	if err != nil {
		return nil, fmt.Errorf("error formatting time: %v", err)
	}

	for i := range p {
		idx := p[i]
		endTime, err := formatTime(p[i].To)
		if err != nil {
			return nil, fmt.Errorf("error formatting time: %v", err)
		}

		diff := endTime.Sub(startTime)
		if diff <= duration {
			selectedPeriods = append(selectedPeriods, idx)
		}
	}

	return selectedPeriods, nil
}

func formatTime(s string) (time.Time, error) {
	// remove end Z character
	sCut := s[:len(s)-1]
	// append :00Z
	formattedTime := fmt.Sprint(sCut, ":00Z")

	t, err := time.Parse(time.RFC3339, formattedTime)
	if err != nil {
		return time.Time{}, fmt.Errorf("from field: error converting string to time: %v", err)
	}

	return t, nil
}

// func hold() {
//
// 	from, err := formatTime(idx.From)
// 	if err != nil {
// 		return nil, fmt.Errorf("error formatting from time: %v", err)
// 	}
// 	to, err := formatTime(idx.To)
// 	if err != nil {
// 		return nil, fmt.Errorf("error formatting to time: %v", err)
// 	}
// 	acceptedSlot := slot{
// 		ValidFrom: from,
// 		ValidTo:   to,
// 		Carbon: carbon{
// 			Intensity: idx.Intensity.Forecast,
// 		},
// 	}
//
// 	lowSlots = append(lowSlots, acceptedSlot)
//
// }

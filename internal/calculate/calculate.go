package calculate

import (
	"fmt"
	"time"

	"github.com/lstratta/flatpeak-take-home-task/internal/neso"
)

type slot struct {
	ValidFrom time.Time `json:"valid_from"`
	ValidTo   time.Time `json:"valid_to"`
	Carbon    carbon    `json:"carbon"`
}

type carbon struct {
	Intensity int64 `json:"intensity"`
}

func FilterLowIntensityPeriods(p []neso.Period) ([]neso.Period, error) {
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

func FilterByDuration(p []neso.Period, duration int64, isContinuous bool) ([]neso.Period, error)

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

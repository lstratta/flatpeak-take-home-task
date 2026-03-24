package calculate

import (
	"fmt"
	"log"
	"time"

	"github.com/lstratta/flatpeak-take-home-task/internal/neso"
)

type slots struct {
	ValidFrom time.Time `json:"valid_from"`
	ValidTo   time.Time `json:"valid_to"`
	Carbon    carbon    `json:"carbon"`
}

type carbon struct {
	Intensity int64 `json:"intensity"`
}

func FilterLowIntensitySlots(s *neso.Data) ([]slots, error) {
	var lowSlots []slots

	for i := range s.Data {
		idx := s.Data[i]
		if idx.Intensity.Index != "low" {
			continue
		}

		from, err := formatTime(idx.From)
		if err != nil {
			return nil, fmt.Errorf("error formatting from time: %v", err)
		}
		to, err := formatTime(idx.To)
		if err != nil {
			return nil, fmt.Errorf("error formatting to time: %v", err)
		}
		acceptedSlot := slots{
			ValidFrom: from,
			ValidTo:   to,
			Carbon: carbon{
				Intensity: idx.Intensity.Forecast,
			},
		}

		lowSlots = append(lowSlots, acceptedSlot)
	}

	log.Printf("length of lowSlots: %d", len(lowSlots))

	return lowSlots, nil
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

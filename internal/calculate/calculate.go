package calculate

import (
	"fmt"
	"math"
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

type index string

const (
	veryLowIndex    index         = "very low"
	lowIndex        index         = "low"
	moderateIndex   index         = "moderate"
	highIndex       index         = "high"
	fixedTimePeriod time.Duration = 30
)

func FilterPeriodsByLowestIntensity(p []neso.Period, duration time.Duration) ([]neso.Period, error) {
	var lowPeriods []neso.Period

	for i := range p {
		idx := p[i]

		lowPeriods = append(lowPeriods, idx)
	}

	return lowPeriods, nil
}

func FilterPeriodsByDuration(pArr []neso.Period, duration time.Duration) ([]neso.Period, error) {
	var selectedPeriods []neso.Period

	startTime, err := formatTime(pArr[0].From)
	if err != nil {
		return nil, fmt.Errorf("error formatting time: %v", err)
	}

	for i := range pArr {
		idx := pArr[i]
		endTime, err := formatTime(pArr[i].To)
		if err != nil {
			return nil, fmt.Errorf("error formatting time: %v", err)
		}

		diff := endTime.Sub(startTime)
		if diff <= duration {
			selectedPeriods = append(selectedPeriods, idx)
		}
	}

	// Capture the period any duration over 30 mins overflows into
	timeRemainder := int(duration.Minutes()) % int(fixedTimePeriod)
	if timeRemainder > 0 {
		l := len(selectedPeriods)
		selectedPeriods = append(selectedPeriods, pArr[l])
	}
	return selectedPeriods, nil
}

func CalculateWeightedAverageForTimePeriod() {

}

func CalculateContinuousPeriod(pArr []neso.Period, duration time.Duration) (int64, error) {
	weight := 0.0
	totalIntensity := 0.0
	l := len(pArr)

	timeRemainder := int(duration.Minutes()) % int(fixedTimePeriod)
	if timeRemainder > 0 {
		weight = float64(timeRemainder) / float64(fixedTimePeriod)
		totalIntensity = float64(pArr[l-1].Intensity.Forecast) * weight
		for i := range l {
			// skip the last element of the slice
			if i == l-1 {
				continue
			}
			totalIntensity += float64(pArr[i].Intensity.Forecast)
		}
	} else {
		weight = 1.0
		for _, p := range pArr {
			totalIntensity += float64(p.Intensity.Forecast)
		}
	}

	averageIntensity := totalIntensity / (float64(l) - 1 + weight)

	return int64(math.Round(averageIntensity)), nil
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

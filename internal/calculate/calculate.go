package calculate

import (
	"fmt"
	"math"
	"sort"
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

func FilterPeriodsByLowestIntensity(pArr []neso.Period, duration time.Duration) ([]neso.Period, error) {
	var lowPeriods []neso.Period
	timeSpan := int(math.Ceil(float64(duration.Minutes()) / float64(fixedTimePeriod)))
	indexes := []index{veryLowIndex, lowIndex, moderateIndex, highIndex}

	fmt.Println("timeSpan: ", timeSpan)

	k := 0
	for _, idx := range indexes {
		for _, p := range pArr {

			if idx == index(p.Intensity.Index) {
				lowPeriods = append(lowPeriods, p)
				k++
			}
			if k == timeSpan {
				break
			}
		}
		if k == timeSpan {
			break
		}
	}

	fmt.Println("before sort", lowPeriods)
	sort.Sort(neso.ByDateSorter(lowPeriods))
	fmt.Println("after sort", lowPeriods)
	return lowPeriods, nil
}

func FilterPeriodsByDuration(pArr []neso.Period, duration time.Duration) ([]neso.Period, error) {
	var selectedPeriods []neso.Period

	startTime := pArr[0].From

	for i := range pArr {
		idx := pArr[i]
		endTime := pArr[i].To

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

func CalculateWeightedAverageForTimePeriod(pArr []neso.Period, duration time.Duration) ([]Slot, error) {
	s := []Slot{}

	for _, p := range pArr {

		entry := Slot{
			ValidFrom: p.From,
			ValidTo:   p.To,
			Carbon: Carbon{
				Intensity: p.Intensity.Forecast,
			},
		}
		s = append(s, entry)
	}

	fixedTimePeriodInt64 := int64(fixedTimePeriod)
	durationInt64 := int64(duration.Minutes())

	timeRemainder := durationInt64 % fixedTimePeriodInt64
	if timeRemainder == 0 {
		return s, nil
	}

	weight := timeRemainder / fixedTimePeriodInt64
	l := len(s)
	if l < 1 {
		return nil, fmt.Errorf("slice length too short")
	}

	s[l-1].Carbon.Intensity = s[l-1].Carbon.Intensity * weight

	return s, nil
}

func CalculateContinuousPeriodIntensity(pArr []neso.Period, duration time.Duration) (int64, error) {
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

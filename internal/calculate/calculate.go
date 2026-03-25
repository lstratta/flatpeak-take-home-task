package calculate

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/lstratta/flatpeak-take-home-task/internal/models"
)

const (
	fixedTimePeriod time.Duration = 30
)

func FilterPeriodsByLowestIntensity(pArr []models.Period, duration time.Duration) ([]models.Period, error) {
	var lowPeriods []models.Period
	count := len(pArr)
	for count > 0 {
		minItensityIndex := 0
		for i := 1; i < count; i++ {
			if pArr[i].Intensity.Forecast < pArr[minItensityIndex].Intensity.Forecast {
				minItensityIndex = i
			}
		}

		lowPeriods = append(lowPeriods, pArr[minItensityIndex])
		count--
	}

	sort.Sort(models.ByDateSorter(lowPeriods))
	return lowPeriods, nil
}

func FilterPeriodsByDuration(pArr []models.Period, duration time.Duration) ([]models.Period, error) {
	var selectedPeriods []models.Period

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

func CalculateWeightedAverageForTimePeriodByDuration(pArr []models.Period, duration time.Duration) ([]models.Slot, error) {
	s := []models.Slot{}

	for _, p := range pArr {

		entry := models.Slot{
			ValidFrom: p.From,
			ValidTo:   p.To,
			Carbon: models.Carbon{
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

	weight := float64(timeRemainder) / float64(fixedTimePeriodInt64)
	l := len(s)
	if l < 1 {
		return nil, fmt.Errorf("slice length too short")
	}

	i := float64(s[l-1].Carbon.Intensity) * weight

	s[l-1].Carbon.Intensity = int64(math.Round(i))
	return s, nil
}

func CalculateContinuousPeriodIntensity(pArr []models.Period, duration time.Duration) (*models.Slot, error) {
	weight := 0.0
	totalIntensity := 0.0

	l := len(pArr)
	if l < 1 {
		return nil, fmt.Errorf("array length 0")
	}
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

	intensity := int64(math.Round(averageIntensity))

	s := &models.Slot{
		ValidFrom: pArr[0].From,
		ValidTo:   pArr[len(pArr)-1].To,
		Carbon: models.Carbon{
			Intensity: intensity,
		},
	}

	return s, nil
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

package calculate

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/lstratta/flatpeak-take-home-task/internal/models"
)

const (
	// The length of the time periods provided by NESO (30 minute intervals)
	fixedTimePeriod time.Duration = 1_800_000_000_000 // nanoseconds
)

type CalculationService struct {
}

func NewCalculationService() *CalculationService {
	return &CalculationService{}
}

// FilterPeriodsByLowestIntensity accepts a slice of models.Period and a time.Duration variable, and returns
// []models.Period, or an error. It will filter the lowest intesity time periods.
func (c *CalculationService) FilterPeriodsByLowestIntensity(pArr []models.Period, dur time.Duration) ([]models.Period, error) {
	var lowPeriods []models.Period

	timeSpan := int64(math.Ceil(dur.Minutes() / fixedTimePeriod.Minutes()))
	timeSpanCount := int64(0)

	count := len(pArr)
	for count > 0 {
		minIntensityIndex := 0
		for i := 1; i < count; i++ {
			if pArr[i].Intensity.Forecast < pArr[minIntensityIndex].Intensity.Forecast {
				minIntensityIndex = i
			}
		}

		lowPeriods = append(lowPeriods, pArr[minIntensityIndex])
		pArr = append(pArr[:minIntensityIndex], pArr[minIntensityIndex+1:]...)
		count--
		timeSpanCount++
		if timeSpanCount >= timeSpan {
			break
		}
	}

	sort.Sort(models.ByDateSorter(lowPeriods))
	return lowPeriods, nil
}

// FilterPeriodsByDuration accepts a slice of models.Period and a time.Duration variable, and returns
// []models.Period, or an error. It will filter the number of periods that fit within the specified
// duration.
func (c *CalculationService) FilterPeriodsByDuration(pArr []models.Period, duration time.Duration) ([]models.Period, error) {
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
	timeRemainder := int(duration.Minutes()) % int(fixedTimePeriod.Minutes())
	if timeRemainder > 0 {
		l := len(selectedPeriods)
		selectedPeriods = append(selectedPeriods, pArr[l])
	}

	return selectedPeriods, nil
}

// CalculateWeightedAverage accepts a slice of models.Period and a time.Duration variable, and returns
// []models.Slot, or an error. It will calculate the last weighted average of the last element in the slice.
func (c *CalculationService) CalculateWeightedAverage(pArr []models.Period, duration time.Duration) ([]models.Slot, error) {
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

	fixedTimePeriodInt64 := int64(fixedTimePeriod.Minutes())
	durationInt64 := int64(duration.Minutes())

	// return if there is no partial time
	timeRemainder := durationInt64 % fixedTimePeriodInt64
	if timeRemainder == 0 {
		return s, nil
	}

	weight := float64(timeRemainder) / float64(fixedTimePeriodInt64)
	l := len(s)
	if l < 1 {
		return nil, fmt.Errorf("slice length too short")
	}

	// calculate weight of last element
	i := float64(s[l-1].Carbon.Intensity) * weight
	s[l-1].Carbon.Intensity = int64(math.Round(i))
	timeRemainderDuration, err := time.ParseDuration(fmt.Sprint(strconv.Itoa(int(timeRemainder)), "m"))
	if err != nil {
		return nil, fmt.Errorf("error parsing duration: %v", err)
	}

	// update entry.ValidTo with correct duraion from the start of the last element
	s[l-1].ValidTo = s[l-1].ValidFrom.Add(timeRemainderDuration)

	return s, nil
}

// CalculateContinuousPeriodIntensity accepts a slice of models.Period and a time.Duration variable, and returns
// a pointer to models.Slot, or an error. It will calculate the average of all elements in the slice, taking into
// account partial time periods.
func (c *CalculationService) CalculateContinuousPeriodIntensity(pArr []models.Period, duration time.Duration) (*models.Slot, error) {
	weight := 0.0
	totalIntensity := 0.0

	arrLength := len(pArr)
	if arrLength < 1 {
		return nil, fmt.Errorf("array length 0")
	}

	lastElem := pArr[arrLength-1]

	timeRemainder := int(duration.Minutes()) % int(fixedTimePeriod.Minutes())
	if timeRemainder <= 0 {
		weight = 1.0

		for _, p := range pArr {
			totalIntensity += float64(p.Intensity.Forecast)
		}

		averageIntensity := totalIntensity / (float64(arrLength) - 1 + weight)

		intensity := int64(math.Round(averageIntensity))
		s := &models.Slot{
			ValidFrom: pArr[0].From,
			ValidTo:   lastElem.To,
			Carbon: models.Carbon{
				Intensity: intensity,
			},
		}

		return s, nil

	}

	// handle durations that are not multiples of 30
	// e.g. 13, 45, 61
	weight = float64(timeRemainder) / float64(fixedTimePeriod.Minutes())

	totalIntensity = float64(lastElem.Intensity.Forecast) * weight

	for i := range arrLength {
		// skip the last element of the slice
		if i == arrLength-1 {
			continue
		}
		totalIntensity += float64(pArr[i].Intensity.Forecast)
	}

	timeRemainderDuration, err := time.ParseDuration(fmt.Sprint(strconv.Itoa(int(fixedTimePeriod.Minutes())-timeRemainder), "m"))
	if err != nil {
		return nil, fmt.Errorf("error parsing duration: %v", err)
	}

	averageIntensity := totalIntensity
	if arrLength != 1 {
		averageIntensity = totalIntensity / (float64(arrLength) - 1 + weight)
	}

	intensity := int64(math.Round(averageIntensity))
	s := &models.Slot{
		ValidFrom: pArr[0].From,
		ValidTo:   lastElem.To.Add(-timeRemainderDuration),
		Carbon: models.Carbon{
			Intensity: intensity,
		},
	}
	return s, nil
}

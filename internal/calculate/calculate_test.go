package calculate

import (
	"fmt"
	"testing"
	"time"

	"github.com/lstratta/flatpeak-take-home-task/internal/neso"
)

func Test_FilterPeriodsByLowInensity_ReturnSlots(t *testing.T) {
	var tests = []struct {
		target   int
		duration string
	}{
		{1, "30m"},
		{2, "60m"},
		{5, "150m"},
		{2, "45m"},
		{3, "61m"},
		{7, "184m"},
		{10, "1440m"},
	}

	d := genData()
	for _, tab := range tests {
		dur, err := time.ParseDuration(tab.duration)
		if err != nil {
			t.Errorf("error parsing duration: %v", err)
		}

		pArr, err := FilterPeriodsByLowestIntensity(d.Data, dur)
		if err != nil {
			t.Errorf("error filtering low intensity slots: %v", err)
		}

		actualLen := len(pArr)

		if actualLen != tab.target {
			t.Errorf("required length: %d, actual: %d", tab.target, actualLen)
		}

		fmt.Printf("%v\n", pArr)

	}

}

func Test_FilterPeriodsByDuration_ReturnsSlots(t *testing.T) {
	var tests = []struct {
		target   int
		duration string
	}{
		{1, "30m"},
		{2, "60m"},
		{5, "150m"},
		{2, "45m"},
		{3, "61m"},
	}

	d := genData()

	for _, tab := range tests {
		dur, err := time.ParseDuration(tab.duration)
		if err != nil {
			t.Errorf("error parsing duration: %v", err)
		}

		ls, err := FilterPeriodsByDuration(d.Data, dur)
		if err != nil {
			t.Errorf("error filtering by duration: %v", err)
		}

		actualLen := len(ls)
		if actualLen != tab.target {
			t.Logf("test: target: %d, duration: %s", tab.target, tab.duration)
			t.Errorf("required length: %d, actual: %d", tab.target, actualLen)
		}
	}
}

func Test_CalculateNonContinuousPeriod_ReturnsCorrectAverage(t *testing.T) {

}

func Test_CalculateContinuousPeriod_ReturnsCorrectAverage(t *testing.T) {
	var tests = []struct {
		target   int64
		duration string
	}{
		{56, "30m"},
		{77, "60m"},
		{55, "90m"},
		{70, "45m"},
		{75, "61m"},
	}

	d := genData()
	for _, tab := range tests {
		dur, err := time.ParseDuration(tab.duration)
		if err != nil {
			t.Errorf("error parsing duration: %v", err)
		}

		n, err := FilterPeriodsByDuration(d.Data, dur)
		if err != nil {
			t.Errorf("error FilterPeriodsByDuration: %v", err)
		}

		//if tab.duration == "45m" {
		//	fmt.Printf("%v\n", n)
		//}

		actual, err := CalculateContinuousPeriod(n, dur)
		if err != nil {
			t.Errorf("error calculating continuous period: %v", err)
		}

		if actual != tab.target {
			t.Errorf("averaged intensity - expected: %d, actual: %d", tab.target, actual)
		}
	}
}

func genData() *neso.Data {
	return &neso.Data{
		Data: []neso.Period{
			{
				From: "2026-03-24T10:30Z",
				To:   "2026-03-24T11:00Z",
				Intensity: neso.Intensity{
					Forecast: 56,
					Actual:   0,
					Index:    "low",
				},
			},
			{
				From: "2026-03-24T11:00Z",
				To:   "2026-03-24T11:30Z",
				Intensity: neso.Intensity{
					Forecast: 97,
					Actual:   0,
					Index:    "low",
				},
			},
			{
				From: "2026-03-24T11:30Z",
				To:   "2026-03-24T12:00Z",
				Intensity: neso.Intensity{
					Forecast: 13,
					Actual:   0,
					Index:    "very low",
				},
			},
			{
				From: "2026-03-24T12:00Z",
				To:   "2026-03-24T12:30Z",
				Intensity: neso.Intensity{
					Forecast: 170,
					Actual:   0,
					Index:    "high",
				},
			},
			{
				From: "2026-03-24T12:30Z",
				To:   "2026-03-24T13:00Z",
				Intensity: neso.Intensity{
					Forecast: 83,
					Actual:   0,
					Index:    "low",
				},
			},
			{
				From: "2026-03-24T13:00Z",
				To:   "2026-03-24T13:30Z",
				Intensity: neso.Intensity{
					Forecast: 130,
					Actual:   0,
					Index:    "moderate",
				},
			},
			{
				From: "2026-03-24T13:30Z",
				To:   "2026-03-24T14:00Z",
				Intensity: neso.Intensity{
					Forecast: 83,
					Actual:   0,
					Index:    "low",
				},
			},
			{
				From: "2026-03-24T14:00Z",
				To:   "2026-03-24T14:30Z",
				Intensity: neso.Intensity{
					Forecast: 385,
					Actual:   0,
					Index:    "high",
				},
			},
			{
				From: "2026-03-24T14:30Z",
				To:   "2026-03-24T15:00Z",
				Intensity: neso.Intensity{
					Forecast: 138,
					Actual:   0,
					Index:    "moderate",
				},
			},
			{
				From: "2026-03-24T15:00Z",
				To:   "2026-03-24T15:30Z",
				Intensity: neso.Intensity{
					Forecast: 1,
					Actual:   0,
					Index:    "very low",
				},
			},
		},
	}

}

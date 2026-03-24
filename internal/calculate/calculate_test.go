package calculate

import (
	"testing"
	"time"

	"github.com/lstratta/flatpeak-take-home-task/internal/neso"
)

func Test_Example(t *testing.T) {
	// t.Errorf("test failed")
}

func Test_FilterPeriodsByLowInensity_Returns2Slots(t *testing.T) {
	targetLen := 2
	d := genData()
	ls, err := FilterPeriodsByLowIntensity(d.Data)
	if err != nil {
		t.Errorf("error filtering low intensity slots: %v", err)
	}

	actualLen := len(ls)

	if actualLen != targetLen {
		t.Errorf("required length: %d, actual: %d", targetLen, actualLen)
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

//	func Test_CalculateContinuousPeriod_ReturnsCorrectAverage(t *testing.T) {
//		var tests = []struct {
//			target int64
//		}{
//			{1},
//		}
//
//		d := genData()
//
//		for _, tab := range tests {
//			slot, err := CalculateContinuousPeriod(d.Data)
//			if err != nil {
//				t.Errorf("error calculating continuous period: %v", err)
//			}
//
//			actual := slot.Intensity.Forecast
//
//			if actual != tab.target {
//				t.Errorf("averaged intensity - expected: %d, actual: %d", tab.target, actual)
//			}
//		}
//	}
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
					Forecast: 284,
					Actual:   0,
					Index:    "high",
				},
			},
			{
				From: "2026-03-24T12:00Z",
				To:   "2026-03-24T12:30Z",
				Intensity: neso.Intensity{
					Forecast: 284,
					Actual:   0,
					Index:    "high",
				},
			},
			{
				From: "2026-03-24T12:30Z",
				To:   "2026-03-24T13:00Z",
				Intensity: neso.Intensity{
					Forecast: 284,
					Actual:   0,
					Index:    "high",
				},
			},
			{
				From: "2026-03-24T13:00Z",
				To:   "2026-03-24T13:30Z",
				Intensity: neso.Intensity{
					Forecast: 284,
					Actual:   0,
					Index:    "high",
				},
			},
		},
	}

}

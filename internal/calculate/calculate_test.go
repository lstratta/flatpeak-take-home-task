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

func Test_FilterPeriodsByDuration_WhereContinuousIsFalse_Returns1Slot(t *testing.T) {
	d := genData()
	targetLen := 1
	dur, err := time.ParseDuration("30m")
	if err != nil {
		t.Errorf("error parsing duration: %v", err)
	}

	ls, err := FilterPeriodsByDuration(d.Data, dur)
	if err != nil {
		t.Errorf("error filtering by duration: %v", err)
	}

	actualLen := len(ls)
	if actualLen != targetLen {
		t.Errorf("required length: %d, actual: %d", targetLen, actualLen)
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
					Forecast: 284,
					Actual:   0,
					Index:    "high",
				},
			},
		},
	}

}

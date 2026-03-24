package calculate

import (
	"testing"

	"github.com/lstratta/flatpeak-take-home-task/internal/neso"
)

func Test_Example(t *testing.T) {
	// t.Errorf("test failed")
}

func Test_FilterLowInensitySlots_Returns2Slots(t *testing.T) {
	targetLen := 2
	d := genData()
	ls, err := FilterLowIntensitySlots(d)
	if err != nil {
		t.Errorf("error filtering low intensity slots: %v", err)
	}

	actualLen := len(ls)

	if actualLen != targetLen {
		t.Errorf("required length: %d, actual: %d", targetLen, actualLen)
	}

}

func Test_FilterByDurationWhereContinuousIsFalse_Returns1Slot(t *testing.T) {
	d := genData()
	targetLen := 1

	ls, err := FilterByDuration(d, 30, false)
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
		Data: []neso.Neso{
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
				To:   "2026-03-24T11:00Z",
				Intensity: neso.Intensity{
					Forecast: 284,
					Actual:   0,
					Index:    "high",
				},
			},
		},
	}

}

package calculate

import (
	"testing"

	"github.com/lstratta/flatpeak-take-home-task/internal/neso"
)

func Test_Example(t *testing.T) {
	// t.Errorf("test failed")
}

func Test_FilterLowInensitySlots_Returns2Slots(t *testing.T) {
	target := 2

	ls, err := FilterLowIntensitySlots(genData())
	if err != nil {
		t.Errorf("error filtering low intensity slots: %v", err)
	}

	actual := len(ls)

	if actual != 2 {
		t.Errorf("required length: %d, actual: %d", target, actual)
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

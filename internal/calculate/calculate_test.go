package calculate

import (
	"testing"
	"time"

	"github.com/lstratta/flatpeak-take-home-task/internal/models"
)

func Test_FilterPeriodsByLowestInensity_ReturnSlots(t *testing.T) {
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
		{12, "1440m"},
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

func Test_CalculateWeightedAverageForLastPeriodInSlice_ReturnsCorrectIntensity(t *testing.T) {
	var tests = []struct {
		target   int64
		duration string
	}{
		{1, "30m"},
		{18, "60m"},
		{83, "90m"},
		{51, "45m"},
		{3, "61m"},
	}

	d := genData()
	for _, tab := range tests {
		dur, err := time.ParseDuration(tab.duration)
		if err != nil {
			t.Errorf("error parsing duration: %v", err)
		}

		pArr, err := FilterPeriodsByLowestIntensity(d.Data, dur)
		if err != nil {
			t.Errorf("error FilterPeriodsByLowestIntensity: %v", err)
		}

		slots, err := CalculateWeightedAverage(pArr, dur)
		if err != nil {
			t.Errorf("error calculating continuous period: %v", err)
		}

		// if tab.duration == "45m" {
		// 	jsonData, _ := json.MarshalIndent(slots, "", "  ")
		// 	fmt.Println(string(jsonData))
		// }
		// if tab.duration == "61m" {
		// 	jsonData, _ := json.MarshalIndent(slots, "", "  ")
		// 	fmt.Println(string(jsonData))
		// }

		l := len(slots)
		if l < 1 {
			t.Errorf("length too short")
		}
		actual := slots[l-1].Carbon.Intensity

		if actual != tab.target {
			t.Errorf("averaged intensity - expected: %d, actual: %d", tab.target, actual)
		}
	}

}

func Test_CalculateContinuousPeriod_ReturnsCorrectAverageIntensity(t *testing.T) {
	var tests = []struct {
		target   int64
		duration string
	}{
		{56, "30m"},
		{77, "60m"},
		{55, "90m"},
		{56, "45m"},
		{76, "61m"},
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

		actual, err := CalculateContinuousPeriodIntensity(n, dur)
		if err != nil {
			t.Errorf("error calculating continuous period: %v", err)
		}

		if actual.Carbon.Intensity != tab.target {
			t.Errorf("averaged intensity - expected: %d, actual: %d", tab.target, actual.Carbon.Intensity)
		}
	}
}

func genData() *models.Data {
	return &models.Data{
		Data: []models.Period{
			{
				From: time.Date(2026, time.Month(3), 24, 10, 30, 00, 00, time.Now().UTC().Location()),
				To:   time.Date(2026, time.Month(3), 24, 11, 00, 00, 00, time.Now().UTC().Location()),
				Intensity: models.Intensity{
					Forecast: 56,
					Actual:   0,
					Index:    "low",
				},
			},
			{
				From: time.Date(2026, time.Month(3), 24, 11, 00, 00, 00, time.Now().UTC().Location()),
				To:   time.Date(2026, time.Month(3), 24, 11, 30, 00, 00, time.Now().UTC().Location()),
				Intensity: models.Intensity{
					Forecast: 97,
					Actual:   0,
					Index:    "low",
				},
			},
			{
				From: time.Date(2026, time.Month(3), 24, 11, 30, 00, 00, time.Now().UTC().Location()),
				To:   time.Date(2026, time.Month(3), 24, 12, 00, 00, 00, time.Now().UTC().Location()),
				Intensity: models.Intensity{
					Forecast: 13,
					Actual:   0,
					Index:    "very low",
				},
			},
			{
				From: time.Date(2026, time.Month(3), 24, 12, 00, 00, 00, time.Now().UTC().Location()),
				To:   time.Date(2026, time.Month(3), 24, 12, 30, 00, 00, time.Now().UTC().Location()),
				Intensity: models.Intensity{
					Forecast: 170,
					Actual:   0,
					Index:    "high",
				},
			},
			{
				From: time.Date(2026, time.Month(3), 24, 12, 30, 00, 00, time.Now().UTC().Location()),
				To:   time.Date(2026, time.Month(3), 24, 13, 00, 00, 00, time.Now().UTC().Location()),
				Intensity: models.Intensity{
					Forecast: 83,
					Actual:   0,
					Index:    "low",
				},
			},
			{
				From: time.Date(2026, time.Month(3), 24, 13, 00, 00, 00, time.Now().UTC().Location()),
				To:   time.Date(2026, time.Month(3), 24, 13, 30, 00, 00, time.Now().UTC().Location()),
				Intensity: models.Intensity{
					Forecast: 130,
					Actual:   0,
					Index:    "moderate",
				},
			},
			{
				From: time.Date(2026, time.Month(3), 24, 13, 30, 00, 00, time.Now().UTC().Location()),
				To:   time.Date(2026, time.Month(3), 24, 14, 00, 00, 00, time.Now().UTC().Location()),
				Intensity: models.Intensity{
					Forecast: 83,
					Actual:   0,
					Index:    "low",
				},
			},
			{
				From: time.Date(2026, time.Month(3), 24, 14, 00, 00, 00, time.Now().UTC().Location()),
				To:   time.Date(2026, time.Month(3), 24, 14, 30, 00, 00, time.Now().UTC().Location()),
				Intensity: models.Intensity{
					Forecast: 385,
					Actual:   0,
					Index:    "high",
				},
			},
			{
				From: time.Date(2026, time.Month(3), 24, 14, 30, 00, 00, time.Now().UTC().Location()),
				To:   time.Date(2026, time.Month(3), 24, 15, 00, 00, 00, time.Now().UTC().Location()),
				Intensity: models.Intensity{
					Forecast: 138,
					Actual:   0,
					Index:    "moderate",
				},
			},
			{
				From: time.Date(2026, time.Month(3), 24, 15, 00, 00, 00, time.Now().UTC().Location()),
				To:   time.Date(2026, time.Month(3), 24, 15, 30, 00, 00, time.Now().UTC().Location()),
				Intensity: models.Intensity{
					Forecast: 1,
					Actual:   0,
					Index:    "very low",
				},
			},
			{
				From: time.Date(2026, time.Month(3), 24, 15, 30, 00, 00, time.Now().UTC().Location()),
				To:   time.Date(2026, time.Month(3), 24, 16, 00, 00, 00, time.Now().UTC().Location()),
				Intensity: models.Intensity{
					Forecast: 18,
					Actual:   0,
					Index:    "very low",
				},
			},
			{
				From: time.Date(2026, time.Month(3), 24, 16, 00, 00, 00, time.Now().UTC().Location()),
				To:   time.Date(2026, time.Month(3), 24, 16, 30, 00, 00, time.Now().UTC().Location()),
				Intensity: models.Intensity{
					Forecast: 101,
					Actual:   0,
					Index:    "moderate",
				},
			},
		},
	}

}

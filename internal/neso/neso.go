package neso

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lstratta/flatpeak-take-home-task/internal/models"
)

const (
	NESO_API            string = "https://api.carbonintensity.org.uk"
	NESO_INTENSITY_PATH string = "intensity"
	NESO_FW24H          string = "fw24h"
)

type response struct {
	Data []interval
}

type interval struct {
	From      string           `json:"from"`
	To        string           `json:"to"`
	Intensity models.Intensity `json:"intensity"`
}

// Fetch data from NESO API
func GetNesoData() (*models.Data, error) {
	from := time.Now().Format(time.RFC3339) // RFC3339 = "2006-01-02T15:04:05Z07:00"
	url := fmt.Sprintf("%s/%s/%s/%s", NESO_API, NESO_INTENSITY_PATH, from, NESO_FW24H)

	res, err := makeRequestGET(url)
	if err != nil {
		return nil, fmt.Errorf("error making GET request to NESO: %v", err)
	}

	d, err := convertResToData(res)
	if err != nil {
		return nil, fmt.Errorf("error converting response: %v", err)
	}

	return d, nil
}

func convertResToData(res *http.Response) (*models.Data, error) {
	r := response{}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error decoding json data: %v", err)
	}

	for i := range r.Data {
		newDateFormat := r.Data[i].From[:len(r.Data[i].From)-1] // remove the Z
		newDateFormat = newDateFormat + ":00Z"                  // add in the missing seconds
		r.Data[i].From = newDateFormat

		newDateFormat = r.Data[i].To[:len(r.Data[i].To)-1] // remove the Z
		newDateFormat = newDateFormat + ":00Z"             // add in the missing seconds
		r.Data[i].To = newDateFormat
	}

	data, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling data: %v", err)
	}

	d := models.Data{}
	err = json.Unmarshal(data, &d)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling data: %v", err)
	}

	return &d, nil
}

func makeRequestGET(url string) (*http.Response, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error execuing GET request: %v", err)
	}
	return res, nil
}

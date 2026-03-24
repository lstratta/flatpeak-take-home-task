package neso

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Data struct {
	Data []Neso `json:"data"`
}

type Neso struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Intensity Intensity `json:"intensity"`
}

type Intensity struct {
	Forecast int64  `json:"forecast"`
	Actual   int64  `json:"actual"`
	Index    string `json:"index"`
}

const (
	NESO_API            string = "https://api.carbonintensity.org.uk"
	NESO_INTENSITY_PATH string = "intensity"
	NESO_FW24H          string = "fw24h"
)

func GetNesoData() (*Data, error) {
	from := time.Now().Format(time.RFC3339) // RFC3339 = "2006-01-02T15:04:05Z07:00"
	url := fmt.Sprintf("%s/%s/%s/%s", NESO_API, NESO_INTENSITY_PATH, from, NESO_FW24H)
	fmt.Println(url)

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

func convertResToData(res *http.Response) (*Data, error) {
	d := Data{}

	if err := json.NewDecoder(res.Body).Decode(&d); err != nil {
		return nil, fmt.Errorf("error decoding json data: %v", err)
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

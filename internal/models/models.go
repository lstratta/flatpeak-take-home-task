package models

import "time"

type Slot struct {
	ValidFrom time.Time `json:"valid_from"`
	ValidTo   time.Time `json:"valid_to"`
	Carbon    Carbon    `json:"carbon"`
}

type Carbon struct {
	Intensity int64 `json:"intensity"`
}

type Data struct {
	Data []Period `json:"data"`
}

type Period struct {
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
	Intensity Intensity `json:"intensity"`
}

type Intensity struct {
	Forecast int64  `json:"forecast"`
	Actual   int64  `json:"actual"`
	Index    string `json:"index"`
}

type ByDateSorter []Period

func (b ByDateSorter) Len() int           { return len(b) }
func (b ByDateSorter) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByDateSorter) Less(i, j int) bool { return b[i].From.Before(b[j].From) }

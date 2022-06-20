package jsonmodels

import (
	"fmt"
	"time"
)

type EventModel struct {
	Desc string `json:"desc"`
	Date string `json:"date"`
}

type EventModelUpdate struct {
	ID   int    `json:"id"`
	Desc string `json:"desc,omitempty"`
	Date string `json:"date,omitempty"`
}

type EventModelDelete struct {
	ID string `json:"id"`
}

const layout = "2006-01-02"

func ParseDate(s string) (time.Time, error) {
	if s == "null" {
		return time.Time{}, fmt.Errorf("time doesn't exists")
	}
	result, err := time.Parse(layout, s)
	return result, err
}

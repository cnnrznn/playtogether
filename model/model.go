package model

import (
	"fmt"

	"github.com/google/uuid"
)

type Ping struct {
	Activity string  `json:"activity"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Expire   int     `json:"expire"`
	RangeKM  int     `json:"range_km"`

	// other stuff like age, intensity
}

type Game struct {
	Id       uuid.UUID `json:"id"`
	Activity string    `json:"activity"`
	Lat      float64   `json:"lat"`
	Lon      float64   `json:"lon"`
	Players  []string  `json:"players"`
}

type Area struct {
	LatMin, LatMax float64
	LonMin, LonMax float64
}

func (a Area) String() string {
	return fmt.Sprintf("Lat: [%v, %v], Lon: [%v, %v]",
		a.LatMin, a.LatMax,
		a.LonMin, a.LonMax,
	)
}

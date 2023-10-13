package model

import (
	"fmt"

	"github.com/google/uuid"
)

type Ping struct {
	ID       uuid.UUID `json:"id"`
	Player   uuid.UUID `json:"player"`
	Activity string    `json:"activity"`
	Lat      float64   `json:"lat"`
	Lon      float64   `json:"lon"`
	Expire   int       `json:"expire"`
	RangeKM  int       `json:"range_km"`

	// other stuff like age, intensity
}

type Player struct {
	ID uuid.UUID `json:"id"`
}

type Game struct {
	ID       uuid.UUID   `json:"id"`
	Activity string      `json:"activity"`
	Lat      float64     `json:"lat"`
	Lon      float64     `json:"lon"`
	Players  []uuid.UUID `json:"players"`
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

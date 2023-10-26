package model

import (
	"fmt"

	"github.com/google/uuid"
)

type PlayRequest struct {
	User     uuid.UUID `json:"user"`     // user who initiated the request
	Size     int       `json:"size"`     // number of people in group
	Activity string    `json:"activity"` // sport or activity type
	Lat      float64   `json:"lat"`      // location
	Lon      float64   `json:"lon"`
	Start    int       `json:"start"`    // start time
	End      int       `json:"end"`      // end time
	RangeKM  int       `json:"range_km"` // distance away they'd travel

	// other stuff like age, intensity
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

type Game struct {
	Version      uuid.UUID
	ID           uuid.UUID
	PlayRequests []uuid.UUID
	Status       GameStatus
	Activity     string
	Lat          float64
	Lon          float64
}

type GameStatus int

const (
	PENDING GameStatus = iota
	ACTIVE
)

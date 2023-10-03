package play

import (
	"fmt"
	"math"

	"github.com/jftuga/geodist"
	_ "github.com/lib/pq"

	"github.com/cnnrznn/playtogether/db"
	"github.com/cnnrznn/playtogether/model"
)

var (
	initDone bool = false
)

type Response struct {
	Found bool       `json:"found"`
	Game  model.Game `json:"game,omitempty"`
}

type Area struct {
	latMin, latMax float64
	lonMin, lonMax float64
}

func (a Area) String() string {
	return fmt.Sprintf("Lat: [%v, %v], Lon: [%v, %v]",
		a.latMin, a.latMax,
		a.lonMin, a.lonMax,
	)
}

func Update(ping model.Ping) (*Response, error) {
	// calculate lat,lon range for game or other players
	calculateArea(ping)

	// First, check for games already going on in the area

	// If no games found, put player into players DB and try to create a game with the new player information

	// If game created, put game into games table and send alerts to players

	return nil, fmt.Errorf("not implemented")
}

func calculateArea(ping model.Ping) (*Area, error) {
	latDelta := calculateDegree(ping, true)
	lonDelta := calculateDegree(ping, false)

	return &Area{
		latMin: ping.Lat - latDelta,
		latMax: ping.Lat + latDelta,
		lonMin: ping.Lon - lonDelta,
		lonMax: ping.Lon + lonDelta,
	}, nil
}

func calculateDegree(ping model.Ping, latitude bool) float64 {
	const tolerance float64 = 0.001
	var target float64 = float64(ping.RangeKM)

	var minDeg float64 = 0.00
	var maxDeg float64 = 1.00
	var midDeg float64
	var midCoord geodist.Coord

	// iteratively binary search to find degree in tolerance
	for {
		midDeg = (minDeg + maxDeg) / 2
		switch latitude {
		case true:
			midCoord = geodist.Coord{
				Lat: ping.Lat + midDeg,
				Lon: ping.Lon,
			}
		default:
			midCoord = geodist.Coord{
				Lat: ping.Lat,
				Lon: ping.Lon + midDeg,
			}
		}

		// calculate distance to mid
		_, km, _ := geodist.VincentyDistance(
			geodist.Coord{
				Lat: ping.Lat,
				Lon: ping.Lon,
			},
			midCoord,
		)

		if math.Abs(target-km) < tolerance {
			break
		}

		if target-km < 0 { // point is farther away than target
			maxDeg = midDeg
		} else {
			minDeg = midDeg
		}
	}

	return midDeg
}

func Init() error {
	if !initDone {
		err := db.Init()
		if err != nil {
			return err
		}

		initDone = true
	}

	return nil
}

func Run() error {
	// every 5m, scan and delete rows in db past expiration
	if err := Init(); err != nil {
		return err
	}

	return runExpire()
}

func runExpire() error {
	return nil
}

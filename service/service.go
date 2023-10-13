package service

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/jftuga/geodist"

	"github.com/cnnrznn/playtogether/db"
	"github.com/cnnrznn/playtogether/model"
)

var (
	initDone bool = false
)

type Response struct {
	Found bool         `json:"found"`
	Games []model.Game `json:"games,omitempty"`
}

func GetPlayerGames(player model.Player) []model.Game {
	return db.LoadPlayerGames(player)
}

func Update(ping model.Ping) (*Response, error) {
	area := calculateArea(ping)

	ping.ID = uuid.New()
	id, err := db.StorePing(ping)
	if err != nil {
		return nil, err
	}
	ping.ID = *id

	// First, check for games already going on in the area
	games, err := db.LoadGamesByArea(ping, area)
	if err != nil {
		return nil, err
	}

	if len(games) > 0 {
		for _, game := range games {
			game.Players = append(game.Players, ping.Player)
			db.StorePlayerGame(ping, game)
		}

		return &Response{
			Found: true,
			Games: games,
		}, nil
	}

	// Load all players in area for activity
	players, err := db.LoadPingsByArea(ping.Activity, area)
	if err != nil {
		return nil, err
	}

	// verify there are enough players to create a game, namely
	// 1. iterate over all players and calculate distance to new player
	// 2. if <activity::threshold> players are in <rangeKM> of new player, create game
	filteredPlayers := []model.Ping{}
	for _, player := range players {
		_, km, _ := geodist.VincentyDistance(
			geodist.Coord{Lat: player.Lat, Lon: player.Lon},
			geodist.Coord{Lat: ping.Lat, Lon: ping.Lon},
		)
		if km < float64(player.RangeKM) {
			filteredPlayers = append(filteredPlayers, player)
		}
	}

	if atThreshold(ping.Activity, filteredPlayers) {
		filteredPlayerIDs := []uuid.UUID{}
		for _, ping := range filteredPlayers {
			filteredPlayerIDs = append(filteredPlayerIDs, ping.Player)
		}
		game := model.Game{
			ID:       uuid.New(),
			Activity: ping.Activity,
			Lat:      ping.Lat,
			Lon:      ping.Lon,
			Players:  filteredPlayerIDs,
		}

		// put new game in DB
		err := db.StoreGame(game)
		if err != nil {
			return nil, err
		}

		for _, player := range filteredPlayers {
			db.StorePlayerGame(player, game)
		}

		return &Response{
			Found: true,
			Games: []model.Game{game},
		}, nil
	}

	return &Response{
		Found: false,
	}, nil
}

func calculateArea(ping model.Ping) model.Area {
	latDelta := calculateDegree(ping, true)
	lonDelta := calculateDegree(ping, false)

	return model.Area{
		LatMin: ping.Lat - latDelta,
		LatMax: ping.Lat + latDelta,
		LonMin: ping.Lon - lonDelta,
		LonMax: ping.Lon + lonDelta,
	}
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
	if err := Init(); err != nil {
		return err
	}

	return runExpire()
}

func runExpire() error {
	ticker := time.NewTicker(1 * time.Minute)

	for range ticker.C {
		err := db.Expire()
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func atThreshold(activity string, players []model.Ping) bool {
	switch activity {
	case "volleyball":
		return len(players) >= 4
	default:
		return len(players) >= 2
	}
}

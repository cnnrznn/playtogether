package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/cnnrznn/playtogether/model"
)

const (
	DB_CONN = "hostname=localhost user=pt database=playtogether sslmode=disable"
)

var (
	db       *sql.DB = nil
	initDone bool    = false
)

func Init() error {
	if !initDone {

		conn, err := sql.Open(
			"postgres",
			fmt.Sprintf(
				"%v password='%v'", DB_CONN, os.Getenv("DB_PASSWD"),
			),
		)
		if err != nil {
			return err
		}

		db = conn
		initDone = true
	}

	return nil
}

func GetGames(ping model.Ping, area model.Area) ([]model.Game, error) {
	games := []model.Game{}

	rows, err := db.Query(`
		SELECT (id, lat, lon) FROM games
			WHERE
				activity = $5 AND
				lat < $1 and lat > $2 AND lon < $3 AND lon > $4
	`, area.LatMax, area.LatMin, area.LonMax, area.LonMin,
		ping.Activity)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		game := model.Game{}

		err := rows.Scan(
			&game.Id,
			&game.Lat,
			&game.Lon,
		)
		if err != nil {
			return nil, err
		}

		games = append(games, game)
	}

	return games, nil
}

func NewPlayer(model.Ping, model.Area) *model.Game {
	return nil
}

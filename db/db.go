package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/cnnrznn/playtogether/model"
	"github.com/google/uuid"
)

const (
	DB_CONN = "host=localhost user=pt database=playtogether sslmode=disable"
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

func LoadGames(ping model.Ping, area model.Area) ([]model.Game, error) {
	games := []model.Game{}

	rows, err := db.Query(`
		SELECT id, lat, lon FROM games
			WHERE
				activity = $5 AND
				lat < $1 AND lat > $2 AND lon < $3 AND lon > $4`,
		area.LatMax, area.LatMin, area.LonMax, area.LonMin, ping.Activity,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func LoadPings(activity string, area model.Area) ([]model.Ping, error) {
	pings := []model.Ping{}

	rows, err := db.Query(`
		SELECT player, lat, lon, range_km FROM ping
		WHERE
			activity = $1 AND
			lat < $2 AND lat > $3 AND lon < $4 AND lon > $5`,
		activity, area.LatMax, area.LatMin, area.LonMax, area.LonMin,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		ping := model.Ping{}

		err := rows.Scan(
			&ping.Player,
			&ping.Lat,
			&ping.Lon,
			&ping.RangeKM,
		)
		if err != nil {
			return nil, err
		}

		pings = append(pings, ping)
	}

	return pings, nil
}

func StorePing(ping model.Ping) error {
	id := uuid.New()

	result, err := db.Exec(`
			INSERT INTO ping (id, player, lat, lon, range_km, expire, activity) VALUES
			($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT ON CONSTRAINT unq_player_activity
			DO UPDATE SET
				lat=$3, lon=$4, range_km=$5, expire=$6`,
		id, ping.Player, ping.Lat, ping.Lon, ping.RangeKM, ping.Expire, ping.Activity,
	)
	if err != nil {
		return err
	}
	if n, err := result.RowsAffected(); err != nil || n != 1 {
		return fmt.Errorf("could not insert ping in table: %w", err)
	}

	return nil
}

func StoreGame(game model.Game) error {
	result, err := db.Exec(`
		INSERT INTO games (id, lat, lon, activity) VALUES
			($1, $2, $3, $4)`,
		game.Id, game.Lat, game.Lon, game.Activity,
	)
	if err != nil {
		return err
	}
	if n, err := result.RowsAffected(); err != nil || n != 1 {
		return fmt.Errorf("could not insert game in table: %w", err)
	}

	return nil
}

func Expire() {
	db.Exec(`
		DELETE FROM ping
		WHERE expire < $1`,
		time.Now().Unix(),
	)
}

func StorePlayerGame(playerID uuid.UUID, game model.Game) {

}

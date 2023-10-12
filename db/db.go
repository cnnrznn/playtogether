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
		SELECT id, player, lat, lon, range_km FROM ping
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
			&ping.ID,
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
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := tx.Query(`
		DELETE FROM ping
		WHERE expire < $1
		RETURNING player, id, game`,
		time.Now().Unix(),
	)
	if err != nil {
		tx.Rollback()
		return
	}
	defer result.Close()

	pings := []model.Ping{}

	for result.Next() {
		var playerID, pingID uuid.UUID
		err := result.Scan(&playerID, &pingID)
		if err != nil {
			tx.Rollback()
			return
		}
		pings = append(pings, model.Ping{
			Player: playerID,
			ID:     pingID,
		})
	}
	if result.Err() != nil {
		tx.Rollback()
		return
	}

	for _, ping := range pings {
		_, err := tx.Exec(`
			DELETE FROM player2game
			WHERE ping=$1`,
			ping.ID)
		if err != nil {
			tx.Rollback()
			return
		}
	}

	tx.Commit()
}

func StorePlayerGame(ping model.Ping, game model.Game) {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = tx.Exec(`
			INSERT INTO player2game (player, game, ping)
			VALUES
				($1, $2, $3)
			ON CONFLICT DO NOTHING`,
		ping.Player, game.Id, ping.ID,
	)
	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
}

func LoadPlayerGames(player model.Player) []model.Game {
	result := []model.Game{}
	gameIDs := []uuid.UUID{}

	rows, err := db.Query(`
		SELECT game from player2game WHERE player=$1`,
		player.ID)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var gameID uuid.UUID
		if err := rows.Scan(&gameID); err != nil {
			fmt.Println(err)
			break
		}

		gameIDs = append(gameIDs, gameID)
	}
	if rows.Err() != nil {
		fmt.Println(rows.Err())
		return nil
	}

	for _, gameID := range gameIDs {
		var game model.Game

		row := db.QueryRow(`SELECT id, lat, lon, activity from games WHERE id=$1`, gameID)
		if row.Err() != nil {
			fmt.Println(err)
			continue
		}

		if err := row.Scan(
			&game.Id,
			&game.Lat,
			&game.Lon,
			&game.Activity,
		); err != nil {
			fmt.Println(err)
			continue
		}

		result = append(result, game)
	}

	return result
}

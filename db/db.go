package db

import (
	"database/sql"
	"encoding/json"
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

func LoadGamesByArea(ping model.Ping, area model.Area) ([]model.Game, error) {
	games := []model.Game{}

	rows, err := db.Query(`
		SELECT id, lat, lon, activity FROM games
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
			&game.Activity,
		)
		if err != nil {
			return nil, err
		}

		games = append(games, game)
	}

	return games, nil
}

func LoadPingsByArea(activity string, area model.Area) ([]model.Ping, error) {
	pings := []model.Ping{}

	rows, err := db.Query(`
		SELECT id, player, lat, lon, activity, range_km, expire FROM ping
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
			&ping.Activity,
			&ping.RangeKM,
			&ping.Expire,
		)
		if err != nil {
			return nil, err
		}

		pings = append(pings, ping)
	}

	return pings, nil
}

func StorePing(ping model.Ping) (*uuid.UUID, error) {
	result, err := db.Exec(`
			INSERT INTO ping (id, player, lat, lon, range_km, expire, activity)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT ON CONSTRAINT unq_player_activity
			DO UPDATE SET
				lat=$3, lon=$4, range_km=$5, expire=$6`,
		ping.ID, ping.Player, ping.Lat, ping.Lon, ping.RangeKM, ping.Expire, ping.Activity,
	)
	if err != nil {
		return nil, err
	}
	if n, err := result.RowsAffected(); err != nil || n != 1 {
		return nil, fmt.Errorf("could not insert ping in table: %w", err)
	}

	row := db.QueryRow(`
		SELECT id from ping
		WHERE player=$1 AND activity=$2`,
		ping.Player, ping.Activity)

	var pingID uuid.UUID

	err = row.Scan(&pingID)
	if err != nil {
		return nil, err
	}

	return &pingID, nil
}

func StoreGame(game model.Game) error {
	bs, err := json.Marshal(game.Players)
	if err != nil {
		return err
	}

	result, err := db.Exec(`
		INSERT INTO games (id, lat, lon, activity, players)
		VALUES ($1, $2, $3, $4, $5)`,
		game.Id, game.Lat, game.Lon, game.Activity, bs,
	)
	if err != nil {
		return err
	}
	if n, err := result.RowsAffected(); err != nil || n != 1 {
		return fmt.Errorf("could not insert game in table: %w", err)
	}

	return nil
}

func Expire() error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	result, err := tx.Query(`
		DELETE FROM ping
		WHERE expire < $1
		RETURNING player, id`,
		time.Now().Unix(),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer result.Close()

	pings := []model.Ping{}

	for result.Next() {
		var playerID, pingID uuid.UUID
		err := result.Scan(&playerID, &pingID)
		if err != nil {
			tx.Rollback()
			return err
		}
		pings = append(pings, model.Ping{
			Player: playerID,
			ID:     pingID,
		})
	}
	if result.Err() != nil {
		tx.Rollback()
		return result.Err()
	}

	for _, ping := range pings {
		_, err := tx.Exec(`
			DELETE FROM player2game
			WHERE ping=$1`,
			ping.ID)
		if err != nil {
			tx.Rollback()
			return err
		}

		// remove player from game and delete game
	}

	tx.Commit()

	return nil
}

func StorePlayerGame(ping model.Ping, game model.Game) {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = tx.Exec(`
			INSERT INTO player2game (player, game, ping)
			VALUES ($1, $2, $3)
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
		var playersBS []byte
		var players []uuid.UUID

		row := db.QueryRow(`SELECT id, lat, lon, activity, players from games WHERE id=$1`, gameID)
		if row.Err() != nil {
			fmt.Println(err)
			continue
		}

		if err := row.Scan(
			&game.Id,
			&game.Lat,
			&game.Lon,
			&game.Activity,
			&playersBS,
		); err != nil {
			fmt.Println(err)
			continue
		}

		err := json.Unmarshal(playersBS, &players)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		game.Players = players

		result = append(result, game)
	}

	return result
}

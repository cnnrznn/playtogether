package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/cnnrznn/playtogether/model"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
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

func UpsertPlayRequest(pr model.PlayRequest) error {
	if _, err := db.Exec(`
		INSERT INTO playrequest (user_id, size, activity, lat, lon, start_time, end_time, range_km)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id)
			DO UPDATE SET size=$2, activity=$3, lat=$4, lon=$5, start_time=$6, end_time=$7, range_km=$8`,
		pr.User, pr.Size, pr.Activity, pr.Lat, pr.Lon, pr.Start, pr.End, pr.RangeKM); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	return nil
}

func LoadPlayRequestUser(userID uuid.UUID) (*model.PlayRequest, error) {
	res := db.QueryRow(`
		SELECT user_id, size, activity, lat, lon, start_time, end_time, range_km FROM playrequest
		WHERE user=$1`,
		userID)

	var pr model.PlayRequest

	if err := res.Scan(
		&pr.User,
		&pr.Size,
		&pr.Activity,
		&pr.Lat,
		&pr.Lon,
		&pr.Start,
		&pr.End,
		&pr.RangeKM,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &pr, nil
}

func LoadPlayRequestArea(pr model.PlayRequest, area model.Area) ([]model.PlayRequest, error) {
	rows, err := db.Query(`
		SELECT user_id, size activity, lat, lon, start_time, end_time, range_km
		FROM playrequest
		WHERE activity=$1 AND
			lat > $2 AND lat < $3 AND lon > $4 AND lon < $5`,
		pr.Activity, area.LatMin, area.LatMax, area.LonMin, area.LonMax,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.PlayRequest

	for rows.Next() {
		var pr model.PlayRequest

		if err := rows.Scan(
			&pr.User, &pr.Size, &pr.Activity, &pr.Lat, &pr.Lon, &pr.Start, &pr.End, &pr.RangeKM,
		); err != nil {
			return nil, err
		}

		result = append(result, pr)
	}

	return result, nil
}

/*
func LoadGamesByArea(ping model.Ping, area model.Area) ([]model.Game, error) {
	games := []model.Game{}

	rows, err := db.Query(`
		SELECT id, lat, lon, activity, players FROM games
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
		playersBS := []byte{}
		players := make(map[uuid.UUID]struct{})

		err := rows.Scan(
			&game.ID,
			&game.Lat,
			&game.Lon,
			&game.Activity,
			&playersBS,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(playersBS, &players); err != nil {
			return nil, err
		}

		game.Players = players
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
	result := db.QueryRow(`
			INSERT INTO ping (id, player, lat, lon, range_km, expire, activity)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT ON CONSTRAINT unq_player_activity
			DO UPDATE SET
				lat=$3, lon=$4, range_km=$5, expire=$6
			RETURNING id`,
		ping.ID, ping.Player, ping.Lat, ping.Lon, ping.RangeKM, ping.Expire, ping.Activity,
	)

	var pingID uuid.UUID

	err := result.Scan(&pingID)
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
		game.ID, game.Lat, game.Lon, game.Activity, bs,
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
		row := tx.QueryRow(`
			DELETE FROM player2game
			WHERE ping=$1 AND player=$2
			RETURNING game`,
			ping.ID, ping.Player)
		var gameID uuid.UUID
		if err := row.Scan(&gameID); err != nil {
			tx.Rollback()
			return err
		}

		if err := RemovePlayerFromGame(tx, ping.Player, gameID); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}

func RemovePlayerFromGame(tx *sql.Tx, playerID, gameID uuid.UUID) error {
	// load game's player list
	row := tx.QueryRow(`
		SELECT players FROM games
		WHERE id=$1`,
		gameID)
	var playerBS []byte
	if err := row.Scan(
		&playerBS,
	); err != nil {
		return err
	}

	players := make(map[uuid.UUID]struct{})
	if err := json.Unmarshal(playerBS, &players); err != nil {
		return err
	}

	delete(players, playerID)

	if len(players) == 0 {
		if _, err := tx.Exec(`
			DELETE FROM games
			WHERE id=$1`, gameID,
		); err != nil {
			return err
		}
		return nil
	}

	playerBS, err := json.Marshal(players)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE games
		SET players=$1
		WHERE id=$2`,
		playerBS, gameID)
	if err != nil {
		return err
	}

	return nil
}

func StorePlayerGame(ping model.Ping, game model.Game) error {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = tx.Exec(`
			INSERT INTO player2game (player, game, ping)
			VALUES ($1, $2, $3)
			ON CONFLICT DO NOTHING`,
		ping.Player, game.ID, ping.ID,
	)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}

	playersBS, err := json.Marshal(game.Players)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		UPDATE games
		SET players=$1
		WHERE id=$2`,
		playersBS, game.ID)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
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
		var players = make(map[uuid.UUID]struct{})

		row := db.QueryRow(`SELECT id, lat, lon, activity, players from games WHERE id=$1`, gameID)
		if row.Err() != nil {
			fmt.Println(err)
			continue
		}

		if err := row.Scan(
			&game.ID,
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
*/

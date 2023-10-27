package db

import (
	"database/sql"
	"encoding/json"
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

func initInternal(
	connString string,
) error {
	if !initDone {

		conn, err := sql.Open(
			"postgres",
			fmt.Sprintf(
				"%v password='%v'", connString, os.Getenv("DB_PASSWD"),
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

func Init() error {
	return initInternal(DB_CONN)
}

func StoreGamePlayer(gameID, playRequestID uuid.UUID, status string) error {
	return nil
}

func LoadGamePlayers(gameID uuid.UUID) (map[uuid.UUID]struct{}, error) {
	return nil, nil
}

func StoreNewGame(game model.Game) error {
	versionID := uuid.New()
	game.ID = uuid.New()
	game.Status = model.CREATED

	blob, err := json.Marshal(game.PlayRequests)
	if err != nil {
		return err
	}

	if _, err := db.Exec(`
		INSERT INTO game (version, id, status, play_requests, activity, lat, lon)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		versionID, game.ID, game.Status, blob, game.Activity, game.Lat, game.Lon,
	); err != nil {
		return err
	}
	return nil
}

func UpdateGame(game model.Game) error {
	blob, err := json.Marshal(game.PlayRequests)
	if err != nil {
		return err
	}

	newVersion := uuid.New()

	row := db.QueryRow(`
		UPDATE game
		SET version=$1, status=$2, play_requests=$3
		WHERE id=$4 AND version=$5
		RETURNING id`,
		newVersion, game.Status, blob,
		game.ID, game.Version,
	)

	return row.Err()
}

func LoadGame(id uuid.UUID) (*model.Game, error) {
	res := db.QueryRow(`
		SELECT version, id, status, play_requests, activity, lat, lon
		FROM game
		WHERE id=$1`,
		id,
	)

	var game model.Game
	var playRequestBS []byte

	if err := res.Scan(
		&game.Version,
		&game.ID,
		&game.Status,
		&playRequestBS,
		&game.Activity,
		&game.Lat,
		&game.Lon,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(playRequestBS, &game.PlayRequests); err != nil {
		return nil, err
	}

	return &game, nil
}

func UpsertPlayRequest(pr model.PlayRequest) error {
	id := uuid.New()

	if _, err := db.Exec(`
		INSERT INTO playrequest (id, user_id, size, activity, lat, lon, start_time, end_time, range_km)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT ON CONSTRAINT unq_player_activity
			DO UPDATE SET id=$1, size=$3, lat=$5, lon=$6, start_time=$7, end_time=$8, range_km=$9`,
		id, pr.User, pr.Size, pr.Activity, pr.Lat, pr.Lon, pr.Start, pr.End, pr.RangeKM); err != nil {
		return err
	}
	return nil
}

func LoadPlayRequestUserActivity(userID uuid.UUID, activity string) (*model.PlayRequest, error) {
	res := db.QueryRow(`
		SELECT id, user_id, size, activity, lat, lon, start_time, end_time, range_km
		FROM playrequest
		WHERE user_id=$1 AND activity=$2`,
		userID, activity)

	var pr model.PlayRequest

	if err := res.Scan(
		&pr.ID,
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
		SELECT id, user_id, size, activity, lat, lon, start_time, end_time, range_km
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
			&pr.ID,
			&pr.User,
			&pr.Size,
			&pr.Activity,
			&pr.Lat,
			&pr.Lon,
			&pr.Start,
			&pr.End,
			&pr.RangeKM,
		); err != nil {
			return nil, err
		}

		result = append(result, pr)
	}

	return result, nil
}

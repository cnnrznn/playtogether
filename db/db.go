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

func GetGames(ping model.Ping, area model.Area) []model.Game {
	return nil
}

func NewPlayer(model.Ping, model.Area) *model.Game {
	return nil
}

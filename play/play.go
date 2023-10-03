package play

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/cnnrznn/playtogether/model"
)

const (
	DB_CONN = "hostname=localhost user=pt database=playtogether sslmode=disable"
)

type PlayService struct {
	db       *sql.DB
	initDone bool
}

func Update(ping model.Ping) error {
	// First, check for games already going on in the area

	// If no games found, put player into players DB and try to create a game with the new player information

	// If game created, put game into games table and send alerts to players

	return nil
}

func (s *PlayService) Init() error {
	if !s.initDone {
		conn, err := sql.Open(
			"postgres",
			fmt.Sprintf(
				"%v password='%v'", DB_CONN, os.Getenv("DB_PASSWD"),
			),
		)
		if err != nil {
			return err
		}

		s.db = conn
	}

	return nil
}

func (s *PlayService) Run() error {
	// every 5m, scan and delete rows in db past expiration
	if err := s.Init(); err != nil {
		return err
	}

	go s.runExpire()

	return nil
}

func (s *PlayService) runExpire() error {
	return nil
}

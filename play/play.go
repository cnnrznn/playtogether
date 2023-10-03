package play

import (
	"fmt"

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

func Update(ping model.Ping) (*Response, error) {
	// First, check for games already going on in the area

	// If no games found, put player into players DB and try to create a game with the new player information

	// If game created, put game into games table and send alerts to players

	return nil, fmt.Errorf("not implemented")
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

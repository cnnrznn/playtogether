package game

import (
	"github.com/cnnrznn/playtogether/db"
	"github.com/cnnrznn/playtogether/model"
	"github.com/google/uuid"
)

func Create(game model.Game) error {
	game.ID = uuid.New()
	return db.StoreGame(game)
}

func Confirm(gameID, prID uuid.UUID) error {
	game, err := db.LoadGame(gameID)
	if err != nil {
		return err
	}

	game.PlayRequests[prID] = struct{}{}

	if err := db.UpdateGame(*game); err != nil {
		return err
	}
	return nil
}

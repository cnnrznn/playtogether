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
	// Store a confirmation mapping
	// gameID --> prID
	// in new table

	// A game is "pending" when everyone on the pr list confirms
	return nil
}

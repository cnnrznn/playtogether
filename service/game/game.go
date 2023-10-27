package game

import (
	"github.com/cnnrznn/playtogether/db"
	"github.com/cnnrznn/playtogether/model"
	"github.com/google/uuid"
)

func Create(game model.Game) error {
	game.ID = uuid.New()
	game.Status = model.CREATED
	return db.StoreGame(game)
}

func Confirm(gameID, prID uuid.UUID) error {
	// Store a confirmation mapping
	// gameID --> prID
	// in new table
	if err := db.StoreGamePlayer(gameID, prID, "confirmed"); err != nil {
		return err
	}

	game, err := db.LoadGame(gameID)
	if err != nil {
		return err
	}

	prs, err := db.LoadGamePlayers(gameID)
	if err != nil {
		return err
	}

	// check if prs is the same set as game
	// if yes, tell db to mark game and players as pending for the game
	confirmed := true
	for prID := range game.PlayRequests {
		if _, ok := prs[prID]; !ok {
			confirmed = false
			break
		}
	}

	if confirmed {
		game.Status = model.PENDING
	}

	if err := db.StoreGame(*game); err != nil {
		return err
	}

	return nil
}

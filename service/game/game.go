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

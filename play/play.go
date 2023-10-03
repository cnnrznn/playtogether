package play

import "github.com/cnnrznn/playtogether/model"

func Update(ping model.Ping) error {
	// First, check for games already going on in the area

	// If no games found, put player into players DB and try to create a game with the new player information

	// If game created, put game into games table and send alerts to players

	return nil
}

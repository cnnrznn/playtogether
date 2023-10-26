package api

import "github.com/google/uuid"

type GameRequest struct {
	Action         string      `json:"action"`
	PlayRequestIDs []uuid.UUID `json:"play_request_ids"`
	GameID         string      `json:"game_id"`
}

package api

import "github.com/google/uuid"

type GameRequest struct {
	Action         string      `json:"action"`
	PlayRequestIDs []uuid.UUID `json:"play_request_ids"`
	PlayRequestID  uuid.UUID   `json:"play_request_id"`
	GameID         uuid.UUID   `json:"game_id"`
}

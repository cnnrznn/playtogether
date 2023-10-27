// Package play handles play requests and the filtering logic to find other play requests
package play

import (
	"github.com/cnnrznn/playtogether/db"
	"github.com/cnnrznn/playtogether/model"
	"github.com/cnnrznn/playtogether/service/area"
	"github.com/google/uuid"
)

// CreatePlayRequest creates a new request or updates an existing request from the same user
func CreatePlayRequest(pr model.PlayRequest) error {
	return db.UpsertPlayRequest(pr)
}

func GetPlayRequests(
	user uuid.UUID,
	activity string,
) ([]model.PlayRequest, error) {
	// Get user's play request from DB
	userPR, err := db.LoadPlayRequestUserActivity(user, activity)
	if err != nil {
		return nil, err
	}

	if userPR == nil {
		return nil, nil
	}

	// Query for all play requests that match activity + range
	a := area.Calculate(userPR.Lat, userPR.Lon, userPR.RangeKM)
	prs, err := db.LoadPlayRequestArea(*userPR, a)
	if err != nil {
		return nil, err
	}

	// Return list of self + matches to user
	result := []model.PlayRequest{*userPR}
	for _, pr := range prs {
		if pr.User != userPR.User {
			result = append(result, pr)
		}
	}

	// Frontend should display a timeline for each match and quantity, distance
	return result, nil
}

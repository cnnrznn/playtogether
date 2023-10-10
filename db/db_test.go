package db

import (
	"testing"
	"time"

	"github.com/cnnrznn/playtogether/model"
	"github.com/google/uuid"
)

func TestDuplicatePlayerActivity(t *testing.T) {
	Init()

	player := uuid.New()

	ping := model.Ping{
		Player:   player,
		Activity: "volleyball",
		Lat:      0.00,
		Lon:      0.00,
		Expire:   int(time.Now().Unix()),
		RangeKM:  20,
	}

	err := StorePing(ping)
	if err != nil {
		t.Error(err)
	}

	err = StorePing(ping)
	if err != nil {
		t.Error(err)
	}
}

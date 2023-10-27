package db

import (
	"reflect"
	"testing"

	"github.com/cnnrznn/playtogether/model"
	"github.com/go-test/deep"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	TEST_DB_CONN = "host=localhost user=pt database=playtogether sslmode=disable"
)

func test_before() error {
	return initInternal(TEST_DB_CONN)
}

func test_after() error {
	// offer a routine here for dropping all table data
	return nil
}

func TestStoreAndLoadGame(t *testing.T) {
	if err := test_before(); err != nil {
		t.Error(err)
		return
	}
	defer test_after()

	playRequests := map[uuid.UUID]struct{}{
		uuid.New(): {},
		uuid.New(): {},
		uuid.New(): {},
	}

	game := model.Game{
		ID:           uuid.New(),
		PlayRequests: playRequests,
		Status:       model.CREATED,
		Activity:     "test_volleyball",
		Lat:          -1.00,
		Lon:          1.00,
	}

	if err := StoreNewGame(game); err != nil {
		t.Error(err)
		return
	}

	gameLoad, err := LoadGame(game.ID)
	if err != nil {
		t.Error(err)
		return
	}

	game.Version = gameLoad.Version // version is internally set by the storage layer
	//bs, _ := json.MarshalIndent(game, "-", "  ")
	//bs2, _ := json.MarshalIndent(gameLoad, "-", "  ")
	//fmt.Println(string(bs))
	//fmt.Println(string(bs2))

	if diff := deep.Equal(game, *gameLoad); diff != nil {
		t.Error(diff)
	}

	assert.True(t, reflect.DeepEqual(game, *gameLoad))
}

func TestUpsertPlayRequest(t *testing.T) {
	test_before()
	defer test_after()

	playRequest := model.PlayRequest{
		User:     uuid.New(),
		Size:     3,
		Activity: "test_volleyball",
		Lat:      -1.00,
		Lon:      1.00,
		Start:    1000,
		End:      1100,
		RangeKM:  20,
	}

	if err := UpsertPlayRequest(playRequest); err != nil {
		t.Error(err)
		return
	}

	prLoad, err := LoadPlayRequestUserActivity(playRequest.User, playRequest.Activity)
	if err != nil {
		t.Error(err)
		return
	}

	playRequest.ID = prLoad.ID

	assert.True(t, reflect.DeepEqual(*prLoad, playRequest))

	playRequest.Size = 10
	oldID := prLoad.ID

	if err := UpsertPlayRequest(playRequest); err != nil {
		t.Error(err)
		return
	}

	prLoad2, err := LoadPlayRequestUserActivity(playRequest.User, playRequest.Activity)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEqual(t, oldID, prLoad2.ID)
	assert.Equal(t, 10, prLoad2.Size)

}

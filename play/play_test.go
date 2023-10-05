package play

import (
	"fmt"
	"testing"
	"time"

	"github.com/cnnrznn/playtogether/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAreaCalculation(t *testing.T) {
	ping := model.Ping{
		Lat:     0,
		Lon:     0,
		RangeKM: 5,
	}

	area := calculateArea(ping)
	fmt.Println(area)

	pingCO := model.Ping{
		Lat:     40.0150,
		Lon:     105.2705,
		RangeKM: 5,
	}

	areaCO := calculateArea(pingCO)
	fmt.Println(areaCO)
}

func TestPing(t *testing.T) {
	ping := model.Ping{
		Player:   uuid.New(),
		Activity: "volleyball",
		Lat:      40.0150,
		Lon:      105.2705,
		Expire:   int(time.Now().Unix()),
		RangeKM:  10,
	}

	res, err := Update(ping)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, res.Found, false)
}

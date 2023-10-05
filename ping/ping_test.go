package ping

import (
	"fmt"
	"testing"

	"github.com/cnnrznn/playtogether/model"
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

// Package area contains utility functions for range and search calculations
package area

import (
	"math"

	"github.com/cnnrznn/playtogether/model"
	"github.com/jftuga/geodist"
)

// Utility function for area calculation

func Calculate(lat, lon float64, rangeKM int) model.Area {
	latDelta := calculateDegree(lat, lon, rangeKM, true)
	lonDelta := calculateDegree(lat, lon, rangeKM, false)

	return model.Area{
		LatMin: lat - latDelta,
		LatMax: lat + latDelta,
		LonMin: lon - lonDelta,
		LonMax: lon + lonDelta,
	}
}

func calculateDegree(lat, lon float64, rangeKM int, latitude bool) float64 {
	const tolerance float64 = 0.001
	var target float64 = float64(rangeKM)

	var minDeg float64 = 0.00
	var maxDeg float64 = 1.00
	var midDeg float64
	var midCoord geodist.Coord

	// iteratively binary search to find degree in tolerance
	for {
		midDeg = (minDeg + maxDeg) / 2
		switch latitude {
		case true:
			midCoord = geodist.Coord{
				Lat: lat + midDeg,
				Lon: lon,
			}
		default:
			midCoord = geodist.Coord{
				Lat: lat,
				Lon: lon + midDeg,
			}
		}

		// calculate distance to mid
		_, km, _ := geodist.VincentyDistance(
			geodist.Coord{
				Lat: lat,
				Lon: lon,
			},
			midCoord,
		)

		if math.Abs(target-km) < tolerance {
			break
		}

		if target-km < 0 { // point is farther away than target
			maxDeg = midDeg
		} else {
			minDeg = midDeg
		}
	}

	return midDeg
}

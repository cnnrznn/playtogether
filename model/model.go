package model

type Ping struct {
	Activity string  `json:"activity"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Expire   int     `json:"expire"`
	RangeKM  int     `json:"range_km"`

	// other stuff like age, intensity
}

type Game struct {
	Activity string   `json:"activity"`
	Lat      float64  `json:"lat"`
	Lon      float64  `json:"lon"`
	Players  []string `json:"players"`
}

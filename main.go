package main

import (
	"fmt"
	"os"

	"github.com/cnnrznn/playtogether/api"
	"github.com/cnnrznn/playtogether/play"
)

func main() {
	ps := &play.PlayService{}
	go ps.Run()

	err := api.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

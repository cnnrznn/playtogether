package main

import (
	"fmt"
	"os"

	"github.com/cnnrznn/playtogether/api"
)

func main() {
	err := api.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

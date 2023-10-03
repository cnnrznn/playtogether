package main

import (
	"log"
	"sync"

	"github.com/cnnrznn/playtogether/api"
	"github.com/cnnrznn/playtogether/play"
)

func main() {
	errs := make(chan error)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		errs <- play.Run()
		wg.Done()
	}()

	go func() {
		errs <- api.Run()
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(errs)
	}()

	for err := range errs {
		if err != nil {
			log.Fatal(err)
		}
	}
}

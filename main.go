package main

import (
	"log"
	"sync"

	"github.com/cnnrznn/playtogether/api"
)

func main() {
	// TODO use a proper logging library

	errs := make(chan error)

	wg := &sync.WaitGroup{}
	wg.Add(1)

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

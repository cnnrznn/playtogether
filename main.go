package main

import (
	"log"
	"sync"

	"github.com/cnnrznn/playtogether/api"
	"github.com/cnnrznn/playtogether/service"
)

func main() {
	// TODO use a proper logging library

	errs := make(chan error, 1)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	errs <- service.Init()

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

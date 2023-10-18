package service

import (
	"github.com/cnnrznn/playtogether/db"
)

var (
	initDone bool = false
)

func Init() error {
	if !initDone {
		err := db.Init()
		if err != nil {
			return err
		}

		initDone = true
	}

	return nil
}

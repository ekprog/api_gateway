package job

import (
	"go.uber.org/dig"
	"microservice/app/core"
)

var log core.Logger
var di *dig.Container

type Job interface {
	Run() error
}

func Init(logger core.Logger, di_ *dig.Container) error {
	log = logger
	di = di_
	return nil
}

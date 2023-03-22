package core

import "go.uber.org/dig"

var di *dig.Container

func init() {
	di = dig.New()
}

func GetDI() *dig.Container {
	return di
}

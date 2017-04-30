package svc

import "github.com/drgomesp/cargo/container"

const (
	DI_LOGGER = "logger"
	DI_CONFIG = "config"
)

var Container *container.Container

func init() {
	Container = container.New()

	Container.Set(DI_LOGGER, NewLiveLogger())
}

package empty

import (
	"helloworld/shared/driver"
	"helloworld/shared/infrastructure/config"
	"helloworld/shared/infrastructure/logger"
)

type gateway struct {
	log     logger.Logger
	appData driver.ApplicationData
	config  *config.Config
}

// NewGateway ...
func NewGateway(log logger.Logger, appData driver.ApplicationData, config *config.Config) *gateway {

	return &gateway{
		log:     log,
		appData: appData,
		config:  config,
	}
}

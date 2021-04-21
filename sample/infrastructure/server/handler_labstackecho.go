package server

import (
	"github.com/labstack/echo/v4"
)

// LabstackEchoHandler will define basic HTTP configuration with gracefully shutdown
type LabstackEchoHandler struct {
	GracefullyShutdown
	Router *echo.Echo
}

func NewLabstackEchoHandler(address string) LabstackEchoHandler {

	router := echo.New()

	return LabstackEchoHandler{
		GracefullyShutdown: NewGracefullyShutdown(router, address),
		Router:             router,
	}

}

// RunApplication is implementation of RegistryContract.RunApplication()
func (r *LabstackEchoHandler) RunApplication() {
	r.RunWithGracefullyShutdown()
}

package server

import (
	"net/http"
)

// NetHTTPHandler will define basic HTTP configuration with gracefully shutdown
type NetHTTPHandler struct {
	GracefullyShutdown
	Router *http.ServeMux
}

func NewNetHTTPHandler(address string) NetHTTPHandler {

	router := http.NewServeMux()

	return NetHTTPHandler{
		GracefullyShutdown: NewGracefullyShutdown(router, address),
		Router:             router,
	}

}

// RunApplication is implementation of RegistryContract.RunApplication()
func (r *NetHTTPHandler) RunApplication() {
	r.RunWithGracefullyShutdown()
}

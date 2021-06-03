package server

import (
	"github.com/go-chi/chi/v5"
)

// GoChiHandler will define basic HTTP configuration with gracefully shutdown
type GoChiHandler struct {
	GracefullyShutdown
	Router *chi.Mux
}

func NewGoChiHandler(address string) GoChiHandler {

	router := chi.NewMux()

	return GoChiHandler{
		GracefullyShutdown: NewGracefullyShutdown(router, address),
		Router:             router,
	}

}

// RunApplication is implementation of RegistryContract.RunApplication()
func (r *GoChiHandler) RunApplication() {
	r.RunWithGracefullyShutdown()
}

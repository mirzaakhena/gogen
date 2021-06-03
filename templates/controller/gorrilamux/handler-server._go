package server

import (
	"github.com/gorilla/mux"
)

// GorrilaMuxHandler will define basic HTTP configuration with gracefully shutdown
type GorrilaMuxHandler struct {
	GracefullyShutdown
	Router *mux.Router
}

func NewGorrilaMuxHandler(address string) GorrilaMuxHandler {

	router := mux.NewRouter()

	return GorrilaMuxHandler{
		GracefullyShutdown: NewGracefullyShutdown(router, address),
		Router:             router,
	}

}

// RunApplication is implementation of RegistryContract.RunApplication()
func (r *GorrilaMuxHandler) RunApplication() {
	r.RunWithGracefullyShutdown()
}

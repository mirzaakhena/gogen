package server

import (
	"github.com/julienschmidt/httprouter"
)

// JulienSchmidtHandler will define basic HTTP configuration with gracefully shutdown
type JulienSchmidtHandler struct {
	GracefullyShutdown
	Router *httprouter.Router
}

func NewJulienSchmidtHandler(address string) JulienSchmidtHandler {

	router := httprouter.New()

	return JulienSchmidtHandler{
		GracefullyShutdown: NewGracefullyShutdown(router, address),
		Router:             router,
	}

}

// RunApplication is implementation of RegistryContract.RunApplication()
func (r *JulienSchmidtHandler) RunApplication() {
	r.RunWithGracefullyShutdown()
}

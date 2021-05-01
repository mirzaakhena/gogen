package server

import (
	"accounting/infrastructure/log"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// GracefullyShutdown will handle http server with gracefully shutdown mechanism
type GracefullyShutdown struct {
	httpServer *http.Server
}

func NewGracefullyShutdown(handler http.Handler, address string) GracefullyShutdown {
	return GracefullyShutdown{
		httpServer: &http.Server{
			Addr:    address,
			Handler: handler,
		},
	}
}

// RunWithGracefullyShutdown is ...
func (r *GracefullyShutdown) RunWithGracefullyShutdown() {

	ctx := log.Context(context.Background(), "RunWithGracefullyShutdown")

	go func() {
		if err := r.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(ctx, "listen: %s", err)
			os.Exit(1)
		}
	}()

	log.Info(ctx, "server is running in %s", r.httpServer.Addr)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info(ctx, "Shutting down server")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := r.httpServer.Shutdown(ctx); err != nil {
		log.Error(ctx, "Server forced to shutdown: %s", err.Error())
		os.Exit(1)
	}

	log.Info(ctx, "Server stoped")

}

package server

import (
	"accounting/infrastructure/log2"
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

	ctx := log2.Context(context.Background(), "RunWithGracefullyShutdown")

	go func() {
		if err := r.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log2.Error(ctx, "listen: %s", err)
			os.Exit(1)
		}
	}()

	log2.Info(ctx, "server is running ...")

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log2.Info(ctx, "Shutting down server...")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := r.httpServer.Shutdown(ctx); err != nil {
		log2.Error(ctx, "Server forced to shutdown:", err)
		os.Exit(1)
	}

	log2.Info(ctx, "Server stoped.")

}

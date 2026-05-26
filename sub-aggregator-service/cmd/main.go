package main

import (
	"context"
	alogger "github.com/AndreSS-ntp/logger"
	"github.com/AndreSS-ntp/subscription-aggregator/sub-aggregator-service/internal/config"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := alogger.NewLogger()
	baseCtx := alogger.WithLogger(context.Background(), logger)
	logger.Info(baseCtx, "Service is running...")
	ctx, cancel := withGracefulShutdown(baseCtx)
	defer cancel()

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:        config.IP_port,
		Handler:     mux,
		BaseContext: func(net.Listener) context.Context { return ctx },
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Error(ctx, "HTTP server error: "+err.Error())
			return
		}
	}()

	// sub_aggregator_service.Testing()

	<-ctx.Done()

	err := server.Shutdown(context.Background())
	if err != nil {
		logger.Error(ctx, "error shutting down server: "+err.Error())
	}

	logger.Info(ctx, "Service stopped.")
}

func withGracefulShutdown(baseCtx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(baseCtx)
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-exit
		alogger.FromContext(ctx).Warn(ctx, "Shitting down service...")
		cancel()
	}()
	return ctx, cancel
}

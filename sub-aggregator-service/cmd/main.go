package main

import (
	"context"
	alogger "github.com/AndreSS-ntp/logger"
	"github.com/AndreSS-ntp/subscription-aggregator/sub-aggregator-service/internal/app"
	"github.com/AndreSS-ntp/subscription-aggregator/sub-aggregator-service/internal/config"
	"github.com/AndreSS-ntp/subscription-aggregator/sub-aggregator-service/internal/repository"
	"github.com/AndreSS-ntp/subscription-aggregator/sub-aggregator-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
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

	dbpool, err := pgxpool.New(ctx, config.DB_URL)
	if err != nil {
		logger.Error(ctx, "unable to create connection pool: "+err.Error())
		return
	}
	defer dbpool.Close()

	database := repository.NewRepository(dbpool)
	serv := service.NewService(database)
	application := app.NewApp(serv)
	mux := http.NewServeMux()

	for pattern, command := range application.Commands {
		mux.HandleFunc(pattern, alogger.HandlerWithLogger(logger, command.Handler))
	}

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

	<-ctx.Done()

	err = server.Shutdown(context.Background())
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

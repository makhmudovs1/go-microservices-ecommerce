package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/server"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	dsn := os.Getenv("POSTGRES_DSN")
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		panic(err)
	}
	srv, _ := server.New(pool, logger)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Run server on goroutine
	go func() {
		logger.Info("Starting server", zap.String("addr", ":8080"))
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error: %v", zap.Error(err))
		}
	}()
	<-quit
	logger.Info("Shutdown signal received")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Stop(ctx); err != nil {
		logger.Error("Server shutdown error", zap.Error(err))
	} else {
		logger.Info("Server gracefully shutdown")
	}
	pool.Close()
	logger.Info("Pool closed")
}

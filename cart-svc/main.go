package main

import (
	"context"
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/repository"
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/server"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := repository.InitPostgres(ctx); err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	// creating http.Server with any dependencies
	srv := server.New()

	go func() {
		if err := srv.Run(); err != nil {
			log.Printf("server error: %v", err)
			stop()
		}
	}()
	<-ctx.Done()
	log.Println("shutting down gracefully...")
	srv.Stop()
}

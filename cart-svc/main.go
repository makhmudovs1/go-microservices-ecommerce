package main

import (
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/server"
	"log"
	"net/http"
	"os"
)

func main() {
	dsn := os.Getenv("POSTGRES_DSN")
	srv, err := server.New(dsn)
	if err != nil {
		log.Fatalf("cannot create server: %v", err)
	}

	log.Println("Starting server on :8080")
	if err := srv.Run(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

package server

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/handler"
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/handler/cart"
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/repository"
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/service"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func New(pool *pgxpool.Pool) (*Server, error) {
	// 1) connecting to DB
	dbPool := pool
	// 2) repository
	repo := repository.NewCartRepository(dbPool)

	// 3) service
	svc := service.NewCartService(repo)

	// 4) HTTP-handler
	h := cart.NewHandler(svc)

	// 5) Set up routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/cart", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.AddToCart(w, r)
		case http.MethodGet:
			h.GetCart(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	httpSrv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	return &Server{
		httpServer: httpSrv,
	}, nil
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

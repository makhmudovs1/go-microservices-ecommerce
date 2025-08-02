package server

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/handler/cart"
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/repository"
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/service"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func New(pool *pgxpool.Pool, logger *zap.Logger) (*Server, error) {
	// 1) connecting to DB
	dbPool := pool
	// 2) repository
	repo := repository.NewCartRepository(dbPool)

	// 3) service
	svc := service.NewCartService(repo, logger)

	// 4) HTTP-handler
	h := cart.NewHandler(svc, logger)

	// 5) Set up routes
	mux := http.NewServeMux()
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

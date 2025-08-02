package server

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/handler/cart"
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/repository"
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/service"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Info("http request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
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
	handler := LoggingMiddleware(logger)(mux)
	httpSrv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
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

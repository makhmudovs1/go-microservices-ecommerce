package server

import (
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/handler"
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/handler/cart"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func New() *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/cart", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			cart.AddToCart(w, r)
		case http.MethodGet:
			cart.GetCart(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	return &Server{
		httpServer: srv,
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.httpServer.Close()
}

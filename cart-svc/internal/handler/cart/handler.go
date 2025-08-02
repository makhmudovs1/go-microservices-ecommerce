package cart

import (
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/service"
	"go.uber.org/zap"
)

type Handler struct {
	svc    service.CartService
	logger *zap.Logger
}

func NewHandler(svc service.CartService, logger *zap.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
	}
}

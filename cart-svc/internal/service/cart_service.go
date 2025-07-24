package service

import (
	"context"
	"fmt"
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/models"
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/repository"
)

type CartService interface {
	AddItem(ctx context.Context, userID string, sku string, qty int) (models.Cart, error)
	GetCartByUserID(ctx context.Context, userID string) (models.Cart, error)
}

type cartService struct {
	repo repository.CartRepository
}

func NewCartService(repo repository.CartRepository) CartService {
	return &cartService{
		repo: repo,
	}
}

func (s *cartService) AddItem(ctx context.Context, userID string, sku string, qty int) (models.Cart, error) {
	if userID == "" {
		return models.Cart{}, fmt.Errorf("userID is required")
	}
	if sku == "" {
		return models.Cart{}, fmt.Errorf("sku is required")
	}
	if qty <= 0 {
		return models.Cart{}, fmt.Errorf("qty must be > 0")
	}
	err := s.repo.AddItemToCart(ctx, userID, sku, qty)
	if err != nil {
		return models.Cart{}, err
	}
	return s.repo.GetCartByUserID(ctx, userID)
}

func (s *cartService) GetCartByUserID(ctx context.Context, userID string) (models.Cart, error) {
	return s.repo.GetCartByUserID(ctx, userID)
}

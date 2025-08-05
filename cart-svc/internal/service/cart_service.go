package service

import (
	"context"
	"fmt"
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/models"
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/repository"
	"go.uber.org/zap"
)

type CartService interface {
	AddItem(ctx context.Context, userID string, sku string, qty int) (models.Cart, error)
	GetCartByUserID(ctx context.Context, userID string) (models.Cart, error)
	RemoveItem(ctx context.Context, userID, sku string) (models.Cart, error)
}

type cartService struct {
	repo   repository.CartRepository
	logger *zap.Logger
}

func NewCartService(repo repository.CartRepository, logger *zap.Logger) CartService {
	return &cartService{
		repo:   repo,
		logger: logger,
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

func (s *cartService) RemoveItem(ctx context.Context, userID, sku string) (models.Cart, error) {
	if userID == "" {
		return models.Cart{}, fmt.Errorf("userID is required")
	}
	if sku == "" {
		return models.Cart{}, fmt.Errorf("sku is required")
	}
	if err := s.repo.RemoveItemFromCart(ctx, userID, sku); err != nil {
		return models.Cart{}, err
	}
	return s.repo.GetCartByUserID(ctx, userID)
}

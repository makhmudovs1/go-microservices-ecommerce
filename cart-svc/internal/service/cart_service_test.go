package service

import (
	"context"
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/models"
	"testing"
)

type args struct {
	userID string
	sku    string
	qty    int
}

type mockCartRepository struct {
	addItemErr    error
	getCartResult models.Cart
	getCartErr    error
}

func (m *mockCartRepository) AddItemToCart(ctx context.Context, userID, sku string, qty int) error {
	return m.addItemErr
}

func (m *mockCartRepository) GetCartByUserID(ctx context.Context, userID string) (models.Cart, error) {
	return m.getCartResult, m.getCartErr
}

func (m *mockCartRepository) RemoveItemFromCart(ctx context.Context, userID, sku string) error {
	return m.getCartErr
}

func TestCartService_AddItem(t *testing.T) {
	testCases := []struct {
		name       string
		repo       *mockCartRepository
		args       args
		wantErr    bool
		wantResult models.Cart
	}{
		{
			name: "success",
			repo: &mockCartRepository{
				getCartResult: models.Cart{
					UserId: "testuser",
					Items:  []models.CartItem{{SKU: "item-1", Qty: 2}},
				},
			},
			wantErr:    false,
			wantResult: models.Cart{UserId: "testuser", Items: []models.CartItem{{SKU: "item-1", Qty: 2}}},
		},
		{
			name: "repo addItem error",
			repo: &mockCartRepository{
				addItemErr: context.DeadlineExceeded,
			},
			wantErr:    true,
			wantResult: models.Cart{},
		},
		{
			name: "repo getCartErr",
			repo: &mockCartRepository{
				getCartErr: context.DeadlineExceeded,
			},
			wantErr:    true,
			wantResult: models.Cart{},
		},
		{
			name:    "userID is empty",
			repo:    &mockCartRepository{},
			args:    args{userID: "", sku: "item-1", qty: 2},
			wantErr: true,
		},
		{
			name:    "SKU is empty",
			repo:    &mockCartRepository{},
			args:    args{userID: "user", sku: "", qty: 2},
			wantErr: true,
		},
		{
			name:    "qty is zero",
			repo:    &mockCartRepository{},
			args:    args{userID: "user", sku: "item-1", qty: 0},
			wantErr: true,
		},
		{
			name:    "qty is negative",
			repo:    &mockCartRepository{},
			args:    args{userID: "user", sku: "item-1", qty: -5},
			wantErr: true,
		},
		{
			name: "success",
			repo: &mockCartRepository{
				getCartResult: models.Cart{
					UserId: "testuser",
					Items:  []models.CartItem{{SKU: "item-1", Qty: 2}},
				},
			},
			args:       args{userID: "testuser", sku: "item-1", qty: 2},
			wantErr:    false,
			wantResult: models.Cart{UserId: "testuser", Items: []models.CartItem{{SKU: "item-1", Qty: 2}}},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			svc := &cartService{repo: tt.repo}
			result, err := svc.AddItem(context.Background(), tt.args.userID, tt.args.sku, tt.args.qty)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				if result.UserId != tt.wantResult.UserId || len(result.Items) != len(tt.wantResult.Items) {
					t.Fatalf("unexpected result: %+v", result)
				}
			}
		})
	}
}

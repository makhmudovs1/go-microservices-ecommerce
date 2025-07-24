package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/makhmudovs1/go-microservices-ecommerce/cart-svc/internal/models"
)

type CartRepository interface {
	AddItemToCart(ctx context.Context, userID, sku string, qty int) error
	GetCartByUserID(ctx context.Context, userID string) (models.Cart, error)
}

type cartRepository struct {
	db *pgxpool.Pool
}

func NewCartRepository(db *pgxpool.Pool) CartRepository {
	return &cartRepository{
		db: db,
	}
}

// AddItemToCart

func (r *cartRepository) AddItemToCart(ctx context.Context, userID, sku string, qty int) error {
	var cartID int64

	err := r.db.QueryRow(ctx, "SELECT id FROM cart WHERE user_id = $1", userID).Scan(&cartID)
	if err == pgx.ErrNoRows {
		// there are no carts
		err = r.db.QueryRow(ctx, "INSERT INTO cart (user_id) VALUES ($1) RETURNING id", userID).Scan(&cartID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx,
		`INSERT INTO cart_item (cart_id, sku, qty)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (cart_id, sku)
		 DO UPDATE SET qty = cart_item.qty + EXCLUDED.qty`,
		cartID, sku, qty,
	)
	return err
}

// GetCartByUserId

func (r *cartRepository) GetCartByUserID(ctx context.Context, userID string) (models.Cart, error) {
	var cartID int64
	cart := models.Cart{
		UserId: userID,
		Items:  []models.CartItem{},
	}
	err := r.db.QueryRow(ctx, "SELECT id FROM cart WHERE user_id = $1", userID).Scan(&cartID)
	if err == pgx.ErrNoRows {
		return models.Cart{}, err
	}
	rows, err := r.db.Query(ctx,
		"SELECT sku, qty FROM cart_item WHERE cart_id = $1", cartID)
	defer rows.Close()
	for rows.Next() {
		var item models.CartItem
		err = rows.Scan(&item.SKU, &item.Qty)
		if err != nil {
			return models.Cart{}, err
		}
		cart.Items = append(cart.Items, item)
	}
	if err := rows.Err(); err != nil {
		return models.Cart{}, err
	}
	return cart, err
}

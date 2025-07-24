package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	dsn := os.Getenv("POSTGRES_TEST_DSN")
	var err error
	pool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("failed to connect to test db: %v", err)
	}
	code := m.Run()
	pool.Close()
	os.Exit(code)
}

func setupTestDB(ctx context.Context, t *testing.T) CartRepository {
	_, err := pool.Exec(ctx, "TRUNCATE cart_item, cart RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}
	return NewCartRepository(pool)
}

func TestConnectToTestDB(t *testing.T) {
	dsn := os.Getenv("POSTGRES_TEST_DSN")
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	err = pool.Ping(context.Background())
	if err != nil {
		t.Fatalf("failed to ping: %v", err)
	}
}

func TestAddItemToCart(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(ctx, t)
	var (
		userID = "user1"
		sku    = "item1"
		qty    = -1
	)
	err := db.AddItemToCart(ctx, userID, sku, qty)
	if err == nil {
		t.Fatalf("expected error when adding qty < 1, got nil")
	}
}

func TestAddItemToCart_Positive(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(ctx, t)

	userID := "user1"
	sku := "item1"
	qty := 2

	err := db.AddItemToCart(ctx, userID, sku, qty)
	if err != nil {
		t.Fatalf("unexpected error when adding item: %v", err)
	}
	cart, err := db.GetCartByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("failed to get cart: %v", err)
	}
	if len(cart.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(cart.Items))
	}
	if cart.Items[0].SKU != sku || cart.Items[0].Qty != qty {
		t.Fatalf("unexpected item data: %+v", cart.Items[0])
	}
}

func TestAddItemToCart_IncrementQty(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(ctx, t)

	userID := "user1"
	sku := "item1"
	qty := 2

	_ = db.AddItemToCart(ctx, userID, sku, qty)
	_ = db.AddItemToCart(ctx, userID, sku, qty)

	cart, err := db.GetCartByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("failed to get cart: %v", err)
	}
	if len(cart.Items) != 1 {
		t.Fatalf("expected 1 item in cart, got %d", len(cart.Items))
	}
	if cart.Items[0].Qty != 4 {
		t.Fatalf("expected qty=2, got %d", cart.Items[0].Qty)
	}
}

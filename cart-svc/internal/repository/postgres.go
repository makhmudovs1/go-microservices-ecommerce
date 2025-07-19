package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
)

var pool *pgxpool.Pool

func InitPostgres(ctx context.Context) error {
	dsn := os.Getenv("POSTGRES_DSN")
	var err error
	pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		return err
	}
	// test of connect
	if err := pool.Ping(ctx); err != nil {
		return err
	}
	log.Println("Connected to PostgreSQL")
	return nil
}

func GetPool() *pgxpool.Pool {
	return pool
}

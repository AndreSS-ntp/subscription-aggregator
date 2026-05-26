package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db}
}

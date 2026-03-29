package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"hello-world-go/internal/model"
)

type ItemRepository struct {
	db *pgxpool.Pool
}

func NewItemRepository(db *pgxpool.Pool) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) List(ctx context.Context) ([]model.Item, error) {
	rows, err := r.db.Query(ctx, "SELECT id, name, created_at FROM items ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *ItemRepository) Create(ctx context.Context, name string) (model.Item, error) {
	var item model.Item
	err := r.db.QueryRow(ctx,
		"INSERT INTO items (name) VALUES ($1) RETURNING id, name, created_at",
		name,
	).Scan(&item.ID, &item.Name, &item.CreatedAt)
	return item, err
}

func (r *ItemRepository) GetByID(ctx context.Context, id int64) (model.Item, error) {
	var item model.Item
	err := r.db.QueryRow(ctx,
		"SELECT id, name, created_at FROM items WHERE id = $1",
		id,
	).Scan(&item.ID, &item.Name, &item.CreatedAt)
	return item, err
}

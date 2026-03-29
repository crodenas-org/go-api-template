package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"hello-world-go/internal/handler"
	"hello-world-go/internal/repository"
)

func New(db *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/hello", handler.Hello)

	itemHandler := handler.NewItemHandler(repository.NewItemRepository(db))
	r.Get("/items", itemHandler.List)
	r.Post("/items", itemHandler.Create)
	r.Get("/items/{id}", itemHandler.GetByID)

	return r
}

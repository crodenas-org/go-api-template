package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/jackc/pgx/v5/pgxpool"

	"hello-world-go/internal/handler"
	appmiddleware "hello-world-go/internal/middleware"
	"hello-world-go/internal/repository"
)

func New(db *pgxpool.Pool, verifier *oidc.IDTokenVerifier) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// Public routes
	r.Get("/hello", handler.Hello)

	// Authenticated routes
	r.Group(func(r chi.Router) {
		r.Use(appmiddleware.Authenticate(verifier))

		itemHandler := handler.NewItemHandler(repository.NewItemRepository(db))
		r.Get("/items", itemHandler.List)
		r.Get("/items/{id}", itemHandler.GetByID)

		// Role-protected routes
		r.Group(func(r chi.Router) {
			r.Use(appmiddleware.RequireRole("items.write"))
			r.Post("/items", itemHandler.Create)
		})
	})

	return r
}

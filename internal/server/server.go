package server

import (
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	httpswagger "github.com/swaggo/http-swagger/v2"

	_ "hello-world-go/docs"
	"hello-world-go/internal/handler"
	appmiddleware "hello-world-go/internal/middleware"
	"hello-world-go/internal/repository"
)

func New(db *pgxpool.Pool, verifier *oidc.IDTokenVerifier) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173", // Vite dev server
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		MaxAge:         300,
	}))

	// Docs — Swagger UI (BaseLayout hides the Explore bar)
	r.Get("/docs/*", httpswagger.Handler(
		httpswagger.UIConfig(map[string]string{
			"layout": `"BaseLayout"`,
		}),
	))

	// Allow CORS preflight on all routes
	r.Options("/*", func(w http.ResponseWriter, r *http.Request) {})

	// Public routes
	r.Get("/hello", handler.Hello)

	// Authenticated routes
	r.Group(func(r chi.Router) {
		r.Use(appmiddleware.Authenticate(verifier))

		itemHandler := handler.NewItemHandler(repository.NewItemRepository(db))

		r.Group(func(r chi.Router) {
			r.Use(appmiddleware.RequireRole("items.read"))
			r.Get("/items", itemHandler.List)
			r.Get("/items/{id}", itemHandler.GetByID)
		})

		r.Group(func(r chi.Router) {
			r.Use(appmiddleware.RequireRole("items.write"))
			r.Post("/items", itemHandler.Create)
			r.Put("/items/{id}", itemHandler.Update)
			r.Delete("/items/{id}", itemHandler.Delete)
		})
	})

	return r
}

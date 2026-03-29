package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"hello-world-go/internal/handler"
)

func New() http.Handler {
	r := chi.NewRouter()

	r.Get("/hello", handler.Hello)

	return r
}

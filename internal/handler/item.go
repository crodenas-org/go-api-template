package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"hello-world-go/internal/model"
	"hello-world-go/internal/repository"
)

type ItemHandler struct {
	repo *repository.ItemRepository
}

func NewItemHandler(repo *repository.ItemRepository) *ItemHandler {
	return &ItemHandler{repo: repo}
}

// List godoc
// @Summary      List items
// @Description  Returns all items. Requires items.read role.
// @Tags         items
// @Produce      json
// @Success      200  {array}   model.Item
// @Failure      403  {string}  string
// @Failure      500  {string}  string
// @Security     BearerAuth
// @Router       /items [get]
func (h *ItemHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, "failed to list items", http.StatusInternalServerError)
		return
	}

	// Return an empty array instead of null when there are no items
	if items == nil {
		items = []model.Item{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// CreateRequest is the request body for creating an item.
type CreateRequest struct {
	Name string `json:"name" example:"my item"`
}

// Create godoc
// @Summary      Create an item
// @Description  Creates a new item. Requires items.write role.
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        item  body      CreateRequest  true  "Item to create"
// @Success      201   {object}  model.Item
// @Failure      400   {string}  string
// @Failure      403   {string}  string
// @Failure      500   {string}  string
// @Security     BearerAuth
// @Router       /items [post]
func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	item, err := h.repo.Create(r.Context(), req.Name)
	if err != nil {
		http.Error(w, "failed to create item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// UpdateRequest is the request body for updating an item.
type UpdateRequest struct {
	Name string `json:"name" example:"updated name"`
}

// Update godoc
// @Summary      Update an item
// @Description  Updates an item's name. Requires items.write role.
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        id    path      int            true  "Item ID"
// @Param        item  body      UpdateRequest  true  "Updated fields"
// @Success      200   {object}  model.Item
// @Failure      400   {string}  string
// @Failure      403   {string}  string
// @Failure      404   {string}  string
// @Failure      500   {string}  string
// @Security     BearerAuth
// @Router       /items/{id} [put]
func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	item, err := h.repo.Update(r.Context(), id, req.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "item not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to update item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// Delete godoc
// @Summary      Delete an item
// @Description  Deletes an item by ID. Requires items.write role.
// @Tags         items
// @Param        id   path      int  true  "Item ID"
// @Success      204  {string}  string
// @Failure      400  {string}  string
// @Failure      403  {string}  string
// @Failure      404  {string}  string
// @Failure      500  {string}  string
// @Security     BearerAuth
// @Router       /items/{id} [delete]
func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.repo.Delete(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "item not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to delete item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetByID godoc
// @Summary      Get an item
// @Description  Returns a single item by ID. Requires items.read role.
// @Tags         items
// @Produce      json
// @Param        id   path      int  true  "Item ID"
// @Success      200  {object}  model.Item
// @Failure      400  {string}  string
// @Failure      403  {string}  string
// @Failure      404  {string}  string
// @Failure      500  {string}  string
// @Security     BearerAuth
// @Router       /items/{id} [get]
func (h *ItemHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	item, err := h.repo.GetByID(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "item not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to get item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

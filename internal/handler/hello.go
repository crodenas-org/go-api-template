package handler

import (
	"encoding/json"
	"net/http"
)

// Hello godoc
// @Summary      Hello world
// @Description  Returns a hello world message
// @Tags         general
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /hello [get]
func Hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello, World!"})
}

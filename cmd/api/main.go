package main

import (
	"log"
	"net/http"

	"hello-world-go/internal/server"
)

func main() {
	srv := server.New()

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", srv))
}

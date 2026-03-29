package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"hello-world-go/internal/server"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	tenantID := os.Getenv("AZURE_TENANT_ID")
	clientID := os.Getenv("AZURE_CLIENT_ID")
	if tenantID == "" || clientID == "" {
		log.Fatal("AZURE_TENANT_ID and AZURE_CLIENT_ID are required")
	}

	ctx := context.Background()

	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("database unreachable: %v", err)
	}
	log.Println("database connected")

	issuer := fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", tenantID)
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		log.Fatalf("failed to initialize OIDC provider: %v", err)
	}

	// Azure access tokens use "api://<clientID>" as the audience
	verifier := provider.Verifier(&oidc.Config{
		ClientID: fmt.Sprintf("api://%s", clientID),
	})
	log.Println("OIDC provider initialized")

	srv := server.New(db, verifier)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", srv))
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"hello-world-go/internal/repository"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
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

	repo := repository.NewItemRepository(db)

	s := server.NewMCPServer(
		"hello-world-go",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	s.AddTool(
		mcp.NewTool("list_items",
			mcp.WithDescription("List all items"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			items, err := repo.List(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to list items: %v", err)), nil
			}
			if len(items) == 0 {
				return mcp.NewToolResultText("No items found"), nil
			}
			result := ""
			for _, item := range items {
				result += fmt.Sprintf("- [%d] %s (created %s)\n", item.ID, item.Name, item.CreatedAt.Format("2006-01-02 15:04:05"))
			}
			return mcp.NewToolResultText(result), nil
		},
	)

	s.AddTool(
		mcp.NewTool("create_item",
			mcp.WithDescription("Create a new item"),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Name of the item to create"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name, err := req.RequireString("name")
			if err != nil {
				return mcp.NewToolResultError("name is required"), nil
			}
			item, err := repo.Create(ctx, name)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to create item: %v", err)), nil
			}
			return mcp.NewToolResultText(fmt.Sprintf("Created item [%d] %s", item.ID, item.Name)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("get_item",
			mcp.WithDescription("Get an item by ID"),
			mcp.WithNumber("id",
				mcp.Required(),
				mcp.Description("ID of the item to retrieve"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			idFloat, err := req.RequireFloat("id")
			if err != nil {
				return mcp.NewToolResultError("id is required"), nil
			}
			id, err := strconv.ParseInt(strconv.FormatFloat(idFloat, 'f', 0, 64), 10, 64)
			if err != nil {
				return mcp.NewToolResultError("invalid id"), nil
			}
			item, err := repo.GetByID(ctx, id)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("item not found: %v", err)), nil
			}
			return mcp.NewToolResultText(fmt.Sprintf("[%d] %s (created %s)", item.ID, item.Name, item.CreatedAt.Format("2006-01-02 15:04:05"))), nil
		},
	)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("MCP server error: %v", err)
	}
}

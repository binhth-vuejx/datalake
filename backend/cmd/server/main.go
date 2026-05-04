package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	datasetqueries "github.com/yourorg/datalake-free/internal/application/dataset/queries"
	"github.com/yourorg/datalake-free/internal/infrastructure/postgres"
	"github.com/yourorg/datalake-free/internal/infrastructure/rilldata"
	httpserver "github.com/yourorg/datalake-free/internal/interfaces/http"
)

func main() {
	// Load .env in development
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file, using environment variables")
	}

	ctx := context.Background()

	// ── Infrastructure ────────────────────────────────────────────────────────
	db, err := postgres.Connect(ctx)
	if err != nil {
		log.Fatalf("DB connect: %v", err)
	}
	defer db.Close()
	log.Println("✓ Database connected")

	// ── Repositories (infrastructure implements domain interfaces) ────────────
	
	// Load RillData config
	rillConfigPath := filepath.Join("..", "rill", "datasets.yml")
	datasetRepo, err := rilldata.NewRepository(rillConfigPath, db)
	if err != nil {
		log.Printf("⚠ RillData config not loaded: %v", err)
		log.Println("  Dataset API will not be available")
		datasetRepo = nil
	} else {
		datasets, _ := datasetRepo.LoadConfig(ctx)
		log.Printf("✓ RillData: Loaded %d datasets", len(datasets))
	}

	// ── Application use-cases ─────────────────────────────────────────────────
	deps := httpserver.Dependencies{
		// Database connection for command handlers
		DB: db,
	}
	
	// Add dataset handlers if RillData config loaded successfully
	if datasetRepo != nil {
		deps.ListDatasets = datasetqueries.NewListDatasetsHandler(datasetRepo)
		deps.ExecuteQuery = datasetqueries.NewExecuteQueryHandler(datasetRepo)
	}

	// ── HTTP server ───────────────────────────────────────────────────────────
	router := httpserver.NewRouter(deps)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("✓ Server:    http://localhost:%s", port)
	log.Printf("✓ OpenAPI:   http://localhost:%s/openapi.json", port)
	log.Printf("✓ API Docs:  http://localhost:%s/docs", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}

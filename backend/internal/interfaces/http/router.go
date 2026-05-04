package http

import (
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	datasetqueries "github.com/yourorg/datalake-free/internal/application/dataset/queries"
	datasethandler "github.com/yourorg/datalake-free/internal/interfaces/http/handlers/dataset"
	commandhandler "github.com/yourorg/datalake-free/internal/interfaces/http/handlers/commands"
	"github.com/yourorg/datalake-free/internal/interfaces/http/middleware"
)

type Dependencies struct {
	// Dataset handlers
	ListDatasets  *datasetqueries.ListDatasetsHandler
	ExecuteQuery  *datasetqueries.ExecuteQueryHandler
	
	// Database connection for command handlers
	DB interface{}
}

func NewRouter(deps Dependencies) http.Handler {
	router := chi.NewMux()

	// Global middleware
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(middleware.CORS)
	router.Use(middleware.Auth)

	// Health check
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	// huma API — exposes /openapi.json and /docs automatically
	config := huma.DefaultConfig("Datalake Free API", "1.0.0")
	config.Info.Description = "DDD-structured API powered by huma + Chi"
	api := humachi.New(router, config)

	// Register domain handlers
	
	// Register dataset handlers
	if deps.ListDatasets != nil && deps.ExecuteQuery != nil {
		datasethandler.Register(api, datasethandler.NewHandler(
			deps.ExecuteQuery,
			deps.ListDatasets,
		))
	}
	
	// Register command handlers
	if deps.DB != nil {
		db := deps.DB.(*pgxpool.Pool)
		commandhandler.Register(api, commandhandler.NewHandler(db))
	}

	return router
}

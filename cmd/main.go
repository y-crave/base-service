package main

import (
	"database/sql"
	"log/slog"
	"os"
	"storage-service/internal/config"
	"storage-service/internal/controller"
	"storage-service/internal/logger"
	"storage-service/internal/middleware"
	"storage-service/internal/service"

	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Errorf(err.Error())
		os.Exit(1)
	}
	appLogger := logger.NewLogger(cfg.LogLevel)
	slog.SetDefault(appLogger)

	db, err := sql.Open("postgres", cfg.PostgresDSN)
	if err != nil {
		slog.Default().Error("Failed to connect to DB", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Default().Error("Failed to close DB", err)
		}
	}()

	monitoringService := service.NewMonitoringService(db)
	monitoringController := controller.NewMonitoringController(monitoringService)

	mainRouter := mux.NewRouter()

	apiPrefix := fmt.Sprintf("/api/v1/%s/", cfg.AppName)
	router := mainRouter.PathPrefix(apiPrefix).Subrouter()
	monitoringController.RegisterRoutes(router)

	addr := fmt.Sprintf("%s:%d", cfg.HttpHost, cfg.HttpPort)
	handler := middleware.LoggingMiddleware(router)

	middleware.PrintRoutes(router)
	slog.Default().Info("Server started:", "address", addr)
	slog.Default().Error("Fatal", "error", http.ListenAndServe(addr, handler))
}

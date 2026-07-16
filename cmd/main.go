package main

import (
	"base-service/internal/config"
	"base-service/internal/controller"
	"base-service/internal/logger"
	"base-service/internal/middleware"
	"base-service/internal/repository/subscription_repo"
	"base-service/internal/service"
	"database/sql"
	"log/slog"

	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Default().Error(err.Error())
		os.Exit(1)
	}
	appLogger := logger.NewLogger(cfg.LogLevel)
	slog.SetDefault(appLogger)

	db, err := sql.Open("postgres", cfg.PostgresDSN)
	if err != nil {
		slog.Default().Error("Failed to connect to DB", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Default().Error("Failed to close DB", "error", err)
		}
	}()

	monitoringService := service.NewMonitoringService(db)
	monitoringController := controller.NewMonitoringController(monitoringService)

	subscriptionRepo := subscription_repo.NewSubscriptionRepository(db)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo, service.NewDevMockValidator())
	subscriptionsController := controller.NewSubscriptionsController(subscriptionService)

	mainRouter := mux.NewRouter()

	apiPrefix := fmt.Sprintf("/api/v1/%s/", cfg.AppName)
	router := mainRouter.PathPrefix(apiPrefix).Subrouter()
	monitoringController.RegisterRoutes(router)

	subscriptionsRouter := mainRouter.PathPrefix("/api/v1/subscriptions").Subrouter()
	subscriptionsRouter.Use(middleware.GatewayIdentity(cfg))
	subscriptionsRouter.HandleFunc("/me", subscriptionsController.GetMe).Methods(http.MethodGet)
	subscriptionsRouter.HandleFunc("", subscriptionsController.PostSubscription).Methods(http.MethodPost)

	addr := fmt.Sprintf("%s:%d", cfg.HttpHost, cfg.HttpPort)
	handler := middleware.LoggingMiddleware(router)

	middleware.PrintRoutes(router)
	slog.Default().Info("Server started:", "address", addr)
	slog.Default().Error("Fatal", "error", http.ListenAndServe(addr, handler))
}

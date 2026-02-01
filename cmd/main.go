package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"storage-service/internal/config"
	"storage-service/internal/controller"
	"storage-service/internal/service"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.PostgresDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	monitoringService := service.NewMonitoringService(db)
	monitoringController := controller.NewMonitoringController(monitoringService)

	mainRouter := mux.NewRouter()

	apiPrefix := fmt.Sprintf("/api/v1/%s/", cfg.AppName)
	router := mainRouter.PathPrefix(apiPrefix).Subrouter()
	monitoringController.RegisterRoutes(router)

	addr := fmt.Sprintf("%s:%d", cfg.AppHost, cfg.AppHttpPort)
	handler := config.LoggingMiddleware(router)

	config.PrintRoutes(router)
	log.Printf("Server started: %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}

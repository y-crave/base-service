package controller

import (
	"context"
	"net/http"
	"storage-service/internal/service"
	"time"

	"github.com/gorilla/mux"
)

type MonitoringController struct {
	service service.MonitoringService
}

func NewMonitoringController(service service.MonitoringService) *MonitoringController {
	return &MonitoringController{service: service}
}

func (c *MonitoringController) LivenessProbe(w http.ResponseWriter, r *http.Request) {
	// Простая проверка: если сервер запущен — OK
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("alive"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (c *MonitoringController) ReadinessProbe(w http.ResponseWriter, r *http.Request) {
	// Проверяем подключение к БД
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := c.service.CheckDB(ctx); err != nil {
		http.Error(w, "DB not ready", http.StatusServiceUnavailable)
		return
	}

	// TODO: добавить проверку других зависимостей: Redis, Kafka и т.д.

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("ready"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (c *MonitoringController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/healthz", c.LivenessProbe).Methods(http.MethodGet)
	router.HandleFunc("/ready", c.ReadinessProbe).Methods(http.MethodGet)
}

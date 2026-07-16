package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"base-service/internal/authidentity"
	"base-service/internal/service"
)

type subscriptionServicer interface {
	GetStatus(ctx context.Context, userID uuid.UUID) (*service.SubscriptionStatusDTO, error)
	CreateSubscription(ctx context.Context, userID uuid.UUID, in service.CreateSubscriptionInput) (*service.SubscriptionStatusDTO, error)
}

// SubscriptionsController handles GET/POST /api/v1/subscriptions.
type SubscriptionsController struct {
	svc subscriptionServicer
}

func NewSubscriptionsController(svc subscriptionServicer) *SubscriptionsController {
	return &SubscriptionsController{svc: svc}
}

type createSubscriptionRequest struct {
	Store    string `json:"store"`
	Receipt  string `json:"receipt"`
	PlanCode string `json:"plan_code"`
}

func (c *SubscriptionsController) GetMe(w http.ResponseWriter, r *http.Request) {
	id, ok := authidentity.FromContext(r.Context())
	if !ok {
		writeSubscriptionError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	status, err := c.svc.GetStatus(r.Context(), id.UserID)
	if err != nil {
		slog.Error("subscriptions GetMe failed", "error", err)
		writeSubscriptionError(w, http.StatusInternalServerError, "internal_error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}

func (c *SubscriptionsController) PostSubscription(w http.ResponseWriter, r *http.Request) {
	id, ok := authidentity.FromContext(r.Context())
	if !ok {
		writeSubscriptionError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeSubscriptionError(w, http.StatusBadRequest, "invalid_request")
		return
	}

	status, err := c.svc.CreateSubscription(r.Context(), id.UserID, service.CreateSubscriptionInput{
		Store:    req.Store,
		Receipt:  req.Receipt,
		PlanCode: req.PlanCode,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSubscriptionInvalidInput):
			writeSubscriptionError(w, http.StatusBadRequest, "invalid_request")
		case errors.Is(err, service.ErrSubscriptionPlanNotFound):
			writeSubscriptionError(w, http.StatusNotFound, "plan_not_found")
		case errors.Is(err, service.ErrInvalidReceipt):
			writeSubscriptionError(w, http.StatusUnprocessableEntity, "invalid_receipt")
		default:
			slog.Error("subscriptions PostSubscription failed", "error", err)
			writeSubscriptionError(w, http.StatusInternalServerError, "internal_error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(status)
}

func writeSubscriptionError(w http.ResponseWriter, status int, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(`{"error":"` + code + `"}`))
}

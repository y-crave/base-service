package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"base-service/internal/repository/subscription_repo"
)

var (
	ErrSubscriptionInvalidInput = errors.New("invalid subscription input")
	ErrSubscriptionPlanNotFound = errors.New("plan not found")
)

// SubscriptionService handles subscription status lookups and receipt-backed grants.
type SubscriptionService struct {
	repo      *subscription_repo.SubscriptionRepository
	validator IAPReceiptValidator
}

func NewSubscriptionService(repo *subscription_repo.SubscriptionRepository, validator IAPReceiptValidator) *SubscriptionService {
	return &SubscriptionService{repo: repo, validator: validator}
}

// SubscriptionStatusDTO is the response for GET /subscriptions/me.
type SubscriptionStatusDTO struct {
	IsPremium bool       `json:"is_premium"`
	PlanCode  *string    `json:"plan_code"`
	Status    *string    `json:"status"`
	ExpiresAt *time.Time `json:"expires_at"`
}

// CreateSubscriptionInput is the service input for POST /subscriptions.
type CreateSubscriptionInput struct {
	Store    string
	Receipt  string
	PlanCode string
}

func (s *SubscriptionService) GetStatus(ctx context.Context, userID uuid.UUID) (*SubscriptionStatusDTO, error) {
	sub, err := s.repo.GetActiveSubscription(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("subscription_service.GetStatus: %w", err)
	}
	if sub == nil {
		return &SubscriptionStatusDTO{IsPremium: false}, nil
	}

	dto := &SubscriptionStatusDTO{
		IsPremium: true,
		PlanCode:  &sub.PlanCode,
		Status:    &sub.Status,
	}
	if sub.ExpiresAt.Valid {
		dto.ExpiresAt = &sub.ExpiresAt.Time
	}
	return dto, nil
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, userID uuid.UUID, in CreateSubscriptionInput) (*SubscriptionStatusDTO, error) {
	if in.Receipt == "" || in.PlanCode == "" {
		return nil, ErrSubscriptionInvalidInput
	}
	switch in.Store {
	case "ios", "android", "dev":
	default:
		return nil, ErrSubscriptionInvalidInput
	}

	plan, err := s.repo.GetPlanByCode(ctx, in.PlanCode)
	if err != nil {
		if errors.Is(err, subscription_repo.ErrPlanNotFound) {
			return nil, ErrSubscriptionPlanNotFound
		}
		return nil, fmt.Errorf("subscription_service.CreateSubscription plan lookup: %w", err)
	}

	validated, err := s.validator.ValidateReceipt(ctx, ReceiptValidationInput{
		Store:    in.Store,
		Receipt:  in.Receipt,
		PlanCode: in.PlanCode,
	})
	if err != nil {
		if errors.Is(err, ErrInvalidReceipt) {
			return nil, ErrInvalidReceipt
		}
		return nil, fmt.Errorf("subscription_service.CreateSubscription validate: %w", err)
	}

	expiresAt := validated.ExpiresAt
	if expiresAt == nil {
		fallback := time.Now().UTC().AddDate(0, 0, plan.DurationDays)
		expiresAt = &fallback
	}

	row, err := s.repo.CreateSubscription(ctx, userID, plan.ID, in.Store, validated.TransactionID, sql.NullTime{Time: *expiresAt, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("subscription_service.CreateSubscription: %w", err)
	}

	status := row.Status
	return &SubscriptionStatusDTO{
		IsPremium: true,
		PlanCode:  &plan.Code,
		Status:    &status,
		ExpiresAt: expiresAt,
	}, nil
}

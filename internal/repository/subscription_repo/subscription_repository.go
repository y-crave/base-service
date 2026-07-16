package subscription_repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var ErrPlanNotFound = errors.New("plan not found")

// SubscriptionRepository persists plans and per-user subscription state.
type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// PlanRow is a persisted subscription plan.
type PlanRow struct {
	ID           uuid.UUID
	Code         string
	Name         string
	DurationDays int
}

// SubscriptionRow is a persisted user subscription, joined with its plan.
type SubscriptionRow struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	PlanCode   string
	PlanName   string
	Status     string
	Store      string
	ExpiresAt  sql.NullTime
	CreatedAt  time.Time
}

// GetPlanByCode looks up an active plan by its code (e.g. "premium_monthly").
func (r *SubscriptionRepository) GetPlanByCode(ctx context.Context, code string) (*PlanRow, error) {
	const q = `
SELECT id, code, name, duration_days
FROM plans
WHERE code = $1 AND is_active = true`

	row := &PlanRow{}
	err := r.db.QueryRowContext(ctx, q, code).Scan(&row.ID, &row.Code, &row.Name, &row.DurationDays)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPlanNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("subscription_repository.GetPlanByCode: %w", err)
	}
	return row, nil
}

// GetActiveSubscription returns the user's current active (non-expired) subscription, or nil if none.
func (r *SubscriptionRepository) GetActiveSubscription(ctx context.Context, userID uuid.UUID) (*SubscriptionRow, error) {
	const q = `
SELECT s.id, s.user_id, p.code, p.name, s.status, s.store, s.expires_at, s.created_at
FROM user_subscriptions s
JOIN plans p ON p.id = s.plan_id
WHERE s.user_id = $1 AND s.status = 'active' AND (s.expires_at IS NULL OR s.expires_at > now())
ORDER BY s.created_at DESC
LIMIT 1`

	row := &SubscriptionRow{}
	err := r.db.QueryRowContext(ctx, q, userID).Scan(
		&row.ID, &row.UserID, &row.PlanCode, &row.PlanName, &row.Status, &row.Store, &row.ExpiresAt, &row.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("subscription_repository.GetActiveSubscription: %w", err)
	}
	return row, nil
}

// CreateSubscription records a new subscription grant for userID and expires any previously active ones.
func (r *SubscriptionRepository) CreateSubscription(
	ctx context.Context,
	userID, planID uuid.UUID,
	store, receiptRef string,
	expiresAt sql.NullTime,
) (*SubscriptionRow, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("subscription_repository.CreateSubscription begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `
UPDATE user_subscriptions SET status = 'expired', updated_at = now()
WHERE user_id = $1 AND status = 'active'`, userID)
	if err != nil {
		return nil, fmt.Errorf("subscription_repository.CreateSubscription expire previous: %w", err)
	}

	const q = `
INSERT INTO user_subscriptions (user_id, plan_id, status, store, receipt_ref, expires_at)
VALUES ($1, $2, 'active', $3, $4, $5)
RETURNING id, user_id, status, store, expires_at, created_at`

	row := &SubscriptionRow{UserID: userID}
	err = tx.QueryRowContext(ctx, q, userID, planID, store, receiptRef, expiresAt).Scan(
		&row.ID, &row.UserID, &row.Status, &row.Store, &row.ExpiresAt, &row.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("subscription_repository.CreateSubscription: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("subscription_repository.CreateSubscription commit: %w", err)
	}
	return row, nil
}

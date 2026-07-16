package service

import (
	"context"
	"errors"
	"time"
)

var ErrInvalidReceipt = errors.New("invalid IAP receipt")

// ReceiptValidationInput is what the client submitted for POST /subscriptions.
type ReceiptValidationInput struct {
	Store    string // "ios" | "android" | "dev"
	Receipt  string
	PlanCode string
}

// ValidatedReceipt is what a validator confirmed with the store.
// ExpiresAt is nil when the store doesn't report an expiry and the service
// should fall back to the plan's duration_days.
type ValidatedReceipt struct {
	TransactionID string
	ExpiresAt     *time.Time
}

// IAPReceiptValidator confirms an in-app-purchase receipt with the issuing store.
// Swap in a real StoreKit/Google Play adapter for prod; DevMockValidator is for
// local/dev/QA environments where no store credentials are configured yet.
type IAPReceiptValidator interface {
	ValidateReceipt(ctx context.Context, in ReceiptValidationInput) (ValidatedReceipt, error)
}

// DevMockValidator trusts any non-empty receipt without contacting a real store.
// Never wire this into a production deployment.
type DevMockValidator struct{}

func NewDevMockValidator() *DevMockValidator {
	return &DevMockValidator{}
}

func (v *DevMockValidator) ValidateReceipt(_ context.Context, in ReceiptValidationInput) (ValidatedReceipt, error) {
	if in.Receipt == "" {
		return ValidatedReceipt{}, ErrInvalidReceipt
	}
	return ValidatedReceipt{TransactionID: in.Receipt}, nil
}

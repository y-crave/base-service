package authidentity

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey struct{}

// Identity carries gateway-forwarded OIDC identity (Istio validates JWT upstream).
type Identity struct {
	UserID      uuid.UUID
	Email       string
	DisplayName string
}

func WithContext(ctx context.Context, id Identity) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}

func FromContext(ctx context.Context) (Identity, bool) {
	v := ctx.Value(ctxKey{})
	if v == nil {
		return Identity{}, false
	}
	id, ok := v.(Identity)
	return id, ok
}

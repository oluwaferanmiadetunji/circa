package auth

import (
	"context"
)

type AuthService interface {
	CreatePendingSignup(ctx context.Context, fullName, email string, displayName *string) (*SignupResult, error)
	GenerateNonce(address string) (*NonceResult, error)
}


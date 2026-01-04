package auth

import (
	sqlc "circa/internal/db/sqlc/generated"
	"context"
)

type VerifyTokenResult struct {
	Email       string
	DisplayName *string
	SessionID   string
	NeedsWallet bool // true for signup, false for login
}

type CompleteSignupResult struct {
	User      sqlc.User
	SessionID string
}

type LoginResult struct {
	Message string
}

type GetSessionUserResult struct {
	User sqlc.User
}

type AuthService interface {
	CreatePendingSignup(ctx context.Context, fullName, email string, displayName *string) (*SignupResult, error)
	CreateLoginMagicLink(ctx context.Context, email string) (*LoginResult, error)
	GenerateNonce(ctx context.Context, sessionID, address string, chainID *int64) (*NonceResult, error)
	VerifyToken(ctx context.Context, token string) (*VerifyTokenResult, error)
	GetSignupSession(ctx context.Context, sessionID string) (map[string]interface{}, error)
	CompleteSignup(ctx context.Context, sessionID, address, signature, message string) (*CompleteSignupResult, error)
	GetSessionUser(ctx context.Context, sessionID string) (*GetSessionUserResult, error)
}

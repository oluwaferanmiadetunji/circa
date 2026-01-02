package auth

import (
	"circa/internal/db"
	sqlc "circa/internal/db/sqlc/generated"
	"circa/internal/errors"
	"circa/internal/queue"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

type Service struct {
	store        db.Store
	queueService *queue.Service
	frontendURL  string
	nonceExpiry  time.Duration
}

type NonceResult struct {
	Nonce           string
	ExpiresAt       time.Time
	MessageTemplate *string
}

type SignupResult struct {
	PendingSignup sqlc.PendingSignup
	MagicLink     sqlc.MagicLink
}

func NewService(store db.Store, queueService *queue.Service, frontendURL string, nonceExpiry time.Duration) *Service {
	return &Service{
		store:        store,
		queueService: queueService,
		frontendURL:  frontendURL,
		nonceExpiry:  nonceExpiry,
	}
}

// GenerateNonce generates a one-time nonce for wallet authentication
func (s *Service) GenerateNonce(address string) (*NonceResult, error) {
	// Generate a random 32-byte nonce
	nonceBytes := make([]byte, 32)
	if _, err := rand.Read(nonceBytes); err != nil {
		log.Error().Err(err).Msg("Failed to generate nonce")
		return nil, err
	}

	nonce := "0x" + hex.EncodeToString(nonceBytes)
	expiresAt := time.Now().Add(s.nonceExpiry)

	// TODO: Store nonce in Redis with expiry
	// This will be implemented when we add Redis integration

	messageTemplate := s.buildMessageTemplate(address, nonce)

	return &NonceResult{
		Nonce:           nonce,
		ExpiresAt:       expiresAt,
		MessageTemplate: &messageTemplate,
	}, nil
}

// buildMessageTemplate creates a SIWE-style message template
func (s *Service) buildMessageTemplate(address, nonce string) string {
	// SIWE (Sign-In With Ethereum) style message
	// Format: domain wants you to sign in with your Ethereum account:\n{address}\n\n{message}\n\nURI: {uri}\nVersion: 1\nChain ID: {chainId}\nNonce: {nonce}\nIssued At: {timestamp}
	return "circa wants you to sign in with your Ethereum account:\n" + address + "\n\nSign in to Circa\n\nURI: https://circa.app\nVersion: 1\nNonce: " + nonce + "\nIssued At: " + time.Now().UTC().Format(time.RFC3339)
}

func (s *Service) CreatePendingSignup(ctx context.Context, fullName, email string, displayName *string) (*SignupResult, error) {
	emailText := pgtype.Text{String: email, Valid: true}
	_, err := s.store.GetUserByEmail(ctx, emailText)
	if err == nil {
		return nil, errors.ErrEmailAlreadyExists
	}

	if err != pgx.ErrNoRows {
		log.Error().Err(err).Msg("Failed to check if email exists")
		return nil, err
	}

	pgxStore, ok := s.store.(*db.PGXStore)
	if !ok {
		return nil, errors.ErrInvalidStore
	}

	tx, err := pgxStore.GetDB().Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := pgxStore.Queries.WithTx(tx)

	expiresAt := time.Now().Add(24 * time.Hour)
	pendingSignupParams := sqlc.CreatePendingSignupParams{
		FullName:    pgtype.Text{String: fullName, Valid: true},
		Email:       emailText,
		DisplayName: displayName,
	}
	pendingSignup, err := qtx.CreatePendingSignup(ctx, pendingSignupParams)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create pending signup")
		return nil, err
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		log.Error().Err(err).Msg("Failed to generate token")
		return nil, err
	}

	token := hex.EncodeToString(tokenBytes)
	tokenHash := sha256.Sum256([]byte(token))
	tokenHashHex := hex.EncodeToString(tokenHash[:])

	magicLinkExpiresAt := pgtype.Timestamp{Time: expiresAt, Valid: true}
	magicLinkParams := sqlc.CreateMagicLinkParams{
		PendingSignupID: pendingSignup.ID,
		TokenHash:       tokenHashHex,
		ExpiresAt:       magicLinkExpiresAt,
	}

	magicLink, err := qtx.CreateMagicLink(ctx, magicLinkParams)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create magic link")
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return nil, err
	}

	magicLinkURL := fmt.Sprintf("%s/auth/verify?token=%s", s.frontendURL, token)

	if s.queueService != nil {
		recipientName := fullName
		if displayName != nil && *displayName != "" {
			recipientName = *displayName
		}
		_, err := s.queueService.Enqueue(ctx, "send_magic_link_email", queue.JobPayload{
			"email":          email,
			"name":           recipientName,
			"magic_link_url": magicLinkURL,
		}, nil)
		if err != nil {
			log.Error().Err(err).Msg("Failed to enqueue magic link email")
		}
	}

	return &SignupResult{
		PendingSignup: pendingSignup,
		MagicLink:     magicLink,
	}, nil
}

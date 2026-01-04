package auth

import (
	"circa/internal/db"
	sqlc "circa/internal/db/sqlc/generated"
	"circa/internal/errors"
	"circa/internal/queue"
	circaredis "circa/internal/redis"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

type Service struct {
	store        db.Store
	queueService *queue.Service
	frontendURL  string
	nonceExpiry  time.Duration
	redisClient  *redis.Client
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
	var redisClient *redis.Client
	if circaredis.RedisClient != nil {
		redisClient = circaredis.RedisClient
	}
	return &Service{
		store:        store,
		queueService: queueService,
		frontendURL:  frontendURL,
		nonceExpiry:  nonceExpiry,
		redisClient:  redisClient,
	}
}

func (s *Service) GetSignupSession(ctx context.Context, sessionID string) (map[string]any, error) {
	if s.redisClient == nil {
		return nil, errors.ErrInvalidSession
	}

	sessionKey := fmt.Sprintf("signup_session:%s", sessionID)
	sessionJSON, err := s.redisClient.Get(ctx, sessionKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.ErrInvalidSession
		}
		log.Error().Err(err).Msg("Failed to get signup session from Redis")
		return nil, err
	}

	var sessionData map[string]any
	if err := json.Unmarshal([]byte(sessionJSON), &sessionData); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal session data")
		return nil, err
	}

	return sessionData, nil
}

func (s *Service) GenerateNonce(ctx context.Context, sessionID, address string, chainID *int64) (*NonceResult, error) {
	_, err := s.GetSignupSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	nonceBytes := make([]byte, 32)
	if _, err := rand.Read(nonceBytes); err != nil {
		log.Error().Err(err).Msg("Failed to generate nonce")
		return nil, err
	}

	nonce := "0x" + hex.EncodeToString(nonceBytes)
	expiresAt := time.Now().Add(s.nonceExpiry)

	if s.redisClient != nil {
		nonceData := map[string]any{
			"session_id": sessionID,
			"address":    address,
			"used":       false,
			"created_at": time.Now().Unix(),
			"expires_at": expiresAt.Unix(),
		}

		nonceJSON, err := json.Marshal(nonceData)
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal nonce data")
			return nil, err
		}

		// Store with key format: nonce:{nonce} for quick lookup
		nonceKey := fmt.Sprintf("nonce:%s", nonce)
		expiry := s.nonceExpiry
		err = s.redisClient.Set(ctx, nonceKey, nonceJSON, expiry).Err()
		if err != nil {
			log.Error().Err(err).Msg("Failed to store nonce in Redis")
			return nil, err
		}

		// Also store a mapping from session+address to nonce for validation
		sessionNonceKey := fmt.Sprintf("session_nonce:%s:%s", sessionID, address)
		err = s.redisClient.Set(ctx, sessionNonceKey, nonce, expiry).Err()
		if err != nil {
			log.Error().Err(err).Msg("Failed to store session-nonce mapping in Redis")
			// Don't fail if this fails, the main nonce key is more important
		}
	} else {
		log.Warn().Msg("Redis client not available, nonce will not be persisted")
	}

	messageTemplate := s.buildMessageTemplate(address, nonce, chainID)

	return &NonceResult{
		Nonce:           nonce,
		ExpiresAt:       expiresAt,
		MessageTemplate: &messageTemplate,
	}, nil
}

// buildMessageTemplate creates a SIWE-style message template
func (s *Service) buildMessageTemplate(address, nonce string, chainID *int64) string {
	message := "circa wants you to sign in with your Ethereum account:\n" + address + "\n\nSign in to Circa\n\nURI: " + s.frontendURL + "\nVersion: 1\n"
	if chainID != nil {
		message += fmt.Sprintf("Chain ID: %d\n", *chainID)
	}
	message += "Nonce: " + nonce + "\nIssued At: " + time.Now().UTC().Format(time.RFC3339)
	return message
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

	if err := qtx.InvalidateMagicLinksByEmail(ctx, emailText); err != nil {
		log.Error().Err(err).Msg("Failed to invalidate old magic links")
		return nil, err
	}

	if err := qtx.InvalidatePendingSignupsByEmail(ctx, emailText); err != nil {
		log.Error().Err(err).Msg("Failed to invalidate old pending signups")
		return nil, err
	}

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

		go func() {
			bgCtx := context.Background()
			_, err := s.queueService.Enqueue(bgCtx, "send_magic_link_email", queue.JobPayload{
				"email":          email,
				"name":           recipientName,
				"magic_link_url": magicLinkURL,
				"is_login":       false,
			}, nil)
			if err != nil {
				log.Error().Err(err).Msg("Failed to enqueue magic link email")
			}
		}()
	}

	return &SignupResult{
		PendingSignup: pendingSignup,
		MagicLink:     magicLink,
	}, nil
}

// CreateLoginMagicLink creates a magic link for an existing user to log in
func (s *Service) CreateLoginMagicLink(ctx context.Context, email string) (*LoginResult, error) {
	emailText := pgtype.Text{String: email, Valid: true}

	// Check if user exists
	user, err := s.store.GetUserByEmail(ctx, emailText)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &LoginResult{
				Message: "If an account exists with this email, you will receive a login link.",
			}, nil
		}
		log.Error().Err(err).Msg("Failed to check if user exists")
		return nil, err
	}

	// Generate token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		log.Error().Err(err).Msg("Failed to generate token")
		return nil, err
	}

	token := hex.EncodeToString(tokenBytes)
	tokenHash := sha256.Sum256([]byte(token))
	tokenHashHex := hex.EncodeToString(tokenHash[:])

	// Store login magic link in Redis (expires in 24 hours)
	if s.redisClient != nil {
		loginLinkData := map[string]any{
			"user_id":    user.ID.String(),
			"email":      user.Email.String,
			"address":    user.Address,
			"created_at": time.Now().Unix(),
			"expires_at": time.Now().Add(24 * time.Hour).Unix(),
		}

		loginLinkJSON, err := json.Marshal(loginLinkData)
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal login link data")
			return nil, err
		}

		expiry := 24 * time.Hour
		loginLinkKey := fmt.Sprintf("login_magic_link:%s", tokenHashHex)
		err = s.redisClient.Set(ctx, loginLinkKey, loginLinkJSON, expiry).Err()
		if err != nil {
			log.Error().Err(err).Msg("Failed to store login magic link in Redis")
			return nil, err
		}
	} else {
		log.Warn().Msg("Redis client not available, login magic link will not be persisted")
	}

	magicLinkURL := fmt.Sprintf("%s/auth/verify?token=%s", s.frontendURL, token)

	// Send email
	if s.queueService != nil {
		recipientName := user.Email.String
		if user.DisplayName != nil && *user.DisplayName != "" {
			recipientName = *user.DisplayName
		} else if user.FullName.Valid {
			recipientName = user.FullName.String
		}

		go func() {
			bgCtx := context.Background()
			_, err := s.queueService.Enqueue(bgCtx, "send_magic_link_email", queue.JobPayload{
				"email":          email,
				"name":           recipientName,
				"magic_link_url": magicLinkURL,
				"is_login":       true,
			}, nil)
			if err != nil {
				log.Error().Err(err).Msg("Failed to enqueue login magic link email")
			}
		}()
	}

	return &LoginResult{
		Message: "If an account exists with this email, you will receive a login link.",
	}, nil
}

func (s *Service) VerifyToken(ctx context.Context, token string) (*VerifyTokenResult, error) {
	tokenHash := sha256.Sum256([]byte(token))
	tokenHashHex := hex.EncodeToString(tokenHash[:])

	// First check if it's a login magic link (stored in Redis)
	if s.redisClient != nil {
		loginLinkKey := fmt.Sprintf("login_magic_link:%s", tokenHashHex)
		loginLinkJSON, err := s.redisClient.Get(ctx, loginLinkKey).Result()
		if err == nil {
			// It's a login magic link
			var loginLinkData map[string]any
			if err := json.Unmarshal([]byte(loginLinkJSON), &loginLinkData); err != nil {
				log.Error().Err(err).Msg("Failed to unmarshal login link data")
				return nil, err
			}

			// Delete the login magic link (one-time use)
			s.redisClient.Del(ctx, loginLinkKey)

			// Create main session directly (user already exists)
			mainSessionID := uuid.New().String()
			sessionData := map[string]any{
				"user_id":    loginLinkData["user_id"],
				"address":    loginLinkData["address"],
				"email":      loginLinkData["email"],
				"created_at": time.Now().Unix(),
			}

			sessionJSON, err := json.Marshal(sessionData)
			if err != nil {
				log.Error().Err(err).Msg("Failed to marshal session data")
				return nil, err
			}

			// Main session expires in 7 days
			expiry := 7 * 24 * time.Hour
			err = s.redisClient.Set(ctx, fmt.Sprintf("session:%s", mainSessionID), sessionJSON, expiry).Err()
			if err != nil {
				log.Error().Err(err).Msg("Failed to store main session in Redis")
				return nil, err
			}

			email := loginLinkData["email"].(string)
			var displayName *string
			if dn, ok := loginLinkData["display_name"].(string); ok && dn != "" {
				displayName = &dn
			}

			return &VerifyTokenResult{
				Email:       email,
				DisplayName: displayName,
				SessionID:   mainSessionID,
				NeedsWallet: false, // Login doesn't need wallet
			}, nil
		}
	}

	// If not a login magic link, check for signup magic link
	pgxStore, ok := s.store.(*db.PGXStore)
	if !ok {
		return nil, errors.ErrInvalidStore
	}

	magicLink, err := pgxStore.Queries.GetMagicLinkByTokenHash(ctx, tokenHashHex)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrInvalidToken
		}
		log.Error().Err(err).Msg("Failed to get magic link by token hash")
		return nil, err
	}

	pendingSignup, err := pgxStore.Queries.GetPendingSignupByID(ctx, magicLink.PendingSignupID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrInvalidToken
		}
		log.Error().Err(err).Msg("Failed to get pending signup")
		return nil, err
	}

	tx, err := pgxStore.GetDB().Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := pgxStore.Queries.WithTx(tx)

	_, err = qtx.MarkMagicLinkAsUsed(ctx, magicLink.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to mark magic link as used")
		return nil, err
	}

	now := pgtype.Timestamp{Time: time.Now(), Valid: true}
	_, err = qtx.UpdatePendingSignup(ctx, sqlc.UpdatePendingSignupParams{
		EmailVerifiedAt: now,
		ID:              pendingSignup.ID,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to update pending signup")
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return nil, err
	}

	sessionID := uuid.New().String()
	sessionData := map[string]any{
		"pending_signup_id": pendingSignup.ID.String(),
		"email":             pendingSignup.Email.String,
		"full_name":         pendingSignup.FullName.String,
		"display_name":      pendingSignup.DisplayName,
		"created_at":        time.Now().Unix(),
	}

	if s.redisClient != nil {
		sessionJSON, err := json.Marshal(sessionData)
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal session data")
			return nil, err
		}

		expiry := 30 * time.Minute
		err = s.redisClient.Set(ctx, fmt.Sprintf("signup_session:%s", sessionID), sessionJSON, expiry).Err()
		if err != nil {
			log.Error().Err(err).Msg("Failed to store signup session in Redis")
			return nil, err
		}
	} else {
		log.Warn().Msg("Redis client not available, signup session will not be persisted")
	}

	email := pendingSignup.Email.String
	displayName := pendingSignup.DisplayName

	return &VerifyTokenResult{
		Email:       email,
		DisplayName: displayName,
		SessionID:   sessionID,
		NeedsWallet: true, // Signup needs wallet connection
	}, nil
}

// verifySignature verifies an Ethereum personal_sign signature
func verifySignature(address, message, signature string) (bool, error) {
	// Remove 0x prefix if present
	sig := strings.TrimPrefix(signature, "0x")
	if len(sig) != 130 {
		return false, fmt.Errorf("invalid signature length")
	}

	sigBytes, err := hex.DecodeString(sig)
	if err != nil {
		return false, fmt.Errorf("invalid signature hex: %w", err)
	}

	// Ethereum personal_sign adds a prefix
	msg := accounts.TextHash([]byte(message))

	// Recover public key (sigBytes[64] is the recovery ID)
	if sigBytes[64] >= 27 {
		sigBytes[64] -= 27
	}

	pubKey, err := crypto.SigToPub(msg, sigBytes)
	if err != nil {
		return false, fmt.Errorf("failed to recover public key: %w", err)
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	expectedAddr := common.HexToAddress(address)

	return strings.EqualFold(recoveredAddr.Hex(), expectedAddr.Hex()), nil
}

// CompleteSignup completes the signup process by verifying the wallet signature and creating the user
func (s *Service) CompleteSignup(ctx context.Context, sessionID, address, signature, message string) (*CompleteSignupResult, error) {
	// 1. Verify signup session exists and is email-verified
	sessionData, err := s.GetSignupSession(ctx, sessionID)
	if err != nil {
		log.Warn().
			Str("session_id", sessionID).
			Err(err).
			Msg("Failed to get signup session")
		return nil, err
	}

	pendingSignupIDStr, ok := sessionData["pending_signup_id"].(string)
	if !ok {
		return nil, errors.ErrInvalidSession
	}

	pendingSignupID, err := uuid.Parse(pendingSignupIDStr)
	if err != nil {
		return nil, errors.ErrInvalidSession
	}

	pgxStore, ok := s.store.(*db.PGXStore)
	if !ok {
		return nil, errors.ErrInvalidStore
	}

	// Get pending signup (must be verified at this point)
	pendingSignup, err := pgxStore.Queries.GetVerifiedPendingSignupByID(ctx, pendingSignupID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Warn().
				Str("pending_signup_id", pendingSignupIDStr).
				Str("session_id", sessionID).
				Msg("Verified pending signup not found in database")
			return nil, errors.ErrInvalidSession
		}
		log.Error().Err(err).Msg("Failed to get verified pending signup")
		return nil, err
	}

	// Verify email is verified
	if !pendingSignup.EmailVerifiedAt.Valid {
		log.Warn().
			Str("pending_signup_id", pendingSignupIDStr).
			Str("session_id", sessionID).
			Bool("email_verified", pendingSignup.EmailVerifiedAt.Valid).
			Msg("Pending signup email not verified")
		return nil, errors.ErrInvalidSession
	}

	log.Info().
		Str("pending_signup_id", pendingSignupIDStr).
		Str("session_id", sessionID).
		Str("email", pendingSignup.Email.String).
		Msg("Pending signup found and email verified, proceeding with signature verification")

	// 2. Verify nonce exists, matches session+address, not expired, not used
	if s.redisClient == nil {
		return nil, fmt.Errorf("Redis not available")
	}

	// Extract nonce from message (look for "Nonce: 0x..." pattern)
	nonce := ""
	lines := strings.Split(message, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Nonce: ") {
			nonce = strings.TrimPrefix(line, "Nonce: ")
			break
		}
	}

	if nonce == "" {
		return nil, errors.ErrInvalidNonce
	}

	nonceKey := fmt.Sprintf("nonce:%s", nonce)
	nonceJSON, err := s.redisClient.Get(ctx, nonceKey).Result()
	if err != nil {
		if err == redis.Nil {
			log.Warn().
				Str("nonce", nonce).
				Str("session_id", sessionID).
				Msg("Nonce not found in Redis - may have expired")
			return nil, errors.ErrInvalidNonce
		}
		return nil, err
	}

	var nonceData map[string]any
	if err := json.Unmarshal([]byte(nonceJSON), &nonceData); err != nil {
		return nil, err
	}

	// Verify nonce matches session and address
	if nonceData["session_id"] != sessionID {
		log.Warn().
			Str("nonce", nonce).
			Str("expected_session", sessionID).
			Interface("nonce_session", nonceData["session_id"]).
			Msg("Nonce session mismatch")
		return nil, errors.ErrInvalidNonce
	}
	if nonceData["address"] != strings.ToLower(address) && nonceData["address"] != address {
		log.Warn().
			Str("nonce", nonce).
			Str("expected_address", strings.ToLower(address)).
			Interface("nonce_address", nonceData["address"]).
			Msg("Nonce address mismatch")
		return nil, errors.ErrInvalidNonce
	}

	// Check if nonce is used
	if used, ok := nonceData["used"].(bool); ok && used {
		log.Warn().
			Str("nonce", nonce).
			Msg("Nonce already used")
		return nil, errors.ErrInvalidNonce
	}

	// Mark nonce as used
	nonceData["used"] = true
	nonceJSONUpdated, _ := json.Marshal(nonceData)
	s.redisClient.Set(ctx, nonceKey, nonceJSONUpdated, 0)

	// 3. Verify signature recovers the address
	valid, err := verifySignature(address, message, signature)
	if err != nil || !valid {
		return nil, errors.ErrInvalidToken
	}

	// 4. Check wallet isn't already linked to another user
	_, err = pgxStore.Queries.GetUserByAddress(ctx, strings.ToLower(address))
	if err == nil {
		// User with this address already exists
		return nil, fmt.Errorf("wallet address already in use")
	}
	if err != pgx.ErrNoRows {
		log.Error().Err(err).Msg("Failed to check if address exists")
		return nil, err
	}

	// 5. Check email isn't already bound to another user
	emailText := pgtype.Text{String: pendingSignup.Email.String, Valid: true}
	_, err = pgxStore.Queries.GetUserByEmail(ctx, emailText)
	if err == nil {
		// User with this email already exists - we'll treat this as linking wallet to existing user
		// For now, return error. Can be changed based on policy
		return nil, fmt.Errorf("email already registered")
	}
	if err != pgx.ErrNoRows {
		log.Error().Err(err).Msg("Failed to check if email exists")
		return nil, err
	}

	// 6. Create user record
	tx, err := pgxStore.GetDB().Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := pgxStore.Queries.WithTx(tx)

	user, err := qtx.CreateUser(ctx, sqlc.CreateUserParams{
		FullName:    pendingSignup.FullName,
		Email:       pendingSignup.Email,
		Address:     strings.ToLower(address),
		DisplayName: pendingSignup.DisplayName,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return nil, err
	}

	// 7. Create main session
	mainSessionID := uuid.New().String()
	sessionDataMain := map[string]any{
		"user_id":    user.ID.String(),
		"address":    user.Address,
		"email":      user.Email.String,
		"created_at": time.Now().Unix(),
	}

	if s.redisClient != nil {
		sessionJSON, err := json.Marshal(sessionDataMain)
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal session data")
			return nil, err
		}

		// Main session expires in 7 days
		expiry := 7 * 24 * time.Hour
		err = s.redisClient.Set(ctx, fmt.Sprintf("session:%s", mainSessionID), sessionJSON, expiry).Err()
		if err != nil {
			log.Error().Err(err).Msg("Failed to store main session in Redis")
			return nil, err
		}
	} else {
		log.Warn().Msg("Redis client not available, main session will not be persisted")
	}

	return &CompleteSignupResult{
		User:      user,
		SessionID: mainSessionID,
	}, nil
}

// GetSessionUser retrieves the user associated with a session ID
func (s *Service) GetSessionUser(ctx context.Context, sessionID string) (*GetSessionUserResult, error) {
	if s.redisClient == nil {
		return nil, errors.ErrInvalidSession
	}

	sessionKey := fmt.Sprintf("session:%s", sessionID)
	sessionJSON, err := s.redisClient.Get(ctx, sessionKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.ErrInvalidSession
		}
		log.Error().Err(err).Msg("Failed to get session from Redis")
		return nil, err
	}

	var sessionData map[string]any
	if err := json.Unmarshal([]byte(sessionJSON), &sessionData); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal session data")
		return nil, err
	}

	userIDStr, ok := sessionData["user_id"].(string)
	if !ok {
		return nil, errors.ErrInvalidSession
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.ErrInvalidSession
	}

	pgxStore, ok := s.store.(*db.PGXStore)
	if !ok {
		return nil, errors.ErrInvalidStore
	}

	user, err := pgxStore.Queries.GetUserByID(ctx, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrInvalidSession
		}
		log.Error().Err(err).Msg("Failed to get user by ID")
		return nil, err
	}

	return &GetSessionUserResult{
		User: user,
	}, nil
}

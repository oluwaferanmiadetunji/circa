package auth

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/rs/zerolog/log"
)

type Service struct {
	nonceExpiry time.Duration
}

type NonceResult struct {
	Nonce      string
	ExpiresAt  time.Time
	MessageTemplate *string
}

// NewService creates a new auth service
func NewService(nonceExpiry time.Duration) *Service {
	return &Service{
		nonceExpiry: nonceExpiry,
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
		Nonce:          nonce,
		ExpiresAt:      expiresAt,
		MessageTemplate: &messageTemplate,
	}, nil
}

// buildMessageTemplate creates a SIWE-style message template
func (s *Service) buildMessageTemplate(address, nonce string) string {
	// SIWE (Sign-In With Ethereum) style message
	// Format: domain wants you to sign in with your Ethereum account:\n{address}\n\n{message}\n\nURI: {uri}\nVersion: 1\nChain ID: {chainId}\nNonce: {nonce}\nIssued At: {timestamp}
	return "circa wants you to sign in with your Ethereum account:\n" + address + "\n\nSign in to Circa\n\nURI: https://circa.app\nVersion: 1\nNonce: " + nonce + "\nIssued At: " + time.Now().UTC().Format(time.RFC3339)
}


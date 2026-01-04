package errors

import (
	"errors"
)

// Auth errors
var (
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidStore        = errors.New("store is not a PGXStore")
	ErrInvalidToken        = errors.New("invalid or expired token")
	ErrInvalidSession      = errors.New("invalid or expired session")
	ErrInvalidNonce        = errors.New("invalid or expired nonce")
	ErrInvalidSignature    = errors.New("invalid signature")
	ErrWalletAlreadyLinked = errors.New("wallet address already linked to another user")
)

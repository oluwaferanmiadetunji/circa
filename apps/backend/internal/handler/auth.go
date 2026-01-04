package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"circa/api"
	circaerrors "circa/internal/errors"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime/types"
	"github.com/rs/zerolog/log"
)

func (h *Handler) AuthSignup(ctx echo.Context) error {
	var req api.AuthSignupJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind request")
		return ctx.JSON(400, api.ErrorBadRequest{
			Code:    400,
			Message: "Invalid request body",
		})
	}

	if req.FullName == "" {
		return ctx.JSON(400, api.ErrorBadRequest{
			Code:    400,
			Message: "Full name is required",
		})
	}

	if req.Email == "" {
		return ctx.JSON(400, api.ErrorBadRequest{
			Code:    400,
			Message: "Email is required",
		})
	}

	result, err := h.authService.CreatePendingSignup(ctx.Request().Context(), req.FullName, string(req.Email), req.DisplayName)
	if err != nil {
		if errors.Is(err, circaerrors.ErrEmailAlreadyExists) {
			return ctx.JSON(400, api.ErrorBadRequest{
				Code:    400,
				Message: "Email already exists",
			})
		}
		log.Error().Err(err).Msg("Failed to create pending signup")
		return ctx.JSON(500, api.ErrorInternalServerError{
			Code:    500,
			Message: "Internal server error",
		})
	}

	log.Info().
		Str("pending_signup_id", result.PendingSignup.ID.String()).
		Str("magic_link_id", result.MagicLink.ID.String()).
		Msg("Pending signup and magic link created successfully")

	response := api.AuthSignupResponse{
		Message: "Signup request created successfully",
	}

	return ctx.JSON(200, response)
}

// AuthLogin handles POST /auth/login
func (h *Handler) AuthLogin(ctx echo.Context) error {
	var req api.AuthLoginJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind request")
		return ctx.JSON(400, api.ErrorBadRequest{
			Code:    400,
			Message: "Invalid request body",
		})
	}

	if req.Email == "" {
		return ctx.JSON(400, api.ErrorBadRequest{
			Code:    400,
			Message: "Email is required",
		})
	}

	result, err := h.authService.CreateLoginMagicLink(ctx.Request().Context(), string(req.Email))
	if err != nil {
		log.Error().Err(err).Msg("Failed to create login magic link")
		return ctx.JSON(500, api.ErrorInternalServerError{
			Code:    500,
			Message: "Internal server error",
		})
	}

	response := api.AuthLoginResponse{
		Message: result.Message,
	}

	return ctx.JSON(200, response)
}

// AuthNonce handles POST /auth/nonce
func (h *Handler) AuthNonce(ctx echo.Context) error {
	var req api.AuthNonceJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind request")
		return ctx.JSON(400, api.ErrorBadRequest{
			Code:    400,
			Message: "Invalid request body",
		})
	}

	// Validate address format
	if req.Address == "" {
		return ctx.JSON(400, api.ErrorBadRequest{
			Code:    400,
			Message: "Address is required",
		})
	}

	// Check for circa_signup cookie
	cookie, err := ctx.Cookie("circa_signup")
	if err != nil || cookie == nil || cookie.Value == "" {
		// Log all cookies to help debug
		allCookies := ctx.Cookies()
		log.Warn().
			Err(err).
			Interface("all_cookies", allCookies).
			Str("request_origin", ctx.Request().Header.Get("Origin")).
			Str("request_host", ctx.Request().Host).
			Msg("circa_signup cookie not found or empty when generating nonce")
		return ctx.JSON(401, api.ErrorUnauthorized{
			Code:    401,
			Message: "Valid signup session required. Please verify your email first.",
		})
	}

	sessionID := cookie.Value
	log.Info().
		Str("session_id", sessionID).
		Str("address", req.Address).
		Msg("Retrieved session ID from cookie for nonce generation")

	// Verify session exists in Redis
	sessionData, err := h.authService.GetSignupSession(ctx.Request().Context(), sessionID)
	if err != nil {
		log.Warn().
			Str("session_id", sessionID).
			Err(err).
			Msg("Session not found in Redis when generating nonce - cookie exists but session missing")
		return ctx.JSON(401, api.ErrorUnauthorized{
			Code:    401,
			Message: "Invalid or expired signup session. Please verify your email again.",
		})
	}
	log.Info().
		Str("session_id", sessionID).
		Interface("session_data", sessionData).
		Msg("Session verified in Redis for nonce generation")

	// Generate nonce tied to session and address
	var chainID *int64
	if req.ChainId != nil {
		chainIDVal := int64(*req.ChainId)
		chainID = &chainIDVal
	}
	nonceResult, err := h.authService.GenerateNonce(ctx.Request().Context(), sessionID, req.Address, chainID)
	if err != nil {
		if errors.Is(err, circaerrors.ErrInvalidSession) {
			return ctx.JSON(401, api.ErrorUnauthorized{
				Code:    401,
				Message: "Invalid or expired signup session. Please verify your email again.",
			})
		}
		log.Error().Err(err).Msg("Failed to generate nonce")
		return ctx.JSON(500, api.ErrorInternalServerError{
			Code:    500,
			Message: "Internal server error",
		})
	}

	// Convert to API response
	response := api.AuthNonceResponse{
		Nonce:           nonceResult.Nonce,
		ExpiresAt:       api.Timestamp(nonceResult.ExpiresAt),
		MessageTemplate: nonceResult.MessageTemplate,
	}

	return ctx.JSON(200, response)
}

// AuthVerify handles POST /auth/verify
func (h *Handler) AuthVerify(ctx echo.Context) error {
	var req api.AuthVerifyJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind request")
		return ctx.JSON(400, api.ErrorBadRequest{
			Code:    400,
			Message: "Invalid request body",
		})
	}

	if req.Token == "" {
		return ctx.JSON(400, api.ErrorBadRequest{
			Code:    400,
			Message: "Token is required",
		})
	}

	result, err := h.authService.VerifyToken(ctx.Request().Context(), req.Token)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err)
		if errors.Is(err, circaerrors.ErrInvalidToken) {
			return ctx.JSON(401, api.ErrorUnauthorized{
				Code:    401,
				Message: "Invalid or expired token",
			})
		}
		log.Error().Err(err).Msg("Failed to verify token")
		return ctx.JSON(500, api.ErrorInternalServerError{
			Code:    500,
			Message: "Internal server error",
		})
	}

	// Set appropriate cookie based on whether it's login or signup
	if result.NeedsWallet {
		// Signup session - set signup cookie (expires in 30 minutes)
		cookie := new(http.Cookie)
		cookie.Name = "circa_signup"
		cookie.Value = result.SessionID
		cookie.HttpOnly = true
		cookie.Secure = h.config.IsProduction // Only secure in production
		// Use SameSite=None in production (requires Secure), SameSite=Lax in development
		if h.config.IsProduction {
			cookie.SameSite = http.SameSiteNoneMode
		} else {
			cookie.SameSite = http.SameSiteLaxMode
		}
		cookie.Path = "/"
		cookie.MaxAge = 30 * 60 // 30 minutes in seconds
		// Don't set Domain for localhost - let browser handle it
		ctx.SetCookie(cookie)

		// Verify session was stored in Redis
		sessionData, err := h.authService.GetSignupSession(ctx.Request().Context(), result.SessionID)
		if err != nil {
			log.Error().
				Str("session_id", result.SessionID).
				Err(err).
				Msg("Failed to verify session was stored in Redis after setting cookie")
		} else {
			log.Info().
				Str("session_id", result.SessionID).
				Bool("secure", cookie.Secure).
				Interface("session_data", sessionData).
				Msg("Signup session cookie set and verified in Redis")
		}
	} else {
		// Login session - set main session cookie (expires in 7 days)
		cookie := new(http.Cookie)
		cookie.Name = "circa_session"
		cookie.Value = result.SessionID
		cookie.HttpOnly = true
		cookie.Secure = h.config.IsProduction
		cookie.SameSite = http.SameSiteLaxMode
		cookie.Path = "/"
		cookie.MaxAge = 7 * 24 * 60 * 60 // 7 days in seconds
		ctx.SetCookie(cookie)
	}

	response := api.AuthVerifyResponse{
		Email:       types.Email(result.Email),
		DisplayName: result.DisplayName,
		NeedsWallet: result.NeedsWallet,
	}

	return ctx.JSON(200, response)
}

// AuthSignupComplete handles POST /auth/signup/complete
func (h *Handler) AuthSignupComplete(ctx echo.Context) error {
	var req api.AuthSignupCompleteJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind request")
		return ctx.JSON(400, api.ErrorBadRequest{
			Code:    400,
			Message: "Invalid request body",
		})
	}

	// Validate required fields
	if req.Address == "" {
		return ctx.JSON(400, api.ErrorBadRequest{
			Code:    400,
			Message: "Address is required",
		})
	}
	if req.Signature == "" {
		return ctx.JSON(400, api.ErrorBadRequest{
			Code:    400,
			Message: "Signature is required",
		})
	}
	if req.Message == "" {
		return ctx.JSON(400, api.ErrorBadRequest{
			Code:    400,
			Message: "Message is required",
		})
	}

	// Check for circa_signup cookie
	cookie, err := ctx.Cookie("circa_signup")
	if err != nil || cookie == nil || cookie.Value == "" {
		log.Warn().
			Err(err).
			Str("address", string(req.Address)).
			Msg("circa_signup cookie not found or empty when completing signup")
		return ctx.JSON(401, api.ErrorUnauthorized{
			Code:    401,
			Message: "Valid signup session required. Please verify your email first.",
		})
	}

	sessionID := cookie.Value
	log.Info().
		Str("session_id", sessionID).
		Str("address", string(req.Address)).
		Msg("Retrieved session ID from cookie for signup completion")

	// Verify session exists in Redis before attempting to complete signup
	sessionData, err := h.authService.GetSignupSession(ctx.Request().Context(), sessionID)
	if err != nil {
		log.Warn().
			Str("session_id", sessionID).
			Str("address", string(req.Address)).
			Err(err).
			Msg("Session not found in Redis before completing signup - cookie exists but session missing")
		return ctx.JSON(401, api.ErrorUnauthorized{
			Code:    401,
			Message: "Invalid or expired session. Please verify your email again.",
		})
	}
	log.Info().
		Str("session_id", sessionID).
		Interface("session_data", sessionData).
		Msg("Session verified in Redis, proceeding with signup completion")

	// Complete signup
	result, err := h.authService.CompleteSignup(ctx.Request().Context(), sessionID, string(req.Address), req.Signature, req.Message)
	if err != nil {
		if errors.Is(err, circaerrors.ErrInvalidSession) {
			log.Warn().
				Str("session_id", sessionID).
				Str("address", string(req.Address)).
				Err(err).
				Msg("Invalid or expired signup session during completion")
			return ctx.JSON(401, api.ErrorUnauthorized{
				Code:    401,
				Message: "Invalid or expired session. Please verify your email again.",
			})
		}
		if errors.Is(err, circaerrors.ErrInvalidToken) || errors.Is(err, circaerrors.ErrInvalidNonce) {
			return ctx.JSON(401, api.ErrorUnauthorized{
				Code:    401,
				Message: "Invalid or expired nonce. Please connect your wallet again.",
			})
		}
		if errors.Is(err, circaerrors.ErrInvalidSignature) {
			return ctx.JSON(401, api.ErrorUnauthorized{
				Code:    401,
				Message: "Invalid signature. Please try signing again.",
			})
		}
		if errors.Is(err, circaerrors.ErrWalletAlreadyLinked) {
			return ctx.JSON(400, api.ErrorBadRequest{
				Code:    400,
				Message: err.Error(),
			})
		}
		if strings.Contains(err.Error(), "already") {
			return ctx.JSON(409, api.ErrorBadRequest{
				Code:    409,
				Message: err.Error(),
			})
		}
		log.Error().Err(err).Msg("Failed to complete signup")
		return ctx.JSON(500, api.ErrorInternalServerError{
			Code:    500,
			Message: "Internal server error",
		})
	}

	// Clear the signup cookie
	clearCookie := new(http.Cookie)
	clearCookie.Name = "circa_signup"
	clearCookie.Value = ""
	clearCookie.HttpOnly = true
	clearCookie.Secure = h.config.IsProduction
	clearCookie.SameSite = http.SameSiteLaxMode
	clearCookie.Path = "/"
	clearCookie.MaxAge = -1
	ctx.SetCookie(clearCookie)

	// Set main session cookie (expires in 7 days)
	mainCookie := new(http.Cookie)
	mainCookie.Name = "circa_session"
	mainCookie.Value = result.SessionID
	mainCookie.HttpOnly = true
	mainCookie.Secure = h.config.IsProduction
	mainCookie.SameSite = http.SameSiteLaxMode
	mainCookie.Path = "/"
	mainCookie.MaxAge = 7 * 24 * 60 * 60 // 7 days in seconds
	ctx.SetCookie(mainCookie)

	// Convert user to API response
	user := api.User{
		Id:          result.User.ID,
		Address:     api.Address(result.User.Address),
		CreatedAt:   api.Timestamp(result.User.CreatedAt.Time),
		DisplayName: result.User.DisplayName,
	}
	if result.User.UpdatedAt.Valid {
		updatedAt := api.Timestamp(result.User.UpdatedAt.Time)
		user.UpdatedAt = &updatedAt
	}

	response := api.AuthVerifyWalletResponse{
		User: user,
	}

	return ctx.JSON(200, response)
}

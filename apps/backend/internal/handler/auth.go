package handler

import (
	"errors"

	"circa/api"
	circaerrors "circa/internal/errors"

	"github.com/labstack/echo/v4"
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

	// Generate nonce
	nonceResult, err := h.authService.GenerateNonce(req.Address)
	if err != nil {
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

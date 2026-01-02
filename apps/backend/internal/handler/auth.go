package handler

import (
	"circa/api"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

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

package handler

import (
	"circa/api"
	"circa/internal/config"
	circaerrors "circa/internal/errors"
	"circa/internal/service/auth"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	authService auth.AuthService
	config      config.Config
}

// NewHandler creates a new handler instance
func NewHandler(authService auth.AuthService, cfg config.Config) *Handler {
	return &Handler{
		authService: authService,
		config:      cfg,
	}
}

// AuthLogout handles POST /auth/logout
func (h *Handler) AuthLogout(ctx echo.Context) error {
	// TODO: Implement logout
	return ctx.NoContent(204)
}

// GetMe handles GET /me
func (h *Handler) GetMe(ctx echo.Context) error {
	// Check for circa_session cookie
	cookie, err := ctx.Cookie("circa_session")
	if err != nil || cookie == nil || cookie.Value == "" {
		return ctx.JSON(401, api.ErrorUnauthorized{
			Code:    401,
			Message: "Unauthorized - no valid session",
		})
	}

	sessionID := cookie.Value

	// Get user from session
	result, err := h.authService.GetSessionUser(ctx.Request().Context(), sessionID)
	if err != nil {
		if errors.Is(err, circaerrors.ErrInvalidSession) {
			return ctx.JSON(401, api.ErrorUnauthorized{
				Code:    401,
				Message: "Unauthorized - invalid or expired session",
			})
		}
		log.Error().Err(err).Msg("Failed to get session user")
		return ctx.JSON(500, api.ErrorInternalServerError{
			Code:    500,
			Message: "Internal server error",
		})
	}

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

	return ctx.JSON(200, user)
}

// UpdateMe handles PATCH /me
func (h *Handler) UpdateMe(ctx echo.Context) error {
	// TODO: Implement update me
	return ctx.JSON(501, api.ErrorBadRequest{
		Code:    501,
		Message: "Not implemented",
	})
}

// ListGroups handles GET /groups
func (h *Handler) ListGroups(ctx echo.Context, params api.ListGroupsParams) error {
	// TODO: Implement list groups
	return ctx.JSON(501, api.ErrorBadRequest{
		Code:    501,
		Message: "Not implemented",
	})
}

// CreateGroup handles POST /groups
func (h *Handler) CreateGroup(ctx echo.Context) error {
	// TODO: Implement create group
	return ctx.JSON(501, api.ErrorBadRequest{
		Code:    501,
		Message: "Not implemented",
	})
}

// GetGroup handles GET /groups/{groupId}
func (h *Handler) GetGroup(ctx echo.Context, groupId api.UUID) error {
	// TODO: Implement get group
	return ctx.JSON(501, api.ErrorBadRequest{
		Code:    501,
		Message: "Not implemented",
	})
}

// UpdateGroup handles PATCH /groups/{groupId}
func (h *Handler) UpdateGroup(ctx echo.Context, groupId api.UUID) error {
	// TODO: Implement update group
	return ctx.JSON(501, api.ErrorBadRequest{
		Code:    501,
		Message: "Not implemented",
	})
}

// ListGroupMembers handles GET /groups/{groupId}/members
func (h *Handler) ListGroupMembers(ctx echo.Context, groupId api.UUID) error {
	// TODO: Implement list group members
	return ctx.JSON(501, api.ErrorBadRequest{
		Code:    501,
		Message: "Not implemented",
	})
}

// RemoveGroupMember handles DELETE /groups/{groupId}/members/{memberAddress}
func (h *Handler) RemoveGroupMember(ctx echo.Context, groupId api.UUID, memberAddress api.Address) error {
	// TODO: Implement remove group member
	return ctx.JSON(501, api.ErrorBadRequest{
		Code:    501,
		Message: "Not implemented",
	})
}

// LeaveGroup handles POST /groups/{groupId}/leave
func (h *Handler) LeaveGroup(ctx echo.Context, groupId api.UUID) error {
	// TODO: Implement leave group
	return ctx.JSON(501, api.ErrorBadRequest{
		Code:    501,
		Message: "Not implemented",
	})
}

// ListInvites handles GET /groups/{groupId}/invites
func (h *Handler) ListInvites(ctx echo.Context, groupId api.UUID) error {
	// TODO: Implement list invites
	return ctx.JSON(501, api.ErrorBadRequest{
		Code:    501,
		Message: "Not implemented",
	})
}

// CreateInvite handles POST /groups/{groupId}/invites
func (h *Handler) CreateInvite(ctx echo.Context, groupId api.UUID) error {
	// TODO: Implement create invite
	return ctx.JSON(501, api.ErrorBadRequest{
		Code:    501,
		Message: "Not implemented",
	})
}

// RevokeInvite handles DELETE /groups/{groupId}/invites/{inviteId}
func (h *Handler) RevokeInvite(ctx echo.Context, groupId api.UUID, inviteId api.UUID) error {
	// TODO: Implement revoke invite
	return ctx.JSON(501, api.ErrorBadRequest{
		Code:    501,
		Message: "Not implemented",
	})
}

// ListGroupRounds handles GET /groups/{groupId}/rounds
func (h *Handler) ListGroupRounds(ctx echo.Context, groupId api.UUID, params api.ListGroupRoundsParams) error {
	// TODO: Implement list group rounds
	return ctx.JSON(501, api.ErrorBadRequest{
		Code:    501,
		Message: "Not implemented",
	})
}

// CreateRound handles POST /groups/{groupId}/rounds
func (h *Handler) CreateRound(ctx echo.Context, groupId api.UUID) error {
	// TODO: Implement create round
	return ctx.JSON(501, api.ErrorBadRequest{
		Code:    501,
		Message: "Not implemented",
	})
}

// AcceptInvite handles POST /invites/accept
func (h *Handler) AcceptInvite(ctx echo.Context) error {
	// TODO: Implement accept invite
	return ctx.JSON(501, api.ErrorBadRequest{
		Code:    501,
		Message: "Not implemented",
	})
}

// PreviewInvite handles POST /invites/preview
func (h *Handler) PreviewInvite(ctx echo.Context) error {
	// TODO: Implement preview invite
	return ctx.JSON(501, api.ErrorBadRequest{
		Code:    501,
		Message: "Not implemented",
	})
}

// ListRounds handles GET /rounds
func (h *Handler) ListRounds(ctx echo.Context, params api.ListRoundsParams) error {
	// TODO: Implement list rounds
	return ctx.JSON(501, api.ErrorBadRequest{
		Code:    501,
		Message: "Not implemented",
	})
}

// GetRound handles GET /rounds/{roundId}
func (h *Handler) GetRound(ctx echo.Context, roundId api.UUID) error {
	// TODO: Implement get round
	return ctx.JSON(501, api.ErrorBadRequest{
		Code:    501,
		Message: "Not implemented",
	})
}

// GetRoundActivity handles GET /rounds/{roundId}/activity
func (h *Handler) GetRoundActivity(ctx echo.Context, roundId api.UUID, params api.GetRoundActivityParams) error {
	// TODO: Implement get round activity
	return ctx.JSON(501, api.ErrorBadRequest{
		Code:    501,
		Message: "Not implemented",
	})
}

// GetRoundPeriods handles GET /rounds/{roundId}/periods
func (h *Handler) GetRoundPeriods(ctx echo.Context, roundId api.UUID) error {
	// TODO: Implement get round periods
	return ctx.JSON(501, api.ErrorBadRequest{
		Code:    501,
		Message: "Not implemented",
	})
}

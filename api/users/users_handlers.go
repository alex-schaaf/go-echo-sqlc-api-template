package users

import (
	"app/db"
	"app/lib"
	"app/lib/auth"
	"app/lib/config"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UsersHandler struct {
	Queries *db.Queries
	Config  *config.Config
}

func NewUsersHandler(queries *db.Queries, config *config.Config) *UsersHandler {
	return &UsersHandler{Queries: queries, Config: config}
}

func bindAndValidatePasswordRequest(e echo.Context) (UpdateUserPasswordDto, error) {
	req := UpdateUserPasswordDto{}
	if err := e.Bind(&req); err != nil {
		return req, err
	}
	if req.OldPassword == "" || req.NewPassword == "" {
		return req, echo.NewHTTPError(http.StatusBadRequest, "Old password and new password are required")
	}
	if req.NewPassword == req.OldPassword {
		return req, echo.NewHTTPError(http.StatusBadRequest, "Old password and new password must be different")
	}
	if len(req.NewPassword) < 8 {
		return req, echo.NewHTTPError(http.StatusBadRequest, "Password must be at least 8 characters long")
	}

	return req, nil
}

func (h *UsersHandler) UpdateUserPasswordHandler(e echo.Context) error {
	userID := lib.GetUserID(e)

	req, err := bindAndValidatePasswordRequest(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user, err := h.Queries.GetUserById(e.Request().Context(), int64(userID))
	if err != nil {
		return e.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	if !auth.IsValidPassword(user.PasswordHash, req.OldPassword) {
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid old password"})
	}

	newPasswordHash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}

	if err := h.Queries.UpdateUserPassword(e.Request().Context(), db.UpdateUserPasswordParams{
		ID:           int64(userID),
		PasswordHash: newPasswordHash,
	}); err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update password"})
	}

	return e.NoContent(http.StatusNoContent)
}

func (h *UsersHandler) DeleteUserHandler(e echo.Context) error {
	userID := lib.GetUserID(e)

	if err := h.Queries.DeleteUser(e.Request().Context(), int64(userID)); err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete user"})
	}

	return e.NoContent(http.StatusNoContent)
}

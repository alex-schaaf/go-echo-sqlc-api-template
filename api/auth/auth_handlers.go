package auth

import (
	"app/db"
	"app/lib/auth"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	Queries *db.Queries
}

func NewAuthHandler(queries *db.Queries) *AuthHandler {
	return &AuthHandler{Queries: queries}
}

func (h *AuthHandler) LoginHandler(e echo.Context) error {
	// Parse email and password from the request body
	var loginData LoginDto

	if err := e.Bind(&loginData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	email := loginData.Email
	password := loginData.Password

	// Validate the email and password
	if email == "" || password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email and password are required")
	}

	// Get the user from the database
	user, err := h.Queries.GetUserByEmail(e.Request().Context(), email)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid email or password")
	}

	// Check if password is correct
	if !auth.IsValidPassword(user.PasswordHash, password) {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid email or password")
	}

	// Generate token and add as cookie
	cookie, err := auth.GenerateTokenCookie(user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create token cookie")
	}
	e.SetCookie(&cookie)
	userDto := UserDto{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
	return e.JSON(http.StatusOK, userDto)
}

func (h *AuthHandler) RegisterHandler(e echo.Context) error {
	// Parse email and password from the request body
	var registerData RegisterDto
	if err := e.Bind(&registerData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	if registerData.Email == "" || registerData.Password == "" || registerData.Username == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email, username, and password are required")
	}

	// Check if the email is already in use
	_, err := h.Queries.GetUserByEmail(e.Request().Context(), registerData.Email)
	if err == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Email is already in use")
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(registerData.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password")
	}

	// Create the user
	_, err = h.Queries.CreateUser(e.Request().Context(), db.CreateUserParams{
		Email:        registerData.Email,
		Username:     registerData.Username,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err) // "Failed to create user"
	}

	return e.NoContent(http.StatusCreated)
}

func (h *AuthHandler) LogoutHandler(e echo.Context) error {
	cookie := http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	e.SetCookie(&cookie)
	return e.NoContent(http.StatusNoContent)
}

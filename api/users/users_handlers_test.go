package users

import (
	"app/db"
	"app/lib/auth"
	"app/lib/test"
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

func TestUsersHandler_UpdateUserPasswordHandler(t *testing.T) {
	e, queries := test.SetupTest()
	ctx := context.Background()

	password := "oldPassword"
	passwordHash, _ := auth.HashPassword(password)
	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		Username:     "username",
		Email:        "user@example.com",
		PasswordHash: passwordHash,
	})
	if err != nil {
		t.Fatal(err)
	}

	usersHandler := NewUsersHandler(queries)
	updatePasswordJSON := `{"old_password":"oldPassword","new_password":"newPassword"}`
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(updatePasswordJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id")
	c.SetParamNames("user_id")
	c.SetParamValues(strconv.FormatInt(user.ID, 10))
	c.Set("user_id", strconv.FormatInt(user.ID, 10))

	if err := usersHandler.UpdateUserPasswordHandler(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d but got %d", http.StatusNoContent, rec.Code)
	}
}

func TestBindAndValidatePasswordRequest(t *testing.T) {
	tests := []struct {
		name          string
		requestBody   string
		expectedError bool
	}{
		{
			name:          "ValidPasswords",
			requestBody:   `{"old_password":"oldPassword","new_password":"newPassword"}`,
			expectedError: false,
		},
		{
			name:          "EmptyOldPassword",
			requestBody:   `{"old_password":"","new_password":"newpass"}`,
			expectedError: true,
		},
		{
			name:          "EmptyNewPassword",
			requestBody:   `{"old_password":"oldpass","new_password":""}`,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			_, err := bindAndValidatePasswordRequest(c)
			if (err != nil) != tt.expectedError {
				t.Fatalf("Expected error: %v, got: %v", tt.expectedError, err != nil)
			}
		})
	}
}

func TestGetUserID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "1")

	userID := GetUserID(c)
	if userID != 1 {
		t.Errorf("Expected user ID 1, got %d", userID)
	}
}

func TestUsersHandler_DeleteUserHandler(t *testing.T) {
	e, queries := test.SetupTest()
	ctx := context.Background()

	password := "oldPassword"
	passwordHash, _ := auth.HashPassword(password)
	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		Username:     "username",
		Email:        "user@example.com",
		PasswordHash: passwordHash,
	})
	if err != nil {
		t.Fatal(err)
	}

	usersHandler := NewUsersHandler(queries)
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id")
	c.SetParamNames("user_id")
	c.SetParamValues(strconv.FormatInt(user.ID, 10))
	c.Set("user_id", strconv.FormatInt(user.ID, 10))

	if err := usersHandler.DeleteUserHandler(c); err != nil {
		t.Fatal(err)
	}

	_, err = queries.GetUserById(ctx, user.ID)
	if err == nil {
		t.Error("Expected error but got nil")
	}

}
